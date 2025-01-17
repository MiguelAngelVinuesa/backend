package encode

import (
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/hashes"
)

func Hashes() *zjson.Encoder {
	enc := zjson.AcquireEncoder(1024)
	enc.StartObject()

	enc.StartObjectField("binaries")
	enc.StringField("main", hashes.MainHash)
	enc.StringField("mainMod", hashes.MainModDate.Format(time.RFC3339))
	enc.IntField("mainSize", int(hashes.MainFileSize))
	enc.StringField("rngLib", hashes.RngLibHash)
	enc.StringField("rngLibMod", hashes.RngLibModDate.Format(time.RFC3339))
	enc.IntField("rngLibSize", int(hashes.RngLibFileSize))
	enc.StringField("rngInclude", hashes.RngIncludeHash)
	enc.EndObject()

	enc.StartObjectField("git")
	enc.StringField("gameEngine", hashes.GameEngine)
	enc.StringField("gameConfig", hashes.GameConfig)
	enc.StringField("gameManager", hashes.GameManager)
	enc.StringField("gameService", hashes.GameService)
	enc.EndObject()

	enc.StartObjectField("games")
	for k, v := range hashes.GameHashes {
		enc.StringField(k, v)
	}
	enc.EndObject()

	enc.BoolField("debugMode", config.DebugMode)
	enc.BoolField("success", true)

	enc.EndObject()
	return enc
}
