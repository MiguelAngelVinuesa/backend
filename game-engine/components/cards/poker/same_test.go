package poker

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestSame(t *testing.T) {
	H4 := cards.NewCard(cards.Heart4)
	H6 := cards.NewCard(cards.Heart6)
	H8 := cards.NewCard(cards.Heart8)
	HJ := cards.NewCard(cards.HeartJ)
	HK := cards.NewCard(cards.HeartK)
	C3 := cards.NewCard(cards.Club3)
	C4 := cards.NewCard(cards.Club4)
	C6 := cards.NewCard(cards.Club6)
	C8 := cards.NewCard(cards.Club8)
	S4 := cards.NewCard(cards.Spade4)
	S6 := cards.NewCard(cards.Spade6)
	SK := cards.NewCard(cards.SpadeK)
	D2 := cards.NewCard(cards.Diamond2)
	D4 := cards.NewCard(cards.Diamond4)
	D6 := cards.NewCard(cards.Diamond6)
	D8 := cards.NewCard(cards.Diamond8)

	testCases := []struct {
		name   string
		cards  cards.Cards
		count  int
		want   bool
		same   cards.Cards
		remain cards.Cards
	}{
		{"3, none", cards.Cards{H6, C8, SK}, 2, false, nil, nil},
		{"3, pair", cards.Cards{H6, C8, H8}, 2, true, cards.Cards{H8, C8}, cards.Cards{H6}},
		{"3, three of a kind (1)", cards.Cards{H6, C6, S6}, 2, false, nil, nil},
		{"3, three of a kind (2)", cards.Cards{S4, C4, H4}, 3, true, cards.Cards{S4, H4, C4}, nil},
		{"5, none", cards.Cards{H6, C8, HK, C4, HJ}, 2, false, nil, nil},
		{"5, pair", cards.Cards{H6, C8, HK, HJ, H8}, 2, true, cards.Cards{H8, C8}, cards.Cards{HK, HJ, H6}},
		{"5, 2 pair", cards.Cards{H6, C8, HK, C6, H8}, 2, true, cards.Cards{H8, C8}, cards.Cards{HK, H6, C6}},
		{"5, three of a kind (1)", cards.Cards{SK, H6, HJ, C6, S6}, 2, false, nil, nil},
		{"5, three of a kind (2)", cards.Cards{S4, HJ, C4, SK, H4}, 3, true, cards.Cards{S4, H4, C4}, cards.Cards{SK, HJ}},
		{"5, full house (1)", cards.Cards{H6, C8, D6, C6, H8}, 2, true, cards.Cards{H8, C8}, cards.Cards{H6, C6, D6}},
		{"5, full house (2)", cards.Cards{H6, C8, D8, C6, H8}, 2, true, cards.Cards{H6, C6}, cards.Cards{H8, C8, D8}},
		{"5, full house (3)", cards.Cards{H6, C8, D6, C6, H8}, 3, true, cards.Cards{H6, C6, D6}, cards.Cards{H8, C8}},
		{"5, full house (4)", cards.Cards{H6, C8, D8, C6, H8}, 3, true, cards.Cards{H8, C8, D8}, cards.Cards{H6, C6}},
		{"5, four of a kind (1)", cards.Cards{D6, H6, HJ, C6, S6}, 2, false, nil, nil},
		{"5, four of a kind (2)", cards.Cards{S4, HJ, C4, D4, H4}, 3, false, nil, nil},
		{"5, four of a kind (3)", cards.Cards{S4, D4, C4, SK, H4}, 4, true, cards.Cards{S4, H4, C4, D4}, cards.Cards{SK}},
		{"7, none", cards.Cards{H6, C3, D2, C8, HK, C4, HJ}, 2, false, nil, nil},
		{"7, pair", cards.Cards{H6, C8, C3, D2, HK, HJ, H8}, 2, true, cards.Cards{H8, C8}, cards.Cards{HK, HJ, H6, C3, D2}},
		{"7, 2 pair", cards.Cards{H6, C8, HK, C6, H8, H4, HJ}, 2, true, cards.Cards{H8, C8}, cards.Cards{HK, HJ, H6, C6, H4}},
		{"7, 3 pair", cards.Cards{H6, C8, HK, C6, H8, SK, D2}, 2, true, cards.Cards{SK, HK}, cards.Cards{H8, C8, H6, C6, D2}},
		{"7, three of a kind (1)", cards.Cards{C3, D2, SK, H6, HJ, C6, S6}, 2, false, nil, nil},
		{"7, three of a kind (2)", cards.Cards{S4, HJ, C3, D2, C4, SK, H4}, 3, true, cards.Cards{S4, H4, C4}, cards.Cards{SK, HJ, C3, D2}},
		{"7, four of a kind (1)", cards.Cards{D2, D6, H6, HJ, C6, C3, S6}, 2, false, nil, nil},
		{"7, four of a kind (2)", cards.Cards{C3, S4, D2, HJ, C4, D4, H4}, 3, false, nil, nil},
		{"7, four of a kind (3)", cards.Cards{S4, D4, D2, C4, C3, SK, H4}, 4, true, cards.Cards{S4, H4, C4, D4}, cards.Cards{SK, C3, D2}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRanked(tc.cards)
			require.NotNil(t, r)
			defer r.Release()

			r.findSame()
			got := r.haveSameKind(tc.count)
			assert.Equal(t, tc.want, got)

			if tc.same != nil {
				require.Equal(t, len(tc.same), len(r.ranked))
				for ix, card := range tc.same {
					assert.Equal(t, card, r.ranked[ix])
				}
			}

			if tc.remain != nil {
				r.getRemaining()
				require.Equal(t, len(tc.remain), len(r.remaining))
				for ix, card := range tc.remain {
					assert.Equal(t, card, r.remaining[ix])
				}
			}
		})
	}
}

