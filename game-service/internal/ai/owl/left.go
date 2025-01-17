package owl

type LeftPlayer struct{}

func (p *LeftPlayer) Choices() map[string]string {
	return map[string]string{"selection": "left"}
}
