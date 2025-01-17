package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestSpin_PreventPaylines(t *testing.T) {
	prng := rng.AcquireRNG()
	defer prng.ReturnToPool()

	myset := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wwf1, hf1, wwf2, hf2)
	slots := NewSlots(Grid(5, 3), WithSymbols(myset))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name      string
		direction PayDirection
		remWilds  bool
		reelDupes bool
		indexes   utils.Indexes
		want      bool
	}{
		{
			name:      "ltr, none",
			direction: PayLTR,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "ltr, one",
			direction: PayLTR,
			indexes:   utils.Indexes{1, 4, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, three",
			direction: PayLTR,
			indexes:   utils.Indexes{4, 6, 5, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, wild",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, wilds",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{9, 9, 9, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, wilds, skip 2nd reel",
			direction: PayLTR,
			indexes:   utils.Indexes{1, 2, 3, 9, 5, 9, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "ltr, wilds, not skip 2nd reel",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 2, 3, 9, 5, 9, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "ltr, two and wilds",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{9, 4, 5, 4, 5, 9, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:      true,
		},
		{
			name:      "rtl, none",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "rtl, one",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 6, 2, 3},
			want:      true,
		},
		{
			name:      "rtl, three",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 4, 5, 6},
			want:      true,
		},
		{
			name:      "rtl, wild",
			direction: PayRTL,
			remWilds:  true,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 9},
			want:      true,
		},
		{
			name:      "rtl, wilds",
			direction: PayRTL,
			remWilds:  true,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 9, 9, 9},
			want:      true,
		},
		{
			name:      "rtl, wilds, skip 2nd reel",
			direction: PayRTL,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 9, 9, 6, 1, 2, 3},
		},
		{
			name:      "rtl, wilds, not skip 2nd reel",
			direction: PayRTL,
			remWilds:  true,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 9, 9, 1, 2, 3},
			want:      true,
		},
		{
			name:      "rtl, two and wilds",
			direction: PayRTL,
			remWilds:  true,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 9, 9, 4, 5},
			want:      true,
		},
		{
			name:      "misc (1)",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 4, 3, 1, 2, 3, 9, 2, 3, 4, 10, 5, 8, 2, 7},
			want:      true,
		},
		{
			name:      "misc (2)",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 4, 3, 1, 2, 10, 9, 2, 3, 4, 10, 5, 8, 2, 7},
			want:      true,
		},
		{
			name:      "misc (3)",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 4, 10, 1, 2, 3, 9, 2, 3, 4, 10, 5, 8, 2, 7},
			want:      true,
		},
		{
			name:      "misc (4)",
			direction: PayLTR,
			remWilds:  true,
			indexes:   utils.Indexes{1, 4, 10, 1, 2, 10, 9, 2, 3, 4, 10, 5, 8, 2, 7},
			want:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for ix := 0; ix < 1000; ix++ {
				copy(spin.indexes, tc.indexes)

				got := spin.PreventPaylines(tc.direction, tc.remWilds, tc.reelDupes)
				if tc.want {
					assert.True(t, got)

					if tc.direction == PayLTR || tc.direction == PayBoth {
						assert.NotEqual(t, spin.indexes[0], spin.indexes[3])
						assert.NotEqual(t, spin.indexes[0], spin.indexes[4])
						assert.NotEqual(t, spin.indexes[0], spin.indexes[5])
						assert.NotEqual(t, spin.indexes[0], utils.Index(9))
						assert.NotEqual(t, spin.indexes[0], utils.Index(12))
						assert.NotEqual(t, spin.indexes[1], spin.indexes[3])
						assert.NotEqual(t, spin.indexes[1], spin.indexes[4])
						assert.NotEqual(t, spin.indexes[1], spin.indexes[5])
						assert.NotEqual(t, spin.indexes[1], utils.Index(9))
						assert.NotEqual(t, spin.indexes[1], utils.Index(12))
						assert.NotEqual(t, spin.indexes[2], spin.indexes[3])
						assert.NotEqual(t, spin.indexes[2], spin.indexes[4])
						assert.NotEqual(t, spin.indexes[2], spin.indexes[5])
						assert.NotEqual(t, spin.indexes[2], utils.Index(9))
						assert.NotEqual(t, spin.indexes[2], utils.Index(12))

						if !tc.reelDupes {
							assert.NotEqual(t, spin.indexes[0], spin.indexes[1])
							assert.NotEqual(t, spin.indexes[0], spin.indexes[2])
							assert.NotEqual(t, spin.indexes[1], spin.indexes[2])
						}
					} else {
						assert.NotEqual(t, spin.indexes[9], spin.indexes[12])
						assert.NotEqual(t, spin.indexes[9], spin.indexes[13])
						assert.NotEqual(t, spin.indexes[9], spin.indexes[14])
						assert.NotEqual(t, spin.indexes[9], utils.Index(9))
						assert.NotEqual(t, spin.indexes[9], utils.Index(12))
						assert.NotEqual(t, spin.indexes[10], spin.indexes[12])
						assert.NotEqual(t, spin.indexes[10], spin.indexes[13])
						assert.NotEqual(t, spin.indexes[10], spin.indexes[14])
						assert.NotEqual(t, spin.indexes[10], utils.Index(9))
						assert.NotEqual(t, spin.indexes[10], utils.Index(12))
						assert.NotEqual(t, spin.indexes[11], spin.indexes[12])
						assert.NotEqual(t, spin.indexes[11], spin.indexes[13])
						assert.NotEqual(t, spin.indexes[11], spin.indexes[14])
						assert.NotEqual(t, spin.indexes[11], utils.Index(9))
						assert.NotEqual(t, spin.indexes[11], utils.Index(12))

						if !tc.reelDupes {
							assert.NotEqual(t, spin.indexes[12], spin.indexes[13])
							assert.NotEqual(t, spin.indexes[12], spin.indexes[14])
							assert.NotEqual(t, spin.indexes[13], spin.indexes[14])
						}
					}

				} else {
					assert.False(t, got)
					assert.EqualValues(t, tc.indexes, spin.indexes)
				}
			}
		})
	}
}

