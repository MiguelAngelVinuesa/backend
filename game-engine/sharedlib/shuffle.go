package sharedlib

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lprng -Wl,-rpath=/usr/local/lib
// #include "libprng.h"
import "C"

import (
	"sync"
	"unsafe"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
)

// Shuffler implements the Shuffler interface using the shared library "libprng.so".
type Shuffler struct {
	h C.uintptr_t
}

// AcquireShuffler instantiates a new shuffler from the memory pool.
func AcquireShuffler() interfaces.Shuffler {
	h := C.NewRNG()
	if h == 0 {
		panic("failed to get new RNG handle")
	}

	s := shufflerPool.Get().(*Shuffler)
	s.h = h
	return s
}

// Release implements the Objecter interface.
func (s *Shuffler) Release() {
	if s != nil {
		C.FreeRNG(s.h)
		s.h = 0
		shufflerPool.Put(s)
	}
}

// Shuffle shuffles the slice into a random order using the FisherYates algorithm.
// See https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
func (s *Shuffler) Shuffle(c []int) {
	C.FisherYatesShuffle((*C.longlong)(unsafe.Pointer(&c[0])), C.GoInt(len(c)), s.h)
}

// shufflerPool is the memory pool for Shufflers.
var shufflerPool = sync.Pool{New: func() any { return &Shuffler{} }}
