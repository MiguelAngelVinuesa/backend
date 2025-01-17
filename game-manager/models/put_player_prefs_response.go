package models

import (
	"fmt"
	"io"
	"net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

func UnmarshallPutPlayerPrefsResponse(resp *http.Response) (bool, error) {
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	r := putPlayerPrefsResponsePool.Acquire().(*putPlayerPrefsResponse)
	defer r.Release()

	dec := zjson.AcquireDecoder(b)
	defer dec.Release()

	if !dec.Object(r) || !r.success {
		if err = dec.Error(); err == nil {
			err = fmt.Errorf("invalid response or success==false")
		}
		return false, err
	}

	return true, nil
}

type putPlayerPrefsResponse struct {
	success bool
	pool.Object
}

var putPlayerPrefsResponsePool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &putPlayerPrefsResponse{}
	return r, r.reset
})

func (r *putPlayerPrefsResponse) reset() {
	r.success = false
}

func (r *putPlayerPrefsResponse) DecodeField(dec *zjson.Decoder, key []byte) error {
	if string(key) == "success" {
		if success, ok := dec.Bool(); ok {
			r.success = success
			return nil
		}
		return dec.Error()
	}
	return nil // ignore unknown fields
}
