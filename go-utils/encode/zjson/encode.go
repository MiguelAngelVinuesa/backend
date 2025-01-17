package zjson

import (
	"bytes"
	"strconv"
	"sync"
	"time"
)

// ObjectEncoder is the interface for encoding objects.
type ObjectEncoder interface {
	IsEmpty() bool
	EncodeFields(enc *Encoder)
}

// ArrayEncoder is the interface for encoding arrays.
type ArrayEncoder interface {
	IsEmpty() bool
	EncodeValues(enc *Encoder)
}

// Encoder implements a zero-allocation json encoder.
// It is more than 2x faster than the standard json marshaller and also faster than goccy/go-json,
// however it is not a drop-in replacement and requires handwritten code to marshall structures.
// Encoder does no error checking, os it is quite possible to generate invalid json output.
// Encoder is not safe for use across concurrent go-routines.
type Encoder struct {
	bytes.Buffer
}

// AcquireEncoder instantiates a zero allocation json encoder from the memory pool.
func AcquireEncoder(capacity int) *Encoder {
	e := encoderPool.Get().(*Encoder)
	n := FixBufferSize(capacity)
	if c := e.Cap(); c < n {
		e.Grow(n)
	}
	return e
}

// Release returns the encoder to the memory pool.
func (e *Encoder) Release() {
	if e != nil {
		e.Reset()
		encoderPool.Put(e)
	}
}

// Bytes returns the encoded json data.
func (e *Encoder) Bytes() []byte {
	if l := e.Len(); l > 0 {
		return e.Buffer.Bytes()[:l-1]
	}
	return e.Buffer.Bytes()
}

// Object writes an object.
func (e *Encoder) Object(o ObjectEncoder) {
	e.StartObject()
	o.EncodeFields(e)
	e.EndObject()
}

// ObjectOpt writes an object if it is not nil.
func (e *Encoder) ObjectOpt(o ObjectEncoder) {
	if o == nil || o.IsEmpty() {
		return
	}
	e.StartObject()
	o.EncodeFields(e)
	e.EndObject()
}

// ObjectField writes a field key and an object.
func (e *Encoder) ObjectField(key string, o ObjectEncoder) {
	e.StartObjectField(key)
	o.EncodeFields(e)
	e.EndObject()
}

// ObjectFieldOpt writes a field key and an object if the object is not nil.
func (e *Encoder) ObjectFieldOpt(key string, o ObjectEncoder) {
	if o == nil || o.IsEmpty() {
		return
	}
	e.StartObjectField(key)
	o.EncodeFields(e)
	e.EndObject()
}

// StartObject writes an object start delimiter.
func (e *Encoder) StartObject() {
	e.WriteByte('{')
}

// StartObjectField writes a field key and an object start delimiter.
func (e *Encoder) StartObjectField(key string) {
	e.Key(key)
	e.WriteByte('{')
}

// EndObject writes an object end delimiter.
// The function will panic if the output buffer is still empty.
func (e *Encoder) EndObject() {
	b := e.Buffer.Bytes()
	l := len(b)
	if l > 0 {
		l--
		if b[l] == ',' {
			e.Truncate(l)
		}
	}
	e.WriteString("},")
}

// Array writes an array.
func (e *Encoder) Array(o ArrayEncoder) {
	e.StartArray()
	o.EncodeValues(e)
	e.EndArray()
}

// ArrayOpt writes an array if it is not nil and not empty.
func (e *Encoder) ArrayOpt(o ArrayEncoder) {
	if o == nil || o.IsEmpty() {
		return
	}
	e.StartArray()
	o.EncodeValues(e)
	e.EndArray()
}

// ArrayField writes a field key and an array.
func (e *Encoder) ArrayField(key string, o ArrayEncoder) {
	e.StartArrayField(key)
	o.EncodeValues(e)
	e.EndArray()
}

// ArrayFieldOpt writes a field key and an array if the array is not nil and not empty.
func (e *Encoder) ArrayFieldOpt(key string, o ArrayEncoder) {
	if o == nil || o.IsEmpty() {
		return
	}
	e.StartArrayField(key)
	o.EncodeValues(e)
	e.EndArray()
}

