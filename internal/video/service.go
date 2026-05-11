package video

import (
	"context"
	"errors"
	"feedsystem_video_go/internal/account"
	"fmt"
	"mime/multipart"
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

func (vs *VideoService) UploadVideo(ctx context.Context, accountID uint, title, description, tags string, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader) (*UploadVideoResponse, error) {
	if tags == "" {
		tags = ExtractTags(title + " " + description)
	}

	account, err := vs.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account info: %w", err)
	}

	return vs.uploadService.UploadVideo(ctx, accountID, account.Username, title, description, tags, videoFile, coverFile)
}

// 避免重复点赞
func (vs *VideoService) GetVideo(ctx context.Context, videoID uint, requestAccountID uint) (*GetVideoResponse, error) {
	video, err := vs.videoRepository.FindByID(ctx, videoID)
	if err != nil {
		return nil, err
	}

	if err := vs.videoRepository.IncreaseViewCount(ctx, videoID); err != nil {
	}

	account, err := vs.accountRepo.FindByID(ctx, video.AccountID)
	if err != nil {
		return nil, err
	}

	return &GetVideoResponse{
		ID:           video.ID,
		AccountID:    video.AccountID,
		Username:     account.Username,
		AvatarURL:    account.AvatarURL,
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

func (vs *VideoService) ListVideos(ctx context.Context, accountID uint, page, limit int) (*ListVideoResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	videos, total, err := vs.videoRepository.FindByAccountID(ctx, accountID, page, limit)
	if err != nil {
		return nil, err
	}

	videoResponses := make([]GetVideoResponse, len(videos))
	for i, v := range videos {
		account, err := vs.accountRepo.FindByID(ctx, v.AccountID)
		if err != nil {
			continue
		}
		videoResponses[i] = GetVideoResponse{
			ID:           v.ID,
			AccountID:    v.AccountID,
			Username:     account.Username,
			AvatarURL:    account.AvatarURL,
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
		}
	}

	return &ListVideoResponse{
		Videos: videoResponses,
		Total:  total,
	}, nil
}

func (vs *VideoService) DeleteVideo(ctx context.Context, videoID uint, accountID uint) error {
	video, err := vs.videoRepository.FindByID(ctx, videoID)
	if err != nil {
		return err
	}

	if video.AccountID != accountID {
		return errors.New("unauthorized to delete this video")
	}

	if err := vs.videoRepository.DeleteVideo(ctx, videoID); err != nil {
		return err
	}

	return vs.uploadService.DeleteVideoFiles(video.VideoPath, video.CoverPath)
}
