package controllers

import (
	"net/http"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/infra/logger"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AddressController struct {
	addressUseCase *usecase.AddressUseCase
	logger         *logger.Logger
}

func NewAddressController(addressUseCase *usecase.AddressUseCase) *AddressController {
	return &AddressController{
		addressUseCase: addressUseCase,
		logger:         logger.New(),
	}
}

func (c *AddressController) GetByCustomerId(ctx *gin.Context) {
	customerId := ctx.Param("id")

	addresses, err := c.addressUseCase.GetByCustomerId(customerId)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Addresses not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"addresses": addresses})
}

func (c *AddressController) Update(ctx *gin.Context) {
	customerId := ctx.Param("id")
	addressId := ctx.Param("addressId")

	var updateData domain.NewAddress
	if err := ctx.BindJSON(&updateData); err != nil {
		c.logger.Error("Invalid address data for update", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid address data"})
		return
	}

	updatedAddress, err := c.addressUseCase.Update(customerId, addressId, updateData)
	if err != nil {
		c.logger.Error("Failed to update address", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update address"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, updatedAddress)
}

func (c *AddressController) Delete(ctx *gin.Context) {
	customerId := ctx.Param("id")
	addressId := ctx.Param("addressId")

	err := c.addressUseCase.Delete(customerId, addressId)
	if err != nil {
		c.logger.Error("Failed to delete address", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete address"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

func (c *AddressController) Create(ctx *gin.Context) {
	customerId := ctx.Param("id")

	var newAddress domain.NewAddress
	if err := ctx.BindJSON(&newAddress); err != nil {
		c.logger.Error("Invalid address data for creation", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid address data"})
		return
	}

	createdAddress, err := c.addressUseCase.Create(customerId, newAddress)
	if err != nil {
		c.logger.Error("Failed to create address", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create address"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdAddress)
}
