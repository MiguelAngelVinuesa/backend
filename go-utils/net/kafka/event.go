package kafka

import (
	"context"

	"github.com/goccy/go-json"
)

// OperationKind defines the kind of operation on the entity of an event message.
type OperationKind uint8

// NOTE: always add new operations at the end of the list!
const (
	OpCreate OperationKind = iota + 1
	OpRead
	OpUpdate
	OpDelete
	OpDisable
	OpEnable
	OpBlock
	OpUnblock
)

// Event contains the details of an event message.
type Event struct {
	Producer   string        `json:"producer"`       // ClientID of the sender.
	EntityCode string        `json:"entity"`         // Entity operated on (e.g. database table, etc).
	EntityKey  string        `json:"key,omitempty"`  // Unique key for the entity, or a random string.
	Operation  OperationKind `json:"op"`             // Operation kind.
	Data       any           `json:"data,omitempty"` // Relevant data of the entity that was operated on.
	// hidden fields.
	key []byte
}

// NewEvent instantiates a new event message.
func NewEvent(clientID, entityCode, entityKey string, op OperationKind, data any) *Event {
	return &Event{
		Producer:   clientID,
		EntityCode: entityCode,
		EntityKey:  entityKey,
		Operation:  op,
		Data:       data,
		key:        []byte(entityCode + ":" + entityKey),
	}
}

// NewEventFromJSON instantiates a new event message from the given JSON data.
func NewEventFromJSON(event []byte, dataModel any) (*Event, error) {
	e := &Event{Data: dataModel}
	if err := json.Unmarshal(event, e); err != nil {
		return nil, err
	}
	e.key = []byte(e.EntityCode + ":" + e.EntityKey)
	return e, nil
}

// Produce sends the message to the given message queue producer.
func (e *Event) Produce(ctx context.Context, producer Producer) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return producer.Produce(ctx, e.key, b)
}
