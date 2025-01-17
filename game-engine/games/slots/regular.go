package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// RegularParams contains the parameters for instantiating a new regular slots game.
type RegularParams struct {
	PrngLog     bool
	ReturnFlags bool
	Slots       *slots.Slots
	StartGrids  *StartGrids
	PrngCache   []int
}

// AcquireRegular instantiates a new regular slot machine game from the memory pool.
// It panics if it is not initialized with a valid slots instance.
// Note that currently there can only be one PaidAction!
// The embedded PRNG is instantiated as a non-logging buffer.
func AcquireRegular(params RegularParams) *Regular {
	//log.Printf("\n\n LOGGING NEW REGULAR FROM M POOL \n\n")
	r := regularProducer.Acquire().(*Regular)

	r.slots = params.Slots
	r.startGrids = params.StartGrids
	r.doubleSpin = r.slots.DoubleSpin()
	r.playerChoice = r.slots.PlayerChoice()
	r.prngLog = params.PrngLog
	r.stateResetRequired = false
	r.maxPayoutReached = false
	r.minPayoutReached = false
	r.returnFlags = params.ReturnFlags
	r.totalPayout = 0.0
	r.maxPayout = r.slots.MaxPayout()
	r.reverseWin = r.slots.ReverseWin()

	if r.maxPayout == 0 {
		r.maxPayout = math.MaxFloat64
	}

	r.actionsFirst = acquireActionHandler(r, r.slots.ActionsFirst())
	r.actionsFree = acquireActionHandler(r, r.slots.ActionsFree())
	r.actionsFirstBB = acquireActionHandler(r, r.slots.ActionsFirstBB())
	r.actionsFreeBB = acquireActionHandler(r, r.slots.ActionsFreeBB())

	if params.PrngLog || len(params.PrngCache) >= 2 {
		r.prngBuf = rng.AcquireBuffer(rng.AcquireRNG(), params.PrngLog).WithCache(params.PrngCache)
		r.prng = r.prngBuf
	} else {
		r.prng = rng.AcquireRNG()
	}

	if r.slots.SymbolsState() {
		r.symbolsState = slots.AcquireSymbolsState(r.slots.Symbols(), r.slots.ExcludeFromState()...)
	}

	if r.slots == nil {
		panic(consts.MsgNoSymbolsSlots)
	}

	r.spin = slots.AcquireSpin(r.slots, r.prng)
	r.spin.SetGamer(r)

	if s := r.slots.ScriptedRoundSelector(); s != nil {
		r.initScriptedRounds(s)
	}

	r.prepareRound(0)
	return r
}

// Slots returns the slot machine the game is based on.
func (r *Regular) Slots() *slots.Slots {
	return r.slots
}

// ConfigHash returns the sha256 hash of the game config.
func (r *Regular) ConfigHash() string {
	return r.slots.ConfigHash()
}

// ForSale indicates if it is possible to buy a free/bonus game.
// It returns the relevant bonus buy action if it exists.
func (r *Regular) ForSale(bonusKind uint8) *slots.PaidAction {
	return r.actionsFirstBB.getBonusAction(bonusKind)
}

// IsDoubleSpin returns true if the game has the double-spin feature turned on.
func (r *Regular) IsDoubleSpin() bool {
	return r.doubleSpin
}

// IsReverseWin returns true if the game is a "reverse win" style of game.
func (r *Regular) IsReverseWin() bool {
	return r.reverseWin
}

// AllowPlayerChoices returns true if the game allows players to make choices which direct the game play.
func (r *Regular) AllowPlayerChoices() bool {
	return r.playerChoice
}

// NeedPlayerChoice returns true if the game expects a player choice.
func (r *Regular) NeedPlayerChoice() bool {
	return r.makeChoice
}

// LastSpin creates a dummy result for the last spin.
// This can be used to initialize a front-end with a grid of random symbols when the user starts a completely new game.
// Make sure to call Release() on the result if you are done with it.
func (r *Regular) LastSpin() *results.Result {
	return results.AcquireResult(slots.AcquireSpinResult(r.spin), results.SpinData)
}

// ResultCount returns the number of results so far.
func (r *Regular) ResultCount() int {
	return len(r.results)
}

