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

// Book repository implementation constructor
func NewBookRepositoryImpl(db *gorm.DB) book.Repository {
	return &bookRepositoryImpl{
		db: db,
	}
}

// Implementation of FindAllBook repo contract to find all book
// data in our database
func (r *bookRepositoryImpl) FindAllBook() ([]book.Book, error) {
	var books []book.Book
	if err := r.db.Find(&books).Error; err != nil || len(books) == 0 {
		return nil, errBooksNotFound
	}
	return books, nil
}

// Implementation of GetBooksByIDs repo contract to find books data
// where id is in slice of ids
func (r *bookRepositoryImpl) GetBooksByIDs(ids []uint) ([]book.Book, error) {
	var books []book.Book
	if err := r.db.Where("id in ?", ids).Find(&books).
		Error; err != nil || len(books) == 0 {
		return nil, errBooksNotFound
	}
	return books, nil
}
