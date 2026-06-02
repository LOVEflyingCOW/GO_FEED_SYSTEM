package message

import (
	"context"
	"errors"
	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/sse"

	"gorm.io/gorm"
)

type MessageService struct {
	messageRepository *MessageRepository
	accountRepository *account.AccountRepository
	sseHub            *sse.Hub
}

func NewMessageService(messageRepository *MessageRepository, accountRepository *account.AccountRepository, sseHub *sse.Hub) *MessageService {
	return &MessageService{
		messageRepository: messageRepository,
		accountRepository: accountRepository,
		sseHub:            sseHub,
	}
}

func (ms *MessageService) SendMessage(ctx context.Context, senderID, receiverID uint, content string) (*SendMessageResponse, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}

	if senderID == receiverID {
		return nil, errors.New("cannot send message to yourself")
	}

	_, err := ms.accountRepository.FindByID(ctx, receiverID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("receiver not found")
		}
		return nil, err
	}

	message := &Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
	}

	if err := ms.messageRepository.SendMessage(ctx, message); err != nil {
		return nil, err
	}

	sender, err := ms.accountRepository.FindByID(ctx, senderID)
	if err != nil {
		return nil, err
	}

	// 通过 SSE Hub 向接收方发送实时通知
	if ms.sseHub != nil {
		ms.sseHub.SendTo(receiverID, map[string]interface{}{
			"type":    "message",
			"from_id": senderID,
			"from":    sender.Username,
			"content": content,
		})
	}

	return &SendMessageResponse{
		ID:         message.ID,
		SenderID:   message.SenderID,
		SenderName: sender.Username,
		ReceiverID: message.ReceiverID,
		Content:    message.Content,
		CreatedAt:  message.CreatedAt,
	}, nil
}

func (ms *MessageService) GetMessages(ctx context.Context, accountID, otherID uint, page, limit int) ([]MessageResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	messages, total, err := ms.messageRepository.GetMessages(ctx, accountID, otherID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	err = ms.messageRepository.MarkAsRead(ctx, accountID, otherID)
	if err != nil {
	}

	response := make([]MessageResponse, 0, len(messages))
	for _, m := range messages {
		sender, err := ms.accountRepository.FindByID(ctx, m.SenderID)
		if err != nil {
			continue
		}
		response = append(response, MessageResponse{
			ID:         m.ID,
			SenderID:   m.SenderID,
			SenderName: sender.Username,
			Content:    m.Content,
			IsRead:     m.IsRead,
			CreatedAt:  m.CreatedAt,
		})
	}

	return response, total, nil
}

func (ms *MessageService) GetConversations(ctx context.Context, accountID uint) ([]Conversation, error) {
	conversations, err := ms.messageRepository.GetConversations(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for i := range conversations {
		user, err := ms.accountRepository.FindByID(ctx, conversations[i].UserID)
		if err != nil {
			continue
		}
		conversations[i].Username = user.Username
		conversations[i].AvatarURL = user.AvatarURL
	}

	return conversations, nil
}

func (ms *MessageService) CountUnread(ctx context.Context, accountID uint) (int64, error) {
	return ms.messageRepository.CountUnread(ctx, accountID)
}
