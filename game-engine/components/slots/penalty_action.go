package slots

import (
	"bytes"
	"fmt"
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// PenaltyAction is the action that can impose a penalty on a spin. E.g. for "reverse win" slot games.
// A penalty decreases the total win in some way, as defined by the penalty action.
// Note that PenaltyAction does not have a Triggered() function as they are always executed for a spin.
type PenaltyAction struct {
	SpinAction
	// details for imposing a penalty.
	reduce bool    // reduction penaly.
	divide bool    // division penalty.
	count  uint8   // symbol count to trigger the penalty.
	factor float64 // factor of the penalty.
}

// NewReductionAction instantiates a new reduction penalty action.
func NewReductionAction(symbol utils.Index, count uint8, factor float64) *PenaltyAction {
	a := newPenaltyAction()
	a.reduce = true
	a.symbol = symbol
	a.count = count
	a.factor = factor
	return a.finalizer()
}

// NewDivisionAction instantiates a new division penalty action.
func NewDivisionAction(symbol utils.Index, count uint8, factor float64) *PenaltyAction {
	a := newPenaltyAction()
	a.divide = true
	a.symbol = symbol
	a.count = count
	a.factor = factor
	return a.finalizer()
}

// Penalty implements the SpinActioner.Penalty interface.
func (a *PenaltyAction) Penalty(spin *Spin, res *results.Result) SpinActioner {
	count := spin.CountSymbol(a.symbol)
	if count < a.count {
		return nil
	}

	switch {
	case a.reduce:
		res.AddPenalties(NewSlotReduction(a.symbol, count, a.factor, spin))
	case a.divide:
		res.AddPenalties(NewSlotDivision(a.symbol, count, a.factor, spin))
	}

	return a
}

func newPenaltyAction() *PenaltyAction {
	a := &PenaltyAction{}
	a.init(RegularPenalties, Penalty, reflect.TypeOf(a).String())
	return a
}

func (a *PenaltyAction) finalizer() *PenaltyAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.reduce:
		b.WriteString(",reduce=true")
	case a.divide:
		b.WriteString(",divide=true")
	}

	b.WriteString(",factor=")
	b.WriteString(fmt.Sprintf("%.1f", a.factor))

	a.config = b.String()
	return a
}

// PenaltyActions is a convenience type for a slice of PenaltyAction.
type PenaltyActions []*PenaltyAction

// PurgePenaltyActions returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgePenaltyActions(input PenaltyActions, capacity int) PenaltyActions {
	if cap(input) < capacity {
		return make(PenaltyActions, 0, capacity)
	}
	return input[:0]
}
