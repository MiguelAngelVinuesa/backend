package config

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/aws/consts"
)

// NewConfig instantiates a new generic AWS configuration.
func NewConfig(ctx context.Context, region string) (*aws.Config, error) {
	if region == "" {
		region = consts.DefaultRegion
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	return &cfg, err
}
