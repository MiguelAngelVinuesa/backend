package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// getBonusAction returns the action associated with the given bonus buy kind.
func (h *actionHandler) getBonusAction(bonusKind uint8) *slots.PaidAction {
	list := h.forSale
	switch len(list) {
	case 0:
		return nil
	case 1:
		return list[0].(*slots.PaidAction)
	}

	for ix := range list {
		if a, ok := list[ix].(*slots.PaidAction); ok && a.BonusKind() == bonusKind {
			return a
		}
	}
	return list[0].(*slots.PaidAction)
}

// testPreSpinActions tests the appropriate actions before the first spin.
// It returns true if an instant bonusBuy was awarded.
func (h *actionHandler) testPreSpinActions() bool {
	list, spin, spinData, logEvent := h.preSpin, h.r.spin, h.r.spinData, h.r.logEvent

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			t := a.Triggered(spin)
			if t == nil {
				logEvent(a, false)
				continue
			}

			if spinData != nil {
				spinData.SetTransition(t.FeatureTransition())
			}

			switch t.Result() {
			case slots.InstantBonus:
				if b := t.InstantBonus(spin); b != nil {
					if t.PlayerChoice() {
						h.r.makeChoice = true
					}

					if ib, ok := b.(*results.InstantBonus); ok {
						ib.SetChoices(h.r.choices)
						ib.LogEvent(h.r.getEvent(t, true))
						if h.r.prngLog {
							// add the PRNG log.
							l1, l2 := h.r.prngBuf.Log()
							ib.SetLog(l1[h.r.prngLast:], l2[h.r.prngLast:])
							h.r.prngLast = len(l1)
							h.r.eventLast = h.r.prngLast
						}
						h.r.results = append(h.r.results, results.AcquireResult(ib, results.InstantBonusData))
					}
				}

			case slots.SpecialResult:
				if b := t.BonusSelect(spin); b != nil {
					if bs, ok := b.(*results.BonusSelector); ok {
						bs.SetChoices(h.r.choices)
						bs.LogEvent(h.r.getEvent(t, true))
						if h.r.prngLog {
							// add the PRNG log.
							l1, l2 := h.r.prngBuf.Log()
							bs.SetLog(l1[h.r.prngLast:], l2[h.r.prngLast:])
							h.r.prngLast = len(l1)
							h.r.eventLast = h.r.prngLast
						}
						h.r.results = append(h.r.results, results.AcquireResult(bs, results.BonusSelectorData))
					}
				}

			case slots.BonusGame:
				h.r.bonusGame = t
				logEvent(t, true)

			default:
			}
		}
	}

	return h.r.makeChoice
}

// testGridRevisements processes the revise grid actions.
func (h *actionHandler) testGridRevisements() {
	list, spin, spinData, logEvent := h.reviseGrid, h.r.spin, h.r.spinData, h.r.logEvent

	var count int

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			if t := a.Triggered(spin); t == nil {
				logEvent(a, false)
			} else {
				switch t.Result() {
				case slots.ReelsNudged:
					t.Nudge(spin, spinData)

				case slots.Multiplier:
					spinData.SetMultiplier(spin)

				default:
					count++
				}

				logEvent(t, true)
			}
		}
	}

	if count > 0 {
		spinData.Update(spin)
	}
}

// testBeforeExpansions tests the expansions that need to be executed before testing the paylines.
func (h *actionHandler) testBeforeExpansions() {
	list, spin, spinData, logEvent, awardFreeSpins, prng := h.expandBefore, h.r.spin, h.r.spinData, h.r.logEvent, h.r.awardFreeSpins, h.r.prng

	var expanded bool

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			if t := a.Triggered(spin); t == nil {
				logEvent(a, false)
			} else {
				a.Expand(spin)
				expanded = true

				awardFreeSpins(t.NrOfSpins(prng))
				logEvent(t, true)
			}
		}
	}

	if expanded {
		spinData.SetAfterExpand(spin)
		if spin.HasEffect() {
			spinData.SetEffects(spin)
		}
	}
}

