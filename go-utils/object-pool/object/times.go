package object

import (
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// Times is a convenience type for a slice of timestamps.
// The timestamps are stored as milliseconds since Unix Epoch.
// Timestamps should be converted to UTC before being used.
// It is not safe for concurrent use by multiple go-routines.
type Times []int64

// DefaultTimesCap is the default capacity for a slice of timestamps.
const DefaultTimesCap = 8

// ResetTimes clears the slice or creates a new one if it doesn't have the requested capacity.
func ResetTimes(in Times, min, max int, cl bool) Times {
	if cap(in) < min || (max > min && cap(in) > max) {
		return make(Times, 0, min)
	}
	if cl {
		clear(in)
	}
	return in[:0]
}

// Replace the slice with the given input, overwriting any previous content.
func (o *Times) Replace(in Times) Times {
	if in == nil {
		return nil
	}
	l := len(in)
	m := NormalizeSize(l, DefaultTimesCap)
	n := ResetTimes(*o, m, m, false)[:l]
	copy(n, in)
	return n
}

// IsEmpty returns true if the slice of timestamps is empty or nil.
func (o *Times) IsEmpty() bool {
	return len(*o) == 0
}

// Encode encodes the slice of timestamps to json.
// The timestamps are encoded as milliseconds since Unix Epoch.
func (o *Times) Encode(enc *zjson.Encoder) {
	enc.Array(o)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *Times) EncodeValues(enc *zjson.Encoder) {
	for ix := range *o {
		enc.Int64((*o)[ix])
	}
}

// Decode decodes the slice of timestamps from json.
// The timestamps are decoded as milliseconds since Unix Epoch.
func (o *Times) Decode(dec *zjson.Decoder) error {
	*o = (*o)[:0]
	if dec.Array(o.decodeValue) {
		return nil
	}
	return dec.Error()
}

func (o *Times) decodeValue(dec *zjson.Decoder) error {
	if i, ok := dec.Int64(); ok {
		*o = append(*o, i)
		return nil
	}
	return dec.Error()
}

// TimesManager is a memory pool object for a slice of timestamps.
// Timestamps are stored as Unix epoch milliseconds.
// Timestamps before 1900 are considered invalid and are automatically reset to zero.
// Timestamps will be converted to UTC before being appended.
// Returned timestamps are alwayus in UTC.
// A TimesManager is not safe for concurrent use by multiple go-routines.
type TimesManager struct {
	fullClear bool
	minSize   int
	maxSize   int
	pool.Object
	Items Times // exposed slice of Unix millisecond timestamps.
}

// Append appends values to the slice.
func (o *TimesManager) Append(v ...time.Time) {
	for ix := range v {
		if u := v[ix].UTC().UnixMilli(); u > cutoffMilli {
			o.Items = append(o.Items, u)
		} else {
			o.Items = append(o.Items, cutoffMilli)
		}
	}
}

// Value retrieves the indicated timestamp value.
// If the given index is out of range, the nil timestamp will be returned.
func (o *TimesManager) Value(ix int) time.Time {
	if ix >= 0 && ix < len(o.Items) && o.Items[ix] > cutoffMilli {
		return time.UnixMilli(o.Items[ix]).UTC()
	}
	return time.Time{}
}

// IsEmpty implements the objecter.IsEmpty function.
func (o *TimesManager) IsEmpty() bool {
	return o == nil || len(o.Items) == 0
}

// Encode reimplements the Objecter.Encode interface.
func (o *TimesManager) Encode(enc *zjson.Encoder) {
	enc.Array(&o.Items)
}

// EncodeValues implements the zjson.ArrayEncoder.EncodeValues interface.
func (o *TimesManager) EncodeValues(enc *zjson.Encoder) {
	o.Items.EncodeValues(enc)
}

// Decode implements the Objecter.Decode interface.
func (o *TimesManager) Decode(dec *zjson.Decoder) error {
	o.reset()
	if dec.Array(o.Items.decodeValue) {
		return nil
	}
	return dec.Error()
}

// NewTimesProducer instantiates a new TimesManager producer.
func NewTimesProducer(minSize, maxSize int, fullClear bool) pool.Producer {
	return pool.NewProducer(func() (pool.Objecter, func()) {
		m := &TimesManager{
			fullClear: fullClear,
			minSize:   minSize,
			maxSize:   maxSize,
			Items:     make(Times, 0, minSize),
		}
		return m, m.reset
	})
}

var cutoffMilli = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

// reset clears the slice of timestamps.
func (o *TimesManager) reset() {
	o.Items = ResetTimes(o.Items, o.minSize, o.maxSize, o.fullClear)
}
