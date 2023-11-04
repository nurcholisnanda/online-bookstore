package repository

import (
	"errors"

	"github.com/nurcholisnanda/online-bookstore/domain/order"
	"gorm.io/gorm"
)

var errOrderNotFound = errors.New("order records not found")

type orderRepositoryImpl struct {
	db *gorm.DB
}

// Order repository implementation construction
func NewOrderRepositoryImpl(db *gorm.DB) order.Repository {
	return &orderRepositoryImpl{
		db: db,
	}
}

// Implementation of CreateOrder repo contract to insert
// new order and order item to our database
func (r *orderRepositoryImpl) CreateOrder(order *order.Order, orderItem []*order.OrderItem) error {
	//use transaction to make sure the transaction data is safed (rollback when error, commit when all finished)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range orderItem {
			item.OrderID = order.ID
			if err := tx.Create(item).Error; err != nil {
				return err
			}
		}
		// return nil will commit the whole transaction
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// Implementation of FindUserOrderHistory repo contract to find
// order history from a user by user_id from database
func (r *orderRepositoryImpl) FindUserOrderHistory(userID uint) ([]order.Order, error) {
	var orderHistory []order.Order
	//use transaction to make sure the transaction data is safed (rollback when error, commit when all finished)
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("OrderItems").Preload("OrderItems.Book").
			Where("user_id = ?", userID).
			Find(&orderHistory).Error; err != nil || len(orderHistory) == 0 {
			return errOrderNotFound
		}
		// return nil will commit the whole transaction
		return nil
	})

	if err != nil {
		return nil, err
	}

	return orderHistory, nil
}
