package slots

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// WildAction is an action activated by one or more wild symbols, optionally requiring a hero symbol.
// It may optionally indicate expanding wilds and whether they expand before or after payline matching.
type WildAction struct {
	SpinAction
	payout         bool
	wildCount      uint8
	wildPayout     float64
	expand         bool
	transform      bool
	needHero       bool
	expandToReel   bool
	expandToSymbol bool
	lockReels      bool
	expandGrid     GridOffsets
	effectKind     uint8
	jumping        bool
	jumpSymbols    utils.Indexes
	jumpParams     *JumpParams
}

// NewWildPayoutAction instantiates a new wild payout action.
func NewWildPayoutAction(symbol utils.Index, count uint8, payout float64) *WildAction {
	a := newWild(symbol)
	a.payout = true
	a.stage = ExtraPayouts
	a.result = Payout
	a.wildCount = count
	a.wildPayout = payout
	return a.finalizer()
}

// NewWildExpansion instantiates a new wild expansion action.
func NewWildExpansion(nrOfSpins uint8, altSymbols bool, symbol utils.Index, wildCount uint8, needHero, expandWilds, expandFirst bool) *WildAction {
	a := newWild(symbol)
	a.expand = true

	if expandFirst {
		a.stage = ExpandBefore
	} else {
		a.stage = ExpandAfter
	}

	a.nrOfSpins = nrOfSpins
	a.altSymbols = altSymbols
	a.wildCount = wildCount
	a.needHero = needHero
	a.expandToReel = expandWilds

	return a.finalizer()
}

// NewWildTransform instantiates a new wild transformation action.
// The expansion grid can be nil, in which case a 3x3 expansion grid around the center symbol is assigned automatically.
func NewWildTransform(symbol utils.Index, expandBefore bool, expandGrid GridOffsets) *WildAction {
	a := newWild(symbol)
	a.transform = true

	if expandBefore {
		a.stage = ExpandBefore
	} else {
		a.stage = ExpandAfter
	}

	a.result = Processed
	a.wildCount = 1
	a.expandToSymbol = true

	if len(expandGrid) > 0 {
		a.expandGrid = expandGrid
	} else {
		a.expandGrid = DefaultGridOffsets
	}

	return a.finalizer()
}

// NewJumpingWilds instantiates a new jumping wilds action.
// Make sure to activate it after an appropriate clear action if necessary.
// The function panics if maxStep < minStep or both are zero.
// The function also panics if no symbols are given.
func NewJumpingWilds(dir GridDirection, symbols ...utils.Index) *WildAction {
	if len(symbols) == 0 {
		panic(consts.MsgInvalidJumpParameters)
	}

	a := newWild(0)
	a.jumping = true
	a.stage = TestClearance
	a.result = WildsJumped
	a.jumpSymbols = symbols
	a.jumpParams = &JumpParams{
		direction: dir,
		minJump:   1,
		maxJump:   1,
	}
	return a.finalizer()
}

// WithPlayerChoice can be used to indicate that a player choice must be made before the game continues.
func (a *WildAction) WithPlayerChoice() *WildAction {
	a.playerChoice = true
	return a.finalizer()
}

// WithBombEffect can be used to set the special effects kind for a transform to indicate a bomb explosion.
func (a *WildAction) WithBombEffect() *WildAction {
	a.effectKind = BombEffect
	return a.finalizer()
}

// WithLockReels can be used to lock affected reels after a wild expansion.
func (a *WildAction) WithLockReels() *WildAction {
	a.lockReels = true
	return a.finalizer()
}

// WithJumpSize can be used to set the minimum and maximum jump size for jumping wilds.
// If minJump == 0, the symbol may jump in place.
// Non-rectangular grids with maxJump > 1 may give undesirable results, and is not supported.
// The function will panic if it is not a jumping wild action, or the parameters are invalid or conflicting.
func (a *WildAction) WithJumpSize(min, max uint8) *WildAction {
	if min > 4 || max > 4 || max < min || (min == 0 && max == 0) {
		panic(consts.MsgInvalidJumpParameters)
	}
	a.jumpParams.minJump = min
	a.jumpParams.maxJump = max
	return a.finalizer()
}

