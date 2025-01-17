package main

// #cgo CFLAGS: -I/usr/local/include
// #cgo LDFLAGS: -L/usr/local/lib -lprng -Wl,-rpath=/usr/local/lib
// #include "libprng.h"
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	q := C.GoString(C.GetHash())
	fmt.Printf("%s\n\n", q)

	h := C.NewRNG()
	if h == 0 {
		panic("failed to get new PRNG handle")
	}
	defer C.FreeRNG(h)

	for ix := 0; ix < 10; ix++ {
		n := C.GetUInt32(h)
		fmt.Printf("%d  ", uint32(n))
	}
	fmt.Println()

	for ix := 0; ix < 10; ix++ {
		n := C.GetUInt64(h)
		fmt.Printf("%d  ", uint64(n))
	}
	fmt.Println()

	out := make([]int, 15)
	C.GetIntsN(h, 12, (*C.longlong)(unsafe.Pointer(&out[0])), C.GoInt(len(out)))
	for ix := range out {
		fmt.Printf("%d  ", out[ix])
	}
	fmt.Println()

	c := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52}
	C.FisherYatesShuffle((*C.longlong)(unsafe.Pointer(&c[0])), C.GoInt(len(c)), h)
	for ix := range c {
		fmt.Printf("%d ", c[ix])
	}
	fmt.Println()
}
