package lam

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount  = 12
	reels        = 5
	rows         = 3
	noRepeat     = 2
	direction    = comp.PayLTR
	bonusBuyCost = 150
	maxPayout    = 11000.0

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

	scatter = id11
	wild    = id12

	bonusBuyID       = 1
	bonusChoiceID    = 2
	wilds92aID       = 10
	wilds92bID       = 11
	wilds92cID       = 12
	wilds92dID       = 13
	wilds92eID       = 14
	scatters92aID    = 15
	scatters92bID    = 16
	scatters92cID    = 17
	scatters92dID    = 18
	scatters92eID    = 19
	wilds94aID       = 20
	wilds94bID       = 21
	wilds94cID       = 22
	wilds94dID       = 23
	wilds94eID       = 24
	scatters94aID    = 25
	scatters94bID    = 26
	scatters94cID    = 27
	scatters94dID    = 28
	scatters94eID    = 29
	wilds96aID       = 30
	wilds96bID       = 31
	wilds96cID       = 32
	wilds96dID       = 33
	wilds96eID       = 34
	scatters96aID    = 35
	scatters96bID    = 36
	scatters96cID    = 37
	scatters96dID    = 38
	scatters96eID    = 39
	winlinesID       = 60
	freeGamesFirstID = 61
	freeGamesFreeID  = 62
	doubleMultID     = 71
	stickiesID       = 72

	flagBonusBuy    = 0
	flagBonusChoice = 1

	northWing = 1
	southWing = 2

	freeSpinsFirst = 11
	freeSpinsFree  = 5
)

var (
	// WIP WIP WIP
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 30, 40, 40),
		comp.WithWeights(40, 50, 30, 50, 40),
		comp.WithWeights(40, 20, 30, 20, 40),
		comp.WithWeights(40, 30, 20, 30, 40),
		comp.WithWeights(30, 40, 25, 40, 30),
		comp.WithWeights(30, 40, 25, 40, 30),
		comp.WithWeights(20, 16, 15, 16, 20),
		comp.WithWeights(12, 20, 12, 20, 12),
		comp.WithWeights(12, 20, 12, 20, 12),
		comp.WithWeights(10, 13, 20, 13, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	scatterChances92a = []float64{40, 26.7, 4.4, 1, 0.5}                                                                // base game
	scatterChances92b = []float64{45, 34, 12, 1, 0.5}                                                                   // non-paid, north
	scatterChances92c = []float64{48, 30, 15, 0.5, 0.1}                                                                 // non-paid, south
	scatterChances92d = []float64{40, 20, 10, 2, 1}                                                                     // paid, north
	scatterChances92e = []float64{30, 15, 10, 2, 1}                                                                     // paid, south
	wildChances92a    = []float64{29, 9, 2, 0.5}                                                                        // base game
	wildChances92b    = []float64{59.4, 55, 25, 12, 4}                                                                  // non-paid north
	wildWeights92b    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{35, 20, 15, 7.5, 5})  // non-paid north
	wildChances92c    = []float64{6.77, 0.5}                                                                            // non-paid south
	wildChances92d    = []float64{35, 18, 5, 0.5}                                                                       // paid north
	wildWeights92d    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{35, 20, 25, 12.5, 5}) // paid north
	wildChances92e    = []float64{32, 15, 2}                                                                            // paid south

	// WIP WIP WIP
	weights94 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 40, 40, 40), // l5
		comp.WithWeights(40, 40, 40, 40, 40), // l4
		comp.WithWeights(40, 40, 40, 40, 40), // l3
		comp.WithWeights(32, 32, 34, 36, 36), // l2
		comp.WithWeights(32, 32, 34, 36, 36), // l1
		comp.WithWeights(25, 25, 28, 32, 32), // h5
		comp.WithWeights(23, 23, 26, 30, 30), // h4
		comp.WithWeights(21, 21, 24, 28, 28), // h3
		comp.WithWeights(19, 19, 22, 26, 26), // h2
		comp.WithWeights(16, 16, 20, 24, 24), // h1
		comp.WithWeights(0, 0, 0, 0, 0),      // scatter/bonus
		comp.WithWeights(0, 0, 0, 0, 0),      // wild
	}
	scatterChances94a = []float64{42, 25, 5, 2, 1}                                                                      //
	scatterChances94b = []float64{45, 34, 12, 2, 1}                                                                     // non-paid, north
	scatterChances94c = []float64{50, 40, 8, 2, 1}                                                                      // non-paid, south
	scatterChances94d = []float64{40, 20, 10, 2, 1}                                                                     // paid, north
	scatterChances94e = []float64{50, 40, 8, 2, 1}                                                                      // paid, south
	wildChances94a    = []float64{43, 40, 35, 30, 20}                                                                   //
	wildChances94b    = []float64{66, 55, 25, 12, 4}                                                                    // non-paid north
	wildWeights94b    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{35, 20, 15, 7.5, 5})  // non-paid north
	wildChances94c    = []float64{10}                                                                                   // non-paid south
	wildChances94d    = []float64{35, 18, 5, 0.5}                                                                       // paid north
	wildWeights94d    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{35, 20, 25, 12.5, 5}) // paid north
	wildChances94e    = []float64{10}

	// WIP WIP WIP
	weights96 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 40, 40, 40), // l5
		comp.WithWeights(40, 40, 40, 40, 40), // l4
		comp.WithWeights(40, 40, 40, 40, 40), // l3
		comp.WithWeights(32, 32, 34, 36, 36), // l2
		comp.WithWeights(32, 32, 34, 36, 36), // l1
		comp.WithWeights(25, 25, 28, 32, 32), // h5
		comp.WithWeights(23, 23, 26, 30, 30), // h4
		comp.WithWeights(21, 21, 24, 28, 28), // h3
		comp.WithWeights(19, 19, 22, 26, 26), // h2
		comp.WithWeights(16, 16, 20, 24, 24), // h1
		comp.WithWeights(0, 0, 0, 0, 0),      // scatter/bonus
		comp.WithWeights(0, 0, 0, 0, 0),      // wild
	}
	scatterChances96a = []float64{42, 25, 5, 2, 1}                                                                      //
	scatterChances96b = []float64{45, 34, 12, 2, 1}                                                                     // non-paid, north
	scatterChances96c = []float64{50, 40, 8, 2, 1}                                                                      // non-paid, south
	scatterChances96d = []float64{40, 20, 10, 2, 1}                                                                     // paid, north
	scatterChances96e = []float64{50, 40, 8, 2, 1}                                                                      // paid, south
	wildChances96a    = []float64{43, 40, 35, 30, 20}                                                                   //
	wildChances96b    = []float64{66, 55, 25, 12, 4}                                                                    // non-paid north
	wildWeights96b    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{35, 20, 15, 7.5, 5})  // non-paid north
	wildChances96c    = []float64{10}                                                                                   // non-paid south
	wildChances96d    = []float64{35, 18, 5, 0.5}                                                                       // paid north
	wildWeights96d    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{35, 20, 25, 12.5, 5}) // paid north
	wildChances96e    = []float64{10}
)

