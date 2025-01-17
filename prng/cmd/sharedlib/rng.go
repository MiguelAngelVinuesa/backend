package main

// #include <stdint.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

// NewRNG creates a new RNG and returns a handle to it to the caller.
//export NewRNG
func NewRNG() C.uintptr_t {
	return C.uintptr_t(cgo.NewHandle(rng.NewRNG()))
}

// GetUInt32 uses the RNG to return a random uint32.
//export GetUInt32
func GetUInt32(h C.uintptr_t) uint32 {
	r := cgo.Handle(h).Value().(*rng.RNG)
	return r.Uint32()
}

// GetUInt64 uses the RNG to return a random uint64.
//export GetUInt64
func GetUInt64(h C.uintptr_t) uint64 {
	r := cgo.Handle(h).Value().(*rng.RNG)
	return r.Uint64()
}

// GetIntN uses the RNG to return a random integer from the half open interval [0, n).
//export GetIntN
func GetIntN(h C.uintptr_t, n int) int {
	r := cgo.Handle(h).Value().(*rng.RNG)
	return r.IntN(n)
}

// GetIntsN uses the RNG to fill the out slice with random integers from the half open interval [0, n).
//export GetIntsN
func GetIntsN(h C.uintptr_t, n int, a *int, aLen int) {
	r := cgo.Handle(h).Value().(*rng.RNG)
	out := (*[1 << 28]int)(unsafe.Pointer(a))[:aLen:aLen]
	r.IntsN(n, out)
}

// FreeRNG signals that the caller is finished with the RNG, so we can return it to the memory pool here.
//export FreeRNG
func FreeRNG(h C.uintptr_t) {
	r := cgo.Handle(h).Value().(*rng.RNG)
	r.ReturnToPool()
	cgo.Handle(h).Delete()
}
