package app

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"note-service/internal/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken := c.Request.Header.Get("X-Access-Token")
		if userToken == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.New("empty token"))
			return
		}
		userID, err := jwt.ParseToken(userToken)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, errors.New("jwt parse error"))
			return
		}
		c.Set("userId", userID)
	}
}
