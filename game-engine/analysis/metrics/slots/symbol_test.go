package slots

import (
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewSymbol(t *testing.T) {
	testCases := []struct {
		name       string
		symbolID   utils.Index
		symbolName string
		reelCount  int
		first      []int
		second     []int
		free       []int
		free2      []int
		want       *Symbol
		j          string
	}{
		{
			name:       "empty",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalReels:      []uint64{0, 0, 0, 0, 0},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "one first",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			first:      []int{3},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      1,
				TotalReels:      []uint64{0, 0, 0, 1, 0},
				FirstCount:      1,
				FirstReels:      []uint64{0, 0, 0, 1, 0},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "few first",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			first:      []int{3, 2, 4, 1, 0, 2},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      6,
				TotalReels:      []uint64{1, 1, 2, 1, 1},
				FirstCount:      6,
				FirstReels:      []uint64{1, 1, 2, 1, 1},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "one second",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			second:     []int{3},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      1,
				TotalReels:      []uint64{0, 0, 0, 1, 0},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondCount:     1,
				SecondReels:     []uint64{0, 0, 0, 1, 0},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "few second",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			second:     []int{3, 2, 4, 1, 0, 2},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      6,
				TotalReels:      []uint64{1, 1, 2, 1, 1},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondCount:     6,
				SecondReels:     []uint64{1, 1, 2, 1, 1},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "one free",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			free:       []int{3},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      1,
				TotalReels:      []uint64{0, 0, 0, 1, 0},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeCount:       1,
				FreeReels:       []uint64{0, 0, 0, 1, 0},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "few free",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			free:       []int{3, 1, 0, 2, 4, 2},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      6,
				TotalReels:      []uint64{1, 1, 2, 1, 1},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeCount:       6,
				FreeReels:       []uint64{1, 1, 2, 1, 1},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "one free second",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			free2:      []int{3},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      1,
				TotalReels:      []uint64{0, 0, 0, 1, 0},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondCount: 1,
				FreeSecondReels: []uint64{0, 0, 0, 1, 0},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "few free second",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			free2:      []int{3, 1, 0, 2, 4, 2},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      6,
				TotalReels:      []uint64{1, 1, 2, 1, 1},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondCount: 6,
				FreeSecondReels: []uint64{1, 1, 2, 1, 1},
				Payouts:         []uint64{},
			},
		},
		{
			name:       "many first & free",
			symbolID:   1,
			symbolName: "10",
			reelCount:  5,
			first:      []int{3, 2, 1, 2, 4, 4, 2, 1, 0, 4, 4, 2, 0, 0, 1},
			free:       []int{2, 1, 0, 2, 4, 2, 2, 1, 4, 0, 0, 1, 2, 4, 0},
			want: &Symbol{
				ID:              1,
				Name:            "10",
				TotalCount:      30,
				TotalReels:      []uint64{7, 6, 9, 1, 7},
				FirstCount:      15,
				FirstReels:      []uint64{3, 3, 4, 1, 4},
				SecondReels:     []uint64{0, 0, 0, 0, 0},
				FreeCount:       15,
				FreeReels:       []uint64{4, 3, 5, 0, 3},
				FreeSecondReels: []uint64{0, 0, 0, 0, 0},
				Payouts:         []uint64{},
			},
			j: `{"id":1,"totalCount":30,"firstCount":15,"freeCount":15,"totalReels":[7,6,9,1,7],"firstReels":[3,3,4,1,4],"secondReels":[0,0,0,0,0],"freeReels":[4,3,5,0,3],"freeSecondReels":[0,0,0,0,0],"name":"10"}`,
		},
		{
			name:       "many second & free second",
			symbolID:   2,
			symbolName: "J",
			reelCount:  5,
			second:     []int{3, 2, 1, 2, 4, 4, 2, 1, 0, 4, 4, 2, 0, 0, 1},
			free2:      []int{2, 1, 0, 2, 4, 2, 2, 1, 4, 0, 0, 1, 2, 4, 0},
			want: &Symbol{
				ID:              2,
				Name:            "J",
				TotalCount:      30,
				TotalReels:      []uint64{7, 6, 9, 1, 7},
				FirstReels:      []uint64{0, 0, 0, 0, 0},
				SecondCount:     15,
				SecondReels:     []uint64{3, 3, 4, 1, 4},
				FreeReels:       []uint64{0, 0, 0, 0, 0},
				FreeSecondCount: 15,
				FreeSecondReels: []uint64{4, 3, 5, 0, 3},
				Payouts:         []uint64{},
			},
			j: `{"id":2,"totalCount":30,"firstCount":0,"secondCount":15,"freeCount":0,"freeSecondCount":15,"totalReels":[7,6,9,1,7],"firstReels":[0,0,0,0,0],"secondReels":[3,3,4,1,4],"freeReels":[0,0,0,0,0],"freeSecondReels":[4,3,5,0,3],"name":"J"}`,
		},
		{
			name:       "many all",
			symbolID:   2,
			symbolName: "J",
			reelCount:  5,
			first:      []int{3, 2, 1, 2, 4, 4, 2, 1, 0, 4, 4, 2, 0, 0, 1},
			second:     []int{3, 2, 1, 2, 4, 4, 2, 1, 0, 4, 4, 2, 0, 0, 1},
			free:       []int{2, 1, 0, 2, 4, 2, 2, 1, 4, 0, 0, 1, 2, 4, 0},
			free2:      []int{2, 1, 0, 2, 4, 2, 2, 1, 4, 0, 0, 1, 2, 4, 0},
			want: &Symbol{
				ID:              2,
				Name:            "J",
				TotalCount:      60,
				TotalReels:      []uint64{14, 12, 18, 2, 14},
				FirstCount:      15,
				FirstReels:      []uint64{3, 3, 4, 1, 4},
				SecondCount:     15,
				SecondReels:     []uint64{3, 3, 4, 1, 4},
				FreeCount:       15,
				FreeReels:       []uint64{4, 3, 5, 0, 3},
				FreeSecondCount: 15,
				FreeSecondReels: []uint64{4, 3, 5, 0, 3},
				Payouts:         []uint64{},
			},
			j: `{"id":2,"totalCount":60,"firstCount":15,"secondCount":15,"freeCount":15,"freeSecondCount":15,"totalReels":[14,12,18,2,14],"firstReels":[3,3,4,1,4],"secondReels":[3,3,4,1,4],"freeReels":[4,3,5,0,3],"freeSecondReels":[4,3,5,0,3],"name":"J"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewSymbol(tc.symbolID, tc.symbolName, "", tc.reelCount)
			require.NotNil(t, s)
			defer s.Release()

			for _, r := range tc.first {
				s.IncreaseFirst(r)
			}
			for _, r := range tc.second {
				s.IncreaseSecond(r)
			}
			for _, r := range tc.free {
				s.IncreaseFree(r)
			}
			for _, r := range tc.free2 {
				s.IncreaseSecondFree(r)
			}

			if !tc.want.Equals(s) {
				assert.EqualValues(t, tc.want, s)
			}

			if tc.j != "" {
				j, err := json.Marshal(s)
				require.NoError(t, err)
				assert.Equal(t, tc.j, string(j))
			}

			n := s.Clone().(*Symbol)
			require.NotNil(t, n)
			defer n.Release()

			if !tc.want.Equals(n) {
				assert.EqualValues(t, tc.want, n)
			}

			n.ResetData()
			assert.Zero(t, n.TotalCount)
			assert.Zero(t, n.FirstCount)
			assert.Zero(t, n.SecondCount)
			assert.Zero(t, n.FreeCount)
			assert.Zero(t, n.FreeSecondCount)

			for ix := range n.TotalReels {
				assert.Zero(t, n.TotalReels[ix])
				assert.Zero(t, n.FirstReels[ix])
				assert.Zero(t, n.SecondReels[ix])
				assert.Zero(t, n.FreeReels[ix])
				assert.Zero(t, n.FreeSecondReels[ix])
			}
		})
	}
}

func TestSymbol_IncreaseBonus(t *testing.T) {
	t.Run("increase bonus", func(t *testing.T) {
		s := NewSymbol(5, "five", "", 5)
		require.NotNil(t, s)
		defer s.Release()

		assert.Zero(t, s.BonusCount)
		assert.Zero(t, s.StickyCount)
		assert.Zero(t, s.SuperCount)

		for ix := 0; ix < 11; ix++ {
			s.IncreaseBonus()
		}
		assert.Equal(t, uint64(11), s.BonusCount)
		assert.Zero(t, s.StickyCount)
		assert.Zero(t, s.SuperCount)
	})
}

func TestSymbol_IncreaseSticky(t *testing.T) {
	t.Run("increase sticky", func(t *testing.T) {
		s := NewSymbol(5, "five", "", 5)
		require.NotNil(t, s)
		defer s.Release()

		assert.Zero(t, s.BonusCount)
		assert.Zero(t, s.StickyCount)
		assert.Zero(t, s.SuperCount)

		for ix := 0; ix < 12; ix++ {
			s.IncreaseSticky()
		}
		assert.Zero(t, s.BonusCount)
		assert.Equal(t, uint64(12), s.StickyCount)
		assert.Zero(t, s.SuperCount)
	})
}

func TestSymbol_IncreaseSuper(t *testing.T) {
	t.Run("increase super", func(t *testing.T) {
		s := NewSymbol(5, "five", "", 5)
		require.NotNil(t, s)
		defer s.Release()

		assert.Zero(t, s.BonusCount)
		assert.Zero(t, s.StickyCount)
		assert.Zero(t, s.SuperCount)

		for ix := 0; ix < 13; ix++ {
			s.IncreaseSuper()
		}
		assert.Zero(t, s.BonusCount)
		assert.Zero(t, s.StickyCount)
		assert.Equal(t, uint64(13), s.SuperCount)
	})
}
