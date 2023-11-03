package book

import "gorm.io/gorm"

// Define the Costumer entity
type Book struct {
	gorm.Model
	Author string `gorm:"not null"`
	Title  string `gorm:"not null"`
	Price  int
}
