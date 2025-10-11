package domain

type NewAddress struct {
	Cep              string  `json:"cep"`
	Street           *string `json:"street"`
	Number           *string `json:"number"`
	Neighborhood     *string `json:"neighborhood"`
	City             *string `json:"city"`
	State            *string `json:"state"`
	AditionalDetails *string `json:"aditionalDetails"`
	Distance         *int    `json:"distance"`
	IsDefault        *bool   `json:"isDefault"`
}

type Address struct {
	Id               string  `json:"id"`
	OldId            *string `json:"oldId"`
	Cep              string  `json:"cep"`
	Street           *string `json:"street"`
	Number           *string `json:"number"`
	Neighborhood     *string `json:"neighborhood"`
	City             *string `json:"city"`
	State            *string `json:"state"`
	AditionalDetails *string `json:"aditionalDetails"`
	Distance         *int    `json:"distance"`
	IsDefault        *bool   `json:"isDefault"`
}

type AddressRepository interface {
	FindAll(Pagination) ([]Address, Pagination, error)
	FindById(id string) (*Address, error)
	FindBy(filters map[string]any) ([]Address, error)
	FindByCustomerId(customerId string) ([]Address, error)
	Update(Address) error
	Delete(id string) error
	Create(product NewAddress) (*Address, error)
}
