package service

import (
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
)

type bookService struct {
	bookRepo book.Repository
}

type BookService interface {
	GetBooks() ([]dto.BookResponse, error)
}

func NewBookService(bookRepo book.Repository) BookService {
	return &bookService{
		bookRepo: bookRepo,
	}
}

func (s *bookService) GetBooks() ([]dto.BookResponse, error) {
	books, err := s.bookRepo.FindAllBook()
	if err != nil {
		return nil, err
	}
	var res []dto.BookResponse
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
