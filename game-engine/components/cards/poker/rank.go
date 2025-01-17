package poker

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

// Rank represents a poker hand.
type Rank int8

// Standard poker rankings.
const (
	HighCard Rank = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// String implements the Stringer interface.
func (r Rank) String() string {
	switch r {
	case HighCard:
		return "High card"
	case OnePair:
		return "One pair"
	case TwoPair:
		return "Two pair"
	case ThreeOfAKind:
		return "Three of a kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full house"
	case FourOfAKind:
		return "Four of a kind"
	case StraightFlush:
		return "Straight flush"
	case RoyalFlush:
		return "Royal flush"
	default:
		return "No rank"
	}
}

// Ranked contains the details of a set of ranked cards.
type Ranked struct {
	sorted        cards.Cards
	max           int
	rank          Rank
	ranked        cards.Cards
	remaining     cards.Cards
	ace           *cards.Card
	same          [3][4]*cards.Card
	flush         [5]*cards.Card
	straight      [5]*cards.Card
	counts        [5]int
	straightFlush bool
}

var rankedPool = sync.Pool{
	New: func() interface{} {
		return &Ranked{}
	},
}

// NewRanked returns a Ranked object from the memory pool without ranking the hand.
func NewRanked(hand cards.Cards) *Ranked {
	return rankedPool.Get().(*Ranked).init(hand)
}

// RankHand returns a Ranked object from the memory pool and ranks the hand.
// If there is less than 3 cards in the hand, or more than 7, no ranking is done.
// The function assumes the hand is from a single deck (e.g. no duplicates) and the deck has no jokers.
func RankHand(hand cards.Cards) *Ranked {
	r := rankedPool.Get().(*Ranked).init(hand)
	if r.max < 3 || r.max > 7 {
		return r
	}
	return r.findSame().findFlush().findStraight().analyse().getRemaining()
}

// Release puts a Ranked object back into the memory pool.
func (r *Ranked) Release() {
	if r != nil {
		rankedPool.Put(r)
	}
}

// Reload loads the hand with the given cards and ranks it again.
func (r *Ranked) Reload(hand cards.Cards) *Ranked {
	r.init(hand)
	if r.max < 3 || r.max > 7 {
		return r
	}
	return r.findSame().findFlush().findStraight().analyse().getRemaining()
}

// Rank returns the calculated ank of the cards.
func (r *Ranked) Rank() Rank {
	return r.rank
}

// Ranked returns the set of cards that are part of the calculated rank.
func (r *Ranked) Ranked() cards.Cards {
	return r.ranked
}

// Remaining returns the set of cards that are not part of the calculated rank.
func (r *Ranked) Remaining() cards.Cards {
	return r.remaining
}
