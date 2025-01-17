package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestOnGridContains(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), WithMask(2, 3, 4, 4, 3, 2))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		reels   []int
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "none - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 2, 0, 0},
		},
		{
			name:    "one - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 9, 0, 0},
			want:    true,
		},
		{
			name:    "few - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
			want:    true,
		},
		{
			name:    "many - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
			want:    true,
		},
		{
			name:    "none - reel 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 9, 0, 6, 9, 8, 1, 9, 3, 4, 5, 6, 7, 9, 0, 1, 2, 0, 0},
		},
		{
			name:    "one - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 9, 0, 0},
			want:    true,
		},
		{
			name:    "few - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
			want:    true,
		},
		{
			name:    "many - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := OnGridContains(tc.symbol, tc.reels...)
			require.NotNil(t, f)

			copy(spin.indexes, tc.indexes)

			got := f(spin)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOnGridNotContains(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), WithMask(2, 3, 4, 4, 3, 2))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		reels   []int
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "none - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 2, 0, 0},
			want:    true,
		},
		{
			name:    "one - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 9, 0, 0},
		},
		{
			name:    "few - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
		},
		{
			name:    "many - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
		},
		{
			name:    "none - reel 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 9, 0, 6, 9, 8, 1, 9, 3, 4, 5, 6, 7, 9, 0, 1, 2, 0, 0},
			want:    true,
		},
		{
			name:    "one - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 9, 0, 0},
		},
		{
			name:    "few - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
		},
		{
			name:    "many - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := OnGridNotContains(tc.symbol, tc.reels...)
			require.NotNil(t, f)

			copy(spin.indexes, tc.indexes)

			got := f(spin)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOnGridCount(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), WithMask(2, 3, 4, 4, 3, 2))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		reels   []int
		indexes utils.Indexes
		want    int
	}{
		{
			name:    "none - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 2, 0, 0},
		},
		{
			name:    "one - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 9, 0, 0},
			want:    1,
		},
		{
			name:    "few - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
			want:    3,
		},
		{
			name:    "many - all reels",
			symbol:  9,
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
			want:    8,
		},
		{
			name:    "none - reel 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 9, 0, 6, 9, 8, 1, 9, 3, 4, 5, 6, 7, 9, 0, 1, 2, 0, 0},
		},
		{
			name:    "one - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 9, 4, 5, 0, 6, 9, 8, 1, 9, 3, 4, 5, 9, 7, 8, 0, 1, 9, 0, 0},
			want:    1,
		},
		{
			name:    "few - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
			want:    3,
		},
		{
			name:    "many - reels 1 & 6",
			symbol:  9,
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
			want:    3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := OnGridCount(tc.symbol, tc.want, tc.reels...)
			require.NotNil(t, f)

			copy(spin.indexes, tc.indexes)

			got := f(spin)
			assert.True(t, got)
		})
	}
}

func TestOnGridCounts(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), WithMask(2, 3, 4, 4, 3, 2))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		counts  []int
		reels   []int
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "none - all reels",
			symbol:  9,
			counts:  []int{0, 4},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 2, 0, 0},
			want:    true,
		},
		{
			name:    "one - all reels",
			symbol:  9,
			counts:  []int{0, 4},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 0, 1, 9, 0, 0},
		},
		{
			name:    "few - all reels",
			symbol:  9,
			counts:  []int{0, 4},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
		},
		{
			name:    "many - all reels",
			symbol:  9,
			counts:  []int{0, 4},
			indexes: utils.Indexes{1, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
		},
		{
			name:    "none - reel 1 & 6",
			symbol:  9,
			counts:  []int{0, 4},
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 9, 0, 6, 9, 8, 1, 9, 3, 4, 5, 6, 7, 9, 0, 1, 2, 0, 0},
			want:    true,
		},
		{
			name:    "one - reels 1 & 6",
			symbol:  9,
			counts:  []int{0, 4},
			reels:   []int{1, 6},
			indexes: utils.Indexes{1, 2, 0, 0, 9, 4, 5, 0, 6, 9, 8, 1, 9, 3, 4, 5, 9, 7, 8, 0, 1, 9, 0, 0},
		},
		{
			name:    "few - reels 1 & 6",
			symbol:  9,
			counts:  []int{0, 4},
			reels:   []int{1, 6},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 6, 9, 8, 1, 2, 9, 4, 5, 6, 7, 8, 0, 9, 2, 0, 0},
		},
		{
			name:    "many - reels 1 & 6",
			symbol:  9,
			counts:  []int{0, 4},
			reels:   []int{1, 6},
			indexes: utils.Indexes{9, 9, 0, 0, 9, 4, 5, 0, 9, 7, 9, 1, 2, 3, 9, 5, 9, 7, 8, 0, 9, 9, 0, 0},
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := OnGridCounts(tc.symbol, tc.counts, tc.reels...)
			require.NotNil(t, f)

			copy(spin.indexes, tc.indexes)

			got := f(spin)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestOnGridShape(t *testing.T) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), WithMask(2, 3, 4, 4, 3, 2))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	testCases := []struct {
		name    string
		symbol  utils.Index
		mask    utils.UInt8s
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "not found (1)",
			symbol:  9,
			mask:    utils.UInt8s{0, 1, 20, 21},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 0, 8, 9, 0, 0},
		},
		{
			name:    "not found (2)",
			symbol:  9,
			mask:    utils.UInt8s{0, 1, 20, 21},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 0, 8, 9, 0, 0},
		},
		{
			name:    "not found (3)",
			symbol:  9,
			mask:    utils.UInt8s{0, 1, 20, 21},
			indexes: utils.Indexes{9, 2, 0, 0, 3, 4, 5, 0, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 0, 9, 9, 0, 0},
		},
		{
			name:    "not found (4)",
			symbol:  9,
			mask:    utils.UInt8s{0, 1, 20, 21},
			indexes: utils.Indexes{1, 9, 0, 0, 3, 4, 5, 0, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 0, 9, 9, 0, 0},
		},
		{
			name:    "found",
			symbol:  9,
			mask:    utils.UInt8s{0, 1, 20, 21},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 6, 7, 8, 9, 1, 2, 3, 4, 5, 6, 7, 0, 9, 9, 0, 0},
			want:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := OnGridShape(tc.symbol, tc.mask)
			require.NotNil(t, f)

			copy(spin.indexes, tc.indexes)

			got := f(spin)
			assert.Equal(t, tc.want, got)
		})
	}
}
