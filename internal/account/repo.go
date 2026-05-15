package account

import (
	"context"

	"gorm.io/gorm"
)

// 让 AccountRepository 拥有数据库操作能力
type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (ar *AccountRepository) CreateAccount(ctx context.Context, account *Account) error {
	if err := ar.db.WithContext(ctx).Create(account).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) Rename(ctx context.Context, id uint, newUsername string) error {
	result := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("username", newUsername)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (ar *AccountRepository) RenameWithToken(ctx context.Context, id uint, newUsername string, token string) error {
	result := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Updates(map[string]any{
		"username": newUsername,
		"token":    token,
	})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (ar *AccountRepository) ChangePassword(ctx context.Context, id uint, newPassword string) error {
	result := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("password", newPassword)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (ar *AccountRepository) FindByID(ctx context.Context, id uint) (*Account, error) {
	var account Account
	if err := ar.db.WithContext(ctx).First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (ar *AccountRepository) FindByUsername(ctx context.Context, username string) (*Account, error) {
	var account Account
	if err := ar.db.WithContext(ctx).Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (ar *AccountRepository) Login(ctx context.Context, id uint, token, refreshToken string) error {
	if err := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Updates(map[string]interface{}{
		"token":         token,
		"refresh_token": refreshToken,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) Logout(ctx context.Context, id uint) error {
	if err := ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Updates(map[string]interface{}{
		"token":         "",
		"refresh_token": "",
	}).Error; err != nil {
		return err
	}
	return nil
}

func (ar *AccountRepository) UpdateAvatar(ctx context.Context, accountID uint, avatarURL string) error {
	return ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", accountID).Update("avatar_url", avatarURL).Error
}

func (ar *AccountRepository) UpdateToken(ctx context.Context, id uint, token string) error {
	return ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Update("token", token).Error
}

func (ar *AccountRepository) UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error {
	return ar.db.WithContext(ctx).Model(&Account{}).Where("id = ?", id).Updates(updates).Error
}

func (ar *AccountRepository) FindAll(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	if err := ar.db.WithContext(ctx).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// FindByIDs 批量查询用户（避免N+1查询）
func (ar *AccountRepository) FindByIDs(ctx context.Context, ids []uint) ([]*Account, error) {
	var accounts []*Account
	if err := ar.db.WithContext(ctx).Where("id IN ?", ids).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}
