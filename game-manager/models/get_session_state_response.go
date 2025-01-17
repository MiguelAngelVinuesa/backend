package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func UnmarshallGetSessionStateResponse(resp *http.Response) (*slots.GameState, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := sessionStateResponsePool.Acquire().(*sessionStateResponse)
	defer r.Release()

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || r.state == nil {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("empty response")
		}
		return nil, err
	}

	return r.state.Clone().(*slots.GameState), nil
}

func EmptySessionState() *slots.GameState {
	s, _ := slots.AcquireGameStateFromJSON([]byte{'{', '}'})
	return s
}

type sessionStateResponse struct {
	state *slots.GameState
	pool.Object
}

var sessionStateResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &sessionStateResponse{}
	return r, r.reset
})

func (r *sessionStateResponse) reset() {
	if r.state != nil {
		r.state.Release()
		r.state = nil
	}
}

func (r *sessionStateResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	if string(key) == "state" {
		if state, _, ok := dec.String(); ok {
			var err error
			r.state, err = slots.AcquireGameStateFromJSON(dec.Unescaped(state))
			return err
		}
		return dec.Error()
	}
	return nil // ignore unknown fields
}
