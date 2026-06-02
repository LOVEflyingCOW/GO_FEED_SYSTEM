package like

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware"

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
	accountID, err := middleware.GetAccountID(c)
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
	accountID, err := middleware.GetAccountID(c)
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

	accountID, _ := middleware.GetAccountID(c)

	resp, err := h.likeService.GetLikeStatus(c.Request.Context(), accountID, uint(videoID))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListLikes 获取当前用户点赞列表
func (h *LikeHandler) ListLikes(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	resp, err := h.likeService.ListLikes(c.Request.Context(), accountID, int(cursor), limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
