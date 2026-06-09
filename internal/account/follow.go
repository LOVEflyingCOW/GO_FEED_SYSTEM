package account

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Follow 关注关系实体
type Follow struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	FollowerID  uint `gorm:"column:account_id" json:"account_id"` // 关注者ID（映射到数据库的 account_id 字段）
	FollowingID uint `gorm:"column:target_id" json:"target_id"`   // 被关注者ID（映射到数据库的 target_id 字段）
}

// GetFollowing 获取用户关注的所有用户ID
func (ar *AccountRepository) GetFollowing(ctx context.Context, accountID uint) ([]uint, error) {
	var follows []Follow
	err := ar.db.WithContext(ctx).Where("account_id = ?", accountID).Find(&follows).Error
	if err != nil {
		return nil, err
	}

	followingIDs := make([]uint, 0, len(follows))
	for _, f := range follows {
		followingIDs = append(followingIDs, f.FollowingID)
	}

	return followingIDs, nil
}

// FollowAccount 关注用户
func (ar *AccountRepository) FollowAccount(ctx context.Context, followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	follow := Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return ar.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&follow).Error
}

// UnfollowAccount 取消关注
func (ar *AccountRepository) UnfollowAccount(ctx context.Context, followerID, followingID uint) error {
	result := ar.db.WithContext(ctx).Where("account_id = ? AND target_id = ?", followerID, followingID).Delete(&Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetFollowerCount 获取粉丝数
func (ar *AccountRepository) GetFollowerCount(ctx context.Context, accountID uint) (int64, error) {
	var count int64
	err := ar.db.WithContext(ctx).Model(&Follow{}).Where("target_id = ?", accountID).Count(&count).Error
	return count, err
}

// GetFollowingCount 获取关注数
func (ar *AccountRepository) GetFollowingCount(ctx context.Context, accountID uint) (int64, error) {
	var count int64
	err := ar.db.WithContext(ctx).Model(&Follow{}).Where("account_id = ?", accountID).Count(&count).Error
	return count, err
}
