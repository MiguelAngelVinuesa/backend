package magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// ResultCountRange matches nr of results between min and max.
func ResultCountRange(min, max int) Matcher {
	return func(res results.Results) bool {
		l := len(res)
		return l >= min && (max < min || l <= max)
	}
}

// SpinResults instantiates a new array of results that contain only real spins.
// E.g. for Owl Kingdom it would filter out the special bonus results from the input slice.
func SpinResults(res results.Results) results.Results {
	out := make(results.Results, 0, len(res))
	for ix := range res {
		r := res[ix]
		if _, ok := r.Data.(*slots.SpinResult); ok {
			out = append(out, r)
		}
	}
	return out
}
