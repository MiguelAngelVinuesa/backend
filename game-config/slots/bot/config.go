package bot

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount  = 10
	symbolCount2 = symbolCount + 1
	reels        = 5
	rows         = 3
	direction    = comp.PayLTR
	noRepeat     = rows - 1
	maxPayout    = 20000.0

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

	redBook    = id10
	yellowBook = id11

	// action ids.
	// NOTE! each must be unique.
	// **DO NOT REORDER** - they are stored in game results in DB - use higher number for new actions.
	scatterFirstID       = 1
	scatterFree1ID       = 2
	scatterFree2ID       = 3
	scatterFree3ID       = 4
	scatterFree4ID       = 5
	scatterFree5ID       = 6
	bonusPage1ID         = 11
	bonusPage2ID         = 12
	bonusPage3ID         = 13
	bonusPage4ID         = 14
	bonusPage5ID         = 15
	redYellow92ID        = 20
	redYellow94ID        = 21
	redYellow96ID        = 22
	yellowFree1ID        = 30
	yellowFree2ID        = 31
	yellowFree3ID        = 32
	yellowFree1bID       = 35
	yellowFree2bID       = 36
	yellowFree3bID       = 37
	removeBonus1ID       = 40
	removeBonus2ID       = 41
	removeBonus3ID       = 42
	removeBonus4ID       = 43
	removeBonus5ID       = 44
	regPayoutsID         = 50
	removePayoutsFirstID = 51
	removePayoutsFree1ID = 52
	removePayoutsFree2ID = 53
	scatterGames8ID      = 60
	scatterGames10ID     = 61
	scatterGames12ID     = 62
	scatterPayouts3ID    = 65
	scatterPayouts4ID    = 66
	scatterPayouts5ID    = 67
	hotReelsID           = 70
	bonusPayoutsID       = 75
	payoutBandsID        = 90
	yellowBookID         = 91
	bonusPageID          = 92
	freeSpinsID          = 93

	flagPayoutBand = 0
	flagYellowBook = 1
	flagBonusPage  = 2
	flagFreeSpins  = 3
)

