package services

type PresignedURLResponse struct {
	SignedURL string `json:"signedUrl"`
	PublicURL string `json:"publicUrl"`
}

type UploadRepository interface {
	GetPresignedURL(objectKey string) (*PresignedURLResponse, error)
}
