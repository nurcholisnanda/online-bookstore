package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/application/service"
)

type Controller struct {
	service service.UserService
}

// User controller constructor
func NewController(service service.UserService) *Controller {
	return &Controller{
		service: service,
	}
}

// Register function to call register service in our program
func (c *Controller) Register(ctx *gin.Context) {
	var req dto.UserRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: req.GetError(err.(validator.ValidationErrors)),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// calling register service to register new user
	err = c.service.Register(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Register success",
	})
}

// Login function to call login service in our program.
func (c *Controller) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: req.GetError(err.(validator.ValidationErrors)),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// calling login service to return token if user credential is correct.
	token, err := c.service.Login(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Login success",
		Data:    token,
	})
}
