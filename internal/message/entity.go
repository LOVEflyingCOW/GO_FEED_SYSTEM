package message

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"column:from_id" json:"sender_id"`
	ReceiverID uint      `gorm:"column:to_id" json:"receiver_id"`
	Content    string    `gorm:"type:text" json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
}

type SendMessageResponse struct {
	ID         uint      `json:"id"`
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type Conversation struct {
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	LastMessage string    `json:"last_message"`
	UnreadCount int64     `json:"unread_count"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MessageResponse struct {
	ID         uint      `json:"id"`
	SenderID   uint      `json:"sender_id"`
	SenderName string    `json:"sender_name"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}
