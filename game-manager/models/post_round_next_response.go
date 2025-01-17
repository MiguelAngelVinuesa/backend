package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func UnmarshallRoundNextResponse(resp *http.Response) (*slots.RoundResult, int64, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	r := roundNextResponsePool.Acquire().(*roundNextResponse)
	defer r.Release()

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || !r.success {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("empty response or success==false")
		}
		return nil, 0, err
	}

	return r.result.Clone().(*slots.RoundResult), r.newBalance, err
}

type roundNextResponse struct {
	success    bool
	newBalance int64
	result     *slots.RoundResult
	roundID    string
	pool.Object
}

func (r *roundNextResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool

	if string(key) == "success" {
		if r.success, ok = dec.Bool(); ok {
			return nil
		}
	} else if string(key) == "newBalance" {
		if r.newBalance, ok = dec.Int64(); ok {
			return nil
		}
	} else if string(key) == "result" {
		if b, _, ok2 := dec.String(); ok2 {
			r.result = slots.AcquireRoundResultFromJSON(dec.Unescaped(b))
			return nil
		}
	} else if string(key) == "roundId" {
		if b, escaped, ok2 := dec.String(); ok2 {
			if escaped {
				r.roundID = string(dec.Unescaped(b))
			} else {
				r.roundID = string(b)
			}
			return nil
		}
	} else {
		return nil // ignore unknown fields
	}

	return dec.Error()
}

var roundNextResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &roundNextResponse{}
	return r, r.reset
})

func (r *roundNextResponse) reset() {
	if r.result != nil {
		r.result.Release()
		r.result = nil
	}

	r.success = false
	r.newBalance = 0
	r.roundID = ""
}
