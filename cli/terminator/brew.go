package terminator

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// IsUnderHomebrew checks if a given binary path is managed under the Homebrew path.
// This function is useful to verify if a binary is installed via Homebrew
// by comparing its location to the Homebrew binary directory.
func IsUnderHomebrew(path string) bool {
	if path == "" {
		return false
	}

	brewExe, err := exec.LookPath("brew")
	if err != nil {
		return false
	}

	brewPrefixBytes, err := exec.Command(brewExe, "--prefix").Output()
	if err != nil {
		return false
	}

	brewBinPrefix := filepath.Join(strings.TrimSpace(string(brewPrefixBytes)), "bin") + string(filepath.Separator)
	return strings.HasPrefix(path, brewBinPrefix)
}

// HasHomebrew checks if Homebrew is installed on the user's system.
// This function determines the presence of Homebrew by looking for the "brew"
// executable in the system's PATH. It is useful to ensure Homebrew dependencies
// can be managed before executing related commands.
func HasHomebrew() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}
