package poker

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

func (r *Ranked) findStraight() *Ranked {
	if r.sorted[0].IsAce() {
		r.ace = r.sorted[0]
	}

	if r.max <= 5 {
		// for ace-low we need an extra iteration.
		for ix := 0; ix <= 1; {
			ix += r.testStraight(r.sorted[ix], r.sorted[ix+1:])
		}
	} else {
		for ix := 0; ix <= 3; ix++ {
			if r.testStraight(r.sorted[ix], r.sorted[ix+1:]) > 1 && r.straightFlush {
				return r
			}
		}
		for ix := 0; ix <= 3; {
			ix += r.testStraight(r.sorted[ix], r.sorted[ix+1:])
		}
	}

	return r
}

func (r *Ranked) testStraight(card *cards.Card, other cards.Cards) int {
	count := 1
	last := card.Ordinal()
	for _, next := range other {
		switch {
		case next.Ordinal() == last:
			continue
		case last == cards.Ace && next.Ordinal() == cards.King, next.Ordinal() == last-cards.Ace:
			count++
			last = next.Ordinal()
		}
	}

	// special case for ace-low.
	var aceLow bool
	if (count == 4 || count == r.max-1) && last == cards.Two && r.ace != nil {
		aceLow = true
		count++
	}

	if count < 5 && count != r.max {
		return 1
	}
	if count > 5 {
		count = 5
	}

	// found it, so save a copy.
	r.counts[4] = count
	r.straight[0] = card
	last = card.Ordinal()
	ix := 1
	for _, next := range other {
		switch {
		case next.Ordinal() == last:
			continue
		case last == cards.Ace && next.Ordinal() == cards.King, next.Ordinal() == last-cards.Ace:
			last = next.Ordinal()
			r.straight[ix] = next
			ix++
		}
		if ix >= count {
			break
		}
	}

	// add the low ace?
	if aceLow && ix < 5 {
		r.straight[ix] = r.ace
	}

	if r.counts[3] > 0 {
		// check if we have a straight flush.
		r.straightFlush = true
		suit := r.flush[0].Suit()
		for iy, c := range r.straight {
			if iy < r.max && c.Suit() != suit {
				if replace := r.findCard(suit, c.Ordinal()); replace != nil {
					r.straight[iy] = replace
				} else {
					r.straightFlush = false
					break
				}
			}
		}
	}

	return count
}

func (r *Ranked) findCard(suit cards.Suit, ordinal cards.Ordinal) *cards.Card {
	for _, c := range r.sorted {
		if c.Suit() == suit && c.Ordinal() == ordinal {
			return c
		}
		if (ordinal == cards.Ace && !c.IsAce()) || (c.Ordinal() != cards.Ace && c.Ordinal() < ordinal) {
			break
		}
	}
	return nil
}

func (r *Ranked) haveStraightFlush() bool {
	count := r.counts[4]
	if (count != 5 && count != r.max) || !r.straightFlush {
		return false
	}

	r.ranked = r.ranked[:count]
	for ix := 0; ix < count; ix++ {
		r.ranked[ix] = r.straight[ix]
	}

	if r.ranked[0].Ordinal() == cards.Ace {
		r.rank = RoyalFlush
	} else {
		r.rank = StraightFlush
	}

	return true
}

func (r *Ranked) haveStraight() bool {
	count := r.counts[4]
	if count != 5 && count != r.max {
		return false
	}

	r.ranked = r.ranked[:count]
	for ix := 0; ix < count; ix++ {
		r.ranked[ix] = r.straight[ix]
	}

	r.rank = Straight
	return true
}
