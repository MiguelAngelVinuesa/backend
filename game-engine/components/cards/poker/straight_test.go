package poker

import (
	"fmt"
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestStraight(t *testing.T) {
	H2 := cards.NewCard(cards.Heart2)
	H3 := cards.NewCard(cards.Heart3)
	H4 := cards.NewCard(cards.Heart4)
	H5 := cards.NewCard(cards.Heart5)
	H6 := cards.NewCard(cards.Heart6)
	H9 := cards.NewCard(cards.Heart9)
	HX := cards.NewCard(cards.HeartX)
	HJ := cards.NewCard(cards.HeartJ)
	HQ := cards.NewCard(cards.HeartQ)
	HK := cards.NewCard(cards.HeartK)
	HA := cards.NewCard(cards.HeartA)
	C2 := cards.NewCard(cards.Club2)
	C3 := cards.NewCard(cards.Club3)
	C4 := cards.NewCard(cards.Club4)
	C5 := cards.NewCard(cards.Club5)
	C6 := cards.NewCard(cards.Club6)
	CX := cards.NewCard(cards.ClubX)
	CJ := cards.NewCard(cards.ClubJ)
	CQ := cards.NewCard(cards.ClubQ)
	CK := cards.NewCard(cards.ClubK)
	CA := cards.NewCard(cards.ClubA)
	S2 := cards.NewCard(cards.Spade2)
	S3 := cards.NewCard(cards.Spade3)
	S4 := cards.NewCard(cards.Spade4)
	S5 := cards.NewCard(cards.Spade5)
	S6 := cards.NewCard(cards.Spade6)
	SX := cards.NewCard(cards.SpadeX)
	SJ := cards.NewCard(cards.SpadeJ)
	SQ := cards.NewCard(cards.SpadeQ)
	SK := cards.NewCard(cards.SpadeK)
	SA := cards.NewCard(cards.SpadeA)

	testCases := []struct {
		name     string
		cards    cards.Cards
		want     bool
		rank     Rank
		straight cards.Cards
		remain   cards.Cards
	}{
		{"3 cards, none", cards.Cards{C4, H5, H9}, false, 0, nil, nil},
		{"3 cards, straight", cards.Cards{S4, H5, C3}, true, Straight, cards.Cards{H5, S4, C3}, nil},
		{"3 cards, straight ace-low", cards.Cards{H2, C3, HA}, true, Straight, cards.Cards{C3, H2, HA}, nil},
		{"3 cards, straight flush", cards.Cards{H4, H5, H6}, true, StraightFlush, cards.Cards{H6, H5, H4}, nil},
		{"3 cards, straight royal flush", cards.Cards{HK, HA, HQ}, true, RoyalFlush, cards.Cards{HA, HK, HQ}, nil},
		{"5 cards, none", cards.Cards{C4, H5, H2, C6, H9}, false, 0, nil, nil},
		{"5 cards, straight", cards.Cards{S4, H5, H2, C3, S6}, true, Straight, cards.Cards{S6, H5, S4, C3, H2}, nil},
		{"5 cards, straight ace-low", cards.Cards{S4, H5, H2, C3, HA}, true, Straight, cards.Cards{H5, S4, C3, H2, HA}, nil},
		{"5 cards, straight flush (1)", cards.Cards{H4, H5, H2, H3, H6}, true, StraightFlush, cards.Cards{H6, H5, H4, H3, H2}, nil},
		{"5 cards, straight flush (2)", cards.Cards{HK, HX, H9, HQ, HJ}, true, StraightFlush, cards.Cards{HK, HQ, HJ, HX, H9}, nil},
		{"5 cards, straight flush ace-low", cards.Cards{H4, H5, H2, H3, HA}, true, StraightFlush, cards.Cards{H5, H4, H3, H2, HA}, nil},
		{"5 cards, straight royal flush", cards.Cards{HK, HX, HA, HQ, HJ}, true, RoyalFlush, cards.Cards{HA, HK, HQ, HJ, HX}, nil},
		{"7 cards, none", cards.Cards{C4, H5, H2, C6, SK, S6, H9}, false, 0, nil, nil},
		{"7 cards, straight (1)", cards.Cards{S4, H5, C6, H2, H4, C3, S6}, true, Straight, cards.Cards{S6, H5, S4, C3, H2}, cards.Cards{C6, H4}},
		{"7 cards, straight (2)", cards.Cards{S4, H5, CA, H2, H4, C3, S6}, true, Straight, cards.Cards{S6, H5, S4, C3, H2}, cards.Cards{CA, H4}},
		{"7 cards, straight ace-low", cards.Cards{S4, H5, H4, SK, H2, C3, HA}, true, Straight, cards.Cards{H5, S4, C3, H2, HA}, cards.Cards{SK, H4}},
		{"7 cards, straight flush (1)", cards.Cards{H4, SK, HA, H5, H2, H3, H6}, true, StraightFlush, cards.Cards{H6, H5, H4, H3, H2}, cards.Cards{HA, SK}},
		{"7 cards, straight flush (2)", cards.Cards{HK, H6, C4, HX, H9, HQ, HJ}, true, StraightFlush, cards.Cards{HK, HQ, HJ, HX, H9}, cards.Cards{H6, C4}},
		{"7 cards, straight flush (3)", cards.Cards{HK, H6, CA, HX, H9, HQ, HJ}, true, StraightFlush, cards.Cards{HK, HQ, HJ, HX, H9}, cards.Cards{CA, H6}},
		{"7 cards, straight flush ace-low (1)", cards.Cards{H4, H5, SA, H2, CA, H3, HA}, true, StraightFlush, cards.Cards{H5, H4, H3, H2, HA}, cards.Cards{SA, CA}},
		{"7 cards, straight flush ace-low (2)", cards.Cards{C4, C5, SA, C2, CA, C3, HA}, true, StraightFlush, cards.Cards{C5, C4, C3, C2, CA}, cards.Cards{SA, HA}},
		{"7 cards, straight flush ace-low (3)", cards.Cards{S4, S5, SA, S2, CA, S3, HA}, true, StraightFlush, cards.Cards{S5, S4, S3, S2, SA}, cards.Cards{HA, CA}},
		{"7 cards, royal flush (1)", cards.Cards{SA, HK, HX, CA, HA, HQ, HJ}, true, RoyalFlush, cards.Cards{HA, HK, HQ, HJ, HX}, cards.Cards{SA, CA}},
		{"7 cards, royal flush (2)", cards.Cards{SA, SK, SX, CA, HA, SQ, SJ}, true, RoyalFlush, cards.Cards{SA, SK, SQ, SJ, SX}, cards.Cards{HA, CA}},
		{"7 cards, royal flush (3)", cards.Cards{SA, CK, CX, CA, HA, CQ, CJ}, true, RoyalFlush, cards.Cards{CA, CK, CQ, CJ, CX}, cards.Cards{SA, HA}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRanked(tc.cards)
			require.NotNil(t, r)
			defer r.Release()

			r.findFlush()
			r.findStraight()

			if tc.straight != nil {
				for ix, card := range tc.straight {
					assert.Equal(t, card, r.straight[ix])
				}
			}

			got := r.haveStraightFlush()
			if !got {
				got = r.haveStraight()
			}
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.rank, r.Rank())

			if tc.remain != nil {
				r.getRemaining()
				for ix, card := range tc.remain {
					assert.Equal(t, card, r.remaining[ix])
				}
			}
		})
	}
}

