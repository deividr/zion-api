package orders

import (
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
)

type GetOrderByIdUseCase struct {
	orderRepository domain.OrderRepository
}

func (uc *GetOrderByIdUseCase) Execute(id string) (*domain.Order, error) {
	order, err := uc.orderRepository.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching order by id: %v", err)
	}

	return order, nil
}

func NewGetOrderByIdUseCase(orderRepository domain.OrderRepository) *GetOrderByIdUseCase {
	return &GetOrderByIdUseCase{
		orderRepository: orderRepository,
	}
}
