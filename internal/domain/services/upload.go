package services

type UploadRepository interface {
	GetPresignedURL(objectKey string) (string, error)
}
