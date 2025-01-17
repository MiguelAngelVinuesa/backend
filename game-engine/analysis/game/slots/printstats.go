package slots

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/metrics/slots"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// TODO Boyan Donchev:
// Temp implementation mimicing the `simulators` service stats printing. Ideally most of that logic will be moved
// to the `game-engine` service so we can invert the flow of doing changes and testing the math - by having the stats
// printing in this service the changes and testing needed in the `game-config` service can be done without the need to
// push and/or redeploy whenever changes are done to a specific games configs, instead running the slot_test for that
// game can be used with an option to print out a text/csv file to do the necessary math checks.

var factors = []float64{0, 0.3, 0.5, 1, 2, 3, 4, 5, 10, 15, 20, 25, 30, 45, 60, 75, 90,
	100, 125, 150, 175, 200, 250, 300, 400, 500, 600, 700, 800, 900,
	1000, 1250, 1500, 1750, 2000, 2500, 3000, 4000, 5000, 5500, 6000, 7000, 8000, 9000,
	10000, 12500, 15000, 17500, 20000, 25000, 30000, 35000, 40000, 45000, 50000}

var formatter = message.NewPrinter(language.English)

func getBandedCounts(bet float64, input map[int64]uint64) ([]uint64, []int64) {
	counts := make([]uint64, len(factors))
	totals := make([]int64, len(factors))

	statsSize := len(factors)

	for k, v := range input {
		factor := float64(k) / bet
		for ix := 0; ix < statsSize; ix++ {
			if factor < factors[ix] || (factor <= factors[ix] && factors[ix] == 0) {
				counts[ix] += v
				totals[ix] += k * int64(v)
				break
			}
		}
	}

	ix := len(counts) - 1
	for ix >= 0 && counts[ix] == 0.0 {
		ix--
	}
	ix++

	return counts[:ix], totals[:ix]
}

func pct(a, b uint64, prec int) string {
	if a == 0 || b == 0 {
		return "0%"
	}
	return fmt.Sprintf("%.*f%%", prec, float64(a)*100.0/float64(b))
}

func revpct(a, b uint64, prec int) string {
	if a == 0 || b == 0 {
		return "0%"
	}
	return fmt.Sprintf("%.*f%%", prec, 100*(1-float64(a)/float64(b)))
}

func printBlock(writer *bufio.Writer, prt Reporter, title string, b SimKVs) {
	cols := len(title) + 10
	for ix := range b {
		row := b[ix]
		if len(row.Values) > cols {
			cols = len(row.Values)
		}
	}
	cols++

	colMax := make([]int, cols)
	for ix := range b {
		row := b[ix]
		if len(row.Key) > colMax[0] {
			colMax[0] = len(row.Key)
		}

		for iy := range row.Values {
			if l := len(row.Values[iy]); l > colMax[iy+1] {
				colMax[iy+1] = l
			}
		}
	}

	lineMax := 2*len(colMax) - 1
	for ix := range colMax {
		lineMax += colMax[ix]
	}

	// writes to file
	_, err := writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("#### %s ####\n", title))
	if err != nil {
		fmt.Println("Error writing block:", err)
	}
	_, err = writer.WriteString(strings.Repeat("=", lineMax) + "\n")
	if err != nil {
		fmt.Println("Error writing block:", err)
	}

	parms := make([]any, 0, 8)
	var f string

	for ix := range b {
		row := b[ix]
		parms = parms[:0]
		f = "%-*s:"
		parms = append(parms, colMax[0], row.Key)

		for iy := range row.Values {
			f += "  %*s"
			parms = append(parms, colMax[iy+1], row.Values[iy])
		}

		// writes to file
		_, err := fmt.Fprintf(writer, f, parms...)
		if err != nil {
			fmt.Println("Error writing formatted text:", err)
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			fmt.Println("Error writing empty line:", err)
		}
	}

	// writes to file
	_, err = writer.WriteString(strings.Repeat("=", lineMax) + "\n")
	if err != nil {
		fmt.Println("Error writing block:", err)
	}
}

