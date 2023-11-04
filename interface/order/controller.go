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

// Order controller constructor
func NewController(service service.OrderService) *Controller {
	return &Controller{
		service: service,
	}
}

// MakeOrder function to call MakeOrder service to add a new order
func (c *Controller) MakeOrder(ctx *gin.Context) {
	var req dto.OrderRequest
	//get userID value using user_id key in our context from token
	userID, found := ctx.Get("user_id")
	if !found {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: errFetchUserIDErr,
		})
		return
	}

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
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

// GetOrderHistory function to call GetOrderHistory service
// and return user order history data as the response
func (c *Controller) GetOrderHistory(ctx *gin.Context) {
	//get userID value using user_id key in our context from token
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
