package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// StickyAction is an action that can change the stickiness of symbols in the spin result.
type StickyAction struct {
	SpinAction
	reset         bool
	bestSymbol    bool
	symbolsChoice bool
	symbols       utils.Indexes
	reels         []int
}

// NewBestSymbolStickyAction instantiates a new best symbol sticky action.
func NewBestSymbolStickyAction() *StickyAction {
	a := newStickyAction()
	a.bestSymbol = true
	return a.finalizer()
}

// NewSymbolsChooseStickyAction instantiates a new symbol choice sticky action.
func NewSymbolsChooseStickyAction() *StickyAction {
	a := newStickyAction()
	a.symbolsChoice = true
	a.result = ChooseSticky
	return a.finalizer()
}

// NewStickySymbolAction instantiates a new sticky symbol action.
// Reels is an optional array to limit the action to the indicated reels.
// Note that reels are 1-based!
func NewStickySymbolAction(symbol utils.Index, reels ...int) *StickyAction {
	a := newStickyAction()
	a.symbol = symbol
	a.reels = reels
	return a.finalizer()
}

// NewStickySymbolsAction instantiates a new sticky symbols action.
func NewStickySymbolsAction(symbols ...utils.Index) *StickyAction {
	a := newStickyAction()
	a.symbols = symbols
	return a.finalizer()
}

// NewResetStickiesAction instantiates a new reset stickies action.
func NewResetStickiesAction() *StickyAction {
	a := newStickyAction()
	a.reset = true
	return a.finalizer()
}

// Triggered implements the SpinActioner.Triggered interface.
func (a *StickyAction) Triggered(spin *Spin) SpinActioner {
	return a.TriggeredWithState(spin, nil)
}

// TriggeredWithState implements the SpinActioner.TriggeredWithState interface.
func (a *StickyAction) TriggeredWithState(spin *Spin, state *SymbolsState) SpinActioner {
	switch {
	case a.bestSymbol:
		// find best symbol and make it sticky.
		a.FindBestSymbol(spin, state)
		return a

	case a.symbolsChoice:
		// test if there are multiple symbols to be selected as sticky.
		if a.HaveStickyChoices(spin) {
			return a
		}

	case a.symbol != 0:
		// make specific symbol sticky.
		var got bool
		if len(a.reels) == 0 {
			for offset := range spin.indexes {
				if spin.indexes[offset] == a.symbol {
					spin.sticky[offset] = true
					spin.payouts[offset] = 0 // make sure it doesn't get removed with cascading reels!
					got = true
				}
			}
		} else {
			for ix := range a.reels {
				reel := a.reels[ix] - 1
				min := reel * spin.rowCount
				max := min + int(spin.mask[reel])
				for offset := min; offset < max; offset++ {
					if spin.indexes[offset] == a.symbol {
						spin.sticky[offset] = true
						spin.payouts[offset] = 0 // make sure it doesn't get removed with cascading reels!
						got = true
					}
				}
			}
		}
		if got {
			return a
		}

	case len(a.symbols) > 0:
		// make multiple symbols sticky.
		var got bool
		for offset := range spin.indexes {
			if a.symbols.Contains(spin.indexes[offset]) {
				spin.sticky[offset] = true
				got = true
			}
		}
		if got {
			return a
		}

	case a.reset:
		spin.ResetSticky()
		return a
	}

	return nil
}

// FindBestSymbol returns the symbol id of the symbol that is best suited to be the sticky symbol.
// The algorithm returns the highest symbol id with the highest matching count.
// The function also updates the sticky flags in the spin result.
func (a *StickyAction) FindBestSymbol(spin *Spin, state *SymbolsState) utils.Index {
	spin.ResetSticky()

	var bestCount, highestCount uint8
	var highestSymbol utils.Index

	unflagged := func(id utils.Index) bool {
		return state == nil || !state.IsFlagged(id)
	}

	for id := spin.GetSymbols().maxID; id > 0; id-- {
		currCount := spin.CountSymbol(id)
		// mark the best unflagged symbol.
		if currCount > bestCount && unflagged(id) {
			bestCount = currCount
			spin.stickySymbol = id
		}
		// remember the best symbol.
		if currCount > highestCount {
			highestCount = currCount
			highestSymbol = id
		}
	}

	// if we did find a symbol with a higher count but already flagged, we must use that one.
	if highestCount > bestCount {
		spin.stickySymbol = highestSymbol
		bestCount = highestCount
	}

	if a.CanTrigger(spin) {
		for ix, id := range spin.indexes {
			if id == spin.stickySymbol {
				spin.sticky[ix] = true
			}
		}
	}

	return spin.stickySymbol
}

// HaveStickyChoices returns whether the grid has multiple choices for a sticky symbol.
func (a *StickyAction) HaveStickyChoices(spin *Spin) bool {
	var count int
	for id := spin.GetSymbols().maxID; id > 0; id-- {
		if spin.CountSymbol(id) > 1 {
			count++
		}
	}
	return count > 1
}

func newStickyAction() *StickyAction {
	a := &StickyAction{}
	a.init(TestStickiness, Sticky, reflect.TypeOf(a).String())
	return a
}

func (a *StickyAction) finalizer() *StickyAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.bestSymbol:
		b.WriteString(",bestSymbol=true")

	case a.symbolsChoice:
		b.WriteString(",symbolsChoice=true")

	case a.symbol != 0:
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))

	case len(a.symbols) > 0:
		b.WriteString(",symbols=")
		j, _ := json.Marshal(a.symbols)
		b.Write(j)

	case a.reset:
		b.WriteString(",reset=true")
	}

	a.config = b.String()
	return a
}
