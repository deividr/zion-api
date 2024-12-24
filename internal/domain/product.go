package domain

type Product struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Value     uint32 `json:"value"`
	UnityType string `json:"unityType"`
}

type ProductRepository interface {
	FindAll() ([]Product, error)
	FindById(id string) (*Product, error)
}
