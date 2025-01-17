package rng_magic

import (
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

type Function interface {
	AND() bool
	OR() bool
	Append(f Function)
	Init(params FunctionInitParams)
	Matches(res rslt.Results) bool
}

type FunctionInitParams struct {
	g          *game.Regular
	newMatcher func(key string, params map[string]any, g *game.Regular) magic.Matcher
	conditions map[string]*magic.Condition
}

type and struct {
	closed bool
	list   []Function
}

func (a *and) AND() bool         { return !a.closed }
func (a *and) OR() bool          { return false }
func (a *and) Append(f Function) { a.list = append(a.list, f) }

func (a *and) Init(params FunctionInitParams) {
	for ix := range a.list {
		a.list[ix].Init(params)
	}
}

func (a *and) Matches(res rslt.Results) bool {
	for ix := range a.list {
		if !a.list[ix].Matches(res) {
			return false
		}
	}
	return true
}

type or struct {
	closed bool
	list   []Function
}

func (o *or) AND() bool         { return false }
func (o *or) OR() bool          { return !o.closed }
func (o *or) Append(f Function) { o.list = append(o.list, f) }

func (o *or) Init(params FunctionInitParams) {
	for ix := range o.list {
		o.list[ix].Init(params)
	}
}

func (o *or) Matches(res rslt.Results) bool {
	for ix := range o.list {
		if o.list[ix].Matches(res) {
			return true
		}
	}
	return false
}

type function struct {
	not     bool
	name    string
	params  []any
	matcher magic.Matcher
}

func (f *function) AND() bool         { return false }
func (f *function) OR() bool          { return false }
func (f *function) Append(_ Function) { panic("cannot append to function") }

func (f *function) Init(params FunctionInitParams) {
	if c := params.conditions[f.name]; c != nil {
		params2 := make(map[string]any)
		for ix, p := range c.Parameters {
			if ix < len(f.params) {
				params2[p.Name] = f.params[ix]
			} else {
				params2[p.Name] = nil
			}
		}
		f.matcher = params.newMatcher(f.name, params2, params.g)
	}
}

func (f *function) Matches(res rslt.Results) bool {
	return f.matcher(res)
}

type RngMagic struct {
	Magic1      []int   `json:"magic1"`
	Magic2      []int   `json:"magic2,omitempty"`
	NrOfResults int     `json:"nrOfResults"`
	TotalPayout float64 `json:"totalPayout"`
}

type RngMagics struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Magics  []*RngMagic `json:"magics"`
}
