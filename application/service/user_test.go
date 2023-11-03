package service

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	mock_auth "github.com/nurcholisnanda/online-bookstore/application/service/mock"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
	mock_repo "github.com/nurcholisnanda/online-bookstore/domain/user/mock"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mock_auth.NewMockAuthentication(ctrl)
	mockRepo := mock_repo.NewMockRepository(ctrl)

	type args struct {
		auth     Authentication
		userRepo user.Repository
	}
	tests := []struct {
		name string
		args args
		want UserService
	}{
		{
			name: "implemented",
			args: args{
				auth:     mockAuth,
				userRepo: mockRepo,
			},
			want: NewUserService(mockAuth, mockRepo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.auth, tt.args.userRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mock_auth.NewMockAuthentication(ctrl)
	mockRepo := mock_repo.NewMockRepository(ctrl)
	mockRepoRecord := func(err error) func(m *mock_repo.MockRepository) {
		return func(m *mock_repo.MockRepository) {
			m.EXPECT().InsertUser(gomock.Any()).Return(err)
		}
	}
	req := seedUserReq()

	type args struct {
		req *dto.UserRequest
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		wantErr bool
	}{
		{
			name: "success register new user",
			s: &userService{
				auth:     mockAuth,
				userRepo: mockRepo,
			},
			args: args{
				req: req,
			},
			wantErr: false,
		},
		{
			name: "error register new user",
			s: &userService{
				auth:     mockAuth,
				userRepo: mockRepo,
			},
			args: args{
				req: &dto.UserRequest{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if strings.Contains(tt.name, "success") {
			mockRepoRecord(nil)(mockRepo)
		} else {
			mockRepoRecord(errors.New("any error"))(mockRepo)
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Register(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("userService.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockAuth := mock_auth.NewMockAuthentication(ctrl)
	mockRepo := mock_repo.NewMockRepository(ctrl)
	mockFindUserByEmail := func(user *user.User, err error) func(m *mock_repo.MockRepository) {
		return func(m *mock_repo.MockRepository) {
			m.EXPECT().FindUserByEmail(gomock.Any()).Return(user, err)
		}
	}
	mockCreateAccessToken := func(token string, err error) func(m *mock_auth.MockAuthentication) {
		return func(m *mock_auth.MockAuthentication) {
			m.EXPECT().CreateAccessToken(gomock.Any()).Return(token, err)
		}
	}
	req := seedLoginReq()
	user := seedUser()
	signedToken := "anytoken"

	type args struct {
		req *dto.LoginRequest
	}
	tests := []struct {
		name    string
		s       *userService
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success login account",
			s: &userService{
				auth:     mockAuth,
				userRepo: mockRepo,
			},
			args: args{
				req: req,
			},
			want:    signedToken,
			wantErr: false,
		},
		{
			name: "fail find user",
			s: &userService{
				auth:     mockAuth,
				userRepo: mockRepo,
			},
			args: args{
				req: req,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "fail create token",
			s: &userService{
				auth:     mockAuth,
				userRepo: mockRepo,
			},
			args: args{
				req: req,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "success login account":
				mockFindUserByEmail(user, nil)(mockRepo)
				mockCreateAccessToken(signedToken, nil)(mockAuth)
			case "fail find user":
				mockFindUserByEmail(nil, errors.New("any error"))(mockRepo)
			case "fail create token":
				mockFindUserByEmail(user, nil)(mockRepo)
				mockCreateAccessToken("", errors.New("any error"))(mockAuth)
			}
			got, err := tt.s.Login(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}
