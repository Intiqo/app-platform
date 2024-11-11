package file

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/Intiqo/app-platform/internal/pkg/config"
)

type awsS3Manager struct {
	cfg    config.AppConfig
	client *s3.Client
}

func NewFileManager(cfg config.AppConfig, awsCfg aws.Config) Manager {
	c := s3.NewFromConfig(awsCfg)
	s := &awsS3Manager{
		cfg:    cfg,
		client: c,
	}
	return s
}

func (m *awsS3Manager) UploadFile(opts Options) (result string, err error) {
	var key, url string
	if opts.Directory != "" {
		key = fmt.Sprintf("%s/%s", opts.Directory, opts.Filename)
	} else {
		key = opts.Filename
	}

	upParams := &s3.PutObjectInput{
		Bucket:      aws.String(opts.Bucket),
		Key:         aws.String(key),
		Body:        opts.File,
		ContentType: aws.String(opts.ContentType),
	}
	_, err = m.client.PutObject(context.TODO(), upParams)
	if err != nil {
		slog.Error("failed to upload file to S3")
		return url, err
	}
	url = m.constructUrlForFile(key, opts)
	return url, nil
}

func (m *awsS3Manager) constructUrlForFile(key string, opts Options) string {
	return fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", opts.Bucket, opts.Region, key)
}
