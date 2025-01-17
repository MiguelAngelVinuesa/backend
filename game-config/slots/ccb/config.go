package ccb

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount  = 13
	reels        = 5
	rows         = 3
	scatterMin   = 3
	flagCount    = 5
	minFreeSpins = 6
	maxFreeSpins = 15
	maxPayout    = 3000.0

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

	bomb = id13

	a1ID    = 1
	a2ID    = 2
	a3r92ID = 3
	a4r92ID = 4
	a5r92ID = 5
	a3r94ID = 6
	a4r94ID = 7
	a5r94ID = 8
	a3r96ID = 9
	a4r96ID = 10
	a5r96ID = 11
	a6ID    = 15
	a7ID    = 16
	a8ID    = 17
	a9ID    = 18
	a10ID   = 19
	a11ID   = 20
)

var (
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(35, 40, 35, 40, 35),
		comp.WithWeights(35, 30, 35, 30, 35),
		comp.WithWeights(35, 30, 35, 30, 35),
		comp.WithWeights(30, 35, 30, 35, 30),
		comp.WithWeights(30, 25, 30, 25, 30),
		comp.WithWeights(25, 30, 25, 30, 25),
		comp.WithWeights(21, 18, 21, 18, 21),
		comp.WithWeights(18, 20, 18, 20, 18),
		comp.WithWeights(14, 10, 14, 10, 14),
		comp.WithWeights(10, 12, 10, 12, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	weights94 = weights92
	weights96 = weights92

	weightsDedupe1 = util.AcquireWeighting().AddWeights(util.Indexes{2, 3, 4}, []float64{180, 27.5, 1})
	weightsDedupe2 = util.AcquireWeighting().AddWeights(util.Indexes{0, 1, 2, 3, 4}, []float64{150, 35, 10, 2.5, 0.5})
	weightsSuperX  = util.AcquireWeighting().AddWeights(
		util.Indexes{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12},
		[]float64{30, 30, 35, 35, 35, 35, 30, 30, 25, 25, 20, 15})

	// WIP WIP WIP
	pctChangeSuperX92   = 4.6
	pctChangeBomb92     = []float64{2.2, 1}
	pctChanceFreeBomb92 = []float64{60, 7, 1}

	// WIP WIP WIP
	pctChangeSuperX94   = 4.7
	pctChangeBomb94     = []float64{2.4, 1}
	pctChanceFreeBomb94 = []float64{62.5, 12, 1}

	// WIP WIP WIP
	pctChangeSuperX96   = 4.7
	pctChangeBomb96     = []float64{2.5, 1}
	pctChanceFreeBomb96 = []float64{64, 15, 2.7, 1}
)

var (
	n01 = comp.WithName("Nine")
	n02 = comp.WithName("Ten")
	n03 = comp.WithName("Jack")
	n04 = comp.WithName("Queen")
	n05 = comp.WithName("King")
	n06 = comp.WithName("Ace")
	n07 = comp.WithName("Cherry")
	n08 = comp.WithName("Orange")
	n09 = comp.WithName("Strawberry")
	n10 = comp.WithName("Watermelon")
	n11 = comp.WithName("Bell")
	n12 = comp.WithName("Diamond")
	n13 = comp.WithName("Bomb")

	r01 = comp.WithResource("l6")
	r02 = comp.WithResource("l5")
	r03 = comp.WithResource("l4")
	r04 = comp.WithResource("l3")
	r05 = comp.WithResource("l2")
	r06 = comp.WithResource("l1")
	r07 = comp.WithResource("h6")
	r08 = comp.WithResource("h5")
	r09 = comp.WithResource("h4")
	r10 = comp.WithResource("h3")
	r11 = comp.WithResource("h2")
	r12 = comp.WithResource("h1")
	r13 = comp.WithResource("bomb")

	// original game pay-table.
	p01 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 50)
	p02 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 50)
	p03 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 50)
	p04 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 50)
	p05 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 50)
	p06 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 50)
	p07 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 60)
	p08 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 60)
	p09 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 70)
	p10 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 80)
	p11 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 90)
	p12 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 7, 10, 12, 15, 20, 30, 40, 100)
	p13 = comp.WithPayouts()

	// GDD pay-table.
	gdd01 = comp.WithPayouts(0, 0, 0.2, 0.5, 1, 1.5, 2, 2.5, 3, 4, 5, 6, 7, 8, 10)
	gdd02 = comp.WithPayouts(0, 0, 0.2, 0.5, 1, 1.5, 2, 2.5, 3, 4, 5, 6, 7, 8, 10)
	gdd03 = comp.WithPayouts(0, 0, 0.2, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15)
	gdd04 = comp.WithPayouts(0, 0, 0.2, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15)
	gdd05 = comp.WithPayouts(0, 0, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20)
	gdd06 = comp.WithPayouts(0, 0, 0.5, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20)
	gdd07 = comp.WithPayouts(0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 30)
	gdd08 = comp.WithPayouts(0, 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 25, 40)
	gdd09 = comp.WithPayouts(0, 0, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 25, 30, 50)
	gdd10 = comp.WithPayouts(0, 0, 2, 3, 4, 5, 6, 7, 8, 10, 20, 25, 30, 50, 75)
	gdd11 = comp.WithPayouts(0, 0, 3, 4, 5, 6, 7, 8, 9, 10, 25, 30, 50, 75, 100)
	gdd12 = comp.WithPayouts(0, 0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 75, 100, 200)
	gdd13 = comp.WithPayouts()

	ids        = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12, id13}
	names      = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12, n13}
	resources  = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11, r12, r13}
	payouts    = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11, p12, p13}
	gddPayouts = [symbolCount]comp.SymbolOption{gdd01, gdd02, gdd03, gdd04, gdd05, gdd06, gdd07, gdd08, gdd09, gdd10, gdd11, gdd12, gdd13}

	centersSuperX = comp.GridOffsets{{1, 1}, {2, 1}, {3, 1}}
	gridSuperX    = comp.GridOffsets{{-1, -1}, {-1, 1}, {0, 0}, {1, -1}, {1, 1}}
	gridBomb      = comp.GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

	actionsA   comp.SpinActions
	actionsB92 comp.SpinActions
	actionsB94 comp.SpinActions
	actionsB96 comp.SpinActions
	actionsC   comp.SpinActions
	a11        *comp.StateAction

	actions92all comp.SpinActions
	actions94all comp.SpinActions
	actions96all comp.SpinActions
)

