package pool

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

// Producer is the interface for instantiating new memory pool objects.
// Use the Acquire functions to instantiate a new object.
// Call the object Release() function to return it to the pool when it goes out of scope.
// Producer is safe for concurrent use by multiple go-routines.
type Producer interface {
	Acquire() Objecter                        // instantiate a new object from the memory pool.
	AcquireFromJSON([]byte) (Objecter, error) // acquire a new object and initialize it from the JSON data.
	// hidden
	release(Objecter) // return an object to the memory pool.
}

// NewProducer instantiates a new memory pool object producer.
// It will instantiate an object and panic if the returned object does not adhere to the Objecter interface.
func NewProducer(new func() (Objecter, func())) Producer {
	p := &producer{}

	// create the first object to verify the new() function is valid.
	o1, f1 := new()
	if o1 == nil {
		panic("failed to instantiate new object")
	}

	// initialize the memory pool.
	p.pool = sync.Pool{New: func() any {
		o2, f2 := new()
		o2.setInternal(p, o2, f2) // also sets ref-count to 1.
		return o2
	}}

	// set the first object internals, and release it.
	// we do it here after the memory pool has been intialized!
	o1.setInternal(p, o1, f1) // also sets ref-count to 1.
	o1.Release()

	return p
}

// Acquire implements the Producer.Acquire interface.
func (p *producer) Acquire() Objecter {
	o := p.pool.Get().(Objecter)
	o.initRefs()
	return o
}

// AcquireFromJSON implements the Producer.AcquireFromJSON interface.
// This function will panic if Decode() interface is not reimplemented on the object.
func (p *producer) AcquireFromJSON(data []byte) (Objecter, error) {
	o := p.Acquire()
	dec := zjson.AcquireDecoder(data)
	defer dec.Release()
	if err := o.Decode(dec); err != nil {
		o.Release()
		return nil, err
	}
	return o, nil
}

// release implements the Producer.release interface.
func (p *producer) release(o Objecter) {
	p.pool.Put(o)
}

// producer is the underlying structure for a memory pool object Producer.
type producer struct {
	pool sync.Pool // hidden!
}
