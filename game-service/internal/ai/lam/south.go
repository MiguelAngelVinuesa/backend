package lam

type SouthPlayer struct{}

func (p *SouthPlayer) Choices() map[string]string {
	return map[string]string{"wing": "south"}
}