var (
	symbols *comp.SymbolSet

	slots92 *comp.Slots
	slots94 *comp.Slots
	slots96 *comp.Slots

	slots92param game.RegularParams
	slots94param game.RegularParams
	slots96param game.RegularParams
)

func initAction() {
	a1 := comp.NewDeduplicationAction(weightsDedupe1).WithSpinKinds([]comp.SpinKind{comp.FirstSpin, comp.FirstFreeSpin})
	a1.Describe(a1ID, "deduplicate symbols - first spin")
	a2 := comp.NewDeduplicationAction(weightsDedupe2).WithSpinKinds([]comp.SpinKind{comp.SecondSpin, comp.SecondFreeSpin})
	a2.Describe(a2ID, "deduplicate symbols - second spin")

	a3r92 := comp.NewGenerateShapeAction(pctChangeSuperX92, gridSuperX, centersSuperX, weightsSuperX).WithSpinKinds([]comp.SpinKind{comp.FirstSpin, comp.FirstFreeSpin})
	a3r92.Describe(a3r92ID, "generate super-x - RTP 92")
	a4r92 := comp.NewGenerateSymbolAction(bomb, pctChangeBomb92).WithSpinKinds([]comp.SpinKind{comp.SecondSpin})
	a4r92.Describe(a4r92ID, "generate bomb - second spin - RTP 92")
	a5r92 := comp.NewGenerateSymbolAction(bomb, pctChanceFreeBomb92).WithSpinKinds([]comp.SpinKind{comp.SecondFreeSpin})
	a5r92.Describe(a5r92ID, "generate bomb - second free spin - RTP 92")

	a3r94 := comp.NewGenerateShapeAction(pctChangeSuperX94, gridSuperX, centersSuperX, weightsSuperX).WithSpinKinds([]comp.SpinKind{comp.FirstSpin, comp.FirstFreeSpin})
	a3r94.Describe(a3r94ID, "generate super-x - RTP 94")
	a4r94 := comp.NewGenerateSymbolAction(bomb, pctChangeBomb94).WithSpinKinds([]comp.SpinKind{comp.SecondSpin})
	a4r94.Describe(a4r94ID, "generate bomb - second spin - RTP 94")
	a5r94 := comp.NewGenerateSymbolAction(bomb, pctChanceFreeBomb94).WithSpinKinds([]comp.SpinKind{comp.SecondFreeSpin})
	a5r94.Describe(a5r94ID, "generate bomb - second free spin - RTP 94")

	a3r96 := comp.NewGenerateShapeAction(pctChangeSuperX96, gridSuperX, centersSuperX, weightsSuperX).WithSpinKinds([]comp.SpinKind{comp.FirstSpin, comp.FirstFreeSpin})
	a3r96.Describe(a3r96ID, "generate super-x - RTP 96")
	a4r96 := comp.NewGenerateSymbolAction(bomb, pctChangeBomb96).WithSpinKinds([]comp.SpinKind{comp.SecondSpin})
	a4r96.Describe(a4r96ID, "generate bomb - second spin - RTP 96")
	a5r96 := comp.NewGenerateSymbolAction(bomb, pctChanceFreeBomb96).WithSpinKinds([]comp.SpinKind{comp.SecondFreeSpin})
	a5r96.Describe(a5r96ID, "generate bomb - second free spin - RTP 96")

	a6 := comp.NewBestSymbolStickyAction()
	a6.Describe(a6ID, "select best sticky symbol")
	a7 := comp.NewSymbolsChooseStickyAction()
	a7.Describe(a7ID, "sticky symbol choices array")
	a8 := comp.NewSuperShapeAction(gridSuperX, centersSuperX)
	a8.Describe(a8ID, "super-x detection")
	a9 := comp.NewWildTransform(bomb, true, gridBomb).WithBombEffect()
	a9.Describe(a9ID, "bomb expansion")
	a10 := comp.NewAllScatterAction(scatterMin)
	a10.Describe(a10ID, "all scatter payouts")

	a11 = comp.NewFlagSymbolsAction(flagCount, false, comp.WithFreeSpins(0, minFreeSpins, maxFreeSpins))
	a11.Describe(a11ID, "award free spins from flags")

	actionsA = comp.SpinActions{a1, a2}
	actionsB92 = comp.SpinActions{a3r92, a4r92, a5r92}
	actionsB94 = comp.SpinActions{a3r94, a4r94, a5r94}
	actionsB96 = comp.SpinActions{a3r96, a4r96, a5r96}
	// actionsB42 = comp.SpinActions{a3r42, a4r42, a5r42}
	actionsC = comp.SpinActions{a6, a7, a8, a9, a10, a11}

	actions92all = append(append(actionsA, actionsB92...), actionsC...)
	actions94all = append(append(actionsA, actionsB94...), actionsC...)
	actions96all = append(append(actionsA, actionsB96...), actionsC...)
}

