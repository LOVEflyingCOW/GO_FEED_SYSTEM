package search

import (
	"context"
)

type HotSearchItem struct {
	Keyword string `json:"keyword"`
	Count   int64  `json:"count"`
}

type SearchService struct {
	searchRepository *SearchRepository
}

func NewSearchService(searchRepository *SearchRepository) *SearchService {
	return &SearchService{searchRepository: searchRepository}
}

func (ss *SearchService) RecordSearch(ctx context.Context, accountID uint, keyword string) error {
	return ss.searchRepository.RecordSearch(ctx, accountID, keyword)
}

func (ss *SearchService) GetSearchHistory(ctx context.Context, accountID uint) ([]string, error) {
	return ss.searchRepository.GetSearchHistory(ctx, accountID, 10)
}

func (ss *SearchService) ClearSearchHistory(ctx context.Context, accountID uint) error {
	return ss.searchRepository.ClearSearchHistory(ctx, accountID)
}

func (ss *SearchService) GetHotSearches(ctx context.Context) ([]HotSearchItem, error) {
	return ss.searchRepository.GetHotSearches(ctx, 10)
}
