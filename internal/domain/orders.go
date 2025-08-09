package domain

import "time"

type NewOrder struct {
	Number       string     `json:"number"`
	PickupDate   time.Time  `json:"pickupDate"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	CustomerId   string     `json:"customerId"`
	EmployeeId   string     `json:"employeeId"`
	OrderLocal   *string    `json:"orderLocal"`
	Observations *string    `json:"observations"`
	IsPickedUp   *bool      `json:"isPickedUp"`
}

type Order struct {
	Id           string     `json:"id"`
	Number       string     `json:"number"`
	PickupDate   time.Time  `json:"withdraw"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	CustomerId   string     `json:"customerId"`
	EmployeeId   string     `json:"employeeId"`
	OrderLocal   *string    `json:"orderLocal"`
	Observations *string    `json:"observations"`
	IsPickedUp   *bool      `json:"isPickedUp"`
}

type OrderProduct struct {
	Id        string `json:"id"`
	OrderId   string `json:"orderId"`
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
	UnityType string `json:"unityType"`
	Price     int    `json:"price"`
}

type FullOrderProduct struct {
	OrderProduct OrderProduct
	SubProducts  []OrderSubProduct
}

type OrderSubProduct struct {
	Id             string `json:"id"`
	OrderProductId string `json:"orderProductId"`
	ProductId      string `json:"productId"`
}

type FindAllOrderFilters struct {
	PickupDate *time.Time
	CustomerId string
	ProductId  string
}

type OrderRepository interface {
	FindAll(Pagination, FindAllOrderFilters) ([]Order, Pagination, error)
	FindById(id string) (*Order, error)
	Update(Order) error
	Delete(id string) error
	Create(order NewOrder) (*Order, error)
}
