package handler

import (
	"qiniudemo/config"
	"qiniudemo/model"
	"qiniudemo/service"

	"github.com/gin-gonic/gin"
)

func Review(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.ReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.Error(c, -1, "请提供有效的 pr_url 参数")
			return
		}

		result, err := service.AnalyzePR(cfg, req.PRURL)
		if err != nil {
			model.Error(c, -1, err.Error())
			return
		}

		model.Success(c, result)
	}
}
