package log

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestColorString(t *testing.T) {
	cases := []struct {
		Color string
		S     string
		C     color
		Want  string
	}{
		{
			Color: "white",
			S:     "test",
			C:     colorWhite,
			Want:  fmt.Sprintf("%s%s%s", colorWhite, "test", colorReset),
		},
		{
			Color: "cyan",
			S:     "test",
			C:     colorCyan,
			Want:  fmt.Sprintf("%s%s%s", colorCyan, "test", colorReset),
		},
		{
			Color: "yellow",
			S:     "test",
			C:     colorYellow,
			Want:  fmt.Sprintf("%s%s%s", colorYellow, "test", colorReset),
		},
		{
			Color: "red",
			S:     "test",
			C:     colorRed,
			Want:  fmt.Sprintf("%s%s%s", colorRed, "test", colorReset),
		},
		{Color: "reset", S: "test", C: colorReset, Want: "test"},
	}

	for _, c := range cases {
		t.Run(c.Color, func(t *testing.T) {
			require.Equal(t, colorString(c.S, c.C), c.Want)
		})
	}

	// TODO: add test if we can print colored output

	t.Run("no color", func(t *testing.T) {
		noColor = true
		defer func() { noColor = false }()

		require.Equal(t, colorString("test", colorCyan), "test")
	})
}
