package prompter

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Prompter defines an interface for user input interactions.
type Prompter interface {
	Select(message, defaultValue string, options []string) (int, error)
	MultiSelect(message, defaultValue string, options []string) ([]int, error)
	Input(message, defaultValue string) (string, error)
	Confirm(message string, defaultValue bool) (bool, error)
}

// New creates and returns a new Prompter instance.
func New() Prompter {
	return &surveyPrompter{}
}

type surveyPrompter struct {
}

// ask is a helper function to prompt the user and capture the response.
func (p *surveyPrompter) ask(q survey.Prompt, response interface{}) error {
	err := survey.AskOne(q, response)
	if err != nil {
		return fmt.Errorf("prompt error: %w", err)
	}
	return nil
}

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
func (p *surveyPrompter) Select(message, defaultValue string, options []string) (int, error) {
	var result int
	err := p.ask(&survey.Select{
		Message:  message,
		Default:  defaultValue,
		Options:  options,
		PageSize: 20,
	}, &result)
	return result, err
}

// MultiSelect prompts the user to select multiple options from a list.
//
// Parameters:
//   - message: The prompt message to display.
//   - defaultValue: The default selected values.
//   - options: The list of options to display.
//
// Returns:
//   - A slice of indices representing the selected options.
//   - An error, if any.
func (p *surveyPrompter) MultiSelect(message, defaultValue string, options []string) ([]int, error) {
	var result []int
	err := p.ask(&survey.MultiSelect{
		Message:  message,
		Default:  defaultValue,
		Options:  options,
		PageSize: 20,
	}, &result)
	return result, err
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
func (p *surveyPrompter) Input(message, defaultValue string) (string, error) {
	var result string
	err := p.ask(&survey.Input{
		Message: message,
		Default: defaultValue,
	}, &result)
	return result, err
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
func (p *surveyPrompter) Confirm(message string, defaultValue bool) (bool, error) {
	var result bool
	err := p.ask(&survey.Confirm{
		Message: message,
		Default: defaultValue,
	}, &result)
	return result, err
}
