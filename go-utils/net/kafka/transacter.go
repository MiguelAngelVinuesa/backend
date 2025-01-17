package kafka

import (
	"context"
	"fmt"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Transacter interface {
	Producer
}

// NewTransacter instantiates a new transactional message consumer/producer with the given options.
// The transacter runs a go routine to poll and fetch messages from the topic queue(s), and calls the client interface for each retrieved message.
// The client is responsible for producing messages in the same transaction using the returned interface.
// Cancel the done context to stop the go routine when an orderly shutdown is requested.
// The WithBrokers, WithTopics & WithClient options are mandatory for a functioning transacter!
func NewTransacter(producerID, produceTopic, transactionID string, opts ...ConsumerOption) (Transacter, error) {
	t := &transacter{
		producerID:    producerID,
		producerTopic: produceTopic,
		transactionID: transactionID,
	}
	for ix := range opts {
		opts[ix](&t.consumer)
	}

	if t.group == "" || len(t.topics) == 0 {
		return nil, fmt.Errorf("invalid parameters")
	}

	if t.logger != nil {
		t.logger.Info(msgConsumeInit, fieldGroup, t.group, fieldTopics, t.topics, fieldBrokers, t.brokers)
		t.logger.Info(msgProduceInit, fieldClient, t.producerID, fieldTopics, t.producerTopic, fieldBrokers, t.brokers)
	}

	o := make([]kgo.Opt, 0, 5)
	o = append(o, kgo.SeedBrokers(t.brokers...))
	o = append(o, kgo.DefaultProduceTopic(t.producerTopic))
	o = append(o, kgo.TransactionalID(t.transactionID))
	o = append(o, kgo.FetchIsolationLevel(kgo.ReadCommitted()))
	o = append(o, kgo.ConsumerGroup(t.group))
	o = append(o, kgo.ConsumeTopics(t.topics...))
	o = append(o, kgo.RequireStableFetchOffsets())

	var err error
	if t.sess, err = kgo.NewGroupTransactSession(o...); err != nil {
		if t.logger != nil {
			t.logger.Error(msgTransactInitFail, fieldClient, t.producerID, fieldTopics, t.producerTopic, fieldError, err)
		}
		return nil, err
	}

	go t.doConsume()

	return t, nil
}

// Produce sends a new message to the producer queue.
func (t *transacter) Produce(ctx context.Context, key, message []byte) error {
	t.sess.Produce(ctx, t.MakeRecord(key, message), t.promise.Promise())
	return nil
}

// ProduceRecords sends one or more messages to the producer queue.
func (t *transacter) ProduceRecords(ctx context.Context, records ...*kgo.Record) error {
	for ix := range records {
		t.sess.Produce(ctx, records[ix], t.promise.Promise())
	}
	return nil
}
func (t *transacter) WasPinged() bool { return true }

func (t *transacter) MakeRecord(key, message []byte) *kgo.Record {
	return kgo.KeySliceRecord(key, message)
}

func (t *transacter) doConsume() {
	if t.logger != nil {
		t.logger.Info(msgConsumeStart, fieldGroup, t.group, fieldTopics, t.topics)
	}

	for {
		fetches := t.sess.PollFetches(t.done)
		if errs := fetches.Errors(); len(errs) > 0 {
			t.sess.Close()

			if t.logger != nil {
				if len(errs) == 1 && errs[0].Err == t.done.Err() {
					t.logger.Info(msgConsumeStop, fieldGroup, t.group, fieldTopics, t.topics)
					t.logger.Info(msgProduceStop, fieldClient, t.producerID, fieldTopics, t.producerTopic)
				} else {
					t.logger.Error(msgConsumePollFail, fieldGroup, t.group, fieldTopics, t.topics, fieldErrors, errs)
				}
			}
			return
		}

		if err := t.sess.Begin(); err != nil {
			t.logger.Error(msgTransactBeginFail, fieldGroup, t.group, fieldTopics, t.topics, fieldError, err)
			panic(err) // critical error!
		}

		t.promise = kgo.AbortingFirstErrPromise(t.sess.Client())
		fetches.EachRecord(func(record *kgo.Record) {
			if err := t.client.Consume(record.Key, record.Value, record.Timestamp); err != nil && t.logger != nil {
				t.logger.Error(msgTransactProcessFail, fieldGroup, t.group, fieldTopics, t.topics, fieldError, err)
			}
		})

		if _, err := t.sess.End(context.Background(), t.promise.Err() == nil); err != nil {
			t.logger.Error(msgTransactEndFail, fieldGroup, t.group, fieldTopics, t.topics, fieldError, err)
			panic(err) // critical error!
		}
	}
}

type transacter struct {
	producerID    string
	producerTopic string
	transactionID string
	sess          *kgo.GroupTransactSession
	promise       *kgo.FirstErrPromise
	consumer
}
