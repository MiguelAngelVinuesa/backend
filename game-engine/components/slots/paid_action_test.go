package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPaidTrigger(t *testing.T) {
	testCases := []struct {
		name          string
		kind          SpinActionResult
		nrOfSpins     uint8
		betMultiplier int
		symbol        utils.Index
		triggerCount  uint8
	}{
		{"small game", FreeSpins, 3, 10, 11, 3},
		{"big game", FreeSpins, 10, 100, 11, 3},
		{"bonus game", BonusGame, 0, 100, 14, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := NewPaidAction(tc.kind, tc.nrOfSpins, tc.betMultiplier, tc.symbol, tc.triggerCount)
			require.NotNil(t, w)

			assert.Equal(t, tc.kind, w.result)
			assert.Equal(t, tc.nrOfSpins, w.nrOfSpins)
			assert.Equal(t, tc.symbol, w.symbol)
			assert.Equal(t, tc.betMultiplier, w.betMultiplier)
			assert.Equal(t, tc.triggerCount, w.triggerCount)
		})
	}
}

type paidTriggerGame struct{}

func (g *paidTriggerGame) RequireParams() bool { return false }
func (g *paidTriggerGame) Run(result *results.Result, _ ...interface{}) *results.Result {
	return result
}
