package domain

type CategoryProduct struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryProductRepository interface {
	FindAll() ([]CategoryProduct, error)
	FindById(id string) (*CategoryProduct, error)
	Update(CategoryProduct) error
	Delete(id string) error
	Create(CategoryProduct) (*CategoryProduct, error)
}
