package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	MaxReels     = 12
	MaxRows      = 10
	MaxTiles     = 100
	MaxNeighbors = 8
)

// GridDirection indicates the direction between neighboring grid tiles.
type GridDirection uint8

const (
	GridLeft GridDirection = iota + 1
	GridRight
	GridUp
	GridDown
	GridLeftUp
	GridLeftDown
	GridRightUp
	GridRightDown
	GridHorizontal // combines both Left + Right.
	GridVertical   // combines both Up + Down.
	GridDiagonal   // combines all 4 diagonal directions.
	GridAny        // combines all directions.
)

// String implements the Stringer() interface.
func (d GridDirection) String() string {
	switch d {
	case GridLeft:
		return "left"
	case GridRight:
		return "right"
	case GridUp:
		return "up"
	case GridDown:
		return "down"
	case GridLeftUp:
		return "left+up"
	case GridLeftDown:
		return "left+down"
	case GridRightUp:
		return "right+up"
	case GridRightDown:
		return "right+down"
	case GridHorizontal:
		return "horizontal"
	case GridVertical:
		return "vertical"
	case GridDiagonal:
		return "diagonal"
	case GridAny:
		return "any"
	default:
		return "[invalid]"
	}
}

// IsHorizontal returns true if the direction is horizontal.
func (d GridDirection) IsHorizontal() bool {
	return d == GridLeft || d == GridRight || d == GridHorizontal
}

// IsVertical returns true if the direction is vertical.
func (d GridDirection) IsVertical() bool {
	return d == GridUp || d == GridDown || d == GridVertical
}

// IsDiagonal returns true if the direction is diagonal.
func (d GridDirection) IsDiagonal() bool {
	return d == GridLeftDown || d == GridLeftUp || d == GridRightDown || d == GridRightUp || d == GridDiagonal
}

// GridNeighbors is a convenience type for a map of neighboring tiles with their directions.
type GridNeighbors map[uint8]GridDirection

// RandomNeighbor retrieves a random tile from the list of neighbors, and returns its offset.
// The center tile is excluded from the result.
func (n GridNeighbors) RandomNeighbor(prng interfaces.Generator) uint8 {
	ix := prng.IntN(len(n))
	for k := range n {
		if ix--; ix < 0 {
			return k
		}
	}
	return 0 // can never happen!
}

// GridDefinition contains the details of a grid.
type GridDefinition struct {
	haveMask     bool            // indicates if a grid mask was supplied during instantiation.
	reels        int             // number of reels in the grid.
	rows         int             // number of rows in the grid.
	gridSize     int             // number of tiles in the grid.
	offsetStep   int             // prime number for stepping through offsets, guranteeing all tiles will be touched.
	mask         utils.UInt8s    // the grid mask.
	stepsOffGrid utils.UInt8s    // number of steps for each tile to jump off grid.
	neighbors    []GridNeighbors // slice of neighboring tiles for each tile.
	withoutSelf  []GridNeighbors // slice of neighboring tiles for each tile excluding the tile itself.
}

// NewGridDefinition instantiates a new grid definition.
// This function is "expensive". It should only be called once for each possible grid during app initialization.
// The function panics if invalid parameters are given.
func NewGridDefinition(reels, rows int, mask utils.UInt8s) *GridDefinition {
	if reels <= 0 || rows <= 0 || reels > MaxReels || rows > MaxRows || reels*rows > MaxTiles {
		panic(consts.MsgInvalidGridSize)
	}
	if l := len(mask); l != 0 && l != reels {
		panic(consts.MsgInvalidGridMask)
	}

	for reel := range mask {
		m := int(mask[reel])
		if m == 0 || m > rows {
			panic(consts.MsgInvalidGridMask)
		}
	}

	gridSize := reels * rows
	g := &GridDefinition{
		haveMask:     len(mask) > 0,
		reels:        reels,
		rows:         rows,
		gridSize:     gridSize,
		offsetStep:   7,
		mask:         mask,
		neighbors:    make([]GridNeighbors, gridSize),
		withoutSelf:  make([]GridNeighbors, gridSize),
		stepsOffGrid: make(utils.UInt8s, gridSize),
	}

	half := gridSize / 2
	for ix := range primes {
		if p := primes[ix]; p > half {
			g.offsetStep = p
			break
		}
	}

	if !g.haveMask {
		// make sure there's always a mask.
		for reel := 0; reel < g.reels; reel++ {
			g.mask = append(g.mask, uint8(rows))
		}
	}

	// initialize neighboring tiles and calculate steps off grid.
	for ix := 0; ix < g.gridSize; ix++ {
		g.neighbors[ix] = make(GridNeighbors, MaxNeighbors)
		g.withoutSelf[ix] = make(GridNeighbors, MaxNeighbors)
	}
	return g.init()
}

