package results

import (
	"log"
	"math"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
)

type ResultDataKind uint8

const (
	SpinData ResultDataKind = iota + 1
	InstantBonusData
	BonusSelectorData
	BonusWheelData
	BonusCallColorData
	BonusCallSuitData
)

// AcquireResult instantiates a new result from the memory pool.
// Note that Result takes ownership of the Data; DO NOT call Release() on it!
// Payouts are cloned, so it's safe to call Release() on those.
func AcquireResult(data interfaces.Objecter2, kind ResultDataKind, payouts ...Payout) *Result {
	r := resultPool.Acquire().(*Result)
	r.DataKind = kind
	r.Data = data

	l := len(payouts)
	size := object.NormalizeSize(l, 8)
	r.Payouts = PurgePayouts(r.Payouts, size)
	r.Animations = PurgeAnimations(r.Animations, 8)

	r.AddPayouts(payouts...)
	return r
}

// AddState adds the game state to the result.
// If a state was already added, it will be returned to the memory pool.
// Result takes ownership of the state, so make sure to provide a clone or deep copy of the state!
func (r *Result) AddState(state interfaces.Objecter2) *Result {
	if r.State != nil {
		r.State.Release()
	}
	r.State = state
	return r
}

// AddPayouts adds one or more payouts to the result.
// Note that Result takes ownership of the added payouts. DO NOT call Release() on them!
func (r *Result) AddPayouts(list ...Payout) *Result {
	for ix := range list {
		p := list[ix]
		r.Total += p.Total()
		r.Payouts = append(r.Payouts, p)
	}
	if r.Total > 5000 {
		_ = r
	}
	return r
}

// AddPenalties adds one or more penalties to the result.
// Note that Result takes ownership of the added penalties. DO NOT call Release() on them!
func (r *Result) AddPenalties(list ...Penalty) *Result {
	for ix := range list {
		p := list[ix]
		switch p.Kind() {
		case SlotReduction:
			r.Total -= p.Reduce()
		}
		r.Penalties = append(r.Penalties, p)
	}
	return r
}

// AddAnimations adds one or more animation events to the result.
// Note that Result takes ownership of the added events; DO NOT call Release() on them!
func (r *Result) AddAnimations(list ...Animator) *Result {
	r.Animations = append(r.Animations, list...)
	return r
}

// SetMessage implements the Localizer.SetMessage interface.
func (r *Result) SetMessage(msg string) {
	r.Message = msg
}

// ReleasePayouts removes all payouts.
func (r *Result) ReleasePayouts() {
	r.Payouts = ReleasePayouts(r.Payouts)
	r.Penalties = ReleasePenalties(r.Penalties)
	r.Total = 0.0
}

// AwardFreeGames sets the awarded free games for this result.
func (r *Result) AwardFreeGames(freeGames uint8) {
	r.AwardedFreeGames += uint64(freeGames)
}

// SetFreeGames sets the remaining free games counter for this result.
func (r *Result) SetFreeGames(freeGames uint64) {
	r.FreeGames = freeGames
}

// EncodeFields implements the zjson Encoder interface.
func (r *Result) EncodeFields(enc *zjson.Encoder) {
	enc.Uint64FieldOpt("awardedFreeGames", r.AwardedFreeGames)
	enc.Uint64Field("freeGames", r.FreeGames)
	enc.FloatField("total", math.Round(r.Total*100)/100, 'g', -1)
	enc.StringFieldOpt("message", r.Message)
	enc.Uint8FieldOpt("dataKind", uint8(r.DataKind))
	enc.ObjectFieldOpt("data", r.Data)
	enc.ObjectFieldOpt("state", r.State)

	enc.StartArrayField("payouts")
	for ix := range r.Payouts {
		enc.Object(r.Payouts[ix])
	}
	for ix := range r.Penalties {
		enc.Object(r.Penalties[ix].AsPayout())
	}
	enc.EndArray()

	if len(r.Animations) > 0 {
		enc.StartArrayField("animations")
		for ix := range r.Animations {
			enc.Object(r.Animations[ix])
		}
		enc.EndArray()
	}
}

// Encode2 implements the zjson Encoder interface with alternate encoding for the data.
func (r *Result) Encode2(enc *zjson.Encoder) {
	enc.Uint64FieldOpt("awardedFreeGames", r.AwardedFreeGames)
	enc.Uint64Field("freeGames", r.FreeGames)
	enc.FloatField("total", math.Round(r.Total*100)/100, 'g', -1)
	enc.StringFieldOpt("message", r.Message)
	enc.Uint8FieldOpt("dataKind", uint8(r.DataKind))

	if r.Data != nil {
		enc.StartObjectField("data")
		r.Data.Encode2(enc)
		enc.EndObject()
	}

	enc.ObjectFieldOpt("state", r.State)

	enc.StartArrayField("payouts")
	for ix := range r.Payouts {
		enc.Object(r.Payouts[ix])
	}
	for ix := range r.Penalties {
		enc.Object(r.Penalties[ix].AsPayout())
	}
	enc.EndArray()

	if len(r.Animations) > 0 {
		enc.StartArrayField("animations")
		for ix := range r.Animations {
			enc.Object(r.Animations[ix])
		}
		enc.EndArray()
	}
}

