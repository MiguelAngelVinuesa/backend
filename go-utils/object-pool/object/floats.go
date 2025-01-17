package object

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Floats is a convenience type for a slice of floating point values.
// It is not safe for concurrent use by multiple go-routines.
type Floats []float64

// DefaultFloatsCap is the default capacity for a slice of floating point values.
const DefaultFloatsCap = 8

// ResetFloats clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetFloats(in Floats, min, max int, cl bool) Floats {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Floats, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Floats) Replace(in Floats) Floats {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultFloatsCap)
	n := ResetFloats(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of floating point values is empty or nil.
func (o *Floats) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of floating point values to json.
func (o *Floats) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
// By default, it uses format 'g' and precision 4.
func (o *Floats) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Float((*o)[ix], 'g', -1)
	}
}

// EncodeFormat encodes the slice of floating point values to json using the given format & precision.
func (o *Floats) EncodeFormat(enc *zjson.Encoder, format byte, prec int) {
	enc.StartArray()
	for ix := range *o {
		enc.Float((*o)[ix], format, prec)
	}
	enc.EndArray()
}

// Decode decodes the slice of floating point values from json.
func (o *Floats) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Floats) decodeValue(dec *zjson.Decoder) error {
	if f, ok := dec.Float(); ok {
		*o = append(*o, f)
		return nil
	}
	return dec.Error()
}

// FloatsManager is a memory pool object for a slice of floating point values.
// A FloatsManager is not safe for concurrent use by multiple go-routines.
type FloatsManager struct {
	fullClear bool
	format    byte
	prec      int
	minSize   int
	maxSize   int
	pool.Object
	Items Floats // exposed slice of floats.
}

// Append appends values to the slice.
func (o *FloatsManager) Append(v ...float64) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty function.
func (o *FloatsManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *FloatsManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *FloatsManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *FloatsManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	dec.Array(o.Items.decodeValue)
	return dec.Error()
}

// NewFloatsProducer instantiates a new FloatsManager producer.
func NewFloatsProducer(minSize, maxSize int, fullClear bool, format byte, prec int) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &FloatsManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			format:    format,
			prec:      prec,
			Items:     make(Floats, 0, minSize),
		}
		return m, m.reset
	})
}

// reset clears the array of floats.
func (o *FloatsManager) reset() {
	o.Items = ResetFloats(o.Items, o.minSize, o.maxSize, o.fullClear)
}
