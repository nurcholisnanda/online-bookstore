package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"github.com/nurcholisnanda/online-bookstore/domain/order"
)

var (
	errOrderItemsRequired    = errors.New("order items is required")
	errAtLeastOneQtyRequired = errors.New("at least one qty is required")
)

type orderService struct {
	bookRepo  book.Repository
	orderRepo order.Repository
}

//go:generate mockgen -source=order.go -destination=mock/order.go -package=mock
type OrderService interface {
	MakeOrder(*dto.OrderRequest, uint) error
	GetOrderHistory(uint) ([]*dto.OrderHistory, error)
}

// Order service constructor
func NewOrderService(bookRepo book.Repository, orderRepo order.Repository) OrderService {
	return &orderService{
		bookRepo:  bookRepo,
		orderRepo: orderRepo,
	}
}

// MakeOrder service will check whether book in order request is exist.
// Moreover, will call CreateOrder repository contract.
func (s *orderService) MakeOrder(req *dto.OrderRequest, userID uint) error {
	orderReq := &order.Order{
		UserID: userID,
	}

	var orderItems []*order.OrderItem
	bookIds := make([]uint, 0)
	mapBookIds := make(map[uint]struct{})

	if len(req.OrderItems) < 1 {
		return errOrderItemsRequired
	}

	// Iteration to get bookIds and parse into orderItems
	for _, item := range req.OrderItems {
		if item.Quantity < 1 {
			return errAtLeastOneQtyRequired
		}

		bookIds = append(bookIds, item.BookID)
		orderItems = append(orderItems, &order.OrderItem{
			BookID:   item.BookID,
			Quantity: item.Quantity,
		})
	}

	//check whether the books in order request exists in our database
	books, err := s.bookRepo.GetBooksByIDs(bookIds)
	if err != nil {
		return err
	}

	//check whether there are any Ids not exist in database
	if len(books) != len(orderItems) {
		notFoundIds := make([]string, 0)
		for _, book := range books {
			mapBookIds[book.ID] = struct{}{}
		}
		for _, item := range orderItems {
			if _, ok := mapBookIds[item.BookID]; !ok {
				notFoundIds = append(notFoundIds, strconv.Itoa(int(item.BookID)))
			}
		}
		return fmt.Errorf("books with ids: [%s] not found", strings.Join(notFoundIds, ","))
	}

	if err := s.orderRepo.CreateOrder(orderReq, orderItems); err != nil {
		return err
	}

	return nil
}

// GetOrderHistory service will call FindUserOrderHistory in our repository
// contract and return the response as OrderHistory dto to be controlled
// in our order controller
func (s *orderService) GetOrderHistory(userID uint) (res []*dto.OrderHistory, err error) {
	orders, err := s.orderRepo.FindUserOrderHistory(userID)
	if err != nil {
		return nil, err
	}
	for i, order := range orders {
		res = append(res, &dto.OrderHistory{
			ID: order.ID,
		})
		for _, item := range order.OrderItems {
			res[i].Items = append(res[i].Items, &dto.OrderItemResponse{
				Book: dto.BookResponse{
					ID:     item.BookID,
					Author: item.Book.Author,
					Title:  item.Book.Title,
					Price:  item.Book.Price,
				},
				Quantity: item.Quantity,
			})
		}
	}
	return res, err
}