// TotalPayout returns the total payout of the round or MaxPayout if it was reached.
func (r *Regular) TotalPayout() float64 {
	if r.maxPayoutReached {
		return r.maxPayout
	}
	return r.totalPayout
}

// MaxPayout returns the maximum payout for a round.
func (r *Regular) MaxPayout() float64 {
	return r.maxPayout
}

// MaxPayoutReached indicates if the maximum payout has been reached during the round.
func (r *Regular) MaxPayoutReached() bool {
	return r.maxPayoutReached
}

// MinPayoutReached indicates if the minimum payout has been reached during a "reverse win" round.
func (r *Regular) MinPayoutReached() bool {
	return r.minPayoutReached
}

// BonusSymbol returns the currently selected bonus symbol, or MaxIndex if not selected.
func (r *Regular) BonusSymbol() utils.Index {
	return r.spin.BonusSymbol()
}

// PrngLog returns the log (if any) of the embedded PRNG.
func (r *Regular) PrngLog() ([]int, []int) {
	return r.prngBuf.Log()
}

// SpinState returns a clone of the spin state for the game or nil if no spin state is kept.
// Make sure to call Release() on the result if you are done with it.
func (r *Regular) SpinState() *slots.SpinState {
	if r.spinState == nil {
		return nil
	}
	return r.spinState.Clone().(*slots.SpinState)
}

// SymbolsState returns a clone of the current symbols state for the game or nil if no symbols state is kept.
// Make sure to call Release() on the result if you are done with it.
func (r *Regular) SymbolsState() *slots.SymbolsState {
	if r.symbolsState == nil {
		return nil
	}
	return r.symbolsState.Clone().(*slots.SymbolsState)
}

func (r *Regular) initScriptedRounds(s *slots.ScriptedRoundSelector) {
	r.scriptSelector = s
	r.scriptWeights = s.GetWeighting()

	m := s.GetScripts()

	var maxID int
	for _, v := range m {
		if v.ID() > maxID {
			maxID = v.ID()
		}
	}

	maxID++
	reserve := utils.FixArraySize(maxID, 32)

	if r.scripts == nil {
		r.scripts = make(map[int]*slots.ScriptedRound, reserve)
	}

	r.scriptHandlers1 = purgeActionHandlers(r.scriptHandlers1, reserve)[:maxID]
	r.scriptHandlers2 = purgeActionHandlers(r.scriptHandlers2, reserve)[:maxID]

	for _, v := range m {
		r.scripts[v.ID()] = v
		actions1, actions2 := v.Actions()
		if actions1 != nil {
			r.scriptHandlers1[v.ID()] = acquireActionHandler(r, actions1)
		}
		if actions2 != nil {
			r.scriptHandlers2[v.ID()] = acquireActionHandler(r, actions2)
		}
	}
}

// RestoreState restores an externally saved game state back into the game.
// It makes a deep copy of the given states, so the caller remains responsible for calling Release() on them.
func (r *Regular) RestoreState(spinState *slots.SpinState, symbolsState *slots.SymbolsState) *Regular {
	if r.spinState != nil {
		r.spinState.Release()
		r.spinState = nil
	}
	if spinState != nil {
		r.spinState = spinState.Clone().(*slots.SpinState)
	}

	if r.symbolsState != nil {
		r.symbolsState.Release()
		r.symbolsState = nil
	}
	if symbolsState != nil {
		r.symbolsState = symbolsState.Clone().(*slots.SymbolsState)
	}

	return r
}

// Round performs a spin and awards free spins/bonus games from those spins and returns the calculated results.
// Note that the returned slice of results will remain valid until the function is called again.
// Round retains ownership of these results, so DO NOT call Release() on them!
// If bonusBuy > 0 the function will "buy" a bonus game;
// e.g. it will force activation of the related PaidAction with a grid of symbols that has no matching paylines.
func (r *Regular) Round(bonusBuy uint8) results.Results {
	r.prepareRound(bonusBuy)
	return r.doRound()
}

// RoundChoices is the same as Round except that it also tests the given player choices.
// This is used when there is an instant bonus requiring a player choice before the first spin.
func (r *Regular) RoundChoices(bonusBuy uint8, choices map[string]string) results.Results {
	r.prepareRound(bonusBuy)
	r.testChoices(slots.EndpointRound, choices)
	return r.doRound()
}

