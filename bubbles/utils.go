package bubbles

import (
	"os"

	"github.com/muesli/termenv"
	"golang.org/x/term"
)

func ClearScreen() {
	termenv.ClearScreen()
}

func TermSize() (w, h int) {
	w, h, _ = term.GetSize(int(os.Stdin.Fd()))
	return
}

func TermWidth() (w int) {
	w, _, _ = term.GetSize(int(os.Stdin.Fd()))
	return
}

func TermHeight() (h int) {
	_, h, _ = term.GetSize(int(os.Stdin.Fd()))
	return
}
