package encode

import (
	"math"

	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	mngr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
)

func BuildRoundNextResponse(roundID string, endBalance int64, game *game.Regular, nrOfResults int, roundResult *mngr.RoundResult, i18n *models.PrefetchI18n) *zjson.Encoder {
	bet, win := roundResult.Bet, roundResult.Win
	startBalance := endBalance - roundResult.TotalWin
	bonusWin := roundResult.BonusWin
	totalWin := roundResult.ProgressiveWin
	spinWin := roundResult.SpinWin

	winFactor := math.Round(float64(win)*100/float64(bet)) / 100
	bonusFactor := math.Round(float64(bonusWin)*100/float64(bet)) / 100
	totalFactor := math.Round(float64(totalWin)*100/float64(bet)) / 100
	spinFactor := math.Round(float64(spinWin)*100/float64(bet)) / 100

	round := roundResult.Result()
	defer round.Release()

	round.Animations = game.BuildAnimations(round)
	isSpin := round.DataKind == rslt.SpinData

	enc := zjson.AcquireEncoder(2048)
	enc.StartObject()

	enc.StartObjectField("roundData")
	if roundID != "" {
		enc.StringField("roundId", roundID)
	}
	if nrOfResults > 0 {
		enc.IntField("nrOfResults", nrOfResults)
	}

	enc.Uint8Field("dataKind", uint8(round.DataKind))
	enc.Int64Field("balanceBefore", startBalance+roundResult.BalanceBefore)
	enc.Int64Field("balanceAfter", startBalance+roundResult.BalanceAfter)
	enc.Int64Field("bet", bet)
	enc.Int64Field("win", win)
	enc.FloatField("winFactor", winFactor, 'f', 2)
	enc.IntBoolField("realWin", win > bet)
	enc.Int64FieldOpt("bonusWin", bonusWin)
	enc.FloatField("bonusFactor", bonusFactor, 'f', 2)
	enc.IntBoolField("realBonusWin", bonusWin > bet)
	enc.Int64FieldOpt("spinWin", spinWin)
	enc.FloatFieldOpt("spinFactor", spinFactor, 'f', 2)
	enc.IntBoolField("realSpinWin", spinWin > bet)
	enc.Int64Field("totalWin", totalWin)
	enc.FloatField("totalFactor", totalFactor, 'f', 2)
	enc.IntBoolField("realTotalWin", totalWin > bet)
	enc.IntBoolFieldOpt("maxPayout", roundResult.MaxPayout > 0)
	enc.EndObject()

	if isSpin {
		enc.StartObjectField("spinData")
		round.Encode2(enc) // do not encode the PRNG log!
		enc.EndObject()
	} else {
		enc.StartObjectField("data")
		round.Encode2(enc) // do not encode the PRNG log!
		enc.EndObject()
	}

	enc.BoolField("success", true)

	enc.EndObject()
	return enc
}
