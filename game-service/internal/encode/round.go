package encode

import (
	"math"

	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	mngr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
)

func BuildRoundResponse(game *game.Regular, round *mngr.Round, i18n *models.PrefetchI18n) *zjson.Encoder {
	requireChoice := game.AllowPlayerChoices() && game.NeedPlayerChoice()
	startBalance := round.PlayerBalance() - round.TotalWin()
	before, after := round.Balances(0)
	win := after - before
	bet := round.Bet()
	totalWin := round.ProgressiveWin(0)

	result := round.RoundResults()[0]
	first := round.Results()[0]
	first.Animations = game.BuildAnimations(first)
	isSpin := first.DataKind == rslt.SpinData
	count := int64(len(round.RoundResults()))
	spinWin := result.SpinWin

	if i18n != nil && first.Total > 0 {
		addLocalizedMessages(bet, first, i18n)
	}

	winFactor := math.Round(float64(win)*100/float64(bet)) / 100
	totalFactor := math.Round(float64(totalWin)*100/float64(bet)) / 100
	spinFactor := math.Round(float64(spinWin)*100/float64(bet)) / 100

	enc := zjson.AcquireEncoder(2048)
	enc.StartObject()

	enc.StartObjectField("roundData")
	enc.StringField("roundId", round.RoundID())
	enc.Int64Field("roundSeq", round.GameState().RoundSeq())
	enc.Uint8Field("dataKind", uint8(first.DataKind))
	enc.Int64Field("nrOfResults", count)
	enc.Int64Field("balanceBefore", startBalance+before)
	enc.Int64Field("balanceAfter", startBalance+after)
	enc.Int64Field("bet", bet)
	enc.Int64Field("win", after-before)
	enc.FloatField("winFactor", winFactor, 'f', 2)
	enc.IntBoolField("realWin", after-before > bet)
	enc.Int64FieldOpt("spinWin", spinWin)
	enc.FloatFieldOpt("spinFactor", spinFactor, 'f', 2)
	enc.IntBoolField("realSpinWin", spinWin > bet)
	enc.Int64Field("totalWin", totalWin)
	enc.FloatField("totalFactor", totalFactor, 'f', 2)
	enc.IntBoolField("realTotalWin", totalWin > bet)
	enc.IntBoolFieldOpt("maxPayout", result.MaxPayout > 0)
	enc.IntBoolFieldOpt("requireChoice", requireChoice)
	enc.EndObject()

	enc.StartObjectField("spinData")
	if isSpin {
		first.Encode2(enc) // do not encode the PRNG log!
	}
	enc.EndObject()

	if !isSpin {
		enc.StartObjectField("data")
		first.Encode2(enc) // do not encode the PRNG log!
		enc.EndObject()
	}

	enc.BoolField("success", true)

	enc.EndObject()
	return enc
}
