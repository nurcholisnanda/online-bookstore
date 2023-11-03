package repository

import (
	"reflect"
	"testing"

	"github.com/nurcholisnanda/online-bookstore/domain/order"
	"gorm.io/gorm"
)

func TestNewOrderRepositoryImpl(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	tests := []struct {
		name string
		args args
		want order.Repository
	}{
		{
			name: "implemented",
			args: args{
				db: db,
			},
			want: NewOrderRepositoryImpl(db),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderRepositoryImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderRepositoryImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_orderRepositoryImpl_CreateOrder(t *testing.T) {
	type args struct {
		order     *order.Order
		orderItem []*order.OrderItem
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	book := seedBook(db)
	user, _ := seedUser(db)

	tests := []struct {
		name    string
		r       *orderRepositoryImpl
		args    args
		wantErr bool
	}{
		{
			name: "success create order",
			r: &orderRepositoryImpl{
				db: db,
			},
			args: args{
				order: &order.Order{
					Model:  gorm.Model{},
					UserID: user.ID,
				},
				orderItem: []*order.OrderItem{
					{
						Model:    gorm.Model{},
						BookID:   book[0].Model.ID,
						Quantity: 1,
					},
					{
						Model:    gorm.Model{},
						BookID:   book[1].Model.ID,
						Quantity: 2,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error create order",
			r: &orderRepositoryImpl{
				db: db,
			},
			args: args{
				order:     &order.Order{},
				orderItem: []*order.OrderItem{},
			},
			wantErr: true,
		},
		{
			name: "error create order item",
			r: &orderRepositoryImpl{
				db: db,
			},
			args: args{
				order: &order.Order{
					Model:  gorm.Model{},
					UserID: user.ID,
				},
				orderItem: []*order.OrderItem{
					{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.CreateOrder(tt.args.order, tt.args.orderItem); (err != nil) != tt.wantErr {
				t.Errorf("orderRepositoryImpl.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_orderRepositoryImpl_FindUserOrderHistory(t *testing.T) {
	type args struct {
		userID uint
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}
	user, _ := seedUser(db)
	orders, _ := seedOrder(db)

	tests := []struct {
		name             string
		r                *orderRepositoryImpl
		args             args
		wantOrderHistory []order.Order
		wantErr          bool
	}{
		{
			name: "success find order history",
			r: &orderRepositoryImpl{
				db: db,
			},
			args: args{
				userID: user.ID,
			},
			wantOrderHistory: orders,
			wantErr:          false,
		},
		{
			name: "error fetching order history",
			r: &orderRepositoryImpl{
				db: db,
			},
			args: args{
				userID: 0,
			},
			wantOrderHistory: nil,
			wantErr:          true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOrderHistory, err := tt.r.FindUserOrderHistory(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("orderRepositoryImpl.FindUserOrderHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOrderHistory, tt.wantOrderHistory) {
				t.Errorf("orderRepositoryImpl.FindUserOrderHistory() = %v, want %v", gotOrderHistory, tt.wantOrderHistory)
			}
		})
	}
}
