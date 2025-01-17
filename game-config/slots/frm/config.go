package frm

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount        = 10
	reels              = 7
	rows               = 7
	centerTile         = 24
	bonusBuyCost       = 150
	maxPayout          = 10000.0
	freeSpins3   uint8 = 10
	freeSpins4   uint8 = 12
	freeSpins5   uint8 = 15

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

	scatter = id09
	wild    = id10

	bonusBuyID        = 1
	clearProgressID   = 5
	clearFlagsID      = 6
	scatters92ID      = 10
	scatters94ID      = 12
	scatters96ID      = 14
	clusterPayoutsID  = 20
	clusterFreeID     = 30
	cluster20ID       = 40
	cluster40ID       = 41
	cluster80ID       = 42
	cluster160ID      = 43
	wild20ID          = 50
	wild40ID          = 51
	wild80ID          = 52
	progress160usedID = 60
	progress80usedID  = 61
	progress40usedID  = 62
	progress20usedID  = 63
	progressMeterID   = 70
	stickyWildsID     = 75
	freeGames1ID      = 80
	freeGames2ID      = 81
	freeGames3ID      = 82
	progress20ID      = 90
	progress40ID      = 91
	progress80ID      = 92
	progress160ID     = 93
	clearPayoutsID    = 100
	jumpingWildsID    = 101
	resetStickiesID   = 105

	flagBonusBuy = 0
	flagLife1    = 1
	flagLife2    = 2
	flagLife3    = 3
	flagLife4    = 4
)

var (
	// WIP WIP WIP
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 30, 30, 30, 40, 40),
		comp.WithWeights(40, 50, 30, 30, 30, 50, 40),
		comp.WithWeights(40, 20, 30, 30, 30, 20, 40),
		comp.WithWeights(40, 30, 20, 20, 20, 30, 40),
		comp.WithWeights(30, 40, 25, 25, 25, 40, 30),
		comp.WithWeights(30, 40, 25, 25, 25, 40, 30),
		comp.WithWeights(12, 20, 12, 12, 12, 20, 12),
		comp.WithWeights(12, 20, 12, 12, 12, 20, 12),
		comp.WithWeights(0, 0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0, 0),
	}
	scatterChances92 = []float64{35, 90, 90, 5, 1}

	// WIP WIP WIP
	weights94 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 30, 30, 30, 40, 40),
		comp.WithWeights(40, 50, 30, 30, 30, 50, 40),
		comp.WithWeights(40, 20, 30, 30, 30, 20, 40),
		comp.WithWeights(40, 30, 20, 20, 20, 30, 40),
		comp.WithWeights(30, 40, 25, 25, 25, 40, 30),
		comp.WithWeights(30, 40, 25, 25, 25, 40, 30),
		comp.WithWeights(12, 20, 12, 12, 12, 20, 12),
		comp.WithWeights(12, 20, 12, 12, 12, 20, 12),
		comp.WithWeights(0, 0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0, 0),
	}
	scatterChances94 = []float64{35, 90, 90, 5, 1}

	// WIP WIP WIP
	weights96 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 30, 30, 30, 40, 40),
		comp.WithWeights(40, 50, 30, 30, 30, 50, 40),
		comp.WithWeights(40, 20, 30, 30, 30, 20, 40),
		comp.WithWeights(40, 30, 20, 20, 20, 30, 40),
		comp.WithWeights(30, 40, 25, 25, 25, 40, 30),
		comp.WithWeights(30, 40, 25, 25, 25, 40, 30),
		comp.WithWeights(12, 20, 12, 12, 12, 20, 12),
		comp.WithWeights(12, 20, 12, 12, 12, 20, 12),
		comp.WithWeights(0, 0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0, 0),
	}
	scatterChances96 = []float64{35, 90, 90, 5, 1}
)

