package controllers

import (
	"net/http"
	"strconv"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	useCase *usecase.CustomerUseCase
}

func NewCustomerController(useCase *usecase.CustomerUseCase) *CustomerController {
	return &CustomerController{useCase: useCase}
}

func (c *CustomerController) GetAll(ctx *gin.Context) {
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

	customers, pagination, err := c.useCase.GetAll(domain.Pagination{Limit: limit, Page: page}, domain.FindAllCustomerFilters{Name: ctx.Query("name")})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"customers": customers, "pagination": pagination})
}

func (c *CustomerController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	customer, err := c.useCase.GetById(id)

	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, customer)
}

func (c *CustomerController) Update(ctx *gin.Context) {
	var customer domain.Customer
	if err := ctx.BindJSON(&customer); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err := c.useCase.Update(customer)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func (c *CustomerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.useCase.Delete(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func (c *CustomerController) Create(ctx *gin.Context) {
	var newCustomer domain.NewCustomer
	if err := ctx.BindJSON(&newCustomer); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	createdCustomer, err := c.useCase.Create(newCustomer)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, createdCustomer)
}
