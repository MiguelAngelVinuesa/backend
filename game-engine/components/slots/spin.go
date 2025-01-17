package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// Spinner is the interface to fill (part of) a grid.
type Spinner interface {
	Spin(spin *Spin, indexes utils.Indexes)
}

// SpinKind represents the kind of spin.
type SpinKind uint8

const (
	// RegularSpin represents a regular spin.
	RegularSpin SpinKind = iota
	// FirstSpin represents the first spin in a double-spin game.
	FirstSpin
	// SecondSpin represents the second spin in a double-spin game.
	SecondSpin
	// SuperSpin represents a super spin based on a grid shape (semi-free spin).
	SuperSpin
	// RefillSpin represents a spin that re-fills part of the grid (semi-free spin).
	RefillSpin
	// FreeSpin represents a regular free spin.
	FreeSpin
	// FirstFreeSpin represents the first spin of a free spin in a double-spin game.
	FirstFreeSpin
	// SecondFreeSpin represents the second spin of a free spin in a double-spin game.
	SecondFreeSpin
)

const (
	NoEffect uint8 = iota
	BombEffect
)

// IsFirst returns if the spin kind indicates the first spin of a double spin.
func (s SpinKind) IsFirst() bool {
	return s == FirstSpin || s == FirstFreeSpin
}

// IsSecond returns if the spin kind indicates the second spin of a double spin.
func (s SpinKind) IsSecond() bool {
	return s == SecondSpin || s == SecondFreeSpin
}

// AcquireSpin instantiates a spin result from the memory pool.
// It "spins" all reels of the slot machine and calculates the initial result.
func AcquireSpin(slots *Slots, prng interfaces.Generator) *Spin {
	reels, rows := slots.reelCount, slots.rowCount

	s := spinsProducer.Acquire().(*Spin)
	s.multiplierNeedsWild = slots.multiplierOnWildsOnly
	s.prng = prng
	s.slots = slots
	s.reelCount = reels
	s.rowCount = rows
	s.symbols = slots.symbols
	s.altSymbols = slots.altSymbols
	s.paylines = slots.paylines
	s.gridDef = slots.gridDef
	s.mask = s.gridDef.mask
	s.spinner = slots.spinner
	s.refiller = slots.refiller

	s.reels = PurgeReels(s.reels, reels)[:reels]
	s.altReels = PurgeReels(s.altReels, reels)
	s.locked = utils.PurgeBools(s.locked, reels)[:reels]
	s.hot = utils.PurgeBools(s.hot, reels)[:reels]

	size := slots.gridDef.gridSize
	m := object.NormalizeSize(size, 8)
	s.sticky = utils.PurgeBools(s.sticky, m)[:size]
	s.superShape = utils.PurgeBools(s.superShape, m)[:size]
	s.payouts = utils.PurgeUInt8s(s.payouts, m)[:size]
	s.effects = utils.PurgeUInt8s(s.effects, m)[:size]
	s.jumps = utils.PurgeUInt8s(s.jumps, m)[:size]
	s.indexes = utils.PurgeIndexes(s.indexes, m)[:size]
	s.injections = utils.PurgeIndexes(s.injections, m)[:size]
	s.multipliers = utils.PurgeUInt16s(s.multipliers, m)

	for reel := range s.reels {
		s.reels[reel] = AcquireReel(uint8(reel), uint8(rows), s.symbols, slots.noRepeat)
	}
	if s.altSymbols != nil {
		s.altReels = s.altReels[:reels]
		for reel := range s.altReels {
			s.altReels[reel] = AcquireReel(uint8(reel), uint8(rows), s.altSymbols, slots.noRepeat)
		}
	}

	clear(s.locked)
	clear(s.indexes)
	clear(s.superShape)

	s.ResetSpin()
	s.initialSpin()
	s.MismatchPaylines()
	s.CountSpecials()

	s.spinSeq = 0
	return s
}

// SetSpinner sets up a special spinner function to fill the grid.
// By default, or when this function has been called with nil, the built-in spin function will be used.
// The built-in spin function uses the symbol weights for each reel to pick random symbols.
// This option can be used to inject pre-determined data into the game, or select from pre-calculated start grids.
func (s *Spin) SetSpinner(spinner Spinner) *Spin {
	s.spinner = spinner
	return s
}

