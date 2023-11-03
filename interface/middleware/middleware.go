package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/online-bookstore/application/dto"
	"github.com/nurcholisnanda/online-bookstore/application/service"
)

type Middleware struct {
	auth service.Authentication
}

var (
	errAccessTokenRequired = "access token is required"
	errInvalidToken        = "access token is invalid"
)

// NewMiddleware returns a wrapper around middleware client.
func NewMiddleware(auth service.Authentication) *Middleware {
	return &Middleware{
		auth: auth,
	}
}

func (m *Middleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: errAccessTokenRequired,
			})
			return
		}

		splitToken := strings.Split(tokenStr, " ")
		if len(splitToken) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: errInvalidToken,
			})
			return
		}

		userID, err := m.auth.ValidateToken(splitToken[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