// Debug is a debug function to perform a Round with pre-calculated symbols for the first spin.
// This can be used for testing specific paylines, actions, free games/bonus game.
// It can also be used in simulators to force the PaidTrigger without the overhead of "fixing" the symbol grid.
// Results cannot be guaranteed if len(indexes) < reelCount * rowCount.
func (r *Regular) Debug(indexes utils.Indexes, bonusBuy uint8, resume bool, choices map[string]string) results.Results {
	if resume {
		r.prepareResume(choices)
	} else {
		r.prepareRound(bonusBuy)
		r.testChoices(slots.EndpointRound, choices)
	}

	r.debugInitial = true
	r.spin.Debug(indexes)
	r.getResults()
	return r.results
}

// Scripted is a debug function to perform a Round with the given script identifier.
// If the script exists it will be used to perform the round, otherwise a normal round is performed.
func (r *Regular) Scripted(id int, bonusBuy uint8, resume bool, choices map[string]string) results.Results {
	if resume {
		r.prepareResume(choices)
	} else {
		r.prepareRound(bonusBuy)
		r.testChoices(slots.EndpointRound, choices)
	}

	r.debugScript = true
	r.spin.SetDebug(true)
	r.scriptID = id
	return r.doRound()
}

// PrngCache loads the PRNG cache; it then plays a round in debug mode, or, resumes a round in debug mode.
// This can be used for testing an entire round using the PRNG log output from an earlier round,
// and to check that the results are exactly the same.
func (r *Regular) PrngCache(cache []int, bonusBuy uint8, resume bool, choices map[string]string) results.Results {
	if resume {
		r.prepareResume(choices)
	} else {
		r.prepareRound(bonusBuy)
		r.testChoices(slots.EndpointRound, choices)
	}

	if r.prngBuf != nil {
		r.debugPRNG = true
		r.spin.SetDebug(true)
		r.prngBuf.WithCache(cache)
	}
	return r.doRound()
}

// RoundResume resumes the game after a player choice.
func (r *Regular) RoundResume(choices map[string]string) results.Results {
	r.prepareResume(choices)
	return r.doRound()
}

func (r *Regular) prepareResume(choices map[string]string) {
	r.resuming = false

	// prepare the spin round first, before restoring state and setting player choices.
	r.prepareRound(r.bonusBuy)

	// restore free spins counter & spin state.
	if r.spinState != nil {
		r.freeStarted = true
		r.freeSpins = r.spinState.FreeSpins()
		r.spin.RestoreState(r.spinState)

		if r.freeSpins > 0 {
			r.spin.SetKind(slots.FreeSpin)
			r.resuming = true
		}
	}

	// test player choices.
	r.testChoices(slots.EndpointRoundResume, choices)
}

func (r *Regular) testChoices(endpoint slots.EndpointKind, choices map[string]string) {
	if len(choices) > 0 {
		r.choices = choices
		handler := r.getHandler()

		for ix := range handler.playerChoices {
			if a := handler.playerChoices[ix]; a.CanTestChoices(endpoint) {
				r.logEvent(a, a.TestChoices(r.spin, choices) != nil)
			}
		}
	}
}

// prepareRound prepares the regular slot machine game for a new game round.
func (r *Regular) prepareRound(bonusBuy uint8) {
	r.spin.ResetSpin()

	if bonusBuy > 0 {
		if r.paidAction = r.actionsFirstBB.getBonusAction(bonusBuy); r.paidAction != nil {
			r.paidFirst = true
			r.bonusBuy = bonusBuy
			r.spin.SetBonusBuy(bonusBuy)
		}
	}

	r.locked = r.locked[:0]
	r.debugInitial = false
	r.debugPRNG = false
	r.debugScript = false
	r.freeSpins = 0
	r.scriptID = 0
	r.needRefill = false
	r.superSpin = false
	r.makeChoice = false
	r.totalPayout = 0.0
	r.maxPayoutReached = false
	r.minPayoutReached = false
	r.bonusGamePlayed = false
	r.freeStarted = false
	r.bonusGame = nil

	r.results = results.ReleaseResults(r.results)

	if r.doubleSpin {
		if r.spinState != nil {
			r.spin.SetKind(slots.SecondSpin)
			r.spin.RestoreState(r.spinState)
			if r.startGrids != nil {
				r.startGrids.SetLastGrid(r.spinState.StartGrid())
			}
		} else {
			r.spin.SetKind(slots.FirstSpin)
		}
	}

	if r.prngLog {
		r.prngBuf.ResetLog()
		r.prngLast = 0
		r.eventLast = 0
	}
}

