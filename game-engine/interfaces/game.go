package interfaces

// Gamer is the interface for retrieving game round details.
type Gamer interface {
	ResultCount() int
	TotalPayout() float64
}
