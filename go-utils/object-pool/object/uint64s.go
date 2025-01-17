package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Uint64s is a convenience type for a slice of uint64.
// It is not safe for concurrent use by multiple go-routines.
type Uint64s []uint64

// DefaultUint64sCap is the default capacity for a slice of uint64 values.
const DefaultUint64sCap = 8

// ResetUint64s clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetUint64s(in Uint64s, min, max int, cl bool) Uint64s {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Uint64s, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Uint64s) Replace(in Uint64s) Uint64s {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultUint64sCap)
	n := ResetUint64s(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of uint64 values is empty or nil.
func (o *Uint64s) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of uint64 values to json.
func (o *Uint64s) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Uint64s) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Uint64((*o)[ix])
	}
}

// Decode decodes the slice of uint64 values from json.
func (o *Uint64s) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Uint64s) decodeValue(dec *zjson.Decoder) error {
	if i, ok := dec.Uint64(); ok {
		*o = append(*o, i)
		return nil
	}
	return dec.Error()
}

// Uint64sManager is a memory pool object for a slice of uint64 values.
// An Uint64sManager is not safe for concurrent use by multiple go-routines.
type Uint64sManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Uint64s // exposed slice of integers.
}

// Append appends values to the slice.
func (o *Uint64sManager) Append(v ...uint64) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty function.
func (o *Uint64sManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *Uint64sManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues function.
func (o *Uint64sManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *Uint64sManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	dec.Array(o.Items.decodeValue)
	return dec.Error()
}

// NewUint64sProducer instantiates a new Uint64sManager producer.
func NewUint64sProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &Uint64sManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Uint64s, 0, minSize),
		}
		return m, m.reset
	})
}

// reset implements the Objecter.reset function.
func (o *Uint64sManager) reset() {
	o.Items = ResetUint64s(o.Items, o.minSize, o.maxSize, o.fullClear)
}

var clearUInt64s = make([]uint64, 128)
