package service

import (
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
	"github.com/nurcholisnanda/online-bookstore/domain/order"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedBookRes() []dto.BookResponse {
	books := []dto.BookResponse{
		{
			ID:     1,
			Author: "author1",
			Title:  "title1",
			Price:  15,
		},
		{
			ID:     2,
			Author: "author2",
			Title:  "title2",
			Price:  20,
		},
	}
	return books
}

func seedBook() []book.Book {
	books := []book.Book{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Author: "author1",
			Title:  "title1",
			Price:  15,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Author: "author2",
			Title:  "title2",
			Price:  20,
		},
	}
	return books
}

func seedUserReq() *dto.UserRequest {
	req := &dto.UserRequest{
		Name:     "name1",
		Email:    "email1",
		Password: "password1",
	}
	return req
}

func seedLoginReq() *dto.LoginRequest {
	req := &dto.LoginRequest{
		Email:    "email1",
		Password: "password1",
	}
	return req
}

func seedUser() *user.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}

	user := &user.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "name1",
		Email:    "email1",
		Password: string(hashedPassword),
	}
	return user
}

func seedOrderReq() *dto.OrderRequest {
	req := &dto.OrderRequest{
		OrderItems: []dto.OrderItemRequest{
			{
				BookID:   1,
				Quantity: 1,
			},
			{
				BookID:   2,
				Quantity: 2,
			},
		},
	}
	return req
}

func seedOrderHistory() []*dto.OrderHistory {
	res := []*dto.OrderHistory{
		{
			ID: 1,
			Items: []*dto.OrderItemResponse{
				{
					Book:     seedBookRes()[0],
					Quantity: 1,
				},
			},
		},
		{
			ID: 2,
			Items: []*dto.OrderItemResponse{
				{
					Book:     seedBookRes()[1],
					Quantity: 2,
				},
			},
		},
	}
	return res
}

func seedOrders() []order.Order {
	orders := []order.Order{
		{
			Model: gorm.Model{
				ID: 1,
			},
			UserID: 1,
			OrderItems: []order.OrderItem{
				{
					Model: gorm.Model{
						ID: 1,
					},
					OrderID:  1,
					Book:     seedBook()[0],
					BookID:   1,
					Quantity: 1,
				},
			},
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			UserID: 1,
			OrderItems: []order.OrderItem{
				{
					Model: gorm.Model{
						ID: 2,
					},
					OrderID:  1,
					Book:     seedBook()[1],
					BookID:   2,
					Quantity: 2,
				},
			},
		},
	}
	return orders
}
