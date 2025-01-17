package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func UnmarshallGetRoundStateResponse(sessionID, roundID string, resp *http.Response) (*slots.RoundState, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	r := roundStateResponsePool.Acquire().(*roundStateResponse)
	defer r.Release()

	r.sessionID = sessionID
	r.roundID = roundID

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || r.state == nil {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("empty response")
		}
		return nil, err
	}

	return r.state.Clone().(*slots.RoundState), nil
}

type roundStateResponse struct {
	state     *slots.RoundState
	sessionID string
	roundID   string
	pool.Object
}

var roundStateResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &roundStateResponse{}
	return r, r.reset
})

func (r *roundStateResponse) reset() {
	if r.state != nil {
		r.state.Release()
		r.state = nil
	}
	r.sessionID = ""
	r.roundID = ""
}

func (r *roundStateResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	if string(key) == "state" {
		if state, _, ok := dec.String(); ok {
			var err error
			r.state, err = slots.AcquireRoundStateFromJSON(r.sessionID, r.roundID, dec.Unescaped(state))
			return err
		}
		return dec.Error()
	}
	return nil // ignore unknown fields
}
