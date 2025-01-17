package results

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// PrngLog can hold the event & PRNG log for a game result.
type PrngLog struct {
	scriptID  int
	eventLog  Events
	rngInLog  []int
	rngOutLog []int
}

// Initialize initializes the event & PRNG log.
func (l *PrngLog) Initialize() {
	l.scriptID = 0
	l.eventLog = make(Events, 0, 32)
	l.rngInLog = make([]int, 0, 32)
	l.rngOutLog = make([]int, 0, 32)
}

// Reset resets the event & PRNG log.
func (l *PrngLog) Reset() {
	l.scriptID = 0
	l.eventLog = ReleaseEvents(l.eventLog)
	l.rngInLog = utils.PurgeInts(l.rngInLog, 32)
	l.rngOutLog = utils.PurgeInts(l.rngOutLog, 32)
}

// Equals is used internally for unit-tests!
func (l *PrngLog) Equals(other *PrngLog) bool {
	if l.scriptID != other.scriptID ||
		len(l.eventLog) != len(other.eventLog) ||
		!reflect.DeepEqual(l.rngInLog, other.rngInLog) ||
		!reflect.DeepEqual(l.rngOutLog, other.rngOutLog) {
		return false
	}

	for ix := range l.eventLog {
		if !l.eventLog[ix].Equals(other.eventLog[ix]) {
			return false
		}
	}

	return true
}

// SetScriptID sets the scripted round scenario id.
func (r *PrngLog) SetScriptID(id int) {
	r.scriptID = id
}

// ScriptID returns the scripted round scenario id.
func (r *PrngLog) ScriptID() int {
	return r.scriptID
}

// Log returns the event & PRNG log.
func (r *PrngLog) Log() (Events, []int, []int) {
	return r.eventLog, r.rngInLog, r.rngOutLog
}

// LogEvent adds an event to the event log.
func (r *PrngLog) LogEvent(e *Event) {
	r.eventLog = append(r.eventLog, e)
}

// CloneEvent copies an event to the event log.
func (r *PrngLog) CloneEvent(e *Event) {
	r.eventLog = append(r.eventLog, e.Clone().(*Event))
}

// SetLog overwrites the PRNG log.
func (r *PrngLog) SetLog(l1, l2 []int) {
	r.rngInLog = l1
	r.rngOutLog = l2
}

// EncodeEventLog implements the Objecter2.Encode2 interface.
func (r *PrngLog) EncodeEventLog(enc *zjson.Encoder) {
	enc.IntFieldOpt("scriptID", r.scriptID)

	if len(r.eventLog) > 0 {
		enc.StartArrayField("events")
		for ix := range r.eventLog {
			enc.Object(r.eventLog[ix])
		}
		enc.EndArray()
	}

	if len(r.rngInLog) > 0 {
		enc.StartArrayField("rngIn")
		for ix := range r.rngInLog {
			enc.Int64(int64(r.rngInLog[ix]))
		}
		enc.EndArray()
	}

	if len(r.rngOutLog) > 0 {
		enc.StartArrayField("rngOut")
		for ix := range r.rngOutLog {
			enc.Int64(int64(r.rngOutLog[ix]))
		}
		enc.EndArray()
	}
}

// DecodeScriptID can be used to decode the scriptID.
func (r *PrngLog) DecodeScriptID(dec *zjson.Decoder) bool {
	var ok bool
	r.scriptID, ok = dec.Int()
	return ok
}

// DecodeEventLog can be used to decode the event log.
func (r *PrngLog) DecodeEventLog(dec *zjson.Decoder) error {
	r.eventLog = PurgeEvents(r.eventLog, 32)
	e := EventProducer.Acquire().(*Event)
	if ok := dec.Object(e); ok {
		r.eventLog = append(r.eventLog, e)
		return nil
	}
	return dec.Error()
}

// DecodeRngIn can be used to decode the PRNG input log.
func (r *PrngLog) DecodeRngIn(dec *zjson.Decoder) error {
	r.rngInLog = utils.PurgeInts(r.rngInLog, 32)
	if i, ok := dec.Int(); ok {
		r.rngInLog = append(r.rngInLog, i)
		return nil
	}
	return dec.Error()
}

// DecodeRngOut can be used to decode the PRNG output log.
func (r *PrngLog) DecodeRngOut(dec *zjson.Decoder) error {
	r.rngOutLog = utils.PurgeInts(r.rngOutLog, 32)
	if i, ok := dec.Int(); ok {
		r.rngOutLog = append(r.rngOutLog, i)
		return nil
	}
	return dec.Error()
}
