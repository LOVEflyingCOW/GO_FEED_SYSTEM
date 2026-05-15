package like

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService *LikeService
}

func NewLikeHandler(likeService *LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

// LikeVideo 点赞视频
func (h *LikeHandler) LikeVideo(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	resp, err := h.likeService.LikeVideo(c.Request.Context(), accountID, uint(videoID))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UnlikeVideo 取消点赞
func (h *LikeHandler) UnlikeVideo(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	resp, err := h.likeService.UnlikeVideo(c.Request.Context(), accountID, uint(videoID))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetLikeStatus 获取点赞状态
func (h *LikeHandler) GetLikeStatus(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	accountID, _ := getAccountID(c)

	resp, err := h.likeService.GetLikeStatus(c.Request.Context(), accountID, uint(videoID))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListLikes 获取用户点赞列表
func (h *LikeHandler) ListLikes(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.likeService.ListLikes(c.Request.Context(), uint(accountID), page, limit)
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
