package router

import (
	"database/sql"

	"qiniudemo/config"
	"qiniudemo/handler"
	"qiniudemo/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config, db *sql.DB) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.CORS(),
	)

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.Health)

		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register(db))
			auth.POST("/login", handler.Login(db, cfg))
		}

		// review 允许未登录访问，OptionalAuth 提取已登录用户的 ID 以便自动入库
		api.POST("/review", middleware.OptionalAuth(cfg), handler.Review(cfg, db))
		api.POST("/review/stream", middleware.OptionalAuth(cfg), handler.ReviewStream(cfg, db))

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg))
		{
			protected.GET("/reviews", handler.Reviews(db))
		}
	}

	return r
}
