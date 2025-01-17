package magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// AND matches if all tests match.
func AND(tests ...Matcher) Matcher {
	return func(res results.Results) bool {
		for ix := range tests {
			if !tests[ix](res) {
				return false
			}
		}
		return true
	}
}

// OR matches if one of the tests matches.
func OR(tests ...Matcher) Matcher {
	return func(res results.Results) bool {
		for ix := range tests {
			if tests[ix](res) {
				return true
			}
		}
		return false
	}
}

// NOT matches if the test does not match.
func NOT(test Matcher) Matcher {
	return func(res results.Results) bool {
		return !test(res)
	}
}