// doRound performs a spin round.
func (r *Regular) doRound() results.Results {
	if r.scriptSelector != nil && r.scriptSelector.BonusBuyAllowed(r.bonusBuy) && r.scriptSelector.Triggered(r.spin) {
		// produce a random scripted round.
		for ix := 0; ix < 10; ix++ {
			id := int(r.scriptWeights.RandomIndex(r.prng))
			if s := r.scripts[id]; s != nil && s.BonusBuyAllowed(r.bonusBuy) {
				r.scriptID = id
				break
			}
		}
	}

	handler := r.getHandler()
	if handler.testPreSpinActions() {
		return r.results
	}

	if !r.resuming && r.bonusGame == nil {
		r.spin.Spin()
	}

	if r.paidFirst {
		if r.slots.ClusterPays() {
			for ix := range handler.regularPayouts {
				if a, ok := handler.regularPayouts[ix].(*slots.PayoutAction); ok {
					if a.HaveClusterPayouts() {
						a.RemoveClusterPays(r.spin)
					}
				}
			}
		} else {
			r.spin.MismatchPaylines()
		}

		r.spin.ForcePaidAction(r.paidAction)
		r.paidFirst = false
	}

	r.getResults()

	r.resuming = false
	return r.results
}

// getResults compiles the results of the current spin and subsequent free spins/bonus games.
func (r *Regular) getResults() {
	if !r.resuming && r.bonusGame == nil {
		r.getResultCurrent()
	}

	for !r.makeChoice && !r.maxPayoutReached && !r.minPayoutReached && (r.bonusGame != nil || r.freeSpins > 0 || r.needRefill) {
		if r.bonusGame == nil {
			r.playFreeSpin()
		} else {
			r.playBonusGame()
		}
	}

	if r.makeChoice {
		// when there's a player choice request, we need to remember the spin state.
		r.spinState = slots.AcquireSpinState(r.spin)
	} else if !r.doubleSpin && r.spinState != nil {
		r.spinState.Release()
		r.spinState = nil
	}

	if r.stateResetRequired {
		r.symbolsState.ResetState()
		r.stateResetRequired = false
		// reset the state in the last result as well!
		r.currResult.AddState(r.symbolsState.DeepCopy())
	}

	if r.reverseWin {
		// fix any division penalties in the results.
		r.totalPayout = results.FixPenalties(r.results)

		// check max payout, as we ignore it during the free spins.
		if r.totalPayout >= r.maxPayout {
			r.maxPayoutReached = true
			r.spinData.SetMaxPayout()
		}
	}

	// reset temporary round spin data & result.
	r.spinData = nil
	r.currResult = nil
}

// getResultCurrent compiles the results of the current spin.
func (r *Regular) getResultCurrent() {
	r.testActions()

	if r.currResult != nil {
		r.results = append(r.results, r.currResult)
	}
	r.debugInitial = false // make sure revise actions happen during free spins

	// for double spins, set or reset the spin state here before doing free spins!
	if r.doubleSpin {
		if r.spinState == nil && !r.superSpin {
			r.spinState = slots.AcquireSpinState(r.spin)
		} else if r.spinState != nil {
			r.spinState.Release()
			r.spinState = nil
		}
	}
}