// SetRefiller sets up a special refill function to fill empty positions on the grid.
// By default, or when the function is called with nil input, the built-in refill function will be used.
// The built-in refill function uses the symbol weights for each reel to pick random symbols.
// This option can be used to inject pre-determined data into the game, or select from pre-calculated follow-up grids.
func (s *Spin) SetRefiller(refiller Spinner) *Spin {
	s.refiller = refiller
	return s
}

// SetGamer sets up the game round interface.
func (s *Spin) SetGamer(gamer interfaces.Gamer) {
	s.gamer = gamer
}

// SetDebug sets/resets the debug indicator.
func (s *Spin) SetDebug(debug bool) {
	s.debug = debug
}

// SetBonusBuy sets the bonus buy flag (if configured).
func (s *Spin) SetBonusBuy(bonus uint8) {
	if s.slots.bonusBuy && s.slots.flagBB > 0 {
		s.roundFlags[s.slots.flagBB] = int(bonus)
	}
}

// SetRoundFlag sets a specific round flag.
func (s *Spin) SetRoundFlag(flag int, value int) {
	s.roundFlags[flag] = value
}

// GridSize returns the number of reels and rows for the grid.
func (s *Spin) GridSize() (int, int) {
	return s.reelCount, s.rowCount
}

// ReelSize returns the number of rows for a specific reel.
// The reel is 0-based here!
func (s *Spin) ReelSize(reel uint8) int {
	return int(s.mask[reel])
}

// PRNG returns the random number generator.
func (s *Spin) PRNG() interfaces.Generator {
	return s.prng
}

// Reels returns the currently selected reels.
func (s *Spin) Reels() Reels {
	if s.altActive {
		return s.altReels
	}
	return s.reels
}

func (s *Spin) SpinSeq() uint64 {
	return s.spinSeq
}

func (s *Spin) ResultCount() int {
	return s.gamer.ResultCount()
}

func (s *Spin) FreeSpins() uint64 {
	return s.freeSpins
}

func (s *Spin) TotalPayout() float64 {
	return s.gamer.TotalPayout()
}

// Debug initializes the spin as a debug spin with an initial array.
// The outcome is not guaranteed if len(indexes) < reelCount * rowCount.
func (s *Spin) Debug(indexes utils.Indexes) {
	s.debug = true
	s.debugInitial = true
	s.spinSeq++

	s.LockReels()
	s.resetPayouts()

	copy(s.indexes, indexes)

	s.CountSpecials()
}

// TestChance2 generates a random number to test against the given chance (max 2 decimals).
// As special cases, chance <= 0 will always return false, and chance >= 100 will always return true,
// thus eliminating a round trip to the PRNG.
func (s *Spin) TestChance2(chance float64) bool {
	if chance > 99.99 {
		return true
	}
	return float64(s.prng.IntN(10000)) < chance*100
}

// TestChance4 generates a random number to test against the given chance (max 4 decimals).
// As special cases, chance <= 0 will always return false, and chance >= 100 will always return true,
// thus eliminating a round trip to the PRNG.
func (s *Spin) TestChance4(chance float64) bool {
	if chance > 99.9999 {
		return true
	}
	return float64(s.prng.IntN(1000000)) < chance*10000
}

// ResetSpin resets the spin to its initial state.
func (s *Spin) ResetSpin() {
	s.kind = RegularSpin
	s.superSymbol = utils.MaxIndex
	s.bonusSymbol = utils.MaxIndex
	s.altActive = false
	s.expanded = false
	s.freeSpins = 0
	s.spinSeq = 0
	s.progressLevel = 0
	s.multiplier = 0.0

	s.resetHot()
	s.ResetSticky()
	s.ResetEffects()
	s.resetMultipliers()

	clear(s.roundFlags)
}

// RestoreState restores the state of a previous spin.
func (s *Spin) RestoreState(state *SpinState) {
	s.stickySymbol = state.stickySymbol

	if len(state.indexes) == len(s.indexes) {
		copy(s.indexes, state.indexes)
	}

	if len(state.sticky) == len(s.sticky) {
		copy(s.sticky, state.sticky)
	}
}

// SetKind sets the kind of spin.
func (s *Spin) SetKind(kind SpinKind) {
	s.kind = kind
}

// Kind returns the spin kind.
func (s *Spin) Kind() SpinKind {
	return s.kind
}

// GetSymbols returns the active symbol set.
func (s *Spin) GetSymbols() *SymbolSet {
	if s.altActive {
		return s.altSymbols
	}
	return s.symbols
}

