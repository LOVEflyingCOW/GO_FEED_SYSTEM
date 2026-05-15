// package main

// import (
// 	"feedsystem_video_go/internal/account"
// 	"feedsystem_video_go/internal/config"
// 	"feedsystem_video_go/internal/db"
// 	"feedsystem_video_go/internal/middleware/redis"
// 	"log"
// 	"os"
// )

// func main() {
// 	configPath := os.Getenv("CONFIG_PATH")
// 	if configPath == "" {
// 		configPath = "configs/config.yaml"
// 	}
// 	log.Printf("Loading config from %s", configPath)
// 	cfg, usedDefault, err := config.LoadLocalDev(configPath)
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}
// 	if usedDefault {
// 		log.Printf("Config File %s not found, using default local config", configPath)
// 	} else {
// 		log.Printf("Config loaded from file: %s", configPath)
// 	}

// 	sqlDB, err := db.NewDB(cfg.Database)
// 	if err != nil {
// 		log.Fatalf("Failed to connect database: %v", err)
// 	}
// 	defer db.CloseDB(sqlDB)
// 	log.Printf("Database connected")

// 	if err := db.AutoMigrate(sqlDB); err != nil {
// 		log.Fatalf("Failed to migrate database: %v", err)
// 	}
// 	log.Printf("Database migrated")

// 	redisClient := redis.NewClient(cfg.Redis)
// 	if err := redisClient.Ping(nil); err != nil {
// 		log.Fatalf("Failed to connect Redis: %v", err)
// 	}
// 	defer redisClient.Close()
// 	log.Printf("Redis connected")

// 	accountRepo := account.NewAccountRepository(sqlDB)
// 	accountService := account.NewAccountService(accountRepo, redisClient)
// 	accountHandler := account.NewAccountHandler(accountService)
// 	log.Printf("Account module initialized")

//		log.Printf("Server is running on port %d", cfg.Server.Port)
//		log.Fatal("Server started")
//	}

// package main

// import (
// 	"context"
// 	"feedsystem_video_go/internal/account"
// 	"feedsystem_video_go/internal/auth"
// 	"feedsystem_video_go/internal/config"
// 	"feedsystem_video_go/internal/db"
// 	"feedsystem_video_go/internal/middleware/redis"
// 	"fmt"
// 	"log"
// 	"os"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	configPath := os.Getenv("CONFIG_PATH")
// 	if configPath == "" {
// 		configPath = "configs/config.yaml"
// 	}
// 	log.Printf("Loading config from %s", configPath)
// 	cfg, usedDefault, err := config.LoadLocalDev(configPath)
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}
// 	if usedDefault {
// 		log.Printf("Config File %s not found, using default local config", configPath)
// 	} else {
// 		log.Printf("Config loaded from file: %s", configPath)
// 	}

// 	sqlDB, err := db.NewDB(cfg.Database)
// 	if err != nil {
// 		log.Fatalf("Failed to connect database: %v", err)
// 	}
// 	defer db.CloseDB(sqlDB)
// 	log.Printf("Database connected")

// 	if err := db.AutoMigrate(sqlDB); err != nil {
// 		log.Fatalf("Failed to migrate database: %v", err)
// 	}
// 	log.Printf("Database migrated")

// 	redisClient := redis.NewClient(cfg.Redis)
// 	if err := redisClient.Ping(context.Background()); err != nil {
// 		log.Fatalf("Failed to connect Redis: %v", err)
// 	}
// 	defer redisClient.Close()
// 	log.Printf("Redis connected")

// 	accountRepo := account.NewAccountRepository(sqlDB)
// 	accountService := account.NewAccountService(accountRepo, redisClient)
// 	accountHandler := account.NewAccountHandler(accountService)
// 	log.Printf("Account module initialized")

// 	// 创建 Gin 路由器
// 	r := gin.Default()

// 	// 注册账户相关路由
// 	accountGroup := r.Group("/api/accounts")
// 	{
// 		accountGroup.POST("/register", accountHandler.CreateAccount)
// 		accountGroup.POST("/login", accountHandler.Login)
// 		accountGroup.GET("/:id", accountHandler.FindByID)
// 		accountGroup.GET("/username/:username", accountHandler.FindByUsername)

