package user

import (
	"github.com/nurcholisnanda/online-bookstore/domain/order"
	"gorm.io/gorm"
)

// Define User struct
type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Orders   []order.Order
}
