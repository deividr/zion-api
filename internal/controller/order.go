package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/deividr/zion-api/internal/domain"
	"github.com/deividr/zion-api/internal/infra/logger"
	"github.com/deividr/zion-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	useCase *usecase.OrderUseCase
	logger  *logger.Logger
}

func NewOrderController(useCase *usecase.OrderUseCase) *OrderController {
	return &OrderController{
		useCase: useCase,
		logger:  logger.New(),
	}
}

func (c *OrderController) GetAll(ctx *gin.Context) {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit params"})
		return
	}

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit page"})
		return
	}

	search := ctx.Query("search")

	pickupDateStart, err := time.Parse(time.RFC3339, ctx.Query("pickupDateStart"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pickupDateStart params"})
		return
	}

	pickupDateEnd, err := time.Parse(time.RFC3339, ctx.Query("pickupDateEnd"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pickupDateEnd params"})
		return
	}

	orders, pagination, err := c.useCase.GetAll(domain.Pagination{Limit: limit, Page: page}, domain.FindAllOrderFilters{Search: &search, PickupDateStart: pickupDateStart, PickupDateEnd: pickupDateEnd})
	if err != nil {
		c.logger.Error("Error fetching orders", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Fetching orders fatal failed"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"orders": orders, "pagination": pagination})
}

func (c *OrderController) Update(ctx *gin.Context) {
	var order domain.Order
	if err := ctx.BindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order data"})
		return
	}

	err := c.useCase.Update(order)
	if err != nil {
		c.logger.Error("Failed to update order", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Order updated successfully"})
}

func (c *OrderController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.useCase.Delete(id)
	if err != nil {
		c.logger.Error("Failed to delete order", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete order"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

func (c *OrderController) Create(ctx *gin.Context) {
	var newOrder domain.NewOrder
	if err := ctx.BindJSON(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order data"})
		return
	}

	createdOrder, err := c.useCase.Create(newOrder)
	if err != nil {
		c.logger.Error("Failed to create order", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order"})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, createdOrder)
}
