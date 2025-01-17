package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// ScriptedRoundSelector represents the selector to determine if and which scripted round is selected for the next round.
type ScriptedRoundSelector struct {
	chance    float64          // pct chance a scripted round will be activated (max 4 decimals).
	flag      int              // indicates the spin flag to update when a script is selected.
	ids       utils.Indexes    // list of scripted round identifiers.
	weights   []float64        // list of scripted round weights.
	scripts   []*ScriptedRound // map of scripted rounds.
	bonusBuys utils.UInt8s     // list of bonus buys script selection is allowed for.
	bonusFlag int              // index of flag indicating bonus buy.
	chancesBB []float64        // list of pct chances (one for each BB).
}

// NewScriptedRoundSelector instantiates a new scripted round selector.
func NewScriptedRoundSelector(chance float64, rounds ...*ScriptedRound) *ScriptedRoundSelector {
	l := len(rounds)
	s := &ScriptedRoundSelector{
		chance:    chance,
		flag:      -1,
		ids:       make(utils.Indexes, l),
		weights:   make([]float64, l),
		scripts:   make([]*ScriptedRound, l),
		bonusFlag: -1,
	}

	for ix := range rounds {
		r := rounds[ix]
		s.ids[ix] = utils.Index(r.id)
		s.weights[ix] = r.weight
		s.scripts[ix] = r
	}

	return s
}

// WithBonusBuys limits the script selector to specific bonus buys.
func (s *ScriptedRoundSelector) WithBonusBuys(bb ...uint8) *ScriptedRoundSelector {
	s.bonusBuys = bb
	return s
}

// WithBonusChances sets the chance of the selector for each specific bonus buy.
func (s *ScriptedRoundSelector) WithBonusChances(flag int, chances ...float64) *ScriptedRoundSelector {
	s.bonusFlag = flag
	s.chancesBB = chances
	return s
}

// WithSpinFlag sets the spin flag to indicate if a script was selected.
func (s *ScriptedRoundSelector) WithSpinFlag(flag int) *ScriptedRoundSelector {
	s.flag = flag
	return s
}

// BonusBuyAllowed returns true if the given bonus buy is 0 or in the list of allowed bonus buys.
func (s *ScriptedRoundSelector) BonusBuyAllowed(bb uint8) bool {
	if bb == 0 {
		return true
	}
	for ix := range s.bonusBuys {
		if bb == s.bonusBuys[ix] {
			return true
		}
	}
	return false
}

// SpinFlag returns the spin flag.
func (s *ScriptedRoundSelector) SpinFlag() int {
	return s.flag
}

// Triggered determines if a scripted round will be activated.
// Which round will be randomly selected is based on the weights.
func (s *ScriptedRoundSelector) Triggered(spin *Spin) bool {
	c := s.chance

	if s.bonusFlag >= 0 {
		if f := spin.roundFlags[s.bonusFlag] - 1; f >= 0 && f < len(s.chancesBB) {
			c = s.chancesBB[f]
		}
	}

	return float64(spin.prng.IntN(1000000)) < c*10000
}

// GetWeighting acquires a weighted index generator.
func (s *ScriptedRoundSelector) GetWeighting() utils.WeightedGenerator {
	return utils.AcquireWeighting().AddWeights(s.ids, s.weights)
}

// GetScripts returns the map of scripted rounds.
func (s *ScriptedRoundSelector) GetScripts() []*ScriptedRound {
	return s.scripts
}

// ScriptedRound represents the actions used to generate a specific scripted round.
type ScriptedRound struct {
	id        int          // script id.
	weight    float64      // script weight.
	actions1  SpinActions  // primary actions.
	actions2  SpinActions  // secondary actions.
	bonusBuys utils.UInt8s // list of bonus buys the script is allowed for.
}

// NewScriptedRound instantiates a new scripted round.
func NewScriptedRound(id int, weight float64, actions1, actions2 []SpinActioner) *ScriptedRound {
	return &ScriptedRound{
		id:       id,
		weight:   weight,
		actions1: actions1,
		actions2: actions2,
	}
}

// WithBonusBuys limits the scripted round to specific bonus buys.
func (s *ScriptedRound) WithBonusBuys(bb ...uint8) *ScriptedRound {
	s.bonusBuys = bb
	return s
}

// ID returns the unique identifier of the scripted round.
func (s *ScriptedRound) ID() int {
	return s.id
}

// Actions returns the list of actions for first and free spins.
func (s *ScriptedRound) Actions() (SpinActions, SpinActions) {
	return s.actions1, s.actions2
}

// BonusBuyAllowed returns true if the given bonus buy is 0 or in the list of allowed bonus buys.
func (s *ScriptedRound) BonusBuyAllowed(bb uint8) bool {
	if bb == 0 {
		return true
	}
	for ix := range s.bonusBuys {
		if bb == s.bonusBuys[ix] {
			return true
		}
	}
	return false
}