// playFreeSpin plays and compiles the result for a free spin or refill spin.
func (r *Regular) playFreeSpin() {
	if r.needRefill {
		if r.superSpin {
			r.spin.SetKind(slots.SuperSpin)
		} else {
			r.spin.SetKind(slots.RefillSpin)
		}
		r.spin.Refill()
		r.needRefill = false
	} else {
		r.freeStarted = true

		if r.doubleSpin {
			// for double spins, we need to set or reset the spin state for every free spin.
			if r.spinState == nil {
				r.spin.SetKind(slots.FirstFreeSpin)
				r.spin.ResetSticky()
				if r.startGrids != nil {
					r.spinState.SetStartGrid(r.startGrids.LastGrid())
				}
				r.spin.Spin()
				r.spinState = slots.AcquireSpinState(r.spin)
			} else {
				r.spin.SetKind(slots.SecondFreeSpin)
				r.spinState.Release()
				r.spinState = nil
				r.freeSpins--
				r.spin.Spin()
			}
		} else {
			r.spin.SetKind(slots.FreeSpin)
			r.freeSpins--
			r.spin.Spin(r.locked...)
		}
	}

	r.spin.SetFreeSpins(r.freeSpins) // this allows actions to see if there are free spins remaining.
	r.testActions()

	// super-shape completed?
	if r.spin.Kind() == slots.SuperSpin && !r.needRefill {
		if r.currResult.AwardedFreeGames == 0 && r.freeSpins > 0 {
			r.freeSpins--
			if r.currResult != nil {
				r.currResult.SetFreeGames(r.freeSpins) // make sure the result has the samne value!
			}
		}

		// reset the state, so we start with a first spin for the next free game.
		if r.spinState != nil {
			r.spinState.Release()
			r.spinState = nil
		}
	}

	if r.currResult != nil {
		r.results = append(r.results, r.currResult)
	}
	r.debugInitial = false // make sure revise actions happen during free spins
}

// playBonusGame plays and compiles the results for a bonus game.
func (r *Regular) playBonusGame() {
	t := r.bonusGame
	if b := t.BonusGame(r.spin); b != nil {
		r.bonusGamePlayed = true

		if bw, ok := b.(*wheel.BonusWheelResult); ok {
			bw.SetChoices(r.choices)
			bw.LogEvent(r.getEvent(t, true))
			if r.prngLog {
				// add the PRNG log.
				l1, l2 := r.prngBuf.Log()
				bw.SetLog(l1[r.prngLast:], l2[r.prngLast:])
				r.prngLast = len(l1)
				r.eventLast = r.prngLast
			}
			r.results = append(r.results, results.AcquireResult(b, results.BonusWheelData))
		}
	}

	r.testActionsBonusGame()
}

// testActions tests all actions on the current state of the spin, establishes a spin result,
// records awarded payouts/free spins/bonus games, performs wild expansions, performs symbol clearing or stickiness, etc.
// For a "paid" round, grid revisements are skipped.
func (r *Regular) testActions() {
	handler := r.getHandler()

	if !r.superSpin && r.doubleSpin && r.spin.Kind().IsFirst() {
		r.spin.ResetSuper() // safe to reset the super symbol here.

		r.initCurrData()

		// first spin of a double, so only fix-up the grid and test for super shapes and/or stickiness.
		if !r.debugInitial {
			handler.testGridRevisements()
		}

		handler.testGridActions()
		handler.testStickiness()

		if r.symbolsState != nil {
			r.currResult.AddState(r.symbolsState.DeepCopy())
		}
		r.currResult.SetFreeGames(r.freeSpins)

		if r.prngLog {
			l1, l2 := r.prngBuf.Log()
			r.spinData.SetLog(l1[r.prngLast:], l2[r.prngLast:])
			r.prngLast = len(l1)
			r.eventLast = r.prngLast
		}

		// store debug and bonus buy flags.
		if r.debugInitial || r.debugPRNG || r.debugScript {
			r.spinData.SetDebug(true)
		}
		if r.paidAction != nil {
			r.spinData.SetBonusBuy(r.bonusBuy, r.paidAction)
		}

		return
	}

	r.initCurrData()

	if (r.paidAction == nil || r.spin.Kind() != slots.FirstSpin) && !r.superSpin && !r.debugInitial {
		handler.testGridRevisements()
	}

	if !r.superSpin {
		handler.testBeforeExpansions()
		if r.doubleSpin {
			// reset stickies here after we performed morphing & expansion actions that need to know!
			r.spin.ResetSticky()
			r.spinData.SetSticky(r.spin)
		}
	}

	handler.testGridActions()

	if !r.superSpin {
		handler.testRegularPayouts()
		handler.testRegularPenalties()
		handler.testInjections()
		handler.testAfterExpansions()
		handler.testExtraPayouts()
		handler.testStateChanges()

		r.spin.ResetSuper() // safe to reset the super symbol now.
	}

	if !r.superSpin {
		handler.testBonuses()
		r.spin.SetFreeSpins(r.freeSpins) // this allows subsequent actions to see if free spins were awarded.

		if !r.doubleSpin {
			handler.testStickiness()
		}
		handler.testClearing()

		if r.slots.CascadingReels() {
			if r.spin.CascadeFloatingSymbols() {
				r.spinData.SetAfterCascade(r.spin)
			}
		}

		if b := r.spin.BonusSymbol(); b != utils.MaxIndex {
			r.spinData.SetBonusSymbol(b)
		}

		r.locked = r.spin.Locked(r.locked)
	}

	if r.symbolsState != nil {
		r.currResult.AddState(r.symbolsState.DeepCopy())
	}
	r.currResult.SetFreeGames(r.freeSpins)

	r.totalPayout += r.currResult.Total
	switch {
	case r.totalPayout >= r.maxPayout:
		if r.reverseWin {
			r.totalPayout = r.maxPayout
		} else {
			r.maxPayoutReached = true
			r.spinData.SetMaxPayout()
		}

	case r.totalPayout <= 0:
		r.totalPayout = 0
		if r.reverseWin && r.spin.SpinSeq() > 1 {
			r.minPayoutReached = true
		}
	}

	if r.bonusGame == nil && !r.makeChoice && !r.maxPayoutReached && !r.minPayoutReached && (r.freeSpins > 0 || r.needRefill) {
		handler.testPreBonus()
	}

	if r.spinData != nil {
		r.recordLogAndDebug()
	}
}

