package repository

import (
	"database/sql"
	"fmt"

	"qiniudemo/model"
)

func SaveReviewLog(db *sql.DB, userID int64, prURL, prTitle, summary string, riskCount int) (*model.ReviewLog, error) {
	var log model.ReviewLog
	err := db.QueryRow(
		`INSERT INTO review_logs (user_id, pr_url, pr_title, summary, risk_count)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, user_id, pr_url, pr_title, summary, risk_count, created_at`,
		userID, prURL, prTitle, summary, riskCount,
	).Scan(&log.ID, &log.UserID, &log.PRURL, &log.PRTitle, &log.Summary, &log.RiskCount, &log.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("保存审查记录失败: %w", err)
	}
	return &log, nil
}

func GetReviewLogsByUserID(db *sql.DB, userID int64, limit, offset int) ([]model.ReviewLog, error) {
	rows, err := db.Query(
		`SELECT id, user_id, pr_url, pr_title, summary, risk_count, created_at
		 FROM review_logs WHERE user_id=$1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("查询审查记录失败: %w", err)
	}
	defer rows.Close()

	var logs []model.ReviewLog
	for rows.Next() {
		var l model.ReviewLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.PRURL, &l.PRTitle, &l.Summary, &l.RiskCount, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("解析审查记录失败: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, nil
}
