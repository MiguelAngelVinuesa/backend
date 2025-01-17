package crw

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount = 6
	reels       = 5
	rows        = 1
	direction   = comp.PayLTR
	maxPayout   = 1000.0

	id01 = 1
	id02 = 2
	id03 = 3
	id04 = 4
	id05 = 5
	id06 = 6

	scatter = id04
	poison  = id05
	skull   = id06

	scatterFirstID   = 1
	fullbonusFree1ID = 11
	injectPoisonID   = 21
	poisonFreeID     = 22
	injectSkullID    = 25
	skullFreeID      = 26
	regPayoutsID     = 31
	scatterGames10ID = 41
	freeSpinsID      = 91

	flagFreeSpins = 0
)

var (
	weights96a = [symbolCount]comp.SymbolOption{
		comp.WithWeights(155, 95, 155, 95, 155),
		comp.WithWeights(75, 85, 70, 85, 75),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}

	firstScatterWeights = []float64{50, 27.5, 7.5, 100, 100}

	poisonWeights = []float64{50}
	skullWeights  = []float64{59.2}
)

var (
	// symbol names.
	n01 = comp.WithName("Cherries")
	n02 = comp.WithName("BAR")
	n03 = comp.WithName("7")
	n04 = comp.WithName("Gift")
	n05 = comp.WithName("Poison")
	n06 = comp.WithName("Skull")

	r01 = comp.WithResource("h3")
	r02 = comp.WithResource("h2")
	r03 = comp.WithResource("h1")
	r04 = comp.WithResource("scatter")
	r05 = comp.WithResource("poison")
	r06 = comp.WithResource("skull")

	// paytable.
	p01 = comp.WithPayouts(0, 0, 0.5, 5, 10)
	p02 = comp.WithPayouts(0, 0, 5, 10, 25)
	p03 = comp.WithPayouts(0, 0, 50, 100, 1000)
	p04 = comp.WithPayouts()
	p05 = p04
	p06 = p04

	// pay-lines.
	p00000 = comp.NewPayline(id01, rows, 0, 0, 0, 0, 0)

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06}
	paylines  = comp.Paylines{p00000}

	flag0 = comp.NewRoundFlag(flagFreeSpins, "free spin count")
	flags = comp.RoundFlags{flag0}
)

var (
	symbols1 *comp.SymbolSet

	actions96first comp.SpinActions
	actions96free  comp.SpinActions
	actions96all   comp.SpinActions

	slots96 *comp.Slots

	slots96params game.RegularParams
)

func initActions() {
	// generate scatter symbols on first spin.
	scatterFirst := comp.NewGenerateSymbolAction(scatter, firstScatterWeights)
	scatterFirst.Describe(scatterFirstID, "insert scatters")

	// inject poison in free spins (-10x).
	injectPoison := comp.NewGenerateSymbolAction(poison, poisonWeights)
	injectPoison.WithTriggerFilters(comp.OnSpinSequenceAbove(2))
	injectPoison.Describe(injectPoisonID, "inject poison")

	// inject skull in free spins (-50%).
	injectSkull := comp.NewGenerateSymbolAction(skull, skullWeights)
	injectSkull.WithTriggerFilters(comp.OnSpinSequenceAbove(2))
	injectSkull.Describe(injectSkullID, "inject skull")

	// award regular payouts.
	regPayouts := comp.NewPaylinesAction()
	regPayouts.Describe(regPayoutsID, "award regular payouts")

	// impose -10x penalty for poison.
	poisonFree := comp.NewReductionAction(poison, 1, 10.0)
	poisonFree.Describe(poisonFreeID, "impose 10x penalty (poison)")

	// impose -50% penalty for skull.
	skullFree := comp.NewDivisionAction(skull, 1, 2)
	skullFree.Describe(skullFreeID, "impose 50% penalty (skull)")

	// award free spins from scatters.
	scatterGames10 := comp.NewScatterFreeSpinsAction(10, false, scatter, 5, false)
	scatterGames10.Describe(scatterGames10ID, "award 10 free spins from 5 scatters")

	// update round flag 3 marking sequence of free spin.
	freeSpinsFlag := comp.NewRoundFlagIncreaseAction(flagFreeSpins)
	freeSpinsFlag.Describe(freeSpinsID, "count number of free spins (flag 0)")

	// award max bonus in first free spin.
	bonusGrid := comp.GridOffsets{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}}
	bonusCenters := comp.GridOffsets{{0, 0}}
	bonusWeights := util.AcquireWeighting().AddWeights(util.Indexes{id03}, []float64{100})
	fullbonusFree1 := comp.NewGenerateShapeAction(100, bonusGrid, bonusCenters, bonusWeights).AllowPrevious()
	fullbonusFree1.WithTriggerFilters(comp.OnSpinSequence(2))
	fullbonusFree1.Describe(fullbonusFree1ID, "award full bonus on first free spin")

	actionsAall := comp.SpinActions{scatterFirst, fullbonusFree1, injectPoison, injectSkull}
	actionsAfirst := comp.SpinActions{scatterFirst}
	actionsAfree := comp.SpinActions{fullbonusFree1, injectPoison, injectSkull}

	actionsBall := comp.SpinActions{regPayouts, poisonFree, skullFree, scatterGames10, freeSpinsFlag}
	actionsBfirst := comp.SpinActions{regPayouts, scatterGames10}
	actionsBfree := comp.SpinActions{regPayouts, poisonFree, skullFree, freeSpinsFlag}

	actions96all = append(actionsAall, actionsBall...)
	actions96first = append(actionsAfirst, actionsBfirst...)
	actions96free = append(actionsAfree, actionsBfree...)
}

func initSlots(target float64, weights1 [symbolCount]comp.SymbolOption, actions1, actions2 []comp.SpinActioner) *comp.Slots {
	ss := make([]*comp.Symbol, symbolCount)
	for ix := range ss {
		if id := ids[ix]; id == scatter {
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights1[ix], comp.WithKind(comp.Scatter))
		} else {
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights1[ix])
		}
	}
	s1 := comp.NewSymbolSet(ss...)

	if symbols1 == nil {
		symbols1 = s1
	}

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.WithSymbols(s1),
		comp.WithPaylines(direction, false, paylines...),
		comp.MaxPayout(maxPayout),
		comp.WithReverseWin(),
		comp.WithRTP(target),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, nil, nil),
	)
}

func init() {
	initActions()

	slots96 = initSlots(96.0, weights96a, actions96first, actions96free)
	slots96params = game.RegularParams{Slots: slots96}
}
