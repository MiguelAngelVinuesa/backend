package events

import (
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net/kafka"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/i18n"
)

// Consumer is the global interface for consuming events.
var Consumer = &consumer{}

// Consume implements the kafka.Consumer interface.
func (c *consumer) Consume(_, msg []byte, created time.Time) error {
	e, err := kafka.NewEventFromJSON(msg, make(map[string]any, 16))
	if err != nil {
		return err
	}

	if e.Operation == kafka.OpCreate || e.Operation == kafka.OpUpdate {
		if data, ok := e.Data.(map[string]any); ok {
			switch e.EntityCode {
			case consts.EvLocaleStr:
				go i18n.ProcessEvent(data)

			case consts.EvMessage:
				if session, m := NewMessageFromEvent(e, created); m != nil && session != "" {
					go AddMessage(session, m)
				}
			}
		}
	}

	return nil
}

type consumer struct{}
