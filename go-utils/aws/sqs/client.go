package sns

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	http2 "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Client is the interface wrapper for communicating with AWS SNS.
// It implements generic functions for all TopGaming use cases.
type Client interface {
	CreateQueue(ctx context.Context, queue string) (string, string, error)
	ReceiveMessages(ctx context.Context, url string, max, timeout int) ([]types.Message, error)
	DeleteMessage(ctx context.Context, url, handle string) error
	DeleteQueue(ctx context.Context, url string) error
	SetPolicy(ctx context.Context, url, policy string) error
}

// NewClient instantiates a new SQS client.
// It sets the HTTP client with a timeout a little over 10 minutes to allow for long polling.
func NewClient(cfg *aws.Config, policy string) Client {
	cfg2 := cfg.Copy()
	cfg2.HTTPClient = http2.NewBuildableClient().WithTimeout(610 * time.Second)

	return &client{
		client: sqs.NewFromConfig(cfg2),
		policy: policy,
	}
}

// CreateQueue create a new SQS queue with the given name and policy.
func (c *client) CreateQueue(ctx context.Context, queue string) (string, string, error) {
	req := &sqs.CreateQueueInput{
		QueueName: aws.String(queue),
		Attributes: map[string]string{
			"FifoQueue":                 "true",
			"MaximumMessageSize":        "16384",
			"MessageRetentionPeriod":    "60",
			"Policy":                    c.policy,
			"ContentBasedDeduplication": "false",
			"DeduplicationScope":        "messageGroup",
			"FifoThroughputLimit":       "perMessageGroupId",
		},
	}

	resp, err := c.client.CreateQueue(ctx, req)
	if err != nil {
		return "", "", err
	}

	url := aws.ToString(resp.QueueUrl)

	req2 := &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(url),
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameQueueArn},
	}

	var resp2 *sqs.GetQueueAttributesOutput
	resp2, err = c.client.GetQueueAttributes(ctx, req2)

	return url, resp2.Attributes[string(types.QueueAttributeNameQueueArn)], nil
}

// ReceiveMessages tries to retrieve one or more messages from the SQS queue with the given URL.
// max defines the maximum number of messages to retrieve.
// If zero, a maximum of 10 is applied; max cannot be larger than 10!
// timeout indicates how long (in seconds) the call should wait for a message to arrive.
// if zero, timeout is set to 60 seconds; timeout cannot be larrger than 600 (e.g. 10 minutes)!
func (c *client) ReceiveMessages(ctx context.Context, url string, max, timeout int) ([]types.Message, error) {
	if max <= 0 || max > 10 {
		max = 10
	}

	if timeout <= 0 {
		timeout = 60 // 1 minute.
	} else if timeout > 600 {
		timeout = 600 // 10 minutes.
	}

	req := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(url),
		MessageAttributeNames: []string{string(types.QueueAttributeNameAll)},
		VisibilityTimeout:     0,
		MaxNumberOfMessages:   int32(max),
		WaitTimeSeconds:       int32(timeout),
	}

	resp, err := c.client.ReceiveMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

// DeleteMessage removes the message with the given handle from the SQS queue with the given url.
// This function must be called after a received message has been fully processed.
func (c *client) DeleteMessage(ctx context.Context, url, handle string) error {
	req := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(url),
		ReceiptHandle: aws.String(handle),
	}

	_, err := c.client.DeleteMessage(ctx, req)
	return err
}

// DeleteQueue removes the SQS queue with the given url.
func (c *client) DeleteQueue(ctx context.Context, url string) error {
	req := &sqs.DeleteQueueInput{
		QueueUrl: aws.String(url),
	}

	_, err := c.client.DeleteQueue(ctx, req)
	return err
}

// SetPolicy updates the policy for the SQS queue with the given url.
func (c *client) SetPolicy(ctx context.Context, url, policy string) error {
	req := &sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(url),
		Attributes: map[string]string{
			"Policy": policy,
		},
	}

	_, err := c.client.SetQueueAttributes(ctx, req)
	return err
}

type client struct {
	client *sqs.Client
	policy string
}
