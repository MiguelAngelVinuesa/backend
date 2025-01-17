package slots

import (
	"bytes"
	"math"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// MultiplierAction represents the action to manipulate symbol multipliers on the grid.
type MultiplierAction struct {
	SpinAction
	gridMultipliers bool
	triggerSymbol   utils.Index
	multiply        uint16
	max             uint16
	firstMultiplier bool
	firstLevel      int
	multiplier      float64
	roundMultiplier bool
	maxLevel        int
	multiplierList  []float64
	multFreeSpins   map[int]uint8
	flag            int
}

// NewGridMultipliersAction instantiates a new action to manipulate grid multipliers.
// It multiplies the multiplier of existing sticky symbols if a non-sticky triggerSymbol is found.
// The resulting multipliers are capped to the given max.
func NewGridMultipliersAction(symbol, triggerSymbol utils.Index, multiply, max uint16) *MultiplierAction {
	a := &MultiplierAction{flag: -1}
	a.init(TestGrid, Multipliers, reflect.TypeOf(a).String())
	a.symbol = symbol
	a.gridMultipliers = true
	a.triggerSymbol = triggerSymbol
	a.multiply = multiply
	a.max = max
	return a.finalizer()
}

// NewFirstMultiplierAction instantiates a new action to initialize the increasing scale of round multipliers.
// The scale mark and initial multiplier are set as indicated.
func NewFirstMultiplierAction(mark int, multiplier float64) *MultiplierAction {
	a := &MultiplierAction{flag: -1}
	a.init(PreBonus, Processed, reflect.TypeOf(a).String())
	a.firstMultiplier = true
	a.firstLevel = mark
	a.multiplier = multiplier
	return a.finalizer()
}

// NewMultiplierScaleAction instantiates a new action to manipulate an increasing scale of round multipliers.
// The scale mark is set to firstMark if it is zero, otherwise it is increased by the number of trigger symbols.
// The round multiplier is set from the list of multipliers according to the new scale mark.
// Mark 1 corresponds with the first element (index 0) in the multipliers list.
func NewMultiplierScaleAction(triggerSymbol utils.Index, firstLevel int, list ...float64) *MultiplierAction {
	a := &MultiplierAction{flag: -1}
	a.init(TestGrid, Multiplier, reflect.TypeOf(a).String())
	a.roundMultiplier = true
	a.triggerSymbol = triggerSymbol
	a.firstLevel = firstLevel
	a.multiplierList = list
	a.maxLevel = len(a.multiplierList)
	return a.finalizer()
}

func (a *MultiplierAction) WithFreeSpins(spins map[int]uint8) *MultiplierAction {
	if !a.roundMultiplier {
		panic("MultiplierAction: invalid configuration")
	}
	a.multFreeSpins = spins
	return a.finalizer()
}

func (a *MultiplierAction) WithFlag(flag int) *MultiplierAction {
	a.flag = flag
	return a
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *MultiplierAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.gridMultipliers:
		return a.gridTrigger(spin)

	case a.firstMultiplier:
		spin.progressLevel = a.firstLevel
		spin.multiplier = a.multiplier

		if a.flag >= 0 {
			spin.roundFlags[a.flag] = spin.progressLevel
		}

		return &a.SpinAction

	case a.roundMultiplier:
		return a.roundTrigger(spin)

	}
	return nil
}

func (a *MultiplierAction) BonusFreeSpins(multiplier float64) uint8 {
	return a.multFreeSpins[int(math.Round(multiplier))]
}

func (a *MultiplierAction) gridTrigger(spin *Spin) SpinActioner {
	if len(spin.multipliers) == 0 {
		return nil
	}

	var found bool
	for ix := range spin.indexes {
		if spin.indexes[ix] == a.triggerSymbol && !spin.sticky[ix] {
			found = true
			break
		}
	}
	if !found {
		return nil
	}

	for ix := range spin.indexes {
		if spin.indexes[ix] == a.symbol && spin.sticky[ix] && spin.multipliers[ix] > 0 {
			m := spin.multipliers[ix] * a.multiply
			if m > a.max {
				m = a.max
			}
			spin.multipliers[ix] = m
		}
	}

	return a
}

func (a *MultiplierAction) roundTrigger(spin *Spin) SpinActioner {
	count := spin.CountSymbol(a.triggerSymbol)
	if count == 0 {
		return nil
	}

	if spin.progressLevel == 0 {
		// first mark.
		spin.progressLevel = a.firstLevel
	} else {
		// repeat mark.
		spin.progressLevel += int(count)
		if spin.progressLevel > a.maxLevel {
			spin.progressLevel = a.maxLevel
		}
	}

	spin.multiplier = a.multiplierList[spin.progressLevel-1]

	if a.flag >= 0 {
		spin.roundFlags[a.flag] = spin.progressLevel
	}

	return a
}

func (a *MultiplierAction) finalizer() *MultiplierAction {
	b := bytes.Buffer{}
	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	switch {
	case a.gridMultipliers:
		b.WriteString(",gridMultipliers=1")
		b.WriteString(",symbol=")
		b.WriteString(strconv.Itoa(int(a.symbol)))
		b.WriteString(",triggerSymbol=")
		b.WriteString(strconv.Itoa(int(a.triggerSymbol)))
		b.WriteString(",multiply=")
		b.WriteString(strconv.Itoa(int(a.multiply)))
		b.WriteString(",max=")
		b.WriteString(strconv.Itoa(int(a.max)))
		if a.flag >= 0 {
			b.WriteString(",flag=")
			b.WriteString(strconv.Itoa(a.flag))
		}

	case a.firstMultiplier:
		b.WriteString(",firstMultiplier=1")
		b.WriteString(",firstMark=")
		b.WriteString(strconv.Itoa(a.firstLevel))
		b.WriteString(",multiplier=")
		b.WriteString(strconv.FormatFloat(a.multiplier, 'g', -1, 64))
		if a.flag >= 0 {
			b.WriteString(",flag=")
			b.WriteString(strconv.Itoa(a.flag))
		}

	case a.roundMultiplier:
		b.WriteString(",roundMultiplier=1")
		b.WriteString(",triggerSymbol=")
		b.WriteString(strconv.Itoa(int(a.triggerSymbol)))
		b.WriteString(",firstMark=")
		b.WriteString(strconv.Itoa(a.firstLevel))
		b.WriteString(",multipliers=")
		j, _ := json.Marshal(a.multiplierList)
		b.WriteString(string(j))
		if len(a.multFreeSpins) > 0 {
			b.WriteString(",freeSpins=")
			j, _ = json.Marshal(a.multFreeSpins)
			b.WriteString(string(j))
		}
		if a.flag >= 0 {
			b.WriteString(",flag=")
			b.WriteString(strconv.Itoa(a.flag))
		}
	}

	a.config = b.String()
	return a
}
