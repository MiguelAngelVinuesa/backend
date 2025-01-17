package slots

import (
	"crypto/sha256"
	"encoding/base64"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	defaultHighestPayout         = false
	defaultCascadingReels        = false
	defaultHotReelsAsBonusSymbol = false
	defaultMultiplierOnWildsOnly = false
	defaultDirection             = PayLTR
	defaultReelCount             = 5
	defaultRowCount              = 3
	defaultNoRepeat              = 0
)

// Slots represents a slot machine.
type Slots struct {
	highestPayout         bool
	cascadingReels        bool
	doubleSpin            bool
	playerChoice          bool
	roundMultiplier       bool
	hotReelsAsBonusSymbol bool
	multiplierOnWildsOnly bool
	progressMeter         bool
	haveExportedFlags     bool
	bonusBuy              bool
	symbolsState          bool
	clusterPays           bool
	reverseWin            bool
	directions            PayDirection
	noRepeat              uint8
	reelCount             int
	rowCount              int
	flagBB                int
	maxPayout             float64
	targetRTP             float64
	symbols               *SymbolSet
	altSymbols            *SymbolSet
	paylines              *PaylineSet
	gridDef               *GridDefinition
	scriptedRoundSelector *ScriptedRoundSelector
	spinner               Spinner
	refiller              Spinner
	configHash            string
	mask                  utils.UInt8s
	actionsFirst          SpinActions
	actionsFree           SpinActions
	actionsFirstBB        SpinActions
	actionsFreeBB         SpinActions
	excludeFromState      utils.Indexes
	roundFlags            RoundFlags
}

// NewSlots instantiates a slot machine.
// Internally it creates a slice of reels with a length equal to the given reelCount.
// The default reel x row grid is 3x3.
// The function panics if the reel or row count is zero or unreasonable high (> 15).
// Slots are immutable once created so are safe to be used across concurrent go-routines.
func NewSlots(opts ...SlotOption) *Slots {
	s := &Slots{
		highestPayout:         defaultHighestPayout,
		cascadingReels:        defaultCascadingReels,
		hotReelsAsBonusSymbol: defaultHotReelsAsBonusSymbol,
		multiplierOnWildsOnly: defaultMultiplierOnWildsOnly,
		directions:            defaultDirection,
		reelCount:             defaultReelCount,
		rowCount:              defaultRowCount,
		noRepeat:              defaultNoRepeat,
	}

	for ix := range opts {
		opts[ix](s)
	}

	s.gridDef = NewGridDefinition(s.reelCount, s.rowCount, s.mask)

	for ix := range s.roundFlags {
		if s.roundFlags[ix].export {
			s.haveExportedFlags = true
			break
		}
	}

	enc := zjson.AcquireEncoder(4096)
	enc.Object(s)
	h := sha256.New()
	h.Write(enc.Bytes())
	enc.Release()
	s.configHash = base64.StdEncoding.EncodeToString(h.Sum(nil))

	return s
}

// ConfigHash returns the config hash.
func (s *Slots) ConfigHash() string {
	return s.configHash
}

// ReelCount returns the number of reels.
func (s *Slots) ReelCount() int {
	return s.reelCount
}

// RowCount returns the number of rows.
func (s *Slots) RowCount() int {
	return s.rowCount
}

// ReelMask returns the reel mask for non-rectangular grids.
func (s *Slots) ReelMask() utils.UInt8s {
	return s.mask
}

// GridDefinition returns the definition for the grid.
func (s *Slots) GridDefinition() *GridDefinition {
	return s.gridDef
}

// Symbols returns the set of symbols.
func (s *Slots) Symbols() *SymbolSet {
	return s.symbols
}

// AltSymbols returns the alternate set of symbols.
func (s *Slots) AltSymbols() *SymbolSet {
	return s.altSymbols
}

// Paylines returns the configured paylines for the slot machine.
func (s *Slots) Paylines() *PaylineSet {
	return s.paylines
}

// CascadingReels returns whether the cascading reels feature is enabled.
func (s *Slots) CascadingReels() bool {
	return s.cascadingReels
}

// ClusterPays returns whether the cluster payouts feature is enabled.
func (s *Slots) ClusterPays() bool {
	return s.clusterPays
}

// DoubleSpin returns whether the double spin feature is enabled.
func (s *Slots) DoubleSpin() bool {
	return s.doubleSpin
}

