package upload

import (
	"errors"
	"testing"
)

type mockUploadRepository struct {
	url string
	err error
}

func (m *mockUploadRepository) GetPresignedURL(objectKey string) (string, error) {
	return m.url, m.err
}

func TestUploadUseCase_Execute(t *testing.T) {
	t.Run("should return a presigned URL", func(t *testing.T) {
		mockRepo := &mockUploadRepository{
			url: "http://example.com/presigned-url",
			err: nil,
		}
		useCase := NewUploadUseCase(mockRepo)

		url, err := useCase.Execute()
		if err != nil {
			t.Errorf("expected no error, but got %v", err)
		}

		if url != "http://example.com/presigned-url" {
			t.Errorf("expected URL %s, but got %s", "http://example.com/presigned-url", url)
		}
	})

	t.Run("should return an error when repository fails", func(t *testing.T) {
		mockRepo := &mockUploadRepository{
			url: "",
			err: errors.New("repository error"),
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
