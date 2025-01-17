package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHand(t *testing.T) {
	testCases := []struct {
		name  string
		cards Cards
	}{
		{"empty", Cards{}},
		{"1 heart", Cards{NewCard(HeartJ)}},
		{"2 hearts", Cards{NewCard(HeartJ), NewCard(HeartK)}},
		{"aces", Cards{NewCard(HeartA), NewCard(SpadeA), NewCard(ClubA), NewCard(DiamondA)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHand(tc.cards...)
			require.NotNil(t, h)
			defer h.Release()

			assert.Equal(t, len(tc.cards), h.Size())
			assert.NotNil(t, h.CopyCards())
			assert.Equal(t, tc.cards, h.CopyCards())
		})
	}
}

func TestHand_Add(t *testing.T) {
	c1 := NewCard(Diamond2)
	c2 := NewCard(Club3)
	c3 := NewCard(Diamond7)
	c4 := NewCard(SpadeJ)
	c5 := NewCard(Club6)
	c6 := NewCard(Heart4)

	testCases := []struct {
		name    string
		initial Cards
		add     Cards
		want    Cards
	}{
		{"nils", nil, nil, nil},
		{"nil + 1", nil, Cards{c1}, Cards{c1}},
		{"nil + 2", nil, Cards{c1, c2}, Cards{c1, c2}},
		{"1 + nil", Cards{c3}, nil, Cards{c3}},
		{"2 + nil", Cards{c4, c5}, nil, Cards{c4, c5}},
		{"1 + 1", Cards{c3}, Cards{c6}, Cards{c3, c6}},
		{"1 + 2", Cards{c4}, Cards{c2, c3}, Cards{c4, c2, c3}},
		{"2 + 1", Cards{c5, c1}, Cards{c6}, Cards{c5, c1, c6}},
		{"2 + 2", Cards{c3, c2}, Cards{c5, c1}, Cards{c3, c2, c5, c1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var h *Hand
			h = NewHand(tc.initial...)
			require.NotNil(t, h)
			defer h.Release()

			h.Add(tc.add...)

			got := h.CopyCards()
			assert.NotNil(t, got)
			if tc.want != nil {
				assert.EqualValues(t, tc.want, got)
			} else {
				assert.Zero(t, len(got))
			}
		})
	}
}

func TestHand_Sorted(t *testing.T) {
	c1 := NewCard(Diamond2)
	c2 := NewCard(Club3)
	c3 := NewCard(Diamond7)
	c4 := NewCard(Spade2)
	c5 := NewCard(Club6)
	c6 := NewCard(Heart4)

	testCases := []struct {
		name  string
		cards Cards
		order CardSorter
		want  Cards
	}{
		{"empty value-suit", nil, SortByValueSuit, Cards{}},
		{"empty suit-value", nil, SortBySuitValue, Cards{}},
		{"single value-suit", Cards{c1}, SortByValueSuit, Cards{c1}},
		{"single suit-value", Cards{c3}, SortBySuitValue, Cards{c3}},
		{"three value-suit", Cards{c4, c1, c6}, SortByValueSuit, Cards{c1, c4, c6}},
		{"three suit-value", Cards{c5, c2, c3}, SortBySuitValue, Cards{c3, c2, c5}},
		{"all value-suit", Cards{c2, c5, c4, c3, c1, c6}, SortByValueSuit, Cards{c1, c4, c2, c6, c5, c3}},
		{"all suit-value", Cards{c5, c1, c6, c2, c4, c3}, SortBySuitValue, Cards{c1, c3, c2, c5, c6, c4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHand(tc.cards...)
			require.NotNil(t, h)
			defer h.Release()

			got := h.Sorted(tc.order)
			assert.EqualValues(t, tc.want, got)
		})
	}
}

func TestHand_Remove(t *testing.T) {
	c1 := NewCard(Diamond2)
	c2 := NewCard(Club3)
	c3 := NewCard(Diamond7)
	c4 := NewCard(SpadeJ)
	c5 := NewCard(Club6)
	c6 := NewCard(Heart4)

	testCases := []struct {
		name    string
		initial Cards
		remove  Cards
		want    Cards
	}{
		{"nils", nil, nil, Cards{}},
		{"nil - 1", nil, Cards{c1}, Cards{}},
		{"nil - 2", nil, Cards{c1, c2}, Cards{}},
		{"1 - nil", Cards{c3}, nil, Cards{c3}},
		{"2 - nil", Cards{c4, c5}, nil, Cards{c4, c5}},
		{"1 - 1", Cards{c3}, Cards{c6}, Cards{c3}},
		{"1 - 1", Cards{c4}, Cards{c4}, Cards{}},
		{"2 - 1", Cards{c5, c1}, Cards{c5}, Cards{c1}},
		{"2 - 2", Cards{c3, c2}, Cards{c2, c3}, Cards{}},
		{"4 - 1", Cards{c3, c1, c5, c2}, Cards{c2}, Cards{c3, c1, c5}},
		{"4 - 2", Cards{c3, c2, c6, c4}, Cards{c4, c3}, Cards{c2, c6}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHand(tc.initial...)
			require.NotNil(t, h)
			defer h.Release()

			h.Remove(tc.remove...)
			got := h.CopyCards()
			if tc.want != nil {
				assert.EqualValues(t, tc.want, got)
			}
		})
	}
}

func TestHand_SortInto(t *testing.T) {
	c1 := NewCard(Diamond2)
	c2 := NewCard(Club3)
	c3 := NewCard(Diamond7)
	c4 := NewCard(Spade2)
	c5 := NewCard(Club6)
	c6 := NewCard(Heart4)

	testCases := []struct {
		name  string
		cards Cards
		order CardSorter
		want  Cards
	}{
		{"empty value-suit", nil, SortByValueSuit, Cards{}},
		{"empty suit-value", nil, SortBySuitValue, Cards{}},
		{"single value-suit", Cards{c1}, SortByValueSuit, Cards{c1}},
		{"single suit-value", Cards{c3}, SortBySuitValue, Cards{c3}},
		{"three value-suit", Cards{c4, c1, c6}, SortByValueSuit, Cards{c1, c4, c6}},
		{"three suit-value", Cards{c5, c2, c3}, SortBySuitValue, Cards{c3, c2, c5}},
		{"all value-suit", Cards{c2, c5, c4, c3, c1, c6}, SortByValueSuit, Cards{c1, c4, c2, c6, c5, c3}},
		{"all suit-value", Cards{c5, c1, c6, c2, c4, c3}, SortBySuitValue, Cards{c1, c3, c2, c5, c6, c4}},
	}

	work := make(Cards, 2)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := NewHand(tc.cards...)
			require.NotNil(t, h)
			defer h.Release()

			got := h.SortInto(tc.order, work)
			assert.EqualValues(t, tc.want, got)
		})
	}
}
