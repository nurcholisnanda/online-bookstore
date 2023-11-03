package dto

type OrderRequest struct {
	OrderItems []OrderItemRequest `json:"items"`
}

type OrderItemRequest struct {
	BookID   uint `json:"book_id"`
	Quantity int  `json:"qty"`
}

type OrderHistory struct {
	ID    uint                 `json:"order_id"`
	Items []*OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	Book     BookResponse `json:"book"`
	Quantity int          `json:"qty"`
}
