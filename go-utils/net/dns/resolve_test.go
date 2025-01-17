package dns

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithResolverEmpty(t *testing.T) {
	t.Run("with resolver empty", func(t *testing.T) {
		f := NewDialerWithResolver("", "", 0)
		assert.Nil(t, f)
	})
}

func TestWithResolverNotEmpty(t *testing.T) {
	t.Run("with resolver not empty", func(t *testing.T) {
		f := NewDialerWithResolver("1.1.1.1:53", "", 0)
		assert.NotNil(t, f)

		conn, err := f(context.Background(), "tcp", "xs4all.nl:80")
		require.NoError(t, err)
		require.NotNil(t, conn)

		conn.Close()
	})
}

func TestWithCacheEmpty(t *testing.T) {
	t.Run("with cache empty", func(t *testing.T) {
		f := NewDialerWithCache("", "", 0, 0, "help.me")
		assert.Nil(t, f)
	})
}

func TestWithCacheNotEmpty(t *testing.T) {
	t.Run("with cache not empty", func(t *testing.T) {
		f := NewDialerWithCache("1.1.1.1:53", "", 0, 0, "gw.dev.topgaming.team")
		assert.NotNil(t, f)

		conn, err := f(context.Background(), "tcp", "gw.dev.topgaming.team:443")
		require.NoError(t, err)
		require.NotNil(t, conn)

		conn.Close()
	})
}
