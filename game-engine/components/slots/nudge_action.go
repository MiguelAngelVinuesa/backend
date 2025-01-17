package slots

import (
	"math"
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NudgeAction is the action containing the details for a reel nudge.
type NudgeAction struct {
	SpinAction
	allowDupes bool          // indicates if the symbol can appear multiple times on a reel.
	count      uint8         // initial symbol count required to enable the trigger.
	location   NudgeLocation // location(s) where nudge can take place.
	chance     int           // the chance for the nudge.
	tease      int           // chance the nudge will be a teaser.
	reels      utils.UInt8s  // reels the nudge is limited to.
}

// NewNudgeAction instantiates a new nudge action.
func NewNudgeAction(symbol utils.Index, count uint8, location NudgeLocation, chance float64) *NudgeAction {
	a := &NudgeAction{
		allowDupes: true,
		count:      count,
		location:   location,
		chance:     int(math.Round(chance * 100)),
	}
	a.init(ReviseGrid, ReelsNudged, reflect.TypeOf(a).String())
	a.symbol = symbol
	return a
}

// WithTease indicates that the nudge result will be a teaser.
func (a *NudgeAction) WithTease(chance float64) *NudgeAction {
	a.tease = int(math.Round(chance * 100))
	return a
}

// GenerateNoDupes indicates the nudge must not generate duplicate symbols on the reels.
func (a *NudgeAction) GenerateNoDupes() *NudgeAction {
	a.allowDupes = false
	return a
}

// WithReels limits the nudge to the specified reels.
// Note that reels are 1-based, so the left most reel has id == 1!
func (a *NudgeAction) WithReels(reels ...uint8) *NudgeAction {
	a.reels = reels
	return a
}

// WithAlternate adds an alternative nudge action when this action doesn't trigger.
func (a *NudgeAction) WithAlternate(alt *NudgeAction) *NudgeAction {
	a.alternate = alt
	return a
}

// Triggered implements the SpinActioner.Triggered interface.
func (a *NudgeAction) Triggered(spin *Spin) SpinActioner {
	if count := spin.CountSymbol(a.symbol); count == a.count {
		if spin.prng.IntN(10000) < a.chance {
			return a
		}
	}
	return nil
}

// Nudge generates the reel nudge and modifies the result accordingly.
func (a *NudgeAction) Nudge(spin *Spin, result *SpinResult) SpinActioner {
	reels := spin.reelCount
	skip := 7
	if reels == 7 {
		skip = 9
	}

	reel := spin.prng.IntN(reels)
	if len(a.reels) > 0 {
		for {
			for ix := range a.reels {
				if a.reels[ix] == uint8(reel)+1 {
					if a.nudged(spin, result, reel) {
						return a
					}
					break
				}
			}
			reel = (reel + skip) % reels
		}
	} else {
		for {
			if a.nudged(spin, result, reel) {
				return a
			}
			reel = (reel + skip) % reels
		}
	}
}

func (a *NudgeAction) nudged(spin *Spin, result *SpinResult, reel int) bool {
	offset := reel * spin.rowCount
	if !a.allowDupes {
		max := offset + spin.rowCount
		for offs := offset; offs < max; offs++ {
			if spin.indexes[offs] == a.symbol {
				return false
			}
		}
	}

	location := a.location
	if location == NudgeVertical {
		if spin.prng.IntN(10000) < 5000 {
			location = NudgeBottom
		} else {
			location = NudgeTop
		}
	}

	offs2 := offset + spin.rowCount - 1
	if location == NudgeBottom {
		offset, offs2 = offs2, offset
	}
	if spin.indexes[offset] == a.symbol {
		return false
	}

	if l := len(spin.indexes); len(result.afterNudge) != l {
		result.afterNudge = result.afterNudge[:l]
		copy(result.afterNudge, spin.indexes)
	}

	tease := spin.prng.IntN(10000) < a.tease

	if !tease {
		if location == NudgeBottom {
			for ; offs2 < offset; offs2++ {
				result.afterNudge[offs2] = result.afterNudge[offs2+1]
			}
		} else {
			for ; offs2 > offset; offs2-- {
				result.afterNudge[offs2] = result.afterNudge[offs2-1]
			}
		}
		result.afterNudge[offset] = a.symbol
	}

	n := AcquireReelNudge(tease, uint8(reel)+1, 1, a.symbol, location)
	result.nudges = append(result.nudges, n)
	return true
}
