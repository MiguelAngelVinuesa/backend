package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// ScatterAction is a trigger activated by scatter symbols.
type ScatterAction struct {
	SpinAction
	payout        bool
	freeSpins     bool
	bonusWheel    bool
	scatterCount  uint8
	scatterPayout float64
	bonusScatter  bool
	bonusLines    uint8
	allScatters   bool
	allMinimum    uint8
	multiSymbols  utils.Indexes
	wheelFlag     int
	wheelWeights  utils.WeightedGenerator
}

// NewScatterPayoutAction instantiates a new scatter action that awards payouts.
// Make sure to use WithAlternate if the same symbol can result in multiple payouts but only one should be executed.
func NewScatterPayoutAction(symbol utils.Index, count uint8, payout float64) *ScatterAction {
	a := newScatterAction()
	a.payout = true
	a.symbol = symbol
	a.scatterCount = count
	a.scatterPayout = payout
	return a.finalizer()
}

// NewScatterFreeSpinsAction instantiates a new scatter action that awards free spins.
// Make sure to use WithAlternate if the same symbol can result in multiple awards but only one should be executed.
func NewScatterFreeSpinsAction(nrOfSpins uint8, altSymbols bool, symbol utils.Index, scatterCount uint8, bonusSymbol bool) *ScatterAction {
	a := newScatterAction()
	a.stage = AwardBonuses
	a.result = FreeSpins
	a.freeSpins = true
	a.nrOfSpins = nrOfSpins
	a.altSymbols = altSymbols
	a.symbol = symbol
	a.bonusSymbol = bonusSymbol
	a.scatterCount = scatterCount
	return a.finalizer()
}

// NewScatterBonusWheelAction instantiates a new scatter action that awards a bonus wheel game.
// Make sure to use WithAlternate if the same symbol can result in multiple awards but only one should be executed.
func NewScatterBonusWheelAction(symbol utils.Index, scatterCount uint8, flag int, weights utils.WeightedGenerator) *ScatterAction {
	a := newScatterAction()
	a.stage = AwardBonuses
	a.result = BonusGame
	a.bonusWheel = true
	a.symbol = symbol
	a.scatterCount = scatterCount
	a.wheelFlag = flag
	a.wheelWeights = weights
	return a.finalizer()
}

// NewAllScatterAction instantiates a new scatter action that indicates all symbols act as scatter.
// The AllScatter action signals that all symbols act as a scatter (such as in the ChaChaBomb slot machine game).
func NewAllScatterAction(minimum uint8) *ScatterAction {
	a := newScatterAction()
	a.allScatters = true
	a.allMinimum = minimum
	return a.finalizer()
}

// NewBonusScatterAction instantiates a new scatter action which uses the bonus symbol as a scatter.
func NewBonusScatterAction(lines uint8) *ScatterAction {
	a := newScatterAction()
	a.stage = AwardBonuses // normally happens after payouts and free game awards!
	a.bonusScatter = true
	a.bonusLines = lines
	return a
}

// WithAlternate can be used to add an alternative action for cases where they need to be mutually exclusive.
// E.g. a scatter symbol that triggers a bonus game when there are three, but it triggers a free spin when there are two.
// Alternative actions can be nested. E.g. 3 symbols triggers A, 2 symbols triggers B, 1 symbol triggers C.
// Circular references are not supported and will likely result in a panic due to stack overflow.
func (a *ScatterAction) WithAlternate(alt *ScatterAction) *ScatterAction {
	a.alternate = alt
	return a.finalizer()
}

// WithMultiSymbols sets the multi symbol feature of the scatter action.
// The function will panic if used for anything other than the FreeSpins result.
func (a *ScatterAction) WithMultiSymbols(symbols ...utils.Index) *ScatterAction {
	if a.result != FreeSpins {
		panic(consts.MsgInvalidActionKind)
	}
	a.multiSymbols = symbols
	return a.finalizer()
}