// testActionsBonusGame tests the appropriate actions after a bonus game,
// and records awarded payouts/free spins/bonus games, etc.
func (r *Regular) testActionsBonusGame() {
	handler := r.getHandler()
	t := r.bonusGame

	handler.testBonuses()

	if r.bonusGame == t {
		// remove it only after new bonuses have been tested!
		// this makes sure we don't get into an endless loop :)
		r.bonusGame = nil
	}

	if r.bonusGame == nil && !r.makeChoice && !r.maxPayoutReached && !r.minPayoutReached && (r.freeSpins > 0 || r.needRefill) {
		handler.testPreBonus()
	}

	if r.spinData != nil {
		r.recordLogAndDebug()
	}
}

func (r *Regular) getHandler() *actionHandler {
	if r.scriptID > 0 && r.scriptSelector != nil {
		// inside a scripted round.
		a1, a2 := r.scriptHandlers1[r.scriptID], r.scriptHandlers2[r.scriptID]
		if r.freeStarted && a2 != nil {
			return a2
		} else if a1 != nil {
			return a1
		}
	}

	if r.freeStarted {
		// free spins during a normal round.
		if r.paidAction != nil {
			return r.actionsFreeBB
		}
		return r.actionsFree
	}

	// all other cases; e.g. first spin during a normal round.
	if r.paidAction != nil {
		return r.actionsFirstBB
	}
	return r.actionsFirst
}

func (r *Regular) recordLogAndDebug() {
	if r.returnFlags {
		r.spinData.SetRoundFlags(r.spin)
	}

	if r.slots.HaveExportedFlags() {
		r.spinData.SetExportFlags(r.spin)
	}

	if r.prngLog {
		// add the PRNG log.
		l1, l2 := r.prngBuf.Log()
		r.spinData.SetLog(l1[r.prngLast:], l2[r.prngLast:])
		r.prngLast = len(l1)
		r.eventLast = r.prngLast
	}

	// record the debug and bonus buy flag.
	if r.debugInitial || r.debugPRNG || r.debugScript {
		r.spinData.SetDebug(true)
	}
	if r.paidAction != nil {
		r.spinData.SetBonusBuy(r.bonusBuy, r.paidAction)
	}
}

func (r *Regular) initCurrData() {
	if r.scriptID > 0 {
		if flag := r.scriptSelector.SpinFlag(); flag >= 0 {
			r.spin.SetRoundFlag(flag, r.scriptID)
		}
	}
	r.spinData = slots.AcquireSpinResult(r.spin)
	r.spinData.PrngLog.SetScriptID(r.scriptID)

	r.currResult = results.AcquireResult(r.spinData, results.SpinData)

	if r.prngLog {
		r.eventLast = r.prngBuf.LogSize()
	}

	if len(r.eventBuf) > 0 {
		for ix := range r.eventBuf {
			r.spinData.CloneEvent(r.eventBuf[ix])
		}
		r.eventBuf = results.ReleaseEvents(r.eventBuf)
	}

	if r.choices != nil {
		r.spinData.SetChoices(r.choices)
		r.choices = nil
	}
}

