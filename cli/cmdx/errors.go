package cmdx

import "strings"

// IsCmdErr checks if the given error is related to a Cobra command error.
//
// This is useful for distinguishing between user errors (e.g., incorrect commands or flags)
// and program errors, allowing the application to display appropriate messages.
//
// Parameters:
//   - err: The error to check.
//
// Returns:
//   - true if the error message contains any known Cobra command error keywords.
//   - false otherwise.
func IsCmdErr(err error) bool {
	if err == nil {
		return false
	}

	// Known Cobra command error keywords
	cmdErrorKeywords := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
	}

	errMessage := err.Error()
	for _, keyword := range cmdErrorKeywords {
		if strings.Contains(errMessage, keyword) {
			return true
		}
	}
	return false
}
