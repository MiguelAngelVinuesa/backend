package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// ClearAction is an action that when triggered clears some positions in the grid of the spin.
type ClearAction struct {
	SpinAction
	clearPayouts bool
	explodeShape GridOffsets
}

// NewClearPayoutsAction instantiates a new clear winlines action.
// This action is generally used in conjunction with cascading reels.
func NewClearPayoutsAction() *ClearAction {
	a := newClearAction()
	a.clearPayouts = true
	return a.finalize()
}

// NewExplodingBombsAction instantiates a new exploding bombs action.
func NewExplodingBombsAction(symbol utils.Index, shape GridOffsets) *ClearAction {
	a := newClearAction()
	a.symbol = symbol
	a.explodeShape = shape
	return a.finalize()
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *ClearAction) Triggered(s *Spin) SpinActioner {
	if a.clearPayouts && a.doClearPayouts(s) {
		return a
	}
	if a.explodeShape != nil && a.doExplodeBombs(s) {
		return a
	}
	return nil
}

// doClearPayouts removes all symbols that are responsible for a payout.
func (a *ClearAction) doClearPayouts(spin *Spin) bool {
	var cleared bool
	for ix := range spin.payouts {
		if spin.payouts[ix] == 1 {
			spin.indexes[ix] = 0
			cleared = true
		}
	}
	return cleared
}

// doExplodeBombs finds all bomb symbols in the spin result and executes their explosion.
func (a *ClearAction) doExplodeBombs(spin *Spin) bool {
	reels, rows := spin.reelCount, spin.rowCount
	var offset int
	var exploded bool
	for reel := 0; reel < reels; reel++ {
		for row := 0; row < rows; row++ {
			if spin.indexes[offset] == a.symbol {
				exploded = true
				for _, offsets := range a.explodeShape {
					reel2, row2 := reel+offsets[0], row+offsets[1]
					if reel2 >= 0 && reel2 < reels && row2 >= 0 && row2 < rows {
						offset2 := reel2*rows + row2
						if (reel2 == reel && row2 == row) || spin.indexes[offset2] != a.symbol {
							spin.indexes[offset2] = 0
						}
					}
				}
			}
			offset++
		}
	}
	return exploded
}

func (a *ClearAction) finalize() *ClearAction {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	if a.clearPayouts {
		b.WriteString(",clearPayouts=true")
	} else {
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",shape=")
		j, _ := json.Marshal(a.explodeShape)
		b.Write(j)
	}

	a.config = b.String()
	return a
}

func newClearAction() *ClearAction {
	a := &ClearAction{}
	a.init(TestClearance, Refill, reflect.TypeOf(a).String())
	return a
}
