package domain

import "fmt"

type DuplicateAddressError struct {
	CustomerID string
	Cep        string
	Number     string
}

func (e *DuplicateAddressError) Error() string {
	return fmt.Sprintf("address with cep %s and number %s already exists for customer %s", e.Cep, e.Number, e.CustomerID)
}

func NewDuplicateAddressError(customerID, cep, number string) *DuplicateAddressError {
	return &DuplicateAddressError{
		CustomerID: customerID,
		Cep:        cep,
		Number:     number,
	}
}
