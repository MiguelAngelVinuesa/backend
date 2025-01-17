package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewHotAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayBoth))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		indexes utils.Indexes
		want    bool
		hot     utils.UInt8s
	}{
		{"none", 8, utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}, false, utils.UInt8s{}},
		{"1", 10, utils.Indexes{1, 2, 3, 4, 5, 10, 1, 2, 3, 4, 5, 6, 1, 2, 3}, true, utils.UInt8s{2}},
		{"2", 12, utils.Indexes{10, 2, 3, 4, 5, 12, 10, 2, 3, 4, 5, 12, 1, 2, 3}, true, utils.UInt8s{2, 4}},
		{"3", 10, utils.Indexes{10, 12, 3, 4, 10, 10, 11, 12, 3, 4, 5, 10, 11, 2, 3}, true, utils.UInt8s{1, 2, 4}},
		{"all", 10, utils.Indexes{1, 10, 10, 4, 5, 10, 11, 12, 10, 4, 5, 10, 10, 12, 3}, true, utils.UInt8s{1, 2, 3, 4, 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewHotAction(tc.symbol)
			require.NotNil(t, a)

			assert.Equal(t, AwardBonuses, a.stage)
			assert.Equal(t, HotReel, a.result)
			assert.Equal(t, tc.symbol, a.symbol)

			spin.resetHot()
			spin.indexes = tc.indexes

			got := a.Triggered(spin) == a
			assert.Equal(t, tc.want, got)

			hot := spin.Hot(make(utils.UInt8s, 5))
			assert.EqualValues(t, tc.hot, hot)
		})
	}
}
