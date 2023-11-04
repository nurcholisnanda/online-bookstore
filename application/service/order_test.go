package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/book"
	mock_book "github.com/nurcholisnanda/online-bookstore/domain/book/mock"
	"github.com/nurcholisnanda/online-bookstore/domain/order"
	mock_order "github.com/nurcholisnanda/online-bookstore/domain/order/mock"
)

func TestNewOrderService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBook := mock_book.NewMockRepository(ctrl)
	mockOrder := mock_order.NewMockRepository(ctrl)

	type args struct {
		bookRepo  book.Repository
		orderRepo order.Repository
	}
	tests := []struct {
		name string
		args args
		want OrderService
	}{
		{
			name: "implemented",
			args: args{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			want: NewOrderService(mockBook, mockOrder),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderService(tt.args.bookRepo, tt.args.orderRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderService_MakeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBook := mock_book.NewMockRepository(ctrl)
	mockOrder := mock_order.NewMockRepository(ctrl)
	getBooksByIDsRecord := func(books []book.Book, err error) func(m *mock_book.MockRepository) {
		return func(m *mock_book.MockRepository) {
			m.EXPECT().GetBooksByIDs(gomock.Any()).Return(books, err)
		}
	}
	createOrderRecord := func(err error) func(m *mock_order.MockRepository) {
		return func(m *mock_order.MockRepository) {
			m.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	books := seedBook()

	type args struct {
		req    *dto.OrderRequest
		userID uint
	}

	orderReq := seedOrderReq()
	tests := []struct {
		name    string
		s       *orderService
		args    args
		wantErr bool
	}{
		{
			name: "success making order",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				req:    orderReq,
				userID: 1,
			},
			wantErr: false,
		},
		{
			name: "fail get book",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				req: &dto.OrderRequest{
					OrderItems: []dto.OrderItemRequest{
						{
							BookID:   3,
							Quantity: 3,
						},
					},
				},
				userID: 1,
			},
			wantErr: true,
		},
		{
			name: "success get book fail create order",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				req:    orderReq,
				userID: 1,
			},
			wantErr: true,
		},
		{
			name: "error no order item",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				req: &dto.OrderRequest{},
			},
			wantErr: true,
		},
		{
			name: "error no qty order",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				req: &dto.OrderRequest{
					OrderItems: []dto.OrderItemRequest{
						{
							BookID:   1,
							Quantity: 0,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "error ids not exist",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				req: &dto.OrderRequest{
					OrderItems: []dto.OrderItemRequest{
						{
							BookID:   1,
							Quantity: 1,
						},
						{
							BookID:   20,
							Quantity: 1,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "success making order":
				getBooksByIDsRecord(books, nil)(mockBook)
				createOrderRecord(nil)(mockOrder)
			case "fail get book":
				getBooksByIDsRecord(nil, errors.New("any error"))(mockBook)
			case "success get book fail create order":
				getBooksByIDsRecord(books, nil)(mockBook)
				createOrderRecord(errors.New("any error"))(mockOrder)
			case "error ids not exist":
				getBooksByIDsRecord(books, nil)(mockBook)
			}
			if err := tt.s.MakeOrder(tt.args.req, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("orderService.MakeOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_orderService_GetOrderHistory(t *testing.T) {
	type args struct {
		userID uint
	}
	ctrl := gomock.NewController(t)
	mockBook := mock_book.NewMockRepository(ctrl)
	mockOrder := mock_order.NewMockRepository(ctrl)
	findUserOrderHistoryRecord := func(orders []order.Order, err error) func(m *mock_order.MockRepository) {
		return func(m *mock_order.MockRepository) {
			m.EXPECT().FindUserOrderHistory(gomock.Any()).Return(orders, err)
		}
	}
	userID := 1
	orders := seedOrders()
	res := seedOrderHistory()

	tests := []struct {
		name    string
		s       *orderService
		args    args
		wantRes []*dto.OrderHistory
		wantErr bool
	}{
		{
			name: "success get order history",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				userID: uint(userID),
			},
			wantRes: res,
			wantErr: false,
		},
		{
			name: "fail get order history",
			s: &orderService{
				bookRepo:  mockBook,
				orderRepo: mockOrder,
			},
			args: args{
				userID: uint(userID),
			},
			wantRes: res,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "success get order history":
				findUserOrderHistoryRecord(orders, nil)(mockOrder)
			case "fail get order history":
				findUserOrderHistoryRecord(nil, errors.New("any error"))(mockOrder)
			}
			gotRes, err := tt.s.GetOrderHistory(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderService.GetOrderHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("orderService.GetOrderHistory() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
