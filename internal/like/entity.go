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
	Likes []*Like `json:"likes"`
	Total int64   `json:"total"`
}
