package magic

type Condition struct {
	Name       string       `json:"name,omitempty"`
	Parameters []*Parameter `json:"parameters,omitempty"`
}

func NewCondition(name string, params ...*Parameter) *Condition {
	return &Condition{Name: name, Parameters: params}
}

func (c *Condition) AddParameter(param *Parameter) *Condition {
	c.Parameters = append(c.Parameters, param)
	return c
}

type Conditions []*Condition
