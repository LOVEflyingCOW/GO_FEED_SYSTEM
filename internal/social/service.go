package social

import (
	"context"
	"errors"
	"feedsystem_video_go/internal/account"

	"gorm.io/gorm"
)

type SocialService struct {
	socialRepository  *SocialRepository
	accountRepository *account.AccountRepository
}

func NewSocialService(socialRepository *SocialRepository, accountRepository *account.AccountRepository) *SocialService {
	return &SocialService{
		socialRepository:  socialRepository,
		accountRepository: accountRepository,
	}
}

func (ss *SocialService) Follow(ctx context.Context, accountID, targetID uint) (*FollowResponse, error) {
	if accountID == targetID {
		return nil, errors.New("cannot follow yourself")
	}

	_, err := ss.accountRepository.FindByID(ctx, targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("target account not found")
		}
		return nil, err
	}

	isFollowing, err := ss.socialRepository.IsFollowing(ctx, accountID, targetID)
	if err != nil {
		return nil, err
	}

	if isFollowing {
		return &FollowResponse{
			AccountID: accountID,
			TargetID:  targetID,
			IsFollow:  true,
		}, nil
	}

	if err := ss.socialRepository.Follow(ctx, accountID, targetID); err != nil {
		return nil, err
	}

	return &FollowResponse{
		AccountID: accountID,
		TargetID:  targetID,
		IsFollow:  true,
	}, nil
}

func (ss *SocialService) Unfollow(ctx context.Context, accountID, targetID uint) (*FollowResponse, error) {
	if err := ss.socialRepository.Unfollow(ctx, accountID, targetID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("not following this account")
		}
		return nil, err
	}

	return &FollowResponse{
		AccountID: accountID,
		TargetID:  targetID,
		IsFollow:  false,
	}, nil
}

func (ss *SocialService) GetFollowers(ctx context.Context, targetID uint, page, limit int) ([]FollowerResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	follows, total, err := ss.socialRepository.GetFollowers(ctx, targetID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	response := make([]FollowerResponse, 0, len(follows))
	for _, f := range follows {
		acc, err := ss.accountRepository.FindByID(ctx, f.AccountID)
		if err != nil {
			continue
		}
		response = append(response, FollowerResponse{
			ID:        acc.ID,
			Username:  acc.Username,
			AvatarURL: acc.AvatarURL,
		})
	}

	return response, total, nil
}

func (ss *SocialService) GetFollowing(ctx context.Context, accountID uint, page, limit int) ([]FollowingResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	follows, total, err := ss.socialRepository.GetFollowing(ctx, accountID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	response := make([]FollowingResponse, 0, len(follows))
	for _, f := range follows {
		acc, err := ss.accountRepository.FindByID(ctx, f.TargetID)
		if err != nil {
			continue
		}
		response = append(response, FollowingResponse{
			ID:        acc.ID,
			Username:  acc.Username,
			AvatarURL: acc.AvatarURL,
		})
	}

	return response, total, nil
}

func (ss *SocialService) GetProfile(ctx context.Context, accountID, requestAccountID uint) (*GetProfileResponse, error) {
	acc, err := ss.accountRepository.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	followerCount, _ := ss.socialRepository.CountFollowers(ctx, accountID)
	followingCount, _ := ss.socialRepository.CountFollowing(ctx, accountID)

	isFollowed := false
	if requestAccountID > 0 && requestAccountID != accountID {
		isFollowed, _ = ss.socialRepository.IsFollowing(ctx, requestAccountID, accountID)
	}

	return &GetProfileResponse{
		AccountID:      accountID,
		Username:       acc.Username,
		AvatarURL:      acc.AvatarURL,
		Bio:            acc.Bio,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
		IsFollowed:     isFollowed,
	}, nil
}

func (ss *SocialService) SearchFriends(ctx context.Context, accountID uint, keyword string, limit int) ([]*account.Account, error) {
	followingIDs, err := ss.socialRepository.GetFollowingIDs(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if len(followingIDs) == 0 {
		return []*account.Account{}, nil
	}

	return ss.accountRepository.SearchByIDs(ctx, followingIDs, keyword, limit)
}
