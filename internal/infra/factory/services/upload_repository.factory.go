package services

import (
	"os"

	"github.com/deividr/zion-api/internal/domain/services"
	"github.com/deividr/zion-api/internal/infra/storage"
)

func NewTigrisUploadRepository() (services.UploadRepository, error) {
	bucket := os.Getenv("BUCKET_NAME")
	endpoint := os.Getenv("AWS_ENDPOINT_URL_S3")
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	return storage.NewTigris(bucket, endpoint, region, accessKey, secretKey)
}
