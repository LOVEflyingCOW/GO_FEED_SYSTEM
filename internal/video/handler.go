package video

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService *VideoService
}

func NewVideoHandler(videoService *VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

// UploadVideo 上传视频
func (h *VideoHandler) UploadVideo(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	title := c.PostForm("title")
	if title == "" {
		apierror.AbortWithError(c, apierror.ErrTitleRequired)
		return
	}
	description := c.PostForm("description")
	tags := c.PostForm("tags")

	videoFile, err := c.FormFile("video")
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrVideoFileRequired)
		return
	}

	coverFile, _ := c.FormFile("cover")

	resp, err := h.videoService.UploadVideo(c.Request.Context(), accountID, title, description, tags, videoFile, coverFile)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetVideo 获取视频详情
func (h *VideoHandler) GetVideo(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	requestAccountID, _ := getAccountID(c)

	resp, err := h.videoService.GetVideo(c.Request.Context(), uint(videoID), requestAccountID)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListVideos 获取用户视频列表
func (h *VideoHandler) ListVideos(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	resp, err := h.videoService.ListVideos(c.Request.Context(), uint(accountID), page, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteVideo 删除视频
func (h *VideoHandler) DeleteVideo(c *gin.Context) {
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

	err = h.videoService.DeleteVideo(c.Request.Context(), uint(videoID), accountID)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
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
