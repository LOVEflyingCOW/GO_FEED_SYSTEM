package middleware

import (
	"feedsystem_video_go/internal/apierror"

	"github.com/gin-gonic/gin"
)

// GetAccountID 从 Gin 上下文获取当前登录用户ID
func GetAccountID(c *gin.Context) (uint, error) {
	accountID, exists := c.Get("accountID")
	if !exists {
		return 0, apierror.ErrUnauthorized
	}
	id, ok := accountID.(uint)
	if !ok {
		return 0, apierror.ErrValidation
	}
	return id, nil
}
