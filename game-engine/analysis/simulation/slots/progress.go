package magic

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// ProgressMeter matches if the overal round progress meter on the last spin is between min and max.
func ProgressMeter(min, max int) Matcher {
	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < 0 {
			return false
		}

		if r, ok := res[len(res)-1].Data.(*comp.SpinResult); ok {
			m := r.ProgressLevel()
			return m >= min && (max < min || m <= max)
		}
		return false
	}
}
