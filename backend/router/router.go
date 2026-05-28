package router

import (
	"qiniudemo/handler"
	"qiniudemo/middleware"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.CORS(),
	)

	api := r.Group("/api/v1")
	{
		api.GET("/health", handler.Health)
	}

	return r
}
