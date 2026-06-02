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

	err := vr.db.WithContext(ctx).Model(&Video{}).Where("author_id = ?", accountID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = vr.db.WithContext(ctx).Where("author_id = ?", accountID).
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

// FindByAccountIDs 批量查询多个用户的视频
func (vr *VideoRepository) FindByAccountIDs(ctx context.Context, accountIDs []uint, limit, offset int) ([]*Video, error) {
	var videos []*Video
	err := vr.db.WithContext(ctx).
		Where("author_id IN ?", accountIDs).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error
	return videos, err
}

// FindByTag 按标签查询视频（精确匹配完整标签）
func (vr *VideoRepository) FindByTag(ctx context.Context, tag string, limit, offset int) ([]*Video, error) {
	var videos []*Video
	err := vr.db.WithContext(ctx).
		Where("tags = ? OR tags LIKE ? OR tags LIKE ? OR tags LIKE ?",
			tag,           // 标签本身
			tag+",%",      // 标签在开头
			"%,"+tag+",%", // 标签在中间
			"%,"+tag).     // 标签在结尾
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error
	return videos, err
}

// Search 搜索视频（按标题、描述、标签模糊匹配）
func (vr *VideoRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]*Video, error) {
	var videos []*Video
	err := vr.db.WithContext(ctx).
		Where("title LIKE ? OR description LIKE ? OR tags LIKE ?",
			"%"+keyword+"%",
			"%"+keyword+"%",
			"%"+keyword+"%").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&videos).Error
	return videos, err
}

// FindByIDs 批量查询视频
func (vr *VideoRepository) FindByIDs(ctx context.Context, ids []uint) ([]*Video, error) {
	var videos []*Video
	if len(ids) == 0 {
		return videos, nil
	}
	err := vr.db.WithContext(ctx).Where("id IN ?", ids).Find(&videos).Error
	return videos, err
}
