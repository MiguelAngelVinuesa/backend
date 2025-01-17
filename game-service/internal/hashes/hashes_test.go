package hashes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashes(t *testing.T) {
	t.Run("hashes", func(t *testing.T) {
		assert.NotEmpty(t, MainFile)
		assert.NotEmpty(t, MainHash)
		assert.NotEmpty(t, RngLibHash)
		assert.NotEmpty(t, RngIncludeHash)
		assert.NotEqual(t, errorHash, MainHash)
		assert.NotEqual(t, errorHash, RngLibHash)
		assert.NotEqual(t, errorHash, RngIncludeHash)
	})
}
