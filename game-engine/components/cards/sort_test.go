package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSortCards(t *testing.T) {
	c1 := NewCard(Diamond2)
	c2 := NewCard(Club3)
	c3 := NewCard(Diamond7)
	c4 := NewCard(Spade2)
	c5 := NewCard(Club6)
	c6 := NewCard(Heart4)
	c7 := NewCard(Diamond3)
	c8 := NewCard(ClubX)
	c9 := NewCard(Spade4)
	c10 := NewCard(HeartK)
	c11 := NewCard(DiamondJ)

	testCases := []struct {
		name  string
		cards Cards
		order CardSorter
		want  Cards
	}{
		{"empty value-suit", nil, SortByValueSuit, Cards{}},
		{"empty suit-value", nil, SortBySuitValue, Cards{}},
		{"empty value-suit", nil, SortByOrdinalSuit, Cards{}},
		{"empty suit-value", nil, SortBySuitOrdinal, Cards{}},
		{"single value-suit", Cards{c1}, SortByValueSuit, Cards{c1}},
		{"single suit-value", Cards{c3}, SortBySuitValue, Cards{c3}},
		{"single value-suit", Cards{c1}, SortByOrdinalSuit, Cards{c1}},
		{"single suit-value", Cards{c3}, SortBySuitOrdinal, Cards{c3}},
		{"three value-suit", Cards{c4, c1, c6}, SortByValueSuit, Cards{c1, c4, c6}},
		{"three suit-value", Cards{c5, c2, c3}, SortBySuitValue, Cards{c3, c2, c5}},
		{"three value-suit", Cards{c4, c1, c6}, SortByOrdinalSuit, Cards{c1, c4, c6}},
		{"three suit-value", Cards{c5, c2, c3}, SortBySuitOrdinal, Cards{c3, c2, c5}},
		{"six value-suit", Cards{c2, c5, c4, c3, c1, c6}, SortByValueSuit, Cards{c1, c4, c2, c6, c5, c3}},
		{"six suit-value", Cards{c5, c1, c6, c2, c4, c3}, SortBySuitValue, Cards{c1, c3, c2, c5, c6, c4}},
		{"six value-suit", Cards{c2, c5, c4, c3, c1, c6}, SortByOrdinalSuit, Cards{c1, c4, c2, c6, c5, c3}},
		{"six suit-value", Cards{c5, c1, c6, c2, c4, c3}, SortBySuitOrdinal, Cards{c1, c3, c2, c5, c6, c4}},
		{"many suit-value", Cards{c5, c10, c1, c6, c7, c9, c2, c8, c4, c11, c3}, SortBySuitOrdinal, Cards{c1, c7, c3, c11, c2, c5, c8, c6, c10, c4, c9}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := SortCards(tc.cards, tc.order)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestSortCardsInto(t *testing.T) {
	work := make(Cards, 0)

	t.Run("sort into", func(t *testing.T) {
		d := NewDeck(StandardDeck())
		require.NotNil(t, d)

		c := d.cards
		require.NotNil(t, c)

		s := SortCardsInto(c, SortByValueSuit, work)
		require.NotNil(t, s)
		require.NotEqualValues(t, s, c)

		s2 := SortCardsInto(c, SortByValueSuit, s)
		require.NotNil(t, s2)
		require.NotEqualValues(t, s2, c)
		require.EqualValues(t, s2, s)
	})
}
