package term

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// IsUnderHomebrew checks whether the given binary is under the homebrew path.
func IsUnderHomebrew(binpath string) bool {
	if binpath == "" {
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
	return strings.HasPrefix(binpath, brewBinPrefix)
}

// HasHomebrew check whether the user has installed brew
func HasHomebrew() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}