// PlayerChoice returns whether the player can make choices.
func (s *Slots) PlayerChoice() bool {
	return s.playerChoice
}

// RoundMultiplier returns whether spin round may accumulate a multiplier.
func (s *Slots) RoundMultiplier() bool {
	return s.roundMultiplier
}

// ProgressMeter returns whether the spin round may display a progress meter.
func (s *Slots) ProgressMeter() bool {
	return s.progressMeter
}

// MaxPayout returns the maximum allowed payout.
// A value of 0 means there is no maximum.
func (s *Slots) MaxPayout() float64 {
	return s.maxPayout
}

// BonusBuy returns whether the game has a bonus buy feature.
func (s *Slots) BonusBuy() bool {
	return s.bonusBuy
}

// SymbolsState returns whether the game has symbols state.
func (s *Slots) SymbolsState() bool {
	return s.symbolsState
}

// ReverseWin returns whether the game has the "reverse win" style.
func (s *Slots) ReverseWin() bool {
	return s.reverseWin
}

// ExcludeFromState returns the symbols ids that must be excluded from the symbols state.
func (s *Slots) ExcludeFromState() utils.Indexes {
	return s.excludeFromState
}

// RTP returns the official target and 2 simulator RTP for the game.
// Depending on the game the 2nd simulator RTP may be zero.
func (s *Slots) RTP() float64 {
	return s.targetRTP
}

// RoundFlags returns the descriptions for the flags used during a spin round.
func (s *Slots) RoundFlags() RoundFlags {
	return s.roundFlags
}

// HaveExportedFlags returns true if one or more flags are marked to be exported.
func (s *Slots) HaveExportedFlags() bool {
	return s.haveExportedFlags
}

// ActionsFirst returns the actions for the slots game for the first spin with no bonus buy.
func (s *Slots) ActionsFirst() SpinActions {
	return s.actionsFirst
}

// ActionsFree returns the actions for the slots game for the free spins with no bonus buy.
func (s *Slots) ActionsFree() SpinActions {
	return s.actionsFree
}

// ActionsFirstBB returns the actions for the slots game for the first spin with a bonus buy or bonus bet active.
func (s *Slots) ActionsFirstBB() SpinActions {
	return s.actionsFirstBB
}

// ActionsFreeBB returns the actions for the slots game for the free spins with a bonus buy or bonus bet active.
func (s *Slots) ActionsFreeBB() SpinActions {
	return s.actionsFreeBB
}

// ScriptedRoundSelector returns the scripted round selector.
func (s *Slots) ScriptedRoundSelector() *ScriptedRoundSelector {
	return s.scriptedRoundSelector
}

