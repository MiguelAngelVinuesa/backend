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

// RNG implements the Generator interface (PRNG) using the shared library "libprng.so".
type RNG struct {
	h C.uintptr_t
}

// AcquireRNG instantiates a new RNG from the memory pool.
func AcquireRNG() interfaces.Generator {
	globalMU.Lock()
	h := C.NewRNG()
	if h == 0 {
		panic("failed to get new RNG handle")
	}

	r := rngPool.Get().(*RNG)
	r.h = h
	globalMU.Unlock()
	return r
}

// ReturnToPool returns the RNG to the memory pool.
func (r *RNG) ReturnToPool() {
	if r != nil {
		globalMU.Lock()
		C.FreeRNG(r.h)
		r.h = 0
		rngPool.Put(r)
		globalMU.Unlock()
	}
}

// Uint32 return a random uint32.
func (r *RNG) Uint32() uint32 {
	globalMU.Lock()
	i := uint32(C.GetUInt32(r.h))
	globalMU.Unlock()
	return i
}

// Uint64 return a random uint64.
func (r *RNG) Uint64() uint64 {
	globalMU.Lock()
	i := uint64(C.GetUInt64(r.h))
	globalMU.Unlock()
	return i
}

// IntN return a random integer in the half open interval [0,n).
func (r *RNG) IntN(n int) int {
	globalMU.Lock()
	i := int(C.GetIntN(r.h, C.GoInt(n)))
	globalMU.Unlock()
	return i
}

// IntsN fills the given out slice with random integers in the half open interval [0,n).
func (r *RNG) IntsN(n int, out []int) {
	globalMU.Lock()
	C.GetIntsN(r.h, C.GoInt(n), (*C.longlong)(unsafe.Pointer(&out[0])), C.GoInt(len(out)))
	globalMU.Unlock()
}

// rngPool is the memory pool for RNGs.
var rngPool = sync.Pool{New: func() interface{} { return &RNG{} }}

var globalMU sync.Mutex
