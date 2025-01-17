package slots

import (
	"math"
	"testing"

	rng2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	sf1  = slots.NewSymbol(1, slots.WithPayouts(0, 0, 3, 6, 12), slots.WithWeights(90, 70, 90, 70, 90, 70))
	sf2  = slots.NewSymbol(2, slots.WithPayouts(0, 0, 3, 6, 15), slots.WithWeights(90, 70, 90, 70, 90, 70))
	sf3  = slots.NewSymbol(3, slots.WithPayouts(0, 0, 6, 7.5, 15), slots.WithWeights(70, 90, 70, 90, 70, 90))
	sf4  = slots.NewSymbol(4, slots.WithPayouts(0, 0, 6, 12, 21), slots.WithWeights(70, 90, 70, 90, 70, 90))
	sf5  = slots.NewSymbol(5, slots.WithPayouts(0, 0, 9, 15, 24), slots.WithWeights(50, 70, 50, 70, 50, 70))
	sf6  = slots.NewSymbol(6, slots.WithPayouts(0, 0, 12, 18, 27), slots.WithWeights(50, 70, 50, 70, 50, 70))
	wff1 = slots.NewSymbol(7, slots.WildFor(0, 1), slots.WithPayouts(0, 0, 3, 6, 15), slots.WithWeights(20, 10, 20, 10, 20, 10))
	wff2 = slots.NewSymbol(8, slots.WildFor(2, 3), slots.WithPayouts(0, 0, 6, 12, 21), slots.WithWeights(10, 20, 10, 20, 10, 20))
	wf1  = slots.NewSymbol(9, slots.WithKind(slots.Wild), slots.WithPayouts(0, 0, 6, 15, 45), slots.WithWeights(0, 8, 8, 8, 0, 0))
	hf1  = slots.NewSymbol(10, slots.WithKind(slots.Hero), slots.WithWeights(2, 2, 2, 2, 2, 2))
	scf1 = slots.NewSymbol(11, slots.WithKind(slots.Scatter), slots.WithPayouts(0, 3, 6, 24, 60), slots.WithWeights(2, 2, 2, 2, 2, 2))
	scf2 = slots.NewSymbol(14, slots.WithKind(slots.Scatter), slots.WithPayouts(0, 3, 6, 24, 60), slots.WithWeights(2, 2, 2, 2, 2, 2))

	set1 = slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wf1, hf1, scf1)

	wheelWeights = utils.AcquireWeighting().AddWeights(utils.Indexes{1, 2, 3}, []float64{65, 25, 10})

	pl5x3x1 = slots.NewPayline(1, 3, 1, 1, 1, 1, 1)
	pl5x3x2 = slots.NewPayline(2, 3, 0, 0, 0, 0, 0)
	pl5x3x3 = slots.NewPayline(3, 3, 2, 2, 2, 2, 2)

	pl6x4x1 = slots.NewPayline(1, 4, 1, 1, 1, 1, 1, 1)
	pl6x4x2 = slots.NewPayline(2, 4, 0, 0, 0, 0, 0, 0)
	pl6x4x3 = slots.NewPayline(3, 4, 2, 2, 2, 2, 2, 2)

	t1  = slots.NewWildExpansion(1, false, 9, 1, false, false, false)
	t2  = slots.NewScatterFreeSpinsAction(1, false, 11, 2, false)
	t3  = slots.NewScatterBonusWheelAction(11, 3, 1, wheelWeights).WithAlternate(t2)
	t4  = slots.NewPaidAction(slots.BonusGame, 0, 100, scf1.ID(), 3)
	t5  = slots.NewWildExpansion(1, false, 9, 1, false, true, true)
	t6  = slots.NewWildExpansion(1, false, 9, 1, true, true, true)
	t7  = slots.NewWildExpansion(1, false, 9, 1, false, true, false)
	t8  = slots.NewWildExpansion(1, false, 9, 1, true, true, false)
	t9  = slots.NewScatterPayoutAction(11, 2, 3)
	t10 = slots.NewScatterPayoutAction(11, 3, 6).WithAlternate(t9)
	t11 = slots.NewWildPayoutAction(9, 3, 5)
	t12 = slots.NewScatterFreeSpinsAction(10, false, 11, 3, true).WithAlternate(t2)
	t13 = slots.NewScatterBonusWheelAction(14, 3, 2, wheelWeights)
	t14 = slots.NewPaylinesAction()
	t15 = slots.NewAllPaylinesAction(true)
	t16 = slots.NewBonusScatterAction(10)
	t17 = slots.NewScatterFreeSpinsAction(8, true, 11, 3, true)
	t18 = slots.NewScatterFreeSpinsAction(10, true, 11, 4, true).WithAlternate(t17)
	t19 = slots.NewScatterFreeSpinsAction(12, true, 11, 5, true).WithAlternate(t18)
	t20 = slots.NewPaylinesAction()
)

