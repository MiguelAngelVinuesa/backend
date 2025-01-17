package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
)

func TestWeightedSymbolReels_insertScatters(t *testing.T) {
	const numRuns = 10000
	const screenHeight = 7
	const screenWidth = 7

	scatwei := []int{7800, 2300, 840, 210, 32, 4, 1}

	wr := &WeightedSymbolReels{}
	wr.scatwei = scatwei
	wr.prng = rng.AcquireRNG()

	// Distribution tracking
	scatterCounts := make(map[int]int)

	// Run the test multiple times
	for run := 0; run < numRuns; run++ {
		// Generate a random scr
		screen := make([][]int, screenWidth)
		for i := range screen {
			screen[i] = make([]int, screenHeight)
			for j := range screen[i] {
				screen[i][j] = wr.prng.IntN(8) // Random number 0-7
			}
		}
		// Call insertScatters
		wr.insertScatters(screen, wr.scatwei)

		// Count scatters (value 8)
		scatterCount := 0
		for _, col := range screen {
			for _, val := range col {
				if val == 8 {
					scatterCount++
				}
			}
		}

		// Record the scatter count
		scatterCounts[scatterCount]++
	}

	// Output results
	t.Logf("Scatter distribution after %d runs:", numRuns)
	for n, count := range scatterCounts {
		t.Logf("%d scatters: %d occurrences (%.2f%%)", n, count, float64(count)*100/float64(numRuns))
	}
}

/*
func TestReelUpdater_FillIndexesFrom2DArray(t *testing.T) {
	// Arrange
	screen := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{7, 8, 9},
		{7, 8, 9},
		{7, 8, 9},
		{7, 8, 9},
	}
	expectedReels := []utils.Indexes{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{7, 8, 9},
		{7, 8, 9},
		{7, 8, 9},
		{7, 8, 9},
	}

	var updater ReelUpdater

	// Act
	updater.FillIndexesFrom2DArray(screen)

	// Assert
	if !reflect.DeepEqual(updater.reels, expectedReels) {
		t.Errorf("Expected reels to be %v, but got %v", expectedReels, updater.reels)
	}
}

func TestReelUpdater_Spin(t *testing.T) {
	// Arrange
	screen := [][]int{
		{1, 2, 3, 4},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4},
	}
	expectedMapped := utils.Indexes{
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(0), utils.Index(0), utils.Index(0),
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(3), utils.Index(0), utils.Index(0),
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(3), utils.Index(2), utils.Index(0),
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(3), utils.Index(2), utils.Index(1),
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(3), utils.Index(2), utils.Index(0),
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(3), utils.Index(0), utils.Index(0),
		utils.Index(7), utils.Index(6), utils.Index(5), utils.Index(4), utils.Index(0), utils.Index(0), utils.Index(0),
	}

	var updater ReelUpdater
	updater.FillIndexesFrom2DArray(screen)

	// Mock dependencies
	spin := &Spin{
		rowCount: 7,
		mask:     utils.UInt8s{4, 5, 6, 7, 6, 5, 4},
	}
	out := make(utils.Indexes, 7*7)

	// Act
	updater.Spin(spin, out)

	// Assert
	if !reflect.DeepEqual(out, expectedMapped) {
		t.Errorf("Expected output to be %v, but got %v", expectedMapped, out)
	}
}

func TestReelUpdater_FillEmptySlots(t *testing.T) {
	// Arrange
	screen := [][]int{
		{1, 2, 3, 4},
		{1, 2, 3, 4, 5},
		{1, 4, 4, 4, 5, 6},
		{1, 4, 4, 4, 6, 6, 6},
		{1, 4, 4, 4, 6, 6},
		{1, 2, 3, 4, 5},
		{1, 2, 3, 4},
	}

	wr := &WeightedSymbolReels{}
	wr.prng = rng.AcquireRNG()
	wr.paytable = [][]int{
		{6, 10, 15, 20, 30, 40, 50, 60, 80, 100, 150, 150, 150, 150, 400, 400, 400, 400, 600, 600, 600, 600, 800, 800, 800, 800, 800, 1000, 1000, 1000, 1000, 1000, 2000},
		{4, 6, 8, 10, 15, 20, 25, 30, 40, 50, 80, 80, 80, 80, 200, 200, 200, 200, 300, 300, 300, 300, 400, 400, 400, 400, 400, 600, 600, 600, 600, 600, 1000},
		{3, 4, 6, 8, 10, 15, 20, 25, 30, 40, 60, 60, 60, 60, 150, 150, 150, 150, 200, 200, 200, 200, 300, 300, 300, 300, 300, 500, 500, 500, 500, 500, 800},
		{2, 3, 4, 6, 8, 10, 15, 20, 25, 30, 40, 40, 40, 40, 100, 100, 100, 100, 150, 150, 150, 150, 250, 250, 250, 250, 250, 400, 400, 400, 400, 400, 600},
		{1, 1, 2, 2, 3, 4, 5, 6, 8, 10, 15, 15, 15, 15, 40, 40, 40, 40, 60, 60, 60, 60, 100, 100, 100, 100, 100, 150, 150, 150, 150, 150, 250},
		{1, 1, 2, 2, 3, 4, 5, 6, 8, 10, 15, 15, 15, 15, 40, 40, 40, 40, 60, 60, 60, 60, 100, 100, 100, 100, 100, 150, 150, 150, 150, 150, 250},
		{1, 1, 2, 2, 3, 4, 5, 6, 8, 10, 15, 15, 15, 15, 40, 40, 40, 40, 60, 60, 60, 60, 100, 100, 100, 100, 100, 150, 150, 150, 150, 150, 250},
		{1, 1, 2, 2, 3, 4, 5, 6, 8, 10, 15, 15, 15, 15, 40, 40, 40, 40, 60, 60, 60, 60, 100, 100, 100, 100, 100, 150, 150, 150, 150, 150, 250},
	}
	wr.swr1 = []int{17, 18, 19, 20, 21, 21, 21, 21, 0}
	wr.alfa = []int{1600, 0, 0, 0, 0}
	wr.beta = []int{10, 10, 10, 10, 10}
	wr.gama = []float32{0.018, 0.016, 0.014, 0.012, 0.010}
	updater := NewReelUpdater()
	wr.reelsUpdater = updater

	dummy := NewReelUpdater()

	tClusters := findClusters(screen, wr.paytable)
	tEmpty := wr.findEmpytSlots(screen, tClusters)
	tHexagon := wr.cascadeInHexagon(screen, tEmpty, 0)
	tClusters = findClusters(tHexagon, wr.paytable)

	// Mock dependencies
	spin := &Spin{
		rowCount: 7,
		mask:     utils.UInt8s{4, 5, 6, 7, 6, 5, 4},
	}
	out := make(utils.Indexes, 7*7)

	// Act
	wr.reelsUpdater.FillIndexesFrom2DArray(tHexagon)
	wr.reelsUpdater.Spin(spin, out)

	dummyOut := make(utils.Indexes, 7*7)
	dummy.FillIndexesFrom2DArray(tHexagon)
	var offs int
	for ix := range dummy.reels {
		end := offs + int(spin.mask[ix])
		mapped, multi := MapReelSymbolsValues(dummy.reels[ix])
		spin.multiplier *= float64(multi)
		copy(dummyOut[offs:end], mapped)
		offs += spin.rowCount
	}

	// Assert
	if !reflect.DeepEqual(out, dummyOut) {
		t.Errorf("Expected output to be %v, but got %v; \n hexagon %v", dummyOut, out, tHexagon)
	}
}
*/
