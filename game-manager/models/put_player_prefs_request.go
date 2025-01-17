package models

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func MarshallPlayerPrefsRequest(sessionID string, prefs map[string]string) (*zjson.Encoder, error) {
	enc := zjson.AcquireEncoder(1024)
	enc.StartObject()
	enc.StringField("sessionId", sessionID)

	enc2 := zjson.AcquireEncoder(1024)
	enc2.StartObject()
	for k, v := range prefs {
		enc2.StringField(k, v)
	}
	enc2.EndObject()
	enc.EscapedBytesStringField("state", enc2.Bytes())
	enc2.Release()

	enc.EndObject()
	return enc, nil
}
