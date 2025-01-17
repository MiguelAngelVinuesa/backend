package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStandardCards(t *testing.T) {
	testCases := []struct {
		name    string
		id      CardID
		ordinal Ordinal
		joker   bool
		suit    Suit
		color   Color
		value   int
	}{
		{"joker0", Joker0, 0, true, 0, Red, 15},
		{"joker1", Joker1, 0, true, 0, Black, 15},
		{"joker2", Joker2, 0, true, 0, Red, 15},
		{"joker3", Joker3, 0, true, 0, Black, 15},
		{"joker4", Joker4, 14, true, 0, Red, 15},
		{"joker5", Joker5, 14, true, 0, Black, 15},
		{"joker6", Joker6, 14, true, 0, Red, 15},
		{"joker7", Joker7, 14, true, 0, Black, 15},
		{"joker8", Joker8, 15, true, 0, Red, 15},
		{"joker9", Joker9, 15, true, 0, Black, 15},
		{"jokerX", JokerX, 15, true, 0, Red, 15},
		{"jokerY", JokerY, 15, true, 0, Black, 15},

		{"diamondA", DiamondA, Ace, false, Diamonds, Red, 1},
		{"diamond2", Diamond2, Two, false, Diamonds, Red, 2},
		{"diamond3", Diamond3, Three, false, Diamonds, Red, 3},
		{"diamond4", Diamond4, Four, false, Diamonds, Red, 4},
		{"diamond5", Diamond5, Five, false, Diamonds, Red, 5},
		{"diamond6", Diamond6, Six, false, Diamonds, Red, 6},
		{"diamond7", Diamond7, Seven, false, Diamonds, Red, 7},
		{"diamond8", Diamond8, Eight, false, Diamonds, Red, 8},
		{"diamond9", Diamond9, Nine, false, Diamonds, Red, 9},
		{"diamondX", DiamondX, Ten, false, Diamonds, Red, 10},
		{"diamondJ", DiamondJ, Jack, false, Diamonds, Red, 11},
		{"diamondQ", DiamondQ, Queen, false, Diamonds, Red, 12},
		{"diamondK", DiamondK, King, false, Diamonds, Red, 13},

		{"clubA", ClubA, Ace, false, Clubs, Black, 1},
		{"club2", Club2, Two, false, Clubs, Black, 2},
		{"club3", Club3, Three, false, Clubs, Black, 3},
		{"club4", Club4, Four, false, Clubs, Black, 4},
		{"club5", Club5, Five, false, Clubs, Black, 5},
		{"club6", Club6, Six, false, Clubs, Black, 6},
		{"club7", Club7, Seven, false, Clubs, Black, 7},
		{"club8", Club8, Eight, false, Clubs, Black, 8},
		{"club9", Club9, Nine, false, Clubs, Black, 9},
		{"clubX", ClubX, Ten, false, Clubs, Black, 10},
		{"clubJ", ClubJ, Jack, false, Clubs, Black, 11},
		{"clubQ", ClubQ, Queen, false, Clubs, Black, 12},
		{"clubK", ClubK, King, false, Clubs, Black, 13},

		{"heartA", HeartA, Ace, false, Hearts, Red, 1},
		{"heart2", Heart2, Two, false, Hearts, Red, 2},
		{"heart3", Heart3, Three, false, Hearts, Red, 3},
		{"heart4", Heart4, Four, false, Hearts, Red, 4},
		{"heart5", Heart5, Five, false, Hearts, Red, 5},
		{"heart6", Heart6, Six, false, Hearts, Red, 6},
		{"heart7", Heart7, Seven, false, Hearts, Red, 7},
		{"heart8", Heart8, Eight, false, Hearts, Red, 8},
		{"heart9", Heart9, Nine, false, Hearts, Red, 9},
		{"heartX", HeartX, Ten, false, Hearts, Red, 10},
		{"heartJ", HeartJ, Jack, false, Hearts, Red, 11},
		{"heartQ", HeartQ, Queen, false, Hearts, Red, 12},
		{"heartK", HeartK, King, false, Hearts, Red, 13},

		{"spadeA", SpadeA, Ace, false, Spades, Black, 1},
		{"spade2", Spade2, Two, false, Spades, Black, 2},
		{"spade3", Spade3, Three, false, Spades, Black, 3},
		{"spade4", Spade4, Four, false, Spades, Black, 4},
		{"spade5", Spade5, Five, false, Spades, Black, 5},
		{"spade6", Spade6, Six, false, Spades, Black, 6},
		{"spade7", Spade7, Seven, false, Spades, Black, 7},
		{"spade8", Spade8, Eight, false, Spades, Black, 8},
		{"spade9", Spade9, Nine, false, Spades, Black, 9},
		{"spadeX", SpadeX, Ten, false, Spades, Black, 10},
		{"spadeJ", SpadeJ, Jack, false, Spades, Black, 11},
		{"spadeQ", SpadeQ, Queen, false, Spades, Black, 12},
		{"spadeK", SpadeK, King, false, Spades, Black, 13},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCard(tc.id)
			require.NotNil(t, c)
			defer c.Release()

			assert.Equal(t, tc.id, c.ID())
			assert.Equal(t, tc.ordinal, c.Ordinal())
			assert.Equal(t, tc.joker, c.IsJoker())
			assert.Equal(t, tc.suit, c.Suit())
			assert.Equal(t, tc.color, c.Color())
			assert.Equal(t, tc.value, c.Value())

			if tc.ordinal == Ace {
				assert.True(t, c.IsAce())
			}
		})
	}
}

