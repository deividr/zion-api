package usecase

import (
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
)

type OrderUseCase struct {
	repo domain.OrderRepository
}

func NewOrderUseCase(repo domain.OrderRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) GetAll(pagination domain.Pagination, filters domain.FindAllOrderFilters) ([]domain.Order, domain.Pagination, error) {
	orders, pagination, err := uc.repo.FindAll(pagination, filters)
	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error fetching orders: %v", err)
	}

	return orders, pagination, nil
}

func (uc *OrderUseCase) GetById(id string) (*domain.Order, error) {
	order, err := uc.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching order by id: %v", err)
	}
	return order, nil
}

func (uc *OrderUseCase) Update(order domain.Order) error {
	err := uc.repo.Update(order)
	if err != nil {
		return fmt.Errorf("error updating order: %v", err)
	}
	return nil
}

func (uc *OrderUseCase) Delete(id string) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting order: %v", err)
	}
	return nil
}

func (uc *OrderUseCase) Create(newOrder domain.NewOrder) (*domain.Order, error) {
	createdOrder, err := uc.repo.Create(newOrder)

	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}

	return createdOrder, nil
}
