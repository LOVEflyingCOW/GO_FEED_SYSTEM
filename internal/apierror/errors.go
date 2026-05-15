package apierror

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 通用错误定义
var (
	ErrUnauthorized               = errors.New("unauthorized")
	ErrValidation                 = errors.New("validation error")
	ErrUsernameTaken              = errors.New("username already exists")
	ErrInvalidPassword            = errors.New("invalid password")
	ErrInvalidToken               = errors.New("invalid token")
	ErrRefreshTokenExpired        = errors.New("refresh token expired")
	ErrAuthorizationHeaderMissing = errors.New("authorization header required")
	ErrInvalidAuthorizationFormat = errors.New("invalid authorization format")
	ErrInvalidID                  = errors.New("invalid id")
	ErrAccountNotFound            = errors.New("account not found")
	ErrInternalServer             = errors.New("internal server error")
)

// 视频相关错误
var (
	ErrVideoNotFound     = errors.New("video not found")
	ErrVideoUploadFailed = errors.New("video upload failed")
	ErrTitleRequired     = errors.New("title is required")
	ErrVideoFileRequired = errors.New("video file is required")
)

// 评论相关错误
var (
	ErrCommentNotFound = errors.New("comment not found")
	ErrContentRequired = errors.New("content is required")
	ErrReplyToNotFound = errors.New("reply_to comment not found")
)

// 点赞相关错误
var (
	ErrLikeNotFound = errors.New("like not found")
)

// ClassifyHTTPStatus 根据错误类型返回对应的 HTTP 状态码
func ClassifyHTTPStatus(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrInvalidToken),
		errors.Is(err, ErrAuthorizationHeaderMissing), errors.Is(err, ErrInvalidAuthorizationFormat):
		return http.StatusUnauthorized
	case errors.Is(err, ErrValidation), errors.Is(err, ErrInvalidPassword),
		errors.Is(err, ErrContentRequired), errors.Is(err, ErrTitleRequired),
		errors.Is(err, ErrVideoFileRequired), errors.Is(err, ErrInvalidID):
		return http.StatusBadRequest
	case errors.Is(err, ErrUsernameTaken):
		return http.StatusConflict
	case errors.Is(err, ErrVideoNotFound), errors.Is(err, ErrCommentNotFound),
		errors.Is(err, ErrReplyToNotFound), errors.Is(err, ErrLikeNotFound),
		errors.Is(err, ErrAccountNotFound), errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// ErrorResponse 统一错误响应
func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// AbortWithError 统一错误处理，直接中止请求并返回错误响应
func AbortWithError(c *gin.Context, err error) {
	status := ClassifyHTTPStatus(err)
	c.AbortWithStatusJSON(status, ErrorResponse(err))
}
