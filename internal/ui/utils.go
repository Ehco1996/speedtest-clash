package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// General stuff for styling the view
var (
	proxyNodeListStyle = lipgloss.NewStyle().Margin(0, 0)
	term               = termenv.EnvColorProfile()
	keyword            = makeFgStyle("211")
	subtle             = makeFgStyle("241")
	dot                = colorFg(" • ", "236")
)

// Color a string's foreground with the given value.
func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}
