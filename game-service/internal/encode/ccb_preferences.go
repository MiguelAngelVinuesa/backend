package encode

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func CcbPreferences(state *slots.SymbolsState) *zjson.Encoder {
	enc := zjson.AcquireEncoder(1024)
	enc.StartObject()
	enc.BoolField("success", true)
	enc.ObjectField("state", state)
	enc.EndObject()
	return enc
}