// WithPlayerChoice can be used to indicate that a player choice must be made before the game continues.
func (a *ScatterAction) WithPlayerChoice() *ScatterAction {
	a.playerChoice = true
	return a.finalizer()
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *ScatterAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.payout, a.bonusWheel:
		if spin.newScatters >= a.scatterCount {
			if spin.CountSymbol(a.symbol) >= a.scatterCount {
				return a
			}
		}
		if a.alternate != nil {
			return a.alternate.Triggered(spin)
		}

	case a.freeSpins:
		if spin.newScatters >= a.scatterCount {
			n := spin.CountSymbol(a.symbol)
			for ix := range a.multiSymbols {
				n += spin.CountSymbol(a.multiSymbols[ix])
			}
			if n >= a.scatterCount {
				return a
			}
		}
		if a.alternate != nil {
			return a.alternate.Triggered(spin)
		}

	case a.allScatters:
		// sort all symbol indexes so we can compare symbol equality sequentially.
		max := uint8(len(spin.indexes))
		ids := make(utils.Indexes, 100)[:max]
		copy(ids, spin.indexes)
		ids = utils.SortIndexes(ids)

		// skip dummy symbols (symbolID == 0).
		var count, ix uint8
		var last utils.Index
		for ; ix < max && ids[ix] == 0; ix++ {
		}

		// count same symbols until end of slice or found three of a kind.
		for ix <= max-count && count < a.allMinimum {
			last = ids[ix]
			count = 1
			for ix++; ix < max && count < a.allMinimum; ix++ {
				if ids[ix] == last {
					count++
				} else {
					break
				}
			}
		}

		// if we have at least one symbol with a 3-count, we're there!
		if count >= a.allMinimum {
			return a
		}

	case a.bonusScatter:
		count := spin.CountBonusSymbol()
		if count < 2 {
			return nil
		}
		if symbol := spin.GetSymbol(spin.bonusSymbol); symbol != nil && symbol.Payable(count) {
			return a
		}

	}

	return nil
}

// Payout adds a payout to the results based on the number of scatter symbols.
func (a *ScatterAction) Payout(spin *Spin, res *results.Result) SpinActioner {
	var payouts int

	switch {
	case a.payout:
		return a.testPayout(spin, a.countSymbol(spin), res)

	case a.allScatters:
		// sort all symbol indexes so we can compare symbol equality sequentially.
		max := uint8(len(spin.indexes))
		ids := make(utils.Indexes, 100)[:max]
		copy(ids, spin.indexes)
		ids = utils.SortIndexes(ids)

		// skip dummy symbols (symbolID == 0).
		var count, ix uint8
		var last utils.Index
		for ; ids[ix] == 0 && ix < max; ix++ {
		}

		test := func() {
			if last > 0 && count >= a.allMinimum {
				if symbol := spin.GetSymbol(last); symbol != nil {
					if payout := symbol.Payout(count); payout > 0 {
						if symbol.id == spin.superSymbol {
							res.AddPayouts(SuperSymbolPayout(payout, spin.getMultiplier(0), last, count, spin))
							// resetData stickies on a super-shape payout!
							spin.ResetSticky()
						} else if a.bombScatter(spin, res, last) {
							res.AddPayouts(BombScatterPayout(payout, spin.getMultiplier(0), last, count, spin))
						} else {
							res.AddPayouts(ScatterSymbolPayout(payout, spin.getMultiplier(0), last, count, spin))
						}
						payouts++
					}
				}
				last = 0
			}
		}

		// count same symbols until end of slice or found required count.
		for ix <= max-count {
			last = ids[ix]
			count = 1
			for ix++; ix < max; ix++ {
				if ids[ix] == last {
					count++
				} else {
					test()
					break
				}
			}
		}
		test()

	case a.bonusScatter:
		count := spin.CountBonusSymbol()
		if count < 2 {
			return nil
		}
		symbol := spin.GetSymbol(spin.bonusSymbol)
		if symbol == nil {
			panic(consts.MsgSymbolNotFound)
		}
		if payout := symbol.Payout(count); payout > 0 {
			res.AddPayouts(BonusSymbolPayout(payout*float64(a.bonusLines), spin.getMultiplier(0), spin.bonusSymbol, count))
			spin.ExpandBonusSymbol()
			payouts++
		}

	}

	if payouts > 0 {
		return a
	}
	return nil
}

