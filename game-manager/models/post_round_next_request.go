package models

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func MarshallRoundNextRequest(sessionID, roundID string, roundState *slots.RoundState, spinSeq int) (*zjson.Encoder, error) {
	enc := zjson.AcquireEncoder(256)
	enc.StartObject()
	enc.StringField("sessionId", sessionID)
	enc.StringField("roundId", roundID)
	enc.IntFieldOpt("spin", spinSeq)

	if roundState != nil {
		enc2 := zjson.AcquireEncoder(4096)
		enc2.Object(roundState)
		enc.EscapedBytesStringField("roundState", enc2.Bytes())
		enc2.Release()
	}

	enc.EndObject()
	return enc, nil
}
