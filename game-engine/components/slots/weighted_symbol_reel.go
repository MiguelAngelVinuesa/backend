package slots

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

type CascadingWeightedSymbolReels struct {
	wr *WeightedSymbolReels
}

func NewCascadingWeightedSymbolReels(wr *WeightedSymbolReels) *CascadingWeightedSymbolReels {
	cwr := &CascadingWeightedSymbolReels{}
	cwr.wr = wr
	return cwr
}

// Spin can be used to fill all reels on the grid with random symbols from the configured symbol reels.
// This function implements the Spinner interface.
func (cwr *CascadingWeightedSymbolReels) Spin(spin *Spin, out utils.Indexes) {
	cwr.wr.Spin(spin, out)
}

// WeighetedSymbolReels represents a dynamic weights reels spinner.
type WeightedSymbolReels struct {
	isFg       bool
	level      int
	wincap     int
	scatwei    []int
	scatweiFG  []int
	wildmulti  []int
	levelbar   []int
	fgspins    []int
	multiwei   []int
	multiweiFG []int
	swr        [][]int
	colwei     [][]int
	colweiFG   [][]int
	colony     [][][]int
	reels      []utils.Indexes
	prng       interfaces.Generator
}

// NewWeightedSymbolReels instantiates new reels with dynamic weights.
func NewWeightedSymbolReels(colony [][][]int, colwei, colweiFG, swr [][]int, multiwei, multiweiFG, fgspins,
	levelbar, wildmulti, scatwei, scatweiFG []int) *WeightedSymbolReels {

	wr := &WeightedSymbolReels{}
	// reg params
	wr.swr = swr
	wr.levelbar = levelbar
	wr.wildmulti = wildmulti

	wr.fgspins = fgspins
	wr.colony = colony
	wr.colwei = colwei
	wr.multiwei = multiwei
	wr.scatwei = scatwei

	// fg params
	wr.colweiFG = colweiFG
	wr.multiweiFG = multiweiFG
	wr.scatweiFG = scatweiFG

	wr.prng = rng.AcquireRNG()
	wr.reels = make([]utils.Indexes, 7)

	return wr
}

type CascadeRefiller struct {
	level    int
	swr1     []int
	alfa     []int
	alfaFG   []int
	beta     []int
	betaFG   []int
	gama     []float32
	gamaFG   []float32
	paytable [][]int
	reels    []utils.Indexes
	prng     interfaces.Generator
}

func NewCascadeRefiller(level int, swr1, alfa, alfaFG, beta, betaFG []int, gama, gamaFG []float32, paytable [][]int) *CascadeRefiller {
	cr := &CascadeRefiller{}
	// base params
	cr.swr1 = swr1
	cr.alfa = alfa
	cr.beta = beta
	cr.gama = gama
	cr.paytable = paytable

	// fg params
	cr.alfaFG = alfaFG
	cr.betaFG = betaFG
	cr.gamaFG = gamaFG

	cr.prng = rng.AcquireRNG()
	cr.reels = make([]utils.Indexes, 7)

	return cr
}

// Spin here is used to refill a grid after cascading.
// This function implements the spinner interface.
func (cr *CascadeRefiller) Spin(spin *Spin, out utils.Indexes) {
	if spin.kind != RefillSpin {
		return
	}
	tHexagon := Fill2DArrayFromIndexes(spin.Indexes())
	tClusters := findClusters(tHexagon, cr.paytable)
	tEmpty := cr.findEmpytSlots(tHexagon, tClusters)
	tHexagon = cr.cascadeInHexagon(tHexagon, tEmpty, cr.level)
	//cr.screen = tHexagon

	cr.reels = FillIndexesFrom2DArray(tHexagon)
	rows, mask := spin.rowCount, spin.mask
	var offs int
	for ix := range cr.reels {
		end := offs + int(mask[ix])
		mapped, multi := MapReelSymbolsValues(cr.reels[ix])
		spin.multiplier *= float64(multi)
		copy(out[offs:end], mapped)
		offs += rows
	}
}

