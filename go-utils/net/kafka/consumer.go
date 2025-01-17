package kafka

import (
	"context"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"
)

// Consumer is the callback interface for processing messages from a topic queue.
type Consumer interface {
	Consume(key, msg []byte, timestamp time.Time) error
}

// Committer is the interface to commit records after they have been processed successfully.
type Committer interface {
	Commit()
}

// NewConsumer instantiates a new message consumer with the given options.
// The consumer runs a go routine to poll and fetch messages from the topic queue(s), and calls the client interface for each retrieved message.
// Cancel the done context to stop the go routine when an orderly shutdown is requested.
// The returned Committer interface can be used to commit processed messages when processing within a consumer group.
// The WithBrokers, WithTopics & WithClient options are mandatory for a functioning consumer!
func NewConsumer(opts ...ConsumerOption) (Committer, error) {
	c := &consumer{}
	for ix := range opts {
		opts[ix](c)
	}

	if c.logger != nil {
		c.logger.Info(msgConsumeInit, fieldGroup, c.group, fieldTopics, c.topics, fieldBrokers, c.brokers)
	}

	o := make([]kgo.Opt, 0, 5)
	o = append(o, kgo.SeedBrokers(c.brokers...))

	if len(c.topics) > 0 {
		o = append(o, kgo.ConsumeTopics(c.topics...))
	}
	if c.group != "" {
		o = append(o, kgo.DisableAutoCommit())
		o = append(o, kgo.ConsumerGroup(c.group))
	}
	if !c.from.IsZero() {
		o = append(o, kgo.ConsumeResetOffset(kgo.NewOffset().AfterMilli(c.from.UnixMilli())))
	}

	var err error
	if c.queue, err = kgo.NewClient(o...); err != nil {
		if c.logger != nil {
			c.logger.Error(msgConsumeInitFail, fieldGroup, c.group, fieldTopics, c.topics, fieldError, err)
		}
		return nil, err
	}

	go c.doConsume()
	return c, nil
}

// Commit commits all uncommitted messages.
// It should be called after the messages have been processed successfully.
// Consuming messages must be done within a consumer group, and the groupCommit parameter must be false.
func (c *consumer) Commit() {
	if c.group != "" && !c.groupCommit {
		if err := c.queue.CommitUncommittedOffsets(context.Background()); err != nil && c.logger != nil {
			c.logger.Error(msgConsumeCommitFail, fieldGroup, c.group, fieldTopics, c.topics, fieldError, err)
		}
	}
}

// ConsumerOption is the function prototype for options when instantiating a new message consumer.
type ConsumerOption func(c *consumer)

// WithBrokers adds the list of brokers to the message consumer.
// This option is mandatory.
func WithBrokers(brokers ...string) ConsumerOption {
	return func(c *consumer) {
		c.brokers = brokers
	}
}

// WithTopics adds one or more topics to read from.
// This option is mandatory.
func WithTopics(topics ...string) ConsumerOption {
	return func(c *consumer) {
		c.topics = topics
	}
}

// WithClient adds the message callback function to the message consumer.
// This option is mandatory.
func WithClient(client Consumer) ConsumerOption {
	return func(c *consumer) {
		c.client = client
	}
}

// WithGroup sets up the message consumer as part of a consumer group.
// The commit parameter should be true to have the consumer commit processed messages automatically.
// Set it to false if you want to control the commit yourself!
func WithGroup(group string, commit bool) ConsumerOption {
	return func(c *consumer) {
		c.group = group
		c.groupCommit = commit
	}
}

// FromTimestamp adds a timestamp to start reading message from.
// Do not use it together with a consumer group!
func FromTimestamp(from time.Time) ConsumerOption {
	return func(c *consumer) {
		c.from = from
	}
}

// WithLogger adds a basic logger to the message consumer.
func WithLogger(logger log.BasicLogger) ConsumerOption {
	return func(c *consumer) {
		c.logger = logger
	}
}

// WithDone adds the context for a clean shutdown of the message consumer.
// Cancel the given context to stop the consumer.
func WithDone(done context.Context) ConsumerOption {
	return func(c *consumer) {
		c.done = done
	}
}

func (c *consumer) doConsume() {
	if c.logger != nil {
		c.logger.Info(msgConsumeStart, fieldGroup, c.group, fieldTopics, c.topics)
	}

	for {
		fetches := c.queue.PollFetches(c.done)
		if errs := fetches.Errors(); len(errs) > 0 {
			c.queue.Close()

			if c.logger != nil {
				if len(errs) == 1 && errs[0].Err == c.done.Err() {
					c.logger.Info(msgConsumeStop, fieldGroup, c.group, fieldTopics, c.topics)
				} else {
					c.logger.Error(msgConsumePollFail, fieldGroup, c.group, fieldTopics, c.topics, fieldErrors, errs)
				}
			}
			return
		}

		fetches.EachRecord(func(record *kgo.Record) {
			if err := c.client.Consume(record.Key, record.Value, record.Timestamp); err != nil && c.logger != nil {
				c.logger.Error(msgConsumeProcessFail, fieldGroup, c.group, fieldTopics, c.topics, fieldError, err)
			}
		})

		if c.group != "" && c.groupCommit {
			if err := c.queue.CommitUncommittedOffsets(context.Background()); err != nil && c.logger != nil {
				c.logger.Error(msgConsumeCommitFail, fieldGroup, c.group, fieldTopics, c.topics, fieldError, err)
			}
		}
	}
}

type consumer struct {
	group       string
	groupCommit bool
	topics      []string
	brokers     []string
	from        time.Time
	client      Consumer
	logger      log.BasicLogger
	done        context.Context
	queue       *kgo.Client
}