// StartArray writes an array start delimiter.
func (e *Encoder) StartArray() {
	e.WriteByte('[')
}

// StartArrayField writes a field key and an array start delimiter.
func (e *Encoder) StartArrayField(key string) {
	e.Key(key)
	e.WriteByte('[')
}

// EndArray writes an array end delimiter.
// The function will panic if the output buffer is still empty.
func (e *Encoder) EndArray() {
	b := e.Buffer.Bytes()
	l := len(b)
	if l > 0 {
		l--
		if b[l] == ',' {
			e.Truncate(l)
		}
	}
	e.WriteString("],")
}

// EscapedStringField writes a field key and an escaped string.
func (e *Encoder) EscapedStringField(key string, s string) {
	e.Key(key)
	e.EscapedString(s)
}

// EscapedStringFieldOpt writes a field key and an escaped string if the string is not empty.
func (e *Encoder) EscapedStringFieldOpt(key string, s string) {
	if s == "" {
		return
	}
	e.Key(key)
	e.EscapedString(s)
}

// StringField writes a field key and a non-escaped string.
func (e *Encoder) StringField(key string, s string) {
	e.Key(key)
	e.String(s)
}

// StringFieldOpt writes a field key and a non-escaped string if the string is not empty.
func (e *Encoder) StringFieldOpt(key string, s string) {
	if s == "" {
		return
	}
	e.Key(key)
	e.String(s)
}

// StringMapField writes the string map.
func (e *Encoder) StringMapField(key string, m map[string]string) {
	e.Key(key)
	e.StringMap(m)
}

// StringMapFieldOpt writes a field key and a non-escaped string if the string is not empty.
func (e *Encoder) StringMapFieldOpt(key string, m map[string]string) {
	if len(m) == 0 {
		return
	}
	e.Key(key)
	e.StringMap(m)
}

// EscapedBytesStringField writes a field key and an escaped string from the given slice of bytes.
func (e *Encoder) EscapedBytesStringField(key string, s []byte) {
	e.Key(key)
	e.EscapedBytesString(s)
}

// EscapedBytesStringFieldOpt writes a field key and an escaped string if the slice is not empty.
func (e *Encoder) EscapedBytesStringFieldOpt(key string, s []byte) {
	if len(s) == 0 {
		return
	}
	e.Key(key)
	e.EscapedBytesString(s)
}

// BytesStringField writes a field key and a non-escaped string from the given slice of bytes.
func (e *Encoder) BytesStringField(key string, s []byte) {
	e.Key(key)
	e.BytesString(s)
}

// BytesStringFieldOpt writes a field key and a non-escaped string if the slice is not empty.
func (e *Encoder) BytesStringFieldOpt(key string, s []byte) {
	if len(s) == 0 {
		return
	}
	e.Key(key)
	e.BytesString(s)
}

// BoolField write a field key and boolean value.
func (e *Encoder) BoolField(key string, b bool) {
	e.Key(key)
	e.Bool(b)
}

// BoolFieldOpt write a field key and boolean value if the boolean is true.
func (e *Encoder) BoolFieldOpt(key string, b bool) {
	if !b {
		return
	}
	e.Key(key)
	e.Bool(b)
}

// IntBoolField write a field key and boolean value as 0 or 1.
func (e *Encoder) IntBoolField(key string, b bool) {
	e.Key(key)
	e.IntBool(b)
}

// IntBoolFieldOpt write a field key and boolean value as 1 if the boolean is true.
func (e *Encoder) IntBoolFieldOpt(key string, b bool) {
	if !b {
		return
	}
	e.Key(key)
	e.IntBool(b)
}

// IntField writes a field key and integer value.
func (e *Encoder) IntField(key string, i int) {
	e.Int64Field(key, int64(i))
}

// IntFieldOpt writes a field key and integer value if the integer is not zero.
func (e *Encoder) IntFieldOpt(key string, i int) {
	if i == 0 {
		return
	}
	e.Int64Field(key, int64(i))
}

// Int8Field writes a field key and 8-bit integer value.
func (e *Encoder) Int8Field(key string, i int8) {
	e.Int64Field(key, int64(i))
}

