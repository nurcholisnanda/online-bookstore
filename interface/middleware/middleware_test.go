package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/application/service"
	"github.com/nurcholisnanda/online-bookstore/application/service/mock"
)

func TestNewMiddleware(t *testing.T) {
	type args struct {
		auth service.Authentication
	}
	ctrl := gomock.NewController(t)
	mockAuth := mock.NewMockAuthentication(ctrl)

	tests := []struct {
		name string
		args args
		want *Middleware
	}{
		{
			name: "ïmplemented",
			args: args{
				auth: mockAuth,
			},
			want: NewMiddleware(mockAuth),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "ïmplemented":

			}
			if got := NewMiddleware(tt.args.auth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMiddleware_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mock.NewMockAuthentication(ctrl)
	var userID uint = 1
	mockValidateToken := func(userID uint, err error) func(m *mock.MockAuthentication) {
		return func(m *mock.MockAuthentication) {
			m.EXPECT().ValidateToken(gomock.Any()).Return(userID, err)
		}
	}
	m := NewMiddleware(mockAuth)

	tests := []struct {
		name      string
		authToken string
		status    int
		userID    any
	}{
		{
			name:      "authorized",
			authToken: fmt.Sprintf("Bearer %s", "validtoken"),
			status:    http.StatusOK,
			userID:    1,
		},
		{
			name:      "invalid input empty token",
			authToken: "",
			status:    http.StatusUnauthorized,
			userID:    0,
		},
		{
			name:      "invalid auth input format",
			authToken: "Bearer",
			status:    http.StatusUnauthorized,
			userID:    0,
		},
		{
			name:      "fail authenticate",
			authToken: fmt.Sprintf("Bearer %s", "invalidtoken"),
			status:    http.StatusUnauthorized,
			userID:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "authorized":
				mockValidateToken(userID, nil)(mockAuth)
			case "fail authenticate":
				mockValidateToken(0, errors.New("any error"))(mockAuth)
			}
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request, _ = http.NewRequest("POST", "/auth/validurl",
				bytes.NewBuffer([]byte{}),
			)
			ctx.Request.Header.Add("Authorization", tt.authToken)
			m.Authenticate()(ctx)
			if ok := reflect.DeepEqual(ctx.Writer.Status(), tt.status); !ok {
				t.Errorf("Middleware.Authenticate() = %v, want %v", w.Result().Status, tt.status)
			}
			if id, found := ctx.Get("user_id"); found {
				if reflect.DeepEqual(id.(uint), tt.userID) {
					t.Errorf("Middleware.Authenticate() = %v, want %v", id.(uint), tt.userID)
				}
			}
		})
	}
}