func initSlots(target float64, weights [symbolCount]comp.SymbolOption, payouts [symbolCount]comp.SymbolOption, actions []comp.SpinActioner) *comp.Slots {
	ss := make([]*comp.Symbol, symbolCount)
	for ix := range ss {
		if ids[ix] == bomb {
			ss[ix] = comp.NewSymbol(ids[ix], names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.WildBomb))
		} else {
			ss[ix] = comp.NewSymbol(ids[ix], names[ix], resources[ix], payouts[ix], weights[ix])
		}
	}
	s := comp.NewSymbolSet(ss...)

	if symbols == nil {
		symbols = s
	}

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.WithSymbols(s),
		comp.DoubleSpin(),
		comp.WithPlayerChoice(),
		comp.MaxPayout(maxPayout),
		comp.WithSymbolsState(bomb),
		comp.WithRTP(target),
		comp.WithActions(actions, actions, nil, nil),
	)
}

func init() {
	initAction()

	slots92 = initSlots(92.0, weights92, gddPayouts, actions92all)
	slots94 = initSlots(94.0, weights94, gddPayouts, actions94all)
	slots96 = initSlots(96.0, weights96, gddPayouts, actions96all)

	slots92param = game.RegularParams{Slots: slots92}
	slots94param = game.RegularParams{Slots: slots94}
	slots96param = game.RegularParams{Slots: slots96}
}
