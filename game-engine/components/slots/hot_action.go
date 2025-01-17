package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// HotAction represents an action to determine hot reels.
type HotAction struct {
	SpinAction
}

// NewHotAction instantiates a new hot action.
func NewHotAction(symbol utils.Index) *HotAction {
	a := &HotAction{}
	a.init(AwardBonuses, HotReel, reflect.TypeOf(a).String())
	a.symbol = symbol
	return a.finalize()
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *HotAction) Triggered(s *Spin) SpinActioner {
	var found bool
	for reel := 0; reel < s.reelCount; reel++ {
		offset := reel * s.rowCount
		for row := 0; row < s.rowCount; row++ {
			if s.indexes[offset] == a.symbol {
				s.HotReel(uint8(reel))
				found = true
				break
			}
			offset++
		}
	}

	if found {
		return a
	}
	return nil
}

func (a *HotAction) finalize() *HotAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())
	b.WriteString(",symbol=")
	b.WriteString(strconv.Itoa(int(a.symbol)))
	a.config = b.String()
	return a
}
