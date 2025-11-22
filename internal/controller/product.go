package controllers

import (
	"net/http"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/infra/logger"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type ProductController struct {
	useCase *usecase.ProductUseCase
	logger  *logger.Logger
}

func NewProductController(useCase *usecase.ProductUseCase) *ProductController {
	return &ProductController{
		useCase: useCase,
		logger:  logger.New(),
	}
}

func (c *ProductController) GetAll(ctx *gin.Context) {
	products, err := c.useCase.GetAll(domain.FindAllProductFilters{Name: ctx.Query("name")})
	if err != nil {
		c.logger.Error("Error fetching products", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Fetching products fatal failed"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"products": products})
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

func (c *ProductController) Update(ctx *gin.Context) {
	var product domain.Product
	if err := ctx.BindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product data"})
		return
	}

	err := c.useCase.Update(product)
	if err != nil {
		c.logger.Error("Failed to update product", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update product"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.useCase.Delete(id)
	if err != nil {
		c.logger.Error("Failed to delete product", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete product"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (c *ProductController) Create(ctx *gin.Context) {
	var newProduct domain.NewProduct
	if err := ctx.BindJSON(&newProduct); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid product data"})
		return
	}

	createdProduct, err := c.useCase.Create(newProduct)
	if err != nil {
		c.logger.Error("Failed to create product", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create product"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdProduct)
}
