package magic

type Parameter struct {
	Name string  `json:"name,omitempty"`
	Kind string  `json:"kind,omitempty"`
	Min  float64 `json:"min,omitempty"`
	Max  float64 `json:"max,omitempty"`
}

func NewStringParam(name string) *Parameter {
	return &Parameter{Name: name, Kind: "STRING"}
}

func NewIntParam(name string, min, max int) *Parameter {
	return &Parameter{Name: name, Kind: "INT", Min: float64(min), Max: float64(max)}
}

func NewFloatParam(name string, min, max float64) *Parameter {
	return &Parameter{Name: name, Kind: "FLOAT", Min: min, Max: max}
}

func NewSymbolParam() *Parameter {
	return &Parameter{Name: fieldSymbol, Kind: "STRING"}
}

func NewSymbolTypeParam() *Parameter {
	return &Parameter{Name: fieldSymbolType, Kind: "STRING"}
}

func NewPaylineParam() *Parameter {
	return &Parameter{Name: fieldPayline, Kind: "INT", Min: 1, Max: 99}
}

func NewSequenceParam() *Parameter {
	return &Parameter{Name: fieldSeq, Kind: "INT", Min: 1, Max: 999}
}

func NewReelsParam() *Parameter {
	return &Parameter{Name: fieldReels, Kind: "REELS"}
}

func NewGridParam() *Parameter {
	return &Parameter{Name: fieldGrid, Kind: "GRID"}
}