func printTable(writer *bufio.Writer, prt Reporter, title string, tab *SimTable, totals bool) {
	colMax := make([]int, len(tab.Keys))
	for ix := range tab.Keys {
		if l := len(tab.Keys[ix]); l > colMax[ix] {
			colMax[ix] = l
		}
	}

	for r := range tab.Values {
		row := tab.Values[r]
		for ix := range row {
			if l := len(row[ix]); l > colMax[ix] {
				colMax[ix] = l
			}
		}
	}

	lineMax := 2*len(colMax) - 2
	for ix := range colMax {
		lineMax += colMax[ix]
	}

	if l := len(title) + 10; l > lineMax {
		lineMax = l
	}

	// writes to file
	_, err := writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("#### %s ####\n", title))
	if err != nil {
		fmt.Println("Error writing block:", err)
	}
	_, err = writer.WriteString(strings.Repeat("=", lineMax) + "\n")
	if err != nil {
		fmt.Println("Error writing block:", err)
	}

	var f string
	parms := make([]any, 0, len(colMax)*2)

	for ix := range tab.Keys {
		if ix == 0 {
			f = "%*s"
		} else {
			f += "  %*s"
		}
		parms = append(parms, colMax[ix], tab.Keys[ix])
	}

	// writes to file
	_, err = fmt.Fprintf(writer, f, parms...)
	if err != nil {
		fmt.Println("Error writing formatted text:", err)
	}

	// writes to file
	_, err = writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(strings.Repeat("-", lineMax) + "\n")
	if err != nil {
		fmt.Println("Error writing block:", err)
	}

	last := len(tab.Values) - 1
	for r := range tab.Values {
		if totals && r == last {
			// writes to file
			_, err := writer.WriteString("\n")
			if err != nil {
				fmt.Println("Error writing empty line:", err)
			}
			_, err = writer.WriteString(strings.Repeat("-", lineMax) + "\n")
			if err != nil {
				fmt.Println("Error writing block:", err)
			}
			_, err = writer.WriteString("\n")
			if err != nil {
				fmt.Println("Error writing empty line:", err)
			}
		}

		parms = parms[:0]
		row := tab.Values[r]

		for ix := range row {
			if ix == 0 {
				f = "%*s"
			} else {
				f += "  %*s"
			}
			parms = append(parms, colMax[ix], row[ix])
		}

		// writes to file
		_, err = fmt.Fprintf(writer, f, parms...)
		if err != nil {
			fmt.Println("Error writing formatted text:", err)
		}
		_, err := writer.WriteString("\n")
		if err != nil {
			fmt.Println("Error writing empty line:", err)
		}
	}

	// writes to file
	_, err = writer.WriteString(strings.Repeat("=", lineMax) + "\n")
	if err != nil {
		fmt.Println("Error writing block:", err)
	}
}

func BetTotals(all *slots.Rounds) SimKVs {
	bets := all.Bets
	return SimKVs{
		{Key: "Total bets", Values: []string{formatter.Sprintf("%.2f€", float64(bets.Total)/100)}},
		{Key: "Minimum bet", Values: []string{formatter.Sprintf("%.2f€", float64(bets.Min)/100)}},
		{Key: "Maximum bet", Values: []string{formatter.Sprintf("%.2f€", float64(bets.Max)/100)}},
	}
}

func BetSpread(all *slots.Rounds) *SimTable {
	out := &SimTable{Keys: []string{"bet", "count"}}
	for k, v := range all.Bets.Counts {
		out.Values = append(out.Values, []string{formatter.Sprintf("%.2f€", float64(k)/100), formatter.Sprintf("%d", v)})
	}
	return out
}

func WinSpread(bet float64, table map[int64]uint64) *SimTable {
	out := &SimTable{Keys: []string{"<win", "<factor", "count", "total amt", "avg payout", "pct"}}

	counts, amounts := getBandedCounts(bet, table)

	var totalCount, totalAmount uint64
	for ix := range counts {
		totalCount += counts[ix]
		totalAmount += uint64(amounts[ix])
	}

	for ix := range counts {
		factor := factors[ix]
		win := int64(factor * bet)
		count, amount := counts[ix], amounts[ix]
		if count == 0 {
			count = 1
		}

		out.Values = append(out.Values, []string{
			formatter.Sprintf("%.2f€", float64(win)/100),
			formatter.Sprintf("%.1fx", factor),
			formatter.Sprintf("%d", counts[ix]),
			formatter.Sprintf("%.2f€", float64(amount)/100),
			formatter.Sprintf("%.2f€", float64(amount)/float64(count)/100),
			pct(counts[ix], totalCount, 6),
		})
	}

	out.Values = append(out.Values, []string{"", "",
		formatter.Sprintf("%d", totalCount),
		formatter.Sprintf("%.2f€", float64(totalAmount)/100),
		formatter.Sprintf("%.2f€", float64(totalAmount)/float64(totalCount)/100),
		pct(totalCount, totalCount, 6),
	})

	return out
}

