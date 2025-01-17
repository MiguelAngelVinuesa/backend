package rng_magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
)

func TestRngStatement(gameNR tg.GameNR, statement string) error {
	conditions, symbols := GameData(gameNR)
	_, err := parse(statement, conditions, symbols)
	return err
}
