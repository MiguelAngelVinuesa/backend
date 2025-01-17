package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBonusSymbolPayout(t *testing.T) {
	t.Run("new bonus symbol payout", func(t *testing.T) {
		p := BonusSymbolPayout(10, 1, 11, 4).(*SpinPayout)
		require.NotNil(t, p)
		defer p.Release()

		assert.Equal(t, results.SlotBonusSymbol, p.kind)
		assert.Equal(t, 10.0, p.factor)
		assert.Equal(t, utils.Index(11), p.symbol)
		assert.Equal(t, uint8(4), p.count)
		require.NotNil(t, p.payRows)
		require.NotNil(t, p.payMap)
		assert.Empty(t, p.payRows)
		assert.Empty(t, p.payMap)
	})
}
