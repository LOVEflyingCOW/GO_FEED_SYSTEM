package like

import "time"

type Like struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AccountID uint      `json:"account_id"`
	VideoID   uint      `json:"video_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeRequest struct {
	VideoID uint `json:"video_id"`
}

type LikeResponse struct {
	VideoID   uint  `json:"video_id"`
	LikeCount int64 `json:"like_count"`
	IsLiked   bool  `json:"is_liked"`
}

type LikeListRequest struct {
	AccountID uint `json:"account_id"`
	Page      int  `json:"page"`
	Limit     int  `json:"limit"`
}

type LikeListResponse struct {
	Videos []*VideoItem `json:"videos"`
	Total  int64        `json:"total"`
}

type VideoItem struct {
	ID           uint   `json:"id"`
	AccountID    uint   `json:"account_id"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	PlayURL      string `json:"play_url"`
	CoverURL     string `json:"cover_url"`
	Duration     int    `json:"duration"`
	Description  string `json:"description"`
	Tags         string `json:"tags"`
	ViewCount    int64  `json:"view_count"`
	LikeCount    int64  `json:"like_count"`
	CommentCount int64  `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
}
