package terminator

import (
	"os"
	"os/exec"
	"strings"
)

// OpenBrowser opens the default web browser at the specified URL.
//
// Parameters:
//   - goos: The operating system name (e.g., "darwin", "windows", or "linux").
//   - url: The URL to open in the web browser.
//
// Returns:
//   - An *exec.Cmd configured to open the URL. Note that you must call `cmd.Run()`
//     or `cmd.Start()` on the returned command to execute it.
//
// Panics:
//   - This function will panic if called without a TTY (e.g., not running in a terminal).
func OpenBrowser(goos, url string) *exec.Cmd {
	if !IsTTY() {
		panic("OpenBrowser called without a TTY")
	}

	exe := "open"
	var args []string

	switch goos {
	case "darwin":
		// macOS: Use the "open" command to open the URL.
		args = append(args, url)
	case "windows":
		// Windows: Use "cmd /c start" to open the URL.
		exe, _ = exec.LookPath("cmd")
		replacer := strings.NewReplacer("&", "^&")
		args = append(args, "/c", "start", replacer.Replace(url))
	default:
		// Linux: Use "xdg-open" or fallback to "wslview" for WSL environments.
		exe = linuxExe()
		args = append(args, url)
	}

	// Create the command to open the browser and set stderr for error reporting.
	cmd := exec.Command(exe, args...)
	cmd.Stderr = os.Stderr
	return cmd
}

// linuxExe determines the appropriate command to open a web browser on Linux.
func linuxExe() string {
	exe := "xdg-open"

	_, err := exec.LookPath(exe)
	if err != nil {
		_, err := exec.LookPath("wslview")
		if err == nil {
			exe = "wslview"
		}
	}

	return exe
}
