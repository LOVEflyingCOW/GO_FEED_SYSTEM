package video

import (
	"errors"
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VideoHandler struct {
	videoService *VideoService
}

func NewVideoHandler(videoService *VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		c.JSON(apierror.ClassifyHTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}
	description := c.PostForm("description")
	tags := c.PostForm("tags")

	videoFile, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video file is required"})
		return
	}

	coverFile, _ := c.FormFile("cover")

	resp, err := h.videoService.UploadVideo(c.Request.Context(), accountID, title, description, tags, videoFile, coverFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *VideoHandler) GetVideo(c *gin.Context) {
	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	requestAccountID, _ := getAccountID(c)

	resp, err := h.videoService.GetVideo(c.Request.Context(), uint(videoID), requestAccountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *VideoHandler) ListVideos(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account_id"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *VideoHandler) DeleteVideo(c *gin.Context) {
	accountID, err := getAccountID(c)
	if err != nil {
		c.JSON(apierror.ClassifyHTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	videoIDStr := c.Param("video_id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video_id"})
		return
	}

	err = h.videoService.DeleteVideo(c.Request.Context(), uint(videoID), accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
			return
		}
		if errors.Is(err, errors.New("unauthorized to delete this video")) {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
}

func getAccountID(c *gin.Context) (uint, error) {
	accountID, exists := c.Get("accountID")
	if !exists {
		return 0, errors.New("account not authenticated")
	}
	return accountID.(uint), nil
}
