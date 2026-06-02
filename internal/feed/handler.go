package feed

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware"

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

	accountID, _ := middleware.GetAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Cursor:    cursor,
		Limit:     limit,
		Type:      feedType,
		Tag:       tag,
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) GetHotFeed(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)

	accountID, _ := middleware.GetAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Limit:     limit,
		Type:      "hot",
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) GetFollowingFeed(c *gin.Context) {
	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
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
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) GetTagFeed(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}

	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	accountID, _ := middleware.GetAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Cursor:    cursor,
		Limit:     limit,
		Type:      "tag",
		Tag:       tag,
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FeedHandler) SearchFeed(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}

	cursorStr := c.DefaultQuery("cursor", "0")
	limitStr := c.DefaultQuery("limit", "20")

	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)
	limit, _ := strconv.Atoi(limitStr)

	accountID, _ := middleware.GetAccountID(c)

	req := FeedRequest{
		AccountID: accountID,
		Cursor:    cursor,
		Limit:     limit,
		Type:      "search",
		Tag:       keyword,
	}

	resp, err := h.feedService.GetFeed(c.Request.Context(), req)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
