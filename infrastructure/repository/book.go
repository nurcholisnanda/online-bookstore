package repository

import (
	"errors"

	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"gorm.io/gorm"
)

var (
	errBooksNotFound = errors.New("book records not found")
)

type bookRepositoryImpl struct {
	db *gorm.DB
}

func NewBookRepositoryImpl(db *gorm.DB) book.Repository {
	return &bookRepositoryImpl{
		db: db,
	}
}

func (r *bookRepositoryImpl) FindAllBook() ([]book.Book, error) {
	var books []book.Book
	if err := r.db.Find(&books).Error; err != nil || len(books) == 0 {
		return nil, errBooksNotFound
	}
	return books, nil
}

func (r *bookRepositoryImpl) GetBooksByIDs(ids []uint) ([]book.Book, error) {
	var books []book.Book
	if err := r.db.Where("id in ?", ids).Find(&books).
		Error; err != nil || len(books) == 0 {
		return nil, errBooksNotFound
	}
	return books, nil
}
