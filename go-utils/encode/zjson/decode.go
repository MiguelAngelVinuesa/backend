package zjson

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// ObjectDecoder is the interface for decoding objects.
type ObjectDecoder interface {
	DecodeField(dec *Decoder, key []byte) error
}

// ArrayValue is the function signature for processing array values.
type ArrayValue func(dec *Decoder) error

// Decoder implements a zero-allocation json decoder.
// Note that it does memory allocations if you use it to decode into string variables.
// It is faster than most json unmarshalers, more than 6x faster than the standard json unmarshaler.
// However, it is not a drop-in replacement and requires handwritten code to unmarshall objects & arrays.
// Decoder is not safe for use across concurrent go-routines.
type Decoder struct {
	buf   []byte
	err   error
	max   int
	ptr   int
	value []byte
}

// AcquireDecoder instantiates a zero allocation json decoder from the memory pool.
func AcquireDecoder(b []byte) *Decoder {
	d := decoderPool.Get().(*Decoder)
	d.buf = b
	d.max = len(b)
	if d.value == nil {
		d.value = make([]byte, 0, 256)
	}
	return d
}

// Release returns the decoder to the memory pool.
func (d *Decoder) Release() {
	if d != nil {
		d.Reset(nil)
		decoderPool.Put(d)
	}
}

// Reset resets the decoder with the given input.
func (d *Decoder) Reset(b []byte) {
	d.buf = b
	d.err = nil
	d.max = len(b)
	d.ptr = 0
	d.value = d.value[:0]
}

// Error returns the last encountered error or nil if no errors occurred.
func (d *Decoder) Error() error {
	return d.err
}

// Object decodes an object and returns true if successful.
func (d *Decoder) Object(o ObjectDecoder) bool {
	if b, ok := d.nextToken(); !ok || b != '{' {
		d.err = fmt.Errorf("start delimiter '{' missing")
		return false
	}

	first := true
	for {
		b, ok := d.nextToken()
		if ok && b == '}' {
			return true
		}
		if !ok {
			d.err = fmt.Errorf("end delimiter '}' missing")
			return false
		}

		if first {
			d.ptr-- // push back
			first = false
		} else {
			if b != ',' {
				d.err = fmt.Errorf("field separator ',' missing")
				return false
			}
		}

		start, end, escaped, ok2 := d.readString()
		if !ok2 {
			return false
		}
		if b, ok = d.nextToken(); !ok || b != ':' {
			d.err = fmt.Errorf("key delimiter ':' missing")
			return false
		}

		if escaped {
			d.err = o.DecodeField(d, d.Unescaped(d.buf[start:end]))
		} else {
			d.err = o.DecodeField(d, d.buf[start:end])
		}
		if d.err != nil {
			return false
		}
	}
}

// Array decodes an array and returns true if successful.
func (d *Decoder) Array(v ArrayValue) bool {
	if b, ok := d.nextToken(); !ok || b != '[' {
		d.err = fmt.Errorf("start delimiter '[' missing")
		return false
	}

	first := true
	for {
		b, ok := d.nextToken()
		if ok && b == ']' {
			return true
		}
		if !ok {
			d.err = fmt.Errorf("end delimiter ']' missing")
			return false
		}

		if first {
			d.ptr-- // push back
			first = false
		} else {
			if b != ',' {
				d.err = fmt.Errorf("field separator ',' missing")
				return false
			}
		}

		if d.err = v(d); d.err != nil {
			return false
		}
	}
}

// String decodes a string value and returns true if successful.
func (d *Decoder) String() ([]byte, bool, bool) {
	if start, end, escaped, ok := d.readString(); ok {
		return d.buf[start:end], escaped, true
	}
	return nil, false, false
}

// StringMap decodes a string map into the given map and returns true if successful.
func (d *Decoder) StringMap(in map[string]string) (map[string]string, bool) {
	out := in
	if out == nil {
		out = make(map[string]string)
	}

	if b, ok := d.nextToken(); !ok || b != '{' {
		d.err = fmt.Errorf("start delimiter '{' missing")
		return out, false
	}

	first := true
	for {
		b, ok := d.nextToken()
		if ok && b == '}' {
			return out, true
		}
		if !ok {
			d.err = fmt.Errorf("end delimiter '}' missing")
			return out, false
		}

		if first {
			d.ptr-- // push back
			first = false
		} else {
			if b != ',' {
				d.err = fmt.Errorf("field separator ',' missing")
				return out, false
			}
		}

		start, end, escaped, ok2 := d.readString()
		if !ok2 {
			return out, false
		}

		if b, ok = d.nextToken(); !ok || b != ':' {
			d.err = fmt.Errorf("key delimiter ':' missing")
			return out, false
		}

		var k string
		if escaped {
			k = string(d.Unescaped(d.buf[start:end]))
		} else {
			k = string(d.buf[start:end])
		}

		start, end, escaped, ok2 = d.readString()
		if !ok2 {
			return out, false
		}

		if escaped {
			out[k] = string(d.Unescaped(d.buf[start:end]))
		} else {
			out[k] = string(d.buf[start:end])
		}
	}
}

