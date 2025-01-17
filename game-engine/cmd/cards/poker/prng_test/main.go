package main

import (
	"fmt"
	"runtime"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards/poker"
)

func main() {
	count := 267569200
	fmt.Printf("5 cards (%d hands):\n", count)
	analyse(count, 5)
	fmt.Println()
	fmt.Printf("7 cards (%d hands):\n", count)
	analyse(count, 7)
}

func analyse(max, draw int) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	threads := 20
	wg := sync.WaitGroup{}

	var mutex sync.Mutex
	totals := newTotals()

	for ix := 0; ix < threads; ix++ {
		wg.Add(1)
		go func(max, draw int, wg *sync.WaitGroup) {
			tot := run(max, draw)
			mutex.Lock()
			for k, v := range tot {
				totals[k] += v
			}
			mutex.Unlock()
			wg.Done()
		}(max/threads, draw, &wg)
	}

	wg.Wait()

	for r := poker.Rank(0); r <= poker.RoyalFlush; r++ {
		v := totals[r]
		fmt.Printf("%16s   frequency %9d   probability %10.6f\n", r.String(), v, float64(v)*100/float64(max))
	}
}

func newTotals() map[poker.Rank]int {
	total := make(map[poker.Rank]int)
	for r := poker.Rank(0); r <= poker.RoyalFlush; r++ {
		total[r] = 0
	}
	return total
}

func run(max, draw int) map[poker.Rank]int {
	d := cards.NewDeck(cards.StandardDeck())
	r := poker.NewRanked(nil)
	totals := newTotals()
	work := make(cards.Cards, draw)

	for iy := 0; iy < max; iy++ {
		d.Shuffle()
		hand := d.DrawInto(work)
		r.Reload(hand)
		totals[r.Rank()]++
	}

	r.Release()
	d.Release()
	return totals
}