// IsEmpty implements the zjson.Encoder interface.
func (s *Slots) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson.Encoder interface.
func (s *Slots) EncodeFields(enc *zjson.Encoder) {
	enc.ObjectField("grid", s.gridDef)
	enc.Uint8Field("noRepeat", s.noRepeat)
	enc.StringField("paylineDirection", s.directions.String())
	enc.FloatField("maxPayout", s.maxPayout, 'f', 2)
	enc.FloatField("targetRTP", s.targetRTP, 'f', 2)
	enc.IntBoolFieldOpt("highestPayout", s.highestPayout)
	enc.IntBoolFieldOpt("cascadingReels", s.cascadingReels)
	enc.IntBoolFieldOpt("doubleSpin", s.doubleSpin)
	enc.IntBoolFieldOpt("playerChoice", s.playerChoice)
	enc.IntBoolFieldOpt("hotReelsAsBonusSymbol", s.hotReelsAsBonusSymbol)
	enc.IntBoolFieldOpt("roundMultiplier", s.roundMultiplier)
	enc.IntBoolFieldOpt("multiplierOnWildsOnly", s.multiplierOnWildsOnly)
	enc.IntBoolFieldOpt("progressMeter", s.progressMeter)
	enc.IntBoolFieldOpt("clusterPays", s.clusterPays)
	enc.IntBoolFieldOpt("reverseWin", s.reverseWin)

	enc.StartArrayField("symbols")
	for ix := range s.symbols.symbols {
		enc.Object(s.symbols.symbols[ix])
	}
	enc.EndArray()

	if s.altSymbols != nil {
		enc.StartArrayField("altSymbols")
		for ix := range s.altSymbols.symbols {
			enc.Object(s.altSymbols.symbols[ix])
		}
		enc.EndArray()
	}

	if s.paylines != nil && len(s.paylines.paylines) > 0 {
		enc.StartArrayField("paylines")
		for ix := range s.paylines.paylines {
			enc.Object(s.paylines.paylines[ix])
		}
		enc.EndArray()
	}

	if len(s.roundFlags) > 0 {
		enc.StartArrayField("flags")
		for ix := range s.roundFlags {
			enc.Object(s.roundFlags[ix])
		}
		enc.EndArray()
	}

	if len(s.actionsFirst) > 0 {
		enc.StartArrayField("actionsFirst")
		for ix := range s.actionsFirst {
			enc.String(s.actionsFirst[ix].Config())
		}
		enc.EndArray()
	}

	if len(s.actionsFree) > 0 {
		enc.StartArrayField("actionsFree")
		for ix := range s.actionsFree {
			enc.String(s.actionsFree[ix].Config())
		}
		enc.EndArray()
	}

	if len(s.actionsFirstBB) > 0 {
		enc.StartArrayField("actionsFirstBB")
		for ix := range s.actionsFirstBB {
			enc.String(s.actionsFirstBB[ix].Config())
		}
		enc.EndArray()
	}

	if len(s.actionsFreeBB) > 0 {
		enc.StartArrayField("actionsFreeBB")
		for ix := range s.actionsFreeBB {
			enc.String(s.actionsFreeBB[ix].Config())
		}
		enc.EndArray()
	}
}

// SlotOption is the function signature for slot machine options.
type SlotOption func(s *Slots)

// Grid initializes the reel- and row-count of the slot machine.
func Grid(reelCount, rowCount uint8) SlotOption {
	return func(s *Slots) {
		s.reelCount = int(reelCount)
		s.rowCount = int(rowCount)
	}
}

// WithMask defines the reel mask for changing the slot machine grid to a non-rectangular shape.
// E.g. this can be used to define a grid like 3-4-5-4-3 sized reels, or any other shape.
// The mask defines how many positions for each reel are used. During a spin the unused positions will remain zero.
// NOTE: make sure to define the mask for all reels even if the last ones are full-sized,
// as the spin object relies on the mask slice having a length equal to the number of reels!
// The slot engine will panic if you ignore this rule!
func WithMask(mask ...uint8) SlotOption {
	return func(s *Slots) {
		s.mask = mask
	}
}

// WithSymbols initializes the symbols for the slot machine.
func WithSymbols(symbols *SymbolSet) SlotOption {
	return func(s *Slots) {
		s.symbols = symbols
	}
}

// WithAltSymbols initializes the alternate symbol set for the slot machine.
func WithAltSymbols(altSymbols *SymbolSet) SlotOption {
	return func(s *Slots) {
		s.altSymbols = altSymbols
	}
}

// NoRepeat signals the PRNG to prevent repeating symbols on the reels for the given number of rows.
func NoRepeat(noRepeat uint8) SlotOption {
	return func(s *Slots) {
		s.noRepeat = noRepeat
	}
}

// WithPaylines set the paylines for the slot machine.
func WithPaylines(directions PayDirection, highestPayout bool, paylines ...*Payline) SlotOption {
	return func(s *Slots) {
		s.directions = directions
		s.highestPayout = highestPayout
		s.paylines = NewPaylineSet(directions, highestPayout, paylines...)
	}
}

// PayDirections set the valid direction(s) for the paylines.
func PayDirections(directions PayDirection) SlotOption {
	return func(s *Slots) {
		s.directions = directions
	}
}

// HighestPayout turns the highest payout feature on.
// When turned on the highest paying combo on a payline will be selected.
// E.g. if two wild symbols plus the next low-level symbol warrant a 1x payout, but,
// the two wilds can be substituted for a high-level symbol that warrants a 2x payout,
// the high-level symbol wins and will be the resulting winline with a count of 2.
func HighestPayout() SlotOption {
	return func(s *Slots) {
		s.highestPayout = true
	}
}

