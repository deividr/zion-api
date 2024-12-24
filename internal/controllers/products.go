package controllers

import (
	"net/http"

	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	useCase *usecase.ProductUseCase
}

func NewProductController(useCase *usecase.ProductUseCase) *ProductController {
	return &ProductController{useCase: useCase}
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	products, err := c.useCase.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, products)
}

func (c *ProductController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.useCase.GetById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, product)
}