// JumpOnSymbols can be used to indicate jumping wilds can land on other symbols.
// The function will panic if it is not a jumping wild action.
func (a *WildAction) JumpOnSymbols() *WildAction {
	a.jumpParams.onSymbols = true
	a.jumpParams.refill = true
	return a.finalizer()
}

// JumpOffGrid can be used to indicate jumping wilds can land off the grid.
// The function will panic if it is not a jumping wild action.
func (a *WildAction) JumpOffGrid(chance float64) *WildAction {
	a.jumpParams.offGrid = chance
	return a.finalizer()
}

// WithClone can be used to indicate jumping wilds will clone instead of just jump.
// The function will panic if it is not a jumping wild action.
func (a *WildAction) WithClone() *WildAction {
	a.jumpParams.clone = true
	return a.finalizer()
}

// WithAlternate can be used to add an alternative action for cases where they need to be mutually exclusive.
// E.g. a wild symbol that triggers a bonus game when there are three, but it triggers a free spin when there are two.
// Alternative actions can be nested. E.g. 3 symbols triggers A, 2 symbols triggers B, 1 symbol triggers C.
// Circular references are not supported and will likely result in a panic due to stack overflow.
func (a *WildAction) WithAlternate(alt *WildAction) *WildAction {
	a.alternate = alt
	return a.finalizer()
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *WildAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.jumping:
		if spin.CountSymbols(a.jumpSymbols) > 0 && a.doJump(spin) {
			return a
		}

	default:
		if spin.newWilds >= a.wildCount {
			if spin.CountSymbol(a.symbol) >= a.wildCount {
				if !a.needHero || spin.newHeroes > 0 {
					return a
				}
			}
		}
	}

	if a.alternate != nil {
		return a.alternate.Triggered(spin)
	}

	return nil
}

// Expand performs the wild expansion.
func (a *WildAction) Expand(spin *Spin) {
	if a.expandToReel {
		a.doExpandToReel(spin)
	}
	if a.expandToSymbol {
		a.doExpandToSymbol(spin)
	}
}

// doExpandToReel expands the wild symbol across unlocked reels containing the wild symbol and locks those reels.
func (a *WildAction) doExpandToReel(spin *Spin) {
	reels, rows := spin.slots.reelCount, spin.slots.rowCount
	for reel := 0; reel < reels; reel++ {
		if !spin.locked[reel] {
			offset := reel * rows
			max := offset + rows
			for offset < max {
				if spin.indexes[offset] == a.symbol {
					for offs := reel * rows; offs < max; offs++ {
						spin.indexes[offs] = a.symbol
					}
					if a.lockReels {
						spin.locked[reel] = true
					}
					break
				}
				offset++
			}
		}
	}
}

// doExpandToSymbol expands the wild symbol using the configured expansion grid with a selected symbol.
func (a *WildAction) doExpandToSymbol(spin *Spin) {
	symbol := spin.stickySymbol
	if spin.symbols.GetSymbol(symbol) == nil {
		panic("this should never happen")
	}

	reels, rows := spin.slots.reelCount, spin.slots.rowCount

	// expand any wild symbols according to the expandGrid
	var offset uint8
	for reel := 0; reel < reels; reel++ {
		if !spin.locked[reel] {
			for row := 0; row < rows; row++ {
				if spin.indexes[offset] == a.symbol {
					for _, offsets := range a.expandGrid {
						reel2, row2 := reel+offsets[0], row+offsets[1]
						if reel2 >= 0 && reel2 < reels && row2 >= 0 && row2 < rows {
							offset2 := reel2*rows + row2
							// don't overwrite other wild symbols!
							if (reel2 == reel && row2 == row) || spin.indexes[offset2] != a.symbol {
								spin.indexes[offset2] = symbol
								spin.effects[offset2] = a.effectKind
							}
						}
					}
				}
				offset++
			}
		}
	}
}

