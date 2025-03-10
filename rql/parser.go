package rql

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

var validNumberOperations = []string{"eq", "neq", "gt", "lt", "gte", "lte"}
var validStringOperations = []string{"eq", "neq", "like", "in", "notin", "notlike", "empty", "notempty"}
var validBoolOperations = []string{"eq", "neq"}
var validDatetimeOperations = []string{"eq", "neq", "gt", "lt", "gte", "lte"}

const TAG = "rql"
const DATATYPE_NUMBER = "number"
const DATATYPE_DATETIME = "datetime"
const DATATYPE_STRING = "string"
const DATATYPE_BOOL = "bool"
const SORT_ORDER_ASC = "asc"
const SORT_ORDER_DESC = "desc"

var validSortOrder = []string{SORT_ORDER_ASC, SORT_ORDER_DESC}

type Query struct {
	Filters []Filter `json:"filters"`
	GroupBy []string `json:"group_by"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Search  string   `json:"search"`
	Sort    []Sort   `json:"sort"`
}

type Filter struct {
	Name     string `json:"name"`
	Operator string `json:"operator"`
	dataType string
	Value    any `json:"value"`
}

type Sort struct {
	Name  string `json:"name"`
	Order string `json:"order"`
}

func ValidateQuery(q *Query, checkStruct interface{}) error {
	val := reflect.ValueOf(checkStruct)

	// validate filters
	for _, filterItem := range q.Filters {
		//validate filter key name
		filterIdx := searchKeyInsideStruct(filterItem.Name, val)
		if filterIdx < 0 {
			return fmt.Errorf("'%s' is not a valid filter key", filterItem.Name)
		}
		structKeyTag := val.Type().Field(filterIdx).Tag.Get(TAG)

		// validate filter key data type
		allowedDataType := getDataTypeOfField(structKeyTag)
		filterItem.dataType = allowedDataType
		switch allowedDataType {
		case DATATYPE_NUMBER:
			err := validateNumberType(filterItem)
			if err != nil {
				return err
			}
		case DATATYPE_BOOL:
			err := validateBoolType(filterItem)
			if err != nil {
				return err
			}
		case DATATYPE_DATETIME:
			err := validateDatetimeType(filterItem)
			if err != nil {
				return err
			}
		case DATATYPE_STRING:
			err := validateStringType(filterItem)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("type '%s' is not recognized", allowedDataType)
		}

		if !isValidOperator(filterItem) {
			return fmt.Errorf("value '%s' for key '%s' is valid string", filterItem.Operator, filterItem.Name)
		}
	}

	err := validateGroupByKeys(q, val)
	if err != nil {
		return err
	}
	return validateSortKey(q, val)
}

func validateNumberType(filterItem Filter) error {
	// check if the type is any of Golang numeric types
	// if not, return error
	switch filterItem.Value.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, int, uint:
		return nil
	default:
		return fmt.Errorf("value %v for key '%s' is not int type", filterItem.Value, filterItem.Name)
	}
}

func validateDatetimeType(filterItem Filter) error {
	// cast the value to datetime
	// if failed, return error
	castedVal, ok := filterItem.Value.(string)
	if !ok {
		return fmt.Errorf("value %s for key '%s' is not a valid ISO datetime string", filterItem.Value, filterItem.Name)
	}
	_, err := time.Parse(time.RFC3339, castedVal)
	if err != nil {
		return fmt.Errorf("value %s for key '%s' is not a valid ISO datetime string", filterItem.Value, filterItem.Name)
	}
	return nil
}

func validateBoolType(filterItem Filter) error {
	// cast the value to bool
	// if failed, return error
	_, ok := filterItem.Value.(bool)
	if !ok {
		return fmt.Errorf("value %v for key '%s' is not bool type", filterItem.Value, filterItem.Name)
	}
	return nil
}

func validateStringType(filterItem Filter) error {
	// cast the value to string
	// if failed, return error
	_, ok := filterItem.Value.(string)
	if !ok {
		return fmt.Errorf("value %s for key '%s' is valid string type", filterItem.Value, filterItem.Name)
	}
	return nil
}

func searchKeyInsideStruct(keyName string, val reflect.Value) int {
	normalizedKey := strings.ToLower(keyName)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		// Check field name
		if strings.ToLower(field.Name) == normalizedKey {
			return i
		}

		// Check rql tag
		if tag, ok := field.Tag.Lookup("rql"); ok {
			// Parse the tag string
			tagParts := strings.Split(tag, ",")
			for _, part := range tagParts {
				if strings.HasPrefix(part, "name=") {
					tagName := strings.TrimPrefix(part, "name=")
					if strings.ToLower(tagName) == normalizedKey {
						return i
					}
				}
			}
		}
	}

	return -1
}

// parse the tag schema which is of the format
// type=int,min=10,max=200
// to extract type else fallback to string
func getDataTypeOfField(tagString string) string {
	res := DATATYPE_STRING
	splitted := strings.Split(tagString, ",")
	for _, item := range splitted {
		kvSplitted := strings.Split(item, "=")
		if len(kvSplitted) == 2 {
			if kvSplitted[0] == "type" {
				return kvSplitted[1]
			}
		}
	}
	//fallback to string if type not found in tag value
	return res
}

func GetDataTypeOfField(fieldName string, checkStruct interface{}) (string, error) {
	val := reflect.ValueOf(checkStruct)
	filterIdx := searchKeyInsideStruct(fieldName, val)
	if filterIdx < 0 {
		return "", fmt.Errorf("'%s' is not a valid field", fieldName)
	}
	structKeyTag := val.Type().Field(filterIdx).Tag.Get(TAG)
	dataType := getDataTypeOfField(structKeyTag)
	if !slices.Contains([]string{DATATYPE_STRING, DATATYPE_BOOL, DATATYPE_NUMBER, DATATYPE_DATETIME}, dataType) {
		return "", fmt.Errorf("invalid datatype '%s' is for field %s", dataType, fieldName)
	}
	return dataType, nil
}

func isValidOperator(filterItem Filter) bool {
	switch filterItem.dataType {
	case DATATYPE_NUMBER:
		return slices.Contains(validNumberOperations, filterItem.Operator)
	case DATATYPE_DATETIME:
		return slices.Contains(validDatetimeOperations, filterItem.Operator)
	case DATATYPE_STRING:
		return slices.Contains(validStringOperations, filterItem.Operator)
	case DATATYPE_BOOL:
		return slices.Contains(validBoolOperations, filterItem.Operator)
	default:
		return false
	}
}

func validateSortKey(q *Query, val reflect.Value) error {
	for _, item := range q.Sort {
		filterIdx := searchKeyInsideStruct(item.Name, val)
		if filterIdx < 0 {
			return fmt.Errorf("'%s' is not a valid sort key", item.Name)
		}
		if !slices.Contains(validSortOrder, item.Order) {
			return fmt.Errorf("'%s' is not a valid sort key", item.Name)
		}
	}
	return nil
}

func validateGroupByKeys(q *Query, val reflect.Value) error {
	for _, item := range q.GroupBy {
		filterIdx := searchKeyInsideStruct(item, val)
		if filterIdx < 0 {
			return fmt.Errorf("'%s' is not a valid sort key", item)
		}
	}
	return nil
}
