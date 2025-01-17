package pool

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

// Objecter is the interface that all memory pool objects must adhere to.
type Objecter interface {
	zjson.ObjectEncoder
	zjson.ObjectDecoder
	Acquire() Objecter           // instantiates a new memory pool object from the same object producer (implemented here).
	Release()                    // returns the object to its memory pool (implemented here).
	Clone() Objecter             // makes a deep copy of the object (implemented here).
	Encode(*zjson.Encoder)       // encode the object (implemented here; reimplement if not encoding to a JSON object).
	Decode(*zjson.Decoder) error // decode the object (implemented here; reimplement if not decoding from a JSON object).
	// hidden (for internal use only!)
	setInternal(Producer, Objecter, func()) // initializes internal reference fields (implemented here).
	initRefs()                              // initializes reference count (implemented here).
}

// Object serves as the base structure for memory pool objects.
// It is reference counted, so multiple pointers to the same object can exist with a different TTL.
// Object has no use on its own.
// Derived objects should reimplement any base functions where the default behaviour is not sufficient.
type Object struct {
	refs     int      // reference counter.
	producer Producer // the producer that produced this.
	self     Objecter // interface to itself. E.g. the derived object, not "Object".
	reset    func()   // reset function (optional).
}

// Acquire instantiates a new memory pool object from the same object producer.
// The function will panic if the object or its internal producer field are undefined.
func (o *Object) Acquire() Objecter {
	return o.producer.Acquire()
}

// Release returns the memory pool object to its memory pool if the reference count reaches 0.
// The function does nothing if the object or one of its internal reference fields are undefined.
// Care must be taken to call Release() only when the object is no longer in use.
func (o *Object) Release() {
	if o != nil && o.self != nil && o.producer != nil {
		// decrease ref-count (if possible).
		if o.refs > 0 {
			o.refs--

			// return to pool when ref-count reaches 0.
			if o.refs == 0 {
				if o.reset != nil {
					o.reset() // use internal reset function to clear the object to its initial state.
				}
				o.producer.release(o.self)
			}
		}
	}
}

// IsEmpty implements the zjson.ObjectEncoder.IsEmpty interface.
// By default, it returns false. Reimplement the function in derived objects if needed.
func (o *Object) IsEmpty() bool { return false }

// Clone makes a reference copy of the object.
// This allows multiple pointers to the object, each with a different TTL.
// Once all references have called Release() the object will be returned to the memory pool.
func (o *Object) Clone() Objecter {
	// as simple as increasing the ref-count and returning itself.
	o.refs++
	return o.self
}

// Encode encodes the object to json using the given encoder.
// By default, it encodes as a JSON object.
// You *must* reimplement the EncodeFields function for Encode to work!
func (o *Object) Encode(enc *zjson.Encoder) {
	enc.Object(o.self)
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
// By default, it panics. Reimplement the function in derived objects if needed.
func (o *Object) EncodeFields(*zjson.Encoder) { panic("EncodeFields() not implemented") }

// Decode decodes the object from json using the given decoder.
// By default, it decodes as a JSON object.
// You *must* reimplement the DecodeField function for Decode to work!
func (o *Object) Decode(dec *zjson.Decoder) error {
	if dec.Object(o.self) {
		return nil
	}
	return dec.Error()
}

// DecodeField implements the zjson.ObjectDecoder.DecodeField interface.
// By default, it panics. Reimplement the function in derived objects if needed.
func (o *Object) DecodeField(*zjson.Decoder, []byte) error { panic("DecodeField() not implemented") }

// setInternal initializes the internal reference fields.
// It is hidden from external use and called only once when the object is instantiated by its Producer.
func (o *Object) setInternal(producer Producer, self Objecter, reset func()) {
	o.producer = producer
	o.self = self
	o.reset = reset
	o.initRefs()
}

// initRefs initializes the reference counter for a "new" object.
func (o *Object) initRefs() {
	o.refs = 1
}
