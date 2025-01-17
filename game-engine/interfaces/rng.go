package interfaces

// Generator is the interface for a PRNG.
type Generator interface {
	ReturnToPool()

	// Uint32 returns a random uint32.
	Uint32() uint32
	// Uint64 returns a random uint64.
	Uint64() uint64
	// IntN returns a random integer in the half open interval [0,n).
	IntN(n int) int
	// IntsN fills the given slice with random integers in the half open interval [0,n).
	IntsN(n int, out []int)
}
