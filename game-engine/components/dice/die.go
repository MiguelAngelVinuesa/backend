package dice

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
)

// Die contains the details of a die.
type Die struct {
	faces  int
	first  int
	values []int
}

// NewStandardDie instantiates a new 6-faced die from the memory pool.
func NewStandardDie() *Die {
	d := diePool.Get().(*Die)
	d.faces, d.first = 6, 1
	return d
}

// NewFacedDie instantiates a new n-faced die from the memory pool.
func NewFacedDie(faces, first int) *Die {
	d := diePool.Get().(*Die)
	d.faces, d.first = faces, first
	return d
}

// NewValuesDie instantiates a new die with specified values from the memory pool.
func NewValuesDie(values ...int) *Die {
	d := diePool.Get().(*Die)
	d.faces, d.values = len(values), values
	return d
}

// Release returns the die to the memory pool.
func (d *Die) Release() {
	if d != nil {
		d.values = nil
		diePool.Put(d)
	}
}

// Roll rolls the die using the given PRNG and returns the result.
func (d *Die) Roll(rng interfaces.Generator) int {
	i := rng.IntN(d.faces)
	if len(d.values) == d.faces {
		return d.values[i]
	}
	return d.first + i
}

var diePool = sync.Pool{New: func() any { return &Die{} }}
