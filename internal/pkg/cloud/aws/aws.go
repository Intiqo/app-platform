package aws

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// NewAWSConfig returns a new aws.Config
func NewAWSConfig(profile string) (cfg aws.Config) {
	var err error
	if profile != "" {
		cfg, err = awsCfg.LoadDefaultConfig(context.TODO(), awsCfg.WithSharedConfigProfile(profile))
	} else {
		cfg, err = awsCfg.LoadDefaultConfig(context.TODO())
	}
	if err != nil {
		log.Print("failed to load default aws config. Trying with environment variables")
		region := os.Getenv("AWS_REGION")
		accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
		secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
		if region == "" || accessKeyId == "" || secretAccessKey == "" {
			log.Fatal("Missing AWS credentials. Please set AWS_REGION, AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables")
		}
		cfg = aws.Config{
			Region:      region,
			Credentials: credentials.NewStaticCredentialsProvider(accessKeyId, secretAccessKey, ""),
		}
	}
	return cfg
}