func TestNewRegular(t *testing.T) {
	testCases := []struct {
		name      string
		actions   slots.SpinActions
		state     bool
		double    bool
		choice    bool
		maxPayout float64
	}{
		{
			name:    "no actions",
			actions: slots.SpinActions{},
		},
		{
			name:    "1 action",
			actions: slots.SpinActions{t1},
		},
		{
			name:    "2 actions",
			actions: slots.SpinActions{t1, t2},
		},
		{
			name:    "paylines + 2 actions with alt",
			actions: slots.SpinActions{t14, t1, t3},
		},
		{
			name:    "all-paylines + 2 actions with alt",
			actions: slots.SpinActions{t15, t1, t3},
		},
		{
			name:    "all-paylines + 2 actions with alt + symbols state",
			actions: slots.SpinActions{t15, t1, t3},
			state:   true,
		},
		{
			name:    "no actions, double spin",
			actions: slots.SpinActions{},
			double:  true,
		},
		{
			name:    "no actions, player choice",
			actions: slots.SpinActions{},
			choice:  true,
		},
		{
			name:      "no actions, max payout",
			actions:   slots.SpinActions{},
			maxPayout: 3000.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := make([]slots.SlotOption, 0, 8)
			opts = append(opts, slots.Grid(5, 3),
				slots.WithSymbols(set1),
				slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3),
				slots.WithActions(tc.actions, nil, nil, nil),
				slots.MaxPayout(tc.maxPayout),
			)

			if tc.double {
				opts = append(opts, slots.DoubleSpin())
			}
			if tc.choice {
				opts = append(opts, slots.WithPlayerChoice())
			}
			if tc.state {
				opts = append(opts, slots.WithSymbolsState())
			}

			s := slots.NewSlots(opts...)
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			assert.Equal(t, s, r.slots)
			assert.NotNil(t, r.spin)
			assert.Nil(t, r.ForSale(1))
			assert.Nil(t, r.ForSale(2))
			assert.NotNil(t, r.LastSpin())
			assert.Equal(t, tc.double, r.IsDoubleSpin())
			assert.Equal(t, tc.choice, r.AllowPlayerChoices())
			assert.Empty(t, r.results)
			assert.Zero(t, r.totalPayout)
			assert.False(t, r.maxPayoutReached)

			if tc.state {
				assert.NotNil(t, r.symbolsState)
			} else {
				assert.Nil(t, r.symbolsState)
			}

			if tc.maxPayout > 0.0 {
				assert.Equal(t, tc.maxPayout, r.maxPayout)
			} else {
				assert.Equal(t, math.MaxFloat64, r.maxPayout)
			}
		})
	}
}

func TestNewRegularFail(t *testing.T) {
	t.Run("new regular fail", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		r := AcquireRegular(RegularParams{})
		require.Nil(t, r)
	})
}

func TestRegular_ForSale(t *testing.T) {
	a := slots.SpinActions{t4, t14}

	t.Run("new regular for sale", func(t *testing.T) {
		s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3), slots.WithActions(a, a, a, a))

		r := AcquireRegular(RegularParams{Slots: s})
		require.NotNil(t, r)
		defer r.Release()

		assert.NotNil(t, r.ForSale(1))

		res := r.Round(1)
		require.NotNil(t, res)
		assert.Equal(t, 1, len(res))
	})
}

func TestRegular_TestBeforeTriggers(t *testing.T) {
	testCases := []struct {
		name      string
		action    *slots.WildAction
		indexes   utils.Indexes
		before    slots.WildActions
		locked    utils.UInt8s
		freeSpins uint64
	}{
		{
			name:    "before no wild",
			action:  t5,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			before:  slots.WildActions{t5},
		},
		{
			name:      "before expand 1",
			action:    t5,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 9, 6, 1, 2, 3},
			before:    slots.WildActions{t5},
			locked:    utils.UInt8s{4},
			freeSpins: 1,
		},
		{
			name:      "before expand 2",
			action:    t5,
			indexes:   utils.Indexes{1, 2, 3, 4, 9, 6, 1, 2, 3, 4, 9, 6, 1, 2, 3},
			before:    slots.WildActions{t5},
			locked:    utils.UInt8s{2, 4},
			freeSpins: 1,
		},
		{
			name:      "before expand 3",
			action:    t5,
			indexes:   utils.Indexes{9, 2, 3, 4, 5, 6, 1, 9, 3, 4, 5, 6, 1, 2, 9},
			before:    slots.WildActions{t5},
			locked:    utils.UInt8s{1, 3, 5},
			freeSpins: 1,
		},
		{
			name:    "before+hero no wild",
			action:  t6,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			before:  slots.WildActions{t6},
		},
		{
			name:    "before+hero no hero",
			action:  t6,
			indexes: utils.Indexes{1, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			before:  slots.WildActions{t6},
		},
		{
			name:      "before+hero expand 1",
			action:    t6,
			indexes:   utils.Indexes{1, 9, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 10},
			before:    slots.WildActions{t6},
			locked:    utils.UInt8s{1},
			freeSpins: 1,
		},
		{
			name:      "before+hero expand 5",
			action:    t6,
			indexes:   utils.Indexes{1, 9, 3, 9, 5, 6, 1, 9, 3, 4, 5, 9, 1, 9, 10},
			before:    slots.WildActions{t6},
			locked:    utils.UInt8s{1, 2, 3, 4, 5},
			freeSpins: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr := tc.action
			tr.WithLockReels()

			a := slots.SpinActions{tr}

			s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithActions(a, a, a, a))
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			assert.EqualValues(t, tc.before, r.actionsFirst.expandBefore)
			assert.Zero(t, len(r.actionsFirst.expandAfter))
			assert.Zero(t, len(r.actionsFirst.extraPayouts))
			assert.Zero(t, len(r.actionsFirst.bonuses))

			r.spin.Debug(tc.indexes)

			r.spinData = slots.AcquireSpinResult(r.spin)
			r.currResult = results.AcquireResult(r.spinData, results.SpinData)
			defer r.currResult.Release()

			r.actionsFirst.testBeforeExpansions()

			r.locked = r.spin.Locked(r.locked)
			if tc.locked != nil {
				assert.EqualValues(t, tc.locked, r.locked)
			} else {
				assert.Zero(t, len(r.locked))
			}

			assert.Equal(t, tc.freeSpins, r.freeSpins)
			assert.Nil(t, r.bonusGame)
		})
	}
}

