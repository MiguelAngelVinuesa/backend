package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards/poker"
)

func main() {
	testFile("../testdata/poker/poker-hand-testing.data")
	testFile("../testdata/poker/poker-hand-training-true.data")
}

func testFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic("cannot open file")
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {
		fields := strings.Split(scan.Text(), ",")
		if len(fields) == 11 {
			list := make(cards.Cards, 5)
			for ix := 0; ix < 5; ix++ {
				list[ix] = convertCard(fields[ix*2], fields[ix*2+1])
			}

			rank := poker.RankHand(list)

			n, _ := strconv.Atoi(fields[10])
			want, got := poker.Rank(n+1), rank.Rank()
			if want != got {
				fmt.Printf("%s != %s\n", got, want)
			}

			rank.Release()

			for _, c := range list {
				c.Release()
			}
		}
	}
}

func convertCard(suit, ordinal string) *cards.Card {
	s, _ := strconv.Atoi(suit)
	o, _ := strconv.Atoi(ordinal)

	switch s {
	case 1:
		s = 32
	case 2:
		s = 48
	case 3:
		s = 0
	case 4:
		s = 16
	default:
		panic("invalid suit in test data")
	}

	id := cards.CardID(o + s)
	return cards.NewCard(id)
}
