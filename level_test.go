package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseLevel(t *testing.T) {
	cases := []struct {
		String string
		Level  Level
		Num    int
	}{
		{"debug", LevelDebug, 0},
		{"info", LevelInfo, 1},
		{"warn", LevelWarn, 2},
		{"error", LevelError, 3},
		{"fatal", LevelFatal, 4},
	}

	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			l, err := ParseLevel(c.String)
			require.NoError(t, err)
			require.Equal(t, c.Level, l)
		})
	}

	t.Run("invalid", func(t *testing.T) {
		l, err := ParseLevel("invalid level here")
		require.Equal(t, ErrInvalidLevel, err)
		require.Equal(t, LevelInvalid, l)
	})
}

func TestMustParseLevel(t *testing.T) {
	cases := []struct {
		String    string
		WantPanic bool
	}{
		{"invalid", true},
		{"debug", false},
		{"info", false},
		{"warn", false},
		{"error", false},
		{"fatal", false},
	}

	for _, c := range cases {
		t.Run(c.String, func(t *testing.T) {
			switch c.WantPanic {
			case true:
				require.Panics(t, func() { MustParseLevel(c.String) })
			case false:
				require.NotPanics(t, func() { MustParseLevel(c.String) })
			}
		})
	}
}

func TestStringLevel(t *testing.T) {
	cases := []struct {
		Name  string
		Level Level
		Want  string
	}{
		{"invalid", LevelInvalid, "invalid"},
		{"debug", LevelDebug, "dbg"},
		{"info", LevelInfo, "inf"},
		{"warn", LevelWarn, "wrn"},
		{"error", LevelError, "err"},
		{"fatal", LevelFatal, "ftl"},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require.Equal(t, c.Level.String(), c.Want)
		})
	}
}

func TestEvalLevel(t *testing.T) {
	cases := []struct {
		Name     string
		Got      Level
		Want     Level
		Expected bool
	}{
		{Name: "Invalid against Debug", Got: LevelInvalid, Want: LevelDebug, Expected: false},
		{Name: "Debug against Info", Got: LevelDebug, Want: LevelInfo, Expected: false},
		{Name: "Info against Warn", Got: LevelInfo, Want: LevelWarn, Expected: false},
		{Name: "Warn against Error", Got: LevelWarn, Want: LevelError, Expected: false},
		{Name: "Error against Fatal", Got: LevelDebug, Want: LevelFatal, Expected: false},
		{Name: "Fatal should always pass", Got: LevelFatal, Want: LevelInvalid, Expected: true},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			switch c.Expected {
			case true:
				require.True(t, evalLevel(c.Got, c.Want))
			case false:
				require.False(t, evalLevel(c.Got, c.Want))
			}
		})
	}
}