func TestRegular_TestAfterTriggers(t *testing.T) {
	testCases := []struct {
		name      string
		action    *slots.WildAction
		indexes   utils.Indexes
		after     slots.WildActions
		locked    utils.UInt8s
		freeSpins uint64
	}{
		{
			name:    "after no wild",
			action:  t7,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			after:   slots.WildActions{t7},
		},
		{
			name:      "after expand 1",
			action:    t7,
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 9, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			after:     slots.WildActions{t7},
			locked:    utils.UInt8s{2},
			freeSpins: 1,
		},
		{
			name:      "after expand 4",
			action:    t7,
			indexes:   utils.Indexes{9, 2, 9, 4, 5, 9, 1, 2, 3, 4, 9, 6, 9, 9, 3},
			after:     slots.WildActions{t7},
			locked:    utils.UInt8s{1, 2, 4, 5},
			freeSpins: 1,
		},
		{
			name:    "after+hero no wild",
			action:  t8,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
			after:   slots.WildActions{t8},
		},
		{
			name:    "after+hero no hero",
			action:  t8,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 9, 2, 3, 4, 5, 6, 1, 2, 3},
			after:   slots.WildActions{t8},
		},
		{
			name:      "after+hero expand 3",
			action:    t8,
			indexes:   utils.Indexes{1, 2, 9, 4, 5, 6, 9, 2, 3, 4, 5, 6, 1, 10, 9},
			after:     slots.WildActions{t8},
			locked:    utils.UInt8s{1, 3, 5},
			freeSpins: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tr := tc.action
			tr.WithLockReels()

			a := slots.SpinActions{tr}

			s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithActions(a, a, a, a))
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			assert.EqualValues(t, tc.after, r.actionsFirst.expandAfter)
			assert.Zero(t, len(r.actionsFirst.expandBefore))
			assert.Zero(t, len(r.actionsFirst.extraPayouts))
			assert.Zero(t, len(r.actionsFirst.bonuses))

			r.spin.Debug(tc.indexes)

			r.spinData = slots.AcquireSpinResult(r.spin)
			r.currResult = results.AcquireResult(r.spinData, results.SpinData)
			defer r.currResult.Release()

			r.actionsFirst.testRegularPayouts()
			r.actionsFirst.testAfterExpansions()

			r.locked = r.spin.Locked(r.locked)
			if tc.locked != nil {
				assert.EqualValues(t, tc.locked, r.locked)
			} else {
				assert.Zero(t, len(r.locked))
			}

			assert.Equal(t, tc.freeSpins, r.freeSpins)
			assert.Nil(t, r.bonusGame)
		})
	}
}

