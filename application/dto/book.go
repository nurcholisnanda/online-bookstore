package dto

//Define BookResponse struct that will be used for http response body
type BookResponse struct {
	ID     uint   `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Price  int    `json:"price"`
}
