package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewClearWinlinesAction(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots1 := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayBoth))
	spin1 := AcquireSpin(slots1, prng)
	defer spin1.Release()

	slots2 := NewSlots(Grid(5, 3), WithSymbols(setF1))
	spin2 := AcquireSpin(slots2, prng)
	defer spin2.Release()

	testCases := []struct {
		name    string
		spin    *Spin
		indexes utils.Indexes
		payouts utils.UInt8s
		want    bool
		after   utils.Indexes
	}{
		{
			name:    "paylines, no wins",
			spin:    spin1,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "paylines, 1 win LTR",
			spin:    spin1,
			indexes: utils.Indexes{1, 2, 3, 4, 2, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			want:    true,
			after:   utils.Indexes{1, 0, 3, 4, 0, 6, 1, 0, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "paylines, 1 win RTL",
			spin:    spin1,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 2, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0},
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 0, 3, 4, 0, 6, 1, 0, 3},
		},
		{
			name:    "paylines, 2 wins BOTH",
			spin:    spin1,
			indexes: utils.Indexes{1, 2, 3, 4, 1, 6, 1, 2, 9, 4, 1, 3, 1, 2, 3},
			payouts: utils.UInt8s{1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
			want:    true,
			after:   utils.Indexes{0, 2, 3, 4, 0, 6, 1, 2, 0, 4, 0, 0, 0, 2, 0},
		},
		{
			name:    "all paylines, no wins",
			spin:    spin2,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "all paylines, 1 win LTR",
			spin:    spin2,
			indexes: utils.Indexes{1, 2, 3, 3, 5, 6, 1, 2, 3, 3, 5, 6, 1, 2, 3},
			payouts: utils.UInt8s{0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1},
			want:    true,
			after:   utils.Indexes{1, 2, 0, 0, 5, 6, 1, 2, 0, 0, 5, 6, 1, 2, 0},
		},
		{
			name:    "all paylines, 2 wins LTR",
			spin:    spin2,
			indexes: utils.Indexes{1, 2, 3, 4, 1, 3, 1, 2, 9, 4, 1, 6, 1, 2, 3},
			payouts: utils.UInt8s{1, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0},
			want:    true,
			after:   utils.Indexes{0, 2, 0, 4, 0, 0, 1, 2, 0, 4, 0, 6, 0, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewClearPayoutsAction()

			tc.spin.indexes = tc.indexes
			tc.spin.payouts = tc.payouts

			got := a.Triggered(tc.spin) == a
			assert.Equal(t, tc.want, got)
			assert.EqualValues(t, tc.after, tc.spin.indexes)
		})
	}
}

func TestNewExplodingBombsAction(t *testing.T) {
	bf1 := NewSymbol(7, WithKind(Bomb))
	symbols := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, bf1)

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(symbols))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	gridPlus := GridOffsets{{-1, 0}, {0, -1}, {0, 0}, {0, 1}, {1, 0}}
	gridX := GridOffsets{{-1, -1}, {-1, 1}, {0, 0}, {1, -1}, {1, 1}}
	grid3x3 := GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	testCases := []struct {
		name    string
		spin    *Spin
		indexes utils.Indexes
		offsets GridOffsets
		want    bool
		after   utils.Indexes
	}{
		{
			name:    "+, no bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridPlus,
			want:    false,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "+, 1 bomb top-left",
			spin:    spin,
			indexes: utils.Indexes{7, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{0, 0, 3, 0, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "+, 1 bomb bottom-left",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 7, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{1, 0, 0, 4, 5, 0, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "+, 1 bomb middle",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 7, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 0, 6, 0, 0, 0, 4, 0, 6, 1, 2, 3},
		},
		{
			name:    "+, 1 bomb top-right",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 2, 3},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 0, 5, 6, 0, 0, 3},
		},
		{
			name:    "+, 1 bomb bottom-right",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 7},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 0, 1, 0, 0},
		},
		{
			name:    "+, 2 bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 7, 3, 4, 7, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{0, 0, 0, 0, 0, 0, 1, 0, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "+, 3 bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 7, 3, 4, 7, 7, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridPlus,
			want:    true,
			after:   utils.Indexes{0, 0, 0, 0, 0, 0, 1, 0, 0, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "X, no bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridX,
			want:    false,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "X, 1 bomb, top-left",
			spin:    spin,
			indexes: utils.Indexes{7, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{0, 2, 3, 4, 0, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "X, 1 bomb, bottom-left",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 7, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{1, 2, 0, 4, 0, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "X, 1 bomb, middle",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 7, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{0, 2, 0, 4, 0, 6, 0, 2, 0, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "X, 1 bomb, top-right",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 2, 3},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 0, 6, 0, 2, 3},
		},
		{
			name:    "X, 1 bomb, bottom-right",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 7},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 0, 6, 1, 2, 0},
		},
		{
			name:    "X, 2 bombs",
			spin:    spin,
			indexes: utils.Indexes{7, 2, 3, 4, 7, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{0, 2, 0, 4, 0, 6, 0, 2, 0, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "X, 3 bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 7, 4, 7, 6, 7, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: gridX,
			want:    true,
			after:   utils.Indexes{0, 2, 0, 4, 0, 6, 0, 2, 0, 4, 0, 6, 1, 2, 3},
		},
		{
			name:    "3x3, no bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: grid3x3,
			want:    false,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "3x3, 1 bomb, top-left",
			spin:    spin,
			indexes: utils.Indexes{7, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{0, 0, 3, 0, 0, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "3x3, 1 bomb, bottom-left",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 7, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{1, 0, 0, 4, 0, 0, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "3x3, 1 bomb, middle",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 7, 6, 1, 2, 3},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "3x3, 1 bomb, top-right",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 7, 2, 3},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 0, 0, 6, 0, 0, 3},
		},
		{
			name:    "3x3, 1 bomb, bottom-right",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 7},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 0, 0, 1, 0, 0},
		},
		{
			name:    "3x3, 2 bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 7, 3, 4, 7, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "3x3, 3 bombs",
			spin:    spin,
			indexes: utils.Indexes{1, 2, 7, 4, 5, 7, 1, 7, 3, 4, 5, 6, 1, 2, 3},
			offsets: grid3x3,
			want:    true,
			after:   utils.Indexes{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewExplodingBombsAction(7, tc.offsets)

			tc.spin.indexes = tc.indexes

			got := a.Triggered(tc.spin) == a
			assert.Equal(t, tc.want, got)
			assert.EqualValues(t, tc.after, tc.spin.indexes)
		})
	}
}
