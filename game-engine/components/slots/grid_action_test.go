package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

var (
	superX      = GridOffsets{{-1, -1}, {-1, 1}, {0, 0}, {1, -1}, {1, 1}}
	superPlus   = GridOffsets{{-1, 0}, {0, -1}, {0, 0}, {0, 1}, {1, 0}}
	centersX    = GridOffsets{{1, 1}, {2, 1}, {3, 1}}
	centersPlus = GridOffsets{{1, 1}, {2, 1}, {3, 1}}
)

func TestNewGridShapeAction(t *testing.T) {
	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name    string
		shape   GridOffsets
		centers GridOffsets
		indexes utils.Indexes
		want    bool
		symbol  utils.Index
		sticky  []bool
	}{
		{
			name:    "SuperX, no hit",
			shape:   superX,
			centers: centersX,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    false,
		},
		{
			name:    "SuperX, almost hit",
			shape:   superX,
			centers: centersX,
			indexes: utils.Indexes{1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 2, 3, 2, 2, 2},
			want:    false,
		},
		{
			name:    "SuperX, hit left",
			shape:   superX,
			centers: centersX,
			indexes: utils.Indexes{1, 2, 1, 4, 1, 6, 1, 2, 1, 4, 5, 6, 2, 2, 3},
			symbol:  1,
			want:    true,
			sticky:  []bool{true, false, true, false, true, false, true, false, true, false, false, false, false, false, false},
		},
		{
			name:    "SuperX, hit middle",
			shape:   superX,
			centers: centersX,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 4, 1, 4, 3, 4, 5, 4, 1, 2, 3},
			symbol:  4,
			want:    true,
			sticky:  []bool{false, false, false, true, false, true, false, true, false, true, false, true, false, false, false},
		},
		{
			name:    "SuperX, hit right",
			shape:   superX,
			centers: centersX,
			indexes: utils.Indexes{5, 4, 5, 4, 5, 2, 6, 2, 6, 4, 6, 3, 6, 5, 6},
			symbol:  6,
			want:    true,
			sticky:  []bool{false, false, false, false, false, false, true, false, true, false, true, false, true, false, true},
		},
		{
			name:    "SuperPlus, no hit",
			shape:   superPlus,
			centers: centersPlus,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    false,
		},
		{
			name:    "SuperPlus, almost hit",
			shape:   superPlus,
			centers: centersPlus,
			indexes: utils.Indexes{1, 2, 3, 2, 2, 2, 1, 1, 3, 2, 2, 6, 1, 2, 3},
			want:    false,
		},
		{
			name:    "SuperPlus, hit left",
			shape:   superPlus,
			centers: centersPlus,
			indexes: utils.Indexes{1, 2, 3, 2, 2, 2, 1, 2, 3, 3, 4, 6, 1, 5, 3},
			symbol:  2,
			want:    true,
			sticky:  []bool{false, true, false, true, true, true, false, true, false, false, false, false, false, false, false},
		},
		{
			name:    "SuperPlus, hit middle",
			shape:   superPlus,
			centers: centersPlus,
			indexes: utils.Indexes{4, 3, 6, 1, 5, 3, 5, 5, 5, 1, 5, 3, 1, 2, 6},
			symbol:  5,
			want:    true,
			sticky:  []bool{false, false, false, false, true, false, true, true, true, false, true, false, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			a := NewSuperShapeAction(tc.shape, tc.centers)
			require.NotNil(t, a)

			assert.Equal(t, TestGrid, a.Stage())
			assert.Equal(t, SuperRefill, a.Result())
			assert.EqualValues(t, tc.shape, a.shape)
			assert.EqualValues(t, tc.centers, a.centers)

			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			spin.indexes = tc.indexes

			a2 := a.Triggered(spin)
			if tc.want {
				require.NotNil(t, a2)
				assert.Equal(t, tc.symbol, spin.superSymbol)
				assert.EqualValues(t, tc.sticky, spin.sticky)
			} else {
				require.Nil(t, a2)
			}

			a3 := a.Triggered(spin)
			require.Nil(t, a3)
		})
	}
}