// Unescaped returns an unescaped version of the input bytes.
func (d *Decoder) Unescaped(data []byte) []byte {
	if l := len(data); l > cap(d.value) {
		l = FixBufferSize(l)
		d.value = make([]byte, 0, l)
	} else {
		d.value = d.value[:0]
	}

	// TODO: support \uxxxx json escape sequences.

	var escape bool
	for ix := 0; ix < len(data); ix++ {
		b := data[ix]
		if escape {
			switch b {
			case 'b':
				b = '\b'
			case 'f':
				b = '\f'
			case 'n':
				b = '\n'
			case 'r':
				b = '\r'
			case 't':
				b = '\t'
			}
			d.value = append(d.value, b)
			escape = false
		} else {
			switch b {
			case '\\':
				escape = true
			default:
				d.value = append(d.value, b)
			}
		}
	}

	return d.value
}

// Bool decodes a boolean value and returns true if successful.
func (d *Decoder) Bool() (bool, bool) {
	if start, end, ok := d.readValue(); ok && end > start {
		l := end - start
		b := d.buf[start:end]
		switch {
		case l == 4 && b[0] == 't' && b[1] == 'r' && b[2] == 'u' && b[3] == 'e':
			return true, true
		case l == 5 && b[0] == 'f' && b[1] == 'a' && b[2] == 'l' && b[3] == 's' && b[4] == 'e':
			return false, true
		}
	}
	d.err = fmt.Errorf("invalid boolean value")
	return false, false
}

// IntBool decodes a boolean value from an integer and returns true if successful.
func (d *Decoder) IntBool() (bool, bool) {
	if start, end, ok := d.readValue(); ok && end-start == 1 {
		b := d.buf[start:end]
		switch {
		case b[0] == '1':
			return true, true
		case b[0] == '0':
			return false, true
		}
	}
	d.err = fmt.Errorf("invalid boolean value")
	return false, false
}

// Int decodes an integer value and returns true if successful.
func (d *Decoder) Int() (int, bool) {
	i, ok := d.readInt64()
	if ok {
		if i >= math.MinInt && i <= math.MaxInt {
			return int(i), true
		}
		d.err = fmt.Errorf("invalid range for int")
	}
	return 0, false
}

// Int8 decodes an 8-bit integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an int8.
func (d *Decoder) Int8() (int8, bool) {
	if i, ok := d.readInt64(); ok {
		if i >= math.MinInt8 && i <= math.MaxInt8 {
			return int8(i), true
		}
		d.err = fmt.Errorf("invalid range for int8")
	}
	return 0, false
}

// Int16 decodes a 16-bit integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an int16.
func (d *Decoder) Int16() (int16, bool) {
	if i, ok := d.readInt64(); ok {
		if i >= math.MinInt16 && i <= math.MaxInt16 {
			return int16(i), true
		}
		d.err = fmt.Errorf("invalid range for int16")
	}
	return 0, false
}

// Int32 decodes a 32-bit integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an int32.
func (d *Decoder) Int32() (int32, bool) {
	if i, ok := d.readInt64(); ok {
		if i >= math.MinInt32 && i <= math.MaxInt32 {
			return int32(i), true
		}
		d.err = fmt.Errorf("invalid range for int32")
	}
	return 0, false
}

// Int64 decodes a 64-bit integer and returns true if successful.
func (d *Decoder) Int64() (int64, bool) {
	return d.readInt64()
}

// Uint decodes an integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an uint.
func (d *Decoder) Uint() (uint, bool) {
	if i, ok := d.readUint64(); ok {
		if i <= math.MaxUint {
			return uint(i), true
		}
		d.err = fmt.Errorf("invalid range for uint")
	}
	return 0, false
}

// Uint8 decodes an 8-bit integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an uint8.
func (d *Decoder) Uint8() (uint8, bool) {
	if i, ok := d.readUint64(); ok {
		if i <= math.MaxUint8 {
			return uint8(i), true
		}
		d.err = fmt.Errorf("invalid range for uint8")
	}
	return 0, false
}

// Uint16 decodes a 16-bit integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an uint16.
func (d *Decoder) Uint16() (uint16, bool) {
	if i, ok := d.readUint64(); ok {
		if i <= math.MaxUint16 {
			return uint16(i), true
		}
		d.err = fmt.Errorf("invalid range for uint16")
	}
	return 0, false
}

// Uint32 decodes a 32-bit integer and returns true if successful.
// The function returns false and sets an error if the input doesn't fit an uint32.
func (d *Decoder) Uint32() (uint32, bool) {
	if i, ok := d.readInt64(); ok {
		if i <= math.MaxUint32 {
			return uint32(i), true
		}
		d.err = fmt.Errorf("invalid range for uint32")
	}
	return 0, false
}