func (r *Regular) awardFreeSpins(free uint8) {
	if free > 0 {
		r.freeSpins += uint64(free)
		if r.currResult != nil {
			r.currResult.AwardFreeGames(free)
		}
	}
}

func (r *Regular) logEvent(a slots.SpinActioner, t bool) {
	e := r.getEvent(a, t)
	if r.spinData == nil {
		r.eventBuf = append(r.eventBuf, e)
	} else {
		r.spinData.LogEvent(e)
	}
}

func (r *Regular) getEvent(a slots.SpinActioner, t bool) *results.Event {
	if r.prngBuf == nil {
		return results.AcquireEvent(a.ID(), t, nil, nil)
	}

	l1, l2 := r.prngBuf.Log()
	last := r.eventLast
	r.eventLast = r.prngBuf.LogSize()
	return results.AcquireEvent(a.ID(), t, l1[last:], l2[last:])
}

// Regular represents a regular slot machine game.
// Regular is not safe for concurrent use across multiple go-routines.
// Keep fields ordered by ascending SizeOf().
type Regular struct {
	doubleSpin         bool                         // double spin feature.
	paidFirst          bool                         // indicates the bonus buy/bonus bet feature.
	playerChoice       bool                         // player choice feature.
	havePaylines       bool                         // indicates there are specific paylines.
	haveAllPaylines    bool                         // indicates all paylines feature is active.
	haveScatterPays    bool                         // indicates there are scatter payouts.
	haveClusterPays    bool                         // indicates there are cluster payouts.
	prngLog            bool                         // indicates to log the PRNG IntN/IntsN input/output values.
	debugInitial       bool                         // indicates a debug round with preloaded initial grid.
	debugPRNG          bool                         // indicates a debug round with PRNG cache.
	debugScript        bool                         // indicates a debug round with scripted round.
	resuming           bool                         // indicates the game resumes after a player choice.
	needRefill         bool                         // indicates a partial filling re-spin is required.
	superSpin          bool                         // indicates re-spins/free spins are private, blocking specific tests.
	makeChoice         bool                         // indicates the player must make a choice.
	stateResetRequired bool                         // indicates the game state requires a reset.
	maxPayoutReached   bool                         // indicates if the maximum payout limit has been reached.
	minPayoutReached   bool                         // indicates if the minimum payout limit has been reached (reverse win).
	bonusGamePlayed    bool                         // indicates if a bonus game was played during the current round.
	returnFlags        bool                         // indicates if round flags should be exported in spin results.
	freeStarted        bool                         // indicates if the free spins have started.
	reverseWin         bool                         // indicates the game is a "reverse win" style of game.
	bonusBuy           uint8                        // the bonus buy type.
	prngLast           int                          // size of the PRNG logs when last checked.
	eventLast          int                          // size of the PRNG logs when last checked for action events.
	scriptID           int                          // ID of a selected scripted round or zero.
	freeSpins          uint64                       // number of free spins remaining.
	maxPayout          float64                      // maximum payout per slots round.
	totalPayout        float64                      // keeps track of the total payout during a round.
	slots              *slots.Slots                 // slots machine config.
	spin               *slots.Spin                  // the spin grid and relevant details.
	paidAction         *slots.PaidAction            // the bonus buy/bet action durign a bonus buy/bet round.
	startGrids         *StartGrids                  // start grids (optional).
	spinState          *slots.SpinState             // the state of a previous spin (used for double spin feature and player choice).
	symbolsState       *slots.SymbolsState          // the state of the symbols (used for CCB).
	spinData           *slots.SpinResult            // current spin outcome.
	bonusWheel         *wheel.BonusWheelResult      // current bonus wheel outcome.
	currResult         *results.Result              // current result.
	prngBuf            *rng.Buffer                  // PRNG buffering support.
	scriptWeights      utils.WeightedGenerator      // weighting for selecting a scripted round ID.
	scriptSelector     *slots.ScriptedRoundSelector // selector for scripted rounds.
	actionsFirst       *actionHandler               // actions for the first spin without bonus buy or bonus bet.
	actionsFree        *actionHandler               // actions for the free spins without bonus buy or bonus bet.
	actionsFirstBB     *actionHandler               // actions for the first spin with bonus buy or bonus bet.
	actionsFreeBB      *actionHandler               // actions for the free spins with bonus buy or bonus bet.
	prng               interfaces.Generator         // single PRNG for all randomness in the game.
	bonusGame          slots.SpinActioner           // the action with the bonus game to play next.
	locked             utils.UInt8s                 // reel lock indicators.
	results            results.Results              // holds the results of the current spin round.
	eventBuf           results.Events               // buffer for events happening before the spin.
	choices            map[string]string            // temporarily holds player choices when resuming a spin round.
	scripts            map[int]*slots.ScriptedRound // map of scripted rounds.
	scriptHandlers1    []*actionHandler             // list of action handlers for first spin of scripted rounds.
	scriptHandlers2    []*actionHandler             // list of action handlers for free spins of scripted rounds.
	hash               string                       // sha256 hash of the game config.
	pool.Object
}

