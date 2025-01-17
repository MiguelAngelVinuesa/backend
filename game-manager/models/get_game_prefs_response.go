package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func UnmarshallGetGamePrefsResponse(resp *http.Response) (string, string, *slots.GamePrefs, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", nil, err
	}

	r := gamePrefsResponsePool.Acquire().(*gamePrefsResponse)
	defer r.Release()

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || r.state == nil {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("empty response")
		}
		return "", "", nil, err
	}

	return r.casinoID, r.playerID, r.state.Clone().(*slots.GamePrefs), nil
}

func EmptyGamePrefs() *slots.GamePrefs {
	s, _ := slots.AcquireGamePrefsFromJSON([]byte{'{', '}'})
	return s
}

type gamePrefsResponse struct {
	state    *slots.GamePrefs
	casinoID string
	playerID string
	pool.Object
}

var gamePrefsResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &gamePrefsResponse{}
	return r, r.reset
})

func (r *gamePrefsResponse) reset() {
	if r.state != nil {
		r.state.Release()
		r.state = nil
		r.casinoID = ""
		r.playerID = ""
	}
}

func (r *gamePrefsResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	var s []byte
	var escaped, ok bool
	var err error

	if string(key) == "state" {
		if s, escaped, ok = dec.String(); ok {
			if escaped {
				r.state, err = slots.AcquireGamePrefsFromJSON(dec.Unescaped(s))
			} else {
				r.state, err = slots.AcquireGamePrefsFromJSON(s)
			}
			if err != nil {
				return err
			}
		}
	} else if string(key) == "casinoId" {
		if s, escaped, ok = dec.String(); ok {
			if escaped {
				r.casinoID = string(dec.Unescaped(s))
			} else {
				r.casinoID = string(s)
			}
		}
	} else if string(key) == "playerId" {
		if s, escaped, ok = dec.String(); ok {
			if escaped {
				r.playerID = string(dec.Unescaped(s))
			} else {
				r.playerID = string(s)
			}
		}
	}

	if !ok {
		return dec.Error()
	}
	return nil
}
