package feed

import (
	"time"
)

type FeedItem struct {
	VideoID      uint      `json:"video_id"`
	AccountID    uint      `json:"account_id"`
	Username     string    `json:"username"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
	Title        string    `json:"title"`
	VideoURL     string    `json:"video_url"`
	CoverURL     string    `json:"cover_url"`
	Duration     int       `json:"duration"`
	Description  string    `json:"description,omitempty"`
	Tags         string    `json:"tags"`
	ViewCount    int64     `json:"view_count"`
	LikeCount    int64     `json:"like_count"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	IsLiked      bool      `json:"is_liked,omitempty"`
}

type FeedRequest struct {
	AccountID uint   `json:"account_id,omitempty"`
	Cursor    int64  `json:"cursor,omitempty"`
	Limit     int    `json:"limit,omitempty"`
	Type      string `json:"type,omitempty"` // "latest", "hot", "following", "tag"
	Tag       string `json:"tag,omitempty"`
}

type FeedResponse struct {
	Items   []FeedItem `json:"items"`
	HasMore bool       `json:"has_more"`
	Next    int64      `json:"next"`
}

type TimelineItem struct {
	VideoID   uint    `json:"video_id"`
	Timestamp int64   `json:"timestamp"`
	Score     float64 `json:"score"`
}
