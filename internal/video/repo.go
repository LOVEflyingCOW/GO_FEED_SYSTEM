package video

import (
	"context"

	"gorm.io/gorm"
)

type VideoRepository struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (vr *VideoRepository) CreateVideo(ctx context.Context, video *Video) error {
	return vr.db.WithContext(ctx).Create(video).Error
}

func (vr *VideoRepository) FindByID(ctx context.Context, id uint) (*Video, error) {
	var video Video
	if err := vr.db.WithContext(ctx).First(&video, id).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (vr *VideoRepository) FindByAccountID(ctx context.Context, accountID uint, page, limit int) ([]*Video, int64, error) {
	var videos []*Video
	var total int64

	offset := (page - 1) * limit

	err := vr.db.WithContext(ctx).Model(&Video{}).Where("account_id = ?", accountID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = vr.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error
	if err != nil {
		return nil, 0, err
	}

	return videos, total, nil
}

func (vr *VideoRepository) DeleteVideo(ctx context.Context, id uint) error {
	result := vr.db.WithContext(ctx).Delete(&Video{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (vr *VideoRepository) UpdateVideo(ctx context.Context, id uint, updates map[string]interface{}) error {
	return vr.db.WithContext(ctx).Model(&Video{}).Where("id = ?", id).Updates(updates).Error
}

func (vr *VideoRepository) IncreaseViewCount(ctx context.Context, videoID uint) error {
	return vr.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (vr *VideoRepository) IncreaseLikeCount(ctx context.Context, videoID uint) error {
	return vr.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

func (vr *VideoRepository) DecreaseLikeCount(ctx context.Context, videoID uint) error {
	return vr.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

func (vr *VideoRepository) IncreaseCommentCount(ctx context.Context, videoID uint) error {
	return vr.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
}

func (vr *VideoRepository) DecreaseCommentCount(ctx context.Context, videoID uint) error {
	return vr.db.WithContext(ctx).Model(&Video{}).Where("id = ?", videoID).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error
}

func (vr *VideoRepository) FindLatest(ctx context.Context, limit int) ([]*Video, error) {
	var videos []*Video
	err := vr.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&videos).Error
	return videos, err
}