// Spin can be used to fill all reels on the grid with the symbols from the saved grid.
// This function implements the Spinner interface.
func (r *ReelUpdater) Spin(spin *Spin, out utils.Indexes) {
	rows, mask := spin.rowCount, spin.mask
	var offs int
	for ix := range r.reels {
		end := offs + int(mask[ix])
		mapped, multi := MapReelSymbolsValues(r.reels[ix])
		spin.multiplier *= float64(multi)
		copy(out[offs:end], mapped)
		offs += rows
	}
}

// Cluster struct represents a cluster with associated attributes.
type Cluster struct {
	Symbol     int
	Size       int
	Prize      int
	Multiplier int
	Position   [][]int
}

/*
type CascadingSpin struct {
	EndWin         int
	EndLevel       int
	EndMeter       int
	MaxMulti       int
	ProdMulti      int
	NCascades      int
	NScatters      int
	CasHexagons    [][][]int
	CasClusters    [][]Cluster
	CasEmptys      [][][]int
	CasWins        []int
	CasWinsAcc     []int
	CasMeterAcc    [][]int
	CasColonyLevel []int
	CasColonyState []int
}
*/
/*
type Round struct {
	Bet         int
	RoundWin    []int
	BgWin       int
	BaseSpin    *CascadingSpin
	FgTrigg     bool
	FgLength    int
	FgWins      []int
	FgSpins     []int
	FreeGames   []*CascadingSpin
	MaxbgMulti  int
	MaxfgMulti  int
	ProdbgMulti int
	ProdfgMulti int
}
*/

// Helper functions for the dynamic weighted reels refill and spin

// createHexagon creates a hexagonal structure of weighted random integers based on swr values.
func (sr *WeightedSymbolReels) createHexagon() [][]int {
	nx := 7
	ny := []int{4, 5, 6, 7, 6, 5, 4}

	tscreen := make([][]int, nx)
	for i := 0; i < nx; i++ {
		col := make([]int, ny[i])
		for j := 0; j < len(col); j++ {
			col[j] = getWeiRand(sr.swr[i], sr.prng)
		}
		tscreen[i] = col
	}

	return tscreen
}

// InsertScatters inserts scatter symbols (represented by '8' in math, '9' in configs) into random positions on the screen.
func (sr *WeightedSymbolReels) insertScatters(screen [][]int, scwei []int) {
	rils := []int{0, 1, 2, 4, 5, 6}

	nsc := getWeiRand(scwei, sr.prng)

	for i := 0; i < nsc; i++ {
		ind := sr.prng.IntN(len(rils))
		ril := rils[ind]
		y := sr.prng.IntN(len(screen[ril]))

		screen[ril][y] = 8

		// Remove the selected rail to avoid reusing it
		rils = append(rils[:ind], rils[ind+1:]...)
	}
}

// GetWeiRand selects an index based on weighted probabilities
func getWeiRand(weights []int, prng interfaces.Generator) int {
	if len(weights) == 0 {
		return -1
	}

	// Sum up weights
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	// If total weight is zero, return -1
	if totalWeight == 0 {
		return -1
	}

	// Generate a random number in the range of totalWeight
	randomValue := prng.IntN(totalWeight)

	// Find the index where the cumulative weight surpasses the random value
	cumulativeSum := 0
	for i, weight := range weights {
		cumulativeSum += weight
		if cumulativeSum > randomValue {
			return i
		}
	}

	return -1
}

