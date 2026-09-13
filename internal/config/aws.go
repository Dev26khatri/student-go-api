package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

// LoadAWSConfig loads AWS settings from the environment, shared AWS files,
// or other sources supported by the AWS SDK.
func LoadAWSConfig() (aws.Config, error) {
	cfg, err := awscfg.LoadDefaultConfig(
		context.Background(),
	)
	if err != nil {
		log.Fatal("Failed to load AWS configuration 📦", err)
	}
	return cfg, nil
}
