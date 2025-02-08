package usecase

import (
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
)

type CustomerUseCase struct {
	repo domain.CustomerRepository
}

func NewCustomerUseCase(repo domain.CustomerRepository) *CustomerUseCase {
	return &CustomerUseCase{repo: repo}
}

func (uc *CustomerUseCase) GetAll(pagination domain.Pagination, filters domain.FindAllCustomerFilters) ([]domain.Customer, domain.Pagination, error) {
	products, pagination, err := uc.repo.FindAll(pagination, filters)

	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("erro ao buscar clientes: %v", err)
	}

	return products, pagination, nil
}

func (uc *CustomerUseCase) GetById(id string) (*domain.Customer, error) {
	product, err := uc.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar cliente: %v", err)
	}
	return product, nil
}

func (uc *CustomerUseCase) Update(product domain.Customer) error {
	err := uc.repo.Update(product)
	if err != nil {
		return fmt.Errorf("erro ao atualizar cliente: %v", err)
	}
	return nil
}

func (uc *CustomerUseCase) Delete(id string) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("erro ao deletar cliente: %v", err)
	}
	return nil
}

func (uc *CustomerUseCase) Create(newCustomer domain.NewCustomer) (*domain.Customer, error) {
	createdCustomer, err := uc.repo.Create(newCustomer)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente: %v", err)
	}

	return createdCustomer, nil
}