// GetSymbol returns the requested symbol or nil if it can't be found.
func (s *Spin) GetSymbol(symbol utils.Index) *Symbol {
	if s.altActive {
		return s.altSymbols.GetSymbol(symbol)
	}
	return s.symbols.GetSymbol(symbol)
}

// Indexes returns the current grid with symbol ids.
func (s *Spin) Indexes() utils.Indexes {
	return s.indexes
}

// Stickies returns the sticky indicators from the current grid.
func (s *Spin) Stickies() []bool {
	return s.sticky
}

// Multipliers returns the current multipliers from the grid.
func (s *Spin) Multipliers() object.Uint16s {
	return s.multipliers
}

// ModifyTile allows to modify a grid tile from outside the "components" module.
func (s *Spin) ModifyTile(offset int, symbol utils.Index, sticky bool, multiplier uint16) {
	l := len(s.indexes)
	if offset < 0 || offset >= l {
		return
	}

	s.indexes[offset] = symbol
	s.sticky[offset] = sticky

	if multiplier > 0 {
		if len(s.multipliers) == 0 {
			s.multipliers = utils.PurgeUInt16s(s.multipliers, l)
			s.multipliers[offset] = multiplier
		}
	}
}

// LockReels locks the specified reels and unlocks all other reels.
// Note that locks are 1-based and not 0-based.
// The function panics if any given lock == 0 or > len(reels).
func (s *Spin) LockReels(locks ...uint8) {
	clear(s.locked)
	for lock := range locks {
		s.locked[locks[lock]-1] = true
	}
}

// HotReel marks the specified reel as hot, using the actual 0-based reel index.
func (s *Spin) HotReel(reel uint8) {
	s.hot[reel] = true
}

// CountSymbol counts the number of occurrences of the given symbol in the last result.
// The function does not check if the given symbol id is valid. Only its occurrence is counted.
func (s *Spin) CountSymbol(symbol utils.Index) uint8 {
	var count uint8
	for ix := range s.indexes {
		if s.indexes[ix] == symbol {
			count++
		}
	}
	return count
}

// CountSymbols counts the number of occurrences of the given symbols in the last result.
// The function does not check if the given symbol ids are valid. Only their occurrence is counted.
func (s *Spin) CountSymbols(symbols utils.Indexes) uint8 {
	var count uint8
	for _, symbol := range symbols {
		for ix := range s.indexes {
			if s.indexes[ix] == symbol {
				count++
			}
		}
	}
	return count
}

// CountSymbolInReels counts the number of occurrences of the given symbol in the given reels.
// The function does not check if the given symbol id is valid. Only its occurrence is counted.
// Note that reels are 1-based!
func (s *Spin) CountSymbolInReels(symbol utils.Index, reels utils.UInt8s) uint8 {
	var count uint8
	for _, reel := range reels {
		m := int(reel) * s.rowCount
		for offset := m - s.rowCount; offset < m; offset++ {
			if s.indexes[offset] == symbol {
				count++
			}
		}
	}
	return count
}

// CountBonusSymbol counts the number of reels in the last spin that the given bonus symbol occurs in.
// The function does not check if the given symbol id is valid. Only its count across the reels is determined.
// If the "hot reels are bonus" feature is on, every hot reel will add to the count.
func (s *Spin) CountBonusSymbol() uint8 {
	var count uint8
	for reel := 0; reel < s.reelCount; reel++ {
		if s.slots.hotReelsAsBonusSymbol && s.hot[reel] {
			count++
		} else {
			offset := reel * s.rowCount
			for m := offset + s.rowCount; offset < m; offset++ {
				if s.indexes[offset] == s.bonusSymbol {
					count++
					break
				}
			}
		}
	}
	return count
}

// ExpandBonusSymbol expands the bonus symbol across the reels that contain the bonus symbol or are hot.
func (s *Spin) ExpandBonusSymbol() {
	var expand bool
	for reel := 0; reel < s.reelCount; reel++ {
		expand = false
		if s.slots.hotReelsAsBonusSymbol && s.hot[reel] {
			expand = true
		} else {
			offset := reel * s.rowCount
			m := offset + s.rowCount
			for ; offset < m; offset++ {
				if s.indexes[offset] == s.bonusSymbol {
					expand = true
					break
				}
			}
		}
		if expand {
			offset := reel * s.rowCount
			m := offset + s.rowCount
			for ; offset < m; offset++ {
				s.indexes[offset] = s.bonusSymbol
			}
			s.expanded = true
		}
	}
}

