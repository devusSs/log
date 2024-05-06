package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

type msg struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     Level                  `json:"level"`
	Msg       string                 `json:"msg"`
	Args      map[string]interface{} `json:"-"`
}

func (m *msg) String() string {
	ts := formatTimestampRFC3339(m.Timestamp)
	l := formatLevel(m.Level)

	if m.Args == nil || len(m.Args) == 0 {
		return fmt.Sprintf(`timestamp=%s level=%s msg="%s"`, ts, l, m.Msg)
	}

	return fmt.Sprintf(
		`timestamp=%s level=%s msg="%s" %v`,
		formatTimestampRFC3339(m.Timestamp),
		formatLevel(m.Level),
		m.Msg,
		formatArgs(m.Args),
	)
}

func (m *msg) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	buf.WriteString(`"timestamp":`)
	jsonValue, _ := json.Marshal(formatTimestampRFC3339(m.Timestamp))
	buf.Write(jsonValue)
	buf.WriteString(",")

	buf.WriteString(`"level":`)
	jsonValue, _ = json.Marshal(m.Level.String())
	buf.Write(jsonValue)
	buf.WriteString(",")

	buf.WriteString(`"msg":`)
	jsonValue, _ = json.Marshal(m.Msg)
	buf.Write(jsonValue)

	for key, value := range m.Args {
		buf.WriteString(`,"`)
		buf.WriteString(key)
		buf.WriteString(`":`)
		jsonValue, _ = json.Marshal(value)
		buf.Write(jsonValue)
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}

func formatTimestampRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

func formatLevel(l Level) string {
	level := strings.ToUpper(l.String())

	if noColor {
		return strings.ToUpper(level)
	}

	switch l {
	case LevelDebug:
		return colorString(level, colorWhite)
	case LevelInfo:
		return colorString(level, colorCyan)
	case LevelWarn:
		return colorString(level, colorYellow)
	case LevelError:
		return colorString(level, colorRed)
	case LevelFatal:
		return colorString(level, colorRed)
	default:
		return colorString(level, colorWhite)
	}
}

func formatArgs(args map[string]interface{}) string {
	if len(args) == 0 {
		return ""
	}

	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	buf := ""

	for i, key := range keys {
		value := args[key]

		kv := fmt.Sprintf("%s=%v ", key, checkStringType(value))

		if len(keys)-1 == i {
			kv = fmt.Sprintf("%s=%v", key, checkStringType(value))
		}

		buf += kv
	}

	return buf
}

func checkStringType(v interface{}) interface{} {
	if v == nil {
		return ""
	}

	if reflect.TypeOf(v).Kind() == reflect.String {
		return fmt.Sprintf(`"%v"`, v)
	}

	return v
}
