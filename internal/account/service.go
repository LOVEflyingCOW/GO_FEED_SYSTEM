package account

import (
	"context"
	"errors"
	"feedsystem_video_go/internal/apierror"
	"feedsystem_video_go/internal/auth"
	"log"
	"strconv"
	"time"

	rediscache "feedsystem_video_go/internal/middleware/redis"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService struct {
	accountRepository *AccountRepository
	cache             *rediscache.Client
}

func NewAccountService(accountRepository *AccountRepository, cache *rediscache.Client) *AccountService {
	return &AccountService{accountRepository: accountRepository, cache: cache}
}

// CreateAccount 创建账户
func (as *AccountService) CreateAccount(ctx context.Context, account *Account) error {
	if len(account.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	account.Password = string(passwordHash)

	if err := as.accountRepository.CreateAccount(ctx, account); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return apierror.ErrUsernameTaken
		}
		return err
	}
	return nil
}

// Rename 修改用户名
func (as *AccountService) Rename(ctx context.Context, accountID uint, newUsername string) (string, error) {
	if newUsername == "" {
		return "", errors.New("new_username is required")
	}

	token, err := auth.GenerateToken(accountID, newUsername)
	if err != nil {
		return "", err
	}

	if err := as.accountRepository.RenameWithToken(ctx, accountID, newUsername, token); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return "", apierror.ErrUsernameTaken
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("account not found")
		}
		return "", err
	}

	if as.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		if err := as.cache.SetBytes(cacheCtx, as.cache.Key("account:%d", accountID), []byte(token), 24*time.Hour); err != nil {
			log.Printf("[WARN] [AccountService] failed to set cache: %v", err)
		}
	}
	return token, nil
}

// ChangePassword 修改密码
func (as *AccountService) ChangePassword(ctx context.Context, username, oldPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	acc, err := as.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(oldPassword)); err != nil {
		return apierror.ErrInvalidPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := as.accountRepository.ChangePassword(ctx, acc.ID, string(passwordHash)); err != nil {
		return err
	}

	return as.Logout(ctx, acc.ID)
}

// FindByID 根据ID查询用户
func (as *AccountService) FindByID(ctx context.Context, id uint) (*Account, error) {
	return as.accountRepository.FindByID(ctx, id)
}

// FindByUsername 根据用户名查询用户
func (as *AccountService) FindByUsername(ctx context.Context, username string) (*Account, error) {
	return as.accountRepository.FindByUsername(ctx, username)
}

// Login 登录
func (as *AccountService) Login(ctx context.Context, username, password string) (string, string, error) {
	acc, err := as.FindByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		return "", "", apierror.ErrInvalidPassword
	}

	accessToken, err := auth.GenerateToken(acc.ID, acc.Username)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.GenerateRefreshToken(acc.ID)
	if err != nil {
		return "", "", err
	}

	if err := as.accountRepository.Login(ctx, acc.ID, accessToken, refreshToken); err != nil {
		return "", "", err
	}

	if as.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		if err := as.cache.SetBytes(cacheCtx, as.cache.Key("account:%d", acc.ID), []byte(accessToken), 24*time.Hour); err != nil {
			log.Printf("[WARN] [AccountService] failed to set cache: %v", err)
		}
		if err := as.cache.SetBytes(cacheCtx, as.cache.Key("account:%d:refresh", acc.ID), []byte(refreshToken), 7*24*time.Hour); err != nil {
			log.Printf("[WARN] [AccountService] failed to set refresh cache: %v", err)
		}
		if err := as.cache.SetBytes(cacheCtx, as.cache.Key("refresh:%s", refreshToken), []byte(strconv.FormatUint(uint64(acc.ID), 10)), 7*24*time.Hour); err != nil {
			log.Printf("[WARN] [AccountService] failed to set refresh lookup: %v", err)
		}
	}

	return accessToken, refreshToken, nil
}

// Logout 登出
func (as *AccountService) Logout(ctx context.Context, accountID uint) error {
	if err := as.accountRepository.Logout(ctx, accountID); err != nil {
		return err
	}

	if as.cache != nil {
		cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
		defer cancel()
		_ = as.cache.Del(cacheCtx, as.cache.Key("account:%d", accountID))
		_ = as.cache.Del(cacheCtx, as.cache.Key("account:%d:refresh", accountID))
	}
	return nil
}

// Refresh 刷新Token
func (as *AccountService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	if as.cache == nil {
		return "", "", apierror.ErrInvalidToken
	}

	cacheCtx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()

	b, err := as.cache.GetBytes(cacheCtx, as.cache.Key("refresh:%s", refreshToken))
	if err != nil {
		return "", "", apierror.ErrInvalidToken
	}

	accountID, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil || accountID == 0 {
		return "", "", apierror.ErrInvalidToken
	}

	acc, err := as.FindByID(ctx, uint(accountID))
	if err != nil {
		return "", "", err
	}

	if acc.RefreshToken != refreshToken {
		return "", "", apierror.ErrInvalidToken
	}

	newToken, err := auth.GenerateToken(acc.ID, acc.Username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := auth.GenerateRefreshToken(acc.ID)
	if err != nil {
		return "", "", err
	}

	if err := as.accountRepository.Login(ctx, acc.ID, newToken, newRefreshToken); err != nil {
		return "", "", err
	}

	if err := as.cache.SetBytes(cacheCtx, as.cache.Key("account:%d", acc.ID), []byte(newToken), 24*time.Hour); err != nil {
		log.Printf("[WARN] [AccountService] failed to set cache: %v", err)
	}

	return newToken, newRefreshToken, nil
}