var (
	// symbol reel weights for first spins.
	weights92a = [symbolCount2]comp.SymbolOption{
		comp.WithWeights(55, 45, 55, 45, 55),
		comp.WithWeights(45, 55, 50, 55, 45),
		comp.WithWeights(50, 20, 50, 20, 50),
		comp.WithWeights(45, 25, 40, 25, 45),
		comp.WithWeights(35, 45, 20, 45, 35),
		comp.WithWeights(18, 28, 21, 28, 18),
		comp.WithWeights(25, 17, 28, 17, 25),
		comp.WithWeights(17, 14, 18, 14, 17),
		comp.WithWeights(10, 14, 12, 14, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	// symbol reel weights for free spins.
	weights92b = [symbolCount2]comp.SymbolOption{
		comp.WithWeights(70, 50, 70, 50, 70),
		comp.WithWeights(50, 70, 50, 70, 50),
		comp.WithWeights(50, 40, 70, 40, 50),
		comp.WithWeights(60, 40, 60, 40, 60),
		comp.WithWeights(40, 60, 40, 60, 40),
		comp.WithWeights(30, 35, 20, 35, 30),
		comp.WithWeights(25, 20, 40, 20, 25),
		comp.WithWeights(17, 12, 10, 12, 17),
		comp.WithWeights(12, 8, 11, 8, 12),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	// weights for selecting bonus symbol.
	weightsBonus92 = []float64{55, 55, 55, 50, 50, 30, 30, 20, 12}
	// chance of red to yellow book morphing.
	yellowChance92 = 12.35

	// symbol reel weights for first spins.
	weights94a = [symbolCount2]comp.SymbolOption{
		comp.WithWeights(55, 45, 55, 45, 55),
		comp.WithWeights(45, 55, 50, 55, 45),
		comp.WithWeights(50, 20, 50, 20, 50),
		comp.WithWeights(45, 25, 40, 25, 45),
		comp.WithWeights(35, 45, 20, 45, 35),
		comp.WithWeights(18, 28, 21, 28, 18),
		comp.WithWeights(25, 17, 28, 17, 25),
		comp.WithWeights(17, 14, 18, 14, 17),
		comp.WithWeights(10, 14, 12, 14, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	// symbol reel weights for free spins.
	weights94b = [symbolCount2]comp.SymbolOption{
		comp.WithWeights(70, 50, 70, 50, 70),
		comp.WithWeights(50, 70, 50, 70, 50),
		comp.WithWeights(50, 40, 70, 40, 50),
		comp.WithWeights(60, 40, 60, 40, 60),
		comp.WithWeights(40, 60, 40, 60, 40),
		comp.WithWeights(30, 35, 20, 35, 30),
		comp.WithWeights(25, 20, 40, 20, 25),
		comp.WithWeights(17, 12, 10, 12, 17),
		comp.WithWeights(13, 8, 12, 8, 13),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	// weights for selecting bonus symbol.
	weightsBonus94 = []float64{55, 55, 55, 50, 50, 30, 30, 20, 12}
	// chance of red to yellow book morphing.
	yellowChance94 = 14.9

	// symbol reel weights for first spins.
	weights96a = [symbolCount2]comp.SymbolOption{
		comp.WithWeights(55, 45, 55, 45, 55),
		comp.WithWeights(45, 55, 50, 55, 45),
		comp.WithWeights(50, 20, 50, 20, 50),
		comp.WithWeights(45, 25, 40, 25, 45),
		comp.WithWeights(35, 45, 20, 45, 35),
		comp.WithWeights(18, 28, 21, 28, 18),
		comp.WithWeights(25, 17, 28, 17, 25),
		comp.WithWeights(17, 14, 18, 14, 17),
		comp.WithWeights(10, 14, 12, 14, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	// symbol reel weights for free spins.
	weights96b = [symbolCount2]comp.SymbolOption{
		comp.WithWeights(70, 50, 70, 50, 70),
		comp.WithWeights(50, 70, 50, 70, 50),
		comp.WithWeights(50, 40, 70, 40, 50),
		comp.WithWeights(60, 40, 60, 40, 60),
		comp.WithWeights(40, 60, 40, 60, 40),
		comp.WithWeights(30, 35, 20, 35, 30),
		comp.WithWeights(25, 20, 40, 20, 25),
		comp.WithWeights(17, 12, 10, 12, 17),
		comp.WithWeights(13, 8, 12, 8, 13),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	// weights for selecting bonus symbol.
	weightsBonus96 = []float64{55, 55, 55, 50, 50, 30, 30, 20, 12}
	// chance of red to yellow book morphing.
	yellowChance96 = 16.1
)

var (
	// weights to generate scatters on first spins.
	firstScatterWeights = []float64{50, 27.5, 4.4, 1, 0.5}

	// weights to determine payout bands flag on each round.
	freeSpinBands = util.AcquireWeighting().AddWeights([]util.Index{1, 2, 3, 4, 5}, []float64{130, 90, 20, 30, 30})

	// weights to generate scatters on free spins.
	freeScatterWeights1 = []float64{35, 15, 5}         // free spin band 1 - low
	freeScatterWeights2 = []float64{44.5, 22.5, 14}    // free spin band 2 - medium/low
	freeScatterWeights3 = []float64{48, 32, 22, 2}     // free spin band 3 - medium
	freeScatterWeights4 = []float64{53, 36, 27, 5, 1}  // free spin band 4 - medium/high
	freeScatterWeights5 = []float64{65, 44, 30, 10, 3} // free spin band 5 - high

	// weights to remove payouts based on total payout factor during first spins.
	removeBandsFirst = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 2.5, RemoveChance: 70.2},
		{MinPayout: 2.5, MaxPayout: 4, RemoveChance: 50},
		{MinPayout: 4, MaxPayout: 30, RemoveChance: 25},
		{MinPayout: 30, MaxPayout: 40, RemoveChance: 32},
		{MinPayout: 40, MaxPayout: 125, RemoveChance: 85},
		{MinPayout: 125, MaxPayout: 200, RemoveChance: 20.5},
		{MinPayout: 200, MaxPayout: 250, RemoveChance: 75},
		{MinPayout: 250, MaxPayout: maxPayout, RemoveChance: 50},
	}

	// weights to remove payouts based on total payout factor during first few free spins.
	removeBandsFree1 = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 30, RemoveChance: 25},
		{MinPayout: 30, MaxPayout: maxPayout, RemoveChance: 75},
	}

	// weights to remove payouts based on total payout factor during later free spins.
	removeBandsFree2 = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 10, RemoveChance: 17.5},
		{MinPayout: 10, MaxPayout: 30, RemoveChance: 12},
		{MinPayout: 30, MaxPayout: 100, RemoveChance: 25},
		{MinPayout: 100, MaxPayout: 250, RemoveChance: 10},
		{MinPayout: 250, MaxPayout: maxPayout, RemoveChance: 1},
	}

	// chance for removing bonus payouts.
	removeBonusChance1 = 75.0
	removeBonusChance2 = 37.5
	removeBonusChance3 = 11.0
	removeBonusChance4 = 5.5
	removeBonusChance5 = 1.0

	// weights to force yellow book.
	forceYellowBook1a = []float64{7}
	forceYellowBook1b = []float64{15}
	forceYellowBook1c = []float64{17}
	forceYellowBook2a = []float64{5}
	forceYellowBook2b = []float64{16}
	forceYellowBook2c = []float64{17}

	// chance to generate full page of bonus symbols during free spins.
	genFullBonusChances1 = []float64{0, 4, 4, 4, 2, 2, 0.5, 0.5, 0.1, 0.1}
	genFullBonusChances2 = []float64{0, 8, 8, 8, 5, 5, 0.7, 0.7, 0.2, 0.2}
	genFullBonusChances3 = []float64{0, 10, 10, 10, 6, 6, 1.2, 1.2, 0.3, 0.3}
	genFullBonusChances4 = []float64{0, 13, 13, 13, 9, 9, 2.2, 2.2, 0.7, 0.7}
	genFullBonusChances5 = []float64{0, 16, 16, 16, 11, 11, 3, 3, 1.5, 1.5}
)

var (
	scatterPayouts = []float64{0, 0, 2, 20, 200}

	// symbol names.
	n01 = comp.WithName("Ten")
	n02 = comp.WithName("Jack")
	n03 = comp.WithName("Queen")
	n04 = comp.WithName("King")
	n05 = comp.WithName("Ace")
	n06 = comp.WithName("Jar")
	n07 = comp.WithName("Ankh")
	n08 = comp.WithName("Sphinx")
	n09 = comp.WithName("Explorer")
	n10 = comp.WithName("Red Book")
	n11 = comp.WithName("Yellow Book")

	r01 = comp.WithResource("l5")
	r02 = comp.WithResource("l4")
	r03 = comp.WithResource("l3")
	r04 = comp.WithResource("l2")
	r05 = comp.WithResource("l1")
	r06 = comp.WithResource("h4")
	r07 = comp.WithResource("h3")
	r08 = comp.WithResource("h2")
	r09 = comp.WithResource("h1")
	r10 = comp.WithResource("red")
	r11 = comp.WithResource("yellow")

	// paytable.
	p01 = comp.WithPayouts(0, 0, 0.5, 2.5, 10)
	p02 = comp.WithPayouts(0, 0, 0.5, 2.5, 10)
	p03 = comp.WithPayouts(0, 0, 0.5, 2.5, 10)
	p04 = comp.WithPayouts(0, 0, 0.5, 4, 15)
	p05 = comp.WithPayouts(0, 0, 0.5, 4, 15)
	p06 = comp.WithPayouts(0, 0.5, 3, 10, 75)
	p07 = comp.WithPayouts(0, 0.5, 3, 10, 75)
	p08 = comp.WithPayouts(0, 0.5, 4, 40, 200)
	p09 = comp.WithPayouts(0, 1, 10, 100, 500)
	p10 = comp.WithScatterPayouts(scatterPayouts...)
	p11 = comp.WithPayouts()

	// pay-lines.
	p11111 = comp.NewPayline(id01, rows, 1, 1, 1, 1, 1)
	p00000 = comp.NewPayline(id02, rows, 0, 0, 0, 0, 0)
	p22222 = comp.NewPayline(id03, rows, 2, 2, 2, 2, 2)
	p01210 = comp.NewPayline(id04, rows, 0, 1, 2, 1, 0)
	p21012 = comp.NewPayline(id05, rows, 2, 1, 0, 1, 2)
	p12221 = comp.NewPayline(id06, rows, 1, 2, 2, 2, 1)
	p10001 = comp.NewPayline(id07, rows, 1, 0, 0, 0, 1)
	p22100 = comp.NewPayline(id08, rows, 2, 2, 1, 0, 0)
	p00122 = comp.NewPayline(id09, rows, 0, 0, 1, 2, 2)
	p21110 = comp.NewPayline(id10, rows, 2, 1, 1, 1, 0)

	ids       = [symbolCount2]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11}
	names     = [symbolCount2]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11}
	resources = [symbolCount2]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11}
	payouts   = [symbolCount2]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11}
	paylines  = comp.Paylines{p11111, p00000, p22222, p01210, p21012, p12221, p10001, p22100, p00122, p21110}

	bonusWeights92 = util.AcquireWeighting().AddWeights([]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09}, weightsBonus92)
	bonusWeights94 = util.AcquireWeighting().AddWeights([]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09}, weightsBonus94)
	bonusWeights96 = util.AcquireWeighting().AddWeights([]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09}, weightsBonus96)

	flag0 = comp.NewRoundFlag(flagPayoutBand, "payout band")
	flag1 = comp.NewRoundFlag(flagPayoutBand, "yellow book landed")
	flag2 = comp.NewRoundFlag(flagPayoutBand, "full bonus page landed")
	flag3 = comp.NewRoundFlag(flagPayoutBand, "free spin count")
	flags = comp.RoundFlags{flag0, flag1, flag2, flag3}
)