func TestTestTwoPair(t *testing.T) {
	H2 := cards.NewCard(cards.Heart2)
	H4 := cards.NewCard(cards.Heart4)
	H6 := cards.NewCard(cards.Heart6)
	H8 := cards.NewCard(cards.Heart8)
	HJ := cards.NewCard(cards.HeartJ)
	HK := cards.NewCard(cards.HeartK)
	C3 := cards.NewCard(cards.Club3)
	C4 := cards.NewCard(cards.Club4)
	C6 := cards.NewCard(cards.Club6)
	C8 := cards.NewCard(cards.Club8)
	S4 := cards.NewCard(cards.Spade4)
	S6 := cards.NewCard(cards.Spade6)
	SK := cards.NewCard(cards.SpadeK)
	D2 := cards.NewCard(cards.Diamond2)
	D4 := cards.NewCard(cards.Diamond4)
	D6 := cards.NewCard(cards.Diamond6)

	testCases := []struct {
		name   string
		cards  cards.Cards
		count  int
		want   bool
		same   cards.Cards
		remain cards.Cards
	}{
		{"3, none", cards.Cards{H6, C8, SK}, 2, false, nil, nil},
		{"3, pair", cards.Cards{H6, C8, H8}, 2, false, nil, nil},
		{"5, none", cards.Cards{H6, C8, HK, C4, HJ}, 2, false, nil, nil},
		{"5, pair", cards.Cards{H6, C8, HK, HJ, H8}, 2, false, nil, nil},
		{"5, 2 pair", cards.Cards{H6, C8, HK, C6, H8}, 2, true, cards.Cards{H8, C8, H6, C6}, cards.Cards{HK}},
		{"5, three of a kind", cards.Cards{SK, H6, HJ, C6, S6}, 2, false, nil, nil},
		{"5, four of a kind", cards.Cards{D6, H6, HJ, C6, S6}, 2, false, nil, nil},
		{"7, none", cards.Cards{H6, C3, D2, C8, HK, C4, HJ}, 2, false, nil, nil},
		{"7, pair", cards.Cards{H6, C8, C3, D2, HK, HJ, H8}, 2, false, nil, nil},
		{"7, 2 pair", cards.Cards{H6, C8, HK, C6, H8, H4, HJ}, 2, true, cards.Cards{H8, C8, H6, C6}, cards.Cards{HK, HJ, H4}},
		{"7, 3 pair", cards.Cards{H6, C8, HK, C6, H8, SK, D2}, 2, true, cards.Cards{SK, HK, H8, C8}, cards.Cards{H6, C6, D2}},
		{"7, three of a kind (1)", cards.Cards{C3, D2, SK, H6, HJ, C6, S6}, 2, false, nil, nil},
		{"7, three of a kind (2)", cards.Cards{C4, S4, SK, H6, H4, C6, S6}, 2, false, nil, nil},
		{"7, four of a kind (1)", cards.Cards{D2, D6, H6, HJ, C6, H2, S6}, 2, false, nil, nil},
		{"7, four of a kind (2)", cards.Cards{S4, D4, D2, C4, C3, SK, H4}, 2, false, nil, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRanked(tc.cards)
			require.NotNil(t, r)
			defer r.Release()

			r.findSame()
			got := r.haveTwiceSame(tc.count)
			assert.Equal(t, tc.want, got)

			if tc.same != nil {
				require.Equal(t, len(tc.same), len(r.ranked))
				for ix, card := range tc.same {
					assert.Equal(t, card, r.ranked[ix])
				}
			}

			if tc.remain != nil {
				r.getRemaining()
				require.Equal(t, len(tc.remain), len(r.remaining))
				for ix, card := range tc.remain {
					assert.Equal(t, card, r.remaining[ix])
				}
			}
		})
	}
}

