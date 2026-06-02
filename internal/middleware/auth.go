package middleware

import (
	"net/http"
	"strings"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("accountID", claims.AccountID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("accountID")
		if !exists {
			c.JSON(apierror.ClassifyHTTPStatus(apierror.ErrUnauthorized), gin.H{"error": apierror.ErrUnauthorized.Error()})
			c.Abort()
			return
		}
		c.Next()
	}
}
