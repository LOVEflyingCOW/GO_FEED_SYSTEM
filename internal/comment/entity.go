package comment

import "time"

// Comment 评论实体，存储评论的基本信息
type Comment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                         // 评论唯一标识
	AccountID uint      `json:"account_id"`                                   // 发布评论的用户ID
	VideoID   uint      `json:"video_id"`                                     // 评论所属视频ID
	Content   string    `gorm:"type:text" json:"content"`                     // 评论内容
	ReplyTo   uint      `json:"reply_to,omitempty"`                           // 回复的目标评论ID，0表示不是回复
	Mentioned string    `gorm:"type:varchar(512)" json:"mentioned,omitempty"` // @提及的用户列表，JSON格式
	CreatedAt time.Time `json:"created_at"`                                   // 评论创建时间
}

// CreateCommentRequest 创建评论请求结构
type CreateCommentRequest struct {
	VideoID uint   `json:"video_id"`           // 评论目标视频ID
	Content string `json:"content"`            // 评论内容
	ReplyTo uint   `json:"reply_to,omitempty"` // 回复的目标评论ID，可选
}

// CreateCommentResponse 创建评论响应结构
type CreateCommentResponse struct {
	ID        uint      `json:"id"`                   // 新创建评论的ID
	AccountID uint      `json:"account_id"`           // 发布评论的用户ID
	Username  string    `json:"username"`             // 发布评论的用户名
	AvatarURL string    `json:"avatar_url,omitempty"` // 用户头像URL
	VideoID   uint      `json:"video_id"`             // 评论所属视频ID
	Content   string    `json:"content"`              // 评论内容
	ReplyTo   uint      `json:"reply_to,omitempty"`   // 回复的目标评论ID
	CreatedAt time.Time `json:"created_at"`           // 评论创建时间
}

// ListCommentRequest 获取评论列表请求结构
type ListCommentRequest struct {
	VideoID uint `json:"video_id"` // 视频ID，获取该视频的评论
	Page    int  `json:"page"`     // 页码，从1开始
	Limit   int  `json:"limit"`    // 每页数量
}

// ListCommentResponse 获取评论列表响应结构
type ListCommentResponse struct {
	Comments []CommentWithUser `json:"comments"` // 评论列表（包含用户信息）
	Total    int64             `json:"total"`    // 评论总数
}

// CommentWithUser 包含用户信息的评论结构
type CommentWithUser struct {
	ID          uint      `json:"id"`                      // 评论ID
	AccountID   uint      `json:"account_id"`              // 用户ID
	Username    string    `json:"username"`                // 用户名
	AvatarURL   string    `json:"avatar_url,omitempty"`    // 用户头像URL
	VideoID     uint      `json:"video_id"`                // 视频ID
	Content     string    `json:"content"`                 // 评论内容
	ReplyTo     uint      `json:"reply_to,omitempty"`      // 回复目标评论ID
	ReplyToUser ReplyUser `json:"reply_to_user,omitempty"` // 被回复用户信息
	CreatedAt   time.Time `json:"created_at"`              // 创建时间
}

// ReplyUser 被回复用户信息结构
type ReplyUser struct {
	ID        uint   `json:"id"`                   // 用户ID
	Username  string `json:"username"`             // 用户名
	AvatarURL string `json:"avatar_url,omitempty"` // 用户头像URL
}

// DeleteCommentRequest 删除评论请求结构
type DeleteCommentRequest struct {
	CommentID uint `json:"comment_id"` // 要删除的评论ID
}
