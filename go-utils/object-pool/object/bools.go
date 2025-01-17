package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Bools is a convenience type for a slice of booleans.
// It is not safe for concurrent use by multiple go-routines.
type Bools []bool

// DefaultBoolsCap is the default capacity for a slice of booleans.
const DefaultBoolsCap = 8

// ResetBools clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetBools(in Bools, min, max int, cl bool) Bools {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Bools, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Bools) Replace(in Bools) Bools {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultBoolsCap)
	n := ResetBools(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of booleans is empty or nil.
func (o *Bools) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of boolenas to json.
func (o *Bools) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Bools) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.IntBool((*o)[ix])
	}
}

// Decode decodes the slice of booleans from json.
func (o *Bools) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Bools) decodeValue(dec *zjson.Decoder) error {
	if b, ok := dec.IntBool(); ok {
		*o = append(*o, b)
		return nil
	}
	return dec.Error()
}

// BoolsManager is a memory pool object for a slice of booleans.
// A BoolsManager is not safe for concurrent use by multiple go-routines.
type BoolsManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Bools // exposed slice if booleans.
}

// Append appends values to the slice.
func (o *BoolsManager) Append(v ...bool) {
	o.Items = append(o.Items, v...)
}

// IsEmpty reimplements the objecter.IsEmpty function.
func (o *BoolsManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *BoolsManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *BoolsManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode reimplements the Objecter.Decode interface.
func (o *BoolsManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	if dec.Array(o.Items.decodeValue) {
		return nil
	}
	return dec.Error()
}

// NewBoolsProducer instantiates a new BoolsManager producer.
func NewBoolsProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &BoolsManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Bools, 0, minSize),
		}
		return m, m.reset
	})
}

// reset clears the array of bools.
func (o *BoolsManager) reset() {
	o.Items = ResetBools(o.Items, o.minSize, o.maxSize, o.fullClear)
}
