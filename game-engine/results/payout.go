package results

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// PayoutKind represents the kind of result.
type PayoutKind uint8

// List of payout kinds.
// Always add new elements after the end of the list as FE depends on the values!
const (
	// SlotWinline is a payout from a winline.
	SlotWinline PayoutKind = iota + 1
	// SlotWilds is a payout from scattered wild symbols.
	SlotWilds
	// SlotScatters is a payout from a scatter symbol.
	SlotScatters
	// SlotBonusSymbol is a payout from the bonus symbol in free spins.
	SlotBonusSymbol
	// SlotCluster is a payout for a cluster of the same symbol.
	SlotCluster
	// SlotBonusGame is a payout for a bonus game.
	SlotBonusGame
	// SlotSuperShape is a payout for a scatter payout from a super shape.
	SlotSuperShape
	// SlotMultiplier is a payout from a multiplier bonus game.
	SlotMultiplier
	// SlotBombScatters is a payout from a scatter symbol after a bomb explosion.
	SlotBombScatters
	// PlayerChoice is a payout from a player choice (e.g. double or nothing).
	PlayerChoice
	// SlotReducePenalty is a penalty with a fixed reduction factor.
	SlotReducePenalty
	// SlotDividePenalty is a penalty that divides the total payout factor by a percentage.
	SlotDividePenalty
)

// String implements the Stringer interface.
func (k PayoutKind) String() string {
	switch k {
	case SlotWinline:
		return "slot winline"
	case SlotWilds:
		return "slot wilds"
	case SlotScatters:
		return "slot scatters"
	case SlotBonusSymbol:
		return "slot bonus symbol"
	case SlotCluster:
		return "slot cluster"
	case SlotBonusGame:
		return "slot bonus game"
	case SlotSuperShape:
		return "slot super shape"
	case SlotMultiplier:
		return "slot multiplier"
	case SlotBombScatters:
		return "slot bomb scatters"
	case PlayerChoice:
		return "player choice"
	case SlotReducePenalty:
		return "slot reduce penalty"
	case SlotDividePenalty:
		return "slot divide penalty"
	default:
		return "[unknown]"
	}
}

// Payout is the interface for a payout.
type Payout interface {
	Kind() PayoutKind
	Total() float64
	pool.Objecter
}

// Payouts is a convenience type for a slice of payouts.
type Payouts []Payout

// PurgePayouts clears the input slice or creates a new slice if its capacity is too low.
// It will automatically release any payouts in the list.
func PurgePayouts(list Payouts, capacity int) Payouts {
	if cap(list) < capacity {
		ReleasePayouts(list)
		return make(Payouts, 0, capacity)
	}
	return ReleasePayouts(list)
}

// ReleasePayouts releases all elements in the given list and returns the cleared input list.
func ReleasePayouts(list Payouts) Payouts {
	if list == nil {
		return nil
	}
	for ix := range list {
		list[ix].Release()
		list[ix] = nil
	}
	return list[:0]
}
