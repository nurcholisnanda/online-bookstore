package service

import (
	"errors"

	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"golang.org/x/crypto/bcrypt"
)

var errWrongPassword = errors.New("wrong password")

type userService struct {
	auth     Authentication
	userRepo user.Repository
}

type UserService interface {
	Register(*dto.UserRequest) error
	Login(*dto.LoginRequest) (string, error)
}

func NewUserService(auth Authentication, userRepo user.Repository) UserService {
	return &userService{
		auth:     auth,
		userRepo: userRepo,
	}
}

func (s *userService) Register(req *dto.UserRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	err = s.userRepo.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Login(req *dto.LoginRequest) (string, error) {
	user, err := s.userRepo.FindUserByEmail(req.Email)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errWrongPassword
	}

	token, err := s.auth.CreateAccessToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
