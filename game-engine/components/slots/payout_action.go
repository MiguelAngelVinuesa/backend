package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// PayoutAction is the action that calculates the regular payouts for a spin. E.g. paylines or cluster payouts.
// It can also be used to remove previously awarded payouts based on the total payout factor!
// Note that PayoutAction does not have a Triggered() function as they are always executed for a spin.
type PayoutAction struct {
	SpinAction
	// details for awarding payouts.
	paylines      bool            // normal paylines.
	allPaylines   bool            // "all" paylines (mega-ways).
	highestPayout bool            // indicates highest payout if payline starts with wilds.
	cluster       *ClusterPayouts // cluster payouts.

	// details for removing payouts.
	removePayouts     bool               // true if payout removal is active.
	removePayoutBands bool               // true if payout band removal is active.
	remMinFactor      float64            // minimum payout factor required (max 2 decimals).
	remMaxFactor      float64            // maximum payout factor required (max 2 decimals).
	remPayoutChance   float64            // chance of this action to be triggered (max 2 decimals).
	remPayoutDir      PayDirection       // LTR, RTL or both.
	remPayoutMech     int                // removal mechanism (1=dedupe first reels, 2=remove bonus, 3=dedupe 2nd/3rd reel).
	remPayoutWilds    bool               // always remove wild symbols.
	remDupes          bool               // duplicate symbols allowed on a reel.
	remBands          []RemovePayoutBand // payout bands & removal chances.
}

// NewPaylinesAction instantiates a new paylines action.
func NewPaylinesAction() *PayoutAction {
	a := newPayoutAction()
	a.paylines = true
	return a.finalizer()
}

// NewAllPaylinesAction instantiates a new all paylines action.
func NewAllPaylinesAction(highestPayout bool) *PayoutAction {
	a := newPayoutAction()
	a.allPaylines = true
	a.highestPayout = highestPayout
	return a.finalizer()
}

// NewClusterPayoutsAction instantiates a new cluster pays action.
func NewClusterPayoutsAction(reels, rows uint8, opts ...ClusterPayoutsOption) *PayoutAction {
	a := newPayoutAction()
	a.cluster = NewClusterPayouts(reels, rows, opts...)
	return a.finalizer()
}

// NewRemovePayoutsAction instantiates a payouts.
// Payouts are removed chance% of the time *and* if the total payout factor is between min and max.
// direction, mechanism, removeWilds and allowDupes determine the operation(s) used to remove the payouts.
// Note: only mechanism 1 is currently supported!
func NewRemovePayoutsAction(min, max, chance float64, direction PayDirection, mechanism int, removeWilds, allowDupes bool) *PayoutAction {
	if mechanism != 1 && mechanism != 3 {
		panic("ReleasePayouts - invalid mechanism")
	}

	a := newPayoutAction()
	a.removePayouts = true
	a.remMinFactor = min
	a.remMaxFactor = max
	a.remPayoutChance = chance
	a.remPayoutDir = direction
	a.remPayoutMech = mechanism
	a.remPayoutWilds = removeWilds
	a.remDupes = allowDupes
	return a.finalizer()
}

// RemovePayoutBand represents a band of payouts to be removed
// MinPayout is considered inclusive but MaxPayout is considered exclusive (min <= x < max).
// MinPayout can be zero, however a grid with a zero payout is ignored as there are no payouts to be removed.
// RemoveChance indicates the % how often the removal takes place.
type RemovePayoutBand struct {
	MinPayout    float64 `json:"min"`
	MaxPayout    float64 `json:"max"`
	RemoveChance float64 `json:"chance"`
}

// NewRemovePayoutBandsAction instantiates a payout bands remove action.
// Note: this action should be included after the regular payouts have been awarded!
// The total payout factor determines which band to test, and if it is within a specified band, any payouts will be removed chance% of the time.
// mechanism, direction, removeWilds and allowDupes determine the operation(s) used to remove the payouts.
// Only mechanism 1 is currently supported!
func NewRemovePayoutBandsAction(mechanism int, direction PayDirection, removeWilds, allowDupes bool, bands []RemovePayoutBand) *PayoutAction {
	if mechanism != 1 && mechanism != 3 {
		panic("RemovePayoutBands - invalid mechanism")
	}

	a := newPayoutAction()
	a.removePayoutBands = true
	a.remPayoutDir = direction
	a.remPayoutMech = mechanism
	a.remPayoutWilds = removeWilds
	a.remDupes = allowDupes
	a.remBands = bands
	return a.finalizer()
}

