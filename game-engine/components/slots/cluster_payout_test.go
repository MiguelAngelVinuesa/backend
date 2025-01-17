package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewClusterPayouts(t *testing.T) {
	testCases := []struct {
		name  string
		reels uint8
		rows  uint8
		mask  []uint8
		kind  ConnectionKind
		conn  []utils.UInt8s
		want  []utils.UInt8s
	}{
		{
			name:  "3x3",
			reels: 3,
			rows:  3,
			kind:  Rectangular,
			want: []utils.UInt8s{
				{1, 3},       // 0
				{0, 2, 4},    // 1
				{1, 5},       // 2
				{0, 4, 6},    // 3
				{3, 1, 5, 7}, // 4
				{4, 2, 8},    // 5
				{3, 7},       // 6
				{6, 4, 8},    // 7
				{7, 5},       // 8
			},
		},
		{
			name:  "5x3",
			reels: 5,
			rows:  3,
			kind:  Rectangular,
			want: []utils.UInt8s{
				{1, 3},         // 0
				{0, 2, 4},      // 1
				{1, 5},         // 2
				{0, 4, 6},      // 3
				{3, 1, 5, 7},   // 4
				{4, 2, 8},      // 5
				{3, 7, 9},      // 6
				{6, 4, 8, 10},  // 7
				{7, 5, 11},     // 8
				{6, 10, 12},    // 9
				{9, 7, 11, 13}, // 10
				{10, 8, 14},    // 11
				{9, 13},        // 12
				{12, 10, 14},   // 13
				{13, 11},       // 14
			},
		},
		{
			name:  "5x5",
			reels: 5,
			rows:  5,
			kind:  Rectangular,
			want: []utils.UInt8s{
				{1, 5},           // 0
				{0, 2, 6},        // 1
				{1, 3, 7},        // 2
				{2, 4, 8},        // 3
				{3, 9},           // 4
				{0, 6, 10},       // 5
				{5, 1, 7, 11},    // 6
				{6, 2, 8, 12},    // 7
				{7, 3, 9, 13},    // 8
				{8, 4, 14},       // 9
				{5, 11, 15},      // 10
				{10, 6, 12, 16},  // 11
				{11, 7, 13, 17},  // 12
				{12, 8, 14, 18},  // 13
				{13, 9, 19},      // 14
				{10, 16, 20},     // 15
				{15, 11, 17, 21}, // 16
				{16, 12, 18, 22}, // 17
				{17, 13, 19, 23}, // 18
				{18, 14, 24},     // 19
				{15, 21},         // 20
				{20, 16, 22},     // 21
				{21, 17, 23},     // 22
				{22, 18, 24},     // 23
				{23, 19},         // 24
			},
		},
		{
			name:  "2|3|2",
			reels: 3,
			rows:  3,
			mask:  []uint8{2, 3, 2},
			kind:  Hexagonal,
			want: []utils.UInt8s{
				{1, 3, 4},          // 0
				{0, 4, 5},          // 1
				{},                 // 2
				{4, 0, 6},          // 3
				{3, 5, 0, 1, 6, 7}, // 4
				{4, 1, 7},          // 5
				{7, 3, 4},          // 6
				{6, 4, 5},          // 7
				{},                 // 8
			},
		},
		{
			name:  "2|3|4|3|2",
			reels: 5,
			rows:  4,
			mask:  []uint8{2, 3, 4, 3, 2},
			kind:  Hexagonal,
			want: []utils.UInt8s{
				{1, 4, 5},               // 0
				{0, 5, 6},               // 1
				{},                      // 2
				{},                      // 3
				{5, 0, 8, 9},            // 4
				{4, 6, 0, 1, 9, 10},     // 5
				{5, 1, 10, 11},          // 6
				{},                      // 7
				{9, 4, 12},              // 8
				{8, 10, 4, 5, 12, 13},   // 9
				{9, 11, 5, 6, 13, 14},   // 10
				{10, 6, 14},             // 11
				{13, 8, 9, 16},          // 12
				{12, 14, 9, 10, 16, 17}, // 13
				{13, 10, 11, 17},        // 14
				{},                      // 15
				{17, 12, 13},            // 16
				{16, 13, 14},            // 17
				{},                      // 18
				{},                      // 19
			},
		},
		{
			name:  "4|5|6|7|6|5|4",
			reels: 7,
			rows:  7,
			mask:  []uint8{4, 5, 6, 7, 6, 5, 4},
			kind:  Hexagonal,
			want: []utils.UInt8s{
				{1, 7, 8},                // 0
				{0, 2, 8, 9},             // 1
				{1, 3, 9, 10},            // 2
				{2, 10, 11},              // 3
				{},                       // 4
				{},                       // 5
				{},                       // 6
				{8, 0, 14, 15},           // 7
				{7, 9, 0, 1, 15, 16},     // 8
				{8, 10, 1, 2, 16, 17},    // 9
				{9, 11, 2, 3, 17, 18},    // 10
				{10, 3, 18, 19},          // 11
				{},                       // 12
				{},                       // 13
				{15, 7, 21, 22},          // 14
				{14, 16, 7, 8, 22, 23},   // 15
				{15, 17, 8, 9, 23, 24},   // 16
				{16, 18, 9, 10, 24, 25},  // 17
				{17, 19, 10, 11, 25, 26}, // 18
				{18, 11, 26, 27},         // 19
				{},                       // 20
				{22, 14, 28},             // 21
				{21, 23, 14, 15, 28, 29}, // 22
				{22, 24, 15, 16, 29, 30}, // 23
				{23, 25, 16, 17, 30, 31}, // 24
				{24, 26, 17, 18, 31, 32}, // 25
				{25, 27, 18, 19, 32, 33}, // 26
				{26, 19, 33},             // 27
				{29, 21, 22, 35},         // 28
				{28, 30, 22, 23, 35, 36}, // 29
				{29, 31, 23, 24, 36, 37}, // 30
				{30, 32, 24, 25, 37, 38}, // 31
				{31, 33, 25, 26, 38, 39}, // 32
				{32, 26, 27, 39},         // 33
				{},                       // 34
				{36, 28, 29, 42},         // 35
				{35, 37, 29, 30, 42, 43}, // 36
				{36, 38, 30, 31, 43, 44}, // 37
				{37, 39, 31, 32, 44, 45}, // 38
				{38, 32, 33, 45},         // 39
				{},                       // 40
				{},                       // 41
				{43, 35, 36},             // 42
				{42, 44, 36, 37},         // 43
				{43, 45, 37, 38},         // 44
				{44, 38, 39},             // 45
				{},                       // 46
				{},                       // 47
				{},                       // 48
			},
		},
		{
			name:  "predefined",
			reels: 5,
			rows:  3,
			conn: []utils.UInt8s{
				{1, 3},       // 0
				{0, 2, 4},    // 1
				{1, 5},       // 2
				{0, 4, 6},    // 3
				{3, 1, 5, 7}, // 4
				{4, 2, 8},    // 5
				{3, 7},       // 6
				{6, 4, 8},    // 7
				{7, 5},       // 8
			},
			want: []utils.UInt8s{
				{1, 3},       // 0
				{0, 2, 4},    // 1
				{1, 5},       // 2
				{0, 4, 6},    // 3
				{3, 1, 5, 7}, // 4
				{4, 2, 8},    // 5
				{3, 7},       // 6
				{6, 4, 8},    // 7
				{7, 5},       // 8
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := make([]ClusterPayoutsOption, 0, 4)
			if len(tc.mask) > 0 {
				opts = append(opts, ClusterGridMask(tc.mask, tc.kind))
			}
			if len(tc.conn) > 0 {
				opts = append(opts, ClusterConnections(tc.conn))
			}

			c := NewClusterPayouts(tc.reels, tc.rows, opts...)
			require.NotNil(t, c)

			assert.Equal(t, tc.reels, c.reels)
			assert.Equal(t, tc.rows, c.rows)
			assert.EqualValues(t, tc.mask, c.mask)
			assert.Equal(t, tc.kind, c.kind)
			assert.EqualValues(t, tc.want, c.connections)
		})
	}
}

