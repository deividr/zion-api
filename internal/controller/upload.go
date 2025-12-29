package controller

import (
	"net/http"

	"github.com/deividr/zion-api/internal/application/use-cases/upload"
	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadUseCase *upload.UploadUseCase
}

func NewUploadController(uploadUseCase *upload.UploadUseCase) *UploadController {
	return &UploadController{
		uploadUseCase: uploadUseCase,
	}
}

func (c *UploadController) GetPresignedURL(ctx *gin.Context) {
	response, err := c.uploadUseCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"signedUrl": response.SignedURL,
		"publicUrl": response.PublicURL,
	})
}