// testGridActions tests the grid for special features.
func (h *actionHandler) testGridActions() {
	list, spin, spinData, logEvent, awardFreeSpins := h.gridActions, h.r.spin, h.r.spinData, h.r.logEvent, h.r.awardFreeSpins

	h.r.superSpin = false
	oldMultiplier := spin.Multiplier()

	var free uint8

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			if t := a.Triggered(spin); t == nil {
				logEvent(a, false)
			} else {
				switch t.Result() {
				case slots.Refill:
					h.r.needRefill = true
					spinData.SetAfterClear(spin)
					spinData.SetSticky(spin)

				case slots.SuperRefill:
					h.r.superSpin = true
					h.r.needRefill = true
					spinData.SetAfterClear(spin)
					spinData.SetSticky(spin)

				case slots.Multiplier:
					spinData.SetMultiplier(spin)
					if m := spin.Multiplier(); m != oldMultiplier {
						if ma, ok := a.(*slots.MultiplierAction); ok {
							free += ma.BonusFreeSpins(m)
						}
					}

				case slots.Multipliers:
					spinData.SetMultipliers(spin)

				case slots.SymbolsInjected:
					spinData.SetAfterInject(spin)
					spinData.SetMultipliers(spin)

				default:
				}

				logEvent(t, true)
			}
		}
	}

	awardFreeSpins(free)
}

// testRegularPayouts tests the regular payout actions, and updates the result with any payouts.
func (h *actionHandler) testRegularPayouts() {
	list, spin, currResult, logEvent := h.regularPayouts, h.r.spin, h.r.currResult, h.r.logEvent
	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			logEvent(a, a.Payout(spin, currResult) != nil)
		}
	}
}

// testRegularPenalties tests the regular penalty actions, and updates the result with any penalties.
func (h *actionHandler) testRegularPenalties() {
	list, spin, currResult, logEvent := h.regularPenalties, h.r.spin, h.r.currResult, h.r.logEvent
	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			logEvent(a, a.Penalty(spin, currResult) != nil)
		}
	}
}

// testInjections tests symbol/cluster injections that need to be executed after testing the paylines.
// If a symbol or cluster is injected, the regular payouts will be re-tested.
func (h *actionHandler) testInjections() {
	list, spin, spinData, currResult, logEvent := h.injections, h.r.spin, h.r.spinData, h.r.currResult, h.r.logEvent

	var retestPayouts bool
	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			if t := a.Triggered(spin); t == nil {
				logEvent(a, false)
			} else {
				switch t.Result() {
				case slots.Multiplier:
					spinData.SetMultiplier(spin)

				case slots.SymbolsInjected:
					spinData.SetAfterInject(spin)
					spinData.SetMultipliers(spin)
					retestPayouts = true

				default:
				}

				logEvent(t, true)
			}
		}
	}

	if retestPayouts {
		currResult.ReleasePayouts()
		h.testRegularPayouts()
		h.testRegularPenalties()
	}
}

// testAfterExpansions tests the expansions that need to be executed after testing the paylines.
// If there is an expansion, the modified grid will be added to the SpinResult as the "after" image.
func (h *actionHandler) testAfterExpansions() {
	list, spin, spinData, currResult, logEvent, prng := h.expandAfter, h.r.spin, h.r.spinData, h.r.currResult, h.r.logEvent, h.r.prng

	var expanded bool

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			if t := a.Triggered(spin); t == nil {
				logEvent(a, false)
			} else {
				a.Expand(spin)
				expanded = true

				free := t.NrOfSpins(prng)
				h.r.freeSpins += uint64(free)
				currResult.AwardFreeGames(free)

				logEvent(t, true)
			}
		}
	}

	if expanded {
		spinData.SetAfterExpand(spin)
		if spin.HasEffect() {
			spinData.SetEffects(spin)
		}
	}
}

