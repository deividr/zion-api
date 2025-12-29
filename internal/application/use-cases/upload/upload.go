package upload

import (
	"github.com/deividr/zion-api/internal/domain/services"
	"github.com/google/uuid"
)

type UploadUseCase struct {
	uploadRepo services.UploadRepository
}

func NewUploadUseCase(uploadRepo services.UploadRepository) *UploadUseCase {
	return &UploadUseCase{
		uploadRepo: uploadRepo,
	}
}

func (uc *UploadUseCase) Execute() (*services.PresignedURLResponse, error) {
	objectKey := uuid.New().String()
	return uc.uploadRepo.GetPresignedURL(objectKey)
}
