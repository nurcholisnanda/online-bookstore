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

func NewController(service service.UserService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) Register(ctx *gin.Context) {
	var req dto.UserRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if val, ok := err.(validator.ValidationErrors); ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: req.GetError(val),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

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

func (c *Controller) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		if val, ok := err.(validator.ValidationErrors); ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Success: false,
				Message: req.GetError(val),
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

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
