package object

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Objects is a convenience type for a slice of memory pool objects.
// It is not safe for concurrent use by multiple go-routines.
type Objects []pool.Objecter

// DefaultObjectsCap is the default capacity for a slice of memory pool objects.
const DefaultObjectsCap = 8

// ResetObjects clears the slice or creates a new one if it doesn't have the requested capacity.
// It first releases all objects from the slice.
func ResetObjects(in Objects, min, max int) Objects {
	in = ReleaseAll(in)
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Objects, 0, min)
	}
	return in
}

// Replace makes a deep copy of the given input into this slice.
// Before copying, it releases all objects from this slice.
func Replace(in, other Objects) Objects {
	if other == nil {
		return ReleaseAll(in)
	}

	l := len(other)
	m := NormalizeSize(l, DefaultObjectsCap)
	n := ResetObjects(in, m, m)[:l]

	for ix := range other {
		if other[ix] != nil {
			n[ix] = other[ix].Clone()
		} else {
			n[ix] = nil
		}
	}
	return n
}

// ReleaseAll releases all objects from the slice and returns the empty result.
func ReleaseAll(in Objects) Objects {
	if in == nil {
		return nil
	}
	for ix := range in {
		if in[ix] != nil {
			in[ix].Release()
			in[ix] = nil
		}
	}
	return in[:0]
}

// IsEmpty returns true if the slice of objects is empty or nil.
func (o *Objects) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of objects to json.
func (o *Objects) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Objects) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Object((*o)[ix])
	}
}

// ObjectsManager is a memory pool object for a slice of memory pool objects.
type ObjectsManager struct {
	minSize int
	maxSize int
	pool.Object
	objectProducer pool.Producer
	Items          Objects // exposed slice of objects.
}

// Append appends objects to the slice.
func (o *ObjectsManager) Append(v ...pool.Objecter) {
	o.Items = append(o.Items, v...)
}

// IsEmpty implements the Objecter.IsEmpty interface.
func (o *ObjectsManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *ObjectsManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *ObjectsManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode reimplements the Objecter.Decode interface.
func (o *ObjectsManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	dec.Array(o.decodeValue)
	return dec.Error()
}

func (o *ObjectsManager) decodeValue(dec *zjson.Decoder) error {
	if n := o.objectProducer.Acquire(); n != nil {
		if dec.Object(n) {
			o.Append(n)
			return nil
		}
		return dec.Error()
	}
	return fmt.Errorf("ObjectsManager.decodeValue: invalid objectProducer")
}

// NewObjectsProducer instantiates a new ObjectsManager producer.
func NewObjectsProducer(minSize, maxSize int, objectProducer pool.Producer) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &ObjectsManager{
			minSize:        minSize,
			maxSize:        maxSize,
			objectProducer: objectProducer,
			Items:          make(Objects, 0, minSize),
		}
		return m, m.reset
	})
}

// reset implements the Objecter.reset function.
func (o *ObjectsManager) reset() {
	o.Items = ResetObjects(ReleaseAll(o.Items), o.minSize, o.maxSize)
}
