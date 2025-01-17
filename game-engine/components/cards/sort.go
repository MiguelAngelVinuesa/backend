package cards

// SortCards returns a copy of the slice of cards which is sorted by the sorter function.
func SortCards(cards Cards, f CardSorter) Cards {
	return SortCardsInto(cards, f, nil)
}

// SortCardsInto sorts the cards using the sorter function and stores the result in the out slice.
// It returns the out slice.
// If the out slice is nil it will be created, otherwise shrunk/grown to hold exactly all cards.
func SortCardsInto(cards Cards, f CardSorter, out Cards) Cards {
	max := len(cards)
	if out == nil || cap(out) < max {
		if max < 8 {
			out = make(Cards, max, 8)
		} else {
			out = make(Cards, max)
		}
	} else {
		out = out[:max]
	}
	if max == 0 {
		return out
	}
	copy(out, cards)
	return quickSort(out, f)
}

// CardSorter is the function signature for card sorting functions.
type CardSorter func(cards Cards, i, j int) bool

// SortByValueSuit sorts the cards by value and suit
func SortByValueSuit(cards Cards, i, j int) bool {
	switch {
	case cards[i].value == cards[j].value:
		return cards[i].suit < cards[j].suit
	default:
		return cards[i].value < cards[j].value
	}
}

// SortBySuitValue sorts the cards by suit and value
func SortBySuitValue(cards Cards, i, j int) bool {
	switch {
	case cards[i].suit == cards[j].suit:
		return cards[i].value < cards[j].value
	default:
		return cards[i].suit < cards[j].suit
	}
}

// SortByOrdinalSuit sorts the cards by ordinal and suit
func SortByOrdinalSuit(cards Cards, i, j int) bool {
	switch {
	case cards[i].ordinal == cards[j].ordinal:
		return cards[i].suit < cards[j].suit
	default:
		return cards[i].ordinal < cards[j].ordinal
	}
}

// SortBySuitOrdinal sorts the cards by ordinal and suit
func SortBySuitOrdinal(cards Cards, i, j int) bool {
	switch {
	case cards[i].suit == cards[j].suit:
		return cards[i].ordinal < cards[j].ordinal
	default:
		return cards[i].suit < cards[j].suit
	}
}

// we use our own quick-sort as the standard sort package uses reflection which results in heap allocations.
func quickSort(cards Cards, f CardSorter) Cards {
	max := len(cards)
	if max < 2 {
		return cards
	}

	left, right, split := 0, max-1, max/2
	cards[right], cards[split] = cards[split], cards[right]

	for i := range cards {
		if f(cards, i, right) {
			cards[left], cards[i] = cards[i], cards[left]
			left++
		}
	}

	cards[left], cards[right] = cards[right], cards[left]

	quickSort(cards[:left], f)
	quickSort(cards[left+1:], f)

	return cards
}
