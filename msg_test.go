package log

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMsgString(t *testing.T) {
	t.Run("msg with nil args", func(t *testing.T) {
		msg := &msg{
			Timestamp: time.Time{},
			Level:     LevelInfo,
			Msg:       "Message",
			Args:      nil,
		}

		expected := fmt.Sprintf(
			`timestamp=%s level=%s msg="%s"`,
			msg.Timestamp.Format(time.RFC3339),
			formatLevel(msg.Level),
			msg.Msg,
		)

		require.Equal(t, msg.String(), expected)
	})

	t.Run("msg with len 0 args", func(t *testing.T) {
		msg := &msg{
			Timestamp: time.Time{},
			Level:     LevelInfo,
			Msg:       "Message",
			Args:      make(map[string]interface{}),
		}

		expected := fmt.Sprintf(
			`timestamp=%s level=%s msg="%s"`,
			formatTimestampRFC3339(msg.Timestamp),
			formatLevel(msg.Level),
			msg.Msg,
		)

		require.Equal(t, msg.String(), expected)
	})

	t.Run("msg with args", func(t *testing.T) {
		args := make(map[string]interface{})
		args["key"] = 42
		args["key2"] = "meaning"

		msg := &msg{
			Timestamp: time.Time{},
			Level:     LevelInfo,
			Msg:       "Message",
			Args:      args,
		}

		expected := fmt.Sprintf(
			`timestamp=%s level=%s msg="%s" %v`,
			formatTimestampRFC3339(msg.Timestamp),
			formatLevel(msg.Level),
			msg.Msg,
			formatArgs(msg.Args),
		)

		require.Equal(t, msg.String(), expected)
	})
}

func TestMsgMarshal(t *testing.T) {
	t.Run("msg with no args", func(t *testing.T) {
		msg := &msg{
			Timestamp: time.Time{},
			Level:     LevelInfo,
			Msg:       "test",
		}

		b, err := msg.Marshal()
		require.NoError(t, err)

		expected := `{"timestamp":"0001-01-01T00:00:00Z","level":"inf","msg":"test"}`
		require.Equal(t, string(b), expected)
	})

	// TODO: this errors because key/value pairs get changed
	t.Run("msg with args", func(t *testing.T) {
		msg := &msg{
			Timestamp: time.Time{},
			Level:     LevelInfo,
			Msg:       "test",
			Args:      map[string]interface{}{"key": "value", "key2": 42},
		}

		b, err := msg.Marshal()
		require.NoError(t, err)

		expected := `{"timestamp":"0001-01-01T00:00:00Z","level":"inf","msg":"test","key":"value","key2":42}`
		require.Equal(t, string(b), expected)
	})
}

func TestFormatTimestampRFC3339(t *testing.T) {
	t.Run("time now pass", func(t *testing.T) {
		ts := time.Now()

		require.Equal(t, formatTimestampRFC3339(ts), ts.Format(time.RFC3339))
	})

	t.Run("zero time pass", func(t *testing.T) {
		ts := time.Time{}

		require.Equal(t, formatTimestampRFC3339(ts), ts.Format(time.RFC3339))
	})
}

func TestFormatLevel(t *testing.T) {
	cases := []struct {
		Name     string
		L        Level
		Expected string
	}{
		{Name: "invalid", L: LevelInvalid, Expected: colorString("INVALID", colorWhite)},
		{Name: "debug", L: LevelDebug, Expected: colorString("DBG", colorWhite)},
		{Name: "info", L: LevelInfo, Expected: colorString("INF", colorCyan)},
		{Name: "warn", L: LevelWarn, Expected: colorString("WRN", colorYellow)},
		{Name: "error", L: LevelError, Expected: colorString("ERR", colorRed)},
		{Name: "fatal", L: LevelFatal, Expected: colorString("FTL", colorRed)},
		{Name: "default", L: 999, Expected: colorString("INVALID", colorWhite)},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require.Equal(t, formatLevel(c.L), c.Expected)
		})
	}

	t.Run("default no color", func(t *testing.T) {
		noColor = true
		defer func() { noColor = false }()

		require.Equal(t, formatLevel(LevelInfo), strings.ToUpper(LevelInfo.String()))
	})
}

func TestFormatArgs(t *testing.T) {
	sampleArgs := make(map[string]interface{})
	sampleArgs["key"] = 42
	sampleArgs["key2"] = "meaning"

	sampleExpected := `key=42 key2="meaning"`

	cases := []struct {
		Name     string
		Args     map[string]interface{}
		Expected string
	}{
		{Name: "nil args", Args: nil, Expected: ""},
		{Name: "args len 0", Args: make(map[string]interface{}), Expected: ""},
		{Name: "with args", Args: sampleArgs, Expected: sampleExpected},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require.Equal(t, formatArgs(c.Args), c.Expected)
		})
	}
}

func TestCheckStringType(t *testing.T) {
	sampleMap := make(map[string]string)
	sampleMap["key"] = "value"

	cases := []struct {
		Name     string
		Arg      interface{}
		Expected interface{}
	}{
		{Name: "nil arg", Arg: nil, Expected: ""},
		{Name: "integer arg", Arg: 42, Expected: 42},
		{Name: "string arg", Arg: "42", Expected: `"42"`},
		{Name: "map arg", Arg: sampleMap, Expected: sampleMap},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require.Equal(t, checkStringType(c.Arg), c.Expected)
		})
	}
}