var (
	wildWeightsSouth = util.AcquireWeighting().AddWeights(util.Indexes{2}, []float64{1})

	n01 = comp.WithName("Soap")
	n02 = comp.WithName("Baton")
	n03 = comp.WithName("Taser")
	n04 = comp.WithName("Shiv")
	n05 = comp.WithName("Handcuffs")
	n06 = comp.WithName("Control Tower")
	n07 = comp.WithName("Corrupt Guard")
	n08 = comp.WithName("North Wing Prisoner")
	n09 = comp.WithName("South Wing Prisoner")
	n10 = comp.WithName("Prison Director")
	n11 = comp.WithName("Sirens")
	n12 = comp.WithName("Wild")

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
	r11 = comp.WithResource("scatter")
	r12 = comp.WithResource("wild")

	p01 = comp.WithPayouts(0, 0, 0.5, 1.5, 5)
	p02 = comp.WithPayouts(0, 0, 0.5, 1.5, 5)
	p03 = comp.WithPayouts(0, 0, 0.5, 1.5, 5)
	p04 = comp.WithPayouts(0, 0, 0.6, 2, 6)
	p05 = comp.WithPayouts(0, 0, 0.6, 2, 6)
	p06 = comp.WithPayouts(0, 0, 0.8, 2.5, 10)
	p07 = comp.WithPayouts(0, 0, 1, 3, 15)
	p08 = comp.WithPayouts(0, 0, 2, 4, 20)
	p09 = comp.WithPayouts(0, 0, 3, 6, 30)
	p10 = comp.WithPayouts(0, 0, 5, 10, 50)
	p11 = comp.WithPayouts()
	p12 = p11

	p11111 = comp.NewPayline(id01, rows, 1, 1, 1, 1, 1)
	p00000 = comp.NewPayline(id02, rows, 0, 0, 0, 0, 0)
	p22222 = comp.NewPayline(id03, rows, 2, 2, 2, 2, 2)
	p01210 = comp.NewPayline(id04, rows, 0, 1, 2, 1, 0)
	p21012 = comp.NewPayline(id05, rows, 2, 1, 0, 1, 2)
	p00100 = comp.NewPayline(id06, rows, 0, 0, 1, 0, 0)
	p22122 = comp.NewPayline(id07, rows, 2, 2, 1, 2, 2)
	p12221 = comp.NewPayline(id08, rows, 1, 2, 2, 2, 1)
	p10001 = comp.NewPayline(id09, rows, 1, 0, 0, 0, 1)
	p10101 = comp.NewPayline(id10, rows, 1, 0, 1, 0, 1)
	p12121 = comp.NewPayline(id11, rows, 1, 2, 1, 2, 1)

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11, r12}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11, p12}
	paylines  = comp.Paylines{p11111, p00000, p22222, p01210, p21012, p00100, p22122, p12221, p10001, p10101, p12121}

	flag0 = comp.NewRoundFlag(flagBonusBuy, "bonus buy")
	flag1 = comp.NewRoundFlag(flagBonusChoice, "wing choice")
	flags = comp.RoundFlags{flag0, flag1}
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
	// bonus buy feature.
	paidBonus := comp.NewPaidAction(comp.FreeSpins, 11, bonusBuyCost, scatter, 3).WithFlag(flagBonusBuy, bonusBuyID)
	paidBonus.Describe(bonusBuyID, "bonus buy feature")

	// player choice feature.
	bonusChoice := comp.NewPlayerChoiceAction(flagBonusChoice, "wing", []string{"north", "south"}, []int{northWing, southWing})
	bonusChoice.WithTestChoicesFilters(comp.OnRoundResume)
	bonusChoice.Describe(bonusChoiceID, "bonus choice feature")

	// generate scatter symbols RTP 92.
	scatters92a := comp.NewGenerateSymbolAction(scatter, scatterChances92a).GenerateNoDupes()
	scatters92a.Describe(scatters92aID, "generate scatters - first spin - RTP 92")
	scatters92b := comp.NewGenerateSymbolAction(scatter, scatterChances92b).GenerateNoDupes()
	scatters92b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	scatters92b.Describe(scatters92bID, "generate scatters - free spins north base - RTP 92")
	scatters92c := comp.NewGenerateSymbolAction(scatter, scatterChances92c).GenerateNoDupes()
	scatters92c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	scatters92c.Describe(scatters92cID, "generate scatters - free spins south base - RTP 92")
	scatters92d := comp.NewGenerateSymbolAction(scatter, scatterChances92d).GenerateNoDupes()
	scatters92d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	scatters92d.Describe(scatters92dID, "generate scatters - free spins north paid - RTP 92")
	scatters92e := comp.NewGenerateSymbolAction(scatter, scatterChances92e).GenerateNoDupes()
	scatters92e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	scatters92e.Describe(scatters92eID, "generate scatters - free spins south paid - RTP 92")

	// generate scatter symbols RTP 94.
	scatters94a := comp.NewGenerateSymbolAction(scatter, scatterChances94a).GenerateNoDupes()
	scatters94a.Describe(scatters94aID, "generate scatters - first spin - RTP 94")
	scatters94b := comp.NewGenerateSymbolAction(scatter, scatterChances94b).GenerateNoDupes()
	scatters94b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	scatters94b.Describe(scatters94bID, "generate scatters - free spins north base - RTP 94")
	scatters94c := comp.NewGenerateSymbolAction(scatter, scatterChances94c).GenerateNoDupes()
	scatters94c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	scatters94c.Describe(scatters94cID, "generate scatters - free spins south base - RTP 94")
	scatters94d := comp.NewGenerateSymbolAction(scatter, scatterChances94d).GenerateNoDupes()
	scatters94d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	scatters94d.Describe(scatters94dID, "generate scatters - free spins north paid - RTP 94")
	scatters94e := comp.NewGenerateSymbolAction(scatter, scatterChances94e).GenerateNoDupes()
	scatters94e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	scatters94e.Describe(scatters94eID, "generate scatters - free spins south paid - RTP 94")

	// generate scatter symbols RTP 96.
	scatters96a := comp.NewGenerateSymbolAction(scatter, scatterChances96a).GenerateNoDupes()
	scatters96a.Describe(scatters96aID, "generate scatters - first spin - RTP 96")
	scatters96b := comp.NewGenerateSymbolAction(scatter, scatterChances96b).GenerateNoDupes()
	scatters96b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	scatters96b.Describe(scatters96bID, "generate scatters - free spins north base - RTP 96")
	scatters96c := comp.NewGenerateSymbolAction(scatter, scatterChances96c).GenerateNoDupes()
	scatters96c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	scatters96c.Describe(scatters96cID, "generate scatters - free spins south base - RTP 96")
	scatters96d := comp.NewGenerateSymbolAction(scatter, scatterChances96d).GenerateNoDupes()
	scatters96d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	scatters96d.Describe(scatters96dID, "generate scatters - free spins north paid - RTP 96")
	scatters96e := comp.NewGenerateSymbolAction(scatter, scatterChances96e).GenerateNoDupes()
	scatters96e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	scatters96e.Describe(scatters96eID, "generate scatters - free spins south paid - RTP 96")

	// generate wild symbols RTP 92.
	wilds92a := comp.NewGenerateSymbolAction(wild, wildChances92a)
	wilds92a.WithTriggerFilters(comp.OnFirstSpin)
	wilds92a.Describe(wilds92aID, "generate wilds - first spin - RTP 92")
	wilds92b := comp.NewGenerateSymbolAction(wild, wildChances92b).WithMultipliers(wildWeights92b)
	wilds92b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	wilds92b.Describe(wilds92bID, "generate wilds - free spins north base - RTP 92")
	wilds92c := comp.NewGenerateSymbolAction(wild, wildChances92c).WithMultipliers(wildWeightsSouth)
	wilds92c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	wilds92c.Describe(wilds92cID, "generate wilds - free spins south base - RTP 92")
	wilds92d := comp.NewGenerateSymbolAction(wild, wildChances92d).WithMultipliers(wildWeights92d)
	wilds92d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	wilds92d.Describe(wilds92dID, "generate wilds - free spins north paid - RTP 92")
	wilds92e := comp.NewGenerateSymbolAction(wild, wildChances92e).WithMultipliers(wildWeightsSouth)
	wilds92e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	wilds92e.Describe(wilds92eID, "generate wilds - free spins south paid - RTP 92")

	// generate wild symbols RTP 94.
	wilds94a := comp.NewGenerateSymbolAction(wild, wildChances94a)
	wilds94a.Describe(wilds94aID, "generate wilds - first spin - RTP 94")
	wilds94b := comp.NewGenerateSymbolAction(wild, wildChances94b).WithMultipliers(wildWeights94b)
	wilds94b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	wilds94b.Describe(wilds94bID, "generate wilds - free spins north base - RTP 94")
	wilds94c := comp.NewGenerateSymbolAction(wild, wildChances94c).WithMultipliers(wildWeightsSouth)
	wilds94c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	wilds94c.Describe(wilds94cID, "generate wilds - free spins south base - RTP 94")
	wilds94d := comp.NewGenerateSymbolAction(wild, wildChances94d).WithMultipliers(wildWeights94d)
	wilds94d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	wilds94d.Describe(wilds94dID, "generate wilds - free spins north paid - RTP 94")
	wilds94e := comp.NewGenerateSymbolAction(wild, wildChances94e).WithMultipliers(wildWeightsSouth)
	wilds94e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	wilds94e.Describe(wilds94eID, "generate wilds - free spins south paid - RTP 94")

	// generate wild symbols RTP 96.
	wilds96a := comp.NewGenerateSymbolAction(wild, wildChances96a)
	wilds96a.Describe(wilds96aID, "generate wilds - first spin - RTP 96")
	wilds96b := comp.NewGenerateSymbolAction(wild, wildChances96b).WithMultipliers(wildWeights96b)
	wilds96b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	wilds96b.Describe(wilds96bID, "generate wilds - free spins north base - RTP 96")
	wilds96c := comp.NewGenerateSymbolAction(wild, wildChances96c).WithMultipliers(wildWeightsSouth)
	wilds96c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	wilds96c.Describe(wilds96cID, "generate wilds - free spins south base - RTP 96")
	wilds96d := comp.NewGenerateSymbolAction(wild, wildChances96d).WithMultipliers(wildWeights96d)
	wilds96d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, northWing))
	wilds96d.Describe(wilds96dID, "generate wilds - free spins north paid - RTP 96")
	wilds96e := comp.NewGenerateSymbolAction(wild, wildChances96e).WithMultipliers(wildWeightsSouth)
	wilds96e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	wilds96e.Describe(wilds96eID, "generate wilds - free spins south paid - RTP 96")

	// calculate winlines.
	winlines := comp.NewPaylinesAction()
	winlines.Describe(winlinesID, "calculate winlines")

	// award free spins.
	freeGamesFirst := comp.NewScatterFreeSpinsAction(freeSpinsFirst, false, scatter, 3, false).WithPlayerChoice()
	freeGamesFirst.Describe(freeGamesFirstID, "award free games - first spin")
	freeGamesFree := comp.NewScatterFreeSpinsAction(freeSpinsFree, false, scatter, 3, false)
	freeGamesFree.Describe(freeGamesFreeID, "award free games - free spins")

	// double existing wild multipliers when new wild lands in South wing.
	doubleMultipliers := comp.NewGridMultipliersAction(wild, wild, 2, 200)
	doubleMultipliers.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	doubleMultipliers.Describe(doubleMultID, "double multipliers - free spins south")

	// make wilds sticky in South wing.
	stickies := comp.NewStickySymbolAction(wild)
	stickies.WithStickyFilters(comp.OnRoundFlagValue(flagBonusChoice, southWing))
	stickies.Describe(stickiesID, "make wilds sticky - free spins south")

	actionsAall := comp.SpinActions{paidBonus, bonusChoice}
	actionsAfirst := comp.SpinActions{paidBonus, bonusChoice}
	actionsAfree := comp.SpinActions{bonusChoice}

	actionsB92all := comp.SpinActions{scatters92a, scatters92b, scatters92c, scatters92d, scatters92e, wilds92a, wilds92b, wilds92c, wilds92d, wilds92e}
	actionsB92first := comp.SpinActions{scatters92a, wilds92a}
	actionsB92free := comp.SpinActions{scatters92b, scatters92c, wilds92b, wilds92c}
	actionsB92freeBB := comp.SpinActions{scatters92d, scatters92e, wilds92d, wilds92e}

	actionsB94all := comp.SpinActions{scatters94a, scatters94b, scatters94c, scatters94d, scatters94e, wilds94a, wilds94b, wilds94c, wilds94d, wilds94e}
	actionsB94first := comp.SpinActions{scatters94a, wilds94a}
	actionsB94free := comp.SpinActions{scatters94b, scatters94c, wilds94b, wilds94c}
	actionsB94freeBB := comp.SpinActions{scatters94d, scatters94e, wilds94d, wilds94e}

	actionsB96all := comp.SpinActions{scatters96a, scatters96b, scatters96c, scatters96d, scatters96e, wilds96a, wilds96b, wilds96c, wilds96d, wilds96e}
	actionsB96first := comp.SpinActions{scatters96a, wilds96a}
	actionsB96free := comp.SpinActions{scatters96b, scatters96c, wilds96b, wilds96c}
	actionsB96freeBB := comp.SpinActions{scatters96d, scatters96e, wilds96d, wilds96e}

	actionsCall := comp.SpinActions{winlines, freeGamesFirst, freeGamesFree, doubleMultipliers, stickies}
	actionsCfirst := comp.SpinActions{winlines, freeGamesFirst}
	actionsCfree := comp.SpinActions{winlines, freeGamesFree, doubleMultipliers, stickies}

	actions92all = append(actionsAall, append(actionsB92all, actionsCall...)...)
	actions92first = append(actionsAfirst, append(actionsB92first, actionsCfirst...)...)
	actions92free = append(actionsAfree, append(actionsB92free, actionsCfree...)...)
	actions92firstBB = actions92first
	actions92freeBB = append(actionsAfree, append(actionsB92freeBB, actionsCfree...)...)

	actions94all = append(actionsAall, append(actionsB94all, actionsCall...)...)
	actions94first = append(actionsAfirst, append(actionsB94first, actionsCfirst...)...)
	actions94free = append(actionsAfree, append(actionsB94free, actionsCfree...)...)
	actions94firstBB = actions94first
	actions94freeBB = append(actionsAfree, append(actionsB94freeBB, actionsCfree...)...)

	actions96all = append(actionsAall, append(actionsB96all, actionsCall...)...)
	actions96first = append(actionsAfirst, append(actionsB96first, actionsCfirst...)...)
	actions96free = append(actionsAfree, append(actionsB96free, actionsCfree...)...)
	actions96firstBB = actions96first
	actions96freeBB = append(actionsAfree, append(actionsB96freeBB, actionsCfree...)...)
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
	symbols = comp.NewSymbolSet(ss...)

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.NoRepeat(noRepeat),
		comp.WithSymbols(symbols),
		comp.WithPaylines(direction, true, paylines...),
		comp.WithPlayerChoice(),
		comp.WithBonusBuy(),
		comp.MaxPayout(maxPayout),
		comp.WithRTP(target),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, actions3, actions4),
	)
}

func init() {
	initActions()

	slots92 = initSlots(92.0, weights92, actions92first, actions92free, actions92firstBB, actions92freeBB)
	slots94 = initSlots(94.0, weights94, actions94first, actions94free, actions94firstBB, actions94freeBB)
	slots96 = initSlots(96.0, weights96, actions96first, actions96free, actions96firstBB, actions96freeBB)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}
