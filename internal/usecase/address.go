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

func (uc *AddressUseCase) GetBy(filters map[string]any) ([]domain.Address, error) {
	address, err := uc.repo.FindBy(filters)
	if err != nil {
		return nil, fmt.Errorf("error fetching address by id: %v", err)
	}
	return address, nil
}

func (uc *AddressUseCase) Update(address domain.Address) error {
	addressDb, errFind := uc.repo.FindBy(map[string]any{"customer_id": address.CustomerId, "is_default": true, "is_deleted": false})
	if errFind != nil {
		return fmt.Errorf("error fetching address by id: %v", errFind)
	}

	if len(addressDb) > 0 && addressDb[0].Id != address.Id {
		addressDb[0].IsDefault = &[]bool{false}[0]
		if err := uc.repo.Update(addressDb[0]); err != nil {
			return fmt.Errorf("error update old default address to not default: %v", errFind)
		}
	}

	err := uc.repo.Update(address)
	if err != nil {
		return fmt.Errorf("error on update address informations: %v", err)
	}

	return nil
}

func (uc *AddressUseCase) Delete(id string) error {
	addressDb, errFind := uc.repo.FindById(id)
	if errFind != nil {
		return fmt.Errorf("error fetching address by id: %v", errFind)
	}

	if addressDb.IsDefault != nil && *addressDb.IsDefault {
		addressesDb, err := uc.repo.FindBy(map[string]any{"customer_id": addressDb.CustomerId, "is_deleted": false})
		if err != nil {
			return fmt.Errorf("error fetching addresses: %v", err)
		}
		if len(addressesDb) > 0 {
			for _, address := range addressesDb {
				if address.Id == id {
					continue
				}

				address.IsDefault = &[]bool{true}[0]
				if err := uc.repo.Update(address); err != nil {
					return fmt.Errorf("error update address to default: %v", err)
				}

				break
			}
		}
	}

	if err := uc.repo.Delete(id); err != nil {
		return fmt.Errorf("error on delete address: %v", err)
	}

	return nil
}

func (uc *AddressUseCase) Create(newAddress domain.NewAddress) (*domain.Address, error) {
	if newAddress.IsDefault != nil && *newAddress.IsDefault {
		addressDb, errFind := uc.repo.FindBy(map[string]any{"customer_id": newAddress.CustomerId, "is_default": true, "is_deleted": false})
		if errFind != nil {
			return nil, fmt.Errorf("error fetching address by id: %v", errFind)
		}

		if len(addressDb) > 0 {
			addressDb[0].IsDefault = &[]bool{false}[0]
			if err := uc.repo.Update(addressDb[0]); err != nil {
				return nil, fmt.Errorf("error update old default address to not default: %v", errFind)
			}
		}
	}

	createdAddress, err := uc.repo.Create(newAddress)
	if err != nil {
		return nil, fmt.Errorf("error on create address: %v", err)
	}

	return createdAddress, nil
}