var (
	symbols1 *comp.SymbolSet
	symbols2 *comp.SymbolSet

	actions92first comp.SpinActions
	actions92free  comp.SpinActions
	actions92all   comp.SpinActions
	actions94first comp.SpinActions
	actions94free  comp.SpinActions
	actions94all   comp.SpinActions
	actions96first comp.SpinActions
	actions96free  comp.SpinActions
	actions96all   comp.SpinActions

	slots92 *comp.Slots
	slots94 *comp.Slots
	slots96 *comp.Slots

	slots92params game.RegularParams
	slots94params game.RegularParams
	slots96params game.RegularParams
)

func initActions() {
	// generate scatter symbols on first spin.
	scatterFirst := comp.NewGenerateSymbolAction(redBook, firstScatterWeights).GenerateNoDupes()
	scatterFirst.Describe(scatterFirstID, "generate scatter symbols - first spin")

	// generate scatter symbols on free spins (may get overridden by full page bonus action a3).
	scatterFree1 := comp.NewGenerateSymbolAction(redBook, freeScatterWeights1).GenerateNoDupes()
	scatterFree2 := comp.NewGenerateSymbolAction(redBook, freeScatterWeights2).GenerateNoDupes()
	scatterFree3 := comp.NewGenerateSymbolAction(redBook, freeScatterWeights3).GenerateNoDupes()
	scatterFree4 := comp.NewGenerateSymbolAction(redBook, freeScatterWeights4).GenerateNoDupes()
	scatterFree5 := comp.NewGenerateSymbolAction(redBook, freeScatterWeights5).GenerateNoDupes()

	scatterFree1.Describe(scatterFree1ID, "generate scatter symbols - free spin - band 1")
	scatterFree2.Describe(scatterFree2ID, "generate scatter symbols - free spin - band 2")
	scatterFree3.Describe(scatterFree3ID, "generate scatter symbols - free spin - band 3")
	scatterFree4.Describe(scatterFree4ID, "generate scatter symbols - free spin - band 4")
	scatterFree5.Describe(scatterFree5ID, "generate scatter symbols - free spin - band 5")

	scatterFree := comp.NewMultiActionFlagValue(flagPayoutBand, 1, scatterFree1, 2, scatterFree2, 3, scatterFree3, 4, scatterFree4, 5, scatterFree5)

	// generate bonus symbol full page payout (may get overridden by yellow book generators a4xxx).
	bonusPage1 := comp.NewGenerateBonusAction(5, genFullBonusChances1)
	bonusPage2 := comp.NewGenerateBonusAction(5, genFullBonusChances2)
	bonusPage3 := comp.NewGenerateBonusAction(5, genFullBonusChances3)
	bonusPage4 := comp.NewGenerateBonusAction(5, genFullBonusChances4)
	bonusPage5 := comp.NewGenerateBonusAction(5, genFullBonusChances5)

	bonusPage1.WithTriggerFilters(comp.OnNotRoundFlagValue(flagBonusPage, 1), comp.OnBetweenFreeSpins(0, 2))
	bonusPage2.WithTriggerFilters(comp.OnNotRoundFlagValue(flagBonusPage, 1), comp.OnBetweenFreeSpins(0, 3))
	bonusPage3.WithTriggerFilters(comp.OnNotRoundFlagValue(flagBonusPage, 1), comp.OnBetweenFreeSpins(0, 4))
	bonusPage4.WithTriggerFilters(comp.OnBetweenFreeSpins(0, 5))
	bonusPage5.WithTriggerFilters(comp.OnBetweenFreeSpins(0, 6))

	bonusPage1.Describe(bonusPage1ID, "generate bonus symbol full page - band 1")
	bonusPage2.Describe(bonusPage2ID, "generate bonus symbol full page - band 2")
	bonusPage3.Describe(bonusPage3ID, "generate bonus symbol full page - band 3")
	bonusPage4.Describe(bonusPage4ID, "generate bonus symbol full page - band 4")
	bonusPage5.Describe(bonusPage5ID, "generate bonus symbol full page - band 5")

	bonusPage := comp.NewMultiActionFlagValue(flagPayoutBand, 1, bonusPage1, 2, bonusPage2, 3, bonusPage3, 4, bonusPage4, 5, bonusPage5)

	// morph red to yellow book.
	redYellow92 := comp.NewMorphSymbolAction(util.UInt8s{2, 3, 4}, []float64{yellowChance92}, redBook).GenerateNoDupes()
	redYellow94 := comp.NewMorphSymbolAction(util.UInt8s{2, 3, 4}, []float64{yellowChance94}, redBook).GenerateNoDupes()
	redYellow96 := comp.NewMorphSymbolAction(util.UInt8s{2, 3, 4}, []float64{yellowChance96}, redBook).GenerateNoDupes()

	redYellow92.WithTriggerFilters(comp.OnNotRoundFlagValue(flagPayoutBand, 1))
	redYellow94.WithTriggerFilters(comp.OnNotRoundFlagValue(flagPayoutBand, 1))
	redYellow96.WithTriggerFilters(comp.OnNotRoundFlagValue(flagPayoutBand, 1))

	redYellow92.Describe(redYellow92ID, "morph red to yellow book - !band 1 - rtp 92")
	redYellow94.Describe(redYellow94ID, "morph red to yellow book - !band 1 - rtp 94")
	redYellow96.Describe(redYellow96ID, "morph red to yellow book - !band 1 - rtp 96")

	// morph or generate yellow book in last free spins for payout band 1.
	yellowFree1 := comp.NewGenOrMorphSymbolAction(yellowBook, redBook, forceYellowBook1a, 2, 3, 4).GenerateNoDupes()
	yellowFree2 := comp.NewGenOrMorphSymbolAction(yellowBook, redBook, forceYellowBook1b, 2, 3, 4).GenerateNoDupes()
	yellowFree3 := comp.NewGenOrMorphSymbolAction(yellowBook, redBook, forceYellowBook1c, 2, 3, 4).GenerateNoDupes()

	yellowFree1.WithTriggerFilters(comp.OnRoundFlagValue(flagPayoutBand, 1), comp.OnNotRoundFlagValue(flagYellowBook, 1), comp.OnRemainingFreeSpins(2))
	yellowFree2.WithTriggerFilters(comp.OnRoundFlagValue(flagPayoutBand, 1), comp.OnNotRoundFlagValue(flagYellowBook, 1), comp.OnRemainingFreeSpins(1))
	yellowFree3.WithTriggerFilters(comp.OnRoundFlagValue(flagPayoutBand, 1), comp.OnNotRoundFlagValue(flagYellowBook, 1), comp.OnRemainingFreeSpins(0))

	yellowFree1.Describe(yellowFree1ID, "morph red to yellow book - band 1 - 3rd to last round")
	yellowFree2.Describe(yellowFree2ID, "morph red to yellow book - band 1 - 2nd to last round")
	yellowFree3.Describe(yellowFree3ID, "morph red to yellow book - band 1 - last round")

	// morph or generate yellow book in last free spins for payout bands 2-5.
	yellowFree1B := comp.NewGenOrMorphSymbolAction(yellowBook, redBook, forceYellowBook2a, 2, 3, 4).GenerateNoDupes()
	yellowFree2B := comp.NewGenOrMorphSymbolAction(yellowBook, redBook, forceYellowBook2b, 2, 3, 4).GenerateNoDupes()
	yellowFree3B := comp.NewGenOrMorphSymbolAction(yellowBook, redBook, forceYellowBook2c, 2, 3, 4).GenerateNoDupes()

	yellowFree1B.WithTriggerFilters(comp.OnNotRoundFlagValue(flagPayoutBand, 1), comp.OnNotRoundFlagValue(flagYellowBook, 1), comp.OnRemainingFreeSpins(2))
	yellowFree2B.WithTriggerFilters(comp.OnNotRoundFlagValue(flagPayoutBand, 1), comp.OnNotRoundFlagValue(flagYellowBook, 1), comp.OnRemainingFreeSpins(1))
	yellowFree3B.WithTriggerFilters(comp.OnNotRoundFlagValue(flagPayoutBand, 1), comp.OnNotRoundFlagValue(flagYellowBook, 1), comp.OnRemainingFreeSpins(0))

	yellowFree1B.Describe(yellowFree1bID, "morph red to yellow book - !band 1 - 3rd to last round")
	yellowFree2B.Describe(yellowFree2bID, "morph red to yellow book - !band 1 - 2nd to last round")
	yellowFree3B.Describe(yellowFree3bID, "morph red to yellow book - !band 1 - last round")

	// remove bonus payouts before they can happen.
	removeBonus1 := comp.NewRemoveBonusPayoutsAction(0, maxPayout, removeBonusChance1, 2, false)
	removeBonus2 := comp.NewRemoveBonusPayoutsAction(50, 5000, removeBonusChance2, 2, false)
	removeBonus3 := comp.NewRemoveBonusPayoutsAction(100, 500, removeBonusChance3, 2, false)
	removeBonus4 := comp.NewRemoveBonusPayoutsAction(100, 500, removeBonusChance4, 2, false)
	removeBonus5 := comp.NewRemoveBonusPayoutsAction(100, 500, removeBonusChance5, 2, false)

	removeBonus1.Describe(removeBonus1ID, "remove bonus payouts (free spins - band 1)")
	removeBonus2.Describe(removeBonus2ID, "remove bonus payouts (free spins - band 2)")
	removeBonus3.Describe(removeBonus3ID, "remove bonus payouts (free spins - band 3)")
	removeBonus4.Describe(removeBonus4ID, "remove bonus payouts (free spins - band 4)")
	removeBonus5.Describe(removeBonus5ID, "remove bonus payouts (free spins - band 5)")

	removeBonus := comp.NewMultiActionFlagValue(flagPayoutBand, 1, removeBonus1, 2, removeBonus2, 3, removeBonus3, 4, removeBonus4, 5, removeBonus5)

	// award regular payouts.
	regPayouts := comp.NewPaylinesAction()
	regPayouts.Describe(regPayoutsID, "award regular payouts")

	// remove payouts by band.
	removePayoutsFirst := comp.NewRemovePayoutBandsAction(1, direction, true, false, removeBandsFirst)
	removePayoutsFree1 := comp.NewRemovePayoutBandsAction(1, direction, true, false, removeBandsFree1)
	removePayoutsFree2 := comp.NewRemovePayoutBandsAction(1, direction, true, false, removeBandsFree2)

	removePayoutsFree1.WithTriggerFilters(comp.OnRoundFlagBelow(flagFreeSpins, 3))
	removePayoutsFree2.WithTriggerFilters(comp.OnRoundFlagAbove(flagFreeSpins, 2))

	removePayoutsFirst.Describe(removePayoutsFirstID, "remove payouts by band (first spins)")
	removePayoutsFree1.Describe(removePayoutsFree1ID, "remove payouts by band (first few free spins)")
	removePayoutsFree2.Describe(removePayoutsFree2ID, "remove payouts by band (free spins)")

	// award free spins from scatters.
	scatterGames8 := comp.NewScatterFreeSpinsAction(8, true, redBook, 3, true).WithMultiSymbols(yellowBook)
	scatterGames8.Describe(scatterGames8ID, "award 8 free spins from 3 scatters")
	scatterGames10 := comp.NewScatterFreeSpinsAction(10, true, redBook, 4, true).WithMultiSymbols(yellowBook).WithAlternate(scatterGames8)
	scatterGames10.Describe(scatterGames10ID, "award 10 free spins from 4 scatters")
	scatterGames12 := comp.NewScatterFreeSpinsAction(12, true, redBook, 5, true).WithMultiSymbols(yellowBook).WithAlternate(scatterGames10)
	scatterGames12.Describe(scatterGames12ID, "award 12 free spins from 5 scatters")

	// award payout from scatters.
	scatterPayouts3 := comp.NewScatterPayoutAction(redBook, 3, scatterPayouts[2])
	scatterPayouts3.Describe(scatterPayouts3ID, "award payout from scatters x3")
	scatterPayouts4 := comp.NewScatterPayoutAction(redBook, 4, scatterPayouts[3]).WithAlternate(scatterPayouts3)
	scatterPayouts4.Describe(scatterPayouts4ID, "award payout from scatters x4")
	scatterPayouts5 := comp.NewScatterPayoutAction(redBook, 5, scatterPayouts[4]).WithAlternate(scatterPayouts4)
	scatterPayouts5.Describe(scatterPayouts5ID, "award payout from scatters x5")

	// mark hot reels during free spins.
	hotReels := comp.NewHotAction(yellowBook) // hot action before bonus payouts!!
	hotReels.Describe(hotReelsID, "mark hot reels")

	// award bonus payout during free spins.
	bonusPayouts := comp.NewBonusScatterAction(uint8(len(paylines)))
	bonusPayouts.Describe(bonusPayoutsID, "award bonus payout")

	// initialize round flag 0 for payout bands.
	payoutBandsFlag := comp.NewRoundFlagWeightedAction(flagPayoutBand, freeSpinBands)
	payoutBandsFlag.Describe(payoutBandsID, "mark payout band (flag 0)")

	// update round flag 1 marking yellow book generated.
	yellowBookFlag := comp.NewRoundFlagSymbolUsedAction(flagYellowBook, yellowBook)
	yellowBookFlag.Describe(yellowBookID, "mark yellow book occured (flag 1)")

	// update round flag 2 marking full page bonus awarded.
	bonusPageFLag := comp.NewRoundFlagFullBonusAction(flagBonusPage)
	bonusPageFLag.Describe(bonusPageID, "mark full page bonus occured (flag 2)")

	// update round flag 3 marking sequence of free spin.
	freeSpinsFlag := comp.NewRoundFlagIncreaseAction(flagFreeSpins)
	freeSpinsFlag.Describe(freeSpinsID, "count number of free spins (flag 3)")

	actionsAall := comp.SpinActions{scatterFirst, scatterFree1, scatterFree2, scatterFree3, scatterFree4, scatterFree5,
		bonusPage1, bonusPage2, bonusPage3, bonusPage4, bonusPage5}
	actionsAfirst := comp.SpinActions{scatterFirst}
	actionsAfree := comp.SpinActions{scatterFree, bonusPage}

	actionsBall := comp.SpinActions{yellowFree1, yellowFree2, yellowFree3, yellowFree1B, yellowFree2B, yellowFree3B,
		removeBonus1, removeBonus2, removeBonus3, removeBonus4, removeBonus5,
		regPayouts, removePayoutsFirst, removePayoutsFree1, removePayoutsFree2,
		scatterGames12, scatterPayouts5, hotReels, bonusPayouts,
		payoutBandsFlag, yellowBookFlag, bonusPageFLag, freeSpinsFlag}
	actionsBfirst := comp.SpinActions{regPayouts, removePayoutsFirst,
		scatterGames12, scatterPayouts5, payoutBandsFlag}
	actionsBfree := comp.SpinActions{yellowFree1, yellowFree2, yellowFree3, yellowFree1B, yellowFree2B, yellowFree3B,
		removeBonus,
		regPayouts, removePayoutsFree1, removePayoutsFree2,
		scatterGames12, scatterPayouts5, hotReels, bonusPayouts,
		yellowBookFlag, bonusPageFLag, freeSpinsFlag}

	actions92all = append(append(actionsAall, redYellow92), actionsBall...)
	actions92first = append(actionsAfirst, actionsBfirst...)
	actions92free = append(append(actionsAfree, redYellow92), actionsBfree...)
	actions94all = append(append(actionsAall, redYellow94), actionsBall...)
	actions94first = append(actionsAfirst, actionsBfirst...)
	actions94free = append(append(actionsAfree, redYellow94), actionsBfree...)
	actions96all = append(append(actionsAall, redYellow96), actionsBall...)
	actions96first = append(actionsAfirst, actionsBfirst...)
	actions96free = append(append(actionsAfree, redYellow96), actionsBfree...)
}

