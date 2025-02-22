package usecase

import (
	"fmt"

	"github.com/deividr/zion-api/internal/domain"
)

type AddressUseCase struct {
	repo domain.AddressRepository
}

func NewAddressUseCase(repo domain.AddressRepository) *AddressUseCase {
	return &AddressUseCase{repo: repo}
}

func (uc *AddressUseCase) GetAll(pagination domain.Pagination) ([]domain.Address, domain.Pagination, error) {
	addresses, pagination, err := uc.repo.FindAll(pagination)

	if err != nil {
		return nil, domain.Pagination{}, fmt.Errorf("error fetching addresses: %v", err)
	}

	return addresses, pagination, nil
}

func (uc *AddressUseCase) GetById(id string) (*domain.Address, error) {
	address, err := uc.repo.FindById(id)
	if err != nil {
		return nil, fmt.Errorf("error fetching address by id: %v", err)
	}
	return address, nil
}

func (uc *AddressUseCase) GetBy(filters map[string]interface{}) ([]domain.Address, error) {
	address, err := uc.repo.FindBy(filters)
	if err != nil {
		return nil, fmt.Errorf("error fetching address by id: %v", err)
	}
	return address, nil
}

func (uc *AddressUseCase) Update(address domain.Address) error {
	err := uc.repo.Update(address)
	if err != nil {
		return fmt.Errorf("error on update address informations: %v", err)
	}
	return nil
}

func (uc *AddressUseCase) Delete(id string) error {
	err := uc.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("error on delete address: rv", err)
	}
	return nil
}

func (uc *AddressUseCase) Create(newAddress domain.NewAddress) (*domain.Address, error) {
	createdAddress, err := uc.repo.Create(newAddress)

	if err != nil {
		return nil, fmt.Errorf("error on create address: %v", err)
	}

	return createdAddress, nil
}
