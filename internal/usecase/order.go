package usecase

import (
	"fmt"
	"time"

	"github.com/deividr/zion-api/internal/domain"
)

type CreateOrderInput struct {
	PickupDate   time.Time             `json:"pickupDate"`
	CustomerId   string                `json:"customerId"`
	AddressId    *string               `json:"addressId"`
	Employee     string                `json:"employee"`
	OrderLocal   *string               `json:"orderLocal"`
	Observations *string               `json:"observations"`
	IsPickedUp   *bool                 `json:"isPickedUp"`
	Products     []domain.OrderProduct `json:"products"`
}

type UpdateOrderInput struct {
	Id           string                `json:"id"`
	PickupDate   time.Time             `json:"pickupDate"`
	AddressId    *string               `json:"addressId"`
	Employee     string                `json:"employee"`
	OrderLocal   *string               `json:"orderLocal"`
	Observations *string               `json:"observations"`
	IsPickedUp   *bool                 `json:"isPickedUp"`
	Products     []domain.OrderProduct `json:"products"`
}

type OrderUseCase struct {
	repo         domain.OrderRepository
	addressRepo  domain.AddressRepository
	customerRepo domain.CustomerRepository
}

func NewOrderUseCase(repo domain.OrderRepository, addressRepo domain.AddressRepository, customerRepo domain.CustomerRepository) *OrderUseCase {
	return &OrderUseCase{repo: repo, addressRepo: addressRepo, customerRepo: customerRepo}
}

func (uc *OrderUseCase) GetAll(pagination domain.Pagination, filters domain.FindAllOrderFilters) ([]domain.Order, domain.Pagination, error) {
	orders, pagination, err := uc.repo.FindAll(pagination, filters)
	if err != nil {
		return []domain.Order{}, domain.Pagination{}, fmt.Errorf("error fetching orders: %v", err)
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

func (uc *OrderUseCase) Update(input UpdateOrderInput) error {
	order := domain.Order{
		Id:           input.Id,
		PickupDate:   input.PickupDate,
		Employee:     input.Employee,
		OrderLocal:   input.OrderLocal,
		Observations: input.Observations,
		IsPickedUp:   input.IsPickedUp,
		Products:     input.Products,
	}

	if input.AddressId != nil {
		address, err := uc.addressRepo.FindById(*input.AddressId)
		if err != nil {
			return fmt.Errorf("address not found: %v", err)
		}
		order.SetAddress(address)
	}

	if err := uc.repo.Update(order); err != nil {
		return fmt.Errorf("error updating order: %v", err)
	}
	return nil
}

func (uc *OrderUseCase) Delete(id string) error {
	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("error deleting order: %v", err)
	}
	return nil
}

func (uc *OrderUseCase) Create(input CreateOrderInput) (*domain.Order, error) {
	customer, err := uc.customerRepo.FindById(input.CustomerId)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %v", err)
	}

	order := domain.Order{
		PickupDate:   input.PickupDate,
		Customer:     *customer,
		Employee:     input.Employee,
		OrderLocal:   input.OrderLocal,
		Observations: input.Observations,
		IsPickedUp:   input.IsPickedUp,
		Products:     input.Products,
	}

	if input.AddressId != nil {
		address, err := uc.addressRepo.FindById(*input.AddressId)
		if err != nil {
			return nil, fmt.Errorf("address not found: %v", err)
		}
		order.SetAddress(address)
	}

	createdOrder, err := uc.repo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}

	return createdOrder, nil
}
