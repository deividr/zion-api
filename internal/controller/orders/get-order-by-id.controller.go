package orders

import (
	"net/http"

	"github.com/deividr/zion-api/internal/application/use-cases/orders"
	"github.com/gin-gonic/gin"
)

type GetOrderByIdController struct {
	useCase orders.GetOrderByIdUseCase
}

func (c *GetOrderByIdController) Handle(ctx *gin.Context) {
	id := ctx.Param("id")

	order, err := c.useCase.Execute(id)
	if err != nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "Order not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, order)
}

func NewGetOrderByIdController(useCase *orders.GetOrderByIdUseCase) *GetOrderByIdController {
	return &GetOrderByIdController{useCase: *useCase}
}
