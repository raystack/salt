package terminal

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"
)

// IsCI reports whether the process is running in a CI environment.
// Checks common environment variables used by GitHub Actions, Travis CI,
// CircleCI, Jenkins, TeamCity, and others.
func IsCI() bool {
	return os.Getenv("CI") != "" || // GitHub Actions, Travis CI, CircleCI, Cirrus CI, GitLab CI, AppVeyor, CodeShip, dsari
		os.Getenv("BUILD_NUMBER") != "" || // Jenkins, TeamCity
		os.Getenv("RUN_ID") != "" // TaskCluster, dsari
}

// OpenBrowser opens the default web browser at the specified URL.
// The goos parameter should be runtime.GOOS (e.g. "darwin", "windows", "linux").
//
// Returns an *exec.Cmd — call cmd.Run() or cmd.Start() to execute it.
// Returns an error if stdout is not a terminal.
func OpenBrowser(goos, url string) (*exec.Cmd, error) {
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return nil, fmt.Errorf("OpenBrowser requires a TTY on stdout")
	}
	return openBrowserCmd(goos, url), nil
}