func findClusters(screen [][]int, paytable [][]int) []Cluster {
	var clusters []Cluster

	for symb := 0; symb < 8; symb++ {
		var skrin [][]int
		var pos [][]int

		for i := 0; i < len(screen); i++ {
			bla := make([]int, len(screen[i]))

			for j := 0; j < len(screen[i]); j++ {
				if screen[i][j] > 10 {
					bla[j] = symb
				} else {
					bla[j] = screen[i][j]
				}

				if bla[j] == symb {
					pos = append(pos, []int{i, j})
				}
			}
			skrin = append(skrin, bla)
		}

		n := 0
		for len(pos) > 0 {
			n++
			if n > 1000 {
				fmt.Println("Infinite loop detected")
				break
			}

			clust := [][]int{pos[0]}
			k := 0
			for k < len(clust) {
				nejbrs := findNejbrs(skrin, clust[k], symb)

				clen := len(clust)
				for _, neighbor := range nejbrs {
					j := 0
					for ; j < clen; j++ {
						if clust[j][0] == neighbor[0] && clust[j][1] == neighbor[1] {
							break
						}
					}
					if j == clen {
						clust = append(clust, neighbor)
					}
				}

				k++
				if k > 1000 {
					fmt.Println("Infinite loop detected")
					break
				}
			}

			if len(clust) > 4 {
				// Define the new position matrix.
				newpos := [][]int{
					{0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0},
					{0, 0, 0, 0},
				}

				multi := 1
				for _, pos := range clust {
					x, y := pos[0], pos[1]
					newpos[x][y] = 1
					if screen[x][y] > 10 {
						multi *= screen[x][y] - 10
					}
				}

				newCluster := Cluster{
					Symbol:     symb,
					Prize:      multi * paytable[symb][len(clust)-5],
					Size:       len(clust),
					Multiplier: multi,
					Position:   newpos,
				}
				clusters = append(clusters, newCluster)
			}

			// Remove elements of `clust` from `pos`
			for _, c := range clust {
				for i := 0; i < len(pos); i++ {
					if c[0] == pos[i][0] && c[1] == pos[i][1] {
						pos = append(pos[:i], pos[i+1:]...)
						break
					}
				}
			}
		}
	}

	return clusters
}

func findNejbrs(screen [][]int, spot []int, symb int) [][]int {
	x, y := spot[0], spot[1]
	var nejbrs [][]int

	// Check top and bottom neighbors
	if y-1 >= 0 && screen[x][y-1] == symb {
		nejbrs = append(nejbrs, []int{x, y - 1})
	}
	if y+1 < len(screen[x]) && screen[x][y+1] == symb {
		nejbrs = append(nejbrs, []int{x, y + 1})
	}

	// Check additional neighbors based on x value
	switch x {
	case 0:
		if screen[1][y] == symb {
			nejbrs = append(nejbrs, []int{1, y})
		}
		if y+1 < len(screen[1]) && screen[1][y+1] == symb {
			nejbrs = append(nejbrs, []int{1, y + 1})
		}
	case 1:
		if screen[2][y] == symb {
			nejbrs = append(nejbrs, []int{2, y})
		}
		if y+1 < len(screen[2]) && screen[2][y+1] == symb {
			nejbrs = append(nejbrs, []int{2, y + 1})
		}
		if y < 4 && screen[0][y] == symb {
			nejbrs = append(nejbrs, []int{0, y})
		}
		if y-1 >= 0 && screen[0][y-1] == symb {
			nejbrs = append(nejbrs, []int{0, y - 1})
		}
	case 2:
		if screen[3][y] == symb {
			nejbrs = append(nejbrs, []int{3, y})
		}
		if y+1 < len(screen[3]) && screen[3][y+1] == symb {
			nejbrs = append(nejbrs, []int{3, y + 1})
		}
		if y < 5 && screen[1][y] == symb {
			nejbrs = append(nejbrs, []int{1, y})
		}
		if y-1 >= 0 && screen[1][y-1] == symb {
			nejbrs = append(nejbrs, []int{1, y - 1})
		}
	case 3:
		if y < 6 && screen[2][y] == symb {
			nejbrs = append(nejbrs, []int{2, y})
		}
		if y-1 >= 0 && screen[2][y-1] == symb {
			nejbrs = append(nejbrs, []int{2, y - 1})
		}
		if y < 6 && screen[4][y] == symb {
			nejbrs = append(nejbrs, []int{4, y})
		}
		if y-1 >= 0 && screen[4][y-1] == symb {
			nejbrs = append(nejbrs, []int{4, y - 1})
		}
	case 4:
		if screen[3][y] == symb {
			nejbrs = append(nejbrs, []int{3, y})
		}
		if y+1 < len(screen[3]) && screen[3][y+1] == symb {
			nejbrs = append(nejbrs, []int{3, y + 1})
		}
		if y < 5 && screen[5][y] == symb {
			nejbrs = append(nejbrs, []int{5, y})
		}
		if y-1 >= 0 && screen[5][y-1] == symb {
			nejbrs = append(nejbrs, []int{5, y - 1})
		}
	case 5:
		if screen[4][y] == symb {
			nejbrs = append(nejbrs, []int{4, y})
		}
		if y+1 < len(screen[4]) && screen[4][y+1] == symb {
			nejbrs = append(nejbrs, []int{4, y + 1})
		}
		if y < 4 && screen[6][y] == symb {
			nejbrs = append(nejbrs, []int{6, y})
		}
		if y-1 >= 0 && screen[6][y-1] == symb {
			nejbrs = append(nejbrs, []int{6, y - 1})
		}
	case 6:
		if screen[5][y] == symb {
			nejbrs = append(nejbrs, []int{5, y})
		}
		if y+1 < len(screen[5]) && screen[5][y+1] == symb {
			nejbrs = append(nejbrs, []int{5, y + 1})
		}
	}

	return nejbrs
}

