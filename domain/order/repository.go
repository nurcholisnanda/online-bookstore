package order

// Interface that contains repository contract for order domain
//go:generate mockgen -source=repository.go -destination=mock/order.go -package=mock
type Repository interface {
	CreateOrder(*Order, []*OrderItem) error
	FindUserOrderHistory(uint) ([]Order, error)
}
