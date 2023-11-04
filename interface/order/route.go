package order

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/online-bookstore/application/service"
	"github.com/nurcholisnanda/online-bookstore/infrastructure/repository"
	"github.com/nurcholisnanda/online-bookstore/interface/middleware"
	"gorm.io/gorm"
)

// Setup Order router group
func AddOrderRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	m := middleware.NewMiddleware(service.NewAuthClient(
		os.Getenv("JWT_SECRET"), repository.NewUserRepositoryImpl(db),
	))
	ctr := NewController(service.NewOrderService(
		repository.NewBookRepositoryImpl(db),
		repository.NewOrderRepositoryImpl(db),
	))

	grp := rg.Group("/orders")
	{
		grp.POST("", m.Authenticate(), ctr.MakeOrder)
		grp.GET("/history", m.Authenticate(), ctr.GetOrderHistory)
	}
}