func RoundSpins(col1 string, counts []uint64) *SimTable {
	out := &SimTable{Keys: []string{col1, "rounds"}}
	for ix, count := range counts {
		if count > 0 {
			out.Values = append(out.Values, []string{
				formatter.Sprintf("%d", ix),
				formatter.Sprintf("%d", count),
			})
		}
	}
	return out
}

func odds(a, b uint64) string {
	if a == 0 || b == 0 {
		return ""
	}

	o := float64(b) / float64(a)
	prec := 2
	switch {
	case o < 1:
		prec = 4
	case o > 100:
		prec = 0
	case o > 10:
		prec = 1
	}

	s := fmt.Sprintf("1 in ~%.*f", prec, o)

	if prec > 0 {
		for strings.HasSuffix(s, "0") {
			s = s[:len(s)-1]
		}
	}

	if strings.HasSuffix(s, ".") {
		s = s[:len(s)-1]
	}

	return s
}

func RoundMultiplierMarks(marks *slots.MinMaxUInt64) *SimTable {
	out := &SimTable{Keys: []string{"mark", "count", "pct", "odds"}}

	list := marks.Counts
	ids := make([]int, 0, len(list))

	var total uint64
	for id := range list {
		total += list[id]
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	for ix := range ids {
		id := uint64(ids[ix])
		out.Values = append(out.Values, []string{
			formatter.Sprintf("%d", id),
			formatter.Sprintf("%d", list[id]),
			pct(list[id], total, 6),
			odds(list[id], total),
		})
	}

	out.Values = append(out.Values, []string{
		"",
		formatter.Sprintf("%d", total),
		pct(total, total, 6),
		"",
	})

	return out
}

func RoundMultipliers(marks *slots.MinMaxFloat64) *SimTable {
	out := &SimTable{Keys: []string{"multiplier", "count", "pct", "odds"}}

	list := marks.Counts
	ids := make([]int, 0, len(list))

	var total uint64
	for id := range list {
		total += list[id]
		ids = append(ids, int(id))
	}
	sort.Ints(ids)

	for ix := range ids {
		id := int64(ids[ix])
		out.Values = append(out.Values, []string{
			formatter.Sprintf("%.2f", float64(id)/float64(marks.Factor)),
			formatter.Sprintf("%d", list[id]),
			pct(list[id], total, 6),
			odds(list[id], total),
		})
	}

	out.Values = append(out.Values, []string{
		"",
		formatter.Sprintf("%d", total),
		pct(total, total, 6),
		"",
	})

	return out
}

func RoundSymbols(totals *Rounds) *SimTable {
	rounds := totals.AllRounds
	symbols := totals.Symbols
	totalRounds := 100.0 / float64(rounds.Count)
	totalNoFree := 100.0 / float64(rounds.Count-totals.FirstTimes)
	totalFree := 100.0 / float64(totals.FirstTimes)

	out := &SimTable{Keys: []string{"name", "rounds", "% total", "without free", "% total", "with free", "% total"}}

	for ix, count1 := range rounds.SymbolsUsed {
		count2 := rounds.SymbolsNoFree[ix]
		count3 := rounds.SymbolsFree[ix]

		if count1 > 0 && symbols[ix] != nil {
			out.Values = append(out.Values, []string{
				symbols[ix].Name,
				formatter.Sprintf("%d", count1),
				formatter.Sprintf("%.2f%%", float64(count1)*totalRounds),
				formatter.Sprintf("%d", count2),
				formatter.Sprintf("%.2f%%", float64(count2)*totalNoFree),
				formatter.Sprintf("%d", count3),
				formatter.Sprintf("%.2f%%", float64(count3)*totalFree),
			})
		}
	}

	return out
}

// RoundsBracketsToFile will be used to printout stats to a txt file
// similar to the implemantation of `PrintRounds` in the `simulators/internal` package
func (r *Rounds) RoundsBracketsToTxt(rtp int, prt Reporter) {
	// STEP 1 - CREATE FILE TO WRITE IN

	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	// Build the path to the Desktop
	fname := fmt.Sprintf("%s-rounds-%d-brackets-%s.txt", r.gameNR.String(), rtp, time.Now().String())
	desktopPath := filepath.Join(homeDir, "Desktop", fname)

	// Create the txt file on the Desktop
	file, err := os.Create(desktopPath)
	if err != nil {
		fmt.Println("Error creating txt file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// STEP 2 - DO THE LOGIC FROM THE SIMULATOR

	printBlock(writer, prt, "bet totals", BetTotals(r.AllRounds))
	printTable(writer, prt, "bet spread", BetSpread(r.AllRounds), false)

	rounds := r.AllRounds

	bet := float64(rounds.Bets.Min+rounds.Bets.Max) / 2.0

	var freeRounds uint64
	var keys []int
	for k, v := range r.BonusRounds {
		if k != slots.NoFreeSpins {
			keys = append(keys, int(k))
			freeRounds += v.Count
		}
	}
	sort.Ints(keys)

	printTable(writer, prt, "overall win statistics", WinSpread(bet, rounds.Wins.Counts), true)

	if v := r.BonusRounds[slots.NoFreeSpins]; v != nil {
		printTable(writer, prt, "win statistics - no free spins", WinSpread(bet, v.Wins.Counts), true)
	}

	for _, ix := range keys {
		k := slots.BonusKind(ix)
		v := r.BonusRounds[k]
		printTable(writer, prt, "win statistics - "+k.String(), WinSpread(bet, v.Wins.Counts), true)
	}

	printTable(writer, prt, "free spins occurrence per round", RoundSpins("free spins", rounds.FreeSpinRounds), false)

	if len(rounds.RefillRounds) > 1 {
		printTable(writer, prt, "refill spins occurrence per round", RoundSpins("refill spins", rounds.RefillRounds), false)
	}

	if len(rounds.SuperRounds) > 1 {
		printTable(writer, prt, "super spins occurrence per round", RoundSpins("super spins", rounds.SuperRounds), false)
	}

	if r.MultiplierMarks.Total > 0 {
		printTable(writer, prt, "round multiplier marks", RoundMultiplierMarks(r.MultiplierMarks), true)
	}

	if r.Multipliers.Total > 0 {
		printTable(writer, prt, "round multipliers", RoundMultipliers(r.Multipliers), true)
	}

	printTable(writer, prt, "symbol occurrence per round", RoundSymbols(r), false)

	fmt.Printf("txt file successfully saved to %s\n", desktopPath)
}

func (r *Rounds) WriteZeroDistributionReport(rtp int, prt Reporter) {
	// STEP 1 - CREATE FILE TO WRITE IN

	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	// Build the path to the Desktop
	fname := fmt.Sprintf("%s-%d-zero-distribution-%s.txt", r.gameNR.String(), rtp, time.Now().String())
	desktopPath := filepath.Join(homeDir, "Desktop", fname)

	// Create the txt file on the Desktop
	file, err := os.Create(desktopPath)
	if err != nil {
		fmt.Println("Error creating txt file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	_, err = writer.WriteString(fmt.Sprintf("Total: %d", r.TotalSpins))
	if err != nil {
		fmt.Println("Error writing total", err)
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("Wins: %d", r.WinCount))
	if err != nil {
		fmt.Println("Error writing wins:", err)
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("Zero Win: %s", revpct(r.WinCount, r.TotalSpins, 6)))
	if err != nil {
		fmt.Println("Error writing zero win %:", err)
	}
	_, err = writer.WriteString("\n\n\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("All rounds total: %d", r.AllRounds.Count))
	if err != nil {
		fmt.Println("Error writing total", err)
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("All rounds zero wins: %d", r.AllRounds.Wins.Counts[0]))
	if err != nil {
		fmt.Println("Error writing wins:", err)
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		fmt.Println("Error writing empty line:", err)
	}
	_, err = writer.WriteString(fmt.Sprintf("all rounds Zero Win: %s", pct(r.AllRounds.Wins.Counts[0], r.AllRounds.Count, 6)))
	if err != nil {
		fmt.Println("Error writing zero win %:", err)
	}

	fmt.Printf("txt file successfully saved to %s\n", desktopPath)
}
