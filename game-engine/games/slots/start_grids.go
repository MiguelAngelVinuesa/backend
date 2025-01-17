package slots

import (
	"encoding/binary"
	"io"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// StartGrids is a spin grid filler that selects a random grid from a set of pre-calculated grids.
type StartGrids struct {
	f               io.ReadSeekCloser
	prng            interfaces.Generator
	startCount      int
	secondaryCount  int
	transformsCount int
	doubleSpin      bool
	lastGrid        int
}

// NewStartGrids instantiates a new start grid filler from the memory pool.
func NewStartGrids(f io.ReadSeekCloser, doubleSpin bool) *StartGrids {
	s := startGridsPool.Get().(*StartGrids)
	s.f = f
	s.prng = rng.AcquireRNG()
	s.doubleSpin = doubleSpin

	// initialize the counts from the file.
	if _, err := s.f.Seek(0, io.SeekStart); err == nil {
		bb := make([]byte, 8)
		if _, err = s.f.Read(bb); err == nil {
			s.startCount = int(binary.LittleEndian.Uint64(bb))
			s.transformsCount = 20 * s.startCount
			if _, err = s.f.Read(bb); err == nil {
				s.secondaryCount = int(binary.LittleEndian.Uint64(bb))
			}
		}
	}

	return s
}

// Release returns the start grid filler to the memory pool.
func (s *StartGrids) Release() {
	if s != nil {
		s.prng.ReturnToPool()
		s.prng = nil
		s.f.Close()
		s.f = nil
		startGridsPool.Put(s)
	}
}

// LastGrid returns the last selected start grid and transformation.
// result / 20 == start grid.
// result % 20 == transformation.
func (s *StartGrids) LastGrid() int {
	return s.lastGrid
}

// SetLastGrid sets the last selected start grid and transformation.
func (s *StartGrids) SetLastGrid(grid int) {
	s.lastGrid = grid
}

// Spin implements the Spinner interface to fill a spin grid.
func (s *StartGrids) Spin(spin *slots.Spin, indexes utils.Indexes) {
	if s.startCount > 0 {
		reels, rows := spin.GridSize()
		l := len(indexes)

		if !s.doubleSpin || !spin.HasSticky() {
			// first spin or regular spin.
			s.lastGrid = s.prng.IntN(s.transformsCount)
			n := s.lastGrid / 20
			offset := 16 + int64(n*(l*(s.secondaryCount+1)))
			if _, err := s.f.Seek(offset, io.SeekStart); err == nil {
				key := make([]byte, l)
				if _, err = s.f.Read(key); err == nil {
					for ix := range indexes {
						indexes[ix] = utils.Index(key[ix])
					}
					transform(s.lastGrid%20, reels, rows, indexes)
					return
				}
			}
		} else if s.secondaryCount > 0 {
			// second spin
			n := s.lastGrid / 20
			n2 := 1 + s.prng.IntN(s.secondaryCount)
			offset := 16 + int64(n*l*(s.secondaryCount+1)+n2*l)
			if _, err := s.f.Seek(offset, io.SeekStart); err == nil {
				key := make([]byte, l)
				if _, err = s.f.Read(key); err == nil {
					for ix := range indexes {
						indexes[ix] = utils.Index(key[ix])
					}
					transform(s.lastGrid%20, reels, rows, indexes)
					return
				}
			}
		}
	}

	// failed to read from the file, so reset the spin to use its built-in function in the future!
	spin.SetSpinner(nil)
	// and use it now also :)
	spin.Builtin()
}

func transform(kind, reels, rows int, indexes utils.Indexes) {
	slots.Shift(reels, rows, kind%5, indexes)

	switch kind / 5 {
	case 1:
		slots.FlipHorizontal(reels, rows, indexes)
	case 2:
		slots.FlipVertical(reels, rows, indexes)
	case 3:
		slots.Rotate(reels, rows, indexes)
	}
}

var startGridsPool = sync.Pool{New: func() interface{} { return &StartGrids{} }}
