package slots

import (
	"bytes"
	"reflect"
	"strconv"

	"github.com/goccy/go-json"
)

type ChoiceAction struct {
	SpinAction
	flag    int
	key     string
	values  []string
	results []int
}

// NewPlayerChoiceAction instantiates a new choice action.
// This action is used for player choice features.
func NewPlayerChoiceAction(flag int, key string, values []string, results []int) *ChoiceAction {
	if flag < 0 || flag >= 16 || len(values) < 1 || len(results) != len(values) || key == "" {
		panic("ChoiceAction: invalid configuration")
	}

	a := &ChoiceAction{}
	a.init(TestPlayerChoice, Processed, reflect.TypeOf(a).String())
	a.flag = flag
	a.key = key
	a.values = values
	a.results = results
	return a.finalize()
}

// TestChoices tests the player choices for the configered key and updates the configured round flag accordingly.
func (a *ChoiceAction) TestChoices(spin *Spin, choices map[string]string) SpinActioner {
	for k := range choices {
		if k == a.key {
			s := choices[k]
			for ix := range a.values {
				if a.values[ix] == s {
					spin.roundFlags[a.flag] = a.results[ix]
					return a
				}
			}
			return nil
		}
	}
	return nil
}

func (a *ChoiceAction) finalize() *ChoiceAction {
	b := bytes.Buffer{}

	b.WriteString("stage=")
	b.WriteString(a.stage.String())
	b.WriteString(",result=")
	b.WriteString(a.result.String())
	b.WriteString(",flag=")
	b.WriteString(strconv.Itoa(a.flag))
	b.WriteString(",key=")
	b.WriteString(a.key)

	b.WriteString(",values=")
	j, _ := json.Marshal(a.values)
	b.Write(j)

	b.WriteString(",results=")
	j, _ = json.Marshal(a.results)
	b.Write(j)

	a.config = b.String()
	return a
}

// ChoiceActions is a convenience type for a slice of player choice actions.
type ChoiceActions []*ChoiceAction
