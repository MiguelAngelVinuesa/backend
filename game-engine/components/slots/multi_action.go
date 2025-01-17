package slots

import (
	"bytes"
	"reflect"
	"sort"
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// MultiActionFlagValue is an action that chooses one of a series of actions based on the value of the given spin flag.
// It only works for the "generic" SpinActioner actions which at least implement the Triggered, TriggeredWithState and/or Payout interface.
type MultiActionFlagValue struct {
	SpinAction
	flag    int
	actions map[int]SpinActioner
}

// NewMultiActionFlagValue instantiates a multi-select flag value action.
func NewMultiActionFlagValue(flag int, params ...any) *MultiActionFlagValue {
	max := len(params)
	if max%2 != 0 {
		panic("invalid MultiActionFlagValue configuration")
	}

	a := &MultiActionFlagValue{flag: flag, actions: make(map[int]SpinActioner, max/2)}

	for ix := 0; ix < max; ix += 2 {
		if action, ok := params[ix+1].(SpinActioner); ok {
			a.actions[conv.IntFromAny(params[ix])] = action
			if ix == 0 {
				a.init(action.Stage(), action.Result(), reflect.TypeOf(a).String())
			} else {
				if action.Stage() != a.stage || action.Result() != a.result {
					panic("invalid MultiActionFlagValue configuration")
				}
			}
		}
	}

	return a.finalize()
}

// GetAction returns the action based on the spin flag.
func (a *MultiActionFlagValue) GetAction(spin *Spin) SpinActioner {
	v := spin.roundFlags[a.flag]
	return a.actions[v]
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *MultiActionFlagValue) Triggered(spin *Spin) SpinActioner {
	v := spin.roundFlags[a.flag]
	if action := a.actions[v]; action != nil && action.CanTrigger(spin) {
		return action.Triggered(spin)
	}
	return nil
}

// TriggeredWithState implements the SpinActioner.TriggeredWithState() interface.
func (a *MultiActionFlagValue) TriggeredWithState(spin *Spin, state *SymbolsState) SpinActioner {
	v := spin.roundFlags[a.flag]
	if action := a.actions[v]; action != nil && action.CanTrigger(spin) {
		return action.TriggeredWithState(spin, state)
	}
	return nil
}

// Payout implements the SpinActioner.Payout() interface.
func (a *MultiActionFlagValue) Payout(spin *Spin, res *results.Result) SpinActioner {
	v := spin.roundFlags[a.flag]
	if action := a.actions[v]; action != nil && action.CanTrigger(spin) {
		return action.Payout(spin, res)
	}
	return nil
}

func (a *MultiActionFlagValue) finalize() *MultiActionFlagValue {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	b.WriteString(",flag=")
	b.WriteString(strconv.Itoa(a.flag))

	keys := make([]int, 0, len(a.actions))
	for k := range a.actions {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	b.WriteString(",actions={")
	for ix, k := range keys {
		v := a.actions[k]
		if ix > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(k))
		b.WriteString(":{")
		b.WriteString(v.Config())
		b.WriteByte('}')
	}
	b.WriteByte('}')

	a.config = b.String()
	return a
}

// MultiActionGridCount is an action that chooses one of a series of actions based on the count of the given symbol in the spin.
// It only works for the "generic" SpinActioner actions which at least implement the Triggered, TriggeredWithState and/or Payout interface.
type MultiActionGridCount struct {
	SpinAction
	symbol  utils.Index
	reels   []int
	actions map[int]SpinActioner
}

// NewMultiActionGridCount instantiates a multi-select flag value action.
func NewMultiActionGridCount(symbol utils.Index, reels []int, params ...any) *MultiActionGridCount {
	max := len(params)
	if max%2 != 0 {
		panic("invalid MultiActionGridCount configuration")
	}

	a := &MultiActionGridCount{symbol: symbol, reels: reels, actions: make(map[int]SpinActioner, max/2)}

	for ix := 0; ix < max; ix += 2 {
		if action, ok := params[ix+1].(SpinActioner); ok {
			a.actions[conv.IntFromAny(params[ix])] = action
			if ix == 0 {
				a.init(action.Stage(), action.Result(), reflect.TypeOf(a).String())
			} else {
				if action.Stage() != a.stage || action.Result() != a.result {
					panic("invalid MultiActionGridCount configuration")
				}
			}
		}
	}

	return a.finalize()
}

// GetAction returns the action based on the spin flag.
func (a *MultiActionGridCount) GetAction(spin *Spin) SpinActioner {
	v := a.count(spin)
	return a.actions[v]
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *MultiActionGridCount) Triggered(spin *Spin) SpinActioner {
	v := a.count(spin)
	if action := a.actions[v]; action != nil && action.CanTrigger(spin) {
		return action.Triggered(spin)
	}
	return nil
}

// TriggeredWithState implements the SpinActioner.TriggeredWithState() interface.
func (a *MultiActionGridCount) TriggeredWithState(spin *Spin, state *SymbolsState) SpinActioner {
	v := a.count(spin)
	if action := a.actions[v]; action != nil && action.CanTrigger(spin) {
		return action.TriggeredWithState(spin, state)
	}
	return nil
}

// Payout implements the SpinActioner.Payout() interface.
func (a *MultiActionGridCount) Payout(spin *Spin, res *results.Result) SpinActioner {
	v := a.count(spin)
	if action := a.actions[v]; action != nil && action.CanTrigger(spin) {
		return action.Payout(spin, res)
	}
	return nil
}

func (a *MultiActionGridCount) count(spin *Spin) int {
	var c int
	if len(a.reels) == 0 {
		for offset := range spin.indexes {
			if spin.indexes[offset] == a.symbol {
				c++
			}
		}
	} else {
		for ix := range a.reels {
			reel := a.reels[ix] - 1
			min := reel * spin.rowCount
			max := min + int(spin.mask[reel])
			for offset := min; offset < max; offset++ {
				if spin.indexes[offset] == a.symbol {
					c++
				}
			}
		}
	}
	return c
}

func (a *MultiActionGridCount) finalize() *MultiActionGridCount {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	b.WriteString(",symbol=")
	b.WriteString(strconv.Itoa(int(a.symbol)))

	if len(a.reels) > 0 {
		j, _ := json.Marshal(a.reels)
		b.WriteString(",reels=")
		b.Write(j)
	}

	keys := make([]int, 0, len(a.actions))
	for k := range a.actions {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	b.WriteString(",actions={")
	for ix, k := range keys {
		v := a.actions[k]
		if ix > 0 {
			b.WriteByte(',')
		}
		b.WriteString(strconv.Itoa(k))
		b.WriteString(":{")
		b.WriteString(v.Config())
		b.WriteByte('}')
	}
	b.WriteByte('}')

	a.config = b.String()
	return a
}

// MultiActionWeighted represents a weighted action selector.
type MultiActionWeighted struct {
	SpinAction
	weights utils.WeightedGenerator
	actions []SpinActioner
}

// NewMultiActionWeighted instantiates a weighted action selector.
func NewMultiActionWeighted(weights utils.WeightedGenerator, actions ...SpinActioner) *MultiActionWeighted {
	a := &MultiActionWeighted{
		weights: weights,
		actions: make([]SpinActioner, 0, len(actions)),
	}

	for ix := 0; ix < len(actions); ix++ {
		action := actions[ix]
		a.actions = append(a.actions, action)
		if ix == 0 {
			a.init(action.Stage(), action.Result(), reflect.TypeOf(a).String())
		} else {
			if action.Stage() != a.stage || action.Result() != a.result {
				panic("invalid NewMultiActionWeighted configuration")
			}
		}
	}

	return a.finalize()
}

// GetAction returns a random action based on the weights.
func (a *MultiActionWeighted) GetAction(spin *Spin) SpinActioner {
	ix := int(a.weights.RandomIndex(spin.prng)) - 1
	if ix < len(a.actions) {
		return a.actions[ix]
	}
	return nil
}

// Triggered implements the SpinActioner.Triggered() interface.
func (a *MultiActionWeighted) Triggered(spin *Spin) SpinActioner {
	action := a.GetAction(spin)
	if action != nil && action.CanTrigger(spin) {
		return action.Triggered(spin)
	}
	return nil
}

// TriggeredWithState implements the SpinActioner.TriggeredWithState() interface.
func (a *MultiActionWeighted) TriggeredWithState(spin *Spin, state *SymbolsState) SpinActioner {
	action := a.GetAction(spin)
	if action != nil && action.CanTrigger(spin) {
		return action.TriggeredWithState(spin, state)
	}
	return nil
}

// Payout implements the SpinActioner.Payout() interface.
func (a *MultiActionWeighted) Payout(spin *Spin, res *results.Result) SpinActioner {
	action := a.GetAction(spin)
	if action != nil && action.CanTrigger(spin) {
		return action.Payout(spin, res)
	}
	return nil
}

func (a *MultiActionWeighted) finalize() *MultiActionWeighted {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())
	b.WriteString(",weights=")
	b.WriteString(a.weights.String())

	b.WriteString(",actions=[")
	for ix := range a.actions {
		v := a.actions[ix]
		if ix > 0 {
			b.WriteByte(',')
		}
		b.WriteByte('{')
		b.WriteString(v.Config())
		b.WriteString("}]")
	}
	b.WriteByte('}')

	a.config = b.String()
	return a
}