func TestRegular_TestPayoutTriggers(t *testing.T) {
	testCases := []struct {
		name      string
		actions   slots.SpinActions
		indexes   utils.Indexes
		payouts   results.Payouts
		freeSpins uint64
	}{
		{
			name:    "no wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "2 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 1, 2, 3, 4, 9, 6, 1, 2, 3},
		},
		{
			name:    "3 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 1, 2, 3, 4, 9, 6, 1, 2, 9},
			payouts: results.Payouts{
				slots.WildSymbolPayoutWithMap(6, 1, 9, 3, utils.UInt8s{3, 10, 14}),
			},
		},
		{
			name:    "4 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 1, 2, 9, 4, 9, 6, 1, 2, 9},
			payouts: results.Payouts{
				slots.WildSymbolPayoutWithMap(15, 1, 9, 4, utils.UInt8s{3, 8, 10, 14}),
			},
		},
		{
			name:    "5 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 1, 9, 9, 4, 9, 6, 1, 2, 9},
			payouts: results.Payouts{
				slots.WildSymbolPayoutWithMap(45, 1, 9, 5, utils.UInt8s{3, 7, 8, 10, 14}),
			},
		},
		{
			name:    "6 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 1, 9, 9, 4, 9, 6, 9, 2, 9},
			payouts: results.Payouts{
				slots.WildSymbolPayoutWithMap(45, 1, 9, 6, utils.UInt8s{3, 7, 8, 10, 12, 14}),
			},
		},
		{
			name:    "no wilds, 2 scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 11, 5, 6, 1, 2, 3, 4, 5, 6, 11, 2, 3},
			payouts: results.Payouts{
				slots.ScatterSymbolPayoutWithMap(3, 1, 11, 2, utils.UInt8s{3, 12}),
			},
		},
		{
			name:    "no wilds, 3 scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{11, 2, 3, 11, 5, 6, 1, 2, 3, 4, 5, 6, 11, 2, 3},
			payouts: results.Payouts{
				slots.ScatterSymbolPayoutWithMap(6, 1, 11, 3, utils.UInt8s{0, 3, 12}),
			},
		},
		{
			name:    "no wilds, 6 scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{11, 2, 3, 11, 5, 6, 11, 11, 3, 4, 5, 6, 11, 2, 11},
			payouts: results.Payouts{
				slots.ScatterSymbolPayoutWithMap(60, 1, 11, 6, utils.UInt8s{0, 3, 6, 7, 12, 14}),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithActions(tc.actions, nil, nil, nil))
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			assert.EqualValues(t, tc.actions, r.actionsFirst.extraPayouts)
			assert.Zero(t, len(r.actionsFirst.expandBefore))
			assert.Zero(t, len(r.actionsFirst.expandAfter))
			assert.Zero(t, len(r.actionsFirst.bonuses))

			r.spin.Debug(tc.indexes)

			r.spinData = slots.AcquireSpinResult(r.spin)
			r.currResult = results.AcquireResult(r.spinData, results.SpinData)
			defer r.currResult.Release()

			r.actionsFirst.testRegularPayouts()
			r.actionsFirst.testExtraPayouts()

			assert.Zero(t, len(r.locked))
			assert.Equal(t, tc.freeSpins, r.freeSpins)
			assert.Nil(t, r.bonusGame)
			assert.Equal(t, len(tc.payouts), len(r.currResult.Payouts))
		})
	}
}

func TestRegular_TestBonusSymbol(t *testing.T) {
	a := slots.SpinActions{slots.NewBonusScatterAction(10)}
	s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithActions(a, nil, nil, nil))

	testCases := []struct {
		name    string
		symbol  utils.Index
		indexes utils.Indexes
		payout  results.Payout
	}{
		{
			name:    "no symbol",
			symbol:  utils.MaxIndex,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "invalid symbol",
			symbol:  15,
			indexes: utils.Indexes{1, 2, 3, 4, 15, 6, 15, 2, 3, 4, 15, 6, 1, 2, 3},
		},
		{
			name:    "no payout",
			symbol:  2,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 3, 3, 4, 5, 6, 1, 10, 3},
		},
		{
			name:    "payout 3 symbols",
			symbol:  4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 4, 2, 3, 4, 5, 6, 1, 2, 3},
			payout:  slots.BonusSymbolPayout(60, 1, 4, 3),
		},
		{
			name:    "payout 4 symbols",
			symbol:  5,
			indexes: utils.Indexes{1, 5, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 5, 3},
			payout:  slots.BonusSymbolPayout(150, 1, 5, 4),
		},
		{
			name:    "payout 5 symbols",
			symbol:  6,
			indexes: utils.Indexes{6, 2, 3, 4, 5, 6, 1, 6, 3, 4, 5, 6, 1, 2, 6},
			payout:  slots.BonusSymbolPayout(270, 1, 6, 5),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			r.spin.ResetSpin()
			r.spin.Debug(tc.indexes)
			r.spin.SetBonusSymbol(tc.symbol, false)

			r.spinData = slots.AcquireSpinResult(r.spin)
			r.currResult = results.AcquireResult(r.spinData, results.SpinData)
			defer r.currResult.Release()

			r.actionsFirst.testRegularPayouts()
			r.actionsFirst.testBonuses()

			if tc.payout == nil {
				require.Zero(t, len(r.currResult.Payouts))
			} else {
				require.Equal(t, 1, len(r.currResult.Payouts))
			}
		})
	}
}