func TestClusterPayouts_FindRectangular3x3(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.UInt8s
		payouts int
	}{
		{
			name:    "no hit",
			indexes: utils.Indexes{1, 2, 1, 2, 1, 2, 1, 2, 1},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "2x 11 hor",
			indexes: utils.Indexes{1, 2, 11, 1, 3, 11, 4, 2, 1},
			want:    utils.UInt8s{0, 0, 1, 0, 0, 1, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "2x 11 ver",
			indexes: utils.Indexes{1, 2, 1, 11, 11, 4, 4, 2, 1},
			want:    utils.UInt8s{0, 0, 0, 1, 1, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "2x 11 2 times",
			indexes: utils.Indexes{11, 11, 1, 2, 3, 11, 4, 5, 11},
			want:    utils.UInt8s{1, 1, 0, 0, 0, 1, 0, 0, 1},
			payouts: 2,
		},
		{
			name:    "3x 1 hor",
			indexes: utils.Indexes{1, 2, 1, 1, 3, 4, 1, 2, 1},
			want:    utils.UInt8s{1, 0, 0, 1, 0, 0, 1, 0, 0},
			payouts: 1,
		},
		{
			name:    "3x 1 ver",
			indexes: utils.Indexes{1, 1, 1, 2, 3, 2, 1, 2, 1},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "3x 1 & 2x 11 hor",
			indexes: utils.Indexes{1, 2, 11, 1, 11, 2, 1, 11, 1},
			want:    utils.UInt8s{1, 0, 0, 1, 1, 0, 1, 1, 0},
			payouts: 2,
		},
		{
			name:    "3x 4 top-left corner",
			indexes: utils.Indexes{4, 4, 1, 4, 3, 4, 1, 2, 4},
			want:    utils.UInt8s{1, 1, 0, 1, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "3x 4 bottom-left corner",
			indexes: utils.Indexes{1, 4, 4, 2, 3, 4, 4, 4, 2},
			want:    utils.UInt8s{0, 1, 1, 0, 0, 1, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "3x 4 top-right corner",
			indexes: utils.Indexes{1, 4, 1, 4, 3, 4, 4, 4, 1},
			want:    utils.UInt8s{0, 0, 0, 1, 0, 0, 1, 1, 0},
			payouts: 1,
		},
		{
			name:    "3x 4 bottom-right corner",
			indexes: utils.Indexes{1, 4, 1, 2, 3, 4, 1, 4, 4},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 1, 0, 1, 1},
			payouts: 1,
		},
		{
			name:    "3x 1&3 hor",
			indexes: utils.Indexes{1, 2, 3, 1, 4, 3, 1, 2, 3},
			want:    utils.UInt8s{1, 0, 1, 1, 0, 1, 1, 0, 1},
			payouts: 2,
		},
		{
			name:    "3x 1&2&3 hor",
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3},
			want:    utils.UInt8s{1, 1, 1, 1, 1, 1, 1, 1, 1},
			payouts: 3,
		},
		{
			name:    "3x 1&3 ver",
			indexes: utils.Indexes{1, 1, 1, 2, 4, 2, 3, 3, 3},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 0, 1, 1, 1},
			payouts: 2,
		},
		{
			name:    "3x 1&2&3 ver",
			indexes: utils.Indexes{1, 1, 1, 2, 2, 2, 3, 3, 3},
			want:    utils.UInt8s{1, 1, 1, 1, 1, 1, 1, 1, 1},
			payouts: 3,
		},
		{
			name:    "4x 1 box",
			indexes: utils.Indexes{2, 1, 1, 2, 1, 1, 1, 2, 3},
			want:    utils.UInt8s{0, 1, 1, 0, 1, 1, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "5x 1 L",
			indexes: utils.Indexes{1, 1, 1, 2, 3, 1, 1, 2, 1},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 0, 0, 1},
			payouts: 1,
		},
		{
			name:    "5x 1 reverse L",
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 1, 1},
			want:    utils.UInt8s{1, 0, 0, 1, 0, 0, 1, 1, 1},
			payouts: 1,
		},
		{
			name:    "5x 1 plus",
			indexes: utils.Indexes{2, 1, 2, 1, 1, 1, 2, 1, 2},
			want:    utils.UInt8s{0, 1, 0, 1, 1, 1, 0, 1, 0},
			payouts: 1,
		},
		{
			name:    "7x 1 U",
			indexes: utils.Indexes{1, 1, 1, 2, 2, 1, 1, 1, 1},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 1, 1, 1},
			payouts: 1,
		},
		{
			name:    "7x 1 C",
			indexes: utils.Indexes{1, 1, 1, 1, 2, 1, 1, 2, 1},
			want:    utils.UInt8s{1, 1, 1, 1, 0, 1, 1, 0, 1},
			payouts: 1,
		},
		{
			name:    "7x 1 reverse C",
			indexes: utils.Indexes{1, 2, 1, 1, 2, 1, 1, 1, 1},
			want:    utils.UInt8s{1, 0, 1, 1, 0, 1, 1, 1, 1},
			payouts: 1,
		},
		{
			name:    "7x 1 H",
			indexes: utils.Indexes{1, 1, 1, 2, 1, 2, 1, 1, 1},
			want:    utils.UInt8s{1, 1, 1, 0, 1, 0, 1, 1, 1},
			payouts: 1,
		},
		{
			name:    "7x 1 I",
			indexes: utils.Indexes{1, 2, 1, 1, 1, 1, 1, 2, 1},
			want:    utils.UInt8s{1, 0, 1, 1, 1, 1, 1, 0, 1},
			payouts: 1,
		},
		{
			name:    "8x 1 O",
			indexes: utils.Indexes{1, 1, 1, 1, 2, 1, 1, 1, 1},
			want:    utils.UInt8s{1, 1, 1, 1, 0, 1, 1, 1, 1},
			payouts: 1,
		},
		{
			name:    "9x 1 all",
			indexes: utils.Indexes{1, 1, 1, 1, 1, 1, 1, 1, 1},
			want:    utils.UInt8s{1, 1, 1, 1, 1, 1, 1, 1, 1},
			payouts: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			game := NewSlots(Grid(3, 3), WithSymbols(setF1))

			spin := AcquireSpin(game, prng)
			defer spin.Release()

			copy(spin.indexes, tc.indexes)

			c := NewClusterPayouts(3, 3)
			require.NotNil(t, c)

			payouts := c.Find(spin, nil)
			assert.EqualValues(t, tc.want, spin.payouts)
			assert.Equal(t, tc.payouts, payouts)
		})
	}
}

