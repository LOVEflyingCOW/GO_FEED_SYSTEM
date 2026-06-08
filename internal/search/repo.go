package search

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type SearchHistory struct {
	ID        uint      `gorm:"primaryKey"`
	AccountID uint      `json:"account_id"`
	Keyword   string    `json:"keyword"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchRepository struct {
	db *gorm.DB
}

func NewSearchRepository(db *gorm.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (sr *SearchRepository) RecordSearch(ctx context.Context, accountID uint, keyword string) error {
	history := &SearchHistory{
		AccountID: accountID,
		Keyword:   keyword,
	}
	return sr.db.WithContext(ctx).Create(history).Error
}

func (sr *SearchRepository) GetSearchHistory(ctx context.Context, accountID uint, limit int) ([]string, error) {
	var histories []SearchHistory
	err := sr.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Order("created_at DESC").
		Limit(limit).
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	keywords := make([]string, 0, len(histories))
	for _, h := range histories {
		keywords = append(keywords, h.Keyword)
	}

	return keywords, nil
}

func (sr *SearchRepository) ClearSearchHistory(ctx context.Context, accountID uint) error {
	return sr.db.WithContext(ctx).
		Where("account_id = ?", accountID).
		Delete(&SearchHistory{}).Error
}

func (sr *SearchRepository) GetHotSearches(ctx context.Context, limit int) ([]HotSearchItem, error) {
	var results []struct {
		Keyword string
		Count   int64
	}

	err := sr.db.WithContext(ctx).
		Model(&SearchHistory{}).
		Select("keyword, COUNT(*) as count").
		Group("keyword").
		Order("count DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	hotSearches := make([]HotSearchItem, 0, len(results))
	for _, r := range results {
		hotSearches = append(hotSearches, HotSearchItem{
			Keyword: r.Keyword,
			Count:   r.Count,
		})
	}

	return hotSearches, nil
}