func (sr *CascadeRefiller) cascadeInHexagon(screen, empty [][]int, level int) [][]int {
	// Initialize the new screen with adjusted values
	newscreen := make([][]int, len(screen))
	for i := range screen {
		tcol := make([]int, len(screen[i]))
		for j := range screen[i] {
			if empty[i][j] == 0 {
				tcol[j] = screen[i][j]
			} else {
				tcol[j] = -1
			}
		}
		newscreen[i] = tcol
	}

	index := []int{3, 2, 4, 1, 5, 0, 6}

	for i := range newscreen {
		// Regular weights
		regWeights := make([]int, len(sr.swr1))
		copy(regWeights, sr.swr1)

		// Adjusted weights
		adjWeights := make([]int, len(sr.swr1))

		if index[i] > 0 {
			for _, sym := range newscreen[index[i]-1] {
				if sym >= 0 && sym < 8 {
					adjWeights[sym]++
				}
			}
		}

		if index[i] < 6 {
			for _, sym := range newscreen[index[i]+1] {
				if sym >= 0 && sym < 8 {
					adjWeights[sym]++
				}
			}
		}

		// Combined symbol weights
		symbWeights := make([]int, len(sr.swr1))
		for j := range sr.swr1 {
			symbWeights[j] = sr.alfa[level]*adjWeights[j] + sr.beta[level]*regWeights[j]
		}

		// Scatter weight
		var sum int
		for _, weight := range symbWeights {
			sum += weight
		}
		if index[i] == 3 {
			symbWeights[8] = 0
		} else {
			symbWeights[8] = int(sr.gama[level] * float32(sum))
		}

		// Dropping down symbols
		tcol := make([]int, len(newscreen[index[i]]))
		for j := range screen[index[i]] {
			if empty[index[i]][j] > 10 {
				tcol[j] = empty[index[i]][j]
			} else if empty[index[i]][j] == 0 && screen[index[i]][j] > 10 {
				tcol[j] = screen[index[i]][j]
			}
		}

		var tnew []int
		for j := range screen[index[i]] {
			if empty[index[i]][j] == 0 && screen[index[i]][j] < 10 {
				tnew = append(tnew, screen[index[i]][j])
			}
		}

		for j := range screen[index[i]] {
			if empty[index[i]][j] == 1 {
				newsym := getWeiRand(symbWeights, sr.prng)
				tnew = append(tnew, newsym)
			}
		}

		k := 0
		for j := range screen[index[i]] {
			if tcol[j] < 10 && k < len(tnew) {
				tcol[j] = tnew[k]
				k++
			}
		}

		newscreen[index[i]] = tcol
	}

	return newscreen
}