// CascadingReels turns the cascading reels feature on.
// This feature is important on both the front-end and back-end.
// For the FE it means the grid is cleared by dropping symbols to the bottom of the screen, and then filled again
// with symbols falling from the top of the screen.
// After a clear operation (winning paylines or bombs), remaining symbols fall down and new symbols fall from
// the top of the screen to fill the empty spaces.
// In the back-end, this feature only plays a role after a clear operation, as remaining symbols on higher rows
// must move to lower rows if the lower row(s) are empty (symbol index == 0) below it.
func CascadingReels(clusterPays bool) SlotOption {
	return func(s *Slots) {
		s.cascadingReels = true
		s.clusterPays = clusterPays
	}
}

func ClusterPays(clusterPays bool) SlotOption {
	return func(s *Slots) {
		s.clusterPays = clusterPays
	}
}

// DoubleSpin indicates that the slot machine implments the double spin feature.
func DoubleSpin() SlotOption {
	return func(s *Slots) {
		s.doubleSpin = true
	}
}

// WithPlayerChoice indicates that the player can make choices during a round.
func WithPlayerChoice() SlotOption {
	return func(s *Slots) {
		s.playerChoice = true
	}
}

// HotReelsAsBonusSymbol indicates that hot reels will automatically count as bonus symbol during free spins.
func HotReelsAsBonusSymbol() SlotOption {
	return func(s *Slots) {
		s.hotReelsAsBonusSymbol = true
	}
}

// WithRoundMultiplier indicates that a multiplier may be accumulated during a game round.
func WithRoundMultiplier() SlotOption {
	return func(s *Slots) {
		s.roundMultiplier = true
	}
}

// WithMultiplierOnWildsOnly indicates the round multiplier only applies if the payline contains a wild.
func WithMultiplierOnWildsOnly() SlotOption {
	return func(s *Slots) {
		s.multiplierOnWildsOnly = true
	}
}

// WithProgressMeter indicates that a progress meter may be visible during a game round.
func WithProgressMeter() SlotOption {
	return func(s *Slots) {
		s.progressMeter = true
	}
}

// WithBonusBuy indicates that the player can initiate a bonus buy.
// The optional flag indicates which round flag indicates the bonus buy was activated for a round.
func WithBonusBuy(flag ...int) SlotOption {
	return func(s *Slots) {
		s.bonusBuy = true
		if len(flag) > 0 {
			s.flagBB = flag[0]
		}
	}
}

// WithSymbolsState indicates that the game should keep symbols state.
func WithSymbolsState(exclude ...utils.Index) SlotOption {
	return func(s *Slots) {
		s.symbolsState = true
		s.excludeFromState = exclude
	}
}

// MaxPayout defines what maximum payout is allowed for a single spin round including awarded free spins/bonus games.
// The value 0 means there is no limit (= very dangerous!).
func MaxPayout(maxPayout float64) SlotOption {
	return func(s *Slots) {
		s.maxPayout = maxPayout
	}
}

// WithRTP adds the official target RTP to the game.
func WithRTP(target float64) SlotOption {
	return func(s *Slots) {
		s.targetRTP = target
	}
}

// WithRoundFlags adds descriptions for the flags used during a spin round.
func WithRoundFlags(flags ...*RoundFlag) SlotOption {
	return func(s *Slots) {
		s.roundFlags = flags
	}
}

// WithActions adds the actions used for the slot machine.
func WithActions(first, free, firstBB, freeBB SpinActions) SlotOption {
	return func(s *Slots) {
		s.actionsFirst = first
		s.actionsFree = free
		s.actionsFirstBB = firstBB
		s.actionsFreeBB = freeBB
	}
}

// WithScriptedRoundSelector adds a scripted round selector to the slot machine.
func WithScriptedRoundSelector(sel *ScriptedRoundSelector) SlotOption {
	return func(s *Slots) {
		s.scriptedRoundSelector = sel
	}
}

// WithReverseWin indicates that the game has a "reverse win" style.
func WithReverseWin() SlotOption {
	return func(s *Slots) {
		s.reverseWin = true
	}
}

// WithSpinner adds a Spinner interface to the slot machine.
func WithSpinner(spinner Spinner) SlotOption {
	return func(s *Slots) {
		s.spinner = spinner
	}
}

// WithRefiller adds a Spinner interface for refilling the grid to the slot machine.
func WithRefiller(refiller Spinner) SlotOption {
	return func(s *Slots) {
		s.refiller = refiller
	}
}