// IsExpanded returns whether some symbol on the grid was expanded.
func (s *Spin) IsExpanded() bool {
	return s.expanded
}

// IsLocked returns whether the given reel index is a locked reel.
// The given reel index must be 0-based.
func (s *Spin) IsLocked(reel uint8) bool {
	return s.locked[reel]
}

// Locked returns the indexes of locked reels if any.
// The given slice is re-used if it is not nil and of sufficient capacity.
// Note that locks are 1-based and not 0-based.
func (s *Spin) Locked(input utils.UInt8s) utils.UInt8s {
	input = utils.PurgeUInt8s(input, len(s.locked))
	for ix := range s.locked {
		if s.locked[ix] {
			input = append(input, uint8(ix+1))
		}
	}
	return input
}

// Hot returns the indexes of hot reels if any.
// The given slice is re-used if it is not nil and of sufficient capacity.
// Note that hot indexes are 1-based and not 0-based; so if reel 0 is hot, the result wil contain the index 1.
func (s *Spin) Hot(input utils.UInt8s) utils.UInt8s {
	input = utils.PurgeUInt8s(input, s.reelCount)
	for ix := range s.hot {
		if s.hot[ix] {
			input = append(input, uint8(ix+1))
		}
	}
	return input
}

// ResetSticky resets all the sticky indicators.
func (s *Spin) ResetSticky() {
	s.stickySymbol = utils.MaxIndex
	clear(s.sticky)
}

// StickySymbol returns the currently selected sticky symbol or MaxIndex if not selected.
func (s *Spin) StickySymbol() utils.Index {
	return s.stickySymbol
}

// HasSticky returns whether there are any positions in the grid flagged as sticky.
func (s *Spin) HasSticky() bool {
	for ix := range s.sticky {
		if s.sticky[ix] {
			return true
		}
	}
	return false
}

// StickyCount returns the number of positions in the grid flagged as sticky.
func (s *Spin) StickyCount() uint8 {
	var count uint8
	for ix := range s.sticky {
		if s.sticky[ix] {
			count++
		}
	}
	return count
}

// Sticky fills the given slice with the sticky indicators or returns an empty slice if no symbols are sticky.
// The given slice is re-used if it is not nil and of sufficient capacity, otherwise a new slice will be allocated.
// Super shape sticky tiles are marked with 2, other sticky tiles with 1.
func (s *Spin) Sticky(input utils.UInt8s) utils.UInt8s {
	l := len(s.sticky)
	m := max(l, 16)
	input = utils.PurgeUInt8s(input, m)

	if s.HasSticky() {
		input = input[:l]
		clear(input)
		for ix := 0; ix < l; ix++ {
			if s.sticky[ix] {
				if s.superShape[ix] {
					input[ix] = 2
				} else {
					input[ix] = 1
				}
			}
		}
	}

	return input
}

// ResetEffects resets all the special effect indicators and symbol injections.
func (s *Spin) ResetEffects() {
	clear(s.effects)
	clear(s.jumps)
	clear(s.injections)
}

// HasEffect returns whether there are any special effects on the grid.
func (s *Spin) HasEffect() bool {
	for ix := range s.effects {
		if s.effects[ix] != NoEffect {
			return true
		}
	}
	return false
}

// Effects fills the given slice with the special effect indicators or returns an empty slice if there are no effects.
// The given slice is re-used if it is not nil and of sufficient capacity, otherwise a new slice will be allocated.
func (s *Spin) Effects(input utils.UInt8s) utils.UInt8s {
	l := len(s.effects)
	m := max(l, 16)
	input = utils.PurgeUInt8s(input, m)

	if s.HasEffect() {
		input = input[:l]
		copy(input, s.effects)
	}
	return input
}

// ScatterMap updates the given slice with the offsets of the symbol matches in the grid.
// The given slice is re-used if it is not nil and of sufficient capacity.
func (s *Spin) ScatterMap(symbol utils.Index, count uint8, input utils.UInt8s) utils.UInt8s {
	if count < 1 {
		panic(consts.MsgInvalidSymbolCount)
	}

	l := int(count)
	m := max(l, 16)
	input = utils.PurgeUInt8s(input, m)[:l]

	var iy int
	for ix := range s.indexes {
		if s.indexes[ix] == symbol {
			input[iy] = uint8(ix)
			iy++
		}
	}
	return input
}

