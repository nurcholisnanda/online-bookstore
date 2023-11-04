package order

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/application/service"
	"github.com/nurcholisnanda/online-bookstore/application/service/mock"
)

func TestNewController(t *testing.T) {
	type args struct {
		service service.OrderService
	}
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockOrderService(ctrl)

	tests := []struct {
		name string
		args args
		want *Controller
	}{
		{
			name: "implemented",
			args: args{
				service: mockSvc,
			},
			want: NewController(mockSvc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewController(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_MakeOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockOrderService(ctrl)
	mockMakeOrder := func(err error) func(m *mock.MockOrderService) {
		return func(m *mock.MockOrderService) {
			m.EXPECT().MakeOrder(gomock.Any(), gomock.Any()).Return(err)
		}
	}
	ctr := NewController(mockSvc)
	var userID uint = 1

	tests := []struct {
		name   string
		req    string
		status int
	}{
		{
			name:   "success making order",
			req:    `{"items":[{"book_id":1,"qty":1},{"book_id":2,"qty":2}]}`,
			status: http.StatusOK,
		},
		{
			name:   "fail user_id key not set",
			status: http.StatusInternalServerError,
		},
		{
			name:   "fail binding JSON",
			req:    `{"items":[{"book_id":1,"qty":1},{"book_id":2,"qty":2]}`,
			status: http.StatusBadRequest,
		},
		{
			name:   "fail making order",
			req:    `{"items":[{"book_id":1,"qty":1},{"book_id":2,"qty":2}]}`,
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request, _ = http.NewRequest("POST", "/v1/orders",
				bytes.NewBuffer([]byte(tt.req)),
			)
			switch tt.name {
			case "success making order":
				ctx.Set("user_id", userID)
				mockMakeOrder(nil)(mockSvc)
			case "fail binding JSON":
				ctx.Set("user_id", userID)
			case "fail making order":
				ctx.Set("user_id", userID)
				mockMakeOrder(errors.New("any error"))(mockSvc)
			}

			ctr.MakeOrder(ctx)
			if ok := reflect.DeepEqual(ctx.Writer.Status(), tt.status); !ok {
				t.Errorf("http status = %v, want %v", ctx.Writer.Status(), tt.status)
			}
		})
	}
}

func TestController_GetOrderHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockOrderService(ctrl)
	mockGetHistory := func(orders []*dto.OrderHistory, err error) func(m *mock.MockOrderService) {
		return func(m *mock.MockOrderService) {
			m.EXPECT().GetOrderHistory(gomock.Any()).Return(orders, err)
		}
	}
	ctr := NewController(mockSvc)
	var userID uint = 1

	type args struct {
		ctx *gin.Context
	}
	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "success making order",
			status: http.StatusOK,
		},
		{
			name:   "fail user_id key not set",
			status: http.StatusInternalServerError,
		},
		{
			name:   "fail making order",
			status: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request, _ = http.NewRequest("GET", "/v1/orders/history", nil)
			switch tt.name {
			case "success making order":
				ctx.Set("user_id", userID)
				mockGetHistory([]*dto.OrderHistory{}, nil)(mockSvc)
			case "fail making order":
				ctx.Set("user_id", userID)
				mockGetHistory(nil, errors.New("any error"))(mockSvc)
			}

			ctr.GetOrderHistory(ctx)
			if ok := reflect.DeepEqual(ctx.Writer.Status(), tt.status); !ok {
				t.Errorf("http status = %v, want %v", ctx.Writer.Status(), tt.status)
			}
		})
	}
}
