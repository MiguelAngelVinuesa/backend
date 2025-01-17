package slots

import (
	"fmt"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SpinResultsReport generates a report of the spin results as a string.
func SpinResultsReport(in results.Results, f SpinDataFormatter, maxPayout float64) []string {
	b := strings.Builder{}
	b.Grow(4096)

	gt := results.GrandTotal2(in, maxPayout)
	b.WriteString(fmt.Sprintf("#### high-payout round; total payout factor: %.2f; total spins %d ####\n", gt, len(in)))

	for ix := range in {
		result := in[ix]

		if data, ok := result.Data.(*SpinResult); ok {
			spinData(result, data, f, &b, ix)
		}
		if data, ok := result.Data.(*results.InstantBonus); ok {
			instantBonusData(data, &b, ix)
		}
		if data, ok := result.Data.(*results.BonusSelector); ok {
			bonusSelectorData(data, &b, ix)
		}
		if data, ok := result.Data.(*wheel.BonusWheelResult); ok {
			bonusWheelData(data, &b, ix)
		}
	}

	out := strings.Split(b.String(), "\n")
	return out[:len(out)-1]
}

func spinData(result *results.Result, data *SpinResult, f SpinDataFormatter, b *strings.Builder, ix int) {
	initial := f(data.initial)

	var multipliers, afterExpand, afterClear, afterNudge []string
	if len(data.multipliers) == len(data.initial) {
		temp := make(utils.Indexes, len(data.multipliers))
		for iy := range data.multipliers {
			temp[iy] = utils.Index(data.multipliers[iy])
		}
		multipliers = f(temp)
	}
	if len(data.afterExpand) == len(data.initial) {
		afterExpand = f(data.afterExpand)
	}
	if len(data.afterClear) == len(data.initial) {
		afterClear = f(data.afterClear)
	}
	if len(data.afterNudge) == len(data.initial) {
		afterNudge = f(data.afterNudge)
	}

	// header line
	b.WriteString("spin  ")
	b.WriteString(fmt.Sprintf("%-*s", len(initial[0]), "initial"))

	if len(multipliers) > 0 {
		b.WriteString("  ")
		b.WriteString(fmt.Sprintf("%-*s", len(multipliers[0]), "multipliers"))
	}
	if len(afterExpand) > 0 {
		b.WriteString("  ")
		b.WriteString(fmt.Sprintf("%-*s", len(afterExpand[0]), "after expand"))
	}
	if len(afterClear) > 0 {
		b.WriteString("  ")
		b.WriteString(fmt.Sprintf("%-*s", len(afterClear[0]), "after clear"))
	}
	if len(afterNudge) > 0 {
		b.WriteString("  ")
		b.WriteString(fmt.Sprintf("%-*s", len(afterNudge[0]), "after nudge"))
	}

	b.WriteByte('\n')

	// line 1
	b.WriteString(fmt.Sprintf("%-6d", ix+1))
	b.WriteString(initial[0])

	if len(multipliers) > 0 {
		b.WriteString("  ")
		b.WriteString(multipliers[0])
	}
	if len(afterExpand) > 0 {
		b.WriteString("  ")
		b.WriteString(afterExpand[0])
	}
	if len(afterClear) > 0 {
		b.WriteString("  ")
		b.WriteString(afterClear[0])
	}
	if len(afterNudge) > 0 {
		b.WriteString("  ")
		b.WriteString(afterNudge[0])
	}

	if data.bonusSymbol != utils.MaxIndex {
		b.WriteString("  bonus symbol ")
		b.WriteString(strconv.Itoa(int(data.bonusSymbol)))
	}
	if data.stickySymbol != utils.MaxIndex {
		b.WriteString("  sticky symbol ")
		b.WriteString(strconv.Itoa(int(data.stickySymbol)))
	}
	if data.superSymbol != utils.MaxIndex {
		b.WriteString("  super symbol ")
		b.WriteString(strconv.Itoa(int(data.superSymbol)))
	}
	if result.AwardedFreeGames > 0 {
		b.WriteString("  free spins awarded ")
		b.WriteString(strconv.Itoa(int(result.AwardedFreeGames)))
	}

	b.WriteByte('\n')

	// remaining rows.
	for line := 1; line < len(initial); line++ {
		b.WriteString("      ")
		b.WriteString(initial[line])

		if len(multipliers) > 0 {
			b.WriteString("  ")
			b.WriteString(multipliers[line])
		}
		if len(afterExpand) > 0 {
			b.WriteString("  ")
			b.WriteString(afterExpand[line])
		}
		if len(afterClear) > 0 {
			b.WriteString("  ")
			b.WriteString(afterClear[line])
		}
		if len(afterNudge) > 0 {
			b.WriteString("  ")
			b.WriteString(afterNudge[line])
		}

		if line == 1 {
			if len(data.lockedReels) > 0 {
				b.WriteString("  locked reels")
				b.WriteString(fmt.Sprintf("%v", data.lockedReels))
			}
			if len(data.hotReels) > 0 {
				b.WriteString("  hot reels")
				b.WriteString(fmt.Sprintf("%v", data.hotReels))
			}
		}

		b.WriteByte('\n')
	}
}

func instantBonusData(data *results.InstantBonus, b *strings.Builder, ix int) {
	choice := data.Choice()
	l := len(choice)
	options := strings.Replace(fmt.Sprintf("%v", data.Options()), " ", ",", -1)

	// header lines
	b.WriteString("instant bonus selection request\n")
	b.WriteString("spin    ")
	b.WriteString(fmt.Sprintf("%-*s", l, "choice"))
	b.WriteString("  ")
	b.WriteString("options\n")

	// line 1
	b.WriteString(fmt.Sprintf("%-6d", ix+1))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("%-*s", l, choice))
	b.WriteString("  ")
	b.WriteString(options)
	b.WriteByte('\n')
}

func bonusSelectorData(data *results.BonusSelector, b *strings.Builder, ix int) {
	res := strings.Replace(fmt.Sprintf("%v", data.Results()), " ", ",", -1)
	l := len(res)

	// header lines
	b.WriteString("selected bonus\n")
	b.WriteString("spin    index  ")
	b.WriteString(fmt.Sprintf("%-*s", l, "results"))
	b.WriteString("  selected")
	b.WriteByte('\n')

	// line 1
	b.WriteString(fmt.Sprintf("%-6d", ix+1))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("%-5d", data.PlayerChoice()))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("%-*s", l, res))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("%-8d", data.Chosen()))
	b.WriteByte('\n')
}

func bonusWheelData(data *wheel.BonusWheelResult, b *strings.Builder, ix int) {
	// header lines
	b.WriteString("bonus wheel\n")
	b.WriteString("spin    result\n")

	// line 1
	b.WriteString(fmt.Sprintf("%-6d", ix+1))
	b.WriteString("  ")
	b.WriteString(fmt.Sprintf("%d", data.Result()))
	b.WriteByte('\n')
}

type SpinDataFormatter func(data utils.Indexes) []string

func SpinDataRectangle(reels, rows int) SpinDataFormatter {
	return func(data utils.Indexes) []string {
		lines := make([]string, rows)
		for row := 0; row < rows; row++ {
			var b strings.Builder
			b.WriteByte('[')
			for reel := 0; reel < reels; reel++ {
				s := data[reel*rows+row]
				if s < 10 {
					b.WriteByte(' ')
				}
				b.WriteString(strconv.Itoa(int(s)))
				if reel < reels-1 {
					b.WriteByte(' ')
				}
			}
			b.WriteByte(']')
			lines[row] = b.String()
		}
		return lines
	}
}
