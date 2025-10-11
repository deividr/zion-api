package domain

import "time"

type NewOrder struct {
	PickupDate   time.Time      `json:"pickupDate"`
	Customer     Customer       `json:"customer"`
	Address      *Address       `json:"address"`
	Employee     string         `json:"employee"`
	OrderLocal   *string        `json:"orderLocal"`
	Observations *string        `json:"observations"`
	IsPickedUp   *bool          `json:"isPickedUp"`
	Products     []OrderProduct `json:"products"`
}

type Order struct {
	NewOrder
	Id        string     `json:"id"`
	Number    string     `json:"number"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type OrderProduct struct {
	Id          string            `json:"id"`
	OrderId     string            `json:"orderId"`
	ProductId   string            `json:"productId"`
	Quantity    int               `json:"quantity"`
	UnityType   string            `json:"unityType"`
	Price       int               `json:"price"`
	Name        string            `json:"name"`
	SubProducts []OrderSubProduct `json:"subProducts"`
}

type OrderSubProduct struct {
	Id             string `json:"id"`
	OrderProductId string `json:"orderProductId"`
	ProductId      string `json:"productId"`
	Name           string `json:"name"`
}

type FindAllOrderFilters struct {
	PickupDateStart time.Time
	PickupDateEnd   time.Time
	Search          *string
}

type OrderRepository interface {
	FindAll(Pagination, FindAllOrderFilters) ([]Order, Pagination, error)
	FindById(id string) (*Order, error)
	Update(Order) error
	Delete(id string) error
	Create(order NewOrder) (*Order, error)
}
