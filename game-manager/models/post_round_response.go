package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

func UnmarshallRoundResponse(resp *http.Response) (string, int64, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	r := roundResponsePool.Acquire().(*roundResponse)
	defer r.Release()

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || !r.success {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("invalid response or success==false")
		}
		return "", 0, err
	}

	return r.roundID, r.playerData.balance, nil
}

type roundResponse struct {
	success    bool
	roundID    string
	playerData playerData
	pool.Object
}

var roundResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &roundResponse{}
	return r, r.reset
})

func (r *roundResponse) reset() {
	r.success = false
	r.roundID = ""
	r.playerData.balance = 0
}

func (r *roundResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool

	if string(key) == "playerData" {
		if ok = dec.Object(&r.playerData); ok {
			return nil
		}
	} else if string(key) == "roundId" {
		if b, escaped, ok2 := dec.String(); ok2 {
			if escaped {
				r.roundID = string(dec.Unescaped(b))
			} else {
				r.roundID = string(b)
			}
		}
	} else if string(key) == "success" {
		if r.success, ok = dec.Bool(); ok {
			return nil
		}
	} else {
		return nil // ignore unknown fields
	}

	return dec.Error()
}

type playerData struct {
	balance int64
}

func (p *playerData) DecodeField(dec *zjson.Decoder, key []byte) error {
	if string(key) == "balance" {
		if i, ok := dec.Int64(); ok {
			p.balance = i
			return nil
		}
		return dec.Error()
	}
	return nil // ignore unknown fields
}