// ReelCount returns the number of reels in the grid.
func (g *GridDefinition) ReelCount() uint8 {
	return uint8(g.reels)
}

// RowCount returns the number of rows in the grid.
func (g *GridDefinition) RowCount() uint8 {
	return uint8(g.rows)
}

// HaveMask returns true if the grid definition was instantiated with a mask.
func (g *GridDefinition) HaveMask() bool {
	return g.haveMask
}

// GridMask returns the grid mask.
// Unless specified during instantiation, it defaults to a slice of size ReelCount() filled with RowCount().
func (g *GridDefinition) GridMask() utils.UInt8s {
	return g.mask
}

// Neighbors returns a map with neighboring tiles for a given tile.
func (g *GridDefinition) Neighbors(offset uint8) GridNeighbors {
	return g.neighbors[offset]
}

// IsOnTheEdge returns true if the given offset is on the edge of the grid.
func (g *GridDefinition) IsOnTheEdge(offset int) bool {
	reel, row := g.ReelRowFromOffset(offset)
	return reel == 0 || row == 0 || reel >= g.reels-1 || row >= int(g.mask[reel])-1
}

// IsValidOffset returns true if the given offset is a valid tile on the grid.
func (g *GridDefinition) IsValidOffset(offset int) bool {
	reel, row := g.ReelRowFromOffset(offset)
	return reel < g.reels && row < int(g.mask[reel])
}

// TilesFromEdge returns the number of tiles which are the given number of steps away from the edge.
func (g *GridDefinition) TilesFromEdge(steps uint8) uint8 {
	var count uint8
	for ix := range g.stepsOffGrid {
		if g.stepsOffGrid[ix] == steps {
			count++
		}
	}
	return count
}

// NextOffset finds the next valid offset from the current one.
// It uses a prime number to step through the tiles, and is guaranteed to touch all tiles eventually.
func (g *GridDefinition) NextOffset(offset uint8) uint8 {
	var offs int
	for {
		offs = int(offset) + g.offsetStep
		if offs > g.gridSize {
			offs -= g.gridSize
		}
		if g.IsValidOffset(offs) {
			break
		}
	}
	return uint8(offs)
}

func (g *GridDefinition) ReelRowFromOffset(offset int) (int, int) {
	reel := offset / g.rows
	row := offset - (reel * g.rows)
	return reel, row
}

// IsNeighbor determines if two offsets in the grid are neighboring tiles.
// If specified, it will take the mask into consideration for determining which tiles are neighbors.
// If the tiles are found to be neighbors, the functiuon returns true and reports the direction from the first to the second.
// As a special case, if the offsets are equal, the function returns true with direction 0.
// In all other cases the function returns false, with direction 0.
func (g *GridDefinition) IsNeighbor(offset1, offset2 uint8) (bool, GridDirection) {
	dir, ok := g.neighbors[offset1][offset2]
	return ok, dir
}

// IsEmpty implements the zjson.Encoder interface.
func (g *GridDefinition) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson.Encoder interface.
func (g *GridDefinition) EncodeFields(enc *zjson.Encoder) {
	enc.IntField("reels", g.reels)
	enc.IntField("rows", g.rows)
	if g.haveMask {
		enc.StartArrayField("mask")
		for ix := range g.mask {
			enc.Int64(int64(g.mask[ix]))
		}
		enc.EndArray()
	}
}

// init builds up the slice of maps with neighboring tiles.
func (g *GridDefinition) init() *GridDefinition {
	reels, rows := g.reels, g.rows
	mask := g.mask

	testRegular := func(offset1, reel2 int) {
		offset2 := reel2 * rows
		max2 := int(mask[reel2])
		for row2 := 0; row2 < max2; row2++ {
			if ok, dir := g.isNeighbor(offset1, offset2); ok {
				g.neighbors[offset1][uint8(offset2)] = dir
				if dir != 0 {
					g.withoutSelf[offset1][uint8(offset2)] = dir
				}
			}
			offset2++
		}
	}

	testMasked := func(offset1, reel2 int) {
		offset2 := reel2 * rows
		max2 := int(mask[reel2])
		for row2 := 0; row2 < max2; row2++ {
			if ok, dir := g.isNeighborMasked(offset1, offset2); ok {
				g.neighbors[offset1][uint8(offset2)] = dir
				if dir != 0 {
					g.withoutSelf[offset1][uint8(offset2)] = dir
				}
			}
			offset2++
		}
	}

	test := testRegular
	if g.haveMask {
		test = testMasked
	}

	for reel1 := 0; reel1 < reels; reel1++ {
		offset1 := reel1 * rows
		max1 := int(mask[reel1])
		for row1 := 0; row1 < max1; row1++ {
			if reel1 > 0 {
				test(offset1, reel1-1)
			}
			test(offset1, reel1)
			if reel1+1 < reels {
				test(offset1, reel1+1)
			}
			offset1++
		}
	}

	for offset := 0; offset < g.gridSize; offset++ {
		g.stepsOffGrid[offset] = g.getStepsOffGrid(offset)
	}

	return g
}

