package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

type StateAction struct {
	SpinAction
	flagSymbols bool
	flagCount   uint8
	minSpins    int
	maxSpins    int
}

// NewFlagSymbolsAction instantiates a new symbols state action.
func NewFlagSymbolsAction(flagCount uint8, altSymbols bool, opts ...StateActionOption) *StateAction {
	a := &StateAction{}
	a.init(TestState, Processed, reflect.TypeOf(a).String())

	a.altSymbols = altSymbols
	a.flagSymbols = true
	a.flagCount = flagCount

	for ix := range opts {
		opts[ix](a)
	}

	return a.finalizer()
}

// StateUpdate updates the game state as configured.
func (a *StateAction) StateUpdate(spin *Spin, symbols *SymbolsState) {
	if a.flagSymbols {
		a.doFlagSymbols(spin, symbols)
	}

	// TODO: other state changes

}

// StateTriggered indicates if the given symbol state triggers this action.
func (a *StateAction) StateTriggered(symbols *SymbolsState) SpinActioner {
	if a.flagSymbols && !symbols.triggered && symbols.AllFlagged() {
		symbols.triggered = true
		return a
	}

	// TODO: other state triggers

	return nil
}

// NrOfSpins returns the number of free spins from the spin action.
func (a *StateAction) NrOfSpins(prng interfaces.Generator) uint8 {
	if a.minSpins > 0 && a.maxSpins > a.minSpins {
		return uint8(a.minSpins + prng.IntN(a.maxSpins-a.minSpins+1))
	}
	return a.nrOfSpins
}

// doFlagSymbols flags symbols if the count >= flagCount.
func (a *StateAction) doFlagSymbols(spin *Spin, symbols *SymbolsState) {
	max := spin.symbols.maxID
	for id := utils.Index(1); id <= max; id++ {
		if spin.symbols.GetSymbol(id) != nil && !symbols.IsFlagged(id) {
			var count uint8
			for iy := range spin.indexes {
				if spin.indexes[iy] == id {
					count++
				}
			}
			if count >= a.flagCount {
				symbols.SetFlagged(id, true)
			}
		}
	}
}

// StateActionOption is the function signature for state action options.
type StateActionOption func(*StateAction)

// WithFreeSpins marks the state action to reward free spins.
func WithFreeSpins(nrOrSpins, minSpins, maxSpins uint8) StateActionOption {
	return func(a *StateAction) {
		a.result = FreeSpins
		a.nrOfSpins = nrOrSpins
		a.minSpins = int(minSpins)
		a.maxSpins = int(maxSpins)
	}
}

func (a *StateAction) finalizer() *StateAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	b.WriteString(",flagSymbols=true")
	b.WriteString(",flagCount=")
	b.WriteString(strconv.Itoa(int(a.flagCount)))
	b.WriteString(",altSymbols=")
	if a.altSymbols {
		b.WriteString("true")
	} else {
		b.WriteString("false")
	}

	if a.result == FreeSpins {
		if a.nrOfSpins > 0 {
			b.WriteString(",nrOfSpins=")
			b.WriteString(strconv.Itoa(int(a.nrOfSpins)))
		}
		if a.minSpins > 0 || a.maxSpins > 0 {
			b.WriteString(",minSpins=")
			b.WriteString(strconv.Itoa(int(a.minSpins)))
			b.WriteString(",maxSpins=")
			b.WriteString(strconv.Itoa(int(a.maxSpins)))
		}
	}

	a.config = b.String()
	return a
}

// StateActions is a convenience type for a slice of state actions.
type StateActions []*StateAction

// PurgeStateActions returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeStateActions(input StateActions, capacity int) StateActions {
	if cap(input) < capacity {
		return make(StateActions, 0, capacity)
	}
	return input[:0]
}
