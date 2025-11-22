package domain

type NewProduct struct {
	Name       string `json:"name"`
	Value      uint32 `json:"value"`
	UnityType  string `json:"unityType"`
	CategoryId string `json:"categoryId"`
}

type Product struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Value      uint32 `json:"value"`
	UnityType  string `json:"unityType"`
	CategoryId string `json:"categoryId"`
}

type FindAllProductFilters struct {
	Name       string
	UnityType  string
	CategoryId string
}

type ProductRepository interface {
	FindAll(FindAllProductFilters) ([]Product, error)
	FindById(id string) (*Product, error)
	Update(Product) error
	Delete(id string) error
	Create(product NewProduct) (*Product, error)
}
