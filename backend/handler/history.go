package handler

import (
	"database/sql"
	"strconv"

	"qiniudemo/model"
	"qiniudemo/repository"

	"github.com/gin-gonic/gin"
)

func Reviews(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			model.Error(c, 401, "未登录")
			return
		}
		uid, ok := raw.(int64)
		if !ok {
			model.Error(c, 401, "令牌无效")
			return
		}

		limit := 20
		offset := 0
		if l := c.Query("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 50 {
				limit = n
			}
		}
		if o := c.Query("offset"); o != "" {
			if n, err := strconv.Atoi(o); err == nil && n >= 0 {
				offset = n
			}
		}

		logs, err := repository.GetReviewLogsByUserID(db, uid, limit, offset)
		if err != nil {
			model.Error(c, -1, err.Error())
			return
		}
		if logs == nil {
			logs = []model.ReviewLog{}
		}

		model.Success(c, logs)
	}
}

// ReviewDetail 获取单条历史记录的完整详情（含 risk_items/suggestions/health_score）
func ReviewDetail(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, exists := c.Get("user_id")
		if !exists {
			model.Error(c, 401, "未登录")
			return
		}
		uid, ok := raw.(int64)
		if !ok {
			model.Error(c, 401, "令牌无效")
			return
		}

		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			model.Error(c, -1, "无效的记录 ID")
			return
		}

		detail, err := repository.GetReviewLogByID(db, id, uid)
		if err != nil {
			model.Error(c, -1, err.Error())
			return
		}

		model.Success(c, detail)
	}
}
