package owl

type RightPlayer struct{}

func (p *RightPlayer) Choices() map[string]string {
	return map[string]string{"selection": "right"}
}
