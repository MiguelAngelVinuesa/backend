package slots

import (
	"bytes"
	"reflect"
	"strconv"
)

type ProgressAction struct {
	SpinAction
	payoutSymbols bool // indicates the progress level is increased by the number of symbols in the payouts.
	reset         bool // indicates the progress meter must be reset.
	startLevel    int  // starting progress level when first triggered.
	maxLevel      int  // maximum level allowed.
}

// NewResetProgressAction instantiates a new action to reset the progress meter.
func NewResetProgressAction(startLevel int) *ProgressAction {
	a := &ProgressAction{}
	a.init(ReviseGrid, Multiplier, reflect.TypeOf(a).String())
	a.reset = true
	a.startLevel = startLevel
	return a.finalizer()
}

// NewPayoutSymbolProgress instantiates a new action to increase the overall progress level by the number of symbols in the payouts.
func NewPayoutSymbolProgress(startLevel, maxLevel int) *ProgressAction {
	a := &ProgressAction{}
	a.init(AwardBonuses, Multiplier, reflect.TypeOf(a).String())
	a.payoutSymbols = true
	a.startLevel = startLevel
	a.maxLevel = maxLevel
	return a.finalizer()
}

// Triggered implements the SpinActioner.Triggered interface.
func (a *ProgressAction) Triggered(spin *Spin) SpinActioner {
	if a.reset {
		spin.progressLevel = a.startLevel
		return a
	}

	var count int
	for ix := range spin.payouts {
		switch spin.payouts[ix] {
		case 0:
		case 1:
			count++
		default:
			count += int(spin.payouts[ix])
			count--
		}
	}
	if count == 0 {
		return nil
	}

	if spin.progressLevel <= 0 {
		spin.progressLevel = a.startLevel
	}
	spin.progressLevel += count

	if a.maxLevel > 0 && spin.progressLevel > a.maxLevel {
		spin.progressLevel = a.maxLevel
	}

	return a
}

func (a *ProgressAction) finalizer() *ProgressAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.reset:
		b.WriteString(",reset=true")
		b.WriteString(",start=")
		b.WriteString(strconv.Itoa(a.startLevel))

	case a.payoutSymbols:
		b.WriteString(",payoutSymbols=true")
		b.WriteString(",start=")
		b.WriteString(strconv.Itoa(a.startLevel))
		b.WriteString(",max=")
		b.WriteString(strconv.Itoa(a.maxLevel))
	}

	a.config = b.String()
	return a
}
