package poker

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

func (r *Ranked) findFlush() *Ranked {
	if r.max <= 5 {
		r.testFlush(r.sorted[0], r.sorted[1:])
	} else {
		for ix := 0; ix <= r.max-5; {
			ix += r.testFlush(r.sorted[ix], r.sorted[ix+1:])
		}
	}
	return r
}

func (r *Ranked) testFlush(card *cards.Card, other cards.Cards) int {
	count := 1
	for _, next := range other {
		if next.Suit() == card.Suit() {
			count++
		}
	}

	if count < 5 && count != r.max {
		return 1
	}
	if count > 5 {
		count = 5
	}

	r.counts[3] = count
	r.flush[0] = card
	ix := 1
	for _, c := range other {
		if c.Suit() == card.Suit() {
			r.flush[ix] = c
			ix++
			if ix >= count {
				break
			}
		}
	}

	return count
}

func (r *Ranked) haveFlush() bool {
	count := r.counts[3]
	if count != 5 && count != r.max {
		return false
	}

	r.ranked = r.ranked[:count]
	for ix := 0; ix < count; ix++ {
		r.ranked[ix] = r.flush[ix]
	}

	r.rank = Flush
	return true
}