var (
	reelMask = util.UInt8s{4, 5, 6, 7, 6, 5, 4}

	injectMultipliers20 = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3}, []float64{50, 35, 15})
	injectMultipliers40 = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3}, []float64{40, 35, 25})
	injectMultipliers80 = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3}, []float64{30, 30, 40})

	n01 = comp.WithName("Banana")
	n02 = comp.WithName("Grapes")
	n03 = comp.WithName("Plum")
	n04 = comp.WithName("Pear")
	n05 = comp.WithName("Lemon")
	n06 = comp.WithName("Orange")
	n07 = comp.WithName("Strawberry")
	n08 = comp.WithName("Jelly")
	n09 = comp.WithName("Bonus")
	n10 = comp.WithName("Wild")

	r01 = comp.WithResource("l4")
	r02 = comp.WithResource("l3")
	r03 = comp.WithResource("l2")
	r04 = comp.WithResource("l1")
	r05 = comp.WithResource("h4")
	r06 = comp.WithResource("h3")
	r07 = comp.WithResource("h2")
	r08 = comp.WithResource("h1")
	r09 = comp.WithResource("bonus")
	r10 = comp.WithResource("wild")

	p01 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.15, 0.2, 0.25, 0.3, 0.4, 0.5, 0.6, 0.8, 1, 1.5, 1.5, 1.5, 1.5, 1.5, 4, 4, 4, 4, 4, 7.5, 7.5, 7.5, 7.5, 7.5, 10, 10, 10, 10, 10, 20, 20, 20)
	p02 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.15, 0.2, 0.25, 0.3, 0.4, 0.5, 0.6, 0.8, 1, 1.5, 1.5, 1.5, 1.5, 1.5, 4, 4, 4, 4, 4, 7.5, 7.5, 7.5, 7.5, 7.5, 10, 10, 10, 10, 10, 20, 20, 20)
	p03 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.15, 0.2, 0.25, 0.3, 0.4, 0.5, 0.6, 0.8, 1, 1.5, 1.5, 1.5, 1.5, 1.5, 4, 4, 4, 4, 4, 7.5, 7.5, 7.5, 7.5, 7.5, 10, 10, 10, 10, 10, 20, 20, 20)
	p04 = comp.WithPayouts(0, 0, 0, 0, 0.1, 0.15, 0.2, 0.25, 0.3, 0.4, 0.5, 0.6, 0.8, 1, 1.5, 1.5, 1.5, 1.5, 1.5, 4, 4, 4, 4, 4, 7.5, 7.5, 7.5, 7.5, 7.5, 10, 10, 10, 10, 10, 20, 20, 20)
	p05 = comp.WithPayouts(0, 0, 0, 0, 0.3, 0.4, 0.6, 0.8, 1, 1.25, 1.5, 2, 2.5, 3, 5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 20, 20, 20, 20, 20, 30, 30, 30, 30, 30, 50, 50, 50)
	p06 = comp.WithPayouts(0, 0, 0, 0, 0.4, 0.6, 0.8, 1, 1.25, 1.5, 2, 2.5, 3, 4, 6, 6, 6, 6, 6, 15, 15, 15, 15, 15, 25, 25, 25, 25, 25, 40, 40, 40, 40, 40, 75, 75, 75)
	p07 = comp.WithPayouts(0, 0, 0, 0, 0.5, 0.75, 1, 1.25, 1.5, 2, 2.5, 3, 4, 5, 7.5, 7.5, 7.5, 7.5, 7.5, 20, 20, 20, 20, 20, 30, 30, 30, 30, 30, 50, 50, 50, 50, 50, 100, 100, 100)
	p08 = comp.WithPayouts(0, 0, 0, 0, 1, 1.5, 2, 2.5, 3, 4, 5, 6, 8, 10, 15, 15, 15, 15, 15, 40, 40, 40, 40, 40, 50, 50, 50, 50, 50, 100, 100, 100, 100, 100, 200, 200, 200)
	p09 = comp.WithPayouts()
	p10 = p09

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10}

	flag0 = comp.NewRoundFlag(flagBonusBuy, "bonus buy")
	flag1 = comp.NewRoundFlag(flagLife1, "free life 1").WithExport()
	flag2 = comp.NewRoundFlag(flagLife2, "free life 2").WithExport()
	flag3 = comp.NewRoundFlag(flagLife3, "free life 3").WithExport()
	flag4 = comp.NewRoundFlag(flagLife4, "free life 4").WithExport()
	flags = comp.RoundFlags{flag0, flag1, flag2, flag3, flag4}
)

