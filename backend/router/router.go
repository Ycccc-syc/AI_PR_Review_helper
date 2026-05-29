package router

import (
	"qiniudemo/config"
	"qiniudemo/handler"
	"qiniudemo/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.CORS(),
	)

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.Health)
		api.POST("/review", handler.Review(cfg))
	}

	return r
}