func TestTestStraight7(t *testing.T) {
	testCases := []struct {
		name   string
		cards  cards.Cards
		remove cards.Cards
	}{
		{
			name:   "ace high",
			cards:  cards.Cards{cards.NewCard(cards.HeartA), cards.NewCard(cards.ClubK), cards.NewCard(cards.DiamondQ), cards.NewCard(cards.SpadeJ), cards.NewCard(cards.HeartX)},
			remove: cards.Cards{},
		},
		{
			name:   "king high",
			cards:  cards.Cards{cards.NewCard(cards.HeartK), cards.NewCard(cards.ClubQ), cards.NewCard(cards.DiamondJ), cards.NewCard(cards.SpadeX), cards.NewCard(cards.Heart9)},
			remove: cards.Cards{cards.NewCard(cards.SpadeA), cards.NewCard(cards.HeartA), cards.NewCard(cards.ClubA), cards.NewCard(cards.DiamondA)},
		},
		{
			name:   "queen high",
			cards:  cards.Cards{cards.NewCard(cards.HeartQ), cards.NewCard(cards.ClubJ), cards.NewCard(cards.DiamondX), cards.NewCard(cards.Spade9), cards.NewCard(cards.Heart8)},
			remove: cards.Cards{cards.NewCard(cards.SpadeK), cards.NewCard(cards.HeartK), cards.NewCard(cards.ClubA), cards.NewCard(cards.DiamondA)},
		},
		{
			name:   "jack high",
			cards:  cards.Cards{cards.NewCard(cards.HeartJ), cards.NewCard(cards.ClubX), cards.NewCard(cards.Diamond9), cards.NewCard(cards.Spade8), cards.NewCard(cards.Heart7)},
			remove: cards.Cards{cards.NewCard(cards.SpadeQ), cards.NewCard(cards.HeartQ), cards.NewCard(cards.ClubQ), cards.NewCard(cards.DiamondQ)},
		},
		{
			name:   "10 high",
			cards:  cards.Cards{cards.NewCard(cards.HeartX), cards.NewCard(cards.Club9), cards.NewCard(cards.Diamond8), cards.NewCard(cards.Spade7), cards.NewCard(cards.Heart6)},
			remove: cards.Cards{cards.NewCard(cards.SpadeJ), cards.NewCard(cards.HeartJ), cards.NewCard(cards.ClubJ), cards.NewCard(cards.DiamondJ)},
		},
		{
			name:   "9 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart9), cards.NewCard(cards.Club8), cards.NewCard(cards.Diamond7), cards.NewCard(cards.Spade6), cards.NewCard(cards.Heart5)},
			remove: cards.Cards{cards.NewCard(cards.SpadeX), cards.NewCard(cards.HeartX), cards.NewCard(cards.ClubX), cards.NewCard(cards.DiamondX)},
		},
		{
			name:   "8 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart8), cards.NewCard(cards.Club7), cards.NewCard(cards.Diamond6), cards.NewCard(cards.Spade5), cards.NewCard(cards.Heart4)},
			remove: cards.Cards{cards.NewCard(cards.Spade9), cards.NewCard(cards.Heart9), cards.NewCard(cards.Club9), cards.NewCard(cards.Diamond9)},
		},
		{
			name:   "7 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart7), cards.NewCard(cards.Club6), cards.NewCard(cards.Diamond5), cards.NewCard(cards.Spade4), cards.NewCard(cards.Heart3)},
			remove: cards.Cards{cards.NewCard(cards.Spade8), cards.NewCard(cards.Heart8), cards.NewCard(cards.Club8), cards.NewCard(cards.Diamond8)},
		},
		{
			name:   "6 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart6), cards.NewCard(cards.Club5), cards.NewCard(cards.Diamond4), cards.NewCard(cards.Spade3), cards.NewCard(cards.Heart2)},
			remove: cards.Cards{cards.NewCard(cards.Spade7), cards.NewCard(cards.Heart7), cards.NewCard(cards.Club7), cards.NewCard(cards.Diamond7)},
		},
		{
			name:   "5 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart5), cards.NewCard(cards.Club4), cards.NewCard(cards.Diamond3), cards.NewCard(cards.Spade2), cards.NewCard(cards.SpadeA)},
			remove: cards.Cards{cards.NewCard(cards.Spade6), cards.NewCard(cards.Heart6), cards.NewCard(cards.Club6), cards.NewCard(cards.Diamond6)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := makeRestrictedDeck(tc.cards, tc.remove)
			require.NotNil(t, d)
			defer d.Release()

			hand := make([]*cards.Card, 7)

			for ix := 0; ix < 10000; ix++ {
				copy(hand, tc.cards)

				d.Shuffle()
				hand[5] = d.Draw()
				hand[6] = d.Draw()

				r := NewRanked(hand)
				require.NotNil(t, r)
				defer r.Release()

				if r.sorted[0].Ordinal() == cards.Ace {
					r.ace = r.sorted[0]
				}

				r.findFlush()
				r.findStraight()
				if !r.haveStraightFlush() {
					r.haveStraight()
				}

				require.Equal(t, Straight, r.rank, fmt.Sprintf("%+v", hand))
			}
		})
	}
}

