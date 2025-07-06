package controllers

import (
	"net/http"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/infra/logger"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CategoryProductController struct {
	useCase *usecase.CategoryProductUseCase
	logger  *logger.Logger
}

func NewCategoryProductController(useCase *usecase.CategoryProductUseCase) *CategoryProductController {
	return &CategoryProductController{
		useCase: useCase,
		logger:  logger.New(),
	}
}

func (c *CategoryProductController) GetAll(ctx *gin.Context) {
	categories, err := c.useCase.GetAll()
	if err != nil {
		c.logger.Error("Error fetching categories", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Fetching categories fatal failed"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, categories)
}

func (c *CategoryProductController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	category, err := c.useCase.GetById(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Category not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, category)
}

func (c *CategoryProductController) Update(ctx *gin.Context) {
	var category domain.CategoryProduct
	if err := ctx.BindJSON(&category); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category data"})
		return
	}

	err := c.useCase.Update(category)
	if err != nil {
		c.logger.Error("Failed to update category", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update category"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Category updated successfully"})
}

func (c *CategoryProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.useCase.Delete(id)
	if err != nil {
		c.logger.Error("Failed to delete category", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete category"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

func (c *CategoryProductController) Create(ctx *gin.Context) {
	var newCategory domain.CategoryProduct
	if err := ctx.BindJSON(&newCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid category data"})
		return
	}

	createdCategory, err := c.useCase.Create(newCategory)
	if err != nil {
		c.logger.Error("Failed to create category", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create category"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdCategory)
}
