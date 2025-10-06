package orders

import (
	"github.com/deividr/zion-api/internal/application/use-cases/orders"
	"github.com/deividr/zion-api/internal/infra/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetOrderByIdUsecaseFactory(db *pgxpool.Pool) *orders.GetOrderByIdUseCase {
	orderRepository := postgres.NewPgOrderRepository(db)
	return orders.NewGetOrderByIdUseCase(orderRepository)
}
