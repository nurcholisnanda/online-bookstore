package book

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/application/service"
)

type Controller struct {
	service service.BookService
}

func NewController(service service.BookService) *Controller {
	return &Controller{
		service: service,
	}
}

func (c *Controller) GetBooks(ctx *gin.Context) {
	books, err := c.service.GetBooks()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    books,
	})
}
