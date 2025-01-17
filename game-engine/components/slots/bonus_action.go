package slots

import (
	"bytes"
	"math"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// BonusAction represents a bonus award action.
type BonusAction struct {
	SpinAction
	instant       bool
	chance        float64
	tease         bool
	choice        string
	options       []string
	selector      bool
	selWeights    utils.WeightedGenerator
	selCount      int
	selChoiceFlag int
	selFlag       int
	wheel         bool
	wheelFlag     int
	wheelWeights  utils.WeightedGenerator
}

// NewInstantBonusAction instantiates an instant bonus award action.
func NewInstantBonusAction(chance float64) *BonusAction {
	a := newBonusAction()
	a.result = InstantBonus
	a.instant = true
	a.chance = chance
	return a.finalizer()
}

// NewBonusSelectorAction instantiates a bonus selection action.
func NewBonusSelectorAction(weights utils.WeightedGenerator, count, choiceFlag, flag int) *BonusAction {
	a := newBonusAction()
	a.result = SpecialResult
	a.selector = true
	a.selWeights = weights
	a.selCount = count
	a.selChoiceFlag = choiceFlag
	a.selFlag = flag
	return a.finalizer()
}

// NewInstantBonusWheelAction instantiates an instant bonus wheel action.
func NewInstantBonusWheelAction(flag int, weights utils.WeightedGenerator) *BonusAction {
	a := newBonusAction()
	a.result = BonusGame
	a.wheel = true
	a.wheelFlag = flag
	a.wheelWeights = weights
	return a.finalizer()
}

// WithTease can be used to indicate that the instant bonus is a teaser and won't actually happen.
func (a *BonusAction) WithTease() *BonusAction {
	a.tease = true
	return a.finalizer()
}

// WithPlayerChoice can be used to indicate that a player choice must be made before the game continues.
func (a *BonusAction) WithPlayerChoice(choice string, options ...string) *BonusAction {
	a.playerChoice = true
	a.choice = choice
	a.options = options
	return a.finalizer()
}

// WithAlternate can be used to add an alternative action for cases where they need to be mutually exclusive.
func (a *BonusAction) WithAlternate(alt *BonusAction) *BonusAction {
	a.alternate = alt
	return a.finalizer()
}

// Triggered implements the Actioner interface.
func (a *BonusAction) Triggered(spin *Spin) SpinActioner {
	switch {
	case a.instant:
		if spin.prng.IntN(10000) < int(math.Round(a.ModifyChance(a.chance, spin)*100)) {
			return a
		}
		if a.alternate != nil {
			return a.alternate.Triggered(spin)
		}

	case a.selector, a.wheel:
		return a // always triggers, as the action happens somewhere else.
	}

	return nil
}

// InstantBonus instantiates an instant bonus result determined by the kind of bonus.
func (a *BonusAction) InstantBonus(_ *Spin) interfaces.Objecter2 {
	switch {
	case a.tease:
		return results.AcquireInstantBonusTeaser()
	case a.playerChoice:
		return results.AcquireInstantBonusChoice(a.choice, a.options...)
	}
	return nil
}

// BonusSelect instantiates a bonus selection result using the provided weights and count.
// It will also update the requested flag with the player selection.
func (a *BonusAction) BonusSelect(spin *Spin) interfaces.Objecter2 {
	player := spin.roundFlags[a.selChoiceFlag]
	if player <= 0 || player > a.selCount {
		player = 1
	}

	out := make(utils.Indexes, 0, 16)[:a.selCount]
	a.selWeights.FillRandom(spin.prng, a.selCount, out)

	if player != 1 {
		out[player-1], out[0] = out[0], out[player-1]
	}

	chosen := out[player-1]

	if a.selFlag >= 0 {
		spin.roundFlags[a.selFlag] = int(chosen)
	}

	return results.AcquireBonusSelectorChoice(uint8(player), chosen, out...)
}

// BonusGame plays a bonus game and returns the result.
func (a *BonusAction) BonusGame(spin *Spin) interfaces.Objecter2 {
	switch {
	case a.wheel:
		w := wheel.AcquireBonusWheel(spin.prng, a.wheelWeights)
		value, result := w.Run(nil)
		w.Release()

		if a.wheelFlag >= 0 {
			spin.roundFlags[a.wheelFlag] = value
		}
		return result

	default:
		return nil
	}
}

// FeatureTransition returns the applicable bonus feature transition kind.
func (a *BonusAction) FeatureTransition() FeatureTransitionKind {
	switch {
	case a.tease:
		return InstantBonusTeaser
	case a.instant:
		return InstantBonusRequested
	case a.selector:
		return InstantBonusResulted
	case a.wheel:
		return BonusWheelTransition
	default:
		return 0
	}
}

func newBonusAction() *BonusAction {
	a := &BonusAction{}
	a.init(PreSpin, Processed, reflect.TypeOf(a).String())
	return a
}

func (a *BonusAction) finalizer() *BonusAction {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	if a.instant {
		b.WriteString(",instant=true")
		b.WriteString(",chance=")
		b.WriteString(strconv.FormatFloat(a.chance, 'g', -1, 64))
	}

	if a.tease {
		b.WriteString(",tease=true")
	}

	if a.playerChoice {
		b.WriteString(",playerChoice=true")
		b.WriteString(",choice=")
		b.WriteString(a.choice)
		b.WriteString(",options=")
		j, _ := json.Marshal(a.options)
		b.WriteString(string(j))
	}

	if a.selector {
		b.WriteString(",selector=true")
		b.WriteString(",weights=")
		b.WriteString(a.selWeights.String())
		b.WriteString(",count=")
		b.WriteString(strconv.Itoa(a.selCount))
		b.WriteString(",choiceFlag=")
		b.WriteString(strconv.Itoa(a.selChoiceFlag))
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.selFlag))
	}

	if a.wheel {
		b.WriteString(",wheel=true")
		b.WriteString(",flag=")
		b.WriteString(strconv.Itoa(a.wheelFlag))
		b.WriteString(",weights=")
		b.WriteString(a.wheelWeights.String())
	}

	if a.alternate != nil {
		b.WriteString(",alternate=true")
	}

	a.config = b.String()
	return a
}