func TestTestFullHouse(t *testing.T) {
	H4 := cards.NewCard(cards.Heart4)
	H6 := cards.NewCard(cards.Heart6)
	H8 := cards.NewCard(cards.Heart8)
	HJ := cards.NewCard(cards.HeartJ)
	HK := cards.NewCard(cards.HeartK)
	C3 := cards.NewCard(cards.Club3)
	C4 := cards.NewCard(cards.Club4)
	C6 := cards.NewCard(cards.Club6)
	C8 := cards.NewCard(cards.Club8)
	S4 := cards.NewCard(cards.Spade4)
	S6 := cards.NewCard(cards.Spade6)
	SK := cards.NewCard(cards.SpadeK)
	D2 := cards.NewCard(cards.Diamond2)
	D6 := cards.NewCard(cards.Diamond6)
	D8 := cards.NewCard(cards.Diamond8)

	testCases := []struct {
		name   string
		cards  cards.Cards
		count  int
		want   bool
		same   cards.Cards
		remain cards.Cards
	}{
		{"5, none", cards.Cards{H6, C8, HK, C4, HJ}, 3, false, nil, nil},
		{"5, pair", cards.Cards{H6, C8, HK, HJ, H8}, 3, false, nil, nil},
		{"5, 2 pair", cards.Cards{H6, C8, HK, C6, H8}, 3, false, nil, nil},
		{"5, three of a kind", cards.Cards{SK, H6, HJ, C6, S6}, 3, false, nil, nil},
		{"5, full house (1)", cards.Cards{H6, C8, D6, C6, H8}, 3, true, cards.Cards{H6, C6, D6, H8, C8}, nil},
		{"5, full house (2)", cards.Cards{H6, C8, D8, C6, H8}, 3, true, cards.Cards{H8, C8, D8, H6, C6}, nil},
		{"5, four of a kind", cards.Cards{D6, H6, HJ, C6, S6}, 3, false, nil, nil},
		{"7, none", cards.Cards{H6, C3, D2, C8, HK, C4, HJ}, 3, false, nil, nil},
		{"7, pair", cards.Cards{H6, C8, C3, D2, HK, HJ, H8}, 3, false, nil, nil},
		{"7, 2 pair", cards.Cards{H6, C8, HK, C6, H8, H4, HJ}, 3, false, nil, nil},
		{"7, 3 pair", cards.Cards{H6, C8, HK, C6, H8, SK, D2}, 3, false, nil, nil},
		{"7, three of a kind", cards.Cards{C3, D2, SK, H6, HJ, C6, S6}, 3, false, nil, nil},
		{"7, full house (1)", cards.Cards{H6, C8, HJ, D6, S4, C6, H8}, 3, true, cards.Cards{H6, C6, D6, H8, C8}, cards.Cards{HJ, S4}},
		{"7, full house (2)", cards.Cards{H6, HK, C8, D8, HJ, C6, H8}, 3, true, cards.Cards{H8, C8, D8, H6, C6}, cards.Cards{HK, HJ}},
		{"7, four of a kind", cards.Cards{D2, D6, H6, HK, C6, SK, S6}, 3, false, nil, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRanked(tc.cards)
			require.NotNil(t, r)
			defer r.Release()

			r.findSame()
			got := r.haveTwiceSame(tc.count)
			assert.Equal(t, tc.want, got)

			if tc.same != nil {
				require.Equal(t, len(tc.same), len(r.ranked))
				for ix, card := range tc.same {
					assert.Equal(t, card, r.ranked[ix])
				}
			}

			if tc.remain != nil {
				r.getRemaining()
				require.Equal(t, len(tc.remain), len(r.remaining))
				for ix, card := range tc.remain {
					assert.Equal(t, card, r.remaining[ix])
				}
			}
		})
	}
}
