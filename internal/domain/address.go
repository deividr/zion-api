package domain

type NewAddress struct {
	CustomerId       string  `json:"customerId"`
	Cep              string  `json:"cep"`
	Street           *string `json:"street"`
	Number           *string `json:"number"`
	Neighborhood     *string `json:"neighborhood"`
	City             *string `json:"city"`
	State            *string `json:"state"`
	AditionalDetails *string `json:"aditionalDetails"`
	Distance         *string `json:"distance"`
}

type Address struct {
	Id               string  `json:"id"`
	OldId            *string `json:"oldId"`
	CustomerId       string  `json:"customerId"`
	Cep              string  `json:"cep"`
	Street           *string `json:"street"`
	Number           *string `json:"number"`
	Neighborhood     *string `json:"neighborhood"`
	City             *string `json:"city"`
	State            *string `json:"state"`
	AditionalDetails *string `json:"aditionalDetails"`
	Distance         *string `json:"distance"`
}

type AddressRepository interface {
	FindAll(Pagination) ([]Address, Pagination, error)
	FindById(id string) (*Address, error)
	FindBy(filters map[string]interface{}) ([]Address, error)
	Update(Address) error
	Delete(id string) error
	Create(product NewAddress) (*Address, error)
}
