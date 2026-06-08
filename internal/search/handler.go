package search

import (
	"net/http"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *SearchService
}

func NewSearchHandler(searchService *SearchService) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// RecordSearch 记录搜索历史
func (h *SearchHandler) RecordSearch(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	var req struct {
		Keyword string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}

	if req.Keyword == "" {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}

	if err := h.searchService.RecordSearch(c.Request.Context(), accountID, req.Keyword); err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "search recorded"})
}

// GetSearchHistory 获取搜索历史
func (h *SearchHandler) GetSearchHistory(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	history, err := h.searchService.GetSearchHistory(c.Request.Context(), accountID)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": history})
}

// ClearSearchHistory 清除搜索历史
func (h *SearchHandler) ClearSearchHistory(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	if err := h.searchService.ClearSearchHistory(c.Request.Context(), accountID); err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "history cleared"})
}

// GetHotSearches 获取热门搜索
func (h *SearchHandler) GetHotSearches(c *gin.Context) {
	hotSearches, err := h.searchService.GetHotSearches(c.Request.Context())
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"hot_searches": hotSearches})
}
