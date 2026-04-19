package terminal

import (
	"os"
	"os/exec"
	"strings"
)

// openBrowserCmd returns an exec.Cmd configured to open the URL.
func openBrowserCmd(goos, url string) *exec.Cmd {
	exe := "open"
	var args []string

	switch goos {
	case "darwin":
		args = append(args, url)
	case "windows":
		exe, _ = exec.LookPath("cmd")
		replacer := strings.NewReplacer("&", "^&")
		args = append(args, "/c", "start", replacer.Replace(url))
	default:
		exe = linuxExe()
		args = append(args, url)
	}

	cmd := exec.Command(exe, args...)
	cmd.Stderr = os.Stderr
	return cmd
}

// linuxExe determines the appropriate command to open a web browser on Linux.
func linuxExe() string {
	exe := "xdg-open"
	if _, err := exec.LookPath(exe); err != nil {
		if _, err := exec.LookPath("wslview"); err == nil {
			exe = "wslview"
		}
	}
	return exe
}
