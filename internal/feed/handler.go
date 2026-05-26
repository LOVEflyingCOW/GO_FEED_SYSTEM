package feed

import (
	"errors"
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"

	"github.com/gin-gonic/gin"
)

type FeedHandler struct {
	feedService *FeedService
}

func NewFeedHandler(feedService *FeedService) *FeedHandler {
	return &FeedHandler{feedService: feedService}
}

func (h *FeedHandler) GetFeed(c *gin.Context) {
	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")
	feedType := c.DefaultQuery("type", "latest")
	tag := c.Query("tag")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	accountID, _ := getAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Cursor:    cursor,
		Limit:     limit,
		Type:      feedType,
		Tag:       tag,
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) GetHotFeed(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	accountID, _ := getAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Limit:     limit,
		Type:      "hot",
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) GetFollowingFeed(c *gin.Context) {
	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	accountID, err := getAccountID(c)
	if err != nil {
		c.JSON(apierror.ClassifyHTTPStatus(err), gin.H{"error": err.Error()})
		return
	}

	req := FeedRequest{
		AccountID: accountID,
		Cursor:    cursor,
		Limit:     limit,
		Type:      "following",
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) GetTagFeed(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tag is required"})
		return
	}

	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	accountID, _ := getAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Cursor:    cursor,
		Limit:     limit,
		Type:      "tag",
		Tag:       tag,
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func getAccountID(c *gin.Context) (uint, error) {
	accountID, exists := c.Get("accountID")
	if !exists {
		return 0, errors.New("account not authenticated")
	}
	return accountID.(uint), nil
}