func TestTestStraightFlush7(t *testing.T) {
	testCases := []struct {
		name   string
		cards  cards.Cards
		remove cards.Cards
	}{
		{
			name:   "king high",
			cards:  cards.Cards{cards.NewCard(cards.HeartK), cards.NewCard(cards.HeartQ), cards.NewCard(cards.HeartJ), cards.NewCard(cards.HeartX), cards.NewCard(cards.Heart9)},
			remove: cards.Cards{cards.NewCard(cards.HeartA)},
		},
		{
			name:   "queen high",
			cards:  cards.Cards{cards.NewCard(cards.HeartQ), cards.NewCard(cards.HeartJ), cards.NewCard(cards.HeartX), cards.NewCard(cards.Heart9), cards.NewCard(cards.Heart8)},
			remove: cards.Cards{cards.NewCard(cards.HeartK)},
		},
		{
			name:   "jack high",
			cards:  cards.Cards{cards.NewCard(cards.HeartJ), cards.NewCard(cards.HeartX), cards.NewCard(cards.Heart9), cards.NewCard(cards.Heart8), cards.NewCard(cards.Heart7)},
			remove: cards.Cards{cards.NewCard(cards.HeartQ)},
		},
		{
			name:   "10 high",
			cards:  cards.Cards{cards.NewCard(cards.HeartX), cards.NewCard(cards.Heart9), cards.NewCard(cards.Heart8), cards.NewCard(cards.Heart7), cards.NewCard(cards.Heart6)},
			remove: cards.Cards{cards.NewCard(cards.HeartJ)},
		},
		{
			name:   "9 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart9), cards.NewCard(cards.Heart8), cards.NewCard(cards.Heart7), cards.NewCard(cards.Heart6), cards.NewCard(cards.Heart5)},
			remove: cards.Cards{cards.NewCard(cards.HeartX)},
		},
		{
			name:   "8 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart8), cards.NewCard(cards.Heart7), cards.NewCard(cards.Heart6), cards.NewCard(cards.Heart5), cards.NewCard(cards.Heart4)},
			remove: cards.Cards{cards.NewCard(cards.Heart9)},
		},
		{
			name:   "7 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart7), cards.NewCard(cards.Heart6), cards.NewCard(cards.Heart5), cards.NewCard(cards.Heart4), cards.NewCard(cards.Heart3)},
			remove: cards.Cards{cards.NewCard(cards.Heart8)},
		},
		{
			name:   "6 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart6), cards.NewCard(cards.Heart5), cards.NewCard(cards.Heart4), cards.NewCard(cards.Heart3), cards.NewCard(cards.Heart2)},
			remove: cards.Cards{cards.NewCard(cards.Heart7)},
		},
		{
			name:   "5 high",
			cards:  cards.Cards{cards.NewCard(cards.Heart5), cards.NewCard(cards.Heart4), cards.NewCard(cards.Heart3), cards.NewCard(cards.Heart2), cards.NewCard(cards.HeartA)},
			remove: cards.Cards{cards.NewCard(cards.Heart6)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := makeRestrictedDeck(tc.cards, tc.remove)
			require.NotNil(t, d)
			defer d.Release()

			hand := make([]*cards.Card, 7)

			for ix := 0; ix < 10000; ix++ {
				copy(hand, tc.cards)

				d.Shuffle()
				hand[5] = d.Draw()
				hand[6] = d.Draw()

				r := NewRanked(hand)
				require.NotNil(t, r)
				defer r.Release()

				r.findFlush()
				r.findStraight()
				r.haveStraightFlush()

				if r.rank != StraightFlush {
					require.Equal(t, StraightFlush, r.rank, fmt.Sprintf("%+v", hand))
				}
			}
		})
	}
}

