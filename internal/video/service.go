package video

import (
	"context"
	"errors"
	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/apierror"
	"fmt"
	"log"
	"mime/multipart"

	"gorm.io/gorm"
)

type VideoService struct {
	videoRepository *VideoRepository
	accountRepo     *account.AccountRepository
	uploadService   *UploadService
	baseURL         string
}

func NewVideoService(videoRepository *VideoRepository, accountRepo *account.AccountRepository, uploadService *UploadService, baseURL string) *VideoService {
	return &VideoService{
		videoRepository: videoRepository,
		accountRepo:     accountRepo,
		uploadService:   uploadService,
		baseURL:         baseURL,
	}
}

// UploadVideo 上传视频
func (vs *VideoService) UploadVideo(ctx context.Context, accountID uint, title, description, tags string, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader) (*UploadVideoResponse, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}

	if tags == "" {
		tags = ExtractTags(title + " " + description)
	}

	acc, err := vs.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	return vs.uploadService.UploadVideo(ctx, accountID, acc.Username, title, description, tags, videoFile, coverFile)
}

// GetVideo 获取视频详情
func (vs *VideoService) GetVideo(ctx context.Context, videoID uint, requestAccountID uint) (*GetVideoResponse, error) {
	video, err := vs.videoRepository.FindByID(ctx, videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrVideoNotFound
		}
		return nil, err
	}

	if err := vs.videoRepository.IncreaseViewCount(ctx, videoID); err != nil {
		log.Printf("[WARN] [VideoService] failed to increase view count: %v", err)
	}

	acc, err := vs.accountRepo.FindByID(ctx, video.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetVideoResponse{
		ID:           video.ID,
		AccountID:    video.AccountID,
		Username:     acc.Username,
		AvatarURL:    acc.AvatarURL,
		Title:        video.Title,
		VideoURL:     vs.baseURL + "/" + video.VideoPath,
		CoverURL:     vs.baseURL + "/" + video.CoverPath,
		Duration:     video.Duration,
		Description:  video.Description,
		Tags:         video.Tags,
		ViewCount:    video.ViewCount + 1,
		LikeCount:    video.LikeCount,
		CommentCount: video.CommentCount,
		CreatedAt:    video.CreatedAt,
	}, nil
}

// ListVideos 获取用户视频列表（优化后，避免N+1查询）
func (vs *VideoService) ListVideos(ctx context.Context, accountID uint, page, limit int) (*ListVideoResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 查询视频（1次查询）
	videos, total, err := vs.videoRepository.FindByAccountID(ctx, accountID, page, limit)
	if err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		return &ListVideoResponse{
			Videos: []GetVideoResponse{},
			Total:  total,
		}, nil
	}

	// 收集所有需要查询的用户ID
	accountIDMap := make(map[uint]bool)
	for _, v := range videos {
		accountIDMap[v.AccountID] = true
	}

	// 将map转换为切片
	accountIDs := make([]uint, 0, len(accountIDMap))
	for id := range accountIDMap {
		accountIDs = append(accountIDs, id)
	}

	// 批量查询所有用户（1次查询）
	accountMap := make(map[uint]*account.Account)
	if len(accountIDs) > 0 {
		accounts, err := vs.accountRepo.FindByIDs(ctx, accountIDs)
		if err != nil {
			log.Printf("[WARN] [VideoService] failed to find accounts: %v", err)
		} else {
			for _, acc := range accounts {
				accountMap[acc.ID] = acc
			}
		}
	}

	// 组装响应（无需额外查询）
	videoResponses := make([]GetVideoResponse, 0, len(videos))
	for _, v := range videos {
		acc, exists := accountMap[v.AccountID]
		if !exists || acc == nil {
			continue
		}
		videoResponses = append(videoResponses, GetVideoResponse{
			ID:           v.ID,
			AccountID:    v.AccountID,
			Username:     acc.Username,
			AvatarURL:    acc.AvatarURL,
			Title:        v.Title,
			VideoURL:     vs.baseURL + "/" + v.VideoPath,
			CoverURL:     vs.baseURL + "/" + v.CoverPath,
			Duration:     v.Duration,
			Description:  v.Description,
			Tags:         v.Tags,
			ViewCount:    v.ViewCount,
			LikeCount:    v.LikeCount,
			CommentCount: v.CommentCount,
			CreatedAt:    v.CreatedAt,
		})
	}

	return &ListVideoResponse{
		Videos: videoResponses,
		Total:  total,
	}, nil
}

// DeleteVideo 删除视频
func (vs *VideoService) DeleteVideo(ctx context.Context, videoID uint, accountID uint) error {
	video, err := vs.videoRepository.FindByID(ctx, videoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.ErrVideoNotFound
		}
		return err
	}

	if video.AccountID != accountID {
		return apierror.ErrUnauthorized
	}

	if err := vs.videoRepository.DeleteVideo(ctx, videoID); err != nil {
		return err
	}

	return vs.uploadService.DeleteVideoFiles(video.VideoPath, video.CoverPath)
}
