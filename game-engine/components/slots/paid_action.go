package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// PaidAction is an action activated by a separate payment.
// The amount deducted is the current bet amount * betMultiplier of the action.
type PaidAction struct {
	SpinAction
	bonusKind     uint8
	triggerCount  uint8
	betMultiplier int
	flag          int
	flagValue     int
	reels         utils.UInt8s
}

// NewPaidAction instantiates a new paid action.
func NewPaidAction(kind SpinActionResult, nrOfSpins uint8, betMultiplier int, symbol utils.Index, triggerCount uint8, reels ...uint8) *PaidAction {
	a := &PaidAction{}
	a.init(PaidOnly, kind, reflect.TypeOf(a).String())
	a.nrOfSpins = nrOfSpins
	a.symbol = symbol
	a.betMultiplier = betMultiplier
	a.triggerCount = triggerCount
	a.flag = -1
	a.reels = reels
	return a.finalizer()
}

// WithFlag sets a flag when the action triggers.
func (a *PaidAction) WithFlag(flag, flagValue int) *PaidAction {
	a.flag = flag
	a.flagValue = flagValue
	return a.finalizer()
}

// WithBonusKind sets the bonus buy kind.
func (a *PaidAction) WithBonusKind(bonusKind uint8) *PaidAction {
	a.bonusKind = bonusKind
	return a.finalizer()
}

// UpdateFlag updates the round flag if possible.
func (a *PaidAction) UpdateFlag(spin *Spin) SpinActioner {
	if a.flag >= 0 && a.flag < len(spin.roundFlags) {
		spin.roundFlags[a.flag] = a.flagValue
		return a
	}
	return nil
}

// BonusKind returns the bonus buy kind.
func (a *PaidAction) BonusKind() uint8 {
	return a.bonusKind
}

// BetMultiplier returns the bet multiplier for the bonus buy feature.
func (a *PaidAction) BetMultiplier() int {
	return a.betMultiplier
}

func (a *PaidAction) finalizer() *PaidAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	if a.bonusKind > 0 {
		b.WriteString(",bonusKind=")
		b.WriteString(strconv.Itoa(int(a.bonusKind)))
	}
	if a.nrOfSpins > 0 {
		b.WriteString(",nrOfSpins=")
		b.WriteString(strconv.Itoa(int(a.nrOfSpins)))
	}
	if a.symbol > 0 {
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
	}
	if a.betMultiplier > 0 {
		b.WriteString(",betMultiplier=")
		b.WriteString(strconv.Itoa(a.betMultiplier))
	}
	if a.triggerCount > 0 {
		b.WriteString(",triggerCount=")
		b.WriteString(strconv.Itoa(int(a.triggerCount)))
	}
	if a.flag > 0 {
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.flag))
	}
	if a.flagValue > 0 {
		b.WriteString(",flagValue=")
		b.WriteString(strconv.Itoa(a.flagValue))
	}
	if len(a.reels) > 0 {
		b.WriteString(",reels=")
		j, _ := json.Marshal(a.reels)
		b.Write(j)
	}

	a.config = b.String()
	return a
}
