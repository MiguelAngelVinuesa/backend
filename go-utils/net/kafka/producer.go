package kafka

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"
)

// Producer is the interface for sending messages to a topic queue.
type Producer interface {
	Produce(ctx context.Context, key, message []byte) error
	ProduceRecords(ctx context.Context, records ...*kgo.Record) error
	WasPinged() bool
	MakeRecord(key, message []byte) *kgo.Record
}

// NewProducer instantiates a new message producer for sending messages to a topic queue.
// The returned interface can be used to query if the brokers are accessible, and to send messages to the topic queue.
// The producer runs a go routine to ping the brokers every 5 seconds.
// Cancel the done context to stop this go routine when an orderly shutdown is requested.
func NewProducer(brokers []string, clientID, topicID string, logger log.BasicLogger, done context.Context) (Producer, error) {
	if logger != nil {
		logger.Info(msgProduceInit, fieldClient, clientID, fieldTopics, topicID, fieldBrokers, brokers)
	}

	p := &producer{
		client:  clientID,
		topic:   topicID,
		brokers: brokers,
		logger:  logger,
		done:    done,
	}

	opts := []kgo.Opt{
		kgo.ClientID(clientID),
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(topicID),
	}

	var err error
	if p.queue, err = kgo.NewClient(opts...); err != nil {
		if logger != nil {
			logger.Error(msgProduceInitFail, fieldClient, clientID, fieldTopics, topicID, fieldError, err)
		}
		return nil, err
	}

	go p.doKeepAlive()
	return p, nil
}

// Produce sends a new message to the topic queue.
func (p *producer) Produce(ctx context.Context, key, message []byte) error {
	err := p.queue.ProduceSync(ctx, p.MakeRecord(key, message)).FirstErr()
	if err != nil && p.logger != nil {
		p.logger.Error(msgProduceSendFail, fieldClient, p.client, fieldTopics, p.topic, fieldError, err)
	}
	return err
}

// ProduceRecords sends one or more messages to the topic queue.
func (p *producer) ProduceRecords(ctx context.Context, records ...*kgo.Record) error {
	err := p.queue.ProduceSync(ctx, records...).FirstErr()
	if err != nil && p.logger != nil {
		p.logger.Error(msgProduceSendFail, fieldClient, p.client, fieldTopics, p.topic, fieldError, err)
	}
	return err
}

// WasPinged returns whether the last ping to the brokers was successful or not.
func (p *producer) WasPinged() bool {
	return atomic.LoadUint32(&p.pinged) != 0
}

func (p *producer) MakeRecord(key, message []byte) *kgo.Record {
	return kgo.KeySliceRecord(key, message)
}

type producer struct {
	client  string
	topic   string
	brokers []string
	logger  log.BasicLogger
	done    context.Context
	queue   *kgo.Client
	pinged  uint32
}

func (p *producer) doKeepAlive() {
	tick := time.NewTicker(50 * time.Millisecond)

	for {
		select {
		case <-p.done.Done():
			p.queue.Close()

			if p.logger != nil {
				p.logger.Info(msgProduceStop, fieldClient, p.client, fieldTopics, p.topic)
			}
			return

		case <-tick.C:
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			err := p.queue.Ping(ctx)
			cancel()

			ok := b32(err)
			if atomic.SwapUint32(&p.pinged, ok) != ok && p.logger != nil {
				if err == nil {
					p.logger.Info(msgProduceStart, fieldClient, p.client, fieldTopics, p.topic)
					tick.Reset(5 * time.Second)
				} else {
					p.logger.Error(msgProducePingFail, fieldClient, p.client, fieldTopics, p.topic, fieldError, err)
					tick.Reset(1 * time.Second)
				}
			}
		}
	}
}

func b32(err error) uint32 {
	if err == nil {
		return 1
	}
	return 0
}
