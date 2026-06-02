package social

import "time"

type Follow struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AccountID uint      `json:"account_id"`
	TargetID  uint      `json:"target_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowRequest struct {
	TargetID uint `json:"target_id"`
}

type FollowResponse struct {
	AccountID uint `json:"account_id"`
	TargetID  uint `json:"target_id"`
	IsFollow  bool `json:"is_follow"`
}

type FollowerResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type FollowingResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type GetProfileResponse struct {
	AccountID      uint   `json:"account_id"`
	Username       string `json:"username"`
	AvatarURL      string `json:"avatar_url,omitempty"`
	Bio            string `json:"bio,omitempty"`
	VideoCount     int64  `json:"video_count"`
	LikeCount      int64  `json:"like_count"`
	FollowerCount  int64  `json:"follower_count"`
	FollowingCount int64  `json:"following_count"`
	IsFollowed     bool   `json:"is_followed"`
}