func (sr *WeightedSymbolReels) insertColony(screen [][]int, clnyType int, wild int, clnies [][][]int) [][]int {
	// Return the screen if colony type is 0
	if clnyType == 0 {
		return screen
	}

	// Initialize the colony with deep copy of clnies[clnyType]
	clny := make([][]int, len(clnies[clnyType]))
	for i := 0; i < len(clnies[clnyType]); i++ {
		bla := make([]int, len(clnies[clnyType][i]))
		copy(bla, clnies[clnyType][i])
		clny[i] = bla
	}

	// Insert colony
	pos := [][]int{}

	for i := 0; i < len(clny); i++ {
		for j := 0; j < len(clny[i]); j++ {
			if clny[i][j] == 0 {
				clny[i][j] = screen[i][j]
			} else if clny[i][j] == 1 {
				if screen[i][j] > 7 {
					clny[i][j] = screen[i][j]
				} else {
					clny[i][j] = screen[3][3]
					pos = append(pos, []int{i, j})
				}
			}
		}
	}

	// Insert wild symbol if wild is greater than 0
	if wild > 0 && len(pos) > 0 {
		ind := sr.prng.IntN(len(pos))
		x, y := pos[ind][0], pos[ind][1]
		clny[x][y] = wild
	}

	return clny
}

func (cr *CascadeRefiller) findEmpytSlots(screen [][]int, clusters []Cluster) [][]int {
	nx := 7
	ny := []int{4, 5, 6, 7, 6, 5, 4}
	empty := make([][]int, nx)

	// Looking for empty slots
	for i := 0; i < nx; i++ {
		tcol := make([]int, ny[i])
		for j := 0; j < ny[i]; j++ {
			val := 0
			for _, cluster := range clusters {
				if cluster.Position[i][j] > 0 {
					val = 1
					break
				}
			}

			if val > 0 && screen[i][j] > 10 {
				val = -1
			}
			tcol[j] = val
		}
		empty[i] = tcol
	}

	// Moving wilds to adjacent spots
	for i := 0; i < nx; i++ {
		for j := 0; j < ny[i]; j++ {
			if empty[i][j] == -1 {
				pos := findNejbrs(empty, []int{i, j}, 1)
				for k := 0; k < len(pos); k++ {
					x, y := pos[k][0], pos[k][1]
					if x == 3 && y == 3 {
						pos = append(pos[:k], pos[k+1:]...)
						break
					}
				}

				if len(pos) > 0 {
					newpos := pos[cr.prng.IntN(len(pos))]
					x, y := newpos[0], newpos[1]
					empty[x][y] = screen[i][j]
					empty[i][j] = 1
				} else {
					empty[i][j] = screen[i][j]
				}
			}
		}
	}

	return empty
}

func (sr *WeightedSymbolReels) Spin(spin *Spin, out utils.Indexes) {
	tHexagon := sr.createHexagon()
	sr.isFg = spin.kind == FreeSpin

	if sr.isFg {
		sr.insertScatters(tHexagon, sr.scatweiFG)
	} else {
		sr.insertScatters(tHexagon, sr.scatwei)
	}

	sr.reels = FillIndexesFrom2DArray(tHexagon)

	rows, mask := spin.rowCount, spin.mask
	var offs int
	for ix := range sr.reels {
		end := offs + int(mask[ix])
		mapped, multi := MapReelSymbolsValues(sr.reels[ix])
		spin.multiplier *= float64(multi)
		copy(out[offs:end], mapped)
		offs += rows
	}
}