// testStateChanges tests the game state change actions.
func (h *actionHandler) testStateChanges() {
	list, spin, symbolsState, currResult, logEvent, prng := h.state, h.r.spin, h.r.symbolsState, h.r.currResult, h.r.logEvent, h.r.prng

	for ix := range list {
		a := list[ix]
		a.StateUpdate(spin, symbolsState)
		if a.CanTrigger(spin) {
			if t := a.StateTriggered(symbolsState); t == nil {
				logEvent(a, false)
			} else {
				switch t.Result() {
				case slots.FreeSpins:
					free := a.NrOfSpins(prng)
					h.r.freeSpins += uint64(free)
					currResult.AwardFreeGames(free)
					if a.AltSymbols() {
						spin.SetBonusSymbol(utils.MaxIndex, true)
					}

				default:
				}

				h.r.stateResetRequired = true
				logEvent(t, true)
			}
		}
	}
}

// testExtraPayouts tests the extra payout actions.
// An example of this would be a ScatterAction with the Payout kind.
func (h *actionHandler) testExtraPayouts() {
	list, spin, currResult, logEvent := h.extraPayouts, h.r.spin, h.r.currResult, h.r.logEvent

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			if t := a.Payout(spin, currResult); t == nil {
				logEvent(a, false)
			} else {
				logEvent(t, true)
			}
		}
	}
}

// testBonuses tests the bonusBuy actions to award any free games, bonusBuy game or bonusBuy symbol payout.
// An optional bonusBuy symbol is activated only after all actions have been tested.
// If a symbol or cluster is injected, the regular payouts will be re-tested.
func (h *actionHandler) testBonuses() {
	list, spin, spinData, currResult, logEvent, awardFreeSpins, prng := h.bonuses, h.r.spin, h.r.spinData, h.r.currResult, h.r.logEvent, h.r.awardFreeSpins, h.r.prng

	var bonusSymbol = utils.MaxIndex
	oldMultiplier := spin.Multiplier()
	var retestPayouts bool

	var altSymbols bool
	var free uint8

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			t := a.Triggered(spin)
			if t == nil {
				logEvent(a, false)
				continue
			}

			switch t.Result() {
			case slots.BonusGame:
				if h.r.bonusGame != t {
					h.r.bonusGame = t
					if spinData != nil {
						spinData.SetTransition(t.FeatureTransition())
					}
					h.r.makeChoice = h.r.playerChoice && t.PlayerChoice()
				}

			case slots.FreeSpins:
				awardFreeSpins(t.NrOfSpins(prng))
				if t.BonusSymbol() && h.r.BonusSymbol() == utils.MaxIndex {
					// save for last!
					bonusSymbol, altSymbols = h.r.slots.Symbols().GetBonusSymbol(prng), t.AltSymbols()
				}
				h.r.makeChoice = h.r.playerChoice && t.PlayerChoice()

			case slots.Payout:
				a.Payout(spin, currResult)
				if spin.IsExpanded() {
					spinData.SetAfterExpand(spin)
				}

			case slots.HotReel:
				spinData.SetHot(spin)

			case slots.Multiplier:
				spinData.SetMultiplier(spin)
				if m := spin.Multiplier(); m != oldMultiplier {
					if ma, ok := a.(*slots.MultiplierAction); ok {
						free += ma.BonusFreeSpins(m)
					}
				}

			case slots.Sticky:
				spinData.SetSticky(spin)

			case slots.SymbolsInjected:
				spinData.SetAfterInject(spin)
				spinData.SetMultipliers(spin)
				retestPayouts = true

			case slots.GridModified:
				spinData.Update(spin)

			default:
			}

			logEvent(t, true)
		}
	}

	if bonusSymbol != utils.MaxIndex {
		// now it's safe to update this!
		// e.g. BonusSymbol payouts are configured last, so after the BonusSymbol is awarded :)
		spin.SetBonusSymbol(bonusSymbol, altSymbols)
	}

	if retestPayouts {
		currResult.ReleasePayouts()
		h.testRegularPayouts()
		h.testRegularPenalties()
	}

	awardFreeSpins(free)
}

