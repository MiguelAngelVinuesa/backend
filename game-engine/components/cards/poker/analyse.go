package poker

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

func (r *Ranked) analyse() *Ranked {
	if r.haveStraightFlush() || r.haveSameKind(4) || r.haveTwiceSame(3) || r.haveFlush() ||
		r.haveStraight() || r.haveSameKind(3) || r.haveTwiceSame(2) || r.haveSameKind(2) {
		return r
	}

	// high card
	r.ranked = r.ranked[:1]
	r.ranked[0] = r.sorted[0]
	r.rank = HighCard
	return r
}

func (r *Ranked) getRemaining() *Ranked {
	for _, c := range r.sorted {
		if r.isRemaining(c) {
			r.remaining = append(r.remaining, c)
		}
	}
	return r
}

func (r *Ranked) isRemaining(card *cards.Card) bool {
	for _, c := range r.ranked {
		if card == c {
			return false
		}
	}
	return true
}
