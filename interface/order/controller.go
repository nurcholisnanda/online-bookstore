package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/application/service"
)

var errFetchUserIDErr = "error fetching user id from context"

type Controller struct {
	service service.OrderService
}

func NewController(service service.OrderService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) MakeOrder(ctx *gin.Context) {
	var req dto.OrderRequest
	userID, found := ctx.Get("user_id")
	if !found {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: errFetchUserIDErr,
		})
		return
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	err = c.service.MakeOrder(&req, userID.(uint))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Success making order",
	})
}

func (c *Controller) GetOrderHistory(ctx *gin.Context) {
	userID, found := ctx.Get("user_id")
	if !found {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: errFetchUserIDErr,
		})
		return
	}

	orderHistory, err := c.service.GetOrderHistory(userID.(uint))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    orderHistory,
	})
}
