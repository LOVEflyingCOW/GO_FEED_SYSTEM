package main

import (
	"context"
	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/auth"
	"feedsystem_video_go/internal/comment"
	"feedsystem_video_go/internal/config"
	"feedsystem_video_go/internal/db"
	"feedsystem_video_go/internal/feed"
	"feedsystem_video_go/internal/like"
	"feedsystem_video_go/internal/message"
	"feedsystem_video_go/internal/middleware/redis"
	mqcomment "feedsystem_video_go/internal/mq/comment"
	mqlike "feedsystem_video_go/internal/mq/like"
	mqpopularity "feedsystem_video_go/internal/mq/popularity"
	"feedsystem_video_go/internal/mq/rabbitmq"
	"feedsystem_video_go/internal/router"
	"feedsystem_video_go/internal/social"
	"feedsystem_video_go/internal/sse"
	"feedsystem_video_go/internal/video"
	"feedsystem_video_go/internal/worker"
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

	// 设置JWT密钥
	auth.SetJWTSecret(cfg.JWT.Secret)

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

	// 连接RabbitMQ
	rabbitMQClient, err := rabbitmq.NewClient(cfg.RabbitMQ)
	if err != nil {
		log.Printf("Warning: Failed to connect RabbitMQ: %v", err)
		rabbitMQClient = nil
	} else {
		defer rabbitMQClient.Close()
		log.Printf("RabbitMQ connected")
	}

	//初始化账户模块
	accountRepo := account.NewAccountRepository(sqlDB)
	accountService := account.NewAccountService(accountRepo, redisClient)
	accountHandler := account.NewAccountHandler(accountService)
	log.Printf("Account module initialized")

	//初始化视频模块
	videoRepo := video.NewVideoRepository(sqlDB)
	uploadService := video.NewUploadService(videoRepo, cfg.Storage.UploadDir, cfg.Storage.BaseURL)

	//初始化点赞模块（用于feed）
	likeRepo := like.NewLikeRepository(sqlDB)

	//初始化Feed模块（用于时间线）
	feedService := feed.NewFeedService(videoRepo, accountRepo, likeRepo, redisClient, cfg.Storage.BaseURL)

	//初始化视频服务（传入时间线服务）
	videoService := video.NewVideoService(videoRepo, accountRepo, uploadService, cfg.Storage.BaseURL, feedService)
	videoHandler := video.NewVideoHandler(videoService)
	log.Printf("Video module initialized")

	//初始化评论模块
	commentRepo := comment.NewCommentRepository(sqlDB)
	commentService := comment.NewCommentService(commentRepo, accountRepo, videoRepo)
	commentHandler := comment.NewCommentHandler(commentService)
	log.Printf("Comment module initialized")

	//初始化点赞模块（完整）
	likeService := like.NewLikeService(likeRepo, videoRepo, redisClient)
	likeHandler := like.NewLikeHandler(likeService)
	log.Printf("Like module initialized")

	//初始化社交模块（关注功能）
	socialRepo := social.NewSocialRepository(sqlDB)
	socialService := social.NewSocialService(socialRepo, accountRepo)
	socialHandler := social.NewSocialHandler(socialService)
	log.Printf("Social module initialized")

	//初始化SSE Hub（实时通知）
	sseHub := sse.NewHub()
	go sseHub.Run()
	log.Printf("SSE Hub initialized")

	//初始化消息模块
	messageRepo := message.NewMessageRepository(sqlDB)
	messageService := message.NewMessageService(messageRepo, accountRepo, sseHub)
	messageHandler := message.NewMessageHandler(messageService)
	log.Printf("Message module initialized")

	//初始化MQ模块
	var likeMQ *mqlike.LikeMQ
	var commentMQ *mqcomment.CommentMQ
	if rabbitMQClient != nil {
		likeMQ, err = mqlike.NewLikeMQ(rabbitMQClient)
		if err != nil {
			log.Printf("Warning: Failed to initialize LikeMQ: %v", err)
		}

		commentMQ, err = mqcomment.NewCommentMQ(rabbitMQClient)
		if err != nil {
			log.Printf("Warning: Failed to initialize CommentMQ: %v", err)
		}

		_, err = mqpopularity.NewPopularityMQ(rabbitMQClient)
		if err != nil {
			log.Printf("Warning: Failed to initialize PopularityMQ: %v", err)
		}
		log.Printf("Message queues initialized")

		// 启动Worker
		ctx := context.Background()

		// Like Worker
		if likeMQ != nil {
			likeWorker := worker.NewLikeWorker(likeService, likeMQ)
			go likeWorker.Run(ctx)
			log.Printf("LikeWorker started")
		}

		// Comment Worker
		if commentMQ != nil {
			commentWorker := worker.NewCommentWorker(sseHub, commentMQ)
			go commentWorker.Run(ctx)
			log.Printf("CommentWorker started")
		}
	}

	// 初始化Feed Handler
	feedHandler := feed.NewFeedHandler(feedService)
	log.Printf("Feed module initialized")

	// 创建路由器
	r := gin.Default()

	// 注册路由（使用router模块）
	router := router.NewRouter(
		accountHandler,
		videoHandler,
		likeHandler,
		commentHandler,
		feedHandler,
		socialHandler,
		messageHandler,
		sseHub,
	)
	router.RegisterRoutes(r)

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	//启动服务
	log.Printf("Server is running on port %d", cfg.Server.Port)
	if err := r.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
