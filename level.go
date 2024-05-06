package log

import (
	"errors"
)

var ErrInvalidLevel error = errors.New("invalid log level")

type Level int

const (
	LevelInvalid Level = iota - 1
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "dbg"
	case LevelInfo:
		return "inf"
	case LevelWarn:
		return "wrn"
	case LevelError:
		return "err"
	case LevelFatal:
		return "ftl"
	default:
		return "invalid"
	}
}

func ParseLevel(level string) (Level, error) {
	l, ok := levelStrings[level]
	if !ok {
		return LevelInvalid, ErrInvalidLevel
	}

	return l, nil
}

func MustParseLevel(level string) Level {
	l, err := ParseLevel(level)
	if err != nil {
		panic(err)
	}

	return l
}

var levelStrings map[string]Level = map[string]Level{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

func evalLevel(got Level, want Level) bool {
	return got >= want
}
