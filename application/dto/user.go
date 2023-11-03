package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Function to validate struct fields based on tag validation
func (r *UserRequest) GetError(err validator.ValidationErrors) string {
	for _, e := range err {
		switch e.Tag() {
		case "required":
			return fmt.Sprintf("This field (%v) is required", e.Param())
		}
	}
	return "Uncaught validation error"
}

func (r *LoginRequest) GetError(err validator.ValidationErrors) string {
	for _, e := range err {
		switch e.Tag() {
		case "required":
			return fmt.Sprintf("This field (%v) is required", e.Param())
		}
	}
	return "Uncaught validation error"
}
