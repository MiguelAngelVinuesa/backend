package rng_magic

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ber"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/bot"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/btr"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ccb"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/fpr"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/frm"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/lam"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/mgd"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/ofg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-config.git/slots/owl"
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/ai"
	ccbai "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/ai/ccb"
	lamai "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/ai/lam"
	owlai "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/ai/owl"
)

func RunRngStatement(g *game.Regular, gameNR tg.GameNR, statement string, count, timeout int) ([]*RngMagic, error) {
	conditions, symbols := GameData(gameNR)
	if len(conditions) == 0 || symbols == nil {
		return nil, fmt.Errorf("no conditions or symbols found")
	}

	matcher, err := parse(statement, conditions, symbols)
	if err != nil {
		return nil, err
	}

	var newMatcher func(key string, params map[string]any, game *game.Regular) magic.Matcher
	switch gameNR {
	case tg.BOTnr:
		newMatcher = bot.MakeMatcher
	case tg.CCBnr:
		newMatcher = ccb.MakeMatcher
	case tg.MGDnr:
		newMatcher = mgd.MakeMatcher
	case tg.LAMnr:
		newMatcher = lam.MakeMatcher
	case tg.OWLnr:
		newMatcher = owl.MakeMatcher
	case tg.FRMnr:
		newMatcher = frm.MakeMatcher
	case tg.OFGnr:
		newMatcher = ofg.MakeMatcher
	case tg.FPRnr:
		newMatcher = fpr.MakeMatcher
	case tg.BTRnr:
		newMatcher = btr.MakeMatcher
	case tg.BERnr:
		newMatcher = ber.MakeMatcher
	default:
		return nil, fmt.Errorf("cannot find MakeMatcher() function")
	}

	matcher.Init(FunctionInitParams{
		g:          g,
		newMatcher: newMatcher,
		conditions: conditions,
	})

	r := &runner{
		count:   count,
		g:       g,
		matcher: matcher,
	}

	switch gameNR {
	case tg.CCBnr:
		r.stickyAI = ccbai.NewLazySmartPlayerHigh()
	case tg.LAMnr:
		r.choiceAI = lamai.NewRandomPlayer()
	case tg.OWLnr:
		r.choiceAI = owlai.NewRandomPlayer()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ctx context.Context, r *runner, wg *sync.WaitGroup) {
		r.run(ctx, wg)
	}(ctx, r, &wg)

	wg.Wait()

	return r.magics, r.err
}

type runner struct {
	count    int
	g        *game.Regular
	matcher  Function
	choiceAI ai.ChoiceMaker
	stickyAI ai.StickyChooser
	magics   []*RngMagic
	err      error
}

func (r *runner) run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var magic1, magic2 []int
	var res rslt.Results

	for len(r.magics) < r.count {
		select {
		case <-ctx.Done():
			return

		default:
			if res2 := r.g.Round(0); len(res2) > 0 {
				for ix := range res2 {
					res = append(res, res2[ix].Clone().(*rslt.Result))
				}

				magic1 = magics(r.g.PrngLog())
				magic2 = nil

				switch {
				case r.g.NeedPlayerChoice():
					var choices map[string]string
					if r.choiceAI != nil {
						choices = r.choiceAI.Choices()
					}

					res2 = r.g.RoundResume(choices)
					if len(res2) == 0 || r.g.SpinState() != nil {
						r.err = fmt.Errorf("player choice messed up")
						return
					}

					for iy := range res2 {
						res = append(res, res2[iy].Clone().(*rslt.Result))
					}
					magic2 = magics(r.g.PrngLog())

				case r.g.IsDoubleSpin() && len(res) == 1:
					if r.stickyAI != nil {
						chooseSticky(res[0], r.stickyAI)
					}

					res2 = r.g.Round(0)
					if len(res2) == 0 || r.g.SpinState() != nil {
						r.err = fmt.Errorf("double spin messed up")
						return
					}

					for iy := range res2 {
						res = append(res, res2[iy].Clone().(*rslt.Result))
					}
					magic2 = magics(r.g.PrngLog())
				}

				if r.matcher.Matches(res) {
					r.magics = append(r.magics, &RngMagic{
						Magic1:      magic1,
						Magic2:      magic2,
						NrOfResults: len(res),
						TotalPayout: rslt.GrandTotal2(res, r.g.MaxPayout()),
					})
				}

				res = rslt.ReleaseResults(res)
			}
		}
	}
}

func magics(l1, l2 []int) []int {
	out := make([]int, 0, len(l1)*2)
	for ix := range l1 {
		out = append(out, l1[ix], l2[ix])
	}
	return out
}

func chooseSticky(result *rslt.Result, c ai.StickyChooser) {
	spin, ok := result.Data.(*comp.SpinResult)
	if !ok {
		panic(fmt.Sprintf("bad input; not a spin result: %+v", result))
	}

	state, ok2 := result.State.(*comp.SymbolsState)
	if !ok2 {
		panic(fmt.Sprintf("bad input; no spin state: %+v", result))
	}

	sticky := c.SelectSticky(spin.Initial(), spin.StickySymbol(), state.Flagged())
	spin.SetChosenSticky(sticky)
}
