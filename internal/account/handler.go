package account

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService *AccountService
}

func NewAccountHandler(accountService *AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// CreateAccount 创建账户
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	if err := h.accountService.CreateAccount(c.Request.Context(), &Account{
		Username: req.Username,
		Password: req.Password,
	}); err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account created"})
}

// Rename 修改用户名
func (h *AccountHandler) Rename(c *gin.Context) {
	var req RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	token, err := h.accountService.Rename(c.Request.Context(), accountID, req.NewUsername)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// UpdateProfile 更新用户资料
func (h *AccountHandler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	if err := h.accountService.UpdateProfile(c.Request.Context(), accountID, req.Username, req.Bio); err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// ChangePassword 修改密码
func (h *AccountHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	if err := h.accountService.ChangePassword(c.Request.Context(), req.Username, req.OldPassword, req.NewPassword); err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// FindByID 根据ID查询用户
func (h *AccountHandler) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}
	acc, err := h.accountService.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, FindByIDResponse{
		ID:        acc.ID,
		Username:  acc.Username,
		AvatarURL: acc.AvatarURL,
		Bio:       acc.Bio,
	})
}

// FindByUsername 根据用户名查询用户
func (h *AccountHandler) FindByUsername(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	acc, err := h.accountService.FindByUsername(c.Request.Context(), username)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, FindByUsernameResponse{ID: acc.ID, Username: acc.Username})
}

// Login 登录
func (h *AccountHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	acc, err := h.accountService.FindByUsername(c.Request.Context(), req.Username)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	accessToken, refreshToken, err := h.accountService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		AccountID:    acc.ID,
		Username:     acc.Username,
	})
}

// Logout 登出
func (h *AccountHandler) Logout(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	if err := h.accountService.Logout(c.Request.Context(), accountID); err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account logged out"})
}

// SearchUsers 搜索用户
func (h *AccountHandler) SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	accounts, err := h.accountService.SearchUsers(c.Request.Context(), keyword, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	type UserItem struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		AvatarURL string `json:"avatar_url"`
	}

	result := make([]UserItem, 0, len(accounts))
	for _, acc := range accounts {
		result = append(result, UserItem{
			ID:        acc.ID,
			Username:  acc.Username,
			AvatarURL: acc.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, result)
}

// Refresh 刷新Token
func (h *AccountHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}
	token, refreshToken, err := h.accountService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token, RefreshToken: refreshToken})
}
