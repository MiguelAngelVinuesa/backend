package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Uint8s is a convenience type for a slice of uint8.
// It is not safe for concurrent use by multiple go-routines.
type Uint8s []uint8

// DefaultUint8sCap is the default capacity for a slice of uint8 values.
const DefaultUint8sCap = 8

// ResetUint8s clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetUint8s(in Uint8s, min, max int, cl bool) Uint8s {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Uint8s, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Uint8s) Replace(in Uint8s) Uint8s {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultUint8sCap)
	n := ResetUint8s(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of uint8 values is empty or nil.
func (o *Uint8s) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of uint8 values to json.
func (o *Uint8s) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Uint8s) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Uint64(uint64((*o)[ix]))
	}
}

// Decode decodes the slice of uint8 values from json.
func (o *Uint8s) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Uint8s) decodeValue(dec *zjson.Decoder) error {
	if i, ok := dec.Uint8(); ok {
		*o = append(*o, i)
		return nil
	}
	return dec.Error()
}

// Uint8sManager is a memory pool object for a slice of uint8 values.
// An Uint8sManager is not safe for concurrent use by multiple go-routines.
type Uint8sManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Uint8s // exposed slice of uint8 values.
}

// Append appends values to the slice.
func (o *Uint8sManager) Append(v ...uint8) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty function.
func (o *Uint8sManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *Uint8sManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues function.
func (o *Uint8sManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *Uint8sManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	dec.Array(o.Items.decodeValue)
	return dec.Error()
}

// NewUint8sProducer instantiates a new Uint8sManager producer.
func NewUint8sProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &Uint8sManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Uint8s, 0, minSize),
		}
		return m, m.reset
	})
}

// reset clears the slice of uint8 values.
func (o *Uint8sManager) reset() {
	o.Items = ResetUint8s(o.Items, o.minSize, o.maxSize, o.fullClear)
}
