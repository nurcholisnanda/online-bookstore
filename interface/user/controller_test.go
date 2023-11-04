package user

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/application/service"
	"github.com/nurcholisnanda/online-bookstore/application/service/mock"
)

func TestNewController(t *testing.T) {
	type args struct {
		service service.UserService
	}
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockUserService(ctrl)

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

func TestController_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockUserService(ctrl)
	registerMock := func(err error) func(m *mock.MockUserService) {
		return func(m *mock.MockUserService) {
			m.EXPECT().Register(gomock.Any()).Return(err)
		}
	}
	ctr := &Controller{
		service: mockSvc,
	}
	var jsonParam string

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "success register",
			status: http.StatusOK,
		},
		{
			name:   "error validation",
			status: http.StatusBadRequest,
		},
		{
			name:   "error binding json",
			status: http.StatusBadRequest,
		},
		{
			name:   "error register service",
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "success register":
				jsonParam = `{"name":"name","email":"email","password":"password"}`
				registerMock(nil)(mockSvc)
			case "error validation":
				jsonParam = `{"name":"","email":"","password":""}`
			case "error binding json":
				jsonParam = `{"name":"","email":"","password":"",}`
			case "error register service":
				jsonParam = `{"name":"name","email":"email","password":"password"}`
				registerMock(errors.New("any error"))(mockSvc)
			}

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request, _ = http.NewRequest("POST", "/v1/users/register",
				bytes.NewBuffer([]byte(jsonParam)),
			)
			ctr.Register(ctx)
			if ok := reflect.DeepEqual(ctx.Writer.Status(), tt.status); !ok {
				t.Errorf("http status = %v, want %v", ctx.Writer.Status(), tt.status)
			}
		})
	}
}

func TestController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSvc := mock.NewMockUserService(ctrl)
	loginMock := func(token string, err error) func(m *mock.MockUserService) {
		return func(m *mock.MockUserService) {
			m.EXPECT().Login(gomock.Any()).Return(token, err)
		}
	}
	ctr := &Controller{
		service: mockSvc,
	}
	token := "validtoken"
	var jsonParam string

	tests := []struct {
		name   string
		status int
	}{
		{
			name:   "success login",
			status: http.StatusOK,
		},
		{
			name:   "error validation",
			status: http.StatusBadRequest,
		},
		{
			name:   "error binding json",
			status: http.StatusBadRequest,
		},
		{
			name:   "error login service",
			status: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "success login":
				jsonParam = `{"email":"email","password":"password"}`
				loginMock(token, nil)(mockSvc)
			case "error validation":
				jsonParam = `{"email":"","password":""}`
			case "error binding json":
				jsonParam = `{"email":"","password":"",}`
			case "error login service":
				jsonParam = `{"email":"email","password":"password"}`
				loginMock("", errors.New("any error"))(mockSvc)
			}

			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Request, _ = http.NewRequest("POST", "/v1/users/login",
				bytes.NewBuffer([]byte(jsonParam)),
			)
			ctr.Login(ctx)
			if ok := reflect.DeepEqual(ctx.Writer.Status(), tt.status); !ok {
				t.Errorf("http status = %v, want %v", ctx.Writer.Status(), tt.status)
			}
		})
	}
}