func TestGridShape_Retriggered(t *testing.T) {
	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	testCases := []struct {
		name     string
		shape    GridOffsets
		centers  GridOffsets
		indexes1 utils.Indexes
		indexes2 utils.Indexes
		want     bool
		symbol   utils.Index
		sticky   []bool
	}{
		{
			name:     "SuperX, no second hit",
			shape:    superX,
			centers:  centersX,
			indexes1: utils.Indexes{1, 2, 3, 4, 5, 4, 6, 4, 7, 4, 2, 4, 1, 3, 5},
			indexes2: utils.Indexes{7, 6, 5, 4, 5, 4, 6, 4, 7, 4, 2, 4, 5, 6, 7},
			want:     false,
			symbol:   4,
			sticky:   []bool{false, false, false, true, false, true, false, true, false, true, false, true, false, false, false},
		},
		{
			name:     "SuperX, 1 second hit",
			shape:    superX,
			centers:  centersX,
			indexes1: utils.Indexes{1, 2, 3, 4, 5, 4, 6, 4, 7, 4, 2, 4, 1, 3, 5},
			indexes2: utils.Indexes{7, 6, 5, 4, 5, 4, 6, 4, 7, 4, 2, 4, 5, 6, 4},
			want:     true,
			symbol:   4,
			sticky:   []bool{false, false, false, true, false, true, false, true, false, true, false, true, false, false, true},
		},
		{
			name:     "SuperX, 2 second hits",
			shape:    superX,
			centers:  centersX,
			indexes1: utils.Indexes{1, 2, 3, 4, 5, 4, 6, 4, 7, 4, 2, 4, 1, 3, 5},
			indexes2: utils.Indexes{7, 4, 5, 4, 5, 4, 6, 4, 7, 4, 2, 4, 5, 6, 4},
			want:     true,
			symbol:   4,
			sticky:   []bool{false, true, false, true, false, true, false, true, false, true, false, true, false, false, true},
		},
		{
			name:     "SuperX with extra, 3 second hits",
			shape:    superX,
			centers:  centersX,
			indexes1: utils.Indexes{1, 2, 3, 4, 5, 4, 6, 4, 7, 4, 2, 4, 1, 3, 1},
			indexes2: utils.Indexes{7, 4, 5, 4, 5, 4, 4, 4, 7, 4, 2, 4, 5, 4, 4},
			want:     true,
			symbol:   4,
			sticky:   []bool{false, true, false, true, false, true, true, true, false, true, false, true, false, true, true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			a := NewSuperShapeAction(tc.shape, tc.centers)
			require.NotNil(t, a)

			spin := AcquireSpin(slots, prng)
			defer spin.Release()
			spin.indexes = tc.indexes1

			a2 := a.Triggered(spin)
			require.NotNil(t, a2)

			spin.indexes = tc.indexes2
			a3 := a.Triggered(spin)
			if tc.want {
				require.NotNil(t, a3)
				assert.Equal(t, tc.symbol, spin.superSymbol)
				assert.EqualValues(t, tc.sticky, spin.sticky)
			} else {
				require.Nil(t, a3)
			}
		})
	}
}

