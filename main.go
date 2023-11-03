package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nurcholisnanda/online-bookstore/infrastructure/persistence"
	"github.com/nurcholisnanda/online-bookstore/interface/book"
	"github.com/nurcholisnanda/online-bookstore/interface/order"
	"github.com/nurcholisnanda/online-bookstore/interface/user"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic("No env gotten")
	}
}

func main() {
	db, err := persistence.NewDatabase()
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate()
	gormDB := db.GetDB()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.JSON(http.StatusOK, map[string]any{
			"message": "Online Bookstore Server",
		})
	})

	v := r.Group("/v1")

	user.AddUserRoutes(v, gormDB)
	book.AddBookRoutes(v, gormDB)
	order.AddOrderRoutes(v, gormDB)

	if err = r.Run(":8080"); err != nil {
		log.Panic(err)
	}
}
