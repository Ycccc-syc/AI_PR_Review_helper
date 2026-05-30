package handler

import (
	"database/sql"
	"encoding/json"
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
				_, saveErr := repository.SaveReviewLog(db, uid, req.PRURL, prTitle, result.Summary, riskCount, result)
				if saveErr != nil {
					logger.Log.Warn("save review log failed", zap.Error(saveErr))
				}
			} else {
				logger.Log.Warn("user_id type assertion failed in Review", zap.Any("type", userID))
			}
		}

		model.Success(c, result)
	}
}

func formatPRNumber(n int) string {
	return fmt.Sprintf("%d", n)
}

// ReviewStream SSE 流式审查，实时推送 LLM 分析文本
func ReviewStream(cfg *config.Config, db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.ReviewRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.Error(c, -1, "请提供有效的 pr_url 参数")
			return
		}

		// 在设置 SSE 响应头之前校验，避免发送 JSON 错误时响应头已提交
		flusher, ok := c.Writer.(gin.ResponseWriter)
		if !ok {
			model.Error(c, -1, "SSE 不支持当前环境")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Writer.WriteHeader(200)

		result, err := service.AnalyzePRStream(cfg, req.PRURL, func(chunk string) {
			data, _ := json.Marshal(map[string]interface{}{
				"type":    "chunk",
				"content": chunk,
			})
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
		})

		if err != nil {
			data, _ := json.Marshal(map[string]string{
				"type":    "error",
				"message": err.Error(),
			})
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
			return
		}

		// 自动入库（仅已登录用户）
		if userID, exists := c.Get("user_id"); exists {
			if uid, ok := userID.(int64); ok {
				pr, _ := repository.ParsePRURL(req.PRURL)
				prTitle := ""
				if pr != nil {
					prTitle = pr.Repo + " #" + formatPRNumber(pr.Number)
				}
				riskCount := len(result.RiskItems)
				_, saveErr := repository.SaveReviewLog(db, uid, req.PRURL, prTitle, result.Summary, riskCount, result)
				if saveErr != nil {
					logger.Log.Warn("save review log failed", zap.Error(saveErr))
				}
			} else {
				logger.Log.Warn("user_id type assertion failed", zap.Any("type", userID))
			}
		}

		data, _ := json.Marshal(map[string]interface{}{
			"type": "result",
			"data": result,
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	}
}

// AskStream SSE 流式追问接口
func AskStream(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.AskRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.Error(c, -1, "请提供 pr_url 和 question 参数")
			return
		}

		flusher, ok := c.Writer.(gin.ResponseWriter)
		if !ok {
			model.Error(c, -1, "SSE 不支持当前环境")
			return
		}

		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Writer.WriteHeader(200)

		answer, err := service.AskFollowUpStream(cfg, req.PRURL, req.Question, req.Context, func(chunk string) {
			data, _ := json.Marshal(map[string]interface{}{
				"type":    "chunk",
				"content": chunk,
			})
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
		})

		if err != nil {
			data, _ := json.Marshal(map[string]string{
				"type":    "error",
				"message": err.Error(),
			})
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
			return
		}

		data, _ := json.Marshal(map[string]interface{}{
			"type":   "result",
			"answer": answer,
		})
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	}
}