func TestRegular_TestOtherTriggers(t *testing.T) {
	w := utils.AcquireWeighting()
	defer w.Release()
	w.AddWeight(sf1.ID(), 120)
	w.AddWeight(sf2.ID(), 100)
	w.AddWeight(sf3.ID(), 80)
	w.AddWeight(sf4.ID(), 60)
	w.AddWeight(sf5.ID(), 40)
	w.AddWeight(sf6.ID(), 20)

	set2 := slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wf1, hf1, scf1, scf2).SetBonusWeights(w)

	testCases := []struct {
		name        string
		actions     slots.SpinActions
		indexes     utils.Indexes
		bonusGame   bool
		freeSpins   uint64
		bonusSymbol bool
	}{
		{
			name:    "nothing",
			actions: slots.SpinActions{t12, t13},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "bonusBuy game",
			actions:   slots.SpinActions{t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 14, 4, 5, 6, 14, 14, 3},
			bonusGame: true,
		},
		{
			name:      "free spins",
			actions:   slots.SpinActions{t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 11, 4, 5, 6, 11, 14, 3},
			freeSpins: 1,
		},
		{
			name:        "free spins with bonusBuy symbol",
			actions:     slots.SpinActions{t12, t13},
			indexes:     utils.Indexes{1, 2, 3, 4, 11, 6, 1, 2, 11, 4, 5, 6, 11, 14, 3},
			freeSpins:   10,
			bonusSymbol: true,
		},
		{
			name:      "bonusBuy game + free spins",
			actions:   slots.SpinActions{t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 11, 5, 6, 11, 2, 14, 4, 5, 6, 14, 14, 3},
			bonusGame: true,
			freeSpins: 1,
		},
		{
			name:        "bonusBuy game + free spins with bonusBuy symbol",
			actions:     slots.SpinActions{t12, t13},
			indexes:     utils.Indexes{14, 2, 3, 11, 5, 6, 11, 2, 14, 4, 11, 6, 14, 14, 11},
			bonusGame:   true,
			freeSpins:   10,
			bonusSymbol: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set2), slots.WithActions(tc.actions, nil, nil, nil))
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			assert.EqualValues(t, tc.actions, r.actionsFirst.bonuses)

			r.spin.Debug(tc.indexes)
			r.prepareRound(0)

			r.spinData = slots.AcquireSpinResult(r.spin)
			r.currResult = results.AcquireResult(r.spinData, results.SpinData)
			defer r.currResult.Release()

			r.actionsFirst.testBonuses()

			if tc.bonusGame {
				assert.NotNil(t, r.bonusGame)
			} else {
				assert.Nil(t, r.bonusGame)
			}

			assert.Equal(t, tc.freeSpins, r.freeSpins)

			if tc.bonusSymbol {
				assert.NotZero(t, r.BonusSymbol())
			} else {
				assert.Equal(t, utils.MaxIndex, r.BonusSymbol())
			}
		})
	}
}

