package dto

//Define OrderRequest struct that will be used for http request body
type OrderRequest struct {
	OrderItems []OrderItemRequest `json:"items"`
}

// Define order item to represent individual books in an order request
type OrderItemRequest struct {
	BookID   uint `json:"book_id"`
	Quantity int  `json:"qty"`
}

//Define OrderHistory struct that will be used for http response body
type OrderHistory struct {
	ID    uint                 `json:"order_id"`
	Items []*OrderItemResponse `json:"items"`
}

// Define order item to represent individual books in an order response
type OrderItemResponse struct {
	Book     BookResponse `json:"book"`
	Quantity int          `json:"qty"`
}
