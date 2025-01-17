package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// FlipHorizontal flips the grid on the horizontal plane.
// It panics if reels * rows != len(grid).
func FlipHorizontal(reels, rows int, grid utils.Indexes) {
	if reels*rows != len(grid) {
		panic(consts.MsgInvalidGridSize)
	}

	for reel := 0; reel < reels; reel++ {
		o1 := reel * rows
		o2 := o1 + rows
		for o2 > o1 {
			o2--
			grid[o1], grid[o2] = grid[o2], grid[o1]
			o1++
		}
	}
}

// FlipVertical flips the grid on the vertical plane.
// It panics if reels * rows != len(grid).
func FlipVertical(reels, rows int, grid utils.Indexes) {
	if reels*rows != len(grid) {
		panic(consts.MsgInvalidGridSize)
	}

	r1 := 0
	r2 := reels
	for r2 > r1 {
		r2--
		o1 := r1 * rows
		o2 := r2 * rows
		for row := 0; row < rows; row++ {
			grid[o1], grid[o2] = grid[o2], grid[o1]
			o1++
			o2++
		}
		r1++
	}
}

// Rotate rotates the grid 180 degrees.
// It panics if reels * rows != len(grid).
func Rotate(reels, rows int, grid utils.Indexes) {
	FlipVertical(reels, rows, grid)
	FlipHorizontal(reels, rows, grid)
}

// Shift shifts the reels count times to the right.
// It panics if reels * rows != len(grid).
func Shift(reels, rows, count int, grid utils.Indexes) {
	if reels*rows != len(grid) {
		panic(consts.MsgInvalidGridSize)
	}

	for ix := 0; ix < count; ix++ {
		o1 := reels * rows
		o2 := o1 - rows
		for reel := reels - 1; reel > 0; reel-- {
			for row := 0; row < rows; row++ {
				o1--
				o2--
				grid[o1], grid[o2] = grid[o2], grid[o1]
			}
		}
	}
}