// SetBonusSymbol forces the bonus symbol to be set.
// This function exists for testing purposes only.
func (s *Spin) SetBonusSymbol(symbol utils.Index, altSymbols bool) {
	s.bonusSymbol = symbol
	if altSymbols {
		s.altActive = true
	}
}

// BonusSymbol returns the currently selected bonus symbol, or MaxIndex if not selected.
func (s *Spin) BonusSymbol() utils.Index {
	return s.bonusSymbol
}

// SuperSymbol returns the currently selected super symbol, or MaxIndex if not selected.
func (s *Spin) SuperSymbol() utils.Index {
	return s.superSymbol
}

// ResetSuper resets the super symbol.
func (s *Spin) ResetSuper() {
	s.superSymbol = utils.MaxIndex
	clear(s.superShape)
}

// SetFreeSpins sets the remaining free spins counter.
func (s *Spin) SetFreeSpins(count uint64) {
	s.freeSpins = count
}

// Multiplier returns the current value of the overall spin round multiplier.
func (s *Spin) Multiplier() float64 {
	return s.multiplier
}

// MismatchPaylines makes sure that none of the paylines will match for the current result by replacing symbols.
func (s *Spin) MismatchPaylines() {
	// remove all wild and split symbols.
	var offset int
	for offset < len(s.indexes) {
		if id := s.indexes[offset]; id != utils.NullIndex {
			if symbol := s.symbols.GetSymbol(id); symbol != nil {
				if symbol.IsWild() || len(symbol.wildFor) > 0 {
					s.indexes[offset] = s.reels[0].weighting.RandomIndex(s.prng)
					offset--
				}
			}
		}
		offset++
	}

	if s.slots.directions == PayLTR || s.slots.directions == PayBoth {
		// replace symbols on reel 1 that already appear on reel 2
		var offset1 uint8
		for offset1 < s.mask[0] {
			var replace bool
			for offset2 := s.rowCount; offset2 < s.rowCount*2; offset2++ {
				if id1, id2 := s.indexes[offset1], s.indexes[offset2]; id1 != 0 && id2 == id1 {
					replace = true
					break
				}
			}
			if replace {
				for {
					id := s.reels[0].weighting.RandomIndex(s.prng)
					symbol := s.symbols.GetSymbol(id)
					if symbol == nil || (!symbol.IsWild() && len(symbol.wildFor) == 0) {
						s.indexes[offset1] = id
						break
					}
				}
			} else {
				offset1++
			}
		}
	}

	if s.slots.directions == PayRTL || s.slots.directions == PayBoth {
		// replace symbols on last reel that already appear on the last but one reel
		max1 := s.reelCount * s.rowCount
		offset1 := max1 - s.rowCount
		max2 := offset1
		for offset1 < max1 {
			var replace bool
			if id := s.indexes[offset1]; id != utils.NullIndex {
				for offset2 := max2 - s.rowCount; offset2 < max2; offset2++ {
					if s.indexes[offset2] == id {
						replace = true
						break
					}
				}
			}
			if replace {
				for {
					symbol := s.symbols.GetSymbol(s.reels[s.reelCount-1].weighting.RandomIndex(s.prng))
					if !symbol.IsWild() && len(symbol.wildFor) == 0 {
						s.indexes[offset1] = symbol.id
						break
					}
				}
			} else {
				offset1++
			}
		}
	}
}

// ForcePaidAction inserts the required symbols into the current result to satisfy the PaidAction.
// This is used to "buy" a bonus game.
func (s *Spin) ForcePaidAction(paid *PaidAction) {
	if paid.flag >= 0 {
		s.roundFlags[paid.flag] = paid.flagValue
	}

	count := paid.triggerCount
	if count == 0 {
		return
	}

	n := s.CountSymbol(paid.symbol)
	if n >= count {
		return
	}

	a := ReviseAction{
		SpinAction:    SpinAction{symbol: paid.symbol},
		symbolChances: make([]float64, 0, count),
		genAllowDupes: false,
		generateReels: paid.reels,
	}

	count -= n
	for count > 0 {
		a.symbolChances = append(a.symbolChances, 100)
		count--
	}

	a.doGenerateSymbol2(s)
}

