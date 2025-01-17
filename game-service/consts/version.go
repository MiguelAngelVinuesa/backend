package consts

import (
	"fmt"
)

// Semantic versioning scheme: 1.0.0-rc.259
const (
	Major = 1 // major changes.
	Minor = 0 // minor changes.
	Patch = 0 // patch changes.
)

var (
	// SemTag is used to define a "pre-release" tag.
	SemTag = "-precert" // e.g. "-alpha1", "-beta2", "-rc3"
	// SemVerShort contains the short string representation.
	SemVerShort = fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
	// SemVerFull contains the full string representation.
	SemVerFull = fmt.Sprintf("%d.%d.%d%s.%s", Major, Minor, Patch, SemTag, SemBuild)
	// SemDate contains the date of the build (will be compiled in).
	SemDate string
	// SemBuild contains the build number (will be complied in).
	SemBuild string
)
