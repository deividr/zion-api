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

func (c *AddressController) GetById(ctx *gin.Context) {
	id := ctx.Param("id")

	address, err := c.addressUseCase.GetById(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Address not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"address": address})
}

func (c *AddressController) Update(ctx *gin.Context) {
	var address domain.Address
	if err := ctx.BindJSON(&address); err != nil {
		c.logger.Error("Invalid address data for update", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid address data"})
		return
	}

	err := c.addressUseCase.Update(address)
	if err != nil {
		c.logger.Error("Failed to update address", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update address"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
}

func (c *AddressController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.addressUseCase.Delete(id)
	if err != nil {
		c.logger.Error("Failed to delete address", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete address"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

func (c *AddressController) Create(ctx *gin.Context) {
	var newAddress domain.NewAddress
	if err := ctx.BindJSON(&newAddress); err != nil {
		c.logger.Error("Invalid address data for creation", err)
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid address data"})
		return
	}

	createdAddress, err := c.addressUseCase.Create(newAddress)
	if err != nil {
		c.logger.Error("Failed to create address", err)
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create address"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdAddress)
}