// NewRemoveBonusPayoutsAction instantiates a bonus payouts remove action.
// Bonus payouts are removed chance% of the time *and* if the bonus payout factor is between min and max.
// allowDupes determine if duplicate symbols on a reel are allowed.
// Note: only mechanism 2 is currently supported!
func NewRemoveBonusPayoutsAction(min, max, chance float64, mechanism int, allowDupes bool) *PayoutAction {
	if mechanism != 2 {
		panic("RemoveBonusPayouts - invalid mechanism")
	}

	a := newPayoutAction()
	a.stage = RegularPayouts
	a.removePayouts = true
	a.remMinFactor = min
	a.remMaxFactor = max
	a.remPayoutChance = chance
	a.remPayoutMech = mechanism
	a.remDupes = allowDupes
	return a.finalizer()
}

// Payout implements the SpinAction.Payout interface.
// It will update the result with any awarded payouts (or remove them).
func (a *PayoutAction) Payout(spin *Spin, res *results.Result) SpinActioner {
	switch {
	case a.paylines:
		if a.testPaylines(spin, res) {
			return a
		}
	case a.allPaylines:
		if a.testAllPaylines(spin, res, a.highestPayout) {
			return a
		}
	case a.cluster != nil:
		if a.cluster.Find(spin, res) > 0 {
			return a
		}
	case a.removePayouts && !spin.debugInitial:
		if a.testRemovePayouts(spin, res) {
			return a
		}
	case a.removePayoutBands && !spin.debugInitial:
		if a.testRemovePayoutBands(spin, res) {
			return a
		}
	}
	return nil
}

// HavePaylines returns true if the payouts are based on specific paylines.
func (a *PayoutAction) HavePaylines() bool {
	return a.paylines
}

// HaveAllPaylines returns true if the all paylines feature is active.
func (a *PayoutAction) HaveAllPaylines() bool {
	return a.allPaylines
}

// HaveClusterPayouts returns true if cluster payouts are active.
func (a *PayoutAction) HaveClusterPayouts() bool {
	return a.cluster != nil
}

// RemoveClusterPays removes cluster payouts.
func (a *PayoutAction) RemoveClusterPays(spin *Spin) {
	if a.cluster != nil {
		a.cluster.RemovePayouts(spin)
	}
}

func newPayoutAction() *PayoutAction {
	a := &PayoutAction{}
	a.init(RegularPayouts, Payout, reflect.TypeOf(a).String())

	a.allPaylines = false
	a.removePayouts = false
	a.removePayoutBands = false
	a.remMinFactor = 0.0
	a.remMaxFactor = 0.0
	a.remPayoutChance = 0.0
	a.remPayoutDir = PayLTR
	a.remPayoutMech = 0
	a.remPayoutWilds = false
	a.remDupes = true

	return a
}

func (a *PayoutAction) finalizer() *PayoutAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.paylines:
		b.WriteString(",paylines=true")

	case a.allPaylines:
		b.WriteString(",allPaylines=true")

	case a.cluster != nil:
		b.WriteString(",clusterPayouts=true")

	case a.removePayouts:
		b.WriteString(",removePayouts=true")
		b.WriteString(",minFactor=")
		b.WriteString(strconv.FormatFloat(a.remMinFactor, 'f', 2, 64))
		b.WriteString(",maxFactor=")
		b.WriteString(strconv.FormatFloat(a.remMaxFactor, 'f', 2, 64))
		b.WriteString(",chance=")
		b.WriteString(strconv.FormatFloat(a.remPayoutChance, 'f', 2, 64))
		b.WriteString(",dir=")
		b.WriteString(strconv.Itoa(int(a.remPayoutDir)))
		b.WriteString(",mech=")
		b.WriteString(strconv.Itoa(int(a.remPayoutMech)))
		b.WriteString(",wilds=")
		if a.remPayoutWilds {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}

	case a.removePayoutBands:
		b.WriteString(",removePayoutBands=true")
		b.WriteString(",bands=")
		j, _ := json.Marshal(a.remBands)
		b.Write(j)
		b.WriteString(",dir=")
		b.WriteString(strconv.Itoa(int(a.remPayoutDir)))
		b.WriteString(",mech=")
		b.WriteString(strconv.Itoa(int(a.remPayoutMech)))
		b.WriteString(",wilds=")
		if a.remPayoutWilds {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	}

	a.config = b.String()
	return a
}

// PayoutActions is a convenience type for a slice of PayoutAction
type PayoutActions []*PayoutAction

// PurgePayoutActions returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgePayoutActions(input PayoutActions, capacity int) PayoutActions {
	if cap(input) < capacity {
		return make(PayoutActions, 0, capacity)
	}
	return input[:0]
}
