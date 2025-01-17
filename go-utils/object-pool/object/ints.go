package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Ints is a convenience type for a slice of int values.
// It is not safe for concurrent use by multiple go-routines.
type Ints []int

// DefaultIntsCap is the default capacity for a slice of int values.
const DefaultIntsCap = 8

// ResetInts clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetInts(in Ints, min, max int, cl bool) Ints {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Ints, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Ints) Replace(in Ints) Ints {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultIntsCap)
	n := ResetInts(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of int values is empty or nil.
func (o *Ints) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of int values to json.
func (o *Ints) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Ints) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Int64(int64((*o)[ix]))
	}
}

// Decode decodes the slice of int values from json.
func (o *Ints) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Ints) decodeValue(dec *zjson.Decoder) error {
	if i, ok := dec.Int(); ok {
		*o = append(*o, i)
		return nil
	}
	return dec.Error()
}

// IntsManager is a memory pool object for a slice of integers.
// An IntsManager is not safe for concurrent use by multiple go-routines.
type IntsManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Ints // exposed slice of integers.
}

// Append appends values to the slice.
func (o *IntsManager) Append(v ...int) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty function.
func (o *IntsManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *IntsManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues function.
func (o *IntsManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *IntsManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	if dec.Array(o.Items.decodeValue) {
		return nil
	}
	return dec.Error()
}

// NewIntsProducer instantiates a new IntsManager producer.
func NewIntsProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &IntsManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Ints, 0, minSize),
		}
		return m, m.reset
	})
}

// reset clears the slice of int values.
func (o *IntsManager) reset() {
	o.Items = ResetInts(o.Items, o.minSize, o.maxSize, o.fullClear)
}
