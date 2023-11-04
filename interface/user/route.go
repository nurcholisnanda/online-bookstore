package user

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/online-bookstore/application/service"
	"github.com/nurcholisnanda/online-bookstore/infrastructure/repository"
	"gorm.io/gorm"
)

// Setup User router group
func AddUserRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	jwtSecret := os.Getenv("JWT_SECRET")
	userRepo := repository.NewUserRepositoryImpl(db)
	ctr := NewController(service.NewUserService(
		service.NewAuthClient(jwtSecret, userRepo),
		userRepo,
	))

	grp := rg.Group("/users")
	{
		grp.POST("/register", ctr.Register)
		grp.POST("/login", ctr.Login)
	}
}