func TestRegular_TestActions(t *testing.T) {
	w := utils.AcquireWeighting()
	defer w.Release()
	w.AddWeight(sf1.ID(), 120)
	w.AddWeight(sf2.ID(), 100)
	w.AddWeight(sf3.ID(), 80)
	w.AddWeight(sf4.ID(), 60)
	w.AddWeight(sf5.ID(), 40)
	w.AddWeight(sf6.ID(), 20)

	set2 := slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wf1, hf1, scf1, scf2).SetBonusWeights(w)

	testCases := []struct {
		name        string
		actions     slots.SpinActions
		indexes     utils.Indexes
		bonus       utils.Index
		locked      utils.UInt8s
		payouts     results.Payouts
		bonusGame   bool
		freeSpins   uint64
		bonusSymbol bool
	}{
		{
			name:    "no extraPayouts",
			actions: slots.SpinActions{t7, t10, t11, t12, t13, t14},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "1 payline",
			actions: slots.SpinActions{t7, t10, t11, t12, t13, t14},
			indexes: utils.Indexes{1, 2, 3, 4, 2, 6, 1, 2, 3, 4, 2, 6, 1, 2, 3},
			payouts: results.Payouts{
				slots.WinlinePayoutFromData(5, 0, 2, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
			},
		},
		{
			name:    "3 paylines",
			actions: slots.SpinActions{t7, t10, t11, t12, t13, t14},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 6, 4, 2, 3},
			payouts: results.Payouts{
				slots.WinlinePayoutFromData(5, 0, 2, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
				slots.WinlinePayoutFromData(2, 0, 1, 4, slots.PayLTR, 2, utils.UInt8s{0, 0, 0, 0, 0}),
				slots.WinlinePayoutFromData(2, 0, 3, 3, slots.PayLTR, 3, utils.UInt8s{2, 2, 2, 2, 2}),
			},
		},
		{
			name:    "expand before no wild",
			actions: slots.SpinActions{t5},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "expand before no wild",
			actions: slots.SpinActions{t5},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "expand before 1",
			actions:   slots.SpinActions{t5},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 9, 6, 1, 2, 3},
			locked:    utils.UInt8s{4},
			freeSpins: 1,
		},
		{
			name:    "expand after no wild",
			actions: slots.SpinActions{t7},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "expand after 2",
			actions:   slots.SpinActions{t7},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 9, 4, 5, 6, 4, 5, 6, 1, 2, 9},
			locked:    utils.UInt8s{2, 5},
			freeSpins: 1,
		},
		{
			name:    "no wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "2 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 5, 6, 4, 4, 9, 6, 1, 2, 3},
		},
		{
			name:    "3 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 6, 4, 5, 4, 9, 6, 1, 2, 9},
			payouts: results.Payouts{
				slots.WildSymbolPayoutWithMap(6, 1, 9, 3, utils.UInt8s{3, 10, 14}),
			},
		},
		{
			name:    "no wilds, 6 scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{11, 2, 3, 5, 11, 6, 1, 11, 11, 4, 5, 6, 11, 2, 11},
			payouts: results.Payouts{
				slots.ScatterSymbolPayoutWithMap(60, 1, 11, 6, utils.UInt8s{0, 4, 7, 8, 12, 14}),
			},
		},
		{
			name:        "bonusBuy no payout",
			bonus:       2,
			indexes:     utils.Indexes{1, 2, 3, 4, 5, 6, 1, 3, 3, 4, 5, 6, 1, 10, 3},
			bonusSymbol: true,
		},
		{
			name:    "bonusBuy payout 3 symbols",
			actions: slots.SpinActions{t16},
			bonus:   4,
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 4, 2, 3, 4, 5, 6, 1, 2, 3},
			payouts: results.Payouts{
				slots.BonusSymbolPayout(6, 1, 4, 3),
			},
			bonusSymbol: true,
		},
		{
			name:      "bonusBuy game",
			actions:   slots.SpinActions{t7, t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 14, 4, 5, 6, 14, 14, 3},
			bonusGame: true,
		},
		{
			name:      "free spins",
			actions:   slots.SpinActions{t7, t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 11, 4, 5, 6, 11, 14, 3},
			freeSpins: 1,
		},
		{
			name:        "free spins with bonusBuy symbol",
			actions:     slots.SpinActions{t7, t12, t13},
			indexes:     utils.Indexes{1, 2, 3, 4, 11, 6, 1, 2, 11, 4, 5, 6, 11, 14, 3},
			freeSpins:   10,
			bonusSymbol: true,
		},
		{
			name:      "bonusBuy game + free spins",
			actions:   slots.SpinActions{t7, t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 11, 5, 6, 11, 2, 14, 4, 5, 6, 14, 14, 3},
			bonusGame: true,
			freeSpins: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set2), slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3), slots.WithActions(tc.actions, nil, nil, nil))
			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			r.prepareRound(0)
			r.spin.Debug(tc.indexes)

			if tc.bonus > 0 {
				r.spin.SetBonusSymbol(tc.bonus, false)
			}

			r.spinData = slots.AcquireSpinResult(r.spin)
			r.currResult = results.AcquireResult(r.spinData, results.SpinData)
			r.currResult.Release()

			r.testActions()

			if tc.locked != nil {
				assert.EqualValues(t, tc.locked, r.locked)
			} else {
				assert.Zero(t, len(r.locked))
			}

			assert.Equal(t, len(tc.payouts), len(r.currResult.Payouts))
			assert.Equal(t, tc.freeSpins, r.freeSpins)

			if tc.bonusGame {
				assert.NotNil(t, r.bonusGame)
			} else {
				assert.Nil(t, r.bonusGame)
			}

			if tc.bonusSymbol {
				assert.NotZero(t, r.BonusSymbol())
			} else {
				assert.Equal(t, utils.MaxIndex, r.BonusSymbol())
			}
		})
	}
}

func TestRegular_GetResult(t *testing.T) {
	t.Run("regular get results", func(t *testing.T) {
		w := utils.AcquireWeighting()
		defer w.Release()
		w.AddWeight(sf1.ID(), 120)
		w.AddWeight(sf2.ID(), 100)
		w.AddWeight(sf3.ID(), 80)
		w.AddWeight(sf4.ID(), 60)
		w.AddWeight(sf5.ID(), 40)
		w.AddWeight(sf6.ID(), 20)

		set2 := slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wf1, hf1, scf1, scf2).SetBonusWeights(w)

		s := slots.NewSlots(
			slots.Grid(5, 3),
			slots.WithSymbols(set2),
			slots.PayDirections(slots.PayBoth),
			slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3),
			slots.WithActions(slots.SpinActions{t7, t12, t13, t14}, nil, nil, nil),
		)

		r := AcquireRegular(RegularParams{Slots: s})
		require.NotNil(t, r)
		defer r.Release()

		for ix := 0; ix < 100; ix++ {
			r.prepareRound(0)
			r.spin.Spin()
			r.getResults()
			assert.GreaterOrEqual(t, len(r.results), 1)
		}
	})
}

