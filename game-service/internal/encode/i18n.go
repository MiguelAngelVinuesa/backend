package encode

import (
	"fmt"
	"math"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	l10n "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/clients/i18n"
)

func addLocalizedMessages(bet int64, r *rslt.Result, params *models.PrefetchI18n) {
	if r.Total == 0 {
		return
	}

	strings, err := i18n.GetStrings(params.Locale, getKeys(params))
	if err != nil {
		return
	}

	loc := l10n.NewLocalizer(params.Locale)

	r.SetMessage(loc.Localize(strings[0], map[string]string{
		"currency": params.Currency,
		"amount":   fmt.Sprintf("%.2f", amount(bet, r.Total))}))

	for i := range r.Animations {
		a := r.Animations[i]

		if p, ok := a.(*comp.PayoutEvent); ok {
			p.SetMessage(loc.Localize(strings[1], map[string]string{
				"currency":      params.Currency,
				"count_symbols": fmt.Sprintf("%d", p.Count()),
				"amount":        fmt.Sprintf("%.2f", amount(bet, p.Factor())),
			}))
		}

		if p, ok := a.(*comp.AllPaylinesEvent); ok {
			p.SetMessage(loc.Localize(strings[2], map[string]string{
				"currency":      params.Currency,
				"amount":        fmt.Sprintf("%.2f", amount(bet, p.Factor())),
				"count_symbols": fmt.Sprintf("%d", p.Count()),
				"count":         fmt.Sprintf("%d", p.Paylines()),
			}))
		}
	}

	// if r.AwardedFreeGames > 0 {
	//
	// }

	// if r.FreeGames > 0 {
	//
	// }
}

func getKeys(params *models.PrefetchI18n) []string {
	keys := make([]string, 0, 8)

	add := func(key, dflt string) {
		if key == "" {
			key = dflt
		}
		keys = append(keys, "plural."+key)
	}

	add(params.TotalMsg, "game-pays")
	add(params.PayoutMsg, "payline-win")
	add(params.AllpaylinesMsg, "allpayline-win")
	add(params.FreeSpinsMsg, "won-free-spins")
	add(params.FreeSpinsLabel, "free-counter")

	return keys
}

func amount(bet int64, factor float64) float64 {
	return math.Round(float64(bet)*factor) / 100
}