// testClearing tests the clearing actions on the grid result.
func (h *actionHandler) testClearing() {
	list, spin, spinData, logEvent := h.clearing, h.r.spin, h.r.spinData, h.r.logEvent

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) && a.CanClear(spin) {
			if t := a.Triggered(spin); t == nil {
				logEvent(a, false)
			} else {
				switch t.Result() {
				case slots.WildsJumped:
					spinData.SetAfterJump(spin)

				default:
					h.r.needRefill = true
					spinData.SetAfterClear(spin)
					if t.AltSymbols() {
						spin.SetBonusSymbol(utils.MaxIndex, true)
					}
				}

				logEvent(t, true)
			}
		}
	}
}

// testStickiness tests the sticky actions on the grid result.
func (h *actionHandler) testStickiness() {
	list, spin, spinData, symbolsState, logEvent := h.sticky, h.r.spin, h.r.spinData, h.r.symbolsState, h.r.logEvent

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) && a.CanSticky(spin) {
			if t := a.TriggeredWithState(spin, symbolsState); t == nil {
				logEvent(a, false)
			} else {
				switch t.Result() {
				case slots.Sticky:
					spinData.SetSticky(spin)
				case slots.ChooseSticky:
					spinData.SetStickyChoices(spin)
				default:
				}
				logEvent(t, true)
			}
		}
	}
}

// testPreBonus tests actions after the initial spin(s) of a round, and before the free spins/bonusBuy games.
func (h *actionHandler) testPreBonus() {
	list, spin, logEvent := h.preBonus, h.r.spin, h.r.logEvent

	for ix := range list {
		if a := list[ix]; a.CanTrigger(spin) {
			logEvent(a, a.Triggered(spin) != nil)
		}
	}
}

// acquireActionHandler instantiates a new action handler from the memory pool.
func acquireActionHandler(r *Regular, actions slots.SpinActions) *actionHandler {
	h := actionHandlerProducer.Acquire().(*actionHandler)
	h.r = r

	for ix := range actions {
		a := actions[ix]

		switch a.Stage() {
		case slots.PreSpin:
			h.preSpin = append(h.preSpin, a)

		case slots.ReviseGrid:
			h.reviseGrid = append(h.reviseGrid, a)

		case slots.ExpandBefore:
			h.expandBefore = append(h.expandBefore, a.(*slots.WildAction))

		case slots.TestGrid:
			h.gridActions = append(h.gridActions, a)

		case slots.RegularPayouts:
			h.regularPayouts = append(h.regularPayouts, a)
			if p, ok := a.(*slots.PayoutAction); ok {
				if p.HavePaylines() {
					h.r.havePaylines = true
				}
				if p.HaveAllPaylines() {
					h.r.haveAllPaylines = true
				}
				if p.HaveClusterPayouts() {
					h.r.haveClusterPays = true
				}
			}

		case slots.RegularPenalties:
			h.regularPenalties = append(h.regularPenalties, a)

		case slots.Injection:
			h.injections = append(h.injections, a)

		case slots.ExpandAfter:
			h.expandAfter = append(h.expandAfter, a.(*slots.WildAction))

		case slots.TestState:
			h.state = append(h.state, a.(*slots.StateAction))

		case slots.ExtraPayouts:
			h.extraPayouts = append(h.extraPayouts, a)
			if sa, ok := a.(*slots.ScatterAction); ok {
				h.r.haveScatterPays = sa.HaveAllScatterPayouts() || sa.CanPayout()
			}
			if wa, ok := a.(*slots.WildAction); ok {
				h.r.haveScatterPays = wa.CanPayout()
			}

		case slots.AwardBonuses:
			h.bonuses = append(h.bonuses, a)

		case slots.TestClearance:
			h.clearing = append(h.clearing, a)

		case slots.TestStickiness:
			h.sticky = append(h.sticky, a)

		case slots.PaidOnly:
			h.forSale = append(h.forSale, a)

		case slots.PreBonus:
			h.preBonus = append(h.preBonus, a)

		case slots.TestPlayerChoice:
			h.playerChoices = append(h.playerChoices, a.(*slots.ChoiceAction))

		default:
			panic(consts.MsgInvalidActionStage)
		}
	}

	return h
}

