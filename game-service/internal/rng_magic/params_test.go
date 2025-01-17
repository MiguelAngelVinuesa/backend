package rng_magic

import (
	"testing"

	"github.com/stretchr/testify/require"

	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
)

func TestCheckParams(t *testing.T) {
	symbols := slots.NewSymbolSet(
		slots.NewSymbol(1, slots.WithName("l4"), slots.WithResource("low-4")),
		slots.NewSymbol(2, slots.WithName("l3"), slots.WithResource("low-3")),
		slots.NewSymbol(3, slots.WithName("l2"), slots.WithResource("low-2")),
		slots.NewSymbol(4, slots.WithName("l1"), slots.WithResource("low-1")),
		slots.NewSymbol(5, slots.WithName("h4"), slots.WithResource("high-4")),
		slots.NewSymbol(6, slots.WithName("h3"), slots.WithResource("high-3")),
		slots.NewSymbol(7, slots.WithName("h2"), slots.WithResource("high-2")),
		slots.NewSymbol(8, slots.WithName("h1"), slots.WithResource("high-1")),
		slots.NewSymbol(9, slots.WithName("w1"), slots.WithResource("wild-1"), slots.WithKind(slots.Wild)),
		slots.NewSymbol(10, slots.WithName("scat"), slots.WithResource("scatter-1"), slots.WithKind(slots.Scatter)),
	)

	testCases := []struct {
		name   string
		cond   *magic.Condition
		params []string
		fail   bool
	}{
		{
			name: "none - none",
			cond: magic.NewCondition("x"),
		},
		{
			name:   "none - 1 int",
			cond:   magic.NewCondition("x"),
			params: []string{"1"},
		},
		{
			name:   "none - 3 ints",
			cond:   magic.NewCondition("x"),
			params: []string{"1", "2", "3"},
		},
		{
			name: "1 string - none",
			cond: magic.NewCondition("x", magic.NewStringParam("s1")),
		},
		{
			name:   "1 string - 1 string",
			cond:   magic.NewCondition("x", magic.NewStringParam("s1")),
			params: []string{"1"},
		},
		{
			name:   "1 string - 3 strings",
			cond:   magic.NewCondition("x", magic.NewStringParam("s1")),
			params: []string{"1", "2", "3"},
		},
		{
			name: "1 int - none",
			cond: magic.NewCondition("x", magic.NewIntParam("i1", 0, 100)),
		},
		{
			name:   "1 int - 1 int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 0, 100)),
			params: []string{"1"},
		},
		{
			name:   "1 int - 3 ints",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 0, 100)),
			params: []string{"1", "2", "3"},
		},
		{
			name:   "1 int - bad int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 0, 100)),
			params: []string{"x"},
			fail:   true,
		},
		{
			name:   "1 int - low int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 1, 100)),
			params: []string{"0"},
			fail:   true,
		},
		{
			name:   "1 int - high int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 1, 100)),
			params: []string{"101"},
			fail:   true,
		},
		{
			name: "2 ints - none",
			cond: magic.NewCondition("x", magic.NewIntParam("i1", 0, 100), magic.NewIntParam("i2", 1, 100)),
		},
		{
			name:   "2 ints - 1 int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 0, 100), magic.NewIntParam("i2", 1, 100)),
			params: []string{"1"},
		},
		{
			name:   "2 ints - 3 ints",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 0, 100), magic.NewIntParam("i2", 1, 100)),
			params: []string{"1", "2", "3"},
		},
		{
			name:   "2 ints - bad int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 0, 100), magic.NewIntParam("i2", 1, 100)),
			params: []string{"10", "x"},
			fail:   true,
		},
		{
			name:   "2 ints - low int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 1, 100), magic.NewIntParam("i2", 1, 100)),
			params: []string{"10", "0"},
			fail:   true,
		},
		{
			name:   "2 ints - high int",
			cond:   magic.NewCondition("x", magic.NewIntParam("i1", 1, 100), magic.NewIntParam("i2", 1, 100)),
			params: []string{"10", "101"},
			fail:   true,
		},
		{
			name: "1 float - none",
			cond: magic.NewCondition("x", magic.NewFloatParam("i1", 0, 100)),
		},
		{
			name:   "1 float - 1 int",
			cond:   magic.NewCondition("x", magic.NewFloatParam("i1", 0, 100)),
			params: []string{"1"},
		},
		{
			name:   "1 float - 3 ints",
			cond:   magic.NewCondition("x", magic.NewFloatParam("i1", 0, 100)),
			params: []string{"1", "2", "3"},
		},
		{
			name:   "1 float - bad float",
			cond:   magic.NewCondition("x", magic.NewFloatParam("i1", 0, 100)),
			params: []string{"x"},
			fail:   true,
		},
		{
			name:   "1 float - low float",
			cond:   magic.NewCondition("x", magic.NewFloatParam("i1", 1, 100)),
			params: []string{"0.99"},
			fail:   true,
		},
		{
			name:   "1 float - high float",
			cond:   magic.NewCondition("x", magic.NewFloatParam("i1", 1, 100)),
			params: []string{"100.01"},
			fail:   true,
		},
		{
			name: "1 symbol - none",
			cond: magic.NewCondition("x", magic.NewSymbolParam()),
		},
		{
			name:   "1 symbol - correct id",
			cond:   magic.NewCondition("x", magic.NewSymbolParam()),
			params: []string{"4"},
		},
		{
			name:   "1 symbol - correct code",
			cond:   magic.NewCondition("x", magic.NewSymbolParam()),
			params: []string{"h2"},
		},
		{
			name:   "1 symbol - correct resource",
			cond:   magic.NewCondition("x", magic.NewSymbolParam()),
			params: []string{"low-2"},
		},
		{
			name:   "1 symbol - wild",
			cond:   magic.NewCondition("x", magic.NewSymbolParam()),
			params: []string{"wild"},
		},
		{
			name:   "1 symbol - scatter",
			cond:   magic.NewCondition("x", magic.NewSymbolParam()),
			params: []string{"scatter"},
		},
		{
			name: "reels last - none",
			cond: magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
		},
		{
			name:   "reels last - symbol",
			cond:   magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
			params: []string{"high-4"},
		},
		{
			name:   "reels last - symbol + reel",
			cond:   magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
			params: []string{"l1", "2"},
		},
		{
			name:   "reels last - symbol + 4 reels",
			cond:   magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
			params: []string{"l4", "1", "2", "3", "4"},
		},
		{
			name:   "reels last - bad reel",
			cond:   magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
			params: []string{"l4", "1", "2", "x", "4"},
			fail:   true,
		},
		{
			name:   "reels last - low reel",
			cond:   magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
			params: []string{"l4", "1", "2", "3", "0"},
			fail:   true,
		},
		{
			name:   "reels last - high reel",
			cond:   magic.NewCondition("x", magic.NewSymbolParam(), magic.NewReelsParam()),
			params: []string{"l4", "1", "11", "3", "4"},
			fail:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := checkParams(tc.cond, symbols, tc.params)
			if tc.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
