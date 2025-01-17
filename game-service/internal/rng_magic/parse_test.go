package rng_magic

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/mgd"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name     string
		function string
		fail     bool
	}{
		{
			name:     "empty",
			function: "",
			fail:     true,
		},
		{
			name:     "bad function",
			function: "max-payout(",
			fail:     true,
		},
		{
			name:     "double op",
			function: "max-payout() and and max-payout()",
			fail:     true,
		},
		{
			name:     "bool before function",
			function: "and max-payout()",
			fail:     true,
		},
		{
			name:     "double function",
			function: "not max-payout() max-payout()",
			fail:     true,
		},
		{
			name:     "NOT without function",
			function: "max-payout() AND NOT",
			fail:     true,
		},
		{
			name:     "paranthesis after func",
			function: "max-payout() (max-payout())",
			fail:     true,
		},
		{
			name:     "mixed ops",
			function: "max-payout() and max-payout() or max-payout()",
			fail:     true,
		},
		{
			name:     "no closing paranthesis",
			function: "(max-payout() and max-payout()",
			fail:     true,
		},
		{
			name:     "no opening paranthesis",
			function: "max-payout() and max-payout())",
			fail:     true,
		},
		{
			name:     "invalid function",
			function: "min-payout()",
			fail:     true,
		},
		{
			name:     "invalid param (1)",
			function: "first-symbol-count(l5,2)",
			fail:     true,
		},
		{
			name:     "invalid param (2)",
			function: "first-symbol-count(l3,x)",
			fail:     true,
		},
		{
			name:     "invalid param (3)",
			function: "first-symbol-count(l3,3,1,2,11)",
			fail:     true,
		},
		{
			name:     "good (1)",
			function: "spin-symbol-count(1,scatter,3,2,3,4,5)",
		},
		{
			name:     "good (2)",
			function: "spin-symbol-count(1,scatter,3,2,3,4,5) and max-payout()",
		},
		{
			name:     "good (3)",
			function: "spin-symbol-count(1,scatter,2,2,3) and spin-symbol-count(1,scatter,0,4,5)",
		},
		{
			name:     "good (4)",
			function: "magic-devil-bonus()",
		},
		{
			name:     "multiple and",
			function: "magic-devil-bonus() and spin-symbol-count(1,scatter,3,2,3,4,5) and max-payout()",
		},
		{
			name:     "multiple or",
			function: "magic-devil-bonus() or spin-symbol-count(1,scatter,3,2,3,4,5) or max-payout()",
		},
		{
			name:     "mixed and/or (1)",
			function: "(magic-devil-bonus() and spin-symbol-count(1,scatter,3,2,3,4,5)) or max-payout()",
		},
		{
			name:     "mixed and/or (2)",
			function: "magic-devil-bonus() and (spin-symbol-count(1,scatter,3,2,3,4,5) or max-payout())",
		},
		{
			name:     "mixed and/or (3)",
			function: "(magic-devil-bonus() and spin-symbol-count(1,scatter,3,2,3,4,5)) or (max-payout() and spin-symbol-count(1,scatter,3,2,3,4,5))",
		},
		{
			name:     "mixed or/and (1)",
			function: "(magic-devil-bonus() or spin-symbol-count(1,scatter,3,2,3,4,5)) and max-payout()",
		},
		{
			name:     "mixed or/and (2)",
			function: "magic-devil-bonus() or (spin-symbol-count(1,scatter,3,2,3,4,5) and max-payout())",
		},
		{
			name:     "mixed or/and (3)",
			function: "(magic-devil-bonus() or spin-symbol-count(1,scatter,3,2,3,4,5)) and (max-payout() or spin-symbol-count(1,scatter,3,2,3,4,5))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parse(tc.function, mgd.Conditions(), mgd.AllSymbols())
			if tc.fail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
