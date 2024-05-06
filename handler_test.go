package log

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandlerString(t *testing.T) {
	t.Run("text handler", func(t *testing.T) {
		require.Equal(t, TextHandler.String(), "text")
	})

	t.Run("json handler", func(t *testing.T) {
		require.Equal(t, JSONHandler.String(), "json")
	})

	t.Run("unknown handler", func(t *testing.T) {
		require.Equal(t, Handler(2).String(), "unknown")
	})
}
