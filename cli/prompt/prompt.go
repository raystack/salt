package prompt

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

// Prompter defines an interface for user input interactions.
type Prompter interface {
	Select(message, defaultValue string, options []string) (int, error)
	MultiSelect(message, defaultValue string, options []string) ([]int, error)
	Input(message, defaultValue string) (string, error)
	Password(message string) (string, error)
	Confirm(message string, defaultValue bool) (bool, error)
}

// New creates and returns a new Prompter instance.
func New() Prompter {
	return &huhPrompter{}
}

type huhPrompter struct{}

// Select prompts the user to select one option from a list.
//
// Parameters:
//   - message: The prompt message to display.
//   - defaultValue: The default selected value.
//   - options: The list of options to display.
//
// Returns:
//   - The index of the selected option.
//   - An error, if any.
func (p *huhPrompter) Select(message, defaultValue string, options []string) (int, error) {
	huhOptions := make([]huh.Option[int], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt, i)
	}

	var result int
	// Find default index
	for i, opt := range options {
		if opt == defaultValue {
			result = i
			break
		}
	}

	err := huh.NewSelect[int]().
		Title(message).
		Options(huhOptions...).
		Value(&result).
		Run()
	if err != nil {
		return 0, fmt.Errorf("prompt error: %w", err)
	}
	return result, nil
}

// MultiSelect prompts the user to select multiple options from a list.
//
// Parameters:
//   - message: The prompt message to display.
//   - defaultValue: The default selected values (unused, kept for interface compat).
//   - options: The list of options to display.
//
// Returns:
//   - A slice of indices representing the selected options.
//   - An error, if any.
func (p *huhPrompter) MultiSelect(message, _ string, options []string) ([]int, error) {
	huhOptions := make([]huh.Option[int], len(options))
	for i, opt := range options {
		huhOptions[i] = huh.NewOption(opt, i)
	}

	var result []int
	err := huh.NewMultiSelect[int]().
		Title(message).
		Options(huhOptions...).
		Value(&result).
		Run()
	if err != nil {
		return nil, fmt.Errorf("prompt error: %w", err)
	}
	return result, nil
}

// Input prompts the user for a text input.
//
// Parameters:
//   - message: The prompt message to display.
//   - defaultValue: The default input value.
//
// Returns:
//   - The user's input as a string.
//   - An error, if any.
func (p *huhPrompter) Input(message, defaultValue string) (string, error) {
	var result string
	err := huh.NewInput().
		Title(message).
		Value(&result).
		Placeholder(defaultValue).
		Run()
	if err != nil {
		return "", fmt.Errorf("prompt error: %w", err)
	}
	if result == "" {
		result = defaultValue
	}
	return result, nil
}

// Password prompts the user for a secret input. The input is masked.
//
// Parameters:
//   - message: The prompt message to display.
//
// Returns:
//   - The user's input as a string.
//   - An error, if any.
func (p *huhPrompter) Password(message string) (string, error) {
	var result string
	err := huh.NewInput().
		Title(message).
		Value(&result).
		EchoMode(huh.EchoModePassword).
		Run()
	if err != nil {
		return "", fmt.Errorf("prompt error: %w", err)
	}
	return result, nil
}

// Confirm prompts the user for a yes/no confirmation.
//
// Parameters:
//   - message: The prompt message to display.
//   - defaultValue: The default confirmation value.
//
// Returns:
//   - A boolean indicating the user's choice.
//   - An error, if any.
func (p *huhPrompter) Confirm(message string, defaultValue bool) (bool, error) {
	result := defaultValue
	err := huh.NewConfirm().
		Title(message).
		Value(&result).
		Run()
	if err != nil {
		return false, fmt.Errorf("prompt error: %w", err)
	}
	return result, nil
}