func (a *WildAction) doJump(spin *Spin) bool {
	jumps := make(GridJumps, 0, 100)

	var offset int
	for reel := 0; reel < spin.reelCount; reel++ {
		max := offset + int(spin.mask[reel])
		for offs := offset; offs < max; offs++ {
			if a.matchesJumpingWild(spin.indexes[offs]) {
				jumps = jumps.TestGridJump(a.jumpParams, uint8(offs), spin)
			}
		}
		offset += spin.rowCount
	}

	return jumps.Jump()
}

func (a *WildAction) matchesJumpingWild(s1 utils.Index) bool {
	for _, s2 := range a.jumpSymbols {
		if s1 == s2 {
			return true
		}
	}
	return false
}

// Payout adds a payout to the results based on the number of wild symbols.
// The function panics if the symbol indicated in the action cannot be found.
func (a *WildAction) Payout(spin *Spin, res *results.Result) SpinActioner {
	count := spin.CountSymbol(a.symbol)
	if count >= a.wildCount {
		res.AddPayouts(WildSymbolPayout(a.wildPayout, spin.getMultiplier(0), a.symbol, count, spin))
		return a
	}

	if a.alternate != nil {
		s := a.alternate.(*WildAction)
		return s.Payout(spin, res)
	}

	return nil
}

// CanPayout returns true if the action can trigger a payout.
func (a *WildAction) CanPayout() bool {
	return a.payout
}

func newWild(symbol utils.Index) *WildAction {
	a := &WildAction{}
	a.init(ExpandAfter, FreeSpins, reflect.TypeOf(a).String())
	a.symbol = symbol
	return a
}

func (a *WildAction) finalizer() *WildAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.payout:
		b.WriteString(",payout=true")
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.wildCount)))

	case a.expand:
		b.WriteString(",expand=true")
		b.WriteString(",nrOfSpins=")
		b.WriteString(strconv.Itoa(int(a.nrOfSpins)))
		b.WriteString(",altSymbols=")
		if a.altSymbols {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.wildCount)))
		b.WriteString(",needHero=")
		if a.needHero {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
		b.WriteString(",expandToReel=")
		if a.expandToReel {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}

	case a.transform:
		b.WriteString(",transform=true")
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(int(a.wildCount)))
		b.WriteString(",grid=")
		j, _ := json.Marshal(a.expandGrid)
		b.Write(j)

	case a.jumping:
		b.WriteString(",jumping=true")
		b.WriteString(",symbols=")
		j, _ := json.Marshal(a.jumpSymbols)
		b.Write(j)
		b.WriteString(",direction=")
		b.WriteString(a.jumpParams.direction.String())
		b.WriteString(",min=")
		b.WriteString(strconv.Itoa(int(a.jumpParams.minJump)))
		b.WriteString(",max=")
		b.WriteString(strconv.Itoa(int(a.jumpParams.maxJump)))
		if a.jumpParams.offGrid > 0 {
			b.WriteString(",offGrid=")
			b.WriteString(fmt.Sprintf("%g", a.jumpParams.offGrid))
		}
		if a.jumpParams.onSymbols {
			b.WriteString(",onSymbols=true")
		}
		if a.jumpParams.clone {
			b.WriteString(",clone=true")
		}
		if a.jumpParams.refill {
			b.WriteString(",refill=true")
		}
	}

	a.config = b.String()
	return a
}

// WildActions is a convenience type for a slice of wild actions.
type WildActions []*WildAction

// PurgeWildActions returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeWildActions(input WildActions, capacity int) WildActions {
	if cap(input) < capacity {
		return make(WildActions, 0, capacity)
	}
	return input[:0]
}
