package rng

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// Buffer represents a buffered PRNG.
// It fully implements the Generator interface.
// For Uint32() and Uint64() it uses the supplied prng directly.
// For IntN and IntsN it uses internal buffers to reduce the go->C->go stack switches for the shared rng lib.
// It also logs all IntN / IntsN calls with the produced outputs.
type Buffer struct {
	withLog bool
	prng    interfaces.Generator
	buffers map[int]*intsBuffer
	inputs  []int
	outputs []int
	cache   []int
	get     func(n int) int
}

// AcquireBuffer instantiates a new buffered PRNG from the memory pool.
// It takes ownership of the supplied prng. *DO NOT* call ReturnToPool on the supplied prng!
// withLog indicates if the buffer should keep a log of all IntN/IntsN input/output values.
func AcquireBuffer(prng interfaces.Generator, withLog bool) *Buffer {
	b := bufferPool.Get().(*Buffer)
	b.withLog = withLog
	b.prng = prng

	if withLog {
		b.inputs = utils.PurgeInts(b.inputs, defaultLogSize)
		b.outputs = utils.PurgeInts(b.outputs, defaultLogSize)
		b.get = b.getLogged
	} else {
		b.get = b.getN
	}

	return b
}

// ReturnToPool returns the buffered prng to the memory pool.
func (b *Buffer) ReturnToPool() {
	b.Release()
}

// Release returns the buffered prng to the memory pool.
func (b *Buffer) Release() {
	if b != nil {
		b.prng.ReturnToPool()
		b.prng = nil
		b.cache = nil
		b.get = nil

		for n := range b.buffers {
			b.buffers[n].Release()
			b.buffers[n] = nil
		}

		if len(b.inputs) > maxLogSize {
			b.inputs = nil
		} else if b.inputs != nil {
			b.inputs = b.inputs[:0]
		}

		if len(b.outputs) > maxLogSize {
			b.outputs = nil
		} else if b.outputs != nil {
			b.outputs = b.outputs[:0]
		}

		bufferPool.Put(b)
	}
}

// WithCache loads the random number cache (debug mode), and sets the retrieval function to retrieve from the cache.
// It does nothing if the cache is not defined or has less than 2 elements.
func (b *Buffer) WithCache(cache []int) *Buffer {
	if len(cache) >= 2 {
		b.cache = cache
		b.get = b.getCached
	}
	return b
}

// Uint32 implements the Generator interface.
func (b *Buffer) Uint32() uint32 {
	return b.prng.Uint32()
}

// Uint64 implements the Generator interface.
func (b *Buffer) Uint64() uint64 {
	return b.prng.Uint64()
}

// IntN implements the Generator interface.
func (b *Buffer) IntN(n int) int {
	return b.get(n)
}

// IntsN implements the Generator interface.
func (b *Buffer) IntsN(n int, out []int) {
	for ix := range out {
		out[ix] = b.get(n)
	}
}

// LogSize returns the size of the current IntN/IntsN log.
func (b *Buffer) LogSize() int {
	return len(b.inputs)
}

// Log returns the log for all IntN/IntsN input/output values during the buffer's lifetime.
// It will only return valid output if logging was turned on and IntN/IntsN was called at least once.
func (b *Buffer) Log() ([]int, []int) {
	return b.inputs, b.outputs
}

// ResetLog resets the IntN/IntsN log.
func (b *Buffer) ResetLog() {
	if b.inputs != nil {
		b.inputs = b.inputs[:0]
	}
	if b.outputs != nil {
		b.outputs = b.outputs[:0]
	}
}

// getN retrieves a random number from the appropriate buffer.
func (b *Buffer) getN(n int) int {
	buffer := b.buffers[n]
	if buffer == nil {
		buffer = newIntsBuffer(b.prng)
		b.buffers[n] = buffer
	}
	return buffer.getN(n)
}

// getLogged retrieves a random number from the appropriate buffer, and logs the input/output value.
func (b *Buffer) getLogged(n int) int {
	out := b.getN(n)
	b.inputs = append(b.inputs, n)
	b.outputs = append(b.outputs, out)
	return out
}

// getCached retrieves a random number from the cache, if the input matches,
// otherwise, it uses the normal retrieval from the appropriate buffer.
// If the cache becomes empty it sets the get function to the appropriate non-cached retrieval function.
func (b *Buffer) getCached(n int) int {
	l := len(b.cache)

	var out int
	if l < 2 || b.cache[0] != n {
		out = -1
		b.cache = b.cache[0:]
		l = 0
	} else {
		out = b.cache[1]
		b.cache = b.cache[2:]
		l -= 2
	}

	if l < 2 {
		if b.withLog {
			b.get = b.getLogged
		} else {
			b.get = b.getN
		}
		if out < 0 {
			return b.get(n)
		}
	}

	return out
}

// bufferCount returns the number of internal buffers.
func (b *Buffer) bufferCount() int {
	return len(b.buffers)
}

// reloadCount return the total reload count of all internal buffers.
func (b *Buffer) reloadCount() int {
	var out int
	for n := range b.buffers {
		out += b.buffers[n].reloads
	}
	return out
}

var (
	bufferPool = sync.Pool{New: func() any { return &Buffer{buffers: make(map[int]*intsBuffer, defaultBuffersSize)} }}
)

const (
	defaultBuffersSize = 16
	defaultLogSize     = 256
	maxLogSize         = 1024
)
