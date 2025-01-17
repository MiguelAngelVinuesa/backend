package encode

import (
	"io"

	"github.com/goccy/go-json"
)

// Goccy is used to optimize JSON (un)marshalling using goccy/go-json.
type Goccy struct{}

// Consume implements the Consumer interface.
func (c *Goccy) Consume(r io.Reader, i interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(i)
}

// Produce implements the Producer interface.
func (c *Goccy) Produce(w io.Writer, i interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(i)
}
