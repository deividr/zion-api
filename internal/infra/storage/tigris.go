package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/deividr/zion-api/internal/domain/services"
)

type Tigris struct {
	client        *s3.PresignClient
	bucket        string
	endpoint      string
	publicURLBase string
}

func NewTigris(bucket, endpoint, publicURLBase, region, accessKey, secretKey string) (*Tigris, error) {
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

	// Se publicURLBase não foi fornecida, inferir do padrão do Tigris
	if publicURLBase == "" {
		publicURLBase = fmt.Sprintf("https://%s.t3.storage.dev", bucket)
	}

	return &Tigris{
		client:        presignClient,
		bucket:        bucket,
		endpoint:      endpoint,
		publicURLBase: publicURLBase,
	}, nil
}

func (t *Tigris) GetPresignedURL(objectKey string) (*services.PresignedURLResponse, error) {
	req, err := t.client.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &t.bucket,
		Key:    &objectKey,
	}, func(po *s3.PresignOptions) {
		po.Expires = 15 * time.Minute
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Construir a URL pública
	publicURL := t.buildPublicURL(objectKey)

	return &services.PresignedURLResponse{
		SignedURL: req.URL,
		PublicURL: publicURL,
	}, nil
}

func (t *Tigris) buildPublicURL(objectKey string) string {
	publicURLBase := strings.TrimSuffix(t.publicURLBase, "/")
	return fmt.Sprintf("%s/%s", publicURLBase, objectKey)
}
