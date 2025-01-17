package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSpinResult(t *testing.T) {
	testCases := []struct {
		name    string
		before  utils.Indexes
		locked1 []bool
		json    string
		after   utils.Indexes
		locked2 []bool
		json2   string
		bonus   utils.Index
		json3   string
	}{
		{
			name:   "3x3",
			before: utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2},
			json:   `{"kind":0,"initial":[0,1,2,0,1,2,0,1,2]}`,
		},
		{
			name:   "3x3, bonus",
			before: utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2},
			json:   `{"kind":0,"initial":[0,1,2,0,1,2,0,1,2]}`,
			bonus:  4,
			json3:  `{"kind":0,"initial":[0,1,2,0,1,2,0,1,2],"bonusSymbol":4}`,
		},
		{
			name:    "5x3 expand before",
			before:  utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 9, 9, 9, 0, 1, 2},
			locked1: []bool{false, false, false, true, false},
			json:    `{"kind":0,"initial":[0,1,2,0,1,2,0,1,2,9,9,9,0,1,2],"lockedReels":[4]}`,
		},
		{
			name:    "5x3 expand after",
			before:  utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 0, 1, 9, 0, 1, 2},
			json:    `{"kind":0,"initial":[0,1,2,0,1,2,0,1,2,0,1,9,0,1,2]}`,
			after:   utils.Indexes{0, 1, 2, 0, 1, 2, 0, 1, 2, 9, 9, 9, 0, 1, 2},
			locked2: []bool{false, false, false, true, false},
			json2:   `{"kind":0,"initial":[0,1,2,0,1,2,0,1,2,0,1,9,0,1,2],"afterExpand":[0,1,2,0,1,2,0,1,2,9,9,9,0,1,2],"flagsExpand":[0,0,0,0,0,0,0,0,0,1,1,0,0,0,0],"lockedReels":[4]}`,
		},
		{
			name:    "5x3 2x expand after, bonus",
			before:  utils.Indexes{0, 9, 2, 0, 1, 2, 9, 1, 2, 0, 1, 2, 0, 1, 2},
			json:    `{"kind":0,"initial":[0,9,2,0,1,2,9,1,2,0,1,2,0,1,2]}`,
			after:   utils.Indexes{9, 9, 9, 0, 1, 2, 9, 9, 9, 0, 1, 2, 0, 1, 2},
			locked2: []bool{true, false, true, false, false},
			json2:   `{"kind":0,"initial":[0,9,2,0,1,2,9,1,2,0,1,2,0,1,2],"afterExpand":[9,9,9,0,1,2,9,9,9,0,1,2,0,1,2],"flagsExpand":[1,0,1,0,0,0,0,1,1,0,0,0,0,0,0],"lockedReels":[1,3]}`,
			bonus:   2,
			json3:   `{"kind":0,"initial":[0,9,2,0,1,2,9,1,2,0,1,2,0,1,2],"afterExpand":[9,9,9,0,1,2,9,9,9,0,1,2,0,1,2],"flagsExpand":[1,0,1,0,0,0,0,1,1,0,0,0,0,0,0],"lockedReels":[1,3],"bonusSymbol":2}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := &Spin{indexes: tc.before}
			if len(tc.locked1) > 0 {
				spin.locked = tc.locked1
			}

			r := AcquireSpinResult(spin)
			require.NotNil(t, r)
			defer r.Release()

			assert.EqualValues(t, tc.before, r.Initial())
			assert.Zero(t, len(r.AfterExpand()))
			assert.Zero(t, r.BonusSymbol())

			enc := zjson.AcquireEncoder(128)
			defer enc.Release()
			enc.Object(r)
			assert.Equal(t, tc.json, string(enc.Bytes()))

			if len(tc.after) > 0 {
				spin.indexes = tc.after
				if len(tc.locked2) > 0 {
					spin.locked = tc.locked2
				}

				r.SetAfterExpand(spin)
				assert.EqualValues(t, tc.before, r.Initial())
				assert.EqualValues(t, tc.after, r.AfterExpand())
				assert.Zero(t, r.BonusSymbol())

				enc.Reset()
				enc.Object(r)
				assert.Equal(t, tc.json2, string(enc.Bytes()))
			}

			if tc.bonus > 0 {
				r.SetBonusSymbol(tc.bonus)
				assert.Equal(t, tc.bonus, r.BonusSymbol())

				enc.Reset()
				enc.Object(r)
				assert.Equal(t, tc.json3, string(enc.Bytes()))
			}
		})
	}
}

func TestNewSpinResultFromData(t *testing.T) {
	testCases := []struct {
		name    string
		initial utils.Indexes
		expand  utils.Indexes
		clear   utils.Indexes
		locked  utils.UInt8s
		hot     utils.UInt8s
		bonus   utils.Index
		sticky  utils.Index
		super   utils.Index
		j       string
	}{
		{
			name:    "simple",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{},
		},
		{
			name:    "with expand",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{1, 2, 3, 5, 5, 5, 1, 2, 3, 5, 5, 5, 1, 2, 3},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{},
		},
		{
			name:    "with clear",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{1, 2, 3, 0, 0, 0, 1, 2, 3, 0, 0, 0, 1, 2, 3},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{},
		},
		{
			name:    "locked reels",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{1, 4},
			hot:     utils.UInt8s{},
		},
		{
			name:    "hot reels",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{2, 3},
		},
		{
			name:    "bonus symbol",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{},
			bonus:   3,
		},
		{
			name:    "sticky symbol",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{},
			sticky:  7,
		},
		{
			name:    "super symbol",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 11, 2, 11, 4, 11, 6, 11, 2, 11},
			expand:  utils.Indexes{},
			clear:   utils.Indexes{},
			locked:  utils.UInt8s{},
			hot:     utils.UInt8s{},
			super:   11,
		},
		{
			name:    "all",
			initial: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			expand:  utils.Indexes{1, 2, 3, 5, 5, 5, 1, 2, 3, 5, 5, 5, 1, 2, 3, 5, 5, 5},
			clear:   utils.Indexes{0, 0, 0, 5, 5, 5, 0, 0, 0, 5, 5, 5, 0, 0, 0},
			locked:  utils.UInt8s{2, 4},
			hot:     utils.UInt8s{1, 3, 5},
			bonus:   9,
			sticky:  5,
			super:   3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := AcquireSpinResultFromData(tc.initial, tc.expand, tc.clear, tc.locked, tc.hot, tc.bonus, tc.sticky, tc.super)
			require.NotNil(t, spin)
			defer spin.Release()

			assert.EqualValues(t, tc.initial, spin.initial)
			assert.EqualValues(t, tc.expand, spin.AfterExpand())
			assert.EqualValues(t, tc.clear, spin.AfterClear())
			assert.EqualValues(t, tc.locked, spin.lockedReels)
			assert.EqualValues(t, tc.hot, spin.hotReels)
			assert.EqualValues(t, tc.bonus, spin.BonusSymbol())
			assert.EqualValues(t, tc.sticky, spin.StickySymbol())
			assert.EqualValues(t, tc.super, spin.SuperSymbol())
			assert.EqualValues(t, tc.initial, spin.initial)
		})
	}
}

func TestSpinResult_SetStickyChoice(t *testing.T) {
	testCases := []struct {
		name    string
		indexes utils.Indexes
		want    []*StickyChoice
	}{
		{
			name:    "none",
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			want: []*StickyChoice{
				{Symbol: 14, Positions: utils.UInt8s{13}},
				{Symbol: 13, Positions: utils.UInt8s{12}},
				{Symbol: 12, Positions: utils.UInt8s{11}},
				{Symbol: 11, Positions: utils.UInt8s{10}},
				{Symbol: 10, Positions: utils.UInt8s{9}},
				{Symbol: 9, Positions: utils.UInt8s{8}},
				{Symbol: 8, Positions: utils.UInt8s{7}},
				{Symbol: 7, Positions: utils.UInt8s{6}},
				{Symbol: 6, Positions: utils.UInt8s{5}},
				{Symbol: 5, Positions: utils.UInt8s{4}},
				{Symbol: 4, Positions: utils.UInt8s{3}},
				{Symbol: 3, Positions: utils.UInt8s{2}},
				{Symbol: 2, Positions: utils.UInt8s{1}},
				{Symbol: 1, Positions: utils.UInt8s{0}},
			},
		},
		{
			name:    "single",
			indexes: utils.Indexes{1, 2, 3, 1, 4, 5, 1, 6, 7, 1, 8, 9, 1, 10, 11},
			want: []*StickyChoice{
				{Symbol: 11, Positions: utils.UInt8s{14}},
				{Symbol: 10, Positions: utils.UInt8s{13}},
				{Symbol: 9, Positions: utils.UInt8s{11}},
				{Symbol: 8, Positions: utils.UInt8s{10}},
				{Symbol: 7, Positions: utils.UInt8s{8}},
				{Symbol: 6, Positions: utils.UInt8s{7}},
				{Symbol: 5, Positions: utils.UInt8s{5}},
				{Symbol: 4, Positions: utils.UInt8s{4}},
				{Symbol: 3, Positions: utils.UInt8s{2}},
				{Symbol: 2, Positions: utils.UInt8s{1}},
				{Symbol: 1, Positions: utils.UInt8s{0, 3, 6, 9, 12}},
			},
		},
		{
			name:    "two",
			indexes: utils.Indexes{1, 2, 3, 1, 4, 5, 1, 2, 7, 2, 8, 9, 1, 10, 11},
			want: []*StickyChoice{
				{Symbol: 11, Positions: utils.UInt8s{14}},
				{Symbol: 10, Positions: utils.UInt8s{13}},
				{Symbol: 9, Positions: utils.UInt8s{11}},
				{Symbol: 8, Positions: utils.UInt8s{10}},
				{Symbol: 7, Positions: utils.UInt8s{8}},
				{Symbol: 5, Positions: utils.UInt8s{5}},
				{Symbol: 4, Positions: utils.UInt8s{4}},
				{Symbol: 3, Positions: utils.UInt8s{2}},
				{Symbol: 2, Positions: utils.UInt8s{1, 7, 9}},
				{Symbol: 1, Positions: utils.UInt8s{0, 3, 6, 12}},
			},
		},
		{
			name:    "couple",
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 4, 5, 6, 1, 1, 3, 4, 7, 8},
			want: []*StickyChoice{
				{Symbol: 8, Positions: utils.UInt8s{14}},
				{Symbol: 7, Positions: utils.UInt8s{13}},
				{Symbol: 6, Positions: utils.UInt8s{8}},
				{Symbol: 5, Positions: utils.UInt8s{7}},
				{Symbol: 4, Positions: utils.UInt8s{6, 12}},
				{Symbol: 3, Positions: utils.UInt8s{2, 5, 11}},
				{Symbol: 2, Positions: utils.UInt8s{1, 4}},
				{Symbol: 1, Positions: utils.UInt8s{0, 3, 9, 10}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			spin := &Spin{symbols: setF1, indexes: tc.indexes}

			r := AcquireSpinResult(spin)
			require.NotNil(t, r)
			defer r.Release()

			assert.EqualValues(t, tc.indexes, r.Initial())

			r.SetStickyChoices(spin)

			bad := len(tc.want) != len(r.stickyChoices)
			if !bad {
				for ix := range tc.want {
					if !tc.want[ix].Equals(r.stickyChoices[ix]) {
						bad = true
					}
				}
			}
			if bad {
				assert.EqualValues(t, tc.want, r.stickyChoices)
			}
		})
	}
}

func BenchmarkAcquireSpinResult3x3(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(3, 3), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		r := AcquireSpinResult(s)
		r.Release()
	}
}

func BenchmarkAcquireSpinResult5x3(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		r := AcquireSpinResult(s)
		r.Release()
	}
}

func BenchmarkAcquireSpinResult6x4(b *testing.B) {
	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(6, 4), WithSymbols(setF1), PayDirections(PayLTR))

	s := AcquireSpin(slots, prng)
	defer s.Release()

	for i := 0; i < b.N; i++ {
		r := AcquireSpinResult(s)
		r.Release()
	}
}
