package slots

import (
	"bytes"
	"math"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// GridAction defines actions performed on the grid before or after regular payouts are processed.
// Most actions must be scheduled before regular payouts.
// Only actions that do not introduce or remove payouts can be scheduled after regular payouts.
// An example of this is the SuperX feature in ChaCha Bomb, which triggers "private" re-spins until the triggering symbol no longer appears.
type GridAction struct {
	SpinAction
	shape      GridOffsets
	centers    GridOffsets
	fakeChance float64
	fakes      []FakeSpin
}

// FakeSpin represents the data for a fake grid, to replace no-payline grids with a teaser.
type FakeSpin struct {
	Indexes        utils.Indexes `json:"indexes,omitempty"`        // replacement grid; symbol 0 means do not replace!
	MatchReels     utils.UInt8s  `json:"matchReels,omitempty"`     // reels that must match with the replacement grid.
	MatchSymbol    utils.Index   `json:"matchSymbol,omitempty"`    // symbol to match on instead of replacement grid.
	MatchInverse   bool          `json:"matchInverse,omitempty"`   // invert the reel matching algoritm.
	ReplaceSymbols utils.Indexes `json:"replaceSymbols,omitempty"` // symbol(s) to replace in the given reels.
	ReplaceReels   utils.UInt8s  `json:"replaceReels,omitempty"`   // reels to replace symbols from (note they are 1-based).
}

// NewSuperShapeAction instantiates a new grid shape action.
// This is used for things like the SuperX feature of ChaCha Bomb.
// It is a combination of setting sticky indicators and clearing part of the grid for a free spin.
// The given shape defines the required layout with the same symbol to trigger.
func NewSuperShapeAction(shape, centers GridOffsets) *GridAction {
	a := newGridAction()
	a.result = SuperRefill
	a.shape = shape
	a.centers = centers
	return a.finalize()
}

// NewShapeRefillAction instantiates a new grid shape action.
// This is used for things like the Wild Respin feature of Magic Devil.
// It is a combination of setting sticky indicators and clearing part of the grid for a free spin.
// The given shape defines the required layout with the symbol to trigger.
func NewShapeRefillAction(symbol utils.Index, shape, centers GridOffsets, altSymbols bool) *GridAction {
	a := newGridAction()
	a.stage = TestClearance
	a.result = Refill
	a.symbol = symbol
	a.altSymbols = altSymbols
	a.shape = shape
	a.centers = centers
	return a.finalize()
}

// NewFakeSpinAction instantiates a new fake spin action.
// This is used to tease players with an impossible situation that implies it almost paid out a lot but pays nothing.
func NewFakeSpinAction(chance float64, fakes ...FakeSpin) *GridAction {
	a := newGridAction()
	a.stage = AwardBonuses
	a.result = GridModified
	a.fakeChance = chance
	a.fakes = fakes
	return a.finalize()
}

// WithAlternate can be used to add an alternative action for cases where they need to be mutually exclusive.
func (a *GridAction) WithAlternate(alt *GridAction) *GridAction {
	a.alternate = alt
	return a.finalize()
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *GridAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.shape != nil:
		if a.result == SuperRefill {
			if a.superShapeTriggered(spin) {
				return a
			}
		} else {
			if a.shapeRefillTriggered(spin) {
				return a
			}
		}

	case a.fakeChance > 0 && a.fakes != nil:
		if spin.prng.IntN(10000) < int(math.Round(a.ModifyChance(a.fakeChance, spin)*100)) {
			switch len(a.fakes) {
			case 1:
				if fake := a.fakes[0]; fake.matches(spin) {
					return a.replaceGrid(spin, fake)
				}

			default:
				for ix := 0; ix < 25; ix++ {
					fake := a.fakes[spin.prng.IntN(len(a.fakes))]
					if fake.matches(spin) {
						return a.replaceGrid(spin, fake)
					}
				}
			}
		} else {
			if alt := a.alternate; alt != nil {
				if alt.CanTrigger(spin) {
					return alt.Triggered(spin)
				}
			}
		}
	}

	return nil
}

// superShapeTriggered tests if a super shape action has triggered.
func (a *GridAction) superShapeTriggered(spin *Spin) bool {
	var found bool

	if spin.superSymbol != utils.MaxIndex {
		// already found it before; need to check for new occurrences.
		for ix := range spin.indexes {
			if !spin.sticky[ix] && spin.indexes[ix] == spin.superSymbol {
				found = true
				spin.sticky[ix] = true
			}
		}

		if found && a.CanClear(spin) {
			// clear other symbols to trigger the refill spin.
			for ix := range spin.indexes {
				if !spin.sticky[ix] {
					spin.indexes[ix] = 0
				}
			}
		}

		return found
	}

	// resetData any previously set sticky flags.
	spin.ResetSticky()

	// test the shape in every possible position to see if it matched.
	var symbol utils.Index
	var center Offsets
	rows := spin.rowCount

	for _, center = range a.centers {
		symbol, found = utils.MaxIndex, true
		for _, offsets := range a.shape {
			reel, row := center[0]+offsets[0], center[1]+offsets[1]
			offset := reel*rows + row
			if symbol == utils.MaxIndex {
				symbol = spin.indexes[offset]
			} else if spin.indexes[offset] != symbol {
				found = false
				break
			}
		}
		if found {
			break
		}
	}

	// only if it's the proper shape and the symbol doesn't occur anywhere else!
	if !found || int(spin.CountSymbol(symbol)) != len(a.shape) {
		return false
	}

	// remember the symbol we found
	spin.superSymbol = symbol

	if a.CanSticky(spin) {
		// mark the shape as sticky.
		for _, offsets := range a.shape {
			reel, row := center[0]+offsets[0], center[1]+offsets[1]
			spin.sticky[reel*rows+row] = true
			spin.superShape[reel*rows+row] = true
		}
	}

	if a.CanClear(spin) {
		// clear other symbols to trigger the refill spin.
		for ix := range spin.indexes {
			if !spin.sticky[ix] {
				spin.indexes[ix] = 0
			}
		}
	}

	return true
}

