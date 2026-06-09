package router

import (
	"feedsystem_video_go/internal/account"
	"feedsystem_video_go/internal/comment"
	"feedsystem_video_go/internal/feed"
	"feedsystem_video_go/internal/like"
	"feedsystem_video_go/internal/message"
	"feedsystem_video_go/internal/middleware"
	"feedsystem_video_go/internal/search"
	"feedsystem_video_go/internal/social"
	"feedsystem_video_go/internal/sse"
	"feedsystem_video_go/internal/video"

	"github.com/gin-gonic/gin"
)

type Router struct {
	accountHandler *account.AccountHandler
	videoHandler   *video.VideoHandler
	likeHandler    *like.LikeHandler
	commentHandler *comment.CommentHandler
	feedHandler    *feed.FeedHandler
	socialHandler  *social.SocialHandler
	messageHandler *message.MessageHandler
	searchHandler  *search.SearchHandler
	sseHub         *sse.Hub
}

func NewRouter(
	accountHandler *account.AccountHandler,
	videoHandler *video.VideoHandler,
	likeHandler *like.LikeHandler,
	commentHandler *comment.CommentHandler,
	feedHandler *feed.FeedHandler,
	socialHandler *social.SocialHandler,
	messageHandler *message.MessageHandler,
	searchHandler *search.SearchHandler,
	sseHub *sse.Hub,
) *Router {
	return &Router{
		accountHandler: accountHandler,
		videoHandler:   videoHandler,
		likeHandler:    likeHandler,
		commentHandler: commentHandler,
		feedHandler:    feedHandler,
		socialHandler:  socialHandler,
		messageHandler: messageHandler,
		searchHandler:  searchHandler,
		sseHub:         sseHub,
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	api := engine.Group("/api/v1")
	api.Use(middleware.CORS())
	api.Use(middleware.AuthMiddleware())

	account := api.Group("/account")
	{
		account.POST("/register", r.accountHandler.CreateAccount)
		account.POST("/login", r.accountHandler.Login)
		account.POST("/refresh", r.accountHandler.Refresh)
		account.GET("/search", r.accountHandler.SearchUsers)

		account.Use(middleware.RequireAuth())
		{
			account.GET("/:id", r.accountHandler.FindByID)
			account.PUT("/rename", r.accountHandler.Rename)
			account.PUT("/password", r.accountHandler.ChangePassword)
			account.PUT("/profile", r.accountHandler.UpdateProfile)
			account.POST("/logout", r.accountHandler.Logout)
			account.POST("/avatar", r.accountHandler.UploadAvatar)
		}
	}

	video := api.Group("/video")
	{
		video.GET("/:id", r.videoHandler.GetVideo)
		video.GET("/user/:account_id", r.videoHandler.ListVideos)
		video.POST("/:id/view", r.videoHandler.ReportView)

		video.Use(middleware.RequireAuth())
		{
			video.POST("/upload", r.videoHandler.UploadVideo)
			video.DELETE("/:id", r.videoHandler.DeleteVideo)
		}
	}

	like := api.Group("/like")
	like.Use(middleware.RequireAuth())
	{
		like.POST("/:video_id", r.likeHandler.LikeVideo)
		like.DELETE("/:video_id", r.likeHandler.UnlikeVideo)
		like.GET("/:video_id", r.likeHandler.GetLikeStatus)
		like.GET("/list", r.likeHandler.ListLikes)
	}

	comment := api.Group("/comment")
	{
		comment.GET("/:video_id/list", r.commentHandler.ListComments)

		comment.Use(middleware.RequireAuth())
		{
			comment.POST("/:video_id", r.commentHandler.CreateComment)
		}
	}
	commentDelete := api.Group("/comment")
	commentDelete.Use(middleware.RequireAuth())
	{
		commentDelete.DELETE("/:comment_id", r.commentHandler.DeleteComment)
	}

	feed := api.Group("/feed")
	{
		feed.GET("/", r.feedHandler.GetFeed)
		feed.GET("/hot", r.feedHandler.GetHotFeed)
		feed.GET("/tag/:tag", r.feedHandler.GetTagFeed)
		feed.GET("/search", r.feedHandler.SearchFeed)

		feed.Use(middleware.RequireAuth())
		{
			feed.GET("/following", r.feedHandler.GetFollowingFeed)
		}
	}

	social := api.Group("/social")
	{
		social.GET("/profile/:account_id", r.socialHandler.GetProfile)
		social.GET("/followers/:target_id", r.socialHandler.GetFollowers)
		social.GET("/following/:account_id", r.socialHandler.GetFollowing)

		social.Use(middleware.RequireAuth())
		{
			social.POST("/follow/:target_id", r.socialHandler.Follow)
			social.DELETE("/unfollow/:target_id", r.socialHandler.Unfollow)
			social.GET("/friends/search", r.socialHandler.SearchFriends)
		}
	}

	message := api.Group("/message")
	message.Use(middleware.RequireAuth())
	{
		message.POST("/send", r.messageHandler.SendMessage)
		message.GET("/conversations", r.messageHandler.GetConversations)
		message.GET("/:other_id", r.messageHandler.GetMessages)
	}

	search := api.Group("/search")
	{
		search.GET("/hot", r.searchHandler.GetHotSearches)

		search.Use(middleware.RequireAuth())
		{
			search.POST("/record", r.searchHandler.RecordSearch)
			search.GET("/history", r.searchHandler.GetSearchHistory)
			search.DELETE("/history", r.searchHandler.ClearSearchHistory)
		}
	}

	api.GET("/sse", r.sseHub.ServeSSE)
}
