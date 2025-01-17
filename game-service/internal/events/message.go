package events

import (
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net/kafka"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/i18n"
)

func NewMessage(id tg.MessageKind, mode tg.DisplayMode, created time.Time, expire int) *Message {
	return &Message{
		messageID:   id,
		displayMode: mode,
		created:     created,
		expires:     created.Add(time.Second * time.Duration(expire)),
	}
}

func NewMessageFromEvent(event *kafka.Event, created time.Time) (string, *Message) {
	if data, ok := event.Data.(map[string]any); ok {
		if env := conv.StringFromAny(data["environmentID"]); env == config.Environment {
			session := conv.StringFromAny(data["sessionID"])
			ttl := time.Duration(conv.IntFromAny(data["sessionTTL"]))

			return session, &Message{
				messageID:   tg.MessageKind(conv.IntFromAny(data["messageID"])),
				displayMode: tg.DisplayMode(conv.IntFromAny(data["displayMode"])),
				created:     created,
				expires:     created.Add(time.Second * ttl),
			}
		}
	}

	return "", nil
}

// Message contains a UI message for a specific session.
type Message struct {
	messageID   tg.MessageKind
	displayMode tg.DisplayMode
	created     time.Time
	expires     time.Time
}

// Encode encodes the message to JSON.
func (m *Message) Encode(enc *zjson.Encoder, locale string) {
	enc.StartObject()
	enc.IntField("time", int(m.created.UnixMilli()))
	enc.IntField("mode", int(m.displayMode))
	enc.IntField("kind", int(m.messageID))
	enc.StringField("msg", i18n.GetMessage(m.messageID, locale))
	enc.EndObject()
}
