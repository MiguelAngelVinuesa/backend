package results

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// PenaltyKind represents the kind of penalty.
type PenaltyKind uint8

// List of penalty kinds.
// Always add new elements after the end of the list as FE depends on the values!
const (
	// SlotReduction is a penalty with a fixed reduction factor.
	SlotReduction PenaltyKind = iota + 1
	// SlotDivision is a penalty that divides the total payout factor by a percentage.
	SlotDivision
)

// String implements the Stringer interface.
func (k PenaltyKind) String() string {
	switch k {
	case SlotReduction:
		return "slot reduction"
	case SlotDivision:
		return "slot division"
	default:
		return "[unknown]"
	}
}

// Penalty is the interface for a penalty.
type Penalty interface {
	AsPayout() Payout
	Kind() PenaltyKind
	Reduce() float64
	Divide() float64
	SetFactor(float64)
	SetMessage(string)
	pool.Objecter
}

// Penalties is a convenience type for a slice of penalties.
type Penalties []Penalty

// PurgePenalties clears the input slice or creates a new slice if its capacity is too low.
// It will automatically release any payouts in the list.
func PurgePenalties(list Penalties, capacity int) Penalties {
	if cap(list) < capacity {
		ReleasePenalties(list)
		return make(Penalties, 0, capacity)
	}
	return ReleasePenalties(list)
}

// ReleasePenalties releases all elements in the given list and returns the cleared input list.
func ReleasePenalties(list Penalties) Penalties {
	if list == nil {
		return nil
	}
	for ix := range list {
		list[ix].Release()
		list[ix] = nil
	}
	return list[:0]
}
