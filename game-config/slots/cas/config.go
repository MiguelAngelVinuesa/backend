package cas

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount = 17
	reels       = 5
	rows        = 6
	direction   = comp.PayLTR
	maxPayout   = 1000.0

	id01 = 1
	id02 = 2
	id03 = 3
	id04 = 4
	id05 = 5
	id06 = 6
	id07 = 7
	id08 = 8
	id09 = 9
	id10 = 10
	id11 = 11
	id12 = 12
	id13 = 13
	id14 = 14
	id15 = 15
	id16 = 16
	id17 = 17

	security = id11
	camera   = id12
	traitor  = id13
	smoke    = id14
	alarm    = id15
	bonus    = id16
	wild     = id17

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
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(30, 15, 40, 20, 35),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}

	firstScatterWeights = []float64{50, 27.5, 7.5, 100, 100}

	poisonWeights = []float64{50}
	skullWeights  = []float64{63.3}
)

var (
	// symbol names.
	n01 = comp.WithName("Cards")
	n02 = comp.WithName("Chips")
	n03 = comp.WithName("Dice")
	n04 = comp.WithName("Cocktail")
	n05 = comp.WithName("Money")
	n06 = comp.WithName("Driver")
	n07 = comp.WithName("Pyroman")
	n08 = comp.WithName("Trickster")
	n09 = comp.WithName("Hacker")
	n10 = comp.WithName("Leader")
	n11 = comp.WithName("Alarm")
	n12 = comp.WithName("Smoke")
	n13 = comp.WithName("Traitor")
	n14 = comp.WithName("Camera")
	n15 = comp.WithName("Security")
	n16 = comp.WithName("Bonus")
	n17 = comp.WithName("Wild")

	r01 = comp.WithResource("l5")
	r02 = comp.WithResource("l4")
	r03 = comp.WithResource("l3")
	r04 = comp.WithResource("l2")
	r05 = comp.WithResource("l1")
	r06 = comp.WithResource("h5")
	r07 = comp.WithResource("h4")
	r08 = comp.WithResource("h3")
	r09 = comp.WithResource("h2")
	r10 = comp.WithResource("h1")
	r11 = comp.WithResource("n5")
	r12 = comp.WithResource("n4")
	r13 = comp.WithResource("n3")
	r14 = comp.WithResource("n2")
	r15 = comp.WithResource("n1")
	r16 = comp.WithResource("bonus")
	r17 = comp.WithResource("wild")

	// paytable.
	p01 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.2, 0.3, 0.5, 1, 2, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10)
	p02 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.2, 0.3, 0.5, 1, 3, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15)
	p03 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.2, 0.3, 0.5, 1, 5, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20)
	p04 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.2, 0.3, 0.5, 1, 7.5, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25)
	p05 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.2, 0.3, 0.5, 1, 10, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30, 30)
	p06 = comp.WithPayouts(0, 0, 0, 0, 0.2, 0.5, 0.8, 1, 2, 15, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 35, 500)
	p07 = comp.WithPayouts(0, 0, 0, 0, 0.2, 0.5, 0.8, 1, 2, 20, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 40, 1000)
	p08 = comp.WithPayouts(0, 0, 0, 0, 0.4, 0.8, 1, 1.5, 3, 25, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 45, 5000)
	p09 = comp.WithPayouts(0, 0, 0, 0, 0.5, 1, 1.5, 2, 4, 30, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 50, 25000)
	p10 = comp.WithPayouts(0, 0, 0, 0, 1, 1.5, 2.5, 4, 5, 50, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 50000)
	p11 = comp.WithPayouts()
	p12 = p11
	p13 = p11
	p14 = p11
	p15 = p11
	p16 = p11
	p17 = p11

	// pay-lines.
	p00000 = comp.NewPayline(id01, rows, 0, 0, 0, 0, 0)

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12, id13, id14, id15, id16, id17}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12, n13, n14, n15, n16, n17}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11, r12, r13, r14, r15, r16, r17}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11, p12, p13, p14, p15, p16, p17}
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
	scatterFirst := comp.NewGenerateSymbolAction(bonus, firstScatterWeights)
	scatterFirst.Describe(scatterFirstID, "insert scatters")

	// inject poison in free spins (-10x).
	injectPoison := comp.NewGenerateSymbolAction(security, poisonWeights)
	injectPoison.WithTriggerFilters(comp.OnSpinSequenceAbove(2))
	injectPoison.Describe(injectPoisonID, "inject poison")

	// inject skull in free spins (-50%).
	injectSkull := comp.NewGenerateSymbolAction(camera, skullWeights)
	injectSkull.WithTriggerFilters(comp.OnSpinSequenceAbove(2))
	injectSkull.Describe(injectSkullID, "inject skull")

	// award regular payouts.
	regPayouts := comp.NewPaylinesAction()
	regPayouts.Describe(regPayoutsID, "award regular payouts")

	// impose -10x penalty for poison.
	poisonFree := comp.NewReductionAction(security, 1, 10.0)
	poisonFree.Describe(poisonFreeID, "impose 10x penalty (poison)")

	// impose -50% penalty for skull.
	skullFree := comp.NewDivisionAction(camera, 1, 2)
	skullFree.Describe(skullFreeID, "impose 50% penalty (skull)")

	// award free spins from scatters.
	scatterGames10 := comp.NewScatterFreeSpinsAction(10, false, bonus, 5, false)
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
		if id := ids[ix]; id == bonus {
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
