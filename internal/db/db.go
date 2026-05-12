package db

import (
	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/like"
	"feedsystem_video_go/internal/video"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// AutoMigrate() ：自动迁移数据库表结构
// 添加了 account.Account 作为需要迁移的模型
// 方便后续添加表而不冗余main
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&account.Account{},
		&video.Video{},
		&like.Like{},
	)
}

func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