/*
	// Cascading
	for colonyCount > 0 {
		colonyCount--

		if colonyState > 0 {
			var wild int
			if colonyState < 4 {
				if sr.isFg {
					wild = sr.wildmulti[getWeiRand(sr.multiweiFG, sr.prng)]
				} else {
					wild = sr.wildmulti[getWeiRand(sr.multiwei, sr.prng)]
				}
			} else {
				wild = 0
			}

			clnytype := 0
			if sr.isFg {
				clnytype = getWeiRand(sr.colweiFG[colonyState], sr.prng)
			} else {
				clnytype = getWeiRand(sr.colwei[colonyState], sr.prng)
			}

			tHexagon := sr.insertColony(tHexagon, clnytype, wild, sr.colony)
			tClusters = findClusters(tHexagon, sr.paytable)

			sr.reelsUpdater.FillIndexesFrom2DArray(tHexagon)
			sr.reelsUpdater.Spin(spin, out)
		}

				for len(tClusters) > 0 {
					casmeter := []int{}
					caswin := 0

					for _, cluster := range tClusters {
						caswin += cluster.Prize
						meterState += cluster.Size
						casmeter = append(casmeter, meterState)

						if colonyLevel < 4 && meterState > sr.levelbar[colonyLevel] {
							colonyLevel++
							colonyCount++
							if colonyLevel < 4 && meterState > sr.levelbar[colonyLevel] {
								colonyLevel++
								colonyCount++
							}
						}
					}
				}

				colonyState++
			}


	// Calculate additional results
	ncascades := len(tcasHexagons)
	_, endlvl, _ := 0, 0, 0
	if ncascades > 1 {
		//endwin = tcasWinsAcc[len(tcasWinsAcc)-1]
		endlvl = tcasColonyLevel[len(tcasColonyLevel)-1]
		//endmeter = tcasMeterAcc[len(tcasMeterAcc)-1][len(tcasMeterAcc[len(tcasMeterAcc)-1])-1]
	}

	// Calculate multipliers and scatters
	nscat, maxmulti, prodmulti := 0, 0, 0
	if endlvl > 0 {
		prodmulti = 1
	}
	for _, row := range tHexagon {
		for _, val := range row {
			if val == 8 {
				nscat++
			}
			if val > 10 {
				mult := val - 10
				if mult > maxmulti {
					maxmulti = mult
				}
				prodmulti *= mult
			}
		}
	}
}
*/

type ReelUpdater struct {
	reels []utils.Indexes
}

func NewReelUpdater() *ReelUpdater {
	ru := &ReelUpdater{}
	ru.reels = make([]utils.Indexes, 7)

	return ru
}

func FillIndexesFrom2DArray(screen [][]int) []utils.Indexes {
	reels := make([]utils.Indexes, 7)
	for i, col := range screen {
		reels[i] = make(utils.Indexes, len(col))
		for j, val := range col {
			el := val
			if el == -1 {
				el = 0
			}
			reels[i][j] = utils.Index(el)
		}
	}

	return reels
}

func MapReelSymbolsValues(reels utils.Indexes) (utils.Indexes, int) {
	mapped := make(utils.Indexes, len(reels))
	copy(mapped, reels)
	multiplier := 1

	for i := range mapped {
		switch mapped[i] {
		case 0:
			mapped[i] = 8
		case 1:
			mapped[i] = 7
		case 2:
			mapped[i] = 6
		case 3:
			mapped[i] = 5
		case 4:
			mapped[i] = 4
		case 5:
			mapped[i] = 3
		case 6:
			mapped[i] = 2
		case 7:
			mapped[i] = 1
		case 8:
			mapped[i] = 9
		case 11, 12, 13, 16, 19:
			mapped[i] = 10
			multiplier *= int(mapped[i]) - 10
		default:
		}
	}

	return mapped, multiplier
}

func Fill2DArrayFromIndexes(in utils.Indexes) [][]int {
	mask := []int{4, 5, 6, 7, 6, 5, 4}

	newScreen := make([][]int, len(mask))
	for i, dim := range mask {
		newScreen[i] = make([]int, dim)
	}

	//multiplierMappings := []uint16{1, 2, 3, 6, 9}

	reverseMapSymbols := func(mapped utils.Indexes) utils.Indexes {
		reversed := make(utils.Indexes, len(mapped))
		copy(reversed, mapped)

		for i := range reversed {
			switch reversed[i] {
			case 8:
				reversed[i] = 0
			case 7:
				reversed[i] = 1
			case 6:
				reversed[i] = 2
			case 5:
				reversed[i] = 3
			case 4:
				reversed[i] = 4
			case 3:
				reversed[i] = 5
			case 2:
				reversed[i] = 6
			case 1:
				reversed[i] = 7
			case 9:
				reversed[i] = 8
			case 10:
				//todo map fully
				reversed[i] = 1 + 10
			default:
			}
		}

		return reversed
	}

	mappedReels := reverseMapSymbols(in)

	reelIndex := 0
	for i, dim := range mask {
		for j := 0; j < dim; j++ {
			if reelIndex < len(mappedReels) {
				newScreen[i][j] = int(mappedReels[reelIndex])
				reelIndex++
			} else {
				newScreen[i][j] = 9
			}
		}
	}

	return newScreen
}