func TestNewFakeSpinAction(t *testing.T) {
	slots := NewSlots(Grid(6, 4), WithMask(2, 3, 4, 4, 3, 2), WithSymbols(setF1))

	testCases := []struct {
		name    string
		chance  float64
		fake    FakeSpin
		indexes utils.Indexes
		want    bool
	}{
		{
			name:    "just grid",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 9, 9, 0, 0}},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 5, 0, 0},
			want:    true,
		},
		{
			name:    "no replace",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 9, 9, 9, 9, 8, 8, 8, 0, 9, 9, 0, 0}, ReplaceSymbols: utils.Indexes{8}, ReplaceReels: utils.UInt8s{2}},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 5, 0, 0},
			want:    true,
		},
		{
			name:    "replace 1",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 9, 9, 0, 0}, ReplaceSymbols: utils.Indexes{8}, ReplaceReels: utils.UInt8s{2}},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 8, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 5, 0, 0},
			want:    true,
		},
		{
			name:    "replace 3",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 9, 9, 9, 9, 8, 8, 8, 0, 9, 9, 0, 0}, ReplaceSymbols: utils.Indexes{8, 9}, ReplaceReels: utils.UInt8s{2}},
			indexes: utils.Indexes{9, 9, 0, 0, 9, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 5, 0, 0},
			want:    true,
		},
		{
			name:    "grid match (1)",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 9, 9, 0, 0}, MatchReels: utils.UInt8s{1, 6}},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 9, 9, 0, 0},
			want:    true,
		},
		{
			name:    "grid match (2)",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 0, 9, 0, 0}, MatchReels: utils.UInt8s{1, 6}},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 9, 0, 0},
			want:    true,
		},
		{
			name:    "grid match (3)",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 0, 0, 0, 0}, MatchReels: utils.UInt8s{1}},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 5, 0, 0},
			want:    true,
		},
		{
			name:    "grid no match",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{9, 9, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 9, 9, 0, 0}, MatchReels: utils.UInt8s{1, 6}},
			indexes: utils.Indexes{9, 9, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 9, 5, 0, 0},
		},
		{
			name:    "grid inverse match",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 0, 0, 0, 0}, MatchReels: utils.UInt8s{1, 6}, MatchInverse: true, MatchSymbol: 9},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 5, 0, 0},
			want:    true,
		},
		{
			name:    "grid inverse no match",
			chance:  0.5,
			fake:    FakeSpin{Indexes: utils.Indexes{0, 0, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 0, 0, 0, 0}, MatchReels: utils.UInt8s{1, 6}, MatchInverse: true, MatchSymbol: 9},
			indexes: utils.Indexes{1, 2, 0, 0, 3, 4, 5, 0, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 0, 4, 9, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			a := NewFakeSpinAction(tc.chance, tc.fake)
			require.NotNil(t, a)

			spin := AcquireSpin(slots, prng)
			defer spin.Release()

			var counts int

			for ix := 0; ix < 1000; ix++ {
				spin.indexes = tc.indexes
				if a.Triggered(spin) != nil {
					counts++

					for iy := range tc.fake.Indexes {
						if id := tc.fake.Indexes[iy]; id > 0 {
							assert.EqualValues(t, id, spin.indexes[iy])
						}
					}

					if len(tc.fake.ReplaceSymbols) > 0 {
						for iy := range tc.fake.ReplaceReels {
							reel := int(tc.fake.ReplaceReels[iy] - 1)
							offset := reel * spin.rowCount
							max := offset + spin.rowCount
							for ; offset < max; offset++ {
								for iz := range tc.fake.ReplaceSymbols {
									assert.NotEqual(t, tc.fake.ReplaceSymbols[iz], spin.indexes[offset])
								}
							}
						}
					}
				}
			}

			if tc.want {
				assert.NotZero(t, counts)
			} else {
				assert.Zero(t, counts)
			}
		})
	}
}

func TestPurgeGridActions(t *testing.T) {
	t1 := NewSuperShapeAction(superX, centersX)
	t2 := NewSuperShapeAction(superPlus, centersPlus)
	t3 := NewSuperShapeAction(superPlus, centersPlus)
	t4 := NewSuperShapeAction(superX, centersX)
	t5 := NewSuperShapeAction(superPlus, centersPlus)
	t6 := NewSuperShapeAction(superX, centersX)
	t7 := NewSuperShapeAction(superX, centersX)

	testCases := []struct {
		name    string
		in      GridActions
		cap     int
		wantCap int
	}{
		{"nil", nil, 5, 5},
		{"empty", GridActions{}, 5, 5},
		{"short", GridActions{t1, t2, t3}, 5, 5},
		{"exact", GridActions{t1, t2, t3, t4, t5}, 5, 5},
		{"long", GridActions{t1, t2, t3, t4, t5, t6, t7}, 5, 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := PurgeGridActions(tc.in, tc.cap)
			require.NotNil(t, out)
			assert.Equal(t, tc.wantCap, cap(out))
			assert.Zero(t, len(out))
		})
	}
}
