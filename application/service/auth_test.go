package service

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"github.com/nurcholisnanda/online-bookstore/domain/user/mock"
)

func TestNewAuthClient(t *testing.T) {
	type args struct {
		secret   string
		userRepo user.Repository
	}
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRepository(ctrl)
	secret := "test secret"

	tests := []struct {
		name string
		args args
		want Authentication
	}{
		{
			name: "implemented",
			args: args{
				secret:   secret,
				userRepo: mockRepo,
			},
			want: NewAuthClient(secret, mockRepo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthClient(tt.args.secret, tt.args.userRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authClient_CreateAccessToken(t *testing.T) {
	type args struct {
		userID uint
	}
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRepository(ctrl)
	secret := "test secret"

	tests := []struct {
		name    string
		c       *authClient
		args    args
		wantErr bool
	}{
		{
			name: "success create access token",
			c: &authClient{
				secret:   secret,
				userRepo: mockRepo,
			},
			args: args{
				userID: 1,
			},
			wantErr: false,
		},
		{
			name: "error create access token",
			c: &authClient{
				secret:   secret,
				userRepo: mockRepo,
			},
			args: args{
				userID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.c.CreateAccessToken(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.CreateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_authClient_ValidateToken(t *testing.T) {
	type args struct {
		requestToken string
	}
	ctrl := gomock.NewController(t)
	mockRepo := mock.NewMockRepository(ctrl)
	findUserByEmailRecord := func(user *user.User, err error) func(m *mock.MockRepository) {
		return func(m *mock.MockRepository) {
			m.EXPECT().FindUserByID(gomock.Any()).Return(user, err)
		}
	}
	var userID uint = 1
	secret := "test secret"
	client := &authClient{
		secret: secret,
	}
	validToken, _ := client.CreateAccessToken(uint(userID))

	tests := []struct {
		name    string
		c       *authClient
		args    args
		want    uint
		wantErr bool
	}{
		{
			name: "success validating token",
			c: &authClient{
				secret:   secret,
				userRepo: mockRepo,
			},
			args: args{
				requestToken: validToken,
			},
			want:    userID,
			wantErr: false,
		},
		{
			name: "error find user",
			c:    client,
			args: args{
				requestToken: "invalid token",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "success validating token":
				findUserByEmailRecord(nil, nil)(mockRepo)
			}
			got, err := tt.c.ValidateToken(tt.args.requestToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authClient.ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
