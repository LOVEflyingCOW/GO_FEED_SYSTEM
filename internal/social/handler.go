package social

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

type SocialHandler struct {
	socialService *SocialService
}

func NewSocialHandler(socialService *SocialService) *SocialHandler {
	return &SocialHandler{socialService: socialService}
}

func (h *SocialHandler) Follow(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	targetIDStr := c.Param("target_id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	resp, err := h.socialService.Follow(c.Request.Context(), accountID, uint(targetID))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SocialHandler) Unfollow(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	targetIDStr := c.Param("target_id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	resp, err := h.socialService.Unfollow(c.Request.Context(), accountID, uint(targetID))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *SocialHandler) GetFollowers(c *gin.Context) {
	targetIDStr := c.Param("target_id")
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	followers, total, err := h.socialService.GetFollowers(c.Request.Context(), uint(targetID), page, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"total":     total,
	})
}

func (h *SocialHandler) GetFollowing(c *gin.Context) {
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

	following, total, err := h.socialService.GetFollowing(c.Request.Context(), uint(accountID), page, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"following": following,
		"total":     total,
	})
}

func (h *SocialHandler) GetProfile(c *gin.Context) {
	accountIDStr := c.Param("account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	requestAccountID, _ := middleware.GetAccountID(c)

	profile, err := h.socialService.GetProfile(c.Request.Context(), uint(accountID), requestAccountID)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *SocialHandler) SearchFriends(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	keyword := c.Query("keyword")
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	accounts, err := h.socialService.SearchFriends(c.Request.Context(), accountID, keyword, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	type FriendItem struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}

	result := make([]FriendItem, 0, len(accounts))
	for _, acc := range accounts {
		result = append(result, FriendItem{
			ID:        acc.ID,
			Username:  acc.Username,
			AvatarURL: acc.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, result)
}