const (
	lockedCap    = 8
	eventBufCap  = 8
	resultsCap   = 16
	resultMaxCap = 256
)

// regularProducer is the memory pool for regular slot machines.
// Make sure to initialize all slices appropriately!
var regularProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &Regular{
		locked:   make(utils.UInt8s, 0, lockedCap),
		results:  results.NewResults(resultsCap),
		eventBuf: make(results.Events, 0, eventBufCap),
	}
	return r, r.reset
})

// reset clear the regular slot machine.
func (r *Regular) reset() {
	if r != nil {
		if r.spin != nil {
			r.spin.Release()
			r.spin = nil
		}

		if r.prng != nil {
			r.prng.ReturnToPool()
			r.prng = nil
		}

		if r.spinState != nil {
			r.spinState.Release()
			r.spinState = nil
		}

		if r.symbolsState != nil {
			r.symbolsState.Release()
			r.symbolsState = nil
		}

		if r.scriptWeights != nil {
			r.scriptWeights.Release()
			r.scriptWeights = nil
		}

		r.doubleSpin = false
		r.paidFirst = false
		r.playerChoice = false
		r.havePaylines = false
		r.haveAllPaylines = false
		r.haveScatterPays = false
		r.haveClusterPays = false
		r.resuming = false
		r.prngLog = false
		r.debugInitial = false
		r.debugPRNG = false
		r.debugScript = false
		r.needRefill = false
		r.superSpin = false
		r.makeChoice = false
		r.stateResetRequired = false
		r.maxPayoutReached = false
		r.minPayoutReached = false
		r.bonusGamePlayed = false
		r.returnFlags = false
		r.freeStarted = false
		r.reverseWin = false

		r.bonusBuy = 0
		r.prngLast = 0
		r.eventLast = 0
		r.scriptID = 0
		r.freeSpins = 0
		r.maxPayout = 0.0
		r.totalPayout = 0.0

		r.slots = nil
		r.paidAction = nil
		r.startGrids = nil
		r.spinData = nil
		r.bonusWheel = nil
		r.currResult = nil
		r.prngBuf = nil
		r.scriptSelector = nil
		r.bonusGame = nil

		if r.actionsFirst != nil {
			r.actionsFirst.Release()
			r.actionsFirst = nil
		}
		if r.actionsFree != nil {
			r.actionsFree.Release()
			r.actionsFree = nil
		}
		if r.actionsFirstBB != nil {
			r.actionsFirstBB.Release()
			r.actionsFirstBB = nil
		}
		if r.actionsFreeBB != nil {
			r.actionsFreeBB.Release()
			r.actionsFreeBB = nil
		}

		r.locked = r.locked[:0]

		r.results = results.ReleaseResults(r.results)
		if cap(r.results) > resultMaxCap {
			r.results = results.NewResults(resultsCap)
		}

		r.eventBuf = results.ReleaseEvents(r.eventBuf)
		r.choices = nil

		clear(r.scripts)
		r.scriptHandlers1 = releaseActionHandlers(r.scriptHandlers1)
		r.scriptHandlers2 = releaseActionHandlers(r.scriptHandlers2)
	}
}
