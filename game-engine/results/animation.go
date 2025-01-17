package results

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// EventKind indicates the kind of animation event.
type EventKind uint8

// List of animation kinds.
// Always add new kinds at the end of the list! FE + unit-tests depend on the numeric value!
const (
	ReelAnticipationEvent EventKind = iota + 1
	BombEvent
	SuperEvent
	ShooterEvent
	ReelEvent
	PayoutEvent
	FreeGameEvent
	AwardEvent
	StickyEvent
	ClearEvent
	CascadeEvent
	RefillEvent
	AllPaylinesEvent
	BonusAnticipationEvent
	ChoiceRequestEvent
	FeatureTransitionEvent
)

// String implements the Stringer interface.
func (k EventKind) String() string {
	switch k {
	case ReelAnticipationEvent:
		return "reelanticipation"
	case BombEvent:
		return "bomb"
	case SuperEvent:
		return "super"
	case ShooterEvent:
		return "shooter"
	case ReelEvent:
		return "reel"
	case PayoutEvent:
		return "payout"
	case FreeGameEvent:
		return "freegame"
	case AwardEvent:
		return "award"
	case StickyEvent:
		return "sticky"
	case ClearEvent:
		return "clear"
	case CascadeEvent:
		return "cascade"
	case RefillEvent:
		return "refill"
	case AllPaylinesEvent:
		return "megaways"
	case BonusAnticipationEvent:
		return "bonusanticipation"
	case ChoiceRequestEvent:
		return "choicerequest"
	case FeatureTransitionEvent:
		return "featuretransition"
	default:
		return "[unknown]"
	}
}

// Animator describes the interface for animation events.
type Animator interface {
	Kind() EventKind
	pool.Objecter
}

// Localizer describes the interface for localizing animation events.
type Localizer interface {
	SetMessage(string)
}

// Animations is a convenience type for a slice of animation events.
type Animations []Animator

// PurgeAnimations clears the input slice or creates a new slice if its capacity is too low.
// It will automatically release any animation events in the list.
func PurgeAnimations(list Animations, capacity int) Animations {
	if cap(list) < capacity {
		return make(Animations, 0, capacity)
	}
	return ReleaseAnimations(list)
}

// ReleaseAnimations releases all elements in the given list and returns the cleared input list.
func ReleaseAnimations(list Animations) Animations {
	if list == nil {
		return nil
	}
	for ix := range list {
		list[ix].Release()
		list[ix] = nil
	}
	return list[:0]
}

// Clone makes a deep copy of a slice of animation events.
func (a Animations) Clone() Animations {
	out := make(Animations, len(a))
	for ix := range out {
		out[ix] = a[ix].Clone().(Animator)
	}
	return out
}
