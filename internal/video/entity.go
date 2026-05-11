package video

import (
	"time"
)

// Video 视频实体，对应数据库 videos 表
type Video struct {
	ID           uint      `gorm:"primaryKey" json:"id"`                   // 视频ID，主键
	AccountID    uint      `json:"account_id"`                             // 发布者账户ID
	Title        string    `gorm:"type:varchar(255)" json:"title"`         // 视频标题
	VideoPath    string    `gorm:"type:varchar(512)" json:"video_path"`    // 视频文件路径
	CoverPath    string    `gorm:"type:varchar(512)" json:"cover_path"`    // 封面图片路径
	Duration     int       `json:"duration"`                               // 视频时长（秒）
	Description  string    `gorm:"type:text" json:"description,omitempty"` // 视频描述
	Tags         string    `gorm:"type:varchar(512)" json:"tags"`          // 标签，逗号分隔
	ViewCount    int64     `json:"view_count"`                             // 观看次数
	LikeCount    int64     `json:"like_count"`                             // 点赞次数
	CommentCount int64     `json:"comment_count"`                          // 评论次数
	CreatedAt    time.Time `json:"created_at"`                             // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`                             // 更新时间
}

// UploadVideoRequest 上传视频请求结构体
type UploadVideoRequest struct {
	Title       string `json:"title"`                 // 视频标题
	Description string `json:"description,omitempty"` // 视频描述（可选）
	Tags        string `json:"tags,omitempty"`        // 标签（可选）
}

// UploadVideoResponse 上传视频响应结构体
type UploadVideoResponse struct {
	VideoID     uint   `json:"video_id"`              // 视频ID
	VideoURL    string `json:"video_url"`             // 视频访问URL
	CoverURL    string `json:"cover_url"`             // 封面访问URL
	Title       string `json:"title"`                 // 视频标题
	Description string `json:"description,omitempty"` // 视频描述
	Tags        string `json:"tags"`                  // 标签
}

// GetVideoRequest 获取单个视频请求结构体
type GetVideoRequest struct {
	VideoID uint `json:"video_id"` // 视频ID
}

// GetVideoResponse 获取视频响应结构体
type GetVideoResponse struct {
	ID           uint      `json:"id"`                    // 视频ID
	AccountID    uint      `json:"account_id"`            // 发布者账户ID
	Username     string    `json:"username"`              // 发布者用户名
	AvatarURL    string    `json:"avatar_url,omitempty"`  // 发布者头像URL（可选）
	Title        string    `json:"title"`                 // 视频标题
	VideoURL     string    `json:"video_url"`             // 视频访问URL
	CoverURL     string    `json:"cover_url"`             // 封面访问URL
	Duration     int       `json:"duration"`              // 视频时长（秒）
	Description  string    `json:"description,omitempty"` // 视频描述（可选）
	Tags         string    `json:"tags"`                  // 标签
	ViewCount    int64     `json:"view_count"`            // 观看次数
	LikeCount    int64     `json:"like_count"`            // 点赞次数
	CommentCount int64     `json:"comment_count"`         // 评论次数
	CreatedAt    time.Time `json:"created_at"`            // 创建时间
}

// ListVideoRequest 获取视频列表请求结构体
type ListVideoRequest struct {
	AccountID uint `json:"account_id"` // 账户ID（可选，用于获取特定用户的视频）
	Page      int  `json:"page"`       // 页码，从1开始
	Limit     int  `json:"limit"`      // 每页数量
}

// ListVideoResponse 获取视频列表响应结构体
type ListVideoResponse struct {
	Videos []GetVideoResponse `json:"videos"` // 视频列表
	Total  int64              `json:"total"`  // 总数量
}

// DeleteVideoRequest 删除视频请求结构体
type DeleteVideoRequest struct {
	VideoID uint `json:"video_id"` // 视频ID
}
