package slots

import (
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

type (
	// SpinActionStage is the stage at which a spin action is tested.
	SpinActionStage uint8
	// SpinActionResult is the result of a spin action.
	SpinActionResult uint8
)

// always add new elements at the end!!
const (
	// PaidOnly indicates the action can only be "bought".
	PaidOnly SpinActionStage = iota + 1
	// ReviseGrid is the stage where the grid can be revised based on certain conditions.
	ReviseGrid
	// ExpandBefore is the stage right before regular payouts are tested.
	ExpandBefore
	// TestGrid is the stage where the grid is tested for special features.
	TestGrid
	// RegularPayouts is the stage where regular payouts are tested and awarded.
	RegularPayouts
	// ExpandAfter is the stage right after regular payouts are tested.
	ExpandAfter
	// TestState is the stage where the game state is tested & updated.
	TestState
	// ExtraPayouts is the stage where additional payouts are awarded.
	ExtraPayouts
	// AwardBonuses is the stage where free spins and/or bonus games are awarded or the bonus symbol expands and pays out.
	AwardBonuses
	// TestStickiness is the stage where sticky symbols are determined.
	TestStickiness
	// TestClearance is the stage where clearances are tested; e.g. clear payouts, exploding bombs.
	TestClearance
	// PreBonus is the stage after the initial spin(s) and just before bonus games will get played.
	PreBonus
	// TestPlayerChoice is the stage where player choices (if any) are processed.
	TestPlayerChoice
	// PreSpin is the stage before any spin(s) are executed.
	PreSpin
	// Injection are tested right after regular payouts.
	Injection
	// RegularPenalties is the stage where common penalties are imposed.
	RegularPenalties
)

// String implements the Stringer interface.
func (s SpinActionStage) String() string {
	switch s {
	case PaidOnly:
		return "PaidOnly"
	case ReviseGrid:
		return "ReviseGrid"
	case ExpandBefore:
		return "ExpandBefore"
	case TestGrid:
		return "TestGrid"
	case RegularPayouts:
		return "Payouts"
	case ExpandAfter:
		return "ExpandAfter"
	case TestState:
		return "State"
	case ExtraPayouts:
		return "ExtraPayouts"
	case AwardBonuses:
		return "Bonuses"
	case TestStickiness:
		return "Stickiness"
	case TestClearance:
		return "Clearance"
	case PreBonus:
		return "PreBonus"
	case TestPlayerChoice:
		return "PlayerChoice"
	case PreSpin:
		return "PreSpin"
	case Injection:
		return "Injection"
	case RegularPenalties:
		return "Penalties"
	default:
		return "???"
	}
}

// always add new elements at the end!!
const (
	// Processed is the action result requiring no further operations.
	Processed SpinActionResult = iota
	// Payout is the action result awarding a direct payout.
	Payout
	// FreeSpins is the action result awarding one or more free spins.
	FreeSpins
	// BonusGame is the action result awarding a bonus game.
	BonusGame
	// Refill is the action result to request a refill (partial free spin).
	Refill
	// SuperRefill is the action result to request a refill for a super shape (partial free spin).
	SuperRefill
	// HotReel is the action result that indicates at least one new reel has been marked as hot.
	HotReel
	// Sticky is the action result that indicates at one or more tiles were marked as sticky.
	Sticky
	// ChooseSticky is the action result that indicates there are multiple choices for the sticky symbol.
	ChooseSticky
	// Multiplier is the action result that indicates the overall multiplier scale has been changed.
	Multiplier
	// Multipliers is the action result that indicates the grid multipliers have been changed.
	Multipliers
	// InstantBonus is the action result that indicates an instant bonus is awarded before the first spin.
	InstantBonus
	// SpecialResult is the action result that indicates a special result is awarded.
	SpecialResult
	// ReelsNudged is the action result that indicates reels may have been nudged.
	ReelsNudged
	// WildsJumped is the action result that indicates wilds may have jumped.
	WildsJumped
	// SymbolsInjected is the action result that indicates a symbol or cluster may have been injected.
	SymbolsInjected
	// GridModified is the action result to signal that the grid was modified.
	GridModified
	// Penalty is the action result imposing a direct penalty.
	Penalty
)

// String implements the Stringer interface.
func (r SpinActionResult) String() string {
	switch r {
	case Processed:
		return "Processed"
	case Payout:
		return "Payout"
	case FreeSpins:
		return "FreeSpins"
	case BonusGame:
		return "BonusGame"
	case Refill:
		return "Refill"
	case SuperRefill:
		return "SuperRefill"
	case HotReel:
		return "HotReel"
	case Sticky:
		return "Sticky"
	case ChooseSticky:
		return "ChooseSticky"
	case Multiplier:
		return "Multiplier"
	case Multipliers:
		return "Multipliers"
	case InstantBonus:
		return "InstantBonus"
	case SpecialResult:
		return "SpecialResult"
	case ReelsNudged:
		return "ReelsNudged"
	case WildsJumped:
		return "WildsJumped"
	case SymbolsInjected:
		return "SymbolsInjected"
	case Penalty:
		return "Penalty"
	default:
		return "???"
	}
}

// SpinActioner is the interface for spin actions.
type SpinActioner interface {
	ID() int                                              // unique id for the action.
	Name() string                                         // name for the action.
	Kind() string                                         // kind of the action.
	Config() string                                       // config for the action.
	Stage() SpinActionStage                               // stage at which the action should be tested.
	Result() SpinActionResult                             // type of result for the action.
	NrOfSpins(prng interfaces.Generator) uint8            // number of free spins if triggered.
	AltSymbols() bool                                     // indicates to use the alternate symbol set for free spins.
	Alternate() SpinActioner                              // returns an (optional) alternate action to take if the action doesn't trigger.
	BonusSymbol() bool                                    // returns true if a bonus symbol must be selected for awarded free spins.
	PlayerChoice() bool                                   // returns true if a player choice should be presented before continuing the game.
	CanTestChoices(endpoint EndpointKind) bool            // return true if all testChoices filters return true or there are no testChoices filters.
	CanTrigger(*Spin) bool                                // return true if all trigger filters return true or there are no trigger filters.
	CanSticky(*Spin) bool                                 // return true if all sticky filters return true or there are no sticky filters.
	CanClear(*Spin) bool                                  // return true if all clear filters return true or there are no clear filters.
	Triggered(*Spin) SpinActioner                         // indicates the action (or an alternate) triggered for the given spin result.
	TriggeredWithState(*Spin, *SymbolsState) SpinActioner // indicates the action (or an alternate) triggered for the given spin result and symbols state.
	Payout(*Spin, *results.Result) SpinActioner           // add a payout to the result.
	Penalty(*Spin, *results.Result) SpinActioner          // add a penalty to the result.
	Nudge(*Spin, *SpinResult) SpinActioner                // add nudge(s) to the spin result.
	InstantBonus(*Spin) interfaces.Objecter2              // returns an instant bonus result when applicable.
	BonusSelect(*Spin) interfaces.Objecter2               // returns a bonus selector result when applicable.
	BonusGame(*Spin) interfaces.Objecter2                 // plays a bonus game when applicable and returns the result.
	FeatureTransition() FeatureTransitionKind             // returns the bonus feature transition kind when applicable.
}

// SpinAction contains generic spin action details.
type SpinAction struct {
	altSymbols         bool
	bonusSymbol        bool
	playerChoice       bool
	nrOfSpins          uint8
	stage              SpinActionStage
	result             SpinActionResult
	symbol             utils.Index
	id                 int
	name               string
	kind               string
	config             string
	alternate          SpinActioner
	chanceModifier     ChanceModifier
	testChoicesFilters []EndpointFilterer
	triggerFilters     []SpinDataFilterer
	stickyFilters      []SpinDataFilterer
	clearFilters       []SpinDataFilterer
}

// init initializes the spin action.
func (a *SpinAction) init(stage SpinActionStage, result SpinActionResult, kind string) *SpinAction {
	a.stage = stage
	a.result = result

	if strings.HasPrefix(kind, "*slots.") {
		a.kind = strings.ReplaceAll(kind, "*slots.", "")
	} else {
		a.kind = kind
	}

	return a
}

// Describe updates the code and name of the action.
func (a *SpinAction) Describe(id int, name string) {
	a.id = id
	a.name = name
}

// ID returns the id of the action.
func (a *SpinAction) ID() int {
	return a.id
}

// Name returns the name of the action.
func (a *SpinAction) Name() string {
	return a.name
}

// Kind returns the kind of the action.
func (a *SpinAction) Kind() string {
	return a.kind
}

// Config returns the configuration string of the action.
func (a *SpinAction) Config() string {
	return a.config
}

// Stage returns the stage at which the spin action should be tested.
func (a *SpinAction) Stage() SpinActionStage {
	return a.stage
}

// Result returns the result of the spin action.
func (a *SpinAction) Result() SpinActionResult {
	return a.result
}

// NrOfSpins returns the number of free spins from the spin action.
func (a *SpinAction) NrOfSpins(_ interfaces.Generator) uint8 {
	return a.nrOfSpins
}

// AltSymbols indicates to use the alternate symbol set for free spins.
func (a *SpinAction) AltSymbols() bool {
	return a.altSymbols
}

// Alternate returns the alternate action for this action if was defined.
func (a *SpinAction) Alternate() SpinActioner {
	return a.alternate
}

// Reschedule changes the stage at which the action is tested.
func (a *SpinAction) Reschedule(stage SpinActionStage) {
	a.stage = stage
}

// BonusSymbol implements the default SpinActioner.BonusSymbol() interface.
func (a *SpinAction) BonusSymbol() bool {
	return a.bonusSymbol
}

// PlayerChoice implements the default SpinActioner.PlayerChoice() interface.
func (a *SpinAction) PlayerChoice() bool {
	return a.playerChoice
}

// Triggered implements the default SpinActioner.Triggered() interface.
// The reason it returns an action instance is so that any alternate action can be reported as the triggering action,
// and the calling function will be able to tell what the "results" of the triggering action are.
func (a *SpinAction) Triggered(*Spin) SpinActioner {
	// by default not triggered
	return nil
}

// TriggeredWithState implements the default SpinActioner.TriggeredWithState() interface.
// The reason it returns an action instance is so that any alternate action can be reported as the triggering action,
// and the calling function will be able to tell what the "results" of the triggering action are.
func (a *SpinAction) TriggeredWithState(*Spin, *SymbolsState) SpinActioner {
	// by default not triggered
	return nil
}

// Payout implements the default SpinActioner.Payout() interface.
// It should return an action instance in derived objects if any payouts were awarded.
func (a *SpinAction) Payout(*Spin, *results.Result) SpinActioner {
	// by default no payouts
	return nil
}

// Penalty implements the default SpinActioner.Penalty() interface.
// It should return an action instance in derived objects if any penalties were imposed.
func (a *SpinAction) Penalty(*Spin, *results.Result) SpinActioner {
	// by default no payouts
	return nil
}

// Nudge implements the default SpinActioner.Nudge() interface.
// It should return an action instance in derived objects if any nudges were added.
func (a *SpinAction) Nudge(*Spin, *SpinResult) SpinActioner {
	// by default no nudges
	return nil
}

// InstantBonus implements the default SpinActioner.InstantBonus() interface.
// It should return an instant bonus result when applicable in derived objects.
func (a *SpinAction) InstantBonus(_ *Spin) interfaces.Objecter2 {
	// by default no instant bonus
	return nil
}

// BonusSelect implements the default SpinActioner.BonusSelect() interface.
// It should return a bonus selector result when applicable in derived objects.
func (a *SpinAction) BonusSelect(_ *Spin) interfaces.Objecter2 {
	// by default no bonus selector
	return nil
}

// BonusGame implements the default SpinActioner.BonusGame() interface.
// It should return a bonus game when applicable in derived objects.
func (a *SpinAction) BonusGame(_ *Spin) interfaces.Objecter2 {
	// by default no bonus game
	return nil
}

// FeatureTransition implements the default SpinActioner.FeatureTransition() interface.
// It should return the appropriate feature transition kind when applicable in derived objects.
func (a *SpinAction) FeatureTransition() FeatureTransitionKind {
	// by default no transition
	return 0
}

// CanTestChoices returns true if all testChoices filters return true or if there are no testChoices filters.
func (a *SpinAction) CanTestChoices(endpoint EndpointKind) bool {
	for ix := range a.testChoicesFilters {
		if !a.testChoicesFilters[ix](endpoint) {
			return false
		}
	}
	return true
}

// CanTrigger returns true if all trigger filters return true or if there are no trigger filters.
func (a *SpinAction) CanTrigger(spin *Spin) bool {
	for ix := range a.triggerFilters {
		if !a.triggerFilters[ix](spin) {
			return false
		}
	}
	return true
}

// CanSticky returns true if all sticky filters return true or if there are no sticky filters.
func (a *SpinAction) CanSticky(spin *Spin) bool {
	for ix := range a.stickyFilters {
		if !a.stickyFilters[ix](spin) {
			return false
		}
	}
	return true
}

// CanClear returns true if all clear filters return true or if there are no clear filters.
func (a *SpinAction) CanClear(spin *Spin) bool {
	for ix := range a.clearFilters {
		if !a.clearFilters[ix](spin) {
			return false
		}
	}
	return true
}

// WithChanceModifier adds a chance modifier function to the action.
func (a *SpinAction) WithChanceModifier(f ChanceModifier) {
	a.chanceModifier = f
}

// ModifyChance calculates a modified chance when a modifier is present, or returns the chance unmodified.
func (a *SpinAction) ModifyChance(chance float64, spin *Spin) float64 {
	if a.chanceModifier == nil {
		return chance
	}
	return a.chanceModifier.Exec(chance, spin)
}

// WithTestChoicesFilters adds one or more testChoices-filters to the action.
func (a *SpinAction) WithTestChoicesFilters(filters ...EndpointFilterer) {
	a.testChoicesFilters = filters
}

// WithTriggerFilters adds one or more trigger-filters to the action.
func (a *SpinAction) WithTriggerFilters(filters ...SpinDataFilterer) {
	a.triggerFilters = filters
}

// WithStickyFilters adds one or more sticky-filters to the action.
func (a *SpinAction) WithStickyFilters(filters ...SpinDataFilterer) {
	a.stickyFilters = filters
}

// WithClearFilters adds one or more clear-filters to the action.
func (a *SpinAction) WithClearFilters(filters ...SpinDataFilterer) {
	a.clearFilters = filters
}

// WithStage modifies the stage to the action.
func (a *SpinAction) WithStage(stage SpinActionStage) {
	a.stage = stage
}

// SpinActions is a slice of spin actions.
type SpinActions []SpinActioner

// PurgeSpinActions returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeSpinActions(input SpinActions, capacity int) SpinActions {
	if cap(input) < capacity {
		return make(SpinActions, 0, capacity)
	}
	return input[:0]
}

// GetMaxID determines the highest action id in the slice of actions.
func (actions SpinActions) GetMaxID() int {
	var out int
	for ix := range actions {
		a := actions[ix]
		if a.ID() > out {
			out = a.ID()
		}
	}
	return out
}
