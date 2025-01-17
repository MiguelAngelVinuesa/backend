package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewNudgeAction(t *testing.T) {
	testCases := []struct {
		name     string
		count    uint8
		location NudgeLocation
		chance   float64
		tease    float64
		reels    utils.UInt8s
		indexes  utils.Indexes
		want     int
	}{
		{
			name:     "top, 25%, no tease, no trigger",
			count:    2,
			location: NudgeTop,
			chance:   25,
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:     "top, 25%, no tease, trigger",
			count:    2,
			location: NudgeTop,
			chance:   25,
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     3,
		},
		{
			name:     "top, 25%, no tease, reels 2/3/4, trigger",
			count:    2,
			location: NudgeTop,
			chance:   25,
			reels:    utils.UInt8s{2, 3, 4},
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     2,
		},
		{
			name:     "top, 25%, tease, no trigger",
			count:    2,
			location: NudgeTop,
			chance:   25,
			tease:    40,
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:     "top, 25%, tease, trigger",
			count:    2,
			location: NudgeTop,
			chance:   25,
			tease:    40,
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     3,
		},
		{
			name:     "top, 25%, tease, reels 2/3/4, trigger",
			count:    2,
			location: NudgeTop,
			chance:   25,
			tease:    5,
			reels:    utils.UInt8s{2, 3, 4},
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     2,
		},
		{
			name:     "bottom, 25%, no tease, no trigger",
			count:    2,
			location: NudgeBottom,
			chance:   25,
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:     "bottom, 25%, no tease, trigger",
			count:    2,
			location: NudgeBottom,
			chance:   25,
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     3,
		},
		{
			name:     "bottom, 25%, no tease, reels 2/3/4, trigger",
			count:    2,
			location: NudgeBottom,
			chance:   25,
			reels:    utils.UInt8s{2, 3, 4},
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     2,
		},
		{
			name:     "bottom, 25%, tease, no trigger",
			count:    2,
			location: NudgeBottom,
			chance:   25,
			tease:    10,
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:     "bottom, 25%, tease, trigger",
			count:    2,
			location: NudgeBottom,
			chance:   25,
			tease:    10,
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     3,
		},
		{
			name:     "bottom, 25%, tease, reels 2/3/4, trigger",
			count:    2,
			location: NudgeBottom,
			chance:   25,
			tease:    50,
			reels:    utils.UInt8s{2, 3, 4},
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     2,
		},
		{
			name:     "vertical, 25%, no tease, no trigger",
			count:    2,
			location: NudgeVertical,
			chance:   25,
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:     "vertical, 25%, no tease, trigger",
			count:    2,
			location: NudgeVertical,
			chance:   25,
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     3,
		},
		{
			name:     "vertical, 25%, no tease, reels 2/3/4, trigger",
			count:    2,
			location: NudgeVertical,
			chance:   50,
			reels:    utils.UInt8s{2, 3, 4},
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 3, 4, 11, 6, 1, 2, 3},
			want:     2,
		},
		{
			name:     "vertical, 25%, tease, no trigger",
			count:    2,
			location: NudgeVertical,
			chance:   25,
			tease:    10,
			indexes:  utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:     "vertical, 25%, tease, trigger",
			count:    2,
			location: NudgeVertical,
			chance:   25,
			tease:    10,
			indexes:  utils.Indexes{1, 2, 11, 4, 5, 6, 1, 2, 11, 4, 5, 6, 1, 2, 3},
			want:     3,
		},
		{
			name:     "vertical, 25%, tease, reels 2/3/4, trigger",
			count:    2,
			location: NudgeVertical,
			chance:   10,
			tease:    2,
			reels:    utils.UInt8s{2, 3, 4},
			indexes:  utils.Indexes{1, 2, 11, 4, 11, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:     2,
		},
	}

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))
	symbol := utils.Index(11)

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewNudgeAction(symbol, tc.count, tc.location, tc.chance).GenerateNoDupes()
			require.NotNil(t, a)

			if tc.tease > 0 {
				a.WithTease(tc.tease)
			}
			if len(tc.reels) > 0 {
				a.WithReels(tc.reels...)
			}

			var teasers int
			reels := make(map[uint8]int, 5)
			locs := make(map[NudgeLocation]int)

			for ix := 0; ix < 10000; ix++ {
				res := AcquireSpinResult(spin)

				copy(spin.indexes, tc.indexes)

				if a2 := a.Triggered(spin); a2 != nil {
					require.Equal(t, a, a2)

					a2.Nudge(spin, res)
					require.Equal(t, 1, len(res.nudges))

					n := res.nudges[0]
					assert.Equal(t, symbol, n.symbol)
					assert.Equal(t, uint8(1), n.size)

					if n.teaser {
						teasers++
					}

					reels[n.reel] = reels[n.reel] + 1
					locs[n.location] = locs[n.location] + 1
				}

				res.Release()
			}

			if tc.want == 0 {
				assert.Zero(t, teasers)
				assert.Zero(t, len(reels))
				assert.Zero(t, len(locs))
			} else {
				if tc.tease > 0 {
					assert.NotZero(t, teasers)
				} else {
					assert.Zero(t, teasers)
				}

				assert.Zero(t, reels[0])
				assert.Zero(t, reels[6])
				assert.Equal(t, tc.want, len(reels))

				switch tc.location {
				case NudgeTop:
					assert.Zero(t, locs[NudgeVertical])
					assert.NotZero(t, locs[NudgeTop])
					assert.Zero(t, locs[NudgeBottom])
				case NudgeBottom:
					assert.Zero(t, locs[NudgeVertical])
					assert.Zero(t, locs[NudgeTop])
					assert.NotZero(t, locs[NudgeBottom])
				case NudgeVertical:
					assert.Zero(t, locs[NudgeVertical])
					assert.NotZero(t, locs[NudgeTop])
					assert.NotZero(t, locs[NudgeBottom])
				}
			}
		})
	}
}