func TestRegular_PrngLog(t *testing.T) {
	t.Run("regular get prng log", func(t *testing.T) {
		w := utils.AcquireWeighting()
		defer w.Release()
		w.AddWeight(sf1.ID(), 120)
		w.AddWeight(sf2.ID(), 100)
		w.AddWeight(sf3.ID(), 80)
		w.AddWeight(sf4.ID(), 60)
		w.AddWeight(sf5.ID(), 40)
		w.AddWeight(sf6.ID(), 20)

		set2 := slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wf1, hf1, scf1, scf2).SetBonusWeights(w)

		s := slots.NewSlots(
			slots.Grid(5, 3),
			slots.WithSymbols(set2),
			slots.PayDirections(slots.PayBoth),
			slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3),
			slots.WithActions(slots.SpinActions{t7, t12, t13, t14}, nil, nil, nil),
		)

		r := AcquireRegular(RegularParams{Slots: s, PrngLog: true})
		require.NotNil(t, r)
		defer r.Release()

		r.prepareRound(0)
		r.spin.Spin()
		r.getResults()
		assert.NotEmpty(t, r.results)

		for ix := range r.results {
			res := r.results[ix]
			require.NotNil(t, res)

			data, ok := res.Data.(*slots.SpinResult)
			require.True(t, ok)
			require.NotNil(t, data)

			_, l1, l2 := data.Log()
			require.NotEmpty(t, l1)
			require.NotEmpty(t, l2)
		}
	})
}

func TestRegular_Debug(t *testing.T) {
	w := utils.AcquireWeighting()
	defer w.Release()
	w.AddWeight(sf1.ID(), 120)
	w.AddWeight(sf2.ID(), 100)
	w.AddWeight(sf3.ID(), 80)
	w.AddWeight(sf4.ID(), 60)
	w.AddWeight(sf5.ID(), 40)
	w.AddWeight(sf6.ID(), 20)

	set2 := slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, wff1, wff2, wf1, hf1, scf1, scf2).SetBonusWeights(w)

	testCases := []struct {
		name        string
		actions     slots.SpinActions
		indexes     utils.Indexes
		locked      bool
		freeSpins   bool
		payouts     bool
		bonusGame   bool
		bonusSymbol bool
	}{
		{
			name:    "no payouts",
			actions: slots.SpinActions{t7, t10, t11, t12, t13, t14},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "1 payout",
			actions: slots.SpinActions{t7, t10, t11, t12, t13, t14},
			indexes: utils.Indexes{1, 2, 3, 4, 2, 6, 1, 2, 3, 4, 2, 6, 1, 2, 3},
			payouts: true,
		},
		{
			name:    "3 payouts",
			actions: slots.SpinActions{t7, t10, t11, t12, t13, t14},
			indexes: utils.Indexes{1, 2, 3, 1, 2, 3, 1, 2, 3, 1, 2, 6, 4, 2, 3},
			payouts: true,
		},
		{
			name:    "expand before no wild",
			actions: slots.SpinActions{t5},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "expand before 1",
			actions:   slots.SpinActions{t5},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 4, 3, 4, 9, 6, 1, 2, 3},
			locked:    true,
			freeSpins: true,
		},
		{
			name:    "expand after no wild",
			actions: slots.SpinActions{t7},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:      "expand after 2",
			actions:   slots.SpinActions{t7},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 9, 6, 5, 4, 4, 5, 6, 1, 2, 9},
			locked:    true,
			freeSpins: true,
		},
		{
			name:    "no wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3},
		},
		{
			name:    "2 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 5, 6, 4, 4, 9, 6, 1, 2, 3},
		},
		{
			name:    "3 wilds, no scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{1, 2, 3, 9, 5, 6, 6, 4, 5, 4, 9, 6, 1, 2, 9},
			payouts: true,
		},
		{
			name:    "no wilds, 6 scatters",
			actions: slots.SpinActions{t10, t11},
			indexes: utils.Indexes{11, 2, 3, 5, 11, 6, 1, 11, 11, 4, 5, 6, 11, 2, 11},
			payouts: true,
		},
		{
			name:      "bonusBuy game",
			actions:   slots.SpinActions{t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 14, 4, 5, 6, 14, 14, 3},
			bonusGame: true,
		},
		{
			name:      "free spins",
			actions:   slots.SpinActions{t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 11, 4, 5, 6, 11, 14, 3},
			freeSpins: true,
		},
		{
			name:        "free spins with bonusBuy symbol",
			actions:     slots.SpinActions{t12, t13},
			indexes:     utils.Indexes{1, 2, 3, 4, 11, 6, 1, 2, 11, 4, 5, 6, 11, 14, 3},
			bonusSymbol: true,
			freeSpins:   true,
		},
		{
			name:      "bonusBuy game + free spins",
			actions:   slots.SpinActions{t12, t13},
			indexes:   utils.Indexes{1, 2, 3, 11, 5, 6, 11, 2, 14, 4, 5, 6, 14, 14, 3},
			bonusGame: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := slots.NewSlots(
				slots.Grid(5, 3),
				slots.WithSymbols(set2),
				slots.PayDirections(slots.PayBoth),
				slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3),
				slots.WithActions(tc.actions, nil, nil, nil),
			)

			r := AcquireRegular(RegularParams{Slots: s})
			require.NotNil(t, r)
			defer r.Release()

			res := r.Debug(tc.indexes, 0, false, nil)
			require.NotNil(t, res)
			require.NotEmpty(t, res)

			if tc.locked {
				assert.NotZero(t, len(r.locked))
			}
			if tc.freeSpins {
				assert.Greater(t, len(res), 1)
			}
			if tc.payouts {
				assert.NotZero(t, len(res[0].Payouts))
			}

			assert.Equal(t, tc.bonusGame, r.bonusGamePlayed)

			if tc.bonusSymbol {
				assert.NotZero(t, r.BonusSymbol())
			}
		})
	}
}

