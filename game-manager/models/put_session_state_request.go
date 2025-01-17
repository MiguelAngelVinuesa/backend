package models

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

func MarshallSessionStateRequest(sessionID string, state *slots.GameState) (*zjson.Encoder, error) {
	enc := zjson.AcquireEncoder(1024)
	enc.StartObject()
	enc.StringField("sessionId", sessionID)

	enc2 := zjson.AcquireEncoder(1024)
	enc2.Object(state)
	enc.EscapedBytesStringField("state", enc2.Bytes())
	enc2.Release()

	enc.EndObject()
	return enc, nil
}
