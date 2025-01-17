package slots

// Offsets contains the reel & row offset for a single expansion position.
// First offset is for the reel, second offset for the row.
type Offsets [2]int

// GridOffsets contains all the reel & row offsets for an expansion grid.
type GridOffsets []Offsets

// DefaultGridOffsets is the default 3x3 expansion grid around the triggering position.
var DefaultGridOffsets = GridOffsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
