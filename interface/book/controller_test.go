package book

import (
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
		service service.BookService
	}
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockBookService(ctrl)

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

func TestController_GetBooks(t *testing.T) {
	type args struct {
		ctx *gin.Context
	}
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockBookService(ctrl)
	mockGetBooks := func(res []dto.BookResponse, err error) func(m *mock.MockBookService) {
		return func(m *mock.MockBookService) {
			m.EXPECT().GetBooks().Return(res, err)
		}
	}

	tests := []struct {
		name   string
		c      *Controller
		args   args
		status int
	}{
		{
			name: "send success response",
			c: &Controller{
				service: mockSvc,
			},
			args: args{
				ctx: &gin.Context{},
			},
			status: http.StatusOK,
		},
		{
			name: "send error response",
			c: &Controller{
				service: mockSvc,
			},
			args: args{
				ctx: &gin.Context{},
			},
			status: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
			tt.args.ctx.Request, _ = http.NewRequest("GET", "/v1/books", nil)

			switch tt.name {
			case "send success response":
				mockGetBooks([]dto.BookResponse{}, nil)(mockSvc)
			case "send error response":
				mockGetBooks(nil, errors.New("any error"))(mockSvc)
			}
			tt.c.GetBooks(tt.args.ctx)
			if ok := reflect.DeepEqual(tt.args.ctx.Writer.Status(), tt.status); !ok {
				t.Errorf("http status = %v, want %v", tt.args.ctx.Writer.Status(), tt.status)
			}
		})
	}
}
