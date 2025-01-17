package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRoundFlag(t *testing.T) {
	testCases := []struct {
		name string
		id   int
		flag string
	}{
		{name: "1", id: 1, flag: "flag 1"},
		{name: "2", id: 2, flag: "free spin counter"},
		{name: "3", id: 3, flag: "payout band"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := NewRoundFlag(tc.id, tc.flag)
			require.NotNil(t, f)
			assert.Equal(t, tc.id, f.ID())
			assert.Equal(t, tc.flag, f.Name())
		})
	}
}