func TestTestRoyalFlush7(t *testing.T) {
	testCases := []struct {
		name   string
		cards  cards.Cards
		remove cards.Cards
	}{
		{
			name:  "spades",
			cards: cards.Cards{cards.NewCard(cards.SpadeA), cards.NewCard(cards.SpadeK), cards.NewCard(cards.SpadeQ), cards.NewCard(cards.SpadeJ), cards.NewCard(cards.SpadeX)},
		},
		{
			name:  "hearts",
			cards: cards.Cards{cards.NewCard(cards.HeartA), cards.NewCard(cards.HeartK), cards.NewCard(cards.HeartQ), cards.NewCard(cards.HeartJ), cards.NewCard(cards.HeartX)},
		},
		{
			name:  "clubs",
			cards: cards.Cards{cards.NewCard(cards.ClubA), cards.NewCard(cards.ClubK), cards.NewCard(cards.ClubQ), cards.NewCard(cards.ClubJ), cards.NewCard(cards.ClubX)},
		},
		{
			name:  "diamonds",
			cards: cards.Cards{cards.NewCard(cards.DiamondA), cards.NewCard(cards.DiamondK), cards.NewCard(cards.DiamondQ), cards.NewCard(cards.DiamondJ), cards.NewCard(cards.DiamondX)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := makeRestrictedDeck(tc.cards, tc.remove)
			require.NotNil(t, d)
			defer d.Release()

			hand := make([]*cards.Card, 7)

			for ix := 0; ix < 10000; ix++ {
				copy(hand, tc.cards)

				d.Shuffle()
				hand[5] = d.Draw()
				hand[6] = d.Draw()

				r := NewRanked(hand)
				require.NotNil(t, r)
				defer r.Release()

				r.findFlush()
				r.findStraight()
				r.haveStraightFlush()

				require.Equal(t, RoyalFlush, r.rank, fmt.Sprintf("%+v", hand))
			}
		})
	}
}

func makeRestrictedDeck(not1, not2 cards.Cards) *cards.Deck {
	d := cards.NewDeck()

	include := func(card *cards.Card) bool {
		for _, not := range not1 {
			if not.ID() == card.ID() {
				return false
			}
		}
		for _, not := range not2 {
			if not.ID() == card.ID() {
				return false
			}
		}
		return true
	}

	for _, id := range cards.DiamondsAll {
		c := cards.NewCard(id)
		if include(c) {
			d.AddCustom(c)
		}
	}

	for _, id := range cards.ClubsAll {
		c := cards.NewCard(id)
		if include(c) {
			d.AddCustom(c)
		}
	}

	for _, id := range cards.HeartsAll {
		c := cards.NewCard(id)
		if include(c) {
			d.AddCustom(c)
		}
	}

	for _, id := range cards.SpadesAll {
		c := cards.NewCard(id)
		if include(c) {
			d.AddCustom(c)
		}
	}

	return d
}
