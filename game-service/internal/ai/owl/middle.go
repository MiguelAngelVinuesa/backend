package owl

type MiddlePlayer struct{}

func (p *MiddlePlayer) Choices() map[string]string {
	return map[string]string{"selection": "middle"}
}