// 		// 需要认证的接口（添加 JWTMiddleware）
// 		accountGroup.POST("/logout", auth.JWTMiddleware(), accountHandler.Logout)
// 		accountGroup.POST("/refresh", accountHandler.Refresh)
// 		accountGroup.POST("/rename", auth.JWTMiddleware(), accountHandler.Rename)
// 		accountGroup.POST("/change-password", auth.JWTMiddleware(), accountHandler.ChangePassword)
// 	}

//		// 启动 HTTP 服务器
//		log.Printf("Server is running on port %d", cfg.Server.Port)
//		if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
//			log.Fatalf("Failed to start server: %v", err)
//		}
//	}

package main

import (
	"context"
	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/auth"
	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/db"
	"feedsystem_video_go/internal/like"
	"feedsystem_video_go/internal/middleware/redis"
	"feedsystem_video_go/internal/video"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	log.Printf("Loading config from %s", configPath)
	cfg, usedDefault, err := config.LoadLocalDev(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if usedDefault {
		log.Printf("Config File %s not found, using default local config", configPath)
	} else {
		log.Printf("Config loaded from file: %s", configPath)
	}

	// 连接数据库
	sqlDB, err := db.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.CloseDB(sqlDB)
	log.Printf("Database connected")

	// 自动迁移
	if err := db.AutoMigrate(sqlDB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Printf("Database migrated")

	// 连接Redis
	redisClient := redis.NewClient(cfg.Redis)
	if err := redisClient.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}
	defer redisClient.Close()
	log.Printf("Redis connected")

	// 初始化账户模块
	accountRepo := account.NewAccountRepository(sqlDB)
	accountService := account.NewAccountService(accountRepo, redisClient)
	accountHandler := account.NewAccountHandler(accountService)
	log.Printf("Account module initialized")

	// 初始化视频模块
	videoRepo := video.NewVideoRepository(sqlDB)
	uploadService := video.NewUploadService(videoRepo, cfg.Storage.UploadDir, cfg.Storage.BaseURL)
	videoService := video.NewVideoService(videoRepo, accountRepo, uploadService, cfg.Storage.BaseURL)
	videoHandler := video.NewVideoHandler(videoService)
	likeRepo := like.NewLikeRepository(sqlDB)
	likeService := like.NewLikeService(likeRepo, videoRepo, redisClient)
	likeHandler := like.NewLikeHandler(likeService)
	log.Printf("Video module initialized")

	r := gin.Default()

	accountGroup := r.Group("/api/accounts")
	{
		accountGroup.POST("/register", accountHandler.CreateAccount)
		accountGroup.POST("/login", accountHandler.Login)
		accountGroup.GET("/:id", accountHandler.FindByID)
		accountGroup.GET("/username/:username", accountHandler.FindByUsername)

		accountGroup.POST("/logout", auth.JWTMiddleware(), accountHandler.Logout)
		accountGroup.POST("/refresh", accountHandler.Refresh)
		accountGroup.POST("/rename", auth.JWTMiddleware(), accountHandler.Rename)
		accountGroup.POST("/change-password", auth.JWTMiddleware(), accountHandler.ChangePassword)
	}

	videoGroup := r.Group("/api/videos")
	{
		videoGroup.POST("/upload", auth.JWTMiddleware(), videoHandler.UploadVideo)
		videoGroup.GET("/:video_id", videoHandler.GetVideo)
		videoGroup.GET("/account/:account_id", videoHandler.ListVideos)
		videoGroup.DELETE("/:video_id", auth.JWTMiddleware(), videoHandler.DeleteVideo)
	}

	likeGroup := r.Group("/api/likes")
	{
		likeGroup.POST("/:video_id", auth.JWTMiddleware(), likeHandler.LikeVideo)
		likeGroup.DELETE("/:video_id", auth.JWTMiddleware(), likeHandler.UnlikeVideo)
		likeGroup.GET("/:video_id", likeHandler.GetLikeStatus)
		likeGroup.GET("/account/:account_id", likeHandler.ListLikes)
	}
	//启动服务
	log.Printf("Server is running on port %d", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
