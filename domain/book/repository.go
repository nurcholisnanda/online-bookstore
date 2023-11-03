package book

// Interface that contains repository contract for book domain
//go:generate mockgen -source=repository.go -destination=mock/book.go -package=mock
type Repository interface {
	FindAllBook() ([]Book, error)
	GetBooksByIDs([]uint) ([]Book, error)
}
