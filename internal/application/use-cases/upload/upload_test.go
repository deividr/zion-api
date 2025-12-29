package upload

import (
	"errors"
	"testing"

	"github.com/deividr/zion-api/internal/domain/services"
)

type mockUploadRepository struct {
	response *services.PresignedURLResponse
	err      error
}

func (m *mockUploadRepository) GetPresignedURL(objectKey string) (*services.PresignedURLResponse, error) {
	return m.response, m.err
}

func TestUploadUseCase_Execute(t *testing.T) {
	t.Run("should return a presigned URL", func(t *testing.T) {
		mockRepo := &mockUploadRepository{
			response: &services.PresignedURLResponse{
				SignedURL: "http://example.com/presigned-url",
				PublicURL: "http://example.com/bucket/object-key",
			},
			err: nil,
		}
		useCase := NewUploadUseCase(mockRepo)

		response, err := useCase.Execute()
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}

		if response.SignedURL != "http://example.com/presigned-url" {
			t.Errorf("expected SignedURL %s, but got %s", "http://example.com/presigned-url", response.SignedURL)
		}

		if response.PublicURL != "http://example.com/bucket/object-key" {
			t.Errorf("expected PublicURL %s, but got %s", "http://example.com/bucket/object-key", response.PublicURL)
		}
	})

	t.Run("should return an error when repository fails", func(t *testing.T) {
		mockRepo := &mockUploadRepository{
			response: nil,
			err:      errors.New("repository error"),
		}
		useCase := NewUploadUseCase(mockRepo)

		_, err := useCase.Execute()

		if err == nil {
			t.Error("expected an error, but got nil")
		}

		if err.Error() != "repository error" {
			t.Errorf("expected error message '%s', but got '%s'", "repository error", err.Error())
		}
	})
}
