package comment

import (
	"context"

	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论Repository实例
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

// CreateComment 创建评论记录
func (cr *CommentRepository) CreateComment(ctx context.Context, comment *Comment) error {
	return cr.db.WithContext(ctx).Create(comment).Error
}

// FindByID 根据ID查询评论
func (cr *CommentRepository) FindByID(ctx context.Context, id uint) (*Comment, error) {
	var comment Comment
	if err := cr.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// FindByVideoID 根据视频ID分页查询评论
// 返回评论列表、总数和错误
func (cr *CommentRepository) FindByVideoID(ctx context.Context, videoID uint, page, limit int) ([]*Comment, int64, error) {
	var comments []*Comment
	var total int64

	offset := (page - 1) * limit

	err := cr.db.WithContext(ctx).Model(&Comment{}).Where("video_id = ?", videoID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = cr.db.WithContext(ctx).Where("video_id = ?", videoID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// DeleteComment 删除评论记录
// 如果记录不存在，返回gorm.ErrRecordNotFound
func (cr *CommentRepository) DeleteComment(ctx context.Context, id uint) error {
	result := cr.db.WithContext(ctx).Delete(&Comment{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// CountByVideoID 统计指定视频的评论数量
func (cr *CommentRepository) CountByVideoID(ctx context.Context, videoID uint) (int64, error) {
	var count int64
	err := cr.db.WithContext(ctx).Model(&Comment{}).Where("video_id = ?", videoID).Count(&count).Error
	return count, err
}

// FindByAccountID 根据用户ID分页查询评论
func (cr *CommentRepository) FindByAccountID(ctx context.Context, accountID uint, page, limit int) ([]*Comment, int64, error) {
	var comments []*Comment
	var total int64

	offset := (page - 1) * limit

	err := cr.db.WithContext(ctx).Model(&Comment{}).Where("account_id = ?", accountID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = cr.db.WithContext(ctx).Where("account_id = ?", accountID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// FindByIDs 批量查询评论（避免N+1查询）
func (cr *CommentRepository) FindByIDs(ctx context.Context, ids []uint) ([]*Comment, error) {
	var comments []*Comment
	if err := cr.db.WithContext(ctx).Where("id IN ?", ids).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