// TestDupes checks if a symbol occurs multiple times in any of the reels.
func (s *Spin) TestDupes(symbol utils.Index) bool {
	for reel := 0; reel < s.reelCount; reel++ {
		offs := reel * s.rowCount
		end := offs + s.rowCount
		count := 0
		for offs < end {
			if s.indexes[offs] == symbol {
				count++
			}
			offs++
		}
		if count > 1 {
			return true
		}
	}
	return false
}

// resetPayouts resets all the payout indicators.
func (s *Spin) resetPayouts() {
	clear(s.payouts)
}

// resetHot resets the hot reel indicators.
func (s *Spin) resetHot() {
	clear(s.hot)
}

// CloneMultipliers makes a deep copy of the multipliers if they exist.
func (s *Spin) CloneMultipliers(in []uint16) []uint16 {
	out := utils.PurgeUInt16s(in, cap(s.indexes))
	if len(s.multipliers) == 0 {
		return out
	}

	out = out[:len(s.indexes)]
	copy(out, s.multipliers)
	return out
}

// resetMultipliers resets the multipliers.
func (s *Spin) resetMultipliers() {
	if len(s.multipliers) == 0 {
		return
	}

	rows := s.rowCount
	var offset int
	for reel := range s.reels {
		m := offset + rows
		if !s.locked[reel] {
			clear(s.multipliers[offset:m])
		}
		offset = m
	}
}

// CascadeFloatingSymbols drops symbols "floating" on higher rows down to the lowest empty position on the reel.
// This ensures that all empty positions are at the top of a reel and all symbols are at the bottom of a reel.
// Empty positions are indicated with the symbolID == 0.
// Sticky symbols do not cascade and remain in the same position.
func (s *Spin) CascadeFloatingSymbols() bool {
	var offset, bottom int
	var cascaded bool

	if len(s.multipliers) == 0 {
		for reel := range s.reels {
			bottom = offset + int(s.mask[reel]) - 1
			upper := bottom

			for bottom >= offset {
				if s.indexes[bottom] == 0 {
					if upper >= bottom {
						upper = bottom - 1
					}
					for upper >= offset && (s.indexes[upper] == 0 || s.sticky[upper]) {
						upper--
					}
					if upper >= offset {
						s.indexes[bottom], s.indexes[upper] = s.indexes[upper], 0
						cascaded = true
						upper--
					} else {
						break
					}
				}
				bottom--
			}

			offset += s.rowCount
		}
	} else {
		for reel := range s.reels {
			bottom = offset + int(s.mask[reel]) - 1
			upper := bottom

			for bottom >= offset {
				if s.indexes[bottom] == 0 {
					if upper >= bottom {
						upper = bottom - 1
					}
					for upper >= offset && (s.indexes[upper] == 0 || s.sticky[upper]) {
						upper--
					}
					if upper >= offset {
						s.indexes[bottom], s.indexes[upper] = s.indexes[upper], 0
						s.multipliers[bottom], s.multipliers[upper] = s.multipliers[upper], 0
						cascaded = true
						upper--
					} else {
						break
					}
				}
				bottom--
			}

			offset += s.rowCount
		}
	}

	return cascaded
}

// Spin "spins" the reels of the slot machine and calculates the result.
// Locked reels will retain their symbols.
// Note that locks are 1-based and not 0-based (same as reels).
// The function panics if any given lock == 0 or > len(reels).
// If the spinner function is set up, it will be used instead of the built-in spin function.
func (s *Spin) Spin(locks ...uint8) {
	s.LockReels(locks...)
	s.initialSpin()
	s.CountSpecials()
}

func (s *Spin) initialSpin() {
	s.resetPayouts()
	s.ResetEffects()

	if s.spinner != nil {
		s.spinSeq++
		s.spinner.Spin(s, s.indexes)
	} else {
		s.Builtin()
	}
}

// Builtin runs the built-in spin mechanism.
func (s *Spin) Builtin() {
	switch {
	case s.HasSticky():
		s.spinSticky()
	case s.altActive:
		s.spin(s.altReels)
	default:
		s.spin(s.reels)
	}
}

// spinSticky preforms a "spin" leaving the sticky symbols in place.
func (s *Spin) spinSticky() {
	s.clearNonSticky()
	s.Refill()
}

// spin performs a "spin" filling unlocked reels with random symbol ids,
// and honouring an optional non-rectangular grid.
func (s *Spin) spin(reels Reels) {
	s.resetMultipliers()

	s.spinSeq++

	var offset int
	for reel := range reels {
		if !s.locked[reel] {
			reels[reel].Spin(s.prng, s.indexes[offset:offset+int(s.mask[reel])])
		}
		offset += s.rowCount
	}
}