func TestClusterPayouts_FindRectangular5x3(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.UInt8s
		payouts int
	}{
		{
			name:    "no hit",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "3x 4",
			indexes: utils.Indexes{4, 2, 3, 4, 4, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			want:    utils.UInt8s{1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "few clusters",
			indexes: utils.Indexes{2, 2, 7, 4, 4, 4, 4, 4, 4, 4, 7, 5, 5, 5, 5},
			want:    utils.UInt8s{0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 1},
			payouts: 2,
		},
		{
			name:    "11x 4 snake",
			indexes: utils.Indexes{7, 7, 7, 4, 5, 7, 7, 7, 7, 7, 5, 6, 7, 7, 7},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1},
			payouts: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			game := NewSlots(Grid(5, 3), WithSymbols(setF1))

			spin := AcquireSpin(game, prng)
			defer spin.Release()

			copy(spin.indexes, tc.indexes)

			c := NewClusterPayouts(5, 3)
			require.NotNil(t, c)

			payouts := c.Find(spin, nil)
			assert.EqualValues(t, tc.want, spin.payouts)
			assert.Equal(t, tc.payouts, payouts)
		})
	}
}

func TestClusterPayouts_FindRectangular5x5(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.UInt8s
		payouts int
	}{
		{
			name:    "no hit",
			indexes: utils.Indexes{1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "3x 4",
			indexes: utils.Indexes{1, 2, 1, 4, 1, 2, 1, 2, 4, 2, 1, 2, 1, 4, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1},
			want:    utils.UInt8s{0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "few clusters",
			indexes: utils.Indexes{1, 2, 1, 4, 1, 6, 6, 6, 4, 2, 1, 6, 6, 4, 1, 2, 1, 2, 1, 5, 1, 2, 1, 5, 5},
			want:    utils.UInt8s{0, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1},
			payouts: 3,
		},
		{
			name:    "15x 4 snake",
			indexes: utils.Indexes{1, 4, 4, 4, 1, 4, 4, 2, 4, 4, 4, 2, 1, 2, 4, 4, 1, 2, 4, 4, 4, 4, 1, 4, 1},
			want:    utils.UInt8s{0, 1, 1, 1, 0, 1, 1, 0, 1, 1, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 0, 1, 0},
			payouts: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			game := NewSlots(Grid(5, 5), WithSymbols(setF1))

			spin := AcquireSpin(game, prng)
			defer spin.Release()

			copy(spin.indexes, tc.indexes)

			c := NewClusterPayouts(5, 5)
			require.NotNil(t, c)

			payouts := c.Find(spin, nil)
			assert.EqualValues(t, tc.want, spin.payouts)
			assert.Equal(t, tc.payouts, payouts)
		})
	}
}

