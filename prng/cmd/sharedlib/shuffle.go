package main

// #include <stdint.h>
import "C"

import (
	"runtime/cgo"
	"unsafe"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/cards"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

// FisherYatesShuffle shuffles the slice into a random order.
// See https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
//export FisherYatesShuffle
func FisherYatesShuffle(a *int, aLen int, h C.uintptr_t) {
	r := cgo.Handle(h).Value().(*rng.RNG)
	c := (*[1 << 28]int)(unsafe.Pointer(a))[:aLen:aLen]
	cards.FisherYatesShuffle(c, r)
}