type actionHandler struct {
	r                *Regular            // the slot machine game handler.
	forSale          slots.SpinActions   // bonusBuy buy actions.
	preSpin          slots.SpinActions   // list of pre-spin actions.
	reviseGrid       slots.SpinActions   // list of revise grid actions.
	expandBefore     slots.WildActions   // list of expand before actions.
	gridActions      slots.SpinActions   // list of grid actions.
	regularPayouts   slots.SpinActions   // list of regular payout actions.
	regularPenalties slots.SpinActions   // list of penalty actions.
	injections       slots.SpinActions   // list of inject symbol/cluster actions.
	expandAfter      slots.WildActions   // list of expand after actions.
	state            slots.StateActions  // list of state testing actions.
	extraPayouts     slots.SpinActions   // list of extra payout actions.
	bonuses          slots.SpinActions   // list of bonusBuy actions.
	sticky           slots.SpinActions   // list of sticky testing actions.
	clearing         slots.SpinActions   // list of clear testing actions.
	preBonus         slots.SpinActions   // list of actions to execute before bonuses are played.
	playerChoices    slots.ChoiceActions // player choice actions.
	pool.Object
}

// actionHandlerProducer is the memory pool for action handlers.
// Make sure to initialize all slices appropriately!
var actionHandlerProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	h := &actionHandler{
		forSale:          make(slots.SpinActions, 0, 4),
		preSpin:          make(slots.SpinActions, 0, 8),
		reviseGrid:       make(slots.SpinActions, 0, 32),
		expandBefore:     make(slots.WildActions, 0, 8),
		gridActions:      make(slots.SpinActions, 0, 8),
		regularPayouts:   make(slots.SpinActions, 0, 8),
		regularPenalties: make(slots.SpinActions, 0, 8),
		injections:       make(slots.SpinActions, 0, 8),
		expandAfter:      make(slots.WildActions, 0, 8),
		state:            make(slots.StateActions, 0, 8),
		extraPayouts:     make(slots.SpinActions, 0, 8),
		bonuses:          make(slots.SpinActions, 0, 16),
		sticky:           make(slots.SpinActions, 0, 8),
		clearing:         make(slots.SpinActions, 0, 8),
		preBonus:         make(slots.SpinActions, 0, 8),
		playerChoices:    make(slots.ChoiceActions, 0, 8),
	}
	return h, h.reset
})

// reset clear the regular slot machine.
func (h *actionHandler) reset() {
	if h != nil {
		h.forSale = h.forSale[:0]
		h.preSpin = h.preSpin[:0]
		h.reviseGrid = h.reviseGrid[:0]
		h.expandBefore = h.expandBefore[:0]
		h.gridActions = h.gridActions[:0]
		h.regularPayouts = h.regularPayouts[:0]
		h.regularPenalties = h.regularPenalties[:0]
		h.injections = h.injections[:0]
		h.expandAfter = h.expandAfter[:0]
		h.state = h.state[:0]
		h.extraPayouts = h.extraPayouts[:0]
		h.bonuses = h.bonuses[:0]
		h.clearing = h.clearing[:0]
		h.sticky = h.sticky[:0]
		h.preBonus = h.preBonus[:0]
		h.playerChoices = h.playerChoices[:0]
	}
}

type actionHandlers []*actionHandler

func purgeActionHandlers(list actionHandlers, capacity int) actionHandlers {
	if cap(list) < capacity {
		releaseActionHandlers(list)
		return make(actionHandlers, 0, capacity)
	}
	return releaseActionHandlers(list)
}

func releaseActionHandlers(list actionHandlers) actionHandlers {
	if list == nil {
		return nil
	}
	for ix := range list {
		if list[ix] != nil {
			list[ix].Release()
			list[ix] = nil
		}
	}
	return list[:0]
}
