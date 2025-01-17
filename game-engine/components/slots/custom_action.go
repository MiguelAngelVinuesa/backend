package slots

import (
	"bytes"
	"reflect"
)

// CustomAction can be used to build a very game-specific action.
type CustomAction struct {
	SpinAction
	trigger          func(*Spin) bool
	triggerWithState func(*Spin, *SymbolsState) bool
	alternate        SpinActioner
}

// NewCustomAction instantiates a new custom action.
func NewCustomAction(stage SpinActionStage, result SpinActionResult, trigger func(*Spin) bool, triggerWithState func(*Spin, *SymbolsState) bool) *CustomAction {
	a := &CustomAction{trigger: trigger, triggerWithState: triggerWithState}
	a.init(stage, result, reflect.TypeOf(a).String())
	return a.finalizer()
}

// WithAlternate sets an alternative action to be called when the action doesn't trigger itself.
func (a *CustomAction) WithAlternate(alt SpinActioner) *CustomAction {
	a.alternate = alt
	return a.finalizer()
}

// Triggered implements the SpinActioner.Triggered interface.
func (a *CustomAction) Triggered(spin *Spin) SpinActioner {
	if a.trigger != nil {
		if a.trigger(spin) {
			return a
		}
	}
	if a.alternate != nil {
		return a.alternate.Triggered(spin)
	}
	return nil
}

// TriggeredWithState implements the SpinActioner.TriggeredWithState interface.
func (a *CustomAction) TriggeredWithState(spin *Spin, state *SymbolsState) SpinActioner {
	if a.triggerWithState != nil {
		if a.triggerWithState(spin, state) {
			return a
		}
	}
	if a.alternate != nil {
		return a.alternate.TriggeredWithState(spin, state)
	}
	return nil
}

func (a *CustomAction) finalizer() *CustomAction {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())

	if a.trigger != nil {
		b.WriteString(",trigger=defined")
	}
	if a.triggerWithState != nil {
		b.WriteString(",triggerWithState=defined")
	}

	a.config = b.String()
	return a
}
