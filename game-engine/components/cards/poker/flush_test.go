package poker

import (
	"testing"

	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"

	"github.com/stretchr/testify/assert"
)

func TestTestFlush(t *testing.T) {
	H4 := cards.NewCard(cards.Heart4)
	H6 := cards.NewCard(cards.Heart6)
	H8 := cards.NewCard(cards.Heart8)
	HJ := cards.NewCard(cards.HeartJ)
	HK := cards.NewCard(cards.HeartK)
	C4 := cards.NewCard(cards.Club4)
	C6 := cards.NewCard(cards.Club6)
	C8 := cards.NewCard(cards.Club8)
	S4 := cards.NewCard(cards.Spade4)
	S6 := cards.NewCard(cards.Spade6)
	SK := cards.NewCard(cards.SpadeK)

	testCases := []struct {
		name   string
		cards  cards.Cards
		want   bool
		flush  cards.Cards
		remain cards.Cards
	}{
		{"3 cards, no", cards.Cards{H8, C8, H6}, false, nil, nil},
		{"3 cards, yes", cards.Cards{H8, HK, H4}, true, cards.Cards{HK, H8, H4}, nil},
		{"5 cards, no (1)", cards.Cards{HJ, C8, HK, H6, H4}, false, nil, nil},
		{"5 cards, no (2)", cards.Cards{C8, C6, H4, HJ, C4}, false, nil, nil},
		{"5 cards, yes", cards.Cards{HJ, H4, H6, H8, HK}, true, cards.Cards{HK, HJ, H8, H6, H4}, nil},
		{"7 cards, no (1)", cards.Cards{HJ, C6, C8, H6, HK, C4, H4}, false, nil, nil},
		{"7 cards, no (2)", cards.Cards{C8, C6, S6, S4, H4, HJ, C4}, false, nil, nil},
		{"7 cards, yes (1)", cards.Cards{H4, C8, H8, H6, C6, HJ, HK}, true, cards.Cards{HK, HJ, H8, H6, H4}, cards.Cards{C8, C6}},
		{"7 cards, yes (2)", cards.Cards{SK, HJ, H4, H6, HK, H8, C4}, true, cards.Cards{HK, HJ, H8, H6, H4}, cards.Cards{SK, C4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRanked(tc.cards)
			require.NotNil(t, r)
			defer r.Release()

			r.findFlush()
			if tc.flush != nil {
				for ix, card := range tc.flush {
					assert.Equal(t, card, r.flush[ix])
				}
			}

			got2 := r.haveFlush()
			assert.Equal(t, tc.want, got2)

			if tc.remain != nil {
				r.getRemaining()
				for ix, card := range tc.remain {
					assert.Equal(t, card, r.remaining[ix])
				}
			}
		})
	}
}
