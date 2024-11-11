package secrets

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// awsSecretsManager is a wrapper around the AWS Secrets Manager
type awsSecretsManager struct {
	awsConfig aws.Config
	client    *secretsmanager.Client
}

// NewAWSSecretsManager returns a new Manager
func NewAWSSecretsManager(awsCfg aws.Config) Manager {
	sm := &awsSecretsManager{
		awsConfig: awsCfg,
	}
	sm.client = secretsmanager.NewFromConfig(sm.awsConfig)
	return sm
}

// GetSecret returns the secret value for the given name
func (s *awsSecretsManager) GetSecret(name string) (result string, err error) {
	valOut, err := s.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	})
	if err != nil {
		return result, err
	}
	if valOut == nil {
		return result, errors.New("could not find the secret")
	}
	return *valOut.SecretString, err
}
