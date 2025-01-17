package plural

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_String(t *testing.T) {
	assert.Equal(t, "zero", Zero.String())
	assert.Equal(t, "one", One.String())
	assert.Equal(t, "two", Two.String())
	assert.Equal(t, "few", Few.String())
	assert.Equal(t, "many", Many.String())
	assert.Equal(t, "other", Other.String())
	assert.Equal(t, "???", Category(99).String())
}
