package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// RoundFlagAction represents actions that operate on the flags for a round of spins (first spin + free spins).
type RoundFlagAction struct {
	SpinAction
	flag            int
	weightedFlag    bool
	weights         utils.WeightedGenerator
	symbolUsedFlag  bool
	symbolCountFlag bool
	fullBonusFlag   bool
	increaseFlag    bool
	decreaseFlag    bool
	shapeDetect     bool
	shapeGrid       GridOffsets
	resetFlags      bool
	setFlag         bool
	flagValue       int
	flags           []int
}

// NewRoundFlagWeightedAction instantiates a new spin round flag action.
// When triggered the indicated spin round flag is updated with a random value based on the weights using the spin PRNG.
// Note that the indexes in the given weighting must fit in an uint8!
func NewRoundFlagWeightedAction(flag int, weights utils.WeightedGenerator) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = PreBonus // only once per round, before the free spins!!!
	a.flag = flag
	a.weightedFlag = true
	a.weights = weights
	return a.finalizer()
}

// NewRoundFlagSymbolUsedAction instantiates a new spin round flag action.
// When the symbol occurs during a spin the indicated spin round flag is set to 1.
func NewRoundFlagSymbolUsedAction(flag int, symbol utils.Index) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = AwardBonuses // must happen after grid revisions!!!
	a.flag = flag
	a.symbolUsedFlag = true
	a.symbol = symbol
	return a.finalizer()
}

// NewRoundFlagSymbolCountAction instantiates a new spin round flag action.
// It sets the round flag to the count of the symbol in the grid.
func NewRoundFlagSymbolCountAction(flag int, symbol utils.Index) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = AwardBonuses // must happen after grid revisions!!!
	a.flag = flag
	a.symbolCountFlag = true
	a.symbol = symbol
	return a.finalizer()
}

// NewRoundFlagFullBonusAction instantiates a new spin round flag action.
// When the bonus symbol occurs across all reels of a spin the indicated spin round flag is set to 1.
func NewRoundFlagFullBonusAction(flag int) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = AwardBonuses // must happen after grid revisions!!!
	a.flag = flag
	a.fullBonusFlag = true
	return a.finalizer()
}

// NewRoundFlagShapeDetect instantiates a new spin round flag action.
// When the given grid shape is detected to contain teh same symbol, the indicated spin round flag is set to 1.
// The symbol is optional; if not zero the shape must be present with this symbol.
func NewRoundFlagShapeDetect(flag int, symbol utils.Index, grid GridOffsets) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = AwardBonuses // must happen after grid revisions!!!
	a.flag = flag
	a.shapeDetect = true
	a.symbol = symbol
	a.shapeGrid = grid
	return a.finalizer()
}

// NewRoundFlagIncreaseAction instantiates a new spin round flag action.
// Every time it triggers it increases the round flag by 1.
func NewRoundFlagIncreaseAction(flag int) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = ReviseGrid // right at the start of each spin!!!
	a.flag = flag
	a.increaseFlag = true
	return a.finalizer()
}

// NewRoundFlagDecreaseAction instantiates a new spin round flag action.
// Every time it triggers it decreases the round flag by 1 until it is zero.
func NewRoundFlagDecreaseAction(flag int) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = ReviseGrid // right at the start of each spin!!!
	a.flag = flag
	a.decreaseFlag = true
	return a.finalizer()
}

// NewRoundFlagsReset instantiates a new spin round flag action.
// Every time it triggers it resets the given flags to zero.
func NewRoundFlagsReset(flags ...int) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = ReviseGrid // right at the start of each spin!!!
	a.flags = flags
	a.resetFlags = true
	return a.finalizer()
}

// NewRoundFlagSet instantiates a new spin round flag action.
// Every time it triggers it sets the given flags to the given value.
func NewRoundFlagSet(flag, value int) *RoundFlagAction {
	a := newRoundFlagAction()
	a.stage = ReviseGrid // right at the start of each spin!!!
	a.flag = flag
	a.setFlag = true
	a.flagValue = value
	return a.finalizer()
}

// Triggered implements the SpinAction.Triggerer interface.
func (a *RoundFlagAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.weightedFlag:
		spin.roundFlags[a.flag] = int(a.weights.RandomIndex(spin.prng))
		return a

	case a.symbolUsedFlag:
		if spin.roundFlags[a.flag] != 1 {
			for ix := range spin.indexes {
				if spin.indexes[ix] == a.symbol {
					spin.roundFlags[a.flag] = 1
					return a
				}
			}
		}

	case a.symbolCountFlag:
		if spin.roundFlags[a.flag] == 0 {
			spin.roundFlags[a.flag] = int(spin.CountSymbol(a.symbol))
			return a
		}

	case a.fullBonusFlag:
		if spin.roundFlags[a.flag] != 1 {
			if int(spin.CountBonusSymbol()) >= spin.reelCount {
				spin.roundFlags[a.flag] = 1
				return a
			}
		}

	case a.shapeDetect:
		if spin.roundFlags[a.flag] != 1 {
			rows := spin.rowCount
			detect := a.symbol
			for ix := range a.shapeGrid {
				p := a.shapeGrid[ix]
				offset := p[0]*rows + p[1]
				if detect == utils.NullIndex {
					detect = spin.indexes[offset]
				} else if spin.indexes[offset] != detect {
					return nil
				}
			}
			spin.roundFlags[a.flag] = 1
			return a
		}

	case a.increaseFlag:
		spin.roundFlags[a.flag]++
		return a

	case a.decreaseFlag:
		if spin.roundFlags[a.flag] > 0 {
			spin.roundFlags[a.flag]--
			return a
		}

	case a.resetFlags:
		for _, f := range a.flags {
			spin.roundFlags[f] = 0
		}
		return a

	case a.setFlag:
		spin.roundFlags[a.flag] = a.flagValue
		return a
	}

	return nil
}

func newRoundFlagAction() *RoundFlagAction {
	a := &RoundFlagAction{}
	a.init(PreBonus, Processed, reflect.TypeOf(a).String())
	return a
}

func (a *RoundFlagAction) finalizer() *RoundFlagAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.weightedFlag:
		b.WriteString(",weightedFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))

	case a.symbolUsedFlag:
		b.WriteString(",symbolUsedFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))

	case a.symbolCountFlag:
		b.WriteString(",symbolCountFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))

	case a.fullBonusFlag:
		b.WriteString(",fullBonusFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))

	case a.increaseFlag:
		b.WriteString(",increaseFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))

	case a.decreaseFlag:
		b.WriteString(",decreaseFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))

	case a.resetFlags:
		b.WriteString(",resetFlags=true")
		b.WriteString(",flags=")
		j, _ := json.Marshal(a.flags)
		b.Write(j)

	case a.setFlag:
		b.WriteString(",setFlag=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))
		b.WriteString(",value=")
		b.WriteString(strconv.Itoa(a.flagValue))
	}

	a.config = b.String()
	return a
}
