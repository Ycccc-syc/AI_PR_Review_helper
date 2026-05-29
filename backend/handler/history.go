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
		userID, _ := c.Get("user_id")
		uid := userID.(int64)

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