// Refill fills cleared locations with new random symbols, honouring the cascading reels feature.
// It basically acts like a "free spin" but with some symbols locked in place,
// and other symbols landing on empty locations lower down on the same reel if the cascading reels feature is active.
func (s *Spin) Refill() {
	s.ResetEffects()

	s.spinSeq++

	switch {
	case s.refiller != nil:
		s.refiller.Spin(s, s.indexes)
	case s.altActive:
		s.refill(s.altReels)
	default:
		s.refill(s.reels)
	}
}

// refill fills empty positions in the grid using the appropriate weighting per reel,
// and honouring an optional non-rectangular grid.
// Any filled positions (symbolID > 0) are untouched.
func (s *Spin) refill(reels Reels) {
	var offset uint8
	for ix, reel := range reels {
		m := offset + s.mask[ix]
		for iy := offset; iy < m; iy++ {
			if s.indexes[iy] == 0 {
				reel.Spin(s.prng, s.indexes[iy:iy+1])
			}
		}
		offset += uint8(s.rowCount)
	}
}

// clearNonSticky clears all symbol positions that are not marked as sticky.
func (s *Spin) clearNonSticky() {
	if len(s.multipliers) == 0 {
		for ix := range s.indexes {
			if !s.sticky[ix] {
				s.indexes[ix] = 0
			}
		}
	} else {
		for ix := range s.indexes {
			if !s.sticky[ix] {
				s.indexes[ix] = 0
				s.multipliers[ix] = 0
			}
		}
	}
}

// CountSpecials counts the wild, hero and scatter symbols on the unlocked reels.
func (s *Spin) CountSpecials() {
	s.newWilds = 0
	s.newHeroes = 0
	s.newScatters = 0

	var offset int
	for reel := range s.locked {
		if !s.locked[reel] {
			s.testSpecialsOnReel(offset, offset+int(s.mask[reel]))
		}
		offset += s.rowCount
	}
}

// testSpecialsOnReel counts the wild, hero and scatter symbols on a specific reel.
func (s *Spin) testSpecialsOnReel(min, max int) {
	for offs := min; offs < max; offs++ {
		if symbol := s.symbols.GetSymbol(s.indexes[offs]); symbol != nil {
			if symbol.IsWild() {
				s.newWilds++
			}
			if symbol.IsScatter() {
				s.newScatters++
			}
			if symbol.IsHero() {
				s.newHeroes++
			}
		}
	}
}

// markWinline marks the matching symbol positions of a winline.
func (s *Spin) markWinline(count uint8, direction PayDirection, rows []uint8) {
	if len(s.payouts) < len(s.indexes) {
		return
	}
	if direction == PayLTR {
		var offset int
		for ix := 0; ix < int(count); ix++ {
			s.payouts[offset+int(rows[ix])] = 1
			offset += s.rowCount
		}
	} else {
		offset := s.rowCount * s.reelCount
		for ix := 0; ix < int(count); ix++ {
			offset -= s.rowCount
			s.payouts[offset+int(rows[s.reelCount-ix-1])] = 1
		}
	}
}

// markAllPayline marks the matching symbol positions of a winning all payline.
func (s *Spin) markAllPayline(count uint8, rows utils.UInt8s) {
	var offset int
	for ix := 0; ix < int(count); ix++ {
		s.payouts[offset+int(rows[ix])] = 1
		offset += s.rowCount
	}
}

func (s *Spin) getMultiplier(wilds int) float64 {
	if s.multiplierNeedsWild && wilds == 0 {
		return 1.0
	}
	return s.multiplier
}

