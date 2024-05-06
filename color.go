package log

import (
	"fmt"
)

var noColor bool = false

type color string

const (
	colorWhite  color = "\033[37m"
	colorCyan   color = "\033[36m"
	colorYellow color = "\033[33m"
	colorRed    color = "\033[31m"
	colorReset  color = "\033[0m"
)

func colorString(s string, c color) string {
	if noColor {
		return s
	}

	// TODO: Find way to check if we can color this

	if c == colorReset {
		return s
	}

	return fmt.Sprintf("%s%s%s", c, s, colorReset)
}
