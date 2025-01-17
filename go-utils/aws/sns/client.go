package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// Client is the interface wrapper for communicating with AWS SNS.
// It implements generic functions for all TopGaming use cases.
type Client interface {
	Publish(ctx context.Context, arn, group, msg, dedup string) (string, error)
	SubscribeSQS(ctx context.Context, arn, sqs string) (string, error)
	UnsubscribeSQS(ctx context.Context, arn string) error
}

// NewClient instantiates a new SNS client.
func NewClient(cfg *aws.Config) Client {
	return &client{
		client: sns.NewFromConfig(*cfg),
	}
}

// Publish publised the message to the SNS arn.
func (c *client) Publish(ctx context.Context, arn, group, msg, dedup string) (string, error) {
	req := &sns.PublishInput{
		TopicArn:               aws.String(arn),
		MessageGroupId:         aws.String(group),
		Message:                aws.String(msg),
		MessageDeduplicationId: aws.String(dedup),
	}

	resp, err := c.client.Publish(ctx, req)
	if err != nil {
		return "", err
	}

	return aws.ToString(resp.MessageId), nil
}

// SubscribeSQS subscribes the SQS arn to the SNS arn.
func (c *client) SubscribeSQS(ctx context.Context, arn, sqs string) (string, error) {
	req := &sns.SubscribeInput{
		Protocol: aws.String("sqs"),
		TopicArn: aws.String(arn),
		Endpoint: aws.String(sqs),
		Attributes: map[string]string{
			"RawMessageDelivery": "true",
		},
	}

	resp, err := c.client.Subscribe(ctx, req)
	if err != nil {
		return "", err
	}

	return aws.ToString(resp.SubscriptionArn), nil
}

// UnsubscribeSQS subscribes the SQS arn to the SNS arn.
func (c *client) UnsubscribeSQS(ctx context.Context, subscription string) error {
	req := &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(subscription),
	}

	_, err := c.client.Unsubscribe(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

type client struct {
	client *sns.Client
}