// Int8FieldOpt writes a field key and 8-bit integer value if the integer is not zero.
func (e *Encoder) Int8FieldOpt(key string, i int8) {
	if i == 0 {
		return
	}
	e.Int64Field(key, int64(i))
}

// Int16Field writes a field key and 16-bit integer value.
func (e *Encoder) Int16Field(key string, i int16) {
	e.Int64Field(key, int64(i))
}

// Int16FieldOpt writes a field key and 16-bit integer value if the integer is not zero.
func (e *Encoder) Int16FieldOpt(key string, i int16) {
	if i == 0 {
		return
	}
	e.Int64Field(key, int64(i))
}

// Int32Field writes a field key and 32-bit integer value.
func (e *Encoder) Int32Field(key string, i int32) {
	e.Int64Field(key, int64(i))
}

// Int32FieldOpt writes a field key and 32-bit integer value if the integer is not zero.
func (e *Encoder) Int32FieldOpt(key string, i int32) {
	if i == 0 {
		return
	}
	e.Int64Field(key, int64(i))
}

// Int64Field writes a field key and 64-bit integer value.
func (e *Encoder) Int64Field(key string, i int64) {
	e.Key(key)
	e.Int64(i)
}

// Int64FieldOpt writes a field key and 64-bit integer value if the integer is not zero.
func (e *Encoder) Int64FieldOpt(key string, i int64) {
	if i == 0 {
		return
	}
	e.Key(key)
	e.Int64(i)
}

// UintField writes a field key and unsigned integer value.
func (e *Encoder) UintField(key string, i uint) {
	e.Uint64Field(key, uint64(i))
}

// UintFieldOpt writes a field key and unsigned integer value if the integer is not zero.
func (e *Encoder) UintFieldOpt(key string, i uint) {
	if i == 0 {
		return
	}
	e.Uint64Field(key, uint64(i))
}

// Uint8Field writes a field key and 8-bit unsigned integer value.
func (e *Encoder) Uint8Field(key string, i uint8) {
	e.Uint64Field(key, uint64(i))
}

// Uint8FieldOpt writes a field key and 8-bit unsigned integer value if the integer is not zero.
func (e *Encoder) Uint8FieldOpt(key string, i uint8) {
	if i == 0 {
		return
	}
	e.Uint64Field(key, uint64(i))
}

// Uint16Field writes a field key and 16-bit unsigned integer value.
func (e *Encoder) Uint16Field(key string, i uint16) {
	e.Uint64Field(key, uint64(i))
}

// Uint16FieldOpt writes a field key and 16-bit unsigned integer value if the integer is not zero.
func (e *Encoder) Uint16FieldOpt(key string, i uint16) {
	if i == 0 {
		return
	}
	e.Uint64Field(key, uint64(i))
}

// Uint32Field writes a field key and 32-bit unsigned integer value.
func (e *Encoder) Uint32Field(key string, i uint32) {
	e.Uint64Field(key, uint64(i))
}

// Uint32FieldOpt writes a field key and 32-bit unsigned integer value if the integer is not zero.
func (e *Encoder) Uint32FieldOpt(key string, i uint32) {
	if i == 0 {
		return
	}
	e.Uint64Field(key, uint64(i))
}

// Uint64Field writes a field key and 64-bit unsigned integer value.
func (e *Encoder) Uint64Field(key string, i uint64) {
	e.Key(key)
	e.Uint64(i)
}

// Uint64FieldOpt writes a field key and 64-bit unsigned integer value if the integer is not zero.
func (e *Encoder) Uint64FieldOpt(key string, i uint64) {
	if i == 0 {
		return
	}
	e.Key(key)
	e.Uint64(i)
}

// FloatField writes a field key and a 64-bit floating point value.
func (e *Encoder) FloatField(key string, i float64, fmt byte, prec int) {
	e.Key(key)
	e.Float(i, fmt, prec)
}

// FloatFieldOpt writes a field key and a 64-bit floating point value if the value is not zero.
func (e *Encoder) FloatFieldOpt(key string, i float64, fmt byte, prec int) {
	if i == 0.0 {
		return
	}
	e.Key(key)
	e.Float(i, fmt, prec)
}

