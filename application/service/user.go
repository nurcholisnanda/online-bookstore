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

//go:generate mockgen -source=user.go -destination=mock/user.go -package=mock
type UserService interface {
	Register(*dto.UserRequest) error
	Login(*dto.LoginRequest) (string, error)
}

// User service constructor
func NewUserService(auth Authentication, userRepo user.Repository) UserService {
	return &userService{
		auth:     auth,
		userRepo: userRepo,
	}
}

// Register service will hashed password and call InsertUser
// in our user repository contract
func (s *userService) Register(req *dto.UserRequest) error {
	//hashing password
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

// Login service will check whether email in login request is exist.
// Moreover, will create access token if the user is exist
// and password request is correct.
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
