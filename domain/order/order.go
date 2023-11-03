package order

import (
	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"gorm.io/gorm"
)

// Define the Costumer struct
type Order struct {
	gorm.Model
	UserID     uint `gorm:"not null"`
	OrderItems []OrderItem
}

// Define the OrderItem struct to represent individual books in an order
type OrderItem struct {
	gorm.Model
	OrderID  uint      `gorm:"not null"`
	Book     book.Book `gorm:"references:ID"`
	BookID   uint      `gorm:"not null"`
	Quantity int       // The number of copies of the book in the order
}
