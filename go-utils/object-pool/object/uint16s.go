package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Uint16s is a convenience type for a slice of uint16.
// It is not safe for concurrent use by multiple go-routines.
type Uint16s []uint16

// DefaultUint16sCap is the default capacity for a slice of uint16 values.
const DefaultUint16sCap = 8

// ResetUint16s clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetUint16s(in Uint16s, min, max int, cl bool) Uint16s {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Uint16s, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Uint16s) Replace(in Uint16s) Uint16s {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultUint16sCap)
	n := ResetUint16s(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of uint16 values is empty or nil.
func (o *Uint16s) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of uint16 values to json.
func (o *Uint16s) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Uint16s) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Uint64(uint64((*o)[ix]))
	}
}

// Decode decodes the slice of uint16 values from json.
func (o *Uint16s) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Uint16s) decodeValue(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		*o = append(*o, i)
		return nil
	}
	return dec.Error()
}

// Uint16sManager is a memory pool object for a slice of uint16 values.
// An Uint16sManager is not safe for concurrent use by multiple go-routines.
type Uint16sManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Uint16s // exposed slice of integers.
}

// Append appends values to the slice.
func (o *Uint16sManager) Append(v ...uint16) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty function.
func (o *Uint16sManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *Uint16sManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues function.
func (o *Uint16sManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *Uint16sManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	dec.Array(o.Items.decodeValue)
	return dec.Error()
}

// NewUint16sProducer instantiates a new Uint16sManager producer.
func NewUint16sProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &Uint16sManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Uint16s, 0, minSize),
		}
		return m, m.reset
	})
}

// reset clears the slice of uint16 values.
func (o *Uint16sManager) reset() {
	o.Items = ResetUint16s(o.Items, o.minSize, o.maxSize, o.fullClear)
}
