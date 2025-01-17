package slots

import (
	"math"
)

// ChanceModifier is the interface for modifying chances based on a math function.
// a = value (chance) to modify.
// x = function to retrieve power value.
type ChanceModifier interface {
	Exec(a float64, spin *Spin) float64
}

// NewPowerFunc instantiates a chance modifier based on the "a*b^x+c" function.
// a    = chance value - supplied during Exec().
// x    = function to retrieve power value.
// b, c = function constants.
func NewPowerFunc(b, c float64, x func(spin *Spin) float64) ChanceModifier {
	return &chanceModifier{f: PowerFunc, b: b, c: c, x: x}
}

// NewDivideFunc instantiates a chance modifier based on the "a*b/(x+c)" function.
// a    = chance value - supplied during Exec().
// x    = function to retrieve power value.
// b, c = function constants.
func NewDivideFunc(b, c float64, x func(spin *Spin) float64) ChanceModifier {
	return &chanceModifier{f: DivideFunc, b: b, c: c, x: x}
}

// NewMultiFunc instantiates a chance modifier based on multiple functions.
func NewMultiFunc(f ...ChanceModifier) ChanceModifier {
	return &multiModifier{f: f}
}

// FunctionKind represents the kind of modification function.
type FunctionKind uint8

const (
	PowerFunc FunctionKind = iota + 1 // a*b^x+c
	DivideFunc
)

type chanceModifier struct {
	f FunctionKind
	b float64
	c float64
	x func(spin *Spin) float64
}

// Exec implements the ChanceModifier interface.
func (c *chanceModifier) Exec(a float64, spin *Spin) float64 {
	switch c.f {
	case PowerFunc:
		return a*math.Pow(c.b, c.x(spin)) + c.c
	case DivideFunc:
		return a * c.b / (c.x(spin) + c.c)
	default:
		return a
	}
}

type multiModifier struct {
	f []ChanceModifier
}

// Exec implements the ChanceModifier interface.
func (m *multiModifier) Exec(a float64, spin *Spin) float64 {
	out := a
	for ix := range m.f {
		out = m.f[ix].Exec(out, spin)
	}
	return out
}