// Result represents the details for a single game result.
// A game round may contain multiple results, such as when there are free spins in a slots game.
type Result struct {
	DataKind         ResultDataKind       `json:"dataKind,omitempty"`         // the specific game-data kind.
	AwardedFreeGames uint64               `json:"awardedFreeGames,omitempty"` // newly awarded free games for this result.
	FreeGames        uint64               `json:"freeGames"`                  // counter for remaining free games after this result.
	Total            float64              `json:"total"`                      // total payout factor for this result.
	Data             interfaces.Objecter2 `json:"data"`                       // the game-specific data.
	State            interfaces.Objecter2 `json:"state,omitempty"`            // the game-specific state.
	Payouts          Payouts              `json:"payouts"`                    // individual payouts for this result.
	Penalties        Penalties            `json:"penalties,omitempty"`        // individual penalties for this result.
	Animations       Animations           `json:"animations,omitempty"`       // animation events for this result.
	Message          string               `json:"message,omitempty"`          // localized message for the total payout.
	pool.Object
}

var (
	resultPool = pool.NewProducer(func() (pool.Objecter, func()) {
		r := &Result{}
		return r, r.reset
	})
)

// reset clears the result.
func (r *Result) reset() {
	if r != nil {
		r.DataKind = 0
		r.AwardedFreeGames = 0
		r.FreeGames = 0

		if r.Data != nil {
			r.Data.Release()
			r.Data = nil
		}
		if r.State != nil {
			r.State.Release()
			r.State = nil
		}

		r.ReleasePayouts() // this clears Penalties & Total as well.

		r.Animations = ReleaseAnimations(r.Animations)
	}
}

// Results is a convenience type for a slice of results.
type Results []*Result

// ReleaseResults releases all results in the slice and returns the empty slice.
func ReleaseResults(list Results) Results {
	if list == nil {
		return nil
	}
	for ix := range list {
		if r := list[ix]; r != nil {
			r.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeResults returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
// It automatically takes care of releasing the results.
func PurgeResults(list Results, capacity int) Results {
	list = ReleaseResults(list)
	if cap(list) < capacity || cap(list) > 256 {
		FreeResults(list)
		return NewResults(capacity)
	}
	return list
}

func CloneResults(results Results) Results {
	l := len(results)
	out := NewResults(l)[:l]
	for ix := range results {
		out[ix] = results[ix].Clone().(*Result)
	}
	return out
}

// FixPenalties recalculates the total payout for all results with division penalties.
// Note: this function must be called only once after all results have been gathered!
func FixPenalties(results Results) float64 {
	var total float64
	for ix := range results {
		r := results[ix]
		total += r.Total
		for iy := range r.Penalties {
			p := r.Penalties[iy]
			if p.Kind() == SlotDivision {
				log.Printf("\n R TOTAL START %v  TOTAL %v \n", r.Total, total)
				if total <= 0 {
					p.SetFactor(0)
					return total
				}
				cut := total / p.Divide()

				integerPart := int(math.Abs(cut))
				count := 0.0
				for integerPart > 0 {
					integerPart /= 10
					count++
				}

				fract := 0.0
				if count <= 3 {
					fract = cut - math.Floor(cut)
					cut -= fract
				} else {
					fract = cut/10*(count-3) - math.Floor(cut/10*(count-3))
					cut -= fract * (10 * (count - 3))
				}

				if cut != 0 {
					total = cut
				}

				r.Total -= total
				log.Printf("\n R TOTAL END %v  TOTAL %v  CUT %v \n", r.Total, total, cut)

				p.SetFactor(-total)
			}
		}
	}
	return total
}

// GrandTotal sums the payout factors in the slice of results.
func GrandTotal(results Results) float64 {
	var total float64
	for ix := range results {
		total += results[ix].Total
	}
	if total <= 0 {
		return 0
	}
	return total
}

// GrandTotal2 sums the payout factors in the slice of results.
// It returns max(maxPayout, sum).
func GrandTotal2(results Results, maxPayout float64) float64 {
	if sum := GrandTotal(results); sum < maxPayout {
		return sum
	}
	return maxPayout
}

// NewResults acquires a slice of results from the memory pool.
func NewResults(capacity int) Results {
	if capacity <= 16 {
		return resultsPool.Get().(Results)
	}
	capacity = ((capacity + 15) / 16) * 16
	return make(Results, 0, capacity)
}

// FreeResults clears the slice of results and returns it to the memory pool.
func FreeResults(results Results) {
	if len(results) >= 16 && len(results) <= 256 {
		resultsPool.Put(ReleaseResults(results))
	}
}

var resultsPool = sync.Pool{New: func() any { return make(Results, 0, 16) }}
