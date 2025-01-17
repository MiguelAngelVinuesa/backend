package slots

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

// NewFilteredSpinner instantiates a new set of filterable spinners.
// This function panics if the given slices are of unequal length.
func NewFilteredSpinner(filters []SpinDataFilterer, sets []Spinner) *FilteredSpinner {
	if len(filters) != len(sets) {
		panic("FilteredSpinner: invalid configuration")
	}

	return &FilteredSpinner{
		filters:  filters,
		spinners: sets,
	}
}

// Add adds a new filter+spinner to the set of filterable spinners.
func (fsr *FilteredSpinner) Add(filter SpinDataFilterer, spinner Spinner) *FilteredSpinner {
	fsr.filters = append(fsr.filters, filter)
	fsr.spinners = append(fsr.spinners, spinner)
	return fsr
}

// Filters returns the configured filters.
func (fsr *FilteredSpinner) Filters() []SpinDataFilterer {
	return fsr.filters
}

// Spinners returns the configured spinners.
func (fsr *FilteredSpinner) Spinners() []Spinner {
	return fsr.spinners
}

// Spin can be used to fill all reels on the grid with random symbols from the configured symbol reels.
// The first spin state filter, that matches the spin state, will run the corresponding spinner.
// This function panics if none of the filters match the current spin state.
func (fsr *FilteredSpinner) Spin(spin *Spin, out utils.Indexes) {
	for ix := range fsr.filters {
		f := fsr.filters[ix]
		if f(spin) {
			fsr.spinners[ix].Spin(spin, out)
			return
		}
	}
	panic("FilteredSpinner: no matching filter found")
}

// FilteredSpinner represents a set of filterable spinners.
type FilteredSpinner struct {
	filters  []SpinDataFilterer
	spinners []Spinner
}