// shapeRefillTriggered tests if a shape refill action has triggered.
func (a *GridAction) shapeRefillTriggered(spin *Spin) bool {
	var found bool

	// resetData any previously set sticky flags.
	spin.ResetSticky()

	// test the shape in every possible position to see if it matched.
	var center Offsets
	rows := spin.rowCount

	for _, center = range a.centers {
		found = true
		for _, offsets := range a.shape {
			reel, row := center[0]+offsets[0], center[1]+offsets[1]
			offset := reel*rows + row
			if spin.indexes[offset] != a.symbol {
				found = false
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return false
	}

	if a.CanSticky(spin) {
		// mark the shape as sticky.
		for _, offsets := range a.shape {
			reel, row := center[0]+offsets[0], center[1]+offsets[1]
			spin.sticky[reel*rows+row] = true
		}
	}

	if a.CanClear(spin) {
		// clear other symbols to trigger the refill spin.
		for ix := range spin.indexes {
			if !spin.sticky[ix] {
				spin.indexes[ix] = 0
			}
		}
	}

	return true
}

func (f *FakeSpin) matches(spin *Spin) bool {
	if f.MatchSymbol > 0 {
		if f.MatchInverse {
			for ix := range f.MatchReels {
				reel := int(f.MatchReels[ix] - 1)
				offset := reel * spin.rowCount
				max := offset + int(spin.mask[reel])
				for ; offset < max; offset++ {
					if spin.indexes[offset] == f.MatchSymbol {
						return false
					}
				}
			}
		} else {
			for ix := range f.MatchReels {
				reel := int(f.MatchReels[ix] - 1)
				offset := reel * spin.rowCount
				max := offset + int(spin.mask[reel])
				for ; offset < max; offset++ {
					if spin.indexes[offset] != f.MatchSymbol {
						return false
					}
				}
			}
		}
	} else {
		if f.MatchInverse {
			for ix := range f.MatchReels {
				reel := int(f.MatchReels[ix] - 1)
				offset := reel * spin.rowCount
				max := offset + int(spin.mask[reel])
				for ; offset < max; offset++ {
					if id := f.Indexes[offset]; id != 0 && spin.indexes[offset] == id {
						return false
					}
				}
			}
		} else {
			for ix := range f.MatchReels {
				reel := int(f.MatchReels[ix] - 1)
				offset := reel * spin.rowCount
				max := offset + int(spin.mask[reel])
				for ; offset < max; offset++ {
					if id := f.Indexes[offset]; id != 0 && spin.indexes[offset] != id {
						return false
					}
				}
			}
		}
	}
	return true
}

func (a *GridAction) replaceGrid(spin *Spin, fake FakeSpin) SpinActioner {
	if len(fake.ReplaceSymbols) > 0 && len(fake.ReplaceReels) > 0 {
		replace := func(symbol utils.Index) bool {
			for _, id := range fake.ReplaceSymbols {
				if symbol == id {
					return true
				}
			}
			return false
		}

		reels := spin.reels
		if spin.altActive {
			reels = spin.altReels
		}

		for _, r := range fake.ReplaceReels {
			reel := int(r - 1)
			offset := reel * spin.rowCount
			max := offset + int(spin.mask[reel])
			for ; offset < max; offset++ {
				for replace(spin.indexes[offset]) {
					spin.indexes[offset] = reels[reel].RandomIndex(spin.prng)
				}
			}
		}
	}

	for ix := range fake.Indexes {
		if id := fake.Indexes[ix]; id > 0 && ix < len(spin.indexes) {
			spin.indexes[ix] = id
		}
	}

	return a
}

func (a *GridAction) finalize() *GridAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.shape != nil:
		if a.result == SuperRefill {
			b.WriteString(",superShape=true")
		} else {

			b.WriteString(",shapeRefill=true,symbol=")
			b.WriteString(strconv.Itoa(int(a.symbol)))
			b.WriteString(",altSymbols=")
			if a.altSymbols {
				b.WriteString("true")
			} else {
				b.WriteString("false")
			}
		}

		b.WriteString(",shape=")
		j, _ := json.Marshal(a.shape)
		b.Write(j)
		b.WriteString(",centers=")
		j, _ = json.Marshal(a.centers)
		b.Write(j)

	case a.fakeChance > 0:
		b.WriteString(",fakeSpin=true,chance=")
		b.WriteString(strconv.FormatFloat(a.fakeChance, 'g', -1, 64))
		b.WriteString(",fakes=")
		j, _ := json.Marshal(a.fakes)
		b.Write(j)
	}

	a.config = b.String()
	return a
}

func newGridAction() *GridAction {
	a := &GridAction{}
	a.init(TestGrid, Processed, reflect.TypeOf(a).String())
	return a
}

// GridActions is a convenience type for a slice of grid actions.
type GridActions []*GridAction

// PurgeGridActions returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeGridActions(input GridActions, capacity int) GridActions {
	if cap(input) < capacity {
		return make(GridActions, 0, capacity)
	}
	return input[:0]
}
