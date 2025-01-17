package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// NewRoundFlag instantiates a new round flag metrics from the memory pool.
func NewRoundFlag(id int, name string) *RoundFlag {
	r := roundFlagPool.Acquire().(*RoundFlag)
	r.ID = id
	r.Name = name
	return r
}

// Merge merges with the given input.
func (r *RoundFlag) Merge(other *RoundFlag) {
	r.Counts.Merge(other.Counts)
	r.CountsFinal.Merge(other.CountsFinal)
}

// Increase updates the round flag metrics.
func (f *RoundFlag) Increase(value int, final bool) {
	f.Counts.Increase(int64(value))
	if final {
		f.CountsFinal.Increase(int64(value))
	}
}

// RoundFlag contains the metrics for a round flag.
type RoundFlag struct {
	ID          int          `json:"id"`
	Counts      *MinMaxInt64 `json:"counts"`
	CountsFinal *MinMaxInt64 `json:"countsFinal"`
	Name        string       `json:"name"`
	pool.Object
}

var roundFlagPool = pool.NewProducer(func() (pool.Objecter, func()) {
	a := &RoundFlag{
		Counts:      AcquireMinMaxInt64(),
		CountsFinal: AcquireMinMaxInt64(),
	}
	return a, a.reset
})

func (f *RoundFlag) reset() {
	if f != nil {
		f.ResetData()
		f.ID = 0
		f.Name = ""
	}
}

// ResetData resets the found flag metrics to initial state.
func (f *RoundFlag) ResetData() {
	f.Counts.ResetData()
	f.CountsFinal.ResetData()
}

// Equals is used internally for unit-tests!
func (f *RoundFlag) Equals(other *RoundFlag) bool {
	return f.ID == other.ID &&
		f.Counts.Equals(other.Counts) &&
		f.CountsFinal.Equals(other.CountsFinal)
}

// RoundFlags is a convenience type for a slice of round flags.
type RoundFlags []*RoundFlag

// ReleaseRoundFlags releases all round flags in the slice and returns the empty slice.
func ReleaseRoundFlags(list RoundFlags) RoundFlags {
	if list == nil {
		return nil
	}
	for ix := range list {
		if f := list[ix]; f != nil {
			f.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeRoundFlags returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
// It automatically takes care of releasing the round flags.
func PurgeRoundFlags(list RoundFlags, capacity int) RoundFlags {
	list = ReleaseRoundFlags(list)
	if cap(list) < capacity {
		return make(RoundFlags, 0, capacity)
	}
	return list[:0]
}
