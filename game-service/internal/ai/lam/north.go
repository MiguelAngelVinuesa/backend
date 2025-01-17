package lam

type NorthPlayer struct{}

func (p *NorthPlayer) Choices() map[string]string {
	return map[string]string{"wing": "north"}
}
