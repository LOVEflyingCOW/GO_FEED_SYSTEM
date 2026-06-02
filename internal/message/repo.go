package message

import (
	"context"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (mr *MessageRepository) SendMessage(ctx context.Context, message *Message) error {
	return mr.db.WithContext(ctx).Create(message).Error
}

func (mr *MessageRepository) GetMessages(ctx context.Context, senderID, receiverID uint, page, limit int) ([]*Message, int64, error) {
	var messages []*Message
	var total int64

	offset := (page - 1) * limit

	err := mr.db.WithContext(ctx).Model(&Message{}).
		Where("(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
			senderID, receiverID, receiverID, senderID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = mr.db.WithContext(ctx).
		Where("(from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)",
			senderID, receiverID, receiverID, senderID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (mr *MessageRepository) GetConversations(ctx context.Context, accountID uint) ([]Conversation, error) {
	rows, err := mr.db.WithContext(ctx).Raw(`
		SELECT 
			CASE WHEN from_id = ? THEN to_id ELSE from_id END as user_id,
			MAX(created_at) as updated_at,
			(SELECT content FROM messages WHERE 
				(from_id = ? AND to_id = CASE WHEN from_id = ? THEN to_id ELSE from_id END) OR
				(from_id = CASE WHEN from_id = ? THEN to_id ELSE from_id END AND to_id = ?)
				ORDER BY created_at DESC LIMIT 1) as last_message,
			SUM(CASE WHEN to_id = ? AND is_read = false THEN 1 ELSE 0 END) as unread_count
		FROM messages
		WHERE from_id = ? OR to_id = ?
		GROUP BY user_id
		ORDER BY updated_at DESC
	`, accountID, accountID, accountID, accountID, accountID, accountID, accountID, accountID).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		err := rows.Scan(&c.UserID, &c.UpdatedAt, &c.LastMessage, &c.UnreadCount)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}

	return conversations, nil
}

func (mr *MessageRepository) MarkAsRead(ctx context.Context, receiverID, senderID uint) error {
	return mr.db.WithContext(ctx).Model(&Message{}).
		Where("to_id = ? AND from_id = ? AND is_read = false", receiverID, senderID).
		Update("is_read", true).Error
}

func (mr *MessageRepository) CountUnread(ctx context.Context, receiverID uint) (int64, error) {
	var count int64
	err := mr.db.WithContext(ctx).Model(&Message{}).
		Where("to_id = ? AND is_read = false", receiverID).
		Count(&count).Error
	return count, err
}
