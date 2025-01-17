package slots

type EndpointKind uint8

const (
	EndpointRound EndpointKind = iota + 1
	EndpointRoundResume
)

// EndpointFilterer represents the function signature to filter actions based on the called endpoint.
type EndpointFilterer func(endpoint EndpointKind) bool

// OnRound filters based on the round endpoint.
func OnRound(endpoint EndpointKind) bool {
	return endpoint == EndpointRound
}

// OnRoundResume filters based on the round/resume endpoint.
func OnRoundResume(endpoint EndpointKind) bool {
	return endpoint == EndpointRoundResume
}
