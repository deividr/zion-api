package orders

import (
	"github.com/deividr/zion-api/internal/application/use-cases/orders"
	ordersController "github.com/deividr/zion-api/internal/controller/orders"
	"github.com/deividr/zion-api/internal/infra/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOrderByIdControllerFactory(db *pgxpool.Pool) *ordersController.GetOrderByIdController {
	useCase := orders.NewGetOrderByIdUseCase(postgres.NewPgOrderRepository(db))
	return ordersController.NewGetOrderByIdController(useCase)
}
