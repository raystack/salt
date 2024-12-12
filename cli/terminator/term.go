package terminator

import (
	"os"

	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
)

// IsTTY checks if the current output is a TTY (teletypewriter) or a Cygwin terminal.
// This function is useful for determining if the program is running in a terminal
// environment, which is important for features like colored output or interactive prompts.
func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}

// IsColorDisabled checks if color output is disabled based on the environment settings.
// This function uses the `termenv` library to determine if the NO_COLOR environment
// variable is set, which is a common way to disable colored output.
func IsColorDisabled() bool {
	return termenv.EnvNoColor()
}

// IsCI checks if the code is running in a Continuous Integration (CI) environment.
// This function checks for common environment variables used by popular CI systems
// like GitHub Actions, Travis CI, CircleCI, Jenkins, TeamCity, and others.
func IsCI() bool {
	return os.Getenv("CI") != "" || // GitHub Actions, Travis CI, CircleCI, Cirrus CI, GitLab CI, AppVeyor, CodeShip, dsari
		os.Getenv("BUILD_NUMBER") != "" || // Jenkins, TeamCity
		os.Getenv("RUN_ID") != "" // TaskCluster, dsari
}
