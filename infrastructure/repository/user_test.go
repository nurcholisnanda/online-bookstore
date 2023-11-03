package repository

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"gorm.io/gorm"
)

func TestNewUserRepositoryImpl(t *testing.T) {
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
		want user.Repository
	}{
		{
			name: "implemented",
			args: args{
				db: db,
			},
			want: NewUserRepositoryImpl(db),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepositoryImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepositoryImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepositoryImpl_InsertUser(t *testing.T) {
	type args struct {
		user *user.User
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	tests := []struct {
		name    string
		r       *userRepositoryImpl
		args    args
		wantErr bool
	}{
		{
			name: "success insert user",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				user: &user.User{
					Model:    gorm.Model{},
					Name:     "name1",
					Email:    "email1",
					Password: "password1",
				},
			},
			wantErr: false,
		},
		{
			name: "error empty user",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				user: nil,
			},
			wantErr: true,
		},
		{
			name: "error insert duplicate email",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				user: &user.User{
					Model:    gorm.Model{},
					Name:     "name1",
					Email:    "email1",
					Password: "password1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.InsertUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepositoryImpl.InsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepositoryImpl_FindUserByID(t *testing.T) {
	type args struct {
		id uint
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	userSample, err := seedUser(db)
	if err != nil {
		fmt.Println("error seeding user")
	}

	tests := []struct {
		name     string
		r        *userRepositoryImpl
		args     args
		wantUser *user.User
		wantErr  bool
	}{
		{
			name: "success find user by id",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				id: userSample.ID,
			},
			wantUser: userSample,
			wantErr:  false,
		},
		{
			name: "error id not found",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				id: 2,
			},
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := tt.r.FindUserByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepositoryImpl.FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("userRepositoryImpl.FindUserByID() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func Test_userRepositoryImpl_FindUserByEmail(t *testing.T) {
	type args struct {
		email string
	}
	db, err := DBConn()
	if err != nil {
		t.Fatalf("want non error, got %#v", err)
	}

	userSample, err := seedUser(db)
	if err != nil {
		fmt.Println("error seeding user")
	}

	tests := []struct {
		name     string
		r        *userRepositoryImpl
		args     args
		wantUser *user.User
		wantErr  bool
	}{
		{
			name: "success find user by email",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				email: userSample.Email,
			},
			wantUser: userSample,
			wantErr:  false,
		},
		{
			name: "error email not found",
			r: &userRepositoryImpl{
				db: db,
			},
			args: args{
				email: "not existed email",
			},
			wantUser: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := tt.r.FindUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepositoryImpl.FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("userRepositoryImpl.FindUserByEmail() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