func TestClusterPayouts_FindHexagonal34543(t *testing.T) {
	mask := []uint8{3, 4, 5, 4, 3}
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    utils.UInt8s
		payouts int
	}{
		{
			name:    "no hit",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 3, 4, 3, 4, 0, 1, 2, 1, 2, 1, 3, 4, 3, 4, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "no hit 2",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 1, 3, 3, 1, 0, 2, 4, 1, 4, 2, 1, 3, 3, 1, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:    "3x 1",
			indexes: utils.Indexes{1, 1, 2, 0, 0, 3, 1, 3, 3, 0, 2, 4, 2, 4, 2, 1, 3, 3, 1, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "3x 1 twice",
			indexes: utils.Indexes{1, 2, 2, 0, 0, 1, 1, 3, 3, 0, 2, 4, 2, 4, 1, 1, 3, 3, 1, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0},
			payouts: 2,
		},
		{
			name:    "5x 1 hor",
			indexes: utils.Indexes{1, 1, 2, 0, 0, 3, 1, 3, 3, 0, 2, 4, 1, 4, 2, 3, 1, 3, 1, 0, 4, 2, 4, 0, 0},
			want:    utils.UInt8s{1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "5x 1 twice hor",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 1, 3, 3, 1, 0, 1, 4, 2, 4, 1, 1, 3, 3, 1, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{1, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0},
			payouts: 2,
		},
		{
			name:    "5x 1 ver",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 3, 4, 3, 4, 0, 3, 2, 5, 2, 3, 1, 4, 3, 1, 0, 1, 1, 1, 0, 0},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0},
			payouts: 1,
		},
		{
			name:    "5x 1 twice ver",
			indexes: utils.Indexes{1, 1, 1, 0, 0, 1, 2, 2, 1, 0, 3, 4, 5, 4, 3, 1, 2, 2, 1, 0, 1, 1, 1, 0, 0},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0},
			payouts: 2,
		},
		{
			name:    "12x 3 small hexagon",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 3, 3, 3, 3, 0, 3, 4, 5, 4, 3, 3, 3, 3, 3, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "12x 1 big hexagon",
			indexes: utils.Indexes{1, 1, 1, 0, 0, 1, 2, 2, 1, 0, 1, 4, 5, 4, 1, 1, 2, 2, 1, 0, 1, 1, 1, 0, 0},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0},
			payouts: 1,
		},
		{
			name:    "7x 3 X",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 3, 4, 4, 3, 0, 1, 3, 3, 3, 1, 3, 4, 4, 3, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "12x 1 snake",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 1, 4, 1, 4, 0, 1, 2, 1, 1, 2, 1, 4, 3, 1, 0, 1, 1, 1, 0, 0},
			want:    utils.UInt8s{1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 0},
			payouts: 1,
		},
		{
			name:    "3 clusters",
			indexes: utils.Indexes{1, 2, 1, 0, 0, 3, 1, 1, 4, 0, 3, 3, 5, 6, 5, 4, 4, 2, 2, 0, 1, 2, 1, 0, 0},
			want:    utils.UInt8s{1, 0, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0},
			payouts: 3,
		},
		{
			name:    "5 clusters",
			indexes: utils.Indexes{1, 2, 2, 0, 0, 3, 3, 2, 4, 0, 3, 3, 1, 1, 1, 3, 4, 4, 5, 0, 4, 5, 5, 0, 0},
			want:    utils.UInt8s{0, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0},
			payouts: 5,
		},
		{
			name:    "5x 1 + 1 wild",
			indexes: utils.Indexes{1, 1, 1, 0, 0, 1, 9, 2, 1, 0, 3, 4, 5, 4, 3, 1, 2, 4, 1, 0, 2, 3, 1, 0, 0},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
		{
			name:    "5x 1 + 3 wilds",
			indexes: utils.Indexes{9, 9, 1, 0, 0, 1, 1, 1, 9, 0, 3, 4, 1, 5, 6, 1, 2, 4, 1, 0, 7, 3, 6, 0, 0},
			want:    utils.UInt8s{1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			payouts: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			game := NewSlots(Grid(5, 5), WithSymbols(setF1), WithMask(mask...))

			spin := AcquireSpin(game, prng)
			defer spin.Release()

			copy(spin.indexes, tc.indexes)

			c := NewClusterPayouts(5, 5, ClusterGridMask(mask, Hexagonal))
			require.NotNil(t, c)

			payouts := c.Find(spin, nil)
			assert.EqualValues(t, tc.want, spin.payouts)
			assert.Equal(t, tc.payouts, payouts)
		})
	}
}

