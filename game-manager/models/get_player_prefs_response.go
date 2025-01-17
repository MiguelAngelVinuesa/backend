package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
	"github.com/goccy/go-json"
)

func UnmarshallGetPlayerPrefsResponse(resp *http.Response) (map[string]string, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := playerPrefsResponsePool.Acquire().(*playerPrefsResponse)
	defer r.Release()

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || r.state == nil {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("empty response")
		}
		return nil, err
	}

	return r.state, nil
}

type playerPrefsResponse struct {
	state map[string]string
	pool.Object
}

var playerPrefsResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &playerPrefsResponse{
		state: make(map[string]string, 8),
	}
	return r, nil
})

func (r *playerPrefsResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	if string(key) == "state" {
		if state, _, ok := dec.String(); ok {
			return json.Unmarshal(dec.Unescaped(state), &r.state)
		}
		return dec.Error()
	}
	return nil // ignore unknown fields
}
