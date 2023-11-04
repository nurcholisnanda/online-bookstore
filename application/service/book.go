package service

import (
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
)

type bookService struct {
	bookRepo book.Repository
}

//go:generate mockgen -source=book.go -destination=mock/book.go -package=mock
type BookService interface {
	GetBooks() ([]dto.BookResponse, error)
}

// Book service constructor
func NewBookService(bookRepo book.Repository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

// GetBooks service will call FindAllBook in our repository contract
// and return the response as BookResponse dto to be controlled
// in our book controller
func (s *bookService) GetBooks() ([]dto.BookResponse, error) {
	books, err := s.bookRepo.FindAllBook()
	if err != nil {
		return nil, err
	}
	var res []dto.BookResponse
	//iterate books data and map it into BookResponse dto
	for _, book := range books {
		res = append(res, dto.BookResponse{
			ID:     book.ID,
			Author: book.Author,
			Title:  book.Title,
			Price:  book.Price,
		})
	}
	return res, nil
}
