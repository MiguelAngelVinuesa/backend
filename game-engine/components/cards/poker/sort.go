package poker

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
)

// SortPoker sorts the cards in an order useful for ranking the hand.
func SortPoker(hand cards.Cards, i, j int) bool {
	switch {
	case hand[i].Ordinal() == hand[j].Ordinal():
		return hand[i].Suit() > hand[j].Suit()
	case hand[i].Ordinal() == cards.Ace:
		return true
	case hand[j].Ordinal() == cards.Ace:
		return false
	default:
		return hand[i].Ordinal() > hand[j].Ordinal()
	}
}
