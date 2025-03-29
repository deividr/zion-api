package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/infra/logger"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	customerUseCase *usecase.CustomerUseCase
	addressUseCase  *usecase.AddressUseCase
	logger          *logger.Logger
}

func NewCustomerController(customerUseCase *usecase.CustomerUseCase, addressUseCase *usecase.AddressUseCase) *CustomerController {
	return &CustomerController{
		customerUseCase: customerUseCase,
		addressUseCase:  addressUseCase,
		logger:          logger.New(),
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

	customers, pagination, err := c.customerUseCase.GetAll(domain.Pagination{Limit: limit, Page: page}, domain.FindAllCustomerFilters{Name: ctx.Query("name")})
	if err != nil {
		c.logger.Error("Error fetching customers", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Fetching customers fatal failed"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"customers": customers, "pagination": pagination})
}

func (c *CustomerController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	customer, err := c.customerUseCase.GetById(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
		return
	}

	addresses, err := c.addressUseCase.GetBy(map[string]interface{}{"customer_id": customer.Id})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Error to get address by customer id"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"customer": customer, "addresses": addresses})
}

func (c *CustomerController) Update(ctx *gin.Context) {
	var customer domain.Customer
	if err := ctx.BindJSON(&customer); err != nil {
		c.logger.Error("Invalid customer data for update", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid customer data"})
		return
	}

	err := c.customerUseCase.Update(customer)
	if err != nil {
		c.logger.Error("Failed to update customer", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update customer"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Customer updated successfully"})
}

func (c *CustomerController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.customerUseCase.Delete(id)
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

	createdCustomer, err := c.customerUseCase.Create(newCustomer)
	if err != nil {
		c.logger.Error("Failed to create customer", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create customer"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdCustomer)
}
