package magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// PlayerChoice matches if the player made the indicated choice.
func PlayerChoice(code, choice string) Matcher {
	return func(res results.Results) bool {
		for ix := 0; ix < 5 && ix < len(res); ix++ {
			if r, ok := res[ix].Data.(*slots.SpinResult); ok {
				if r.PlayerChoices.Choices()[code] == choice {
					return true
				}
			}
		}
		return false
	}
}