func TestCustomCards(t *testing.T) {
	testCases := []struct {
		name    string
		id      CardID
		opts    []CardOption
		ordinal Ordinal
		joker   bool
		suit    Suit
		color   Color
		value   int
	}{
		{"no opts   ", GameCard0 + 0, nil, 64, false, 0, 0, 64},
		{"red       ", GameCard0 + 1, []CardOption{WithColor(Red)}, 65, false, 0, Red, 65},
		{"black     ", GameCard0 + 2, []CardOption{WithColor(Black)}, 66, false, 0, Black, 66},
		{"hearts    ", GameCard0 + 3, []CardOption{WithSuit(Hearts)}, 67, false, Hearts, 0, 67},
		{"diamonds  ", GameCard0 + 4, []CardOption{WithSuit(Diamonds)}, 68, false, Diamonds, 0, 68},
		{"spades    ", GameCard0 + 5, []CardOption{WithSuit(Spades)}, 69, false, Spades, 0, 69},
		{"clubs     ", GameCard0 + 6, []CardOption{WithSuit(Clubs)}, 70, false, Clubs, 0, 70},
		{"joker     ", GameCard0 + 7, []CardOption{WithJoker(true)}, 71, true, 0, 0, 71},
		{"ordinal 99", GameCard0 + 8, []CardOption{WithOrdinal(99)}, 99, false, 0, 0, 72},
		{"value 99  ", GameCard0 + 9, []CardOption{WithValue(99)}, 73, false, 0, 0, 99},
		{"mixed 1   ", GameCard0 + 10, []CardOption{WithJoker(true), WithSuit(Clubs), WithColor(Red), WithValue(99)}, 74, true, Clubs, Red, 99},
		{"mixed 2   ", GameCard0 + 11, []CardOption{WithValue(0), WithColor(Black), WithSuit(Hearts), WithJoker(false)}, 75, false, Hearts, Black, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCard(tc.id, tc.opts...)
			require.NotNil(t, c)
			defer c.Release()

			assert.Equal(t, tc.id, c.ID())
			assert.Equal(t, tc.ordinal, c.Ordinal())
			assert.Equal(t, tc.joker, c.IsJoker())
			assert.Equal(t, tc.suit, c.Suit())
			assert.Equal(t, tc.color, c.Color())
			assert.Equal(t, tc.value, c.Value())
		})
	}
}

func TestResliceCards(t *testing.T) {
	testCases := []struct {
		name string
		in   Cards
		need int
		cap  int
	}{
		{"nil", nil, 5, 5},
		{"empty", Cards{}, 4, 4},
		{"1 element", Cards{NewCard(1)}, 8, 8},
		{"2 elements", Cards{NewCard(1), NewCard(2)}, 4, 4},
		{"4 elements", Cards{NewCard(1), NewCard(2), NewCard(3), NewCard(4)}, 4, 4},
		{"6 elements", Cards{NewCard(1), NewCard(2), NewCard(3), NewCard(4), NewCard(5), NewCard(6)}, 4, 6},
	}

	for _, tc := range testCases {
		s := ResliceCards(tc.in, tc.need)
		require.NotNil(t, s)
		assert.Zero(t, len(s))
		assert.Equal(t, tc.cap, cap(s))
	}
}