func TestClusterPayouts_RemovePayouts(t *testing.T) {
	reels := uint8(5)
	rows := uint8(5)
	mask := utils.UInt8s{3, 4, 5, 4, 3}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	set := NewSymbolSet(
		NewSymbol(1, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(2, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(3, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(4, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(5, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(6, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(7, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(8, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(9, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(10, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(11, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(12, WithPayouts(0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)),
		NewSymbol(13, WithKind(Wild)),
		NewSymbol(14, WithPayouts(0, 0, 1, 2, 3), WithKind(Scatter)),
	)

	slots := NewSlots(Grid(reels, rows), WithMask(mask...), WithSymbols(set), CascadingReels(true))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	action := NewClusterPayouts(reels, rows, ClusterGridMask(mask, Hexagonal))

	testCases := []struct {
		name     string
		indexes  utils.Indexes
		clusters bool
	}{
		{
			name:    "none",
			indexes: utils.Indexes{1, 2, 3, 0, 0, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 4, 5, 6, 7, 0, 1, 2, 3, 0, 0},
		},
		{
			name:     "6x4",
			indexes:  utils.Indexes{1, 2, 3, 0, 0, 4, 4, 6, 7, 0, 4, 4, 4, 4, 5, 4, 5, 6, 7, 0, 1, 2, 3, 0, 0},
			clusters: true,
		},
		{
			name:     "4x7",
			indexes:  utils.Indexes{1, 2, 3, 0, 0, 4, 5, 6, 7, 0, 1, 2, 3, 4, 7, 4, 5, 6, 7, 0, 1, 2, 7, 0, 0},
			clusters: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			copy(spin.indexes, tc.indexes)

			if tc.clusters {
				n := action.Find(spin, nil)
				assert.NotZero(t, n)
			}

			action.RemovePayouts(spin)

			n := action.Find(spin, nil)
			assert.Zero(t, n)
		})
	}
}

func BenchmarkClusterPayouts_FindRectangular3x3none(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(3, 3), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9})

	c := NewClusterPayouts(3, 3)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular3x3one(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(3, 3), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 1, 4, 5, 6, 1, 8, 9})

	c := NewClusterPayouts(3, 3)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular3x3two(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(3, 3), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 1, 4, 4, 1, 8, 4})

	c := NewClusterPayouts(3, 3)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular3x3three(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(3, 3), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 1, 2, 1, 3, 2, 3, 3, 2})

	c := NewClusterPayouts(3, 3)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular5x5none(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(5, 5), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5})

	c := NewClusterPayouts(5, 5)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular5x5one(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(5, 5), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 5, 6, 7, 3, 9, 10, 1, 2, 3, 3, 5, 6, 3, 8, 9, 10, 1, 2, 3, 4, 5})

	c := NewClusterPayouts(5, 5)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular5x5two(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(5, 5), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 5, 6, 7, 3, 9, 10, 1, 2, 3, 3, 4, 6, 3, 8, 9, 4, 1, 2, 3, 4, 4})

	c := NewClusterPayouts(5, 5)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindRectangular5x5four(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	game := NewSlots(Grid(5, 5), WithSymbols(setF1))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 1, 1, 4, 5, 6, 1, 3, 2, 2, 1, 2, 3, 3, 2, 6, 3, 8, 9, 4, 1, 2, 3, 4, 4})

	c := NewClusterPayouts(5, 5)

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal34543none(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	mask := utils.UInt8s{3, 4, 5, 4, 3}
	game := NewSlots(Grid(5, 5), WithSymbols(setF1), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 0, 0, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 4, 5, 6, 7, 0, 1, 2, 3, 0, 0})

	c := NewClusterPayouts(5, 5, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal34543one(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	mask := utils.UInt8s{3, 4, 5, 4, 3}
	game := NewSlots(Grid(5, 5), WithSymbols(setF1), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 0, 0, 4, 5, 6, 7, 0, 1, 2, 3, 4, 5, 4, 2, 2, 7, 0, 1, 2, 3, 0, 0})

	c := NewClusterPayouts(5, 5, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal34543two(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	mask := utils.UInt8s{3, 4, 5, 4, 3}
	game := NewSlots(Grid(5, 5), WithSymbols(setF1), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 0, 0, 4, 4, 6, 7, 0, 1, 4, 4, 4, 5, 4, 2, 2, 7, 0, 1, 2, 3, 0, 0})

	c := NewClusterPayouts(5, 5, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal34543five(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	mask := utils.UInt8s{3, 4, 5, 4, 3}
	game := NewSlots(Grid(5, 5), WithSymbols(setF1), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 0, 0, 4, 4, 3, 3, 0, 1, 1, 4, 4, 5, 1, 2, 2, 5, 0, 1, 2, 5, 0, 0})

	c := NewClusterPayouts(5, 5, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal4567654none(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	payouts := WithPayouts(0, 0, 0.25, 0.5, 1, 2, 4, 5, 6, 8, 10, 12, 15, 17, 20, 22, 25, 25, 25, 30, 30, 30, 30, 30, 40, 40, 40, 40, 40, 50, 50, 50, 50, 50, 75, 100, 200)
	weights := WithWeights(20, 20, 20, 20, 20, 20, 20)
	s01 := NewSymbol(1, payouts, weights)
	s02 := NewSymbol(2, payouts, weights)
	s03 := NewSymbol(3, payouts, weights)
	s04 := NewSymbol(4, payouts, weights)
	s05 := NewSymbol(5, payouts, weights)
	s06 := NewSymbol(6, payouts, weights)
	s07 := NewSymbol(7, payouts, weights)
	s08 := NewSymbol(8, payouts, weights)
	s09 := NewSymbol(9, payouts, weights)
	s10 := NewSymbol(10, payouts, weights)
	s11 := NewSymbol(11, payouts, weights)
	s12 := NewSymbol(12, payouts, weights)

	symbols := NewSymbolSet(s01, s02, s03, s04, s05, s06, s07, s08, s09, s10, s11, s12)

	mask := utils.UInt8s{4, 5, 6, 7, 6, 5, 4}
	game := NewSlots(Grid(7, 7), WithSymbols(symbols), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 0, 0, 0, 5, 6, 7, 8, 9, 0, 0, 1, 2, 3, 4, 5, 6, 0, 5, 6, 7, 8, 9, 10, 11, 1, 2, 3, 4, 5, 6, 0, 5, 6, 7, 8, 9, 0, 0, 1, 2, 3, 4, 0, 0, 0})

	c := NewClusterPayouts(7, 7, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal4567654one(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	payouts := WithPayouts(0, 0, 0.25, 0.5, 1, 2, 4, 5, 6, 8, 10, 12, 15, 17, 20, 22, 25, 25, 25, 30, 30, 30, 30, 30, 40, 40, 40, 40, 40, 50, 50, 50, 50, 50, 75, 100, 200)
	weights := WithWeights(20, 20, 20, 20, 20, 20, 20)
	s01 := NewSymbol(1, payouts, weights)
	s02 := NewSymbol(2, payouts, weights)
	s03 := NewSymbol(3, payouts, weights)
	s04 := NewSymbol(4, payouts, weights)
	s05 := NewSymbol(5, payouts, weights)
	s06 := NewSymbol(6, payouts, weights)
	s07 := NewSymbol(7, payouts, weights)
	s08 := NewSymbol(8, payouts, weights)
	s09 := NewSymbol(9, payouts, weights)
	s10 := NewSymbol(10, payouts, weights)
	s11 := NewSymbol(11, payouts, weights)
	s12 := NewSymbol(12, payouts, weights)

	symbols := NewSymbolSet(s01, s02, s03, s04, s05, s06, s07, s08, s09, s10, s11, s12)
	mask := utils.UInt8s{4, 5, 6, 7, 6, 5, 4}
	game := NewSlots(Grid(7, 7), WithSymbols(symbols), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 0, 0, 0, 5, 6, 7, 8, 9, 0, 0, 5, 2, 3, 4, 5, 6, 0, 5, 6, 7, 8, 9, 10, 11, 1, 2, 3, 4, 5, 6, 0, 5, 6, 7, 8, 9, 0, 0, 1, 2, 3, 4, 0, 0, 0})

	c := NewClusterPayouts(7, 7, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal4567654two(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	payouts := WithPayouts(0, 0, 0.25, 0.5, 1, 2, 4, 5, 6, 8, 10, 12, 15, 17, 20, 22, 25, 25, 25, 30, 30, 30, 30, 30, 40, 40, 40, 40, 40, 50, 50, 50, 50, 50, 75, 100, 200)
	weights := WithWeights(20, 20, 20, 20, 20, 20, 20)
	s01 := NewSymbol(1, payouts, weights)
	s02 := NewSymbol(2, payouts, weights)
	s03 := NewSymbol(3, payouts, weights)
	s04 := NewSymbol(4, payouts, weights)
	s05 := NewSymbol(5, payouts, weights)
	s06 := NewSymbol(6, payouts, weights)
	s07 := NewSymbol(7, payouts, weights)
	s08 := NewSymbol(8, payouts, weights)
	s09 := NewSymbol(9, payouts, weights)
	s10 := NewSymbol(10, payouts, weights)
	s11 := NewSymbol(11, payouts, weights)
	s12 := NewSymbol(12, payouts, weights)

	symbols := NewSymbolSet(s01, s02, s03, s04, s05, s06, s07, s08, s09, s10, s11, s12)
	mask := utils.UInt8s{4, 5, 6, 7, 6, 5, 4}
	game := NewSlots(Grid(7, 7), WithSymbols(symbols), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 0, 0, 0, 5, 6, 7, 8, 9, 0, 0, 5, 2, 3, 4, 5, 6, 0, 5, 6, 7, 8, 9, 10, 11, 1, 2, 3, 4, 5, 6, 0, 5, 6, 7, 4, 4, 0, 0, 1, 2, 3, 4, 0, 0, 0})

	c := NewClusterPayouts(7, 7, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal4567654four(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	payouts := WithPayouts(0, 0, 0.25, 0.5, 1, 2, 4, 5, 6, 8, 10, 12, 15, 17, 20, 22, 25, 25, 25, 30, 30, 30, 30, 30, 40, 40, 40, 40, 40, 50, 50, 50, 50, 50, 75, 100, 200)
	weights := WithWeights(20, 20, 20, 20, 20, 20, 20)
	s01 := NewSymbol(1, payouts, weights)
	s02 := NewSymbol(2, payouts, weights)
	s03 := NewSymbol(3, payouts, weights)
	s04 := NewSymbol(4, payouts, weights)
	s05 := NewSymbol(5, payouts, weights)
	s06 := NewSymbol(6, payouts, weights)
	s07 := NewSymbol(7, payouts, weights)
	s08 := NewSymbol(8, payouts, weights)
	s09 := NewSymbol(9, payouts, weights)
	s10 := NewSymbol(10, payouts, weights)
	s11 := NewSymbol(11, payouts, weights)
	s12 := NewSymbol(12, payouts, weights)

	symbols := NewSymbolSet(s01, s02, s03, s04, s05, s06, s07, s08, s09, s10, s11, s12)
	mask := utils.UInt8s{4, 5, 6, 7, 6, 5, 4}
	game := NewSlots(Grid(7, 7), WithSymbols(symbols), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 0, 0, 0, 5, 2, 2, 8, 9, 0, 0, 5, 2, 3, 4, 4, 4, 0, 5, 6, 7, 8, 9, 10, 11, 1, 2, 3, 4, 5, 6, 0, 5, 6, 7, 4, 4, 0, 0, 1, 2, 3, 4, 0, 0, 0})

	c := NewClusterPayouts(7, 7, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}

func BenchmarkClusterPayouts_FindHexagonal4567654six(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	payouts := WithPayouts(0, 0, 0.25, 0.5, 1, 2, 4, 5, 6, 8, 10, 12, 15, 17, 20, 22, 25, 25, 25, 30, 30, 30, 30, 30, 40, 40, 40, 40, 40, 50, 50, 50, 50, 50, 75, 100, 200)
	weights := WithWeights(20, 20, 20, 20, 20, 20, 20)
	s01 := NewSymbol(1, payouts, weights)
	s02 := NewSymbol(2, payouts, weights)
	s03 := NewSymbol(3, payouts, weights)
	s04 := NewSymbol(4, payouts, weights)
	s05 := NewSymbol(5, payouts, weights)
	s06 := NewSymbol(6, payouts, weights)
	s07 := NewSymbol(7, payouts, weights)
	s08 := NewSymbol(8, payouts, weights)
	s09 := NewSymbol(9, payouts, weights)
	s10 := NewSymbol(10, payouts, weights)
	s11 := NewSymbol(11, payouts, weights)
	s12 := NewSymbol(12, payouts, weights)

	symbols := NewSymbolSet(s01, s02, s03, s04, s05, s06, s07, s08, s09, s10, s11, s12)
	mask := utils.UInt8s{4, 5, 6, 7, 6, 5, 4}
	game := NewSlots(Grid(7, 7), WithSymbols(symbols), WithMask(mask...))

	spin := AcquireSpin(game, prng)
	defer spin.Release()

	copy(spin.indexes, utils.Indexes{1, 2, 3, 4, 0, 0, 0, 5, 2, 2, 8, 9, 0, 0, 5, 2, 3, 4, 4, 4, 0, 5, 6, 7, 8, 9, 10, 11, 1, 2, 8, 8, 5, 4, 0, 5, 6, 7, 4, 4, 0, 0, 2, 1, 1, 1, 0, 0, 0})

	c := NewClusterPayouts(7, 7, ClusterGridMask(mask, Hexagonal))

	for i := 0; i < b.N; i++ {
		c.Find(spin, nil)
	}
}
