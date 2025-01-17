package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Int64s is a convenience type for a slice of int64 values.
// It is not safe for concurrent use by multiple go-routines.
type Int64s []int64

// DefaultInt64sCap is the default capacity for a slice of int64 values.
const DefaultInt64sCap = 8

// ResetInt64s clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetInt64s(in Int64s, min, max int, cl bool) Int64s {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Int64s, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Int64s) Replace(in Int64s) Int64s {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultInt64sCap)
	n := ResetInt64s(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of int64 values is empty or nil.
func (o *Int64s) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of int64 values to json.
func (o *Int64s) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Int64s) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Int64((*o)[ix])
	}
}

// Decode decodes the slice of int64 values from json.
func (o *Int64s) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Int64s) decodeValue(dec *zjson.Decoder) error {
	if i, ok := dec.Int64(); ok {
		*o = append(*o, i)
		return nil
	}
	return dec.Error()
}

// Int64sManager is a memory pool object for a slice of int64 values.
// An Int64sManager is not safe for concurrent use by multiple go-routines.
type Int64sManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Int64s // exposed slice of integers.
}

// Append appends values to the slice.
func (o *Int64sManager) Append(v ...int64) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty function.
func (o *Int64sManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *Int64sManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues function.
func (o *Int64sManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *Int64sManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	if dec.Array(o.Items.decodeValue) {
		return nil
	}
	return dec.Error()
}

// NewInt64sProducer instantiates a new Int64sManager producer.
func NewInt64sProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &Int64sManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Int64s, 0, minSize),
		}
		return m, m.reset
	})
}

// reset clears the slice of int64 values.
func (o *Int64sManager) reset() {
	o.Items = ResetInt64s(o.Items, o.minSize, o.maxSize, o.fullClear)
}
