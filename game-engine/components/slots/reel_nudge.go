package slots

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NudgeLocation indicates which part of a reel can nudge.
type NudgeLocation uint8

const (
	// NudgeTop indicates the reel can nudge at the top.
	NudgeTop = iota + 1
	// NudgeBottom indicates the reel can nudge at the bottom.
	NudgeBottom
	// NudgeVertical indicates the reel can nudge both at the top and the bottom.
	NudgeVertical
)

func AcquireReelNudge(teaser bool, reel, size uint8, symbol utils.Index, location NudgeLocation) *ReelNudge {
	n := reelNudgeProducer.Acquire().(*ReelNudge)
	n.teaser = teaser
	n.reel = reel
	n.size = size
	n.symbol = symbol
	n.location = location
	return n
}

// EncodeFields implements the zjson encoder interface.
func (n *ReelNudge) EncodeFields(enc *zjson.Encoder) {
	enc.IntBoolFieldOpt("teaser", n.teaser)
	enc.Uint8Field("reel", n.reel)
	enc.Uint8Field("size", n.size)
	enc.Uint16Field("symbol", uint16(n.symbol))
	enc.Uint8Field("location", uint8(n.location))
}

// DecodeField implements the zjson decoder interface.
func (r *ReelNudge) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i8 uint8

	if string(key) == "teaser" {
		var b bool
		if b, ok = dec.IntBool(); ok {
			r.teaser = b
		}
	} else if string(key) == "reel" {
		if i8, ok = dec.Uint8(); ok {
			r.reel = i8
		}
	} else if string(key) == "size" {
		if i8, ok = dec.Uint8(); ok {
			r.size = i8
		}
	} else if string(key) == "symbol" {
		var i16 uint16
		if i16, ok = dec.Uint16(); ok {
			r.symbol = utils.Index(i16)
		}
	} else if string(key) == "location" {
		if i8, ok = dec.Uint8(); ok {
			r.location = NudgeLocation(i8)
		}
	} else {
		return fmt.Errorf("ReelNudge.DecodeField: invalid field encountered [%s]", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// ReelNudge contains the details for a reel nudge.
type ReelNudge struct {
	teaser   bool
	reel     uint8
	size     uint8
	symbol   utils.Index
	location NudgeLocation
	pool.Object
}

var reelNudgeProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &ReelNudge{}
	return e, e.reset
})

// reset clears the reel nudge.
func (e *ReelNudge) reset() {
	if e != nil {
		e.teaser = false
		e.reel = 0
		e.size = 0
		e.symbol = 0
		e.location = 0
	}
}

// DeepEqual is used internally for unit tests.
func (n *ReelNudge) DeepEqual(other *ReelNudge) bool {
	return n.teaser == other.teaser &&
		n.reel == other.reel &&
		n.size == other.size &&
		n.symbol == other.symbol &&
		n.location == other.location
}

// ReelNudges is a convenience type for a slice of reel nudges.
type ReelNudges []*ReelNudge

// ReleaseReelNudges releases any reel nudges from the slice and returns an empty slice.
func ReleaseReelNudges(in ReelNudges) ReelNudges {
	if len(in) == 0 {
		return in
	}

	for ix := range in {
		in[ix].Release()
		in[ix] = nil
	}

	return in[:0]
}

// DeepEqual is used internally for unit tests.
func (n ReelNudges) DeepEqual(other ReelNudges) bool {
	if len(n) != len(other) {
		return false
	}
	for ix := range n {
		if !n[ix].DeepEqual(other[ix]) {
			return false
		}
	}
	return true
}
