package log

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	out     io.Writer
	level   Level
	handler Handler
}

func (l *Logger) SetOut(w io.Writer) {
	if w == nil {
		l.out = defaultOut
		return
	}

	l.out = w
}

func (l *Logger) SetLevel(level Level) {
	if level < LevelInvalid || level > LevelFatal {
		l.level = defaultLevel
		return
	}

	if level == LevelInvalid {
		l.level = defaultLevel
		return
	}

	l.level = level
}

func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) SetHandler(handler Handler) {
	if handler != 0 && handler != 1 {
		l.handler = defaultHandler
		return
	}

	l.handler = handler
}

func NewLogger() *Logger {
	return &Logger{
		out:     defaultOut,
		level:   defaultLevel,
		handler: defaultHandler,
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) (int, error) {
	if !evalLevel(LevelDebug, l.level) {
		return 0, nil
	}

	args = formatOddArgs(args...)

	out := msgFromParams(LevelDebug, msg, args...)

	outStr, _ := l.msgToString(out)

	return fmt.Fprintln(l.out, outStr)
}

func (l *Logger) Info(msg string, args ...interface{}) (int, error) {
	if !evalLevel(LevelInfo, l.level) {
		return 0, nil
	}

	args = formatOddArgs(args...)

	out := msgFromParams(LevelInfo, msg, args...)

	outStr, _ := l.msgToString(out)

	return fmt.Fprintln(l.out, outStr)
}

func (l *Logger) Warn(msg string, args ...interface{}) (int, error) {
	if !evalLevel(LevelWarn, l.level) {
		return 0, nil
	}

	args = formatOddArgs(args...)

	out := msgFromParams(LevelWarn, msg, args...)

	outStr, _ := l.msgToString(out)

	return fmt.Fprintln(l.out, outStr)
}

func (l *Logger) Error(msg string, args ...interface{}) (int, error) {
	if !evalLevel(LevelError, l.level) {
		return 0, nil
	}

	args = formatOddArgs(args...)

	out := msgFromParams(LevelError, msg, args...)

	outStr, _ := l.msgToString(out)

	return fmt.Fprintln(l.out, outStr)
}

var exit func(code int) = os.Exit

func (l *Logger) Fatal(msg string, args ...interface{}) {
	args = formatOddArgs(args...)

	out := msgFromParams(LevelFatal, msg, args...)

	outStr, _ := l.msgToString(out)

	_, _ = fmt.Fprintln(l.out, outStr)
	exit(1)
}

func (l *Logger) msgToString(msg *msg) (string, error) {
	switch l.handler {
	case TextHandler:
		return msg.String(), nil
	case JSONHandler:
		b, _ := msg.Marshal()
		return string(b), nil
	default:
		return "", errors.New("invalid handler")
	}
}

var (
	defaultOut     io.Writer = os.Stderr
	defaultLevel   Level     = LevelInfo
	defaultHandler Handler   = TextHandler
)

func formatOddArgs(args ...interface{}) []interface{} {
	if len(args)%2 != 0 {
		argsCopy := make([]interface{}, len(args)+1)

		original := args[len(args)-1]
		argsCopy[len(argsCopy)-2] = "no_key"
		argsCopy[len(argsCopy)-1] = original

		return argsCopy
	}

	return args
}

func msgFromParams(level Level, msgStr string, args ...interface{}) *msg {
	if len(args) > 0 {
		out := &msg{
			Timestamp: time.Now(),
			Level:     level,
			Msg:       msgStr,
			Args:      argsMapFromSlice(args...),
		}

		return out
	}

	return &msg{
		Timestamp: time.Now(),
		Level:     level,
		Msg:       msgStr,
		Args:      nil,
	}
}

func argsMapFromSlice(args ...interface{}) map[string]interface{} {
	if len(args) == 0 {
		return nil
	}

	resultMap := make(map[string]interface{})

	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprintf("%v", args[i])
		resultMap[key] = args[i+1]
	}

	return resultMap
}