func (g *GridDefinition) isNeighbor(offset1, offset2 int) (bool, GridDirection) {
	reel1, row1 := g.ReelRowFromOffset(offset1)
	reel2, row2 := g.ReelRowFromOffset(offset2)
	reelDiff, rowDiff := reel1-reel2, row1-row2

	switch reelDiff {
	case -1:
		switch rowDiff {
		case -1:
			return true, GridRightDown
		case 0:
			return true, GridRight
		case 1:
			return true, GridRightUp
		}

	case 0:
		switch rowDiff {
		case -1:
			return true, GridDown
		case 0:
			return true, 0
		case 1:
			return true, GridUp
		}

	case 1:
		switch rowDiff {
		case -1:
			return true, GridLeftDown
		case 0:
			return true, GridLeft
		case 1:
			return true, GridLeftUp
		}
	}

	return false, 0
}

func (g *GridDefinition) isNeighborMasked(offset1, offset2 int) (bool, GridDirection) {
	reel1, row1 := g.ReelRowFromOffset(offset1)
	reel2, row2 := g.ReelRowFromOffset(offset2)
	reelDiff, rowDiff := reel1-reel2, row1-row2
	sizeDiff := int(g.mask[reel1]) - int(g.mask[reel2])

	switch reelDiff {
	case -1:
		switch sizeDiff {
		case -2:
			switch rowDiff {
			case -2:
				return true, GridRightDown
			case -1:
				return true, GridRight
			case 0:
				return true, GridRightUp
			}

		case -1:
			switch rowDiff {
			case -1:
				return true, GridRightDown
			case 0:
				return true, GridRightUp
			}

		case 0:
			switch rowDiff {
			case -1:
				return true, GridRightDown
			case 0:
				return true, GridRight
			case 1:
				return true, GridRightUp
			}

		case 1:
			switch rowDiff {
			case 0:
				return true, GridRightDown
			case 1:
				return true, GridRightUp
			}

		case 2:
			switch rowDiff {
			case 0:
				return true, GridRightDown
			case 1:
				return true, GridRight
			case 2:
				return true, GridRightUp
			}
		}

	case 0:
		switch rowDiff {
		case -1:
			return true, GridDown
		case 0:
			return true, 0
		case 1:
			return true, GridUp
		}

	case 1:
		switch sizeDiff {
		case -2:
			switch rowDiff {
			case -2:
				return true, GridLeftDown
			case -1:
				return true, GridLeft
			case 0:
				return true, GridLeftUp
			}

		case -1:
			switch rowDiff {
			case -1:
				return true, GridLeftDown
			case 0:
				return true, GridLeftUp
			}

		case 0:
			switch rowDiff {
			case -1:
				return true, GridLeftDown
			case 0:
				return true, GridLeft
			case 1:
				return true, GridLeftUp
			}

		case 1:
			switch rowDiff {
			case 0:
				return true, GridLeftDown
			case 1:
				return true, GridLeftUp
			}

		case 2:
			switch rowDiff {
			case 0:
				return true, GridLeftDown
			case 1:
				return true, GridLeft
			case 2:
				return true, GridLeftUp
			}

		}
	}

	return false, 0
}

func (g *GridDefinition) getStepsOffGrid(offset int) uint8 {
	reel, row := g.ReelRowFromOffset(offset)

	if row >= int(g.mask[reel]) {
		return 0
	}
	if g.IsOnTheEdge(offset) {
		return 1
	}

	out := row + 1
	if down := int(g.mask[reel]) - row; down < out {
		out = down
	}
	if left := reel + 1; left < out {
		out = left
	}
	if right := g.reels - reel; right < out {
		out = right
	}
	return uint8(out)
}

var (
	primes = []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53}
)
