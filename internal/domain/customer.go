package domain

import "time"

type NewCustomer struct {
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Phone2    *string   `json:"phone2"`
	Email     *string   `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Customer struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Phone2    *string   `json:"phone2"`
	Email     *string   `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type FindAllCustomerFilters struct {
	Name  string
	Phone string
	Email string
}

type CustomerRepository interface {
	FindAll(Pagination, FindAllCustomerFilters) ([]Customer, Pagination, error)
	FindById(id string) (*Customer, error)
	Update(Customer) error
	Delete(id string) error
	Create(product NewCustomer) (*Customer, error)
}