// TimestampField writes a field key and a time-stamp value.
func (e *Encoder) TimestampField(key string, t time.Time) {
	e.Key(key)
	e.Timestamp(t)
}

// TimestampFieldOpt writes a field key and a time-stamp value if the value is valid.
func (e *Encoder) TimestampFieldOpt(key string, t time.Time) {
	if t.After(threshold) {
		return
	}
	e.Key(key)
	e.Timestamp(t)
}

// StringMap writes the map as an object with the values escaped.
func (e *Encoder) StringMap(m map[string]string) {
	e.StartObject()
	for k, v := range m {
		e.EscapedStringField(k, v)
	}
	e.EndObject()
}

var threshold = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)

// Key writes a field key.
func (e *Encoder) Key(key string) {
	e.WriteByte('"')
	e.WriteString(key)
	e.WriteString(`":`)
}

// EscapedString writes an escaped string.
func (e *Encoder) EscapedString(s string) {
	e.WriteByte('"')
	for ix := range s {
		b := s[ix]
		switch b {
		case '\\', '"':
			e.WriteByte('\\')
		case '\n':
			e.WriteByte('\\')
			b = 'n'
		case '\r':
			e.WriteByte('\\')
			b = 'r'
		case '\t':
			e.WriteByte('\\')
			b = 't'
		}
		e.WriteByte(b)
	}
	e.WriteString(`",`)
}

// String writes a non-escaped string.
func (e *Encoder) String(s string) {
	e.WriteByte('"')
	e.WriteString(s)
	e.WriteString(`",`)
}

// EscapedBytesString writes an escaped string from the given slice.
func (e *Encoder) EscapedBytesString(s []byte) {
	e.WriteByte('"')
	for ix := range s {
		b := s[ix]
		switch b {
		case '\\', '"':
			e.WriteByte('\\')
		case '\n':
			e.WriteByte('\\')
			b = 'n'
		case '\r':
			e.WriteByte('\\')
			b = 'r'
		case '\t':
			e.WriteByte('\\')
			b = 't'
		}
		e.WriteByte(b)
	}
	e.WriteString(`",`)
}

// BytesString writes a non-escaped string from a slice of bytes.
func (e *Encoder) BytesString(s []byte) {
	e.WriteByte('"')
	e.Write(s)
	e.WriteString(`",`)
}

// Bool writes a boolean value.
func (e *Encoder) Bool(b bool) {
	if b {
		e.WriteString(trueString)
	} else {
		e.WriteString(falseString)
	}
	e.WriteByte(',')
}

// IntBool writes a boolean value as 0 or 1.
func (e *Encoder) IntBool(b bool) {
	if b {
		e.WriteString("1,")
	} else {
		e.WriteString("0,")
	}
}

// Int64 writes a 64-bit integer value.
func (e *Encoder) Int64(i int64) {
	if i >= 0 {
		e.Uint64(uint64(i))
	} else {
		b := make([]byte, 0, 32)
		e.Write(strconv.AppendInt(b, i, 10))
		e.WriteByte(',')
	}
}

// Uint64 writes a 64-bit unsigned integer value.
func (e *Encoder) Uint64(i uint64) {
	switch {
	case i < 10:
		e.WriteByte('0' + byte(i))
	case i < 100:
		j := i / 10
		e.WriteByte('0' + byte(j))
		e.WriteByte('0' + byte(i-j*10))
	default:
		b := make([]byte, 0, 32)
		e.Write(strconv.AppendUint(b, i, 10))
	}
	e.WriteByte(',')
}

// Float writes a 64-bit floating point value.
func (e *Encoder) Float(i float64, fmt byte, prec int) {
	b := make([]byte, 0, 48)
	e.Write(strconv.AppendFloat(b, i, fmt, prec, 64))
	e.WriteByte(',')
}

// Timestamp writes a time-stamp value.
func (e *Encoder) Timestamp(t time.Time) {
	e.String(t.Format(timestampMilli))
}

// Raw adds the given bytes as-is.
func (e *Encoder) Raw(data []byte) {
	e.Write(data)
}

// encoderPool is the memory pool for encoders.
var encoderPool = sync.Pool{New: func() interface{} { return &Encoder{} }}
