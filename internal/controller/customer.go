package controllers

import (
	"net/http"
	"strconv"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/infra/logger"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	useCase *usecase.CustomerUseCase
	logger  *logger.Logger
}

func NewCustomerController(useCase *usecase.CustomerUseCase) *CustomerController {
	return &CustomerController{
		useCase: useCase,
		logger:  logger.New(),
	}
}

func (c *CustomerController) GetAll(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		c.logger.Warn("Invalid limit parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit params"})
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		c.logger.Warn("Invalid page parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit page"})
		return
	}

	customers, pagination, err := c.useCase.GetAll(domain.Pagination{Limit: limit, Page: page}, domain.FindAllCustomerFilters{Name: ctx.Query("name")})
	if err != nil {
		c.logger.Error("Error fetching customers", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Fetching customers fatal failed"})
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
		c.logger.Error("Invalid customer data for update", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid customer data"})
		return
	}

	err := c.useCase.Update(customer)
	if err != nil {
		c.logger.Error("Failed to update customer", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update customer"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func (c *CustomerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.useCase.Delete(id)
	if err != nil {
		c.logger.Error("Failed to delete customer", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete customer"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func (c *CustomerController) Create(ctx *gin.Context) {
	var newCustomer domain.NewCustomer
	if err := ctx.BindJSON(&newCustomer); err != nil {
		c.logger.Error("Invalid customer data for creation", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid customer data"})
		return
	}

	createdCustomer, err := c.useCase.Create(newCustomer)
	if err != nil {
		c.logger.Error("Failed to create customer", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create customer"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdCustomer)
}
