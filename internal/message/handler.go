package message

import (
	"net/http"
	"strconv"

	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *MessageService
}

func NewMessageHandler(messageService *MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	senderID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.AbortWithError(c, apierror.ErrValidation)
		return
	}

	resp, err := h.messageService.SendMessage(c.Request.Context(), senderID, req.ReceiverID, req.Content)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MessageHandler) GetMessages(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	otherIDStr := c.Param("other_id")
	otherID, err := strconv.ParseUint(otherIDStr, 10, 64)
	if err != nil {
		apierror.AbortWithError(c, apierror.ErrInvalidID)
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "50")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	messages, total, err := h.messageService.GetMessages(c.Request.Context(), accountID, uint(otherID), page, limit)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
		"total":    total,
	})
}

func (h *MessageHandler) GetConversations(c *gin.Context) {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	conversations, err := h.messageService.GetConversations(c.Request.Context(), accountID)
	if err != nil {
		apierror.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, conversations)
}
