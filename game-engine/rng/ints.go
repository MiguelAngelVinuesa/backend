package rng

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
)

type intsBuffer struct {
	reloads int
	prng    interfaces.Generator
	buffer  []int
	output  []int
}

func newIntsBuffer(prng interfaces.Generator) *intsBuffer {
	b := intsBufferPool.Get().(*intsBuffer)
	b.prng = prng
	return b
}

func (b *intsBuffer) Release() {
	if b != nil {
		b.reloads = 0
		b.prng = nil
		b.output = nil
		intsBufferPool.Put(b)
	}
}

func (b *intsBuffer) getN(n int) int {
	if len(b.output) == 0 {
		b.prng.IntsN(n, b.buffer)
		b.output = b.buffer
		b.reloads++
	}
	out := b.output[0]
	b.output = b.output[1:]
	return out
}

var intsBufferPool = sync.Pool{New: func() any { return &intsBuffer{buffer: make([]int, defaultCap)} }}

const defaultCap = 32
