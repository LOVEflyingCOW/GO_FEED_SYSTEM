package like

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// 只允许每个用户对每个视频点赞一次
func (lr *LikeRepository) CreateLike(ctx context.Context, like *Like) error {
	return lr.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(like).Error
}

// 删除点赞
func (lr *LikeRepository) DeleteLike(ctx context.Context, accountID, videoID uint) error {
	result := lr.db.WithContext(ctx).Where("account_id = ? AND video_id = ?", accountID, videoID).Delete(&Like{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// 查询点赞
func (lr *LikeRepository) FindLike(ctx context.Context, accountID, videoID uint) (*Like, error) {
	var like Like
	err := lr.db.WithContext(ctx).Where("account_id = ? AND video_id = ?", accountID, videoID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// 查询点赞是否存在
func (lr *LikeRepository) ExistsLike(ctx context.Context, accountID, videoID uint) (bool, error) {
	var count int64
	err := lr.db.WithContext(ctx).Model(&Like{}).Where("account_id = ? AND video_id = ?", accountID, videoID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 查询视频点赞数
func (lr *LikeRepository) CountByVideoID(ctx context.Context, videoID uint) (int64, error) {
	var count int64
	err := lr.db.WithContext(ctx).Model(&Like{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}

// 查询用户点赞列表
func (lr *LikeRepository) FindByAccountID(ctx context.Context, accountID uint, page, limit int) ([]*Like, int64, error) {
	var likes []*Like
	var total int64

	offset := (page - 1) * limit

	err := lr.db.WithContext(ctx).Model(&Like{}).Where("account_id = ?", accountID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = lr.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error
	if err != nil {
		return nil, 0, err
	}

	return likes, total, nil
}
