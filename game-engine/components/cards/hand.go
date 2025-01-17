package cards

import (
	"sync"
)

// Hand represents a "hand" of cards for a single player.
// It is safe to use across multiple go-routines simultaneously.
type Hand struct {
	mutex sync.RWMutex
	cards Cards
}

var handPool = sync.Pool{
	New: func() interface{} {
		return &Hand{}
	},
}

// NewHand instantiates a new hand from the memory pool.
func NewHand(cards ...*Card) *Hand {
	h := handPool.Get().(*Hand)
	max := len(cards)
	if h.cards == nil || cap(h.cards) < max {
		h.cards = make(Cards, max)
	} else {
		h.cards = h.cards[:max]
	}
	copy(h.cards, cards)
	return h
}

// Release puts a hand back into the memory pool.
func (h *Hand) Release() {
	if h != nil {
		if h.cards != nil {
			h.cards = h.cards[:0]
		}
		handPool.Put(h)
	}
}

// Add adds one or more cards to the hand.
// It returns the new number of cards in the hand.
func (h *Hand) Add(cards ...*Card) int {
	h.mutex.Lock()
	h.cards = append(h.cards, cards...)
	l := len(h.cards)
	h.mutex.Unlock()
	return l
}

// Remove removes one or more cards from the hand.
// It returns the new number of cards in the hand.
func (h *Hand) Remove(cards ...*Card) int {
	h.mutex.Lock()

	n := make(Cards, 0, len(h.cards))

	for _, card := range h.cards {
		var remove bool
		for _, c := range cards {
			if c == card {
				remove = true
				break
			}
		}
		if !remove {
			n = append(n, card)
		}
	}

	h.cards = n

	h.mutex.Unlock()
	return len(n)
}

// Clear removes all cards from the hand.
func (h *Hand) Clear() {
	h.mutex.Lock()
	h.cards = h.cards[:0]
	h.mutex.Unlock()
}

// Size returns the current number of cards in the hand.
func (h *Hand) Size() int {
	h.mutex.RLock()
	l := len(h.cards)
	h.mutex.RUnlock()
	return l
}

// CopyCards returns all cards in the hand.
// It creates a new slice which is valid only until the next update of the hand!
func (h *Hand) CopyCards() Cards {
	h.mutex.RLock()
	out := make(Cards, len(h.cards))
	copy(out, h.cards)
	h.mutex.RUnlock()
	return out
}

// Sorted returns the cards in the hand sorted by the given sorter function.
// It creates a new slice which is valid only until the next update of the hand!
func (h *Hand) Sorted(f CardSorter) Cards {
	h.mutex.RLock()
	m := SortCardsInto(h.cards, f, make(Cards, len(h.cards)))
	h.mutex.RUnlock()
	return m
}

// SortInto sorts the cards in the hand using the sorter function and stores the result in the out slice.
// The out slice is created if it is nil, or shrunk/grown to hold all cards in the hand.
// The out slice is valid only until the next update of the hand!
func (h *Hand) SortInto(f CardSorter, out Cards) Cards {
	h.mutex.RLock()
	m := SortCardsInto(h.cards, f, out)
	h.mutex.RUnlock()
	return m
}
