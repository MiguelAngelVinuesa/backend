package main

import "C"

var hash string

// GetHash is the exported function to return the Git repo hash of the PRNG module used to compile the shared lib.
//export GetHash
func GetHash() *C.char {
	return C.CString(hash)
}

// main is only needed to satisfy the build requirements for a shared lib.
func main() {}