func (a *ScatterAction) countSymbol(spin *Spin) uint8 {
	count := spin.CountSymbol(a.symbol)
	for ix := range a.multiSymbols {
		count += spin.CountSymbol(a.multiSymbols[ix])
	}
	return count
}

func (a *ScatterAction) testPayout(spin *Spin, count uint8, res *results.Result) SpinActioner {
	if count >= a.scatterCount {
		res.AddPayouts(ScatterSymbolPayout(a.scatterPayout, spin.getMultiplier(0), a.symbol, count, spin))
		return a
	}

	if a.alternate != nil {
		s := a.alternate.(*ScatterAction)
		return s.testPayout(spin, count, res)
	}

	return nil
}

func (a *ScatterAction) bombScatter(spin *Spin, res *results.Result, symbol utils.Index) bool {
	data, ok := res.Data.(*SpinResult)
	if !ok {
		return false
	}

	if symbol != data.stickySymbol || len(data.afterExpand) == 0 {
		return false
	}

	bomb := spin.GetSymbols().GetBombSymbol()
	if bomb == nil {
		return false
	}

	for _, id := range data.initial {
		if id == bomb.id {
			return true
		}
	}
	return false
}

// CanPayout return true if a payout can be generated.
func (a *ScatterAction) CanPayout() bool {
	return a.payout
}

// HaveAllScatterPayouts return true if all scatter payouts feature is active.
func (a *ScatterAction) HaveAllScatterPayouts() bool {
	return a.allScatters
}

// BonusGame returns a configured bonus game.
func (a *ScatterAction) BonusGame(spin *Spin) interfaces.Objecter2 {
	switch {
	case a.bonusWheel:
		w := wheel.AcquireBonusWheel(spin.prng, a.wheelWeights)
		value, result := w.Run(nil)
		w.Release()

		if a.wheelFlag >= 0 {
			spin.roundFlags[a.wheelFlag] = value
		}

		return result

	default:
		return nil
	}
}

// FeatureTransition returns the applicable bonus feature transition kind.
func (a *ScatterAction) FeatureTransition() FeatureTransitionKind {
	switch {
	case a.bonusWheel:
		return BonusWheelTransition
	default:
		return 0
	}
}

func newScatterAction() *ScatterAction {
	a := &ScatterAction{}
	a.init(ExtraPayouts, Payout, reflect.TypeOf(a).String())
	return a
}

func (a *ScatterAction) finalizer() *ScatterAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.freeSpins:
		b.WriteString(",freeSpins=true")
		b.WriteString(",nrOfSpins=")
		b.WriteString(strconv.Itoa(int(a.nrOfSpins)))
		b.WriteString(",altSymbols=")
		if a.altSymbols {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",bonusSymbol=")
		if a.bonusSymbol {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.scatterCount)))

	case a.bonusWheel:
		b.WriteString(",bonusWheel=true")
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.scatterCount)))
		b.WriteString(",wheel=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.wheelFlag))
		b.WriteString(",weights=")
		b.WriteString(a.wheelWeights.String())

	case a.allScatters:
		b.WriteString(",allScatters=true")
		b.WriteString(",minimum=")
		b.WriteString(strconv.Itoa(int(a.allMinimum)))

	case a.bonusScatter:
		b.WriteString(",bonusScatter=true")

	default:
		b.WriteString(",payout=true")
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.scatterCount)))
	}

	a.config = b.String()
	return a
}
