package poker

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

func (r *Ranked) init(hand cards.Cards) *Ranked {
	r.sorted = cards.SortCardsInto(hand, SortPoker, r.sorted)
	r.max = len(r.sorted)
	r.rank = 0
	r.ace = nil
	r.straightFlush = false
	r.ranked = cards.ResliceCards(r.ranked, 5)
	r.remaining = cards.ResliceCards(r.remaining, 6)
	for ix := range r.counts {
		r.counts[ix] = 0
	}
	for ix := range r.same[0] {
		r.same[0][ix] = nil
		r.same[1][ix] = nil
		r.same[2][ix] = nil
	}
	for ix := range r.flush {
		r.flush[ix] = nil
		r.straight[ix] = nil
	}
	return r
}
