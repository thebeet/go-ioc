package order

type Service interface {
	GetOrderById(id int) *Order
}
