package repository

import (
	"errors"
	"strings"

	"github.com/nurcholisnanda/online-bookstore/domain/user"
	"gorm.io/gorm"
)

var (
	errEmailAlreadyExist = errors.New("email already exist")
	errUserNotFound      = errors.New("user not found")
)

type userRepositoryImpl struct {
	db *gorm.DB
}

// User repository implementation constructor
func NewUserRepositoryImpl(db *gorm.DB) user.Repository {
	return &userRepositoryImpl{
		db: db,
	}
}

// Implementation of InsertUser repo contract to insert
// a new user into our database
func (r *userRepositoryImpl) InsertUser(user *user.User) error {
	if err := r.db.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			return errEmailAlreadyExist
		}
		return err
	}
	return nil
}

// Implementation of FindUserByID repo contract to find
// a user by id from our database
func (r *userRepositoryImpl) FindUserByID(id uint) (*user.User, error) {
	var user user.User
	if err := r.db.Where("id = ?", id).Take(&user).Error; err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Implementation of FindUserByEmail repo contract to find
// a user by email from our database
func (r *userRepositoryImpl) FindUserByEmail(email string) (*user.User, error) {
	var user user.User
	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, errUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
