package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Define UserRequest struct that will be used for http request body
type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Define LoginRequest struct that will be used for http request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Function to validate UserRequest struct fields based on tag validation
func (r *UserRequest) GetError(err validator.ValidationErrors) string {
	for _, e := range err {
		switch e.Tag() {
		case "required":
			return fmt.Sprintf("This field (%v) is required", e.Field())
		}
	}
	return "Uncaught validation error"
}

// Function to validate LoginRequest struct fields based on tag validation
func (r *LoginRequest) GetError(err validator.ValidationErrors) string {
	for _, e := range err {
		switch e.Tag() {
		case "required":
			return fmt.Sprintf("This field (%v) is required", e.Field())
		}
	}
	return "Uncaught validation error"
}
