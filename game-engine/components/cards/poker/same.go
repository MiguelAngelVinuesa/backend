package poker

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

func (r *Ranked) findSame() *Ranked {
	for ix := 0; ix < r.max-1; {
		ix += r.testSame(r.sorted[ix], r.sorted[ix+1:])
	}
	return r
}

func (r *Ranked) testSame(card *cards.Card, other cards.Cards) int {
	count := 1
	o := card.Ordinal()
	for _, next := range other {
		if next.Ordinal() != o {
			break
		}
		count++
	}

	if count == 1 {
		return 1
	}

	switch {
	case r.counts[0] == 0:
		r.saveSeries(0, count, card, other)
	case r.counts[1] == 0:
		r.saveSeries(1, count, card, other)
	case r.counts[2] == 0:
		r.saveSeries(2, count, card, other)
	}

	return count
}

func (r *Ranked) saveSeries(index, count int, card *cards.Card, other cards.Cards) {
	r.counts[index] = count
	r.same[index][0] = card
	for ix := 1; ix < count; ix++ {
		r.same[index][ix] = other[ix-1]
	}
}

func (r *Ranked) haveSameKind(count int) bool {
	first := r.counts[0] == count

	var series int
	if !first {
		if r.counts[1] != count {
			return false
		}
		series = 1
	}

	r.ranked = r.ranked[:count]
	for ix := 0; ix < count; ix++ {
		r.ranked[ix] = r.same[series][ix]
	}

	switch count {
	case 4:
		r.rank = FourOfAKind
	case 3:
		r.rank = ThreeOfAKind
	case 2:
		r.rank = OnePair
	}
	return true
}

func (r *Ranked) haveTwiceSame(count int) bool {
	var first, second int

	switch {
	case r.counts[0] == count:
		if r.counts[1] < 2 {
			return false
		}
		second = 1

	case r.counts[1] == count:
		if r.counts[0] != 2 {
			return false
		}
		first = 1

	case r.counts[2] == count:
		first = 2

	default:
		return false
	}

	r.ranked = r.ranked[:count+2]
	for ix := 0; ix < count; ix++ {
		r.ranked[ix] = r.same[first][ix]
	}
	for ix := 0; ix < 2; ix++ {
		r.ranked[count+ix] = r.same[second][ix]
	}

	if count == 3 {
		r.rank = FullHouse
	} else {
		r.rank = TwoPair
	}
	return true
}
