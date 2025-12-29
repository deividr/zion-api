package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Tigris struct {
	client *s3.PresignClient
	bucket string
}

func NewTigris(bucket, endpoint, region, accessKey, secretKey string) (*Tigris, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	presignClient := s3.NewPresignClient(s3Client)

	return &Tigris{
		client: presignClient,
		bucket: bucket,
	}, nil
}

func (t *Tigris) GetPresignedURL(objectKey string) (string, error) {
	req, err := t.client.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &t.bucket,
		Key:    &objectKey,
	}, func(po *s3.PresignOptions) {
		po.Expires = 15 * time.Minute
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return req.URL, nil
}
