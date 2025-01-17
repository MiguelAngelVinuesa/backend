package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// NewAction instantiates a new action metrics from the memory pool.
func NewAction(id int, name, kind, config string) *Action {
	a := actionPool.Acquire().(*Action)
	a.ID = id
	a.Name = name
	a.Kind = kind
	a.Config = config
	return a
}

// Merge merges with the given input.
func (a *Action) Merge(other *Action) {
	a.TotalCount += other.TotalCount
	a.FirstCount += other.FirstCount
	a.SecondCount += other.SecondCount
	a.FreeCount += other.FreeCount
	a.FreeSecondCount += other.FreeSecondCount
	a.SuperCount += other.SuperCount
	a.RefillCount += other.RefillCount
	a.TotalTriggered += other.TotalTriggered
	a.FirstTriggered += other.FirstTriggered
	a.SecondTriggered += other.SecondTriggered
	a.FreeTriggered += other.FreeTriggered
	a.FreeSecondTriggered += other.FreeSecondTriggered
	a.SuperTriggered += other.SuperTriggered
	a.RefillTriggered += other.RefillTriggered
}

func (a *Action) IncreaseFirst(trigger bool) {
	a.TotalCount++
	a.FirstCount++
	if trigger {
		a.TotalTriggered++
		a.FirstTriggered++
	}
}

func (a *Action) IncreaseSecond(trigger bool) {
	a.TotalCount++
	a.SecondCount++
	if trigger {
		a.TotalTriggered++
		a.SecondTriggered++
	}
}

func (a *Action) IncreaseFree(trigger bool) {
	a.TotalCount++
	a.FreeCount++
	if trigger {
		a.TotalTriggered++
		a.FreeTriggered++
	}
}

func (a *Action) IncreaseSecondFree(trigger bool) {
	a.TotalCount++
	a.FreeSecondCount++
	if trigger {
		a.TotalTriggered++
		a.FreeSecondTriggered++
	}
}

func (a *Action) IncreaseSuper(trigger bool) {
	a.TotalCount++
	a.SuperCount++
	if trigger {
		a.TotalTriggered++
		a.SuperTriggered++
	}
}

func (a *Action) IncreaseRefill(trigger bool) {
	a.TotalCount++
	a.RefillCount++
	if trigger {
		a.TotalTriggered++
		a.RefillTriggered++
	}
}

type Action struct {
	ID                  int    `json:"id,omitempty"`
	TotalCount          uint64 `json:"totalCount,omitempty"`
	FirstCount          uint64 `json:"firstCount,omitempty"`
	SecondCount         uint64 `json:"secondCount,omitempty"`
	FreeCount           uint64 `json:"freeCount,omitempty"`
	FreeSecondCount     uint64 `json:"freeSecondCount,omitempty"`
	SuperCount          uint64 `json:"superCount,omitempty"`
	RefillCount         uint64 `json:"refillCount,omitempty"`
	TotalTriggered      uint64 `json:"totalTriggered,omitempty"`
	FirstTriggered      uint64 `json:"firstTriggered,omitempty"`
	SecondTriggered     uint64 `json:"secondTriggered,omitempty"`
	FreeTriggered       uint64 `json:"freeTriggered,omitempty"`
	FreeSecondTriggered uint64 `json:"freeTriggeredCount,omitempty"`
	SuperTriggered      uint64 `json:"superTriggered,omitempty"`
	RefillTriggered     uint64 `json:"refillTriggered,omitempty"`
	Name                string `json:"name,omitempty"`
	Kind                string `json:"kind,omitempty"`
	Config              string `json:"config,omitempty"`
	pool.Object
}

var actionPool = pool.NewProducer(func() (pool.Objecter, func()) {
	a := &Action{}
	return a, a.reset
})

func (a *Action) reset() {
	if a != nil {
		a.ResetData()

		a.ID = 0
		a.Name = ""
		a.Kind = ""
		a.Config = ""
	}
}

// ResetData resets the action metrics to initial state.
func (a *Action) ResetData() {
	a.TotalCount = 0
	a.FirstCount = 0
	a.SecondCount = 0
	a.FreeCount = 0
	a.FreeSecondCount = 0
	a.SuperCount = 0
	a.RefillCount = 0
	a.TotalTriggered = 0
	a.FirstTriggered = 0
	a.SecondTriggered = 0
	a.FreeTriggered = 0
	a.FreeSecondTriggered = 0
	a.SuperTriggered = 0
	a.RefillTriggered = 0
}

// Equals is used internally for unit tests!
func (a *Action) Equals(other *Action) bool {
	return a.ID == other.ID &&
		a.TotalCount == other.TotalCount &&
		a.FirstCount == other.FirstCount &&
		a.SecondCount == other.SecondCount &&
		a.FreeCount == other.FreeCount &&
		a.FreeSecondCount == other.FreeSecondCount &&
		a.SuperCount == other.SuperCount &&
		a.RefillCount == other.RefillCount &&
		a.TotalTriggered == other.TotalTriggered &&
		a.FirstTriggered == other.FirstTriggered &&
		a.SecondTriggered == other.SecondTriggered &&
		a.FreeTriggered == other.FreeTriggered &&
		a.FreeSecondTriggered == other.FreeSecondTriggered &&
		a.SuperTriggered == other.SuperTriggered &&
		a.RefillTriggered == other.RefillTriggered &&
		a.Name == other.Name &&
		a.Kind == other.Kind &&
		a.Config == other.Config
}

// Actions is a convenience type for a slice of actions.
type Actions []*Action

// ReleaseActions releases all actions and returns an empty slice.
func ReleaseActions(list Actions) Actions {
	if list == nil {
		return nil
	}
	for ix := range list {
		if a := list[ix]; a != nil {
			a.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeActions resets the slice to zero length or returns a new one if its capacity is less than requested.
func PurgeActions(list Actions, capacity int) Actions {
	list = ReleaseActions(list)
	if cap(list) < capacity {
		return make(Actions, 0, capacity)
	}
	return list
}
