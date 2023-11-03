package persistence

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"github.com/nurcholisnanda/online-bookstore/domain/order"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"github.com/nurcholisnanda/online-bookstore/infrastructure/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type database struct {
	user  user.Repository
	book  book.Repository
	order order.Repository
	db    *gorm.DB
}

func DBUrl() string {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_DBName"),
	)
}

func NewDatabase() (*database, error) {
	db, err := gorm.Open(mysql.Open(DBUrl()), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return &database{
		user:  repository.NewUserRepositoryImpl(db),
		book:  repository.NewBookRepositoryImpl(db),
		order: repository.NewOrderRepositoryImpl(db),
		db:    db,
	}, nil
}

func (r *database) GetDB() *gorm.DB {
	return r.db
}

func (r *database) AutoMigrate() error {
	return r.db.AutoMigrate(&user.User{}, &book.Book{}, &order.Order{}, &order.OrderItem{})
}
