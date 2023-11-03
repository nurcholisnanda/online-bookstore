package repository

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"github.com/nurcholisnanda/online-bookstore/domain/order"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DBConn() (*gorm.DB, error) {
	if _, err := os.Stat("./../../.env"); !os.IsNotExist(err) {
		err := godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
	}
	return LocalDatabase()
}

func LocalDatabase() (*gorm.DB, error) {
	host := os.Getenv("TEST_DB_HOST")
	password := os.Getenv("TEST_DB_PASSWORD")
	username := os.Getenv("TEST_DB_USER")
	dbname := os.Getenv("TEST_DB_DBName")
	port, err := strconv.Atoi(os.Getenv("TEST_DB_PORT"))
	if err != nil {
		panic(err)
	}

	DBURL := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		dbname,
	)

	db, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Migrator().DropTable(&book.Book{}, &user.User{}, &order.Order{}, order.OrderItem{})
	db.AutoMigrate(&book.Book{}, &user.User{}, &order.Order{}, order.OrderItem{})

	return db, nil
}

func seedUser(db *gorm.DB) (*user.User, error) {
	user := &user.User{
		Model:    gorm.Model{},
		Name:     "name1",
		Email:    "email1",
		Password: "password1",
	}
	err := db.Create(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func seedBook(db *gorm.DB) []book.Book {
	bookReq := []*book.Book{
		{
			Model:  gorm.Model{},
			Author: "author1",
			Title:  "title1",
			Price:  15,
		},
		{
			Model:  gorm.Model{},
			Author: "author2",
			Title:  "title2",
			Price:  20,
		},
	}
	var books []book.Book
	for _, v := range bookReq {
		err := db.Create(v).Error
		if err != nil {
			return nil
		}
		books = append(books, *v)
	}

	return books
}

func seedOrder(db *gorm.DB) ([]order.Order, error) {
	ord := &order.Order{
		Model:  gorm.Model{},
		UserID: 1,
	}
	err := db.Create(ord).Error
	if err != nil {
		return nil, err
	}
	books := seedBook(db)

	orderItems := []*order.OrderItem{
		{
			Model:    gorm.Model{},
			OrderID:  ord.ID,
			Book:     books[0],
			BookID:   1,
			Quantity: 1,
		},
		{
			Model:    gorm.Model{},
			OrderID:  ord.ID,
			Book:     books[1],
			BookID:   2,
			Quantity: 2,
		},
	}

	for _, v := range orderItems {
		err := db.Create(v).Error
		if err != nil {
			return nil, err
		}
		ord.OrderItems = append(ord.OrderItems, *v)
	}

	orders := []order.Order{
		*ord,
	}

	return orders, nil
}
