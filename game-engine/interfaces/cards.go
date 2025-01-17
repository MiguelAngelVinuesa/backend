package interfaces

// Shuffler extends Objecter with a Shuffle function.
type Shuffler interface {
	Shuffle(c []int)
	Release()
}