// Spin is used to maintain the grid and status of a spin.
// Spin is not safe for concurrent use across multiple go-routines.
// Keep fields ordered by ascending SizeOf().
type Spin struct {
	debug               bool                 // indicates if spin is running in debug mode.
	debugInitial        bool                 // indicates if spin is running in initial array debug mode.
	altActive           bool                 // indicates that the alternate symbol set is in use.
	expanded            bool                 // indicates a symbol was expanded on one or more reels.
	multiplierNeedsWild bool                 // indicates that the round multiplier only applies for winlines with a wild symbol.
	kind                SpinKind             // the spin kind.
	newWilds            uint8                // count of new wild symbols, excluding locked reels.
	newHeroes           uint8                // count of new hero symbols, excluding locked reels.
	newScatters         uint8                // count of new scatter symbols, excluding locked reels.
	bonusSymbol         utils.Index          // randomly selected bonus symbol during free spins.
	superSymbol         utils.Index          // symbol which triggered a super-shape feature.
	stickySymbol        utils.Index          // symbol which is marked as sticky in the grid.
	reelCount           int                  // copy of reel count.
	rowCount            int                  // copy of row count.
	progressLevel       int                  // current level on the overall progress meter.
	freeSpins           uint64               // keeps track of remaining free spins, so actions can make decisions based on this.
	spinSeq             uint64               // keeps track of number of spins performed.
	multiplier          float64              // overall multiplier for a round of spins.
	slots               *Slots               // slot machine config.
	symbols             *SymbolSet           // copy of primary symbol set.
	altSymbols          *SymbolSet           // copy of alternate symbol set.
	gridDef             *GridDefinition      // copy of the grid definition.
	paylines            *PaylineSet          // copy of paylines.
	spinner             Spinner              // spinner function.
	refiller            Spinner              // refiller function.
	gamer               interfaces.Gamer     // interface for a game round.
	prng                interfaces.Generator // local PRNG for all randomness.
	reels               Reels                // primary reels.
	altReels            Reels                // the alternate reels.
	locked              []bool               // indicators for locked reels.
	hot                 []bool               // indicators for hot reels.
	sticky              []bool               // marks sticky symbols in the grid.
	superShape          []bool               // marks super shape in the grid.
	roundFlags          []int                // local flags carried across multiple spins of a round.
	multipliers         []uint16             // optional multipliers for the symbols on the grid.
	mask                utils.UInt8s         // copy of slot machine reels mask.
	payouts             utils.UInt8s         // marks payout symbols in the grid; 1=standard; 2+=wild.
	effects             utils.UInt8s         // special effects indicators.
	jumps               utils.UInt8s         // array of jumping symbol vectors: 0=no jump, 255=off grid, otherwise new offset+1.
	indexes             utils.Indexes        // the symbol grid.
	injections          utils.Indexes        // array of symbol injections.
	pool.Object
}

// spinsProducer is the memory pool for spin rounds.
// Make sure to initialize all slices appropriately!
var spinsProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &Spin{
		reels:       make(Reels, 0, 8),
		altReels:    make(Reels, 0, 8),
		locked:      make([]bool, 0, 8),
		hot:         make([]bool, 0, 8),
		sticky:      make([]bool, 0, 24),
		superShape:  make([]bool, 0, 24),
		roundFlags:  make([]int, 16),
		multipliers: make([]uint16, 0, 24),
		payouts:     make(utils.UInt8s, 0, 16),
		effects:     make(utils.UInt8s, 0, 24),
		jumps:       make(utils.UInt8s, 0, 24),
		indexes:     make(utils.Indexes, 0, 24),
		injections:  make(utils.Indexes, 0, 24),
	}
	return s, s.reset
})

// reset clears the spin rounds.
func (s *Spin) reset() {
	if s != nil {
		s.debug = false
		s.altActive = false
		s.expanded = false
		s.multiplierNeedsWild = false

		s.kind = 0
		s.newWilds = 0
		s.newHeroes = 0
		s.newScatters = 0
		s.bonusSymbol = utils.NullIndex
		s.superSymbol = utils.NullIndex
		s.stickySymbol = utils.NullIndex
		s.reelCount = 0
		s.rowCount = 0
		s.progressLevel = 0
		s.freeSpins = 0
		s.spinSeq = 0
		s.multiplier = 0.0

		s.slots = nil
		s.symbols = nil
		s.altSymbols = nil
		s.gridDef = nil
		s.paylines = nil
		s.spinner = nil
		s.refiller = nil
		s.gamer = nil
		s.prng = nil
		s.mask = nil

		s.reels = ReleaseReels(s.reels)
		s.altReels = ReleaseReels(s.altReels)

		s.locked = s.locked[:0]
		s.hot = s.hot[:0]
		s.payouts = s.payouts[:0]
		s.sticky = s.sticky[:0]
		s.superShape = s.superShape[:0]
		s.effects = s.effects[:0]
		s.indexes = s.indexes[:0]
		s.injections = s.injections[:0]
		s.jumps = s.jumps[:0]

		if len(s.multipliers) > 0 {
			clear(s.multipliers)
			s.multipliers = s.multipliers[:0]
		}
	}
}