// Uint64 decodes a 64-bit integer and returns true if successful.
func (d *Decoder) Uint64() (uint64, bool) {
	return d.readUint64()
}

// Float decodes a floating point number and returns true if successful.
func (d *Decoder) Float() (float64, bool) {
	start, end, ok := d.readValue()
	if !ok {
		d.err = fmt.Errorf("float error: no input")
		return 0, false
	}

	var f float64
	b := d.buf[start:end]
	// ParseFloat requires a string, but we have a byte slice, so use unsafe package to prevent escaping to heap!
	f, d.err = strconv.ParseFloat(*(*string)(unsafe.Pointer(&b)), 64)
	if d.err != nil {
		return 0, false
	}
	return f, true
}

// Timestamp decodes a time-stamp value and returns true if successful.
func (d *Decoder) Timestamp() (time.Time, bool) {
	start, end, escaped, ok := d.readString()
	if !ok || escaped {
		d.err = fmt.Errorf("timetamp error: no or bad input")
		return time.Time{}, false
	}

	var t time.Time
	b := d.buf[start:end]
	// Parse requires a string but, we have a byte slice, so use unsafe package to prevent escaping to heap!
	t, d.err = time.Parse(timestampMilli, *(*string)(unsafe.Pointer(&b)))
	if d.err != nil {
		return time.Time{}, false
	}
	return t, true
}

// readInt64 decodes a 64-bit signed integer and returns true if successful.
func (d *Decoder) readInt64() (int64, bool) {
	start, end, ok := d.readValue()
	if !ok {
		d.err = fmt.Errorf("integer error: no input")
		return 0, false
	}

	if end-start == 1 && d.buf[start] >= '0' && d.buf[start] <= '9' {
		return int64(d.buf[start] - '0'), true
	}

	var f int64
	b := d.buf[start:end]
	// ParseInt requires a string but, we have a byte slice, so use unsafe package to prevent escaping to heap!
	f, d.err = strconv.ParseInt(*(*string)(unsafe.Pointer(&b)), 10, 64)
	if d.err != nil {
		return 0, false
	}
	return f, true
}

// readUint64 decodes a 64-bit unsigned signed integer and returns true if successful.
func (d *Decoder) readUint64() (uint64, bool) {
	start, end, ok := d.readValue()
	if !ok {
		d.err = fmt.Errorf("unsigned integer error: no input")
		return 0, false
	}

	if end-start == 1 && d.buf[start] >= '0' && d.buf[start] <= '9' {
		return uint64(d.buf[start] - '0'), true
	}

	var f uint64
	b := d.buf[start:end]
	// ParseUint requires a string, but we have a byte slice, so use unsafe package to prevent escaping to heap!
	f, d.err = strconv.ParseUint(*(*string)(unsafe.Pointer(&b)), 10, 64)
	if d.err != nil {
		return 0, false
	}

	if f < 0.0 {
		d.err = fmt.Errorf("unsigned integer error: negative value")
		return 0, false
	}
	return f, true
}

// readValue decodes a value into the buffer and returns true if successful.
func (d *Decoder) readValue() (int, int, bool) {
	d.nextToken() // skip wsp; just ignore the returned token
	d.ptr--
	start := d.ptr

	for d.ptr < d.max {
		b := d.buf[d.ptr]
		if b == ',' || b == '}' || b == ']' || b == ' ' || b == '\n' || b == '\r' || b == '\t' {
			break
		}
		d.ptr++
	}

	return start, d.ptr, d.ptr > start
}

// readString reads a string from the input and returns true if successful.
// It returns the start/end position of the string and whether it contains escape sequences.
func (d *Decoder) readString() (int, int, bool, bool) {
	if b, ok := d.nextToken(); !ok || b != '"' {
		d.err = fmt.Errorf("string delimiter '\"' missing")
		return 0, 0, false, false
	}

	var escaped bool
	start := d.ptr

	for d.ptr < d.max {
		b := d.buf[d.ptr]
		d.ptr++

		if b == '"' {
			return start, d.ptr - 1, escaped, true
		}

		if b == '\\' {
			// TODO: support \uxxxx json escape sequences.
			d.ptr++
			escaped = true
		}
	}

	d.err = fmt.Errorf("incomplete string: reached EOF")
	return 0, 0, false, false
}

// nextToken skips all white-space and returns the first non white-space character but does not skip past it.
func (d *Decoder) nextToken() (byte, bool) {
	for d.ptr < d.max {
		b := d.buf[d.ptr]
		d.ptr++
		if b != ' ' && b != '\n' && b != '\r' && b != '\t' {
			return b, true
		}
	}
	return 0, false
}

// decoderPool is the memory pool for decoders.
var decoderPool = sync.Pool{New: func() interface{} { return &Decoder{} }}