var (
	symbols *comp.SymbolSet

	actions92all     comp.SpinActions
	actions92first   comp.SpinActions
	actions92free    comp.SpinActions
	actions92firstBB comp.SpinActions
	actions92freeBB  comp.SpinActions
	actions94all     comp.SpinActions
	actions94first   comp.SpinActions
	actions94free    comp.SpinActions
	actions94firstBB comp.SpinActions
	actions94freeBB  comp.SpinActions
	actions96all     comp.SpinActions
	actions96first   comp.SpinActions
	actions96free    comp.SpinActions
	actions96firstBB comp.SpinActions
	actions96freeBB  comp.SpinActions

	slots92 *comp.Slots
	slots94 *comp.Slots
	slots96 *comp.Slots

	slots92params game.RegularParams
	slots94params game.RegularParams
	slots96params game.RegularParams
)

func initActions() {
	// bonus buy
	bonusBuy := comp.NewPaidAction(comp.FreeSpins, freeSpins5, bonusBuyCost, scatter, 5).WithFlag(flagBonusBuy, bonusBuyID)
	bonusBuy.Describe(bonusBuyID, "bonus buy")

	// clear progress meter and free lives on free spin.
	clearProgress := comp.NewResetProgressAction(0)
	clearProgress.WithTriggerFilters(comp.OnFreeSpin) // leave this in place as it must not be called doing refills!!!
	clearProgress.Describe(clearProgressID, "clear progress meter - free spin")

	clearFlags := comp.NewRoundFlagsReset(flagLife1, flagLife2, flagLife3, flagLife4)
	clearFlags.WithTriggerFilters(comp.OnFreeSpin) // leave this in place as it must not be called doing refills!!!
	clearFlags.Describe(clearFlagsID, "clear free lives - free spin")

	// generate scatter symbols RTP 92.
	scatters92 := comp.NewGenerateSymbolAction(scatter, scatterChances92).GenerateNoDupes()
	scatters92.Describe(scatters92ID, "generate scatters - first spin - RTP 92")

	// generate scatter symbols RTP 94.
	scatters94 := comp.NewGenerateSymbolAction(scatter, scatterChances94).GenerateNoDupes()
	scatters94.Describe(scatters94ID, "generate scatters - first spin - RTP 94")

	// generate scatter symbols RTP 96.
	scatters96 := comp.NewGenerateSymbolAction(scatter, scatterChances96).GenerateNoDupes()
	scatters96.Describe(scatters96ID, "generate scatters - first spin - RTP 96")

	// award cluster payouts.
	clusterPays := comp.NewClusterPayoutsAction(reels, rows, comp.ClusterGridMask(reelMask, comp.Hexagonal))
	clusterPays.Describe(clusterPayoutsID, "cluster payouts")

	// inject cluster on every free spin for the bonus buy.
	clusterFree := comp.NewClusterOffsetInject(5, 5, centerTile)
	clusterFree.Reschedule(comp.TestGrid)
	clusterFree.WithTriggerFilters(comp.OnFreeSpin) // leave this in place as it must not be called doing refills!!!
	clusterFree.Describe(clusterFreeID, "inject cluster - free spin")

	// inject cluster on zero payouts.
	cluster20 := comp.NewClusterOffsetInject(5, 5, centerTile)
	cluster20.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife1, 1))
	cluster20.Describe(cluster20ID, "inject cluster - free life 1")
	cluster40 := comp.NewClusterOffsetInject(5, 10, centerTile)
	cluster40.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife2, 1), comp.OnRoundFlagValue(flagLife1, 2))
	cluster40.Describe(cluster40ID, "inject cluster - free life 2")
	cluster80 := comp.NewClusterOffsetInject(10, 15, centerTile)
	cluster80.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife3, 1), comp.OnRoundFlagValue(flagLife2, 2))
	cluster80.Describe(cluster80ID, "inject cluster - free life 3")
	cluster160 := comp.NewClusterOffsetInject(15, 20, centerTile)
	cluster160.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife4, 1), comp.OnRoundFlagValue(flagLife3, 2))
	cluster160.Describe(cluster160ID, "inject cluster - free life 4")

	// inject wild on zero payouts.
	wild20 := comp.NewSymbolInjectFromEdge(wild, 2).WithMultipliers(injectMultipliers20)
	wild20.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife1, 1))
	wild20.Describe(wild20ID, "inject wild - free life 1")
	wild40 := comp.NewSymbolInjectFromEdge(wild, 2).WithMultipliers(injectMultipliers40)
	wild40.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife2, 1), comp.OnRoundFlagValue(flagLife1, 2))
	wild40.Describe(wild40ID, "inject wild - free life 2")
	wild80 := comp.NewSymbolInjectFromEdge(wild, 2).WithMultipliers(injectMultipliers80)
	wild80.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife3, 1), comp.OnRoundFlagValue(flagLife2, 2))
	wild80.Describe(wild80ID, "inject wild - free life 3")

	// update free life flags when used.
	// Note that we must do it in reverse order, or they will cascade :)
	progress160used := comp.NewRoundFlagIncreaseAction(flagLife4)
	progress160used.Reschedule(comp.Injection)
	progress160used.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife4, 1), comp.OnRoundFlagValue(flagLife3, 2))
	progress160used.Describe(progress160usedID, "mark free life used - free life 4")
	progress80used := comp.NewRoundFlagIncreaseAction(flagLife3)
	progress80used.Reschedule(comp.Injection)
	progress80used.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife3, 1), comp.OnRoundFlagValue(flagLife2, 2))
	progress80used.Describe(progress80usedID, "mark free life used - free life 3")
	progress40used := comp.NewRoundFlagIncreaseAction(flagLife2)
	progress40used.Reschedule(comp.Injection)
	progress40used.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife2, 1), comp.OnRoundFlagValue(flagLife1, 2))
	progress40used.Describe(progress40usedID, "mark free life used - free life 2")
	progress20used := comp.NewRoundFlagIncreaseAction(flagLife1)
	progress20used.Reschedule(comp.Injection)
	progress20used.WithTriggerFilters(comp.OnZeroPayouts(), comp.OnRoundFlagValue(flagLife1, 1))
	progress20used.Describe(progress20usedID, "mark free life used - free life 1")

	// cluster payouts progress meter.
	progressMeter := comp.NewPayoutSymbolProgress(0, 160)
	progressMeter.Describe(progressMeterID, "cluster payouts progress meter")

	// make sure they stick before cascading reels!
	stickyWilds := comp.NewStickySymbolAction(wild)
	stickyWilds.Reschedule(comp.AwardBonuses)
	stickyWilds.Describe(stickyWildsID, "sticky wilds")

	// award free spins.
	awardSpins3 := comp.NewScatterFreeSpinsAction(freeSpins3, false, scatter, 3, false)
	awardSpins3.WithTriggerFilters(comp.OnFirstSpin)
	awardSpins3.Describe(freeGames1ID, "award 3 free spins")
	awardSpins4 := comp.NewScatterFreeSpinsAction(freeSpins4, false, scatter, 4, false).WithAlternate(awardSpins3)
	awardSpins4.WithTriggerFilters(comp.OnFirstSpin)
	awardSpins4.Describe(freeGames2ID, "award 4 free spins")
	awardSpins5 := comp.NewScatterFreeSpinsAction(freeSpins5, false, scatter, 5, false).WithAlternate(awardSpins4)
	awardSpins5.WithTriggerFilters(comp.OnFirstSpin)
	awardSpins5.Describe(freeGames3ID, "award 5 free spins")

	// award free life on progress meter levels.
	progress20 := comp.NewRoundFlagIncreaseAction(flagLife1)
	progress20.Reschedule(comp.AwardBonuses)
	progress20.WithTriggerFilters(comp.OnRoundFlagValue(flagLife1, 0), comp.OnProgressLevelAbove(19))
	progress20.Describe(progress20ID, "progress level 20 - free life 1")
	progress40 := comp.NewRoundFlagIncreaseAction(flagLife2)
	progress40.Reschedule(comp.AwardBonuses)
	progress40.WithTriggerFilters(comp.OnRoundFlagValue(flagLife2, 0), comp.OnProgressLevelAbove(39))
	progress40.Describe(progress40ID, "progress level 40 - free life 2")
	progress80 := comp.NewRoundFlagIncreaseAction(flagLife3)
	progress80.Reschedule(comp.AwardBonuses)
	progress80.WithTriggerFilters(comp.OnRoundFlagValue(flagLife3, 0), comp.OnProgressLevelAbove(79))
	progress80.Describe(progress80ID, "progress level 80 - free life 3")
	progress160 := comp.NewRoundFlagIncreaseAction(flagLife4)
	progress160.Reschedule(comp.AwardBonuses)
	progress160.WithTriggerFilters(comp.OnRoundFlagValue(flagLife4, 0), comp.OnProgressLevelAbove(159))
	progress160.Describe(progress160ID, "progress level 160 - free life 4")

	// cascading reels.
	clearPayouts := comp.NewClearPayoutsAction()
	clearPayouts.Describe(clearPayoutsID, "clear cluster payouts")

	// jumping wilds.
	jumpingWilds := comp.NewJumpingWilds(comp.GridAny, wild)
	jumpingWilds.Describe(jumpingWildsID, "jumping wilds")

	// reset sticky indicators on no payout.
	resetStickies := comp.NewResetStickiesAction()
	resetStickies.WithTriggerFilters(comp.OnZeroPayouts())
	resetStickies.Describe(resetStickiesID, "reset sticky indicators - no payouts")

	actionsAall := comp.SpinActions{bonusBuy, clearProgress, clearFlags}
	actionsAfirst := comp.SpinActions{bonusBuy}
	actionsAfree := comp.SpinActions{bonusBuy, clearProgress, clearFlags}

	actionsBall := comp.SpinActions{clusterPays, clusterFree, cluster20, cluster40, cluster80, cluster160,
		wild20, wild40, wild80, progress160used, progress80used, progress40used, progress20used,
		progressMeter, stickyWilds, awardSpins5, progress20, progress40, progress80, progress160,
		clearPayouts, jumpingWilds, resetStickies}
	actionsBfirst := comp.SpinActions{clusterPays, cluster20, cluster40, cluster80, cluster160,
		wild20, wild40, wild80, progress160used, progress80used, progress40used, progress20used,
		progressMeter, stickyWilds, awardSpins5, progress20, progress40, progress80, progress160,
		clearPayouts, jumpingWilds, resetStickies}
	actionsBfree := comp.SpinActions{clusterPays, clusterFree, cluster20, cluster40, cluster80, cluster160,
		wild20, wild40, wild80, progress160used, progress80used, progress40used, progress20used,
		progressMeter, stickyWilds, progress20, progress40, progress80, progress160,
		clearPayouts, jumpingWilds, resetStickies}

	actions92all = append(append(actionsAall, scatters92), actionsBall...)
	actions92first = append(append(actionsAfirst, scatters92), actionsBfirst...)
	actions92free = append(actionsAfree, actionsBfree...)
	actions92firstBB = append(append(actionsAfirst, scatters92), actionsBfirst...)
	actions92freeBB = append(actionsAfree, actionsBfree...)

	actions94all = append(append(actionsAall, scatters94), actionsBall...)
	actions94first = append(append(actionsAfirst, scatters94), actionsBfirst...)
	actions94free = append(actionsAfree, actionsBfree...)
	actions94firstBB = append(append(actionsAfirst, scatters94), actionsBfirst...)
	actions94freeBB = append(actionsAfree, actionsBfree...)

	actions96all = append(append(actionsAall, scatters96), actionsBall...)
	actions96first = append(append(actionsAfirst, scatters96), actionsBfirst...)
	actions96free = append(actionsAfree, actionsBfree...)
	actions96firstBB = append(append(actionsAfirst, scatters96), actionsBfirst...)
	actions96freeBB = append(actionsAfree, actionsBfree...)
}

func initSlots(target float64, weights [symbolCount]comp.SymbolOption, actions1, actions2, actions3, actions4 []comp.SpinActioner) *comp.Slots {
	ss := make([]*comp.Symbol, symbolCount)
	for ix := range ss {
		id := ids[ix]
		switch id {
		case scatter:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Scatter))
		case wild:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Wild))
		default:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix])
		}
	}
	s := comp.NewSymbolSet(ss...)

	if symbols == nil {
		symbols = s
	}

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.WithMask(reelMask...),
		comp.WithSymbols(s),
		comp.CascadingReels(true),
		comp.WithProgressMeter(),
		comp.WithBonusBuy(),
		comp.MaxPayout(maxPayout),
		comp.WithRTP(target),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, actions3, actions4),
	)
}

func init() {
	initActions()

	slots92 = initSlots(92, weights92, actions92first, actions92free, actions92firstBB, actions92freeBB)
	slots94 = initSlots(94, weights94, actions94first, actions94free, actions94firstBB, actions94freeBB)
	slots96 = initSlots(96, weights96, actions96first, actions96free, actions96firstBB, actions96freeBB)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}
