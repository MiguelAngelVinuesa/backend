package slots

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireStickyChoice instantiates a new sticky symbol choice from the memory pool.
func AcquireStickyChoice(symbol utils.Index, spin *Spin) *StickyChoice {
	c := stickyChoiceProducer.Acquire().(*StickyChoice)
	c.Symbol = symbol

	c.Positions = utils.PurgeUInt8s(c.Positions, 8)
	for ix := range spin.indexes {
		if spin.indexes[ix] == symbol {
			c.Positions = append(c.Positions, uint8(ix))
		}
	}

	return c
}

// IsEmpty implements the zjson encoder interface.
func (c *StickyChoice) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson encoder interface.
func (c *StickyChoice) EncodeFields(enc *zjson.Encoder) {
	enc.Uint16Field("symbol", uint16(c.Symbol))

	enc.StartArrayField("positions")
	for ix := range c.Positions {
		enc.Uint64(uint64(c.Positions[ix]))
	}
	enc.EndArray()
}

// DecodeField implements the zjson decoder interface.
func (c *StickyChoice) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i uint64

	if string(key) == "symbol" {
		if i, ok = dec.Uint64(); ok {
			c.Symbol = utils.Index(i)
		}
	} else if string(key) == "positions" {
		c.Positions = utils.PurgeUInt8s(c.Positions, 8)
		ok = dec.Array(c.decodePositions)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (c *StickyChoice) decodePositions(dec *zjson.Decoder) error {
	if i, ok := dec.Uint8(); ok {
		c.Positions = append(c.Positions, i)
		return nil
	}
	return dec.Error()
}

// StickyChoice represents a symbol which can be chosen as the sticky symbol.
// Keep fields ordered by ascending SizeOf().
type StickyChoice struct {
	Symbol    utils.Index  `json:"symbol"`
	Positions utils.UInt8s `json:"positions"`
	pool.Object
}

// stickyChoiceProducer is the memory pool for spin payouts.
// Make sure to initialize all slices appropriately!
var stickyChoiceProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	c := &StickyChoice{
		Positions: make(utils.UInt8s, 0, 16),
	}
	return c, c.reset
})

// reset clears the sticky choices.
func (c *StickyChoice) reset() {
	if c != nil {
		c.Symbol = utils.NullIndex
		c.Positions = c.Positions[:0]
	}
}

// Equals is used internally for unit-tests!
func (c *StickyChoice) Equals(other *StickyChoice) bool {
	return c.Symbol == other.Symbol && reflect.DeepEqual(c.Positions, other.Positions)
}

// StickyChoices is a convenience type for a slice of sticky choices.
type StickyChoices []*StickyChoice

// ReleaseStickyChoices releases all sticky choices and returns an empty slice.
func ReleaseStickyChoices(list StickyChoices) StickyChoices {
	if list == nil {
		return nil
	}
	for ix := range list {
		if r := list[ix]; r != nil {
			r.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeStickyChoices returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeStickyChoices(list StickyChoices, capacity int) StickyChoices {
	list = ReleaseStickyChoices(list)
	if cap(list) < capacity {
		return make(StickyChoices, 0, capacity)
	}
	return list
}
