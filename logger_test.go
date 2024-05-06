package log

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLoggerSetOut(t *testing.T) {
	l := NewLogger()

	t.Run("out is nil", func(t *testing.T) {
		l.SetOut(nil)

		require.Equal(t, defaultOut, l.out)
	})

	t.Run("out is not nil", func(t *testing.T) {
		l.SetOut(os.Stdout)

		require.Equal(t, os.Stdout, l.out)
	})
}

func TestLoggerSetLevel(t *testing.T) {
	l := NewLogger()

	t.Run("level is invalid", func(t *testing.T) {
		l.SetLevel(LevelInvalid)

		require.Equal(t, defaultLevel, l.level)
	})

	t.Run("level valid", func(t *testing.T) {
		l.SetLevel(LevelWarn)

		require.Equal(t, LevelWarn, l.level)
		require.NotEqual(t, l.level, defaultLevel)
	})

	t.Run("number level invalid", func(t *testing.T) {
		l.SetLevel(99)

		require.Equal(t, LevelInfo, l.GetLevel())
	})
}

func TestLoggerGetLevel(t *testing.T) {
	l := NewLogger()

	t.Run("default level", func(t *testing.T) {
		require.Equal(t, l.GetLevel(), defaultLevel)
	})

	t.Run("non default level", func(t *testing.T) {
		l.SetLevel(LevelWarn)

		require.NotEqual(t, l.GetLevel(), defaultLevel)
		require.Equal(t, l.GetLevel(), LevelWarn)
	})
}

func TestLoggerSetHandler(t *testing.T) {
	l := NewLogger()

	t.Run("default handler", func(t *testing.T) {
		require.Equal(t, l.handler, defaultHandler)
	})

	t.Run("non default handler", func(t *testing.T) {
		l.SetHandler(JSONHandler)
		require.Equal(t, l.handler, JSONHandler)
	})

	t.Run("invalid handler, use default", func(t *testing.T) {
		l.SetHandler(2)
		require.Equal(t, l.handler, defaultHandler)
	})
}

func TestNewLogger(t *testing.T) {
	l := NewLogger()

	require.NotNil(t, l)
	require.Equal(t, l.out, defaultOut)
	require.Equal(t, l.level, defaultLevel)
}

func TestLoggerDebug(t *testing.T) {
	l := NewLogger()

	t.Run("level too low no print", func(t *testing.T) {
		b, _ := l.Debug("test")
		require.Zero(t, b)
	})

	t.Run("level matches print", func(t *testing.T) {
		l.SetLevel(LevelDebug)

		b, _ := l.Debug("test")
		require.NotZero(t, b)
	})
}

func TestLoggerInfo(t *testing.T) {
	l := NewLogger()

	t.Run("level enough print", func(t *testing.T) {
		b, _ := l.Info("test")
		require.NotZero(t, b)
	})

	t.Run("level too low no print", func(t *testing.T) {
		l.SetLevel(LevelWarn)

		b, _ := l.Info("test")
		require.Zero(t, b)
	})
}

func TestLoggerWarn(t *testing.T) {
	l := NewLogger()

	t.Run("level enough print", func(t *testing.T) {
		b, _ := l.Warn("test")
		require.NotZero(t, b)
	})

	t.Run("level too low no print", func(t *testing.T) {
		l.SetLevel(LevelError)

		b, _ := l.Warn("test")
		require.Zero(t, b)
	})
}

func TestLoggerError(t *testing.T) {
	l := NewLogger()

	t.Run("level enough print", func(t *testing.T) {
		b, _ := l.Error("test")
		require.NotZero(t, b)
	})

	t.Run("level too low no print", func(t *testing.T) {
		l.SetLevel(LevelFatal)

		b, _ := l.Error("test")
		require.Zero(t, b)
	})
}

func TestLoggerFatal(t *testing.T) {
	l := NewLogger()

	t.Run("fatal should always print and exit", func(t *testing.T) {
		exit = func(code int) {
			panic("this should panic")
		}

		require.Panics(t, func() { l.Fatal("test") })
	})
}

func TestLoggerMsgToString(t *testing.T) {
	l := NewLogger()

	t.Run("json handler", func(t *testing.T) {
		l.SetHandler(JSONHandler)

		msg := &msg{
			Timestamp: time.Time{},
			Level:     LevelInfo,
			Msg:       "test",
		}

		expected := `{"timestamp":"0001-01-01T00:00:00Z","level":"inf","msg":"test"}`
		b, _ := l.msgToString(msg)
		require.Equal(t, expected, b)
	})

	t.Run("invalid handler", func(t *testing.T) {
		l.handler = 99

		s, err := l.msgToString(nil)
		require.Error(t, err)
		require.Empty(t, s)
	})
}

func TestFormatOddArgs(t *testing.T) {
	cases := []struct {
		Name     string
		Args     []interface{}
		Expected []interface{}
	}{
		{
			Name:     "Non odd args",
			Args:     []interface{}{"key", "value"},
			Expected: []interface{}{"key", "value"},
		},
		{
			Name:     "Odd args",
			Args:     []interface{}{"value"},
			Expected: []interface{}{"no_key", "value"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			require.Equal(t, formatOddArgs(c.Args...), c.Expected)
		})
	}
}

func TestMsgFromParams(t *testing.T) {
	t.Run("msg from params", func(t *testing.T) {
		msg := msgFromParams(LevelInfo, "hello", "hello", "world")

		require.NotNil(t, msg)
	})
}

func TestArgsMapFromSlice(t *testing.T) {
	t.Run("args map from slice nil", func(t *testing.T) {
		require.Nil(t, argsMapFromSlice())
	})

	t.Run("args map from slice not nil", func(t *testing.T) {
		sampleArgs := make([]interface{}, 0, 2)
		sampleArgs = append(sampleArgs, "key")
		sampleArgs = append(sampleArgs, "value")

		sampleMap := map[string]interface{}{
			"key": "value",
		}

		require.NotNil(t, argsMapFromSlice(sampleArgs...))
		require.Equal(t, argsMapFromSlice(sampleArgs...), sampleMap)
	})
}
