package elasticsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		s, err := New("x/", "y", "z")
		require.NoError(t, err)
		require.NotNil(t, s)

		assert.Equal(t, "x", s.url)
		assert.Equal(t, "y", s.user)
		assert.Equal(t, "z", s.pass)
	})
}
