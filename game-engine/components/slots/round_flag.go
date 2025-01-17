package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

// RoundFlag contains the details for a round flag.
type RoundFlag struct {
	id     int
	name   string
	export bool
}

// NewRoundFlag instantiates a new round flag details.
func NewRoundFlag(id int, name string) *RoundFlag {
	return &RoundFlag{id: id, name: name}
}

// WithExport makes the flag exportable in results.
func (f *RoundFlag) WithExport() *RoundFlag {
	f.export = true
	return f
}

// ID returns the identifier of the flag.
func (f *RoundFlag) ID() int {
	return f.id
}

// Name returns the name of the flag.
func (f *RoundFlag) Name() string {
	return f.name
}

// Export returns true if the flag should be exported in results.
func (f *RoundFlag) Export() bool {
	return f.export
}

// IsEmpty implements the zjson.Encoder interface.
func (f *RoundFlag) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson.Encoder interface.
func (f *RoundFlag) EncodeFields(enc *zjson.Encoder) {
	enc.IntField("id", f.id)
	enc.StringField("name", f.name)
}

// RoundFlags is a convenience type for a slice of round flags.
type RoundFlags []*RoundFlag

// GetMaxID returns the highest round flag ID.
func (f RoundFlags) GetMaxID() int {
	var max int
	for ix := range f {
		if id := f[ix].id; id > max {
			max = id
		}
	}
	return max
}
