package prompt_test

import (
	"fmt"

	"github.com/raystack/salt/cli/prompt"
)

func ExampleNew() {
	p := prompt.New()

	// Single select
	idx, err := p.Select("Choose a color", "red", []string{"red", "green", "blue"})
	if err != nil {
		panic(err)
	}
	fmt.Println("selected index:", idx)

	// Text input
	name, err := p.Input("Enter your name", "")
	if err != nil {
		panic(err)
	}
	fmt.Println("name:", name)

	// Confirmation
	ok, err := p.Confirm("Continue?", true)
	if err != nil {
		panic(err)
	}
	fmt.Println("confirmed:", ok)
}