func TestRegular_AltSymbols(t *testing.T) {
	t.Run("alternate symbols", func(t *testing.T) {
		w := utils.AcquireWeighting()
		defer w.Release()
		w.AddWeight(sf1.ID(), 120)
		w.AddWeight(sf2.ID(), 100)
		w.AddWeight(sf3.ID(), 80)
		w.AddWeight(sf4.ID(), 60)
		w.AddWeight(sf5.ID(), 40)
		w.AddWeight(sf6.ID(), 20)

		a := slots.SpinActions{t19}

		set2 := slots.NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, scf1, scf2).SetBonusWeights(w)
		s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set2), slots.WithAltSymbols(set1), slots.WithActions(a, nil, nil, nil))

		r := AcquireRegular(RegularParams{Slots: s})
		require.NotNil(t, r)
		defer r.Release()

		indexes := utils.Indexes{11, 2, 3, 4, 5, 6, 11, 2, 3, 4, 5, 6, 11, 2, 3}
		res := r.Debug(indexes, 0, false, nil)
		require.NotNil(t, res)

		assert.NotEqual(t, utils.MaxIndex, r.BonusSymbol())
		assert.GreaterOrEqual(t, len(res), 9)
	})
}

func BenchmarkRegular_Round5x3(b *testing.B) {
	a := slots.SpinActions{t14}
	s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3), slots.WithActions(a, nil, nil, nil))

	r := AcquireRegular(RegularParams{Slots: s})
	defer r.Release()

	for i := 0; i < b.N; i++ {
		r.Round(0)
	}
}

func BenchmarkRegular_Round5x3debug(b *testing.B) {
	a := slots.SpinActions{t14}
	s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3), slots.WithActions(a, nil, nil, nil))

	r := AcquireRegular(RegularParams{Slots: s})
	defer r.Release()

	for i := 0; i < b.N; i++ {
		r.Debug(utils.Indexes{1, 2, 3, 1, 4, 5, 1, 2, 3, 1, 4, 5, 1, 2, 3}, 0, false, nil)
	}
}

func BenchmarkRegular_Round5x3noRepeat(b *testing.B) {
	a := slots.SpinActions{t14}
	s := slots.NewSlots(slots.Grid(5, 3), slots.WithSymbols(set1), slots.NoRepeat(2), slots.WithPaylines(slots.PayLTR, false, pl5x3x1, pl5x3x2, pl5x3x3), slots.WithActions(a, nil, nil, nil))

	r := AcquireRegular(RegularParams{Slots: s})
	defer r.Release()

	for i := 0; i < b.N; i++ {
		r.Round(0)
	}
}

func BenchmarkRegular_Round6x4(b *testing.B) {
	a := slots.SpinActions{t20}
	s := slots.NewSlots(slots.Grid(6, 4), slots.WithSymbols(set1), slots.WithPaylines(slots.PayLTR, false, pl6x4x1, pl6x4x2, pl6x4x3), slots.WithActions(a, nil, nil, nil))

	r := AcquireRegular(RegularParams{Slots: s})
	defer r.Release()

	for i := 0; i < b.N; i++ {
		r.Round(0)
	}
}

func BenchmarkRegular_Round6x4noRepeat(b *testing.B) {
	a := slots.SpinActions{t20}
	s := slots.NewSlots(slots.Grid(6, 4), slots.WithSymbols(set1), slots.NoRepeat(3), slots.WithPaylines(slots.PayLTR, false, pl6x4x1, pl6x4x2, pl6x4x3), slots.WithActions(a, nil, nil, nil))

	r := AcquireRegular(RegularParams{Slots: s})
	defer r.Release()

	for i := 0; i < b.N; i++ {
		r.Round(0)
	}
}

func init() {
	// use internal RNG for unit-tests
	rng.AcquireRNG = func() interfaces.Generator { return rng2.NewRNG() }
}
