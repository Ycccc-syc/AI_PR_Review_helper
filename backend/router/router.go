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

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg))
		{
			protected.POST("/review", handler.Review(cfg, db))
			protected.GET("/reviews", handler.Reviews(db))
		}
	}

	return r
}
