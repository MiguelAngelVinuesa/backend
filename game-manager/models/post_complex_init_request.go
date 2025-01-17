package models

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func MarshallComplexInitRequest(r *slots.Round, debug, withData bool) (*zjson.Encoder, error) {
	enc := zjson.AcquireEncoder(4096)
	enc.StartObject()
	enc.StringField("sessionId", r.SessionID())
	enc.BoolFieldOpt("debug", debug)

	if withData {
		enc2 := zjson.AcquireEncoder(4096)
		enc2.StartArray()
		res := r.RoundResults()
		for ix := range res {
			enc2.Object(res[ix])
		}
		enc2.EndArray()
		enc.EscapedBytesStringField("result", enc2.Bytes())

		if s := r.GameState(); s != nil {
			enc2.Reset()
			enc2.Object(s)
			enc.EscapedBytesStringField("sessionState", enc2.Bytes())
		}

		if s := slots.AcquireRoundState(r.SessionID(), r.RoundID(), len(r.Results())); s != nil {
			enc2.Reset()
			enc2.Object(s)
			enc.EscapedBytesStringField("roundState", enc2.Bytes())
		}

		enc2.Release()
	}

	enc.EndObject()
	return enc, nil
}
