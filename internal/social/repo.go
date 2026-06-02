package social

import (
	"context"

	"gorm.io/gorm"
)

type SocialRepository struct {
	db *gorm.DB
}

func NewSocialRepository(db *gorm.DB) *SocialRepository {
	return &SocialRepository{db: db}
}

func (sr *SocialRepository) Follow(ctx context.Context, accountID, targetID uint) error {
	follow := &Follow{
		AccountID: accountID,
		TargetID:  targetID,
	}
	return sr.db.WithContext(ctx).Create(follow).Error
}

func (sr *SocialRepository) Unfollow(ctx context.Context, accountID, targetID uint) error {
	result := sr.db.WithContext(ctx).Where("account_id = ? AND target_id = ?", accountID, targetID).Delete(&Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (sr *SocialRepository) IsFollowing(ctx context.Context, accountID, targetID uint) (bool, error) {
	var count int64
	err := sr.db.WithContext(ctx).Model(&Follow{}).Where("account_id = ? AND target_id = ?", accountID, targetID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (sr *SocialRepository) GetFollowers(ctx context.Context, targetID uint, page, limit int) ([]*Follow, int64, error) {
	var follows []*Follow
	var total int64

	offset := (page - 1) * limit

	err := sr.db.WithContext(ctx).Model(&Follow{}).Where("target_id = ?", targetID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = sr.db.WithContext(ctx).Where("target_id = ?", targetID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}

func (sr *SocialRepository) GetFollowing(ctx context.Context, accountID uint, page, limit int) ([]*Follow, int64, error) {
	var follows []*Follow
	var total int64

	offset := (page - 1) * limit

	err := sr.db.WithContext(ctx).Model(&Follow{}).Where("account_id = ?", accountID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = sr.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&follows).Error
	if err != nil {
		return nil, 0, err
	}

	return follows, total, nil
}

func (sr *SocialRepository) CountFollowers(ctx context.Context, targetID uint) (int64, error) {
	var count int64
	err := sr.db.WithContext(ctx).Model(&Follow{}).Where("target_id = ?", targetID).Count(&count).Error
	return count, err
}

func (sr *SocialRepository) CountFollowing(ctx context.Context, accountID uint) (int64, error) {
	var count int64
	err := sr.db.WithContext(ctx).Model(&Follow{}).Where("account_id = ?", accountID).Count(&count).Error
	return count, err
}
