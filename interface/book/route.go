package book

import (
	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/online-bookstore/application/service"
	"github.com/nurcholisnanda/online-bookstore/infrastructure/repository"
	"gorm.io/gorm"
)

func AddBookRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	ctr := NewController(service.NewBookService(repository.NewBookRepositoryImpl(db)))

	grp := rg.Group("/books")
	{
		grp.GET("", ctr.GetBooks)
	}
}
