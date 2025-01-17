package slots

import (
	"fmt"
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireSpinResult instantiates a new spin result from the memory pool using the given spin.
func AcquireSpinResult(spin *Spin) *SpinResult {
	r := spinResultsProducer.Acquire().(*SpinResult)
	if spin == nil {
		return r
	}

	r.debug = spin.debug
	r.kind = spin.kind
	r.initial = utils.CopyIndexes(spin.indexes, r.initial)
	r.multipliers = spin.CloneMultipliers(r.multipliers)
	r.lockedReels = spin.Locked(r.lockedReels)
	r.hotReels = spin.Hot(r.hotReels)
	r.sticky = spin.Sticky(r.sticky)
	r.effects = spin.Effects(r.effects)
	r.bonusSymbol = spin.bonusSymbol
	r.stickySymbol = spin.stickySymbol
	r.superSymbol = spin.superSymbol
	r.progressLevel = spin.progressLevel
	r.multiplier = spin.multiplier

	if spin.slots != nil {
		r.configHash = spin.slots.configHash
	}
	return r
}

// AcquireSpinResultFromData instantiates a new spin result from the memory pool using the given data.
// This is for testing/debug purposes only.
// The result takes ownership of all the given slices.
func AcquireSpinResultFromData(initial, afterExpand, afterClear utils.Indexes, locked, hot utils.UInt8s, bonus, sticky, super utils.Index) *SpinResult {
	r := spinResultsProducer.Acquire().(*SpinResult)
	r.kind = RegularSpin
	r.initial = initial
	r.afterExpand = afterExpand
	r.afterClear = afterClear
	r.lockedReels = locked
	r.hotReels = hot
	r.bonusSymbol = bonus
	r.stickySymbol = sticky
	r.superSymbol = super
	return r
}

// Kind returns the type of the result.
func (r *SpinResult) Kind() SpinKind {
	return r.kind
}

// Initial returns the symbol id grid initially selected.
func (r *SpinResult) Initial() utils.Indexes {
	return r.initial
}

// HaveMultipliers returns true if there is at least one grid multiplier.
func (r *SpinResult) HaveMultipliers() bool {
	for ix := range r.multipliers {
		if utils.ValidMultiplier(float64(r.multipliers[ix])) {
			return true
		}
	}
	return false
}

// Multipliers returns the grid multipliers.
func (r *SpinResult) Multipliers() []uint16 {
	return r.multipliers
}

// AfterExpand returns the symbol id grid after expansions.
func (r *SpinResult) AfterExpand() utils.Indexes {
	return r.afterExpand
}

// AfterClear returns the symbol id grid after clearings.
func (r *SpinResult) AfterClear() utils.Indexes {
	return r.afterClear
}

// AfterCascade returns the symbol id grid after the cascading reels feature.
func (r *SpinResult) AfterCascade() utils.Indexes {
	return r.afterCascade
}

// AfterNudge returns the symbol id grid after nudges.
func (r *SpinResult) AfterNudge() utils.Indexes {
	return r.afterNudge
}

// IsLocked returns whether the given reel index is for a locked reel.
// Reel indexes are 1-based, so the first reel has the index 1.
func (r *SpinResult) IsLocked(reel uint8) bool {
	for _, id := range r.lockedReels {
		if id == reel {
			return true
		}
	}
	return false
}

// Locked returns the locked reels.
func (r *SpinResult) Locked() utils.UInt8s {
	return r.lockedReels
}

// IsHot returns whether the given reel index is for a hot reel.
// Reel indexes are 1-based, so the first reel has the index 1.
func (r *SpinResult) IsHot(reel uint8) bool {
	for _, id := range r.hotReels {
		if id == reel {
			return true
		}
	}
	return false
}

// Hot returns the hot reels.
func (r *SpinResult) Hot() utils.UInt8s {
	return r.hotReels
}

// HaveSticky returns true if there are sticky indicators.
func (r *SpinResult) HaveSticky() bool {
	for ix := range r.sticky {
		if r.sticky[ix] != 0 {
			return true
		}
	}
	return false
}

// Sticky returns the sticky indicators.
func (r *SpinResult) Sticky() utils.UInt8s {
	return r.sticky
}

// HaveClearing returns true if one or more tiles were cleared.
func (r *SpinResult) HaveClearing() bool {
	for ix := range r.afterClear {
		if r.afterClear[ix] != r.initial[ix] {
			return true
		}
	}
	return false
}

// Effects returns the special effect indicators.
func (r *SpinResult) Effects() utils.UInt8s {
	return r.effects
}

// BonusSymbol returns the selected bonus symbol.
func (r *SpinResult) BonusSymbol() utils.Index {
	return r.bonusSymbol
}

// StickySymbol returns the selected sticky symbol.
func (r *SpinResult) StickySymbol() utils.Index {
	return r.stickySymbol
}

// SuperSymbol returns the selected super symbol.
func (r *SpinResult) SuperSymbol() utils.Index {
	return r.superSymbol
}

// ProgressLevel returns the current multiplier scale mark.
func (r *SpinResult) ProgressLevel() int {
	return r.progressLevel
}

// Multiplier returns the current multiplier.
func (r *SpinResult) Multiplier() float64 {
	return r.multiplier
}

// IsDebug returns if the spin result was created in debug mode.
func (r *SpinResult) IsDebug() bool {
	return r.debug
}

// SetDebug sets the debug mode for the result.
func (r *SpinResult) SetDebug(debug bool) {
	r.debug = debug
}

// SetMaxPayout sets the spin to have reached max payout.
func (r *SpinResult) SetMaxPayout() {
	r.maxPayout = true
}

// BonusBuy returns if the spin result is part of a bonus buy round and the bonus buy factor if bonus buy is true.
func (r *SpinResult) BonusBuy() (uint8, float64) {
	return r.bonusBuy, r.buyFactor
}

// SetBonusBuy sets the bonus buy mode for the result.
func (r *SpinResult) SetBonusBuy(bb uint8, paid *PaidAction) {
	r.bonusBuy = bb
	r.buyFactor = float64(paid.betMultiplier)
}

// Update refreshes the result with any changes made to the grid after the result was first instantiated.
func (r *SpinResult) Update(spin *Spin) {
	r.initial = utils.CopyIndexes(spin.indexes, r.initial)
	r.multipliers = spin.CloneMultipliers(r.multipliers)
	r.lockedReels = spin.Locked(r.lockedReels)
	r.hotReels = spin.Hot(r.hotReels)
	r.sticky = spin.Sticky(r.sticky)
	r.bonusSymbol = spin.bonusSymbol
	r.stickySymbol = spin.stickySymbol
	r.superSymbol = spin.superSymbol
	r.progressLevel = spin.progressLevel
	r.multiplier = spin.multiplier
}

// SetMultipliers refreshes the grid multipliers.
func (r *SpinResult) SetMultipliers(spin *Spin) {
	r.multipliers = spin.CloneMultipliers(r.multipliers)
}

// SetAfterNudge adds the grid after a nudge operation on the spin result.
func (r *SpinResult) SetAfterNudge(after *Spin) {
	r.afterNudge = utils.CopyIndexes(after.indexes, r.afterNudge)
}

// SetAfterExpand adds any wild expansions to the spin result along with locked reels.
func (r *SpinResult) SetAfterExpand(after *Spin) {
	r.afterExpand = utils.CopyIndexes(after.indexes, r.afterExpand)
	r.lockedReels = after.Locked(r.lockedReels)
}

// SetAfterInject adds the grid after a symbol/cluster injection on the spin result.
func (r *SpinResult) SetAfterInject(after *Spin) {
	r.injections = ReleaseTiles(r.injections)

	maxM := len(after.multipliers)
	for offset := range after.injections {
		if id := after.injections[offset]; id > 0 {
			var sticky uint8
			if after.sticky[offset] {
				sticky = 1
			}
			var multiplier uint16
			if offset < maxM {
				multiplier = after.multipliers[offset]
			}
			r.injections = append(r.injections, AcquireTile(uint8(offset), id, sticky, multiplier))
		}
	}
}

// HasInjections indicates if the result has injected symbols.
func (r *SpinResult) HasInjections() bool {
	return len(r.injections) > 0
}

// SetAfterClear adds the grid after a clear operation on the spin result.
// A clear operation happens for cascading reels on winning paylines or when bomb symbols explode.
func (r *SpinResult) SetAfterClear(after *Spin) {
	r.afterClear = utils.CopyIndexes(after.indexes, r.afterClear)
}

// SetAfterJump adds the grid after a jump operation on the spin result.
func (r *SpinResult) SetAfterJump(after *Spin) {
	r.jumps = ReleaseTiles(r.jumps)

	haveMultipliers := len(after.multipliers) == len(after.indexes)

	for from := range after.jumps {
		if to := after.jumps[from]; to > 0 {
			if to < 255 {
				to--
			}

			var sticky uint8
			if after.sticky[to] {
				sticky = 1
			}

			var multiplier uint16
			if haveMultipliers {
				multiplier = after.multipliers[to]
			}

			tile := AcquireJumpedTile(uint8(from), to, after.indexes[to], sticky, multiplier)
			r.jumps = append(r.jumps, tile)
		}
	}
}

// HasJumps indicates if the result has jumping symbols.
func (r *SpinResult) HasJumps() bool {
	return len(r.jumps) > 0
}

// SetAfterCascade adds the grid after cascading reels operation on the spin result.
// A clear operation happens for cascading reels on winning paylines or when bomb symbols explode.
func (r *SpinResult) SetAfterCascade(after *Spin) {
	r.afterCascade = utils.CopyIndexes(after.indexes, r.afterCascade)
}

// SetSticky sets the sticky symbol indicators from the given spin, and optionally the sticky symbol or super symbol.
func (r *SpinResult) SetSticky(spin *Spin) {
	r.sticky = spin.Sticky(r.sticky)
	if spin.stickySymbol != utils.MaxIndex {
		r.stickySymbol = spin.stickySymbol
	}
	if spin.superSymbol != utils.MaxIndex {
		r.superSymbol = spin.superSymbol
	}
}

// SetEffects sets the special effect indicators from the given spin.
func (r *SpinResult) SetEffects(spin *Spin) {
	r.effects = spin.Effects(r.effects)
}

// SetChosenSticky sets the sticky symbol to the players choice.
func (r *SpinResult) SetChosenSticky(symbol utils.Index) {
	r.stickySymbol = symbol

	l := len(r.initial)
	max := object.NormalizeSize(l, 16)
	r.sticky = utils.PurgeUInt8s(r.sticky, max)[:l]

	for ix := range r.initial {
		if r.initial[ix] == symbol {
			r.sticky[ix] = 1
		} else {
			r.sticky[ix] = 0
		}
	}
}

// SetHot sets the hot reels from the given spin.
func (r *SpinResult) SetHot(spin *Spin) {
	r.hotReels = spin.Hot(r.hotReels)
}

// SetBonusSymbol adds the chosen bonus symbol to the spin result.
func (r *SpinResult) SetBonusSymbol(bonus utils.Index) {
	r.bonusSymbol = bonus
}

// SetMultiplier updates the progress level and overall multiplier.
func (r *SpinResult) SetMultiplier(spin *Spin) {
	r.progressLevel = spin.progressLevel
	r.multiplier = spin.multiplier
}

// SetStickyChoices records the symbols that can be made sticky for a player to choose from.
func (r *SpinResult) SetStickyChoices(spin *Spin) {
	r.stickyChoices = PurgeStickyChoices(r.stickyChoices, 16)
	for id := spin.GetSymbols().maxID; id > 0; id-- {
		if count := spin.CountSymbol(id); count > 0 {
			r.stickyChoices = append(r.stickyChoices, AcquireStickyChoice(id, spin))
		}
	}
}

// hasSticky returns whether the result has one or more sticky flags turned on.
func (r *SpinResult) hasSticky() bool {
	for ix := range r.sticky {
		if r.sticky[ix] != 0 {
			return true
		}
	}
	return false
}

// SetRoundFlags sets the round flags from the given spin.
func (r *SpinResult) SetRoundFlags(spin *Spin) {
	l := len(spin.roundFlags)
	r.roundFlags = utils.PurgeInts(r.roundFlags, l)[:l]
	copy(r.roundFlags, spin.roundFlags)
}

// RoundFlags retrieves the round flags.
func (r *SpinResult) RoundFlags() []int {
	return r.roundFlags
}

// SetExportFlags sets the export flags.
func (r *SpinResult) SetExportFlags(spin *Spin) {
	r.exportFlags = utils.PurgeInts(r.exportFlags, 16)
	max := len(spin.roundFlags)
	for ix := range spin.slots.roundFlags {
		f := spin.slots.roundFlags[ix]
		if f.export {
			if f.id < max {
				r.exportFlags = append(r.exportFlags, spin.roundFlags[f.id])
			} else {
				r.exportFlags = append(r.exportFlags, 0)
			}
		}
	}
}

// ExportFlags retrieves the export flags.
func (r *SpinResult) ExportFlags() []int {
	return r.exportFlags
}

// SetTransition sets the bonus feature transition type.
func (r *SpinResult) SetTransition(transition FeatureTransitionKind) {
	r.transition = transition
}

// Transition returns the bonus feature transition type.
func (r *SpinResult) Transition() FeatureTransitionKind {
	return r.transition
}

// IsEmpty implements the zjson encoder interface.
func (r *SpinResult) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson encoder interface.
func (r *SpinResult) EncodeFields(enc *zjson.Encoder) {
	r.encode(enc, true)
}

// Encode2 implements the PoolRCZ.Encode2 interface.
func (r *SpinResult) Encode2(enc *zjson.Encoder) {
	r.encode(enc, false)
}

func (r *SpinResult) encode(enc *zjson.Encoder, withLog bool) {
	enc.Uint8Field("kind", uint8(r.kind))
	enc.IntBoolFieldOpt("debug", r.debug)
	enc.IntBoolFieldOpt("maxPayout", r.maxPayout)
	enc.Uint8FieldOpt("bonusBuy", r.bonusBuy)
	enc.Uint8FieldOpt("transition", uint8(r.transition))

	enc.StartArrayField("initial")
	for ix := range r.initial {
		enc.Uint64(uint64(r.initial[ix]))
	}
	enc.EndArray()

	if len(r.multipliers) > 0 {
		enc.StartArrayField("multipliers")
		for ix := range r.multipliers {
			enc.Uint64(uint64(r.multipliers[ix]))
		}
		enc.EndArray()
	}

	if len(r.afterNudge) > 0 {
		enc.StartArrayField("afterNudge")
		for ix := range r.afterNudge {
			enc.Uint64(uint64(r.afterNudge[ix]))
		}
		enc.EndArray()
		enc.StartArrayField("flagsNudge")
		for ix := range r.afterNudge {
			enc.IntBool(r.afterNudge[ix] != r.initial[ix])
		}
		enc.EndArray()
	}

	if len(r.afterExpand) > 0 {
		enc.StartArrayField("afterExpand")
		for ix := range r.afterExpand {
			enc.Uint64(uint64(r.afterExpand[ix]))
		}
		enc.EndArray()
		enc.StartArrayField("flagsExpand")
		for ix := range r.afterExpand {
			enc.IntBool(r.afterExpand[ix] != r.initial[ix])
		}
		enc.EndArray()
	}

	if len(r.injections) > 0 {
		enc.StartArrayField("injections")
		for ix := range r.injections {
			enc.Object(r.injections[ix])
		}
		enc.EndArray()
	}

	if len(r.afterClear) > 0 {
		enc.StartArrayField("afterClear")
		for ix := range r.afterClear {
			enc.Uint64(uint64(r.afterClear[ix]))
		}
		enc.EndArray()
		enc.StartArrayField("flagsClear")
		for ix := range r.afterClear {
			enc.IntBool(r.afterClear[ix] != r.initial[ix])
		}
		enc.EndArray()
	}

	if len(r.jumps) > 0 {
		enc.StartArrayField("jumps")
		for ix := range r.jumps {
			enc.Object(r.jumps[ix])
		}
		enc.EndArray()
	}

	if len(r.afterCascade) > 0 {
		enc.StartArrayField("afterCascade")
		for ix := range r.afterCascade {
			enc.Uint64(uint64(r.afterCascade[ix]))
		}
		enc.EndArray()
		enc.StartArrayField("flagsCascade")
		for ix := range r.afterCascade {
			enc.IntBool(r.afterCascade[ix] != r.initial[ix])
		}
		enc.EndArray()
	}

	if len(r.sticky) > 0 {
		enc.StartArrayField("sticky")
		for ix := range r.sticky {
			enc.Uint64(uint64(r.sticky[ix]))
		}
		enc.EndArray()
	}

	if len(r.effects) > 0 {
		enc.StartArrayField("effects")
		for ix := range r.effects {
			enc.Uint64(uint64(r.effects[ix]))
		}
		enc.EndArray()
	}

	if len(r.lockedReels) > 0 {
		enc.StartArrayField("lockedReels")
		for ix := range r.lockedReels {
			enc.Uint64(uint64(r.lockedReels[ix]))
		}
		enc.EndArray()
	}

	if len(r.hotReels) > 0 {
		enc.StartArrayField("hotReels")
		for ix := range r.hotReels {
			enc.Uint64(uint64(r.hotReels[ix]))
		}
		enc.EndArray()
	}

	if r.bonusSymbol != utils.MaxIndex {
		enc.Uint16FieldOpt("bonusSymbol", uint16(r.bonusSymbol))
	}
	if r.stickySymbol != utils.MaxIndex {
		enc.Uint16FieldOpt("stickySymbol", uint16(r.stickySymbol))
	}
	if r.superSymbol != utils.MaxIndex {
		enc.Uint16FieldOpt("superSymbol", uint16(r.superSymbol))
	}

	enc.IntFieldOpt("progressLevel", r.progressLevel)
	enc.FloatFieldOpt("multiplier", r.multiplier, 'g', -1)
	enc.FloatFieldOpt("buyFactor", r.buyFactor, 'g', -1)

	if len(r.nudges) > 0 {
		enc.StartArrayField("reelNudges")
		for ix := range r.nudges {
			enc.Object(r.nudges[ix])
		}
		enc.EndArray()
	}

	if len(r.stickyChoices) > 0 {
		enc.StartArrayField("stickyChoices")
		for ix := range r.stickyChoices {
			enc.Object(r.stickyChoices[ix])
		}
		enc.EndArray()
	}

	if len(r.roundFlags) > 0 {
		enc.StartArrayField("roundFlags")
		for ix := range r.roundFlags {
			enc.Int64(int64(r.roundFlags[ix]))
		}
		enc.EndArray()
	}

	if len(r.exportFlags) > 0 {
		enc.StartArrayField("exportFlags")
		for ix := range r.exportFlags {
			enc.Int64(int64(r.exportFlags[ix]))
		}
		enc.EndArray()
	}

	r.PlayerChoices.EncodeChoices(enc)
	if withLog {
		r.PrngLog.EncodeEventLog(enc)
	}

	enc.StringFieldOpt("configHash", r.configHash)
}

// DecodeField implements the zjson decoder interface.
func (r *SpinResult) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i16 uint16

	if string(key) == "debug" {
		r.debug, ok = dec.IntBool()
	} else if string(key) == "maxPayout" {
		r.maxPayout, ok = dec.IntBool()
	} else if string(key) == "bonusBuy" {
		r.bonusBuy, ok = dec.Uint8()
	} else if string(key) == "kind" {
		var i8 uint8
		if i8, ok = dec.Uint8(); ok {
			r.kind = SpinKind(i8)
		}
	} else if string(key) == "transition" {
		var i8 uint8
		if i8, ok = dec.Uint8(); ok {
			r.transition = FeatureTransitionKind(i8)
		}
	} else if string(key) == "bonusSymbol" {
		if i16, ok = dec.Uint16(); ok {
			r.bonusSymbol = utils.Index(i16)
		}
	} else if string(key) == "stickySymbol" {
		if i16, ok = dec.Uint16(); ok {
			r.stickySymbol = utils.Index(i16)
		}
	} else if string(key) == "superSymbol" {
		if i16, ok = dec.Uint16(); ok {
			r.superSymbol = utils.Index(i16)
		}
	} else if string(key) == "progressLevel" {
		var i int
		if i, ok = dec.Int(); ok {
			r.progressLevel = i
		}
	} else if string(key) == "multiplierMark" { // ** DEPRECATED **
		var i int
		if i, ok = dec.Int(); ok {
			r.progressLevel = i
		}
	} else if string(key) == "multiplier" {
		var i float64
		if i, ok = dec.Float(); ok {
			r.multiplier = i
		}
	} else if string(key) == "buyFactor" {
		var i float64
		if i, ok = dec.Float(); ok {
			r.buyFactor = i
		}
	} else if string(key) == "initial" {
		r.initial = utils.PurgeIndexes(r.initial, 24)
		ok = dec.Array(r.decodeInitial)
	} else if string(key) == "multipliers" {
		r.multipliers = utils.PurgeUInt16s(r.multipliers, 24)
		ok = dec.Array(r.decodeMultipliers)
	} else if string(key) == "afterNudge" {
		r.afterNudge = utils.PurgeIndexes(r.afterNudge, 24)
		ok = dec.Array(r.decodeAfterNudge)
	} else if string(key) == "flagsNudge" {
		ok = dec.Array(r.ignoreInt64)
	} else if string(key) == "afterExpand" {
		r.afterExpand = utils.PurgeIndexes(r.afterExpand, 24)
		ok = dec.Array(r.decodeAfterExpand)
	} else if string(key) == "flagsExpand" {
		ok = dec.Array(r.ignoreInt64)
	} else if string(key) == "injections" {
		r.injections = ReleaseTiles(r.injections)
		ok = dec.Array(r.decodeInjection)
	} else if string(key) == "afterClear" {
		r.afterClear = utils.PurgeIndexes(r.afterClear, 24)
		ok = dec.Array(r.decodeAfterClear)
	} else if string(key) == "flagsClear" {
		ok = dec.Array(r.ignoreInt64)
	} else if string(key) == "jumps" {
		r.jumps = ReleaseTiles(r.jumps)
		ok = dec.Array(r.decodeJump)
	} else if string(key) == "afterCascade" {
		r.afterCascade = utils.PurgeIndexes(r.afterCascade, 24)
		ok = dec.Array(r.decodeAfterCascade)
	} else if string(key) == "flagsCascade" {
		ok = dec.Array(r.ignoreInt64)
	} else if string(key) == "lockedReels" {
		r.lockedReels = utils.PurgeUInt8s(r.lockedReels, 8)
		ok = dec.Array(r.decodeLockedReels)
	} else if string(key) == "hotReels" {
		r.hotReels = utils.PurgeUInt8s(r.hotReels, 8)
		ok = dec.Array(r.decodeHotReels)
	} else if string(key) == "sticky" {
		r.sticky = utils.PurgeUInt8s(r.sticky, 24)
		ok = dec.Array(r.decodeSticky)
	} else if string(key) == "effects" {
		r.effects = utils.PurgeUInt8s(r.effects, 24)
		ok = dec.Array(r.decodeEffects)
	} else if string(key) == "stickyChoices" {
		if cap(r.stickyChoices) < 16 {
			r.stickyChoices = make([]*StickyChoice, 0, 16)
		} else {
			r.stickyChoices = r.stickyChoices[:0]
		}
		ok = dec.Array(r.decodeStickyChoices)
	} else if string(key) == "roundFlags" {
		r.roundFlags = utils.PurgeInts(r.roundFlags, 16)
		ok = dec.Array(r.decodeRoundFlag)
	} else if string(key) == "exportFlags" {
		r.exportFlags = utils.PurgeInts(r.exportFlags, 16)
		ok = dec.Array(r.decodeExportFlag)
	} else if string(key) == "reelNudges" {
		ok = dec.Array(r.decodeReelNudges)
	} else if string(key) == "playerChoices" {
		ok = r.PlayerChoices.DecodeChoices(dec)
	} else if string(key) == "scriptID" {
		ok = r.PrngLog.DecodeScriptID(dec)
	} else if string(key) == "events" {
		ok = dec.Array(r.PrngLog.DecodeEventLog)
	} else if string(key) == "rngIn" {
		ok = dec.Array(r.PrngLog.DecodeRngIn)
	} else if string(key) == "rngOut" {
		ok = dec.Array(r.PrngLog.DecodeRngOut)
	} else if string(key) == "configHash" {
		var b []byte
		if b, _, ok = dec.String(); ok {
			r.configHash = string(b)
		}
	} else {
		return fmt.Errorf("SpinResult.DecodeField: invalid field encountered [%s]", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) ignoreInt64(dec *zjson.Decoder) error {
	if _, ok := dec.Int64(); ok {
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeInitial(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		r.initial = append(r.initial, utils.Index(i))
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeMultipliers(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		r.multipliers = append(r.multipliers, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeAfterNudge(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		r.afterNudge = append(r.afterNudge, utils.Index(i))
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeAfterExpand(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		r.afterExpand = append(r.afterExpand, utils.Index(i))
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeInjection(dec *zjson.Decoder) error {
	if tile, ok := AcquireTileFromJSON(dec); ok {
		r.injections = append(r.injections, tile)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeAfterClear(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		r.afterClear = append(r.afterClear, utils.Index(i))
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeJump(dec *zjson.Decoder) error {
	if tile, ok := AcquireTileFromJSON(dec); ok {
		r.jumps = append(r.jumps, tile)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeAfterCascade(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		r.afterCascade = append(r.afterCascade, utils.Index(i))
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeLockedReels(dec *zjson.Decoder) error {
	if i, ok := dec.Uint8(); ok {
		r.lockedReels = append(r.lockedReels, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeHotReels(dec *zjson.Decoder) error {
	if i, ok := dec.Uint8(); ok {
		r.hotReels = append(r.hotReels, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeSticky(dec *zjson.Decoder) error {
	if i, ok := dec.Uint8(); ok {
		r.sticky = append(r.sticky, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeRoundFlag(dec *zjson.Decoder) error {
	if i, ok := dec.Int(); ok {
		r.roundFlags = append(r.roundFlags, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeExportFlag(dec *zjson.Decoder) error {
	if i, ok := dec.Int(); ok {
		r.exportFlags = append(r.exportFlags, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeEffects(dec *zjson.Decoder) error {
	if i, ok := dec.Uint8(); ok {
		r.effects = append(r.effects, i)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeReelNudges(dec *zjson.Decoder) error {
	n := reelNudgeProducer.Acquire().(*ReelNudge)
	if ok := dec.Object(n); ok {
		r.nudges = append(r.nudges, n)
		return nil
	}
	return dec.Error()
}

func (r *SpinResult) decodeStickyChoices(dec *zjson.Decoder) error {
	c := stickyChoiceProducer.Acquire().(*StickyChoice)
	if ok := dec.Object(c); ok {
		r.stickyChoices = append(r.stickyChoices, c)
		return nil
	}
	return dec.Error()
}

// MarshalJSON implements the json marshaller interface.
func (r *SpinResult) MarshalJSON() ([]byte, error) {
	enc := zjson.AcquireEncoder(4096)
	defer enc.Release()

	enc.Object(r)
	b := enc.Bytes()

	out := make([]byte, len(b))
	copy(out, b)

	return out, nil
}

// SpinResult contains the display properties of a spin result.
// SpinResult is not safe for concurrent use across multiple go-routines.
// Keep fields ordered by ascending SizeOf().
type SpinResult struct {
	debug         bool                  // indicates a spin resulting from debug mode.
	maxPayout     bool                  // indicates this spin reached the max payout.
	bonusBuy      uint8                 // indicates a spin in bonus buy/bonus bet mode.
	kind          SpinKind              // type of spin result.
	transition    FeatureTransitionKind // optional type of transition.
	bonusSymbol   utils.Index           // the selected bonus symbol for free spins.
	stickySymbol  utils.Index           // the selected sticky symbol.
	superSymbol   utils.Index           // the selected super symbol for grid super shapes.
	progressLevel int                   // mark on the overall muliplier scale.
	multiplier    float64               // overall multiplier.
	buyFactor     float64               // bet multiplier for bonus buy/bonus bet.
	initial       utils.Indexes         // initial symbol grid.
	multipliers   []uint16              // optional multipliers for symbol grid.
	afterNudge    utils.Indexes         // symbol grid after nudge operations.
	afterExpand   utils.Indexes         // symbol grid after expansions.
	afterClear    utils.Indexes         // symbol grid after clear operations.
	afterCascade  utils.Indexes         // symbol grid after cascading reels operation.
	lockedReels   utils.UInt8s          // locked reels after wild expansions (1-based).
	hotReels      utils.UInt8s          // hot reels during free games (1-based).
	sticky        utils.UInt8s          // sticky symbol indicators (same size as grid).
	effects       utils.UInt8s          // special effect indicators (same size as grid).
	injections    Tiles                 // list of injected symbols.
	jumps         Tiles                 // list of jumped symbols.
	nudges        ReelNudges            // list of nudged reels.
	stickyChoices StickyChoices         // list of possible sticky choices for player to choose from.
	roundFlags    []int                 // optional round flags with their latest value.
	exportFlags   []int                 // all exported round flags with their latest value.
	configHash    string                // sha256 hash of the game config.
	results.PlayerChoices
	results.PrngLog
	pool.Object
}

// spinResultsProducer is the memory pool for spin results.
// Make sure to initialize all slices appropriately!
var spinResultsProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &SpinResult{
		initial:       make(utils.Indexes, 0, 24),
		multipliers:   make([]uint16, 0, 24),
		afterNudge:    make(utils.Indexes, 0, 24),
		afterExpand:   make(utils.Indexes, 0, 24),
		afterClear:    make(utils.Indexes, 0, 24),
		afterCascade:  make(utils.Indexes, 0, 24),
		lockedReels:   make(utils.UInt8s, 0, 8),
		hotReels:      make(utils.UInt8s, 0, 8),
		sticky:        make(utils.UInt8s, 0, 24),
		effects:       make(utils.UInt8s, 0, 24),
		injections:    make(Tiles, 0, 16),
		jumps:         make(Tiles, 0, 4),
		nudges:        make(ReelNudges, 0, 2),
		stickyChoices: make(StickyChoices, 0, 2),
		roundFlags:    make([]int, 0, 16),
		exportFlags:   make([]int, 0, 16),
	}
	r.PrngLog.Initialize()
	return r, r.reset
})

// reset clears the spin results.
func (r *SpinResult) reset() {
	if r != nil {
		r.debug = false
		r.maxPayout = false
		r.bonusBuy = 0
		r.kind = 0
		r.transition = 0
		r.bonusSymbol = utils.NullIndex
		r.stickySymbol = utils.NullIndex
		r.superSymbol = utils.NullIndex
		r.progressLevel = 0
		r.multiplier = 0.0
		r.buyFactor = 0.0

		r.initial = r.initial[:0]
		r.multipliers = r.multipliers[:0]
		r.afterNudge = r.afterNudge[:0]
		r.afterExpand = r.afterExpand[:0]
		r.afterClear = r.afterClear[:0]
		r.afterCascade = r.afterCascade[:0]
		r.lockedReels = r.lockedReels[:0]
		r.hotReels = r.hotReels[:0]
		r.sticky = r.sticky[:0]
		r.effects = r.effects[:0]

		r.injections = ReleaseTiles(r.injections)
		r.jumps = ReleaseTiles(r.jumps)
		r.nudges = ReleaseReelNudges(r.nudges)
		r.stickyChoices = ReleaseStickyChoices(r.stickyChoices)
		r.roundFlags = utils.PurgeInts(r.roundFlags, 16)
		r.exportFlags = utils.PurgeInts(r.exportFlags, 16)

		r.configHash = ""

		r.PlayerChoices.Reset()
		r.PrngLog.Reset()
	}
}

// Equals is used internally for unit-tests!
func (r *SpinResult) Equals(other *SpinResult) bool {
	if r.debug != other.debug ||
		r.maxPayout != other.maxPayout ||
		r.kind != other.kind ||
		r.bonusBuy != other.bonusBuy ||
		r.transition != other.transition ||
		r.bonusSymbol != other.bonusSymbol ||
		r.stickySymbol != other.stickySymbol ||
		r.superSymbol != other.superSymbol ||
		r.progressLevel != other.progressLevel ||
		r.multiplier != other.multiplier ||
		r.buyFactor != other.buyFactor ||
		r.configHash != other.configHash ||
		len(r.stickyChoices) != len(other.stickyChoices) ||
		!r.PrngLog.Equals(&other.PrngLog) ||
		!r.PlayerChoices.Equal(&other.PlayerChoices) ||
		!r.injections.DeepEqual(other.injections) ||
		!r.jumps.DeepEqual(other.jumps) ||
		!r.nudges.DeepEqual(other.nudges) ||
		!reflect.DeepEqual(r.initial, other.initial) ||
		!reflect.DeepEqual(r.multipliers, other.multipliers) ||
		!reflect.DeepEqual(r.afterNudge, other.afterNudge) ||
		!reflect.DeepEqual(r.afterExpand, other.afterExpand) ||
		!reflect.DeepEqual(r.afterClear, other.afterClear) ||
		!reflect.DeepEqual(r.afterCascade, other.afterCascade) ||
		!reflect.DeepEqual(r.lockedReels, other.lockedReels) ||
		!reflect.DeepEqual(r.hotReels, other.hotReels) ||
		!reflect.DeepEqual(r.sticky, other.sticky) ||
		!reflect.DeepEqual(r.effects, other.effects) ||
		!reflect.DeepEqual(r.roundFlags, other.roundFlags) ||
		!reflect.DeepEqual(r.exportFlags, other.exportFlags) {
		return false
	}

	for ix := range r.nudges {
		if !r.nudges[ix].DeepEqual(other.nudges[ix]) {
			return false
		}
	}
	for ix := range r.stickyChoices {
		if !r.stickyChoices[ix].Equals(other.stickyChoices[ix]) {
			return false
		}
	}

	return true
}
