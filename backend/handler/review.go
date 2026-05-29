package handler

import (
	"database/sql"
	"fmt"

	"qiniudemo/config"
	"qiniudemo/logger"
	"qiniudemo/model"
	"qiniudemo/repository"
	"qiniudemo/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Review(cfg *config.Config, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.ReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.Error(c, -1, "请提供有效的 pr_url 参数")
			return
		}

		// 解析 PR 信息（先于 LLM 调用，便于入库时拿到 title）
		pr, _ := repository.ParsePRURL(req.PRURL)
		prTitle := ""
		if pr != nil {
			prTitle = pr.Repo + " #" + formatPRNumber(pr.Number)
		}

		result, err := service.AnalyzePR(cfg, req.PRURL)
		if err != nil {
			model.Error(c, -1, err.Error())
			return
		}

		// 自动入库（仅已登录用户）
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(int64); ok {
				riskCount := len(result.RiskItems)
				_, saveErr := repository.SaveReviewLog(db, uid, req.PRURL, prTitle, result.Summary, riskCount)
				if saveErr != nil {
					logger.Log.Warn("save review log failed", zap.Error(saveErr))
				}
			}
		}

		model.Success(c, result)
	}
}

func formatPRNumber(n int) string {
	return fmt.Sprintf("%d", n)
}
