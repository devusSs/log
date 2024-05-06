package log

type Handler int

const (
	TextHandler Handler = iota
	JSONHandler
)

func (h Handler) String() string {
	switch h {
	case TextHandler:
		return "text"
	case JSONHandler:
		return "json"
	default:
		return "unknown"
	}
}
