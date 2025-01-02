package controllers

import (
	"net/http"
	"strconv"

	"github.com/deividr/zion-api/internal/domain"
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
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid limit params"})
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid limit page"})
		return
	}

	products, pagination, err := c.useCase.GetAll(domain.Pagination{Limit: limit, Page: page}, domain.FindAllProductFilters{Name: ctx.Query("name")})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"products": products, "pagination": pagination})
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
