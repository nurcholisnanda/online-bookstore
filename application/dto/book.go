package dto

type BookResponse struct {
	ID     uint   `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
}
