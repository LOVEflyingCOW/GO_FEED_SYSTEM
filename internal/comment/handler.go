package comment

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论控制器，处理HTTP请求
type CommentHandler struct {
	commentService *CommentService
}

// NewCommentHandler 创建评论Handler实例
func NewCommentHandler(commentService *CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

// CreateComment 创建评论
// POST /api/comments
func (h *CommentHandler) CreateComment(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}

	resp, err := h.commentService.CreateComment(c.Request.Context(), accountID, req.VideoID, req.Content, req.ReplyTo)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteComment 删除评论
// DELETE /api/comments/:comment_id
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	commentIDStr := c.Param("comment_id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	err = h.commentService.DeleteComment(c.Request.Context(), uint(commentID), accountID)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

// ListComments 获取评论列表
// GET /api/videos/:video_id/comments?page=1&limit=20
func (h *CommentHandler) ListComments(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.commentService.ListComments(c.Request.Context(), uint(videoID), page, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// getAccountID 从上下文获取用户ID
func getAccountID(c *gin.Context) (uint, error) {
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
