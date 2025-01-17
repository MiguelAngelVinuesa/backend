package results

import (
	"fmt"
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireEvent instantiates a new event from the memory pool.
func AcquireEvent(id int, triggered bool, rngIn, rngOut []int) *Event {
	e := EventProducer.Acquire().(*Event)
	e.id = id
	e.triggered = triggered

	if l := len(rngIn); l > 0 {
		max := object.NormalizeSize(l, 16)
		e.rngIn = utils.PurgeInts(e.rngIn, max)[:l]
		copy(e.rngIn, rngIn)
	}

	if l := len(rngOut); l > 0 {
		max := object.NormalizeSize(l, 16)
		e.rngOut = utils.PurgeInts(e.rngOut, max)[:l]
		copy(e.rngOut, rngOut)
	}

	return e
}

// ID returns the event id.
func (e *Event) ID() int {
	return e.id
}

// Triggered returns the event triggered flag.
func (e *Event) Triggered() bool {
	return e.triggered
}

// RngIn returns the event random number inputs.
func (e *Event) RngIn() []int {
	return e.rngIn
}

// RngOut returns the event random number outputs.
func (e *Event) RngOut() []int {
	return e.rngOut
}

// IsEmpty implements the zjson.Encoder interface.
func (e *Event) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson.Encoder interface.
func (e *Event) EncodeFields(enc *zjson.Encoder) {
	enc.IntField("id", e.id)
	enc.IntBoolField("triggered", e.triggered)

	if len(e.rngIn) > 0 {
		enc.StartArrayField("rngIn")
		for ix := range e.rngIn {
			enc.Int64(int64(e.rngIn[ix]))
		}
		enc.EndArray()
	}

	if len(e.rngOut) > 0 {
		enc.StartArrayField("rngOut")
		for ix := range e.rngOut {
			enc.Int64(int64(e.rngOut[ix]))
		}
		enc.EndArray()
	}
}

// DecodeField implements the zjson.Decoder interface.
func (e *Event) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok, t bool
	var i int

	if string(key) == "id" {
		if i, ok = dec.Int(); ok {
			e.id = i
		}
	} else if string(key) == "triggered" {
		if t, ok = dec.IntBool(); ok {
			e.triggered = t
		}
	} else if string(key) == "rngIn" {
		e.rngIn = utils.PurgeInts(e.rngIn, 16)
		ok = dec.Array(e.decodeRngIn)
	} else if string(key) == "rngOut" {
		e.rngOut = utils.PurgeInts(e.rngOut, 16)
		ok = dec.Array(e.decodeRngOut)
	} else {
		return fmt.Errorf("Event.DecodeField: invalid field encountered [%s]", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (e *Event) decodeRngIn(dec *zjson.Decoder) error {
	if i, ok := dec.Int(); ok {
		e.rngIn = append(e.rngIn, i)
		return nil
	}
	return dec.Error()
}

func (e *Event) decodeRngOut(dec *zjson.Decoder) error {
	if i, ok := dec.Int(); ok {
		e.rngOut = append(e.rngOut, i)
		return nil
	}
	return dec.Error()
}

// Event is used to log events from actions.
// Keep fields ordered by ascending SizeOf().
type Event struct {
	triggered bool
	id        int
	rngIn     []int
	rngOut    []int
	pool.Object
}

// EventProducer is the memory pool for events.
// Make sure to initialize all slices appropriately!
var EventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &Event{
		rngIn:  make([]int, 0, 16),
		rngOut: make([]int, 0, 16),
	}
	return e, e.reset
})

// reset clears the event.
func (e *Event) reset() {
	if e != nil {
		e.triggered = false
		e.id = 0
		e.rngIn = e.rngIn[:0]
		e.rngOut = e.rngOut[:0]
	}
}

// Equals is used internally for unit-tests!
func (e *Event) Equals(other *Event) bool {
	return e.id == other.id &&
		e.triggered == other.triggered &&
		reflect.DeepEqual(e.rngIn, other.rngIn) &&
		reflect.DeepEqual(e.rngOut, other.rngOut)
}

// Events is a convenience type for an array of events.
type Events []*Event

// ReleaseEvents releases the events and returns an empty slice.
func ReleaseEvents(list Events) Events {
	if list == nil {
		return nil
	}
	for ix := range list {
		if e := list[ix]; e != nil {
			e.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeEvents resets the input slice or returns a new one if its capacity is less than required.
func PurgeEvents(list Events, capacity int) Events {
	list = ReleaseEvents(list)
	if cap(list) < capacity {
		return make(Events, 0, capacity)
	}
	return list[:0]
}