func initSlots(target float64, weights1, weights2 [symbolCount2]comp.SymbolOption, bonusWeights util.WeightedGenerator, actions1, actions2 []comp.SpinActioner) *comp.Slots {
	ss := make([]*comp.Symbol, symbolCount)
	for ix := range ss {
		if id := ids[ix]; id == redBook {
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights1[ix], comp.WithKind(comp.WildScatter))
		} else {
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights1[ix])
		}
	}
	s1 := comp.NewSymbolSet(ss...).SetBonusWeights(bonusWeights)

	ss2 := make([]*comp.Symbol, symbolCount2)
	copy(ss2, ss)
	ss2[symbolCount-1] = comp.NewSymbol(redBook, n10, r10, p10, weights2[symbolCount-1], comp.WithKind(comp.WildScatter), comp.MorphInto(yellowBook))
	ss2[symbolCount] = comp.NewSymbol(yellowBook, n11, r11, p11, weights2[symbolCount], comp.WithKind(comp.WildScatter))
	s2 := comp.NewSymbolSet(ss2...)

	if symbols1 == nil {
		symbols1 = s1
		symbols2 = s2
	}

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.NoRepeat(noRepeat),
		comp.WithSymbols(s1),
		comp.WithAltSymbols(s2),
		comp.WithPaylines(direction, true, paylines...),
		comp.MaxPayout(maxPayout),
		comp.HotReelsAsBonusSymbol(),
		comp.WithRTP(target),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, nil, nil),
	)
}

func init() {
	initActions()

	slots92 = initSlots(92.0, weights92a, weights92b, bonusWeights92, actions92first, actions92free)
	slots94 = initSlots(94.0, weights94a, weights94b, bonusWeights94, actions94first, actions94free)
	slots96 = initSlots(96.0, weights96a, weights96b, bonusWeights96, actions96first, actions96free)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}
