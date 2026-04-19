package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// Exporter controls structured JSON output for a command.
// When non-nil, the command should call Write instead of rendering
// human-readable output.
//
// A nil Exporter means the user did not request --json.
type Exporter interface {
	// Fields returns the field names requested via --json.
	Fields() []string
	// Write encodes data as JSON and writes it to the IOStreams.
	// On a TTY it writes indented, colorized JSON; when piped it
	// writes compact JSON.
	Write(ios *IOStreams, data any) error
}

// Exportable may be implemented by structs to control which fields
// are exported and how. When a struct implements this interface,
// the JSON exporter calls ExportData instead of using reflection.
//
//	func (r *Resource) ExportData(fields []string) map[string]any {
//	    return cli.StructExportData(r, fields)
//	}
type Exportable interface {
	ExportData(fields []string) map[string]any
}

// AddJSONFlags adds a --json flag to cmd that accepts a comma-separated
// list of field names. When --json is used, exportTarget is set to a
// non-nil Exporter in a PreRunE hook. The command should check for a
// non-nil Exporter and call Write instead of rendering a table.
//
//	var exporter cli.Exporter
//	cli.AddJSONFlags(cmd, &exporter, []string{"id", "name", "status"})
//
//	// In RunE:
//	if exporter != nil {
//	    return exporter.Write(cli.IO(cmd), results)
//	}
//	cli.Output(cmd).Table(rows)
func AddJSONFlags(cmd *cobra.Command, exportTarget *Exporter, fields []string) {
	cmd.Flags().StringSlice("json", nil, "Output JSON with the specified `fields`")

	// Shell completion for field names.
	_ = cmd.RegisterFlagCompletionFunc("json", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var results []string
		var prefix string
		if idx := strings.LastIndexByte(toComplete, ','); idx >= 0 {
			prefix = toComplete[:idx+1]
			toComplete = toComplete[idx+1:]
		}
		toComplete = strings.ToLower(toComplete)
		for _, f := range fields {
			if strings.HasPrefix(strings.ToLower(f), toComplete) {
				results = append(results, prefix+f)
			}
		}
		sort.Strings(results)
		return results, cobra.ShellCompDirectiveNoSpace
	})

	// Validate field names and set the exporter.
	oldPreRun := cmd.PreRunE
	cmd.PreRunE = func(c *cobra.Command, args []string) error {
		if oldPreRun != nil {
			if err := oldPreRun(c, args); err != nil {
				return err
			}
		}

		jsonFlag := c.Flags().Lookup("json")
		if jsonFlag == nil || !jsonFlag.Changed {
			*exportTarget = nil
			return nil
		}

		requested, _ := c.Flags().GetStringSlice("json")
		allowed := make(map[string]bool, len(fields))
		for _, f := range fields {
			allowed[f] = true
		}
		for _, f := range requested {
			if !allowed[f] {
				sorted := make([]string, len(fields))
				copy(sorted, fields)
				sort.Strings(sorted)
				return fmt.Errorf("unknown JSON field: %q\nAvailable fields:\n  %s", f, strings.Join(sorted, "\n  "))
			}
		}

		*exportTarget = &jsonExporter{fields: requested}
		return nil
	}

	// When --json is passed without arguments, list available fields.
	parentFlagErr := cmd.FlagErrorFunc()
	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		if c == cmd && err.Error() == "flag needs an argument: --json" {
			sorted := make([]string, len(fields))
			copy(sorted, fields)
			sort.Strings(sorted)
			return fmt.Errorf("specify one or more comma-separated fields for --json:\n  %s", strings.Join(sorted, "\n  "))
		}
		if parentFlagErr != nil {
			return parentFlagErr(c, err)
		}
		return err
	})

	// Annotate for help display.
	if len(fields) > 0 {
		if cmd.Annotations == nil {
			cmd.Annotations = map[string]string{}
		}
		cmd.Annotations["help:json-fields"] = strings.Join(fields, ",")
	}
}

// StructExportData extracts the requested fields from a struct using
// case-insensitive field name matching. Use this as a default
// implementation for [Exportable.ExportData]:
//
//	func (r *Resource) ExportData(fields []string) map[string]any {
//	    return cli.StructExportData(r, fields)
//	}
func StructExportData(s any, fields []string) map[string]any {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	data := make(map[string]any, len(fields))
	for _, f := range fields {
		sf := fieldByTag(v, f)
		if !sf.IsValid() {
			sf = fieldByName(v, f)
		}
		if sf.IsValid() && sf.CanInterface() {
			data[f] = sf.Interface()
		}
	}
	return data
}

// jsonExporter is the default Exporter implementation.
type jsonExporter struct {
	fields []string
}

func (e *jsonExporter) Fields() []string {
	return e.fields
}

func (e *jsonExporter) Write(ios *IOStreams, data any) error {
	extracted := e.extractData(reflect.ValueOf(data))

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(extracted); err != nil {
		return err
	}

	w := ios.Out
	if ios.IsStdoutTTY() {
		// Re-encode with indentation for readability.
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, buf.Bytes(), "", "  "); err != nil {
			// Fallback to compact.
			_, err = io.Copy(w, buf)
			return err
		}
		pretty.WriteByte('\n')
		_, err := io.Copy(w, &pretty)
		return err
	}

	_, err := io.Copy(w, buf)
	return err
}

func (e *jsonExporter) extractData(v reflect.Value) any {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if !v.IsNil() {
			return e.extractData(v.Elem())
		}
		return nil
	case reflect.Slice:
		a := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			a[i] = e.extractData(v.Index(i))
		}
		return a
	case reflect.Struct:
		if v.CanAddr() {
			if ex, ok := v.Addr().Interface().(Exportable); ok {
				return ex.ExportData(e.fields)
			}
		}
		if ex, ok := v.Interface().(Exportable); ok {
			return ex.ExportData(e.fields)
		}
		return StructExportData(v.Interface(), e.fields)
	}
	return v.Interface()
}

// fieldByTag finds a struct field whose `json` tag matches the given name.
func fieldByTag(v reflect.Value, name string) reflect.Value {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if idx := strings.IndexByte(tag, ','); idx >= 0 {
			tag = tag[:idx]
		}
		if strings.EqualFold(tag, name) {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

// fieldByName finds a struct field by case-insensitive name match.
func fieldByName(v reflect.Value, name string) reflect.Value {
	return v.FieldByNameFunc(func(s string) bool {
		return strings.EqualFold(name, s)
	})
}
