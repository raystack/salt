package cfg

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	ODPF_CONFIG_DIR = "ODPF_CONFIG_DIR"
	XDG_CONFIG_HOME = "XDG_CONFIG_HOME"
	APP_DATA        = "AppData"
	LOCAL_APP_DATA  = "LocalAppData"
)

func ConfigFile(app string) string {
	file := app + ".yml"
	return filepath.Join(configDir("odpf"), file)
}

func configDir(root string) string {
	var path string
	if a := os.Getenv(ODPF_CONFIG_DIR); a != "" {
		path = a
	} else if b := os.Getenv(XDG_CONFIG_HOME); b != "" {
		path = filepath.Join(b, root)
	} else if c := os.Getenv(APP_DATA); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, root)
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", root)
	}

	if !dirExists(path) {
		_ = os.MkdirAll(filepath.Dir(path), 0755)
	}

	return path
}

func dirExists(path string) bool {
	f, err := os.Stat(path)
	return err == nil && f.IsDir()
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
