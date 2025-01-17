package hashes

import (
	"context"
	"fmt"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net/kafka"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
)

func SendReport(appStart time.Time) {
	done, cancel := context.WithCancel(context.Background())
	defer cancel()

	p, err := kafka.NewProducer(config.MqBrokers, config.ClientID, config.HashesTopic, log.Logger, done)
	if err != nil {
		return
	}

	ts := appStart.UnixMilli()
	key := fmt.Sprintf("%s:%d", config.ClientID, ts)

	data := map[string]any{
		"svcCode":     MainFile,
		"main":        MainHash,
		"rngLib":      RngLibHash,
		"gameEngine":  GameEngine,
		"gameConfig":  GameConfig,
		"gameManager": GameManager,
		"gameService": GameService,
		"debugMode":   config.DebugMode,
	}
	for k, v := range GameHashes {
		data[k] = v
	}

	e := kafka.NewEvent(config.ClientID, "gamesvc", key, kafka.OpCreate, data)
	ctx, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	e.Produce(ctx, p)
	cancel2()
}
