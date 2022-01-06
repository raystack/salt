package term

import (
	"os"

	"github.com/mattn/go-isatty"
)

func IsTTY() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