func TestSpin_PreventBonus(t *testing.T) {
	prng := rng.AcquireRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name      string
		reelDupes bool
		indexes   utils.Indexes
		bonus     utils.Index
		want      bool
	}{
		{
			name:    "none",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			bonus:   7,
		},
		{
			name:    "one",
			indexes: utils.Indexes{1, 4, 3, 4, 5, 6, 1, 2, 3, 4, 7, 6, 1, 2, 3},
			want:    true,
			bonus:   5,
		},
		{
			name:    "three",
			indexes: utils.Indexes{4, 6, 5, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    true,
			bonus:   6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for ix := 0; ix < 1000; ix++ {
				copy(spin.indexes, tc.indexes)
				spin.bonusSymbol = tc.bonus

				got := spin.PreventBonus(tc.reelDupes)
				if tc.want {
					assert.True(t, got)
					if !tc.reelDupes {
						assert.NotEqual(t, spin.indexes[0], spin.indexes[1])
						assert.NotEqual(t, spin.indexes[0], spin.indexes[2])
						assert.NotEqual(t, spin.indexes[1], spin.indexes[2])
						assert.NotEqual(t, spin.indexes[3], spin.indexes[4])
						assert.NotEqual(t, spin.indexes[3], spin.indexes[5])
						assert.NotEqual(t, spin.indexes[4], spin.indexes[5])
						assert.NotEqual(t, spin.indexes[6], spin.indexes[7])
						assert.NotEqual(t, spin.indexes[6], spin.indexes[8])
						assert.NotEqual(t, spin.indexes[7], spin.indexes[8])
						assert.NotEqual(t, spin.indexes[9], spin.indexes[10])
						assert.NotEqual(t, spin.indexes[9], spin.indexes[11])
						assert.NotEqual(t, spin.indexes[10], spin.indexes[11])
						assert.NotEqual(t, spin.indexes[12], spin.indexes[13])
						assert.NotEqual(t, spin.indexes[12], spin.indexes[14])
						assert.NotEqual(t, spin.indexes[13], spin.indexes[14])
					}
				} else {
					assert.False(t, got)
					assert.EqualValues(t, tc.indexes, spin.indexes)
				}
			}
		})
	}
}
