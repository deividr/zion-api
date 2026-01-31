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

func (uc *AddressUseCase) GetByCustomerId(customerId string) ([]domain.Address, error) {
	addresses, err := uc.repo.FindByCustomerId(customerId)
	if err != nil {
		return nil, fmt.Errorf("error fetching addresses by customer id: %v", err)
	}
	return addresses, nil
}

func (uc *AddressUseCase) Update(customerId string, addressId string, updateData domain.NewAddress) (*domain.Address, error) {
	// Verify if the address belongs to the customer
	addresses, err := uc.repo.FindByCustomerId(customerId)
	if err != nil {
		return nil, fmt.Errorf("error on verify customer addresses: %v", err)
	}

	// Validate if the address belongs to the customer
	addressBelongsToCustomer := false
	var currentAddress *domain.Address
	for _, addr := range addresses {
		if addr.Id == addressId {
			addressBelongsToCustomer = true
			currentAddress = &addr
			break
		}
	}

	if !addressBelongsToCustomer {
		return nil, fmt.Errorf("address does not belong to this customer")
	}

	// Check if CEP or number has changed
	cepChanged := currentAddress.Cep != updateData.Cep
	numberChanged := (currentAddress.Number == nil && updateData.Number != nil) ||
		(currentAddress.Number != nil && updateData.Number == nil) ||
		(currentAddress.Number != nil && updateData.Number != nil && *currentAddress.Number != *updateData.Number)

	// If CEP or number changed, we need to delete the old address association and create a new one
	if cepChanged || numberChanged {
		// Delete the old address association
		err := uc.repo.Delete(customerId, addressId)
		if err != nil {
			return nil, fmt.Errorf("error on delete old address: %v", err)
		}

		// Create the new address association
		newAddress, err := uc.repo.Create(customerId, updateData)
		if err != nil {
			return nil, fmt.Errorf("error on create new address: %v", err)
		}

		return newAddress, nil
	}

	// Update the address (only non-identifying fields)
	addressToUpdate := domain.Address{
		Id:               addressId,
		Cep:              updateData.Cep,
		Street:           updateData.Street,
		Number:           updateData.Number,
		Neighborhood:     updateData.Neighborhood,
		City:             updateData.City,
		State:            updateData.State,
		AditionalDetails: updateData.AditionalDetails,
		Distance:         updateData.Distance,
		IsDefault:        updateData.IsDefault,
	}

	err = uc.repo.Update(addressToUpdate)
	if err != nil {
		return nil, fmt.Errorf("error on update address informations: %v", err)
	}

	// If the address is marked as default, remove default flag from other addresses
	if updateData.IsDefault != nil && *updateData.IsDefault {
		err = uc.repo.UpdateDefaultAddress(customerId, addressId)
		if err != nil {
			return nil, fmt.Errorf("error on update default address: %v", err)
		}
	}

	// Fetch and return the updated address
	updatedAddress, err := uc.repo.FindById(addressId)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated address: %v", err)
	}

	updatedAddress.IsDefault = updateData.IsDefault

	return updatedAddress, nil
}

func (uc *AddressUseCase) Delete(customerId string, addressId string) error {
	// Verify if the address belongs to the customer
	addresses, err := uc.repo.FindByCustomerId(customerId)
	if err != nil {
		return fmt.Errorf("error on verify customer addresses: %v", err)
	}

	// Validate if the address belongs to the customer
	addressBelongsToCustomer := false
	for _, addr := range addresses {
		if addr.Id == addressId {
			addressBelongsToCustomer = true
			break
		}
	}

	if !addressBelongsToCustomer {
		return fmt.Errorf("address does not belong to this customer")
	}

	if err := uc.repo.Delete(customerId, addressId); err != nil {
		return fmt.Errorf("error on delete address: %v", err)
	}

	return nil
}

func (uc *AddressUseCase) Create(customerId string, newAddress domain.NewAddress) (*domain.Address, error) {
	createdAddress, err := uc.repo.Create(customerId, newAddress)
	if err != nil {
		return nil, err
	}

	return createdAddress, nil
}
