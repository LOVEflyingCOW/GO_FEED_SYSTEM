package comment

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"

	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/video"

	"gorm.io/gorm"
)

// CommentService 评论业务逻辑层
type CommentService struct {
	commentRepository *CommentRepository
	accountRepository *account.AccountRepository
	videoRepository   *video.VideoRepository
}

// NewCommentService 创建评论Service实例
func NewCommentService(commentRepository *CommentRepository, accountRepository *account.AccountRepository, videoRepository *video.VideoRepository) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
		accountRepository: accountRepository,
		videoRepository:   videoRepository,
	}
}

// CreateComment 创建评论
func (cs *CommentService) CreateComment(ctx context.Context, accountID, videoID uint, content string, replyTo uint) (*CreateCommentResponse, error) {
	// 参数校验
	if content == "" {
		return nil, apierror.ErrContentRequired
	}

	// 校验视频是否存在
	if _, err := cs.videoRepository.FindByID(ctx, videoID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apierror.ErrVideoNotFound
		}
		return nil, err
	}

	// 校验回复目标评论是否存在（如果有回复）
	if replyTo > 0 {
		if _, err := cs.commentRepository.FindByID(ctx, replyTo); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apierror.ErrReplyToNotFound
			}
			return nil, err
		}
	}

	// 提取 @提及用户
	mentioned := extractMentions(content)

	// 创建评论
	comment := &Comment{
		AccountID: accountID,
		VideoID:   videoID,
		Content:   content,
		ReplyTo:   replyTo,
		Mentioned: mentioned,
	}

	if err := cs.commentRepository.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	// 更新视频评论数
	if err := cs.videoRepository.IncreaseCommentCount(ctx, videoID); err != nil {
		log.Printf("[WARN] [CommentService] failed to increase comment count: %v", err)
	}

	// 获取用户信息
	acc, err := cs.accountRepository.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &CreateCommentResponse{
		ID:        comment.ID,
		AccountID: comment.AccountID,
		Username:  acc.Username,
		AvatarURL: acc.AvatarURL,
		VideoID:   comment.VideoID,
		Content:   comment.Content,
		ReplyTo:   comment.ReplyTo,
		CreatedAt: comment.CreatedAt,
	}, nil
}

// DeleteComment 删除评论
func (cs *CommentService) DeleteComment(ctx context.Context, commentID, accountID uint) error {
	// 查询评论
	comment, err := cs.commentRepository.FindByID(ctx, commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apierror.ErrCommentNotFound
		}
		return err
	}

	// 权限校验
	if comment.AccountID != accountID {
		return apierror.ErrUnauthorized
	}

	// 删除评论
	if err := cs.commentRepository.DeleteComment(ctx, commentID); err != nil {
		return err
	}

	// 更新视频评论数
	if err := cs.videoRepository.DecreaseCommentCount(ctx, comment.VideoID); err != nil {
		log.Printf("[WARN] [CommentService] failed to decrease comment count: %v", err)
	}

	return nil
}

// ListComments 获取评论列表
func (cs *CommentService) ListComments(ctx context.Context, videoID uint, page, limit int) (*ListCommentResponse, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 查询评论（1次查询）
	comments, total, err := cs.commentRepository.FindByVideoID(ctx, videoID, page, limit)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return &ListCommentResponse{
			Comments: []CommentWithUser{},
			Total:    total,
		}, nil
	}

	// 收集所有需要查询的用户ID和父评论ID
	accountIDMap := make(map[uint]bool)
	replyToIDMap := make(map[uint]bool)

	for _, c := range comments {
		accountIDMap[c.AccountID] = true
		if c.ReplyTo > 0 {
			replyToIDMap[c.ReplyTo] = true
		}
	}

	// 将map转换为切片
	accountIDs := make([]uint, 0, len(accountIDMap))
	for id := range accountIDMap {
		accountIDs = append(accountIDs, id)
	}

	replyToIDs := make([]uint, 0, len(replyToIDMap))
	for id := range replyToIDMap {
		replyToIDs = append(replyToIDs, id)
	}

	// 批量查询所有用户（1次查询）
	accountMap := make(map[uint]*account.Account)
	if len(accountIDs) > 0 {
		accounts, err := cs.accountRepository.FindByIDs(ctx, accountIDs)
		if err != nil {
			log.Printf("[WARN] [CommentService] failed to find accounts: %v", err)
		} else {
			for _, acc := range accounts {
				accountMap[acc.ID] = acc
			}
		}
	}

	// 批量查询所有父评论（1次查询）
	parentCommentMap := make(map[uint]*Comment)
	if len(replyToIDs) > 0 {
		parentComments, err := cs.commentRepository.FindByIDs(ctx, replyToIDs)
		if err != nil {
			log.Printf("[WARN] [CommentService] failed to find parent comments: %v", err)
		} else {
			for _, pc := range parentComments {
				parentCommentMap[pc.ID] = pc
				// 如果父评论的用户还没查询，添加到查询列表
				if !accountIDMap[pc.AccountID] {
					accountIDs = append(accountIDs, pc.AccountID) //accountids被加长，需要补充map
					accountIDMap[pc.AccountID] = true
				}
			}
		}
	}

	// 补充查询父评论的用户（最多1次查询）
	if len(accountIDs) > len(accountMap) {
		var missingAccountIDs []uint
		for id := range accountIDMap {
			if _, exists := accountMap[id]; !exists {
				missingAccountIDs = append(missingAccountIDs, id)
			}
		}
		if len(missingAccountIDs) > 0 {
			accounts, err := cs.accountRepository.FindByIDs(ctx, missingAccountIDs)
			if err == nil {
				for _, acc := range accounts {
					accountMap[acc.ID] = acc
				}
			}
		}
	}

	// 组装响应（无需额外查询）
	commentResponses := make([]CommentWithUser, 0, len(comments))
	for _, c := range comments {
		acc, exists := accountMap[c.AccountID]
		if !exists || acc == nil {
			continue
		}

		var replyToUser ReplyUser
		if c.ReplyTo > 0 {
			// 去父评论map里找：被回复的那条评论存在吗
			if parentComment, exists := parentCommentMap[c.ReplyTo]; exists {
				// 找到了父评论 → 再拿这个父评论的发布者
				// 去用户map里找：这个人存在吗？
				if replyAccount, exists := accountMap[parentComment.AccountID]; exists {
					replyToUser = ReplyUser{
						ID:        replyAccount.ID,
						Username:  replyAccount.Username,
						AvatarURL: replyAccount.AvatarURL,
					}
				}
			}
		}

		commentResponses = append(commentResponses, CommentWithUser{
			ID:          c.ID,
			AccountID:   c.AccountID,
			Username:    acc.Username,
			AvatarURL:   acc.AvatarURL,
			VideoID:     c.VideoID,
			Content:     c.Content,
			ReplyTo:     c.ReplyTo,
			ReplyToUser: replyToUser,
			CreatedAt:   c.CreatedAt,
		})
	}

	return &ListCommentResponse{
		Comments: commentResponses,
		Total:    total,
	}, nil
}

// extractMentions 提取评论中的 @提及用户
// "@张三 @李四 这个视频太好看了！"提出"张三,李四"
func extractMentions(content string) string {
	//规则
	re := regexp.MustCompile(`@([\p{Han}\w]+)`)
	//提取所有匹配项
	matches := re.FindAllStringSubmatch(content, -1)
	var mentions []string
	for _, match := range matches {
		if len(match) > 1 {
			mentions = append(mentions, match[1])
		}
	}
	return strings.Join(mentions, ",")
}
