package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"qiniudemo/model"
)

func SaveReviewLog(db *sql.DB, userID int64, prURL, prTitle, summary string, riskCount int, result *model.ReviewResponse) (*model.ReviewLog, error) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("序列化审查结果失败: %w", err)
	}

	var log model.ReviewLog
	err = db.QueryRow(
		`INSERT INTO review_logs (user_id, pr_url, pr_title, summary, risk_count, result_json)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, pr_url, pr_title, summary, risk_count, created_at`,
		userID, prURL, prTitle, summary, riskCount, string(resultJSON),
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

// GetReviewLogByID 获取单条记录的完整详情（含 result_json）
func GetReviewLogByID(db *sql.DB, logID, userID int64) (*model.ReviewDetail, error) {
	var l model.ReviewLog
	err := db.QueryRow(
		`SELECT id, user_id, pr_url, pr_title, summary, risk_count, result_json, created_at
		 FROM review_logs WHERE id=$1 AND user_id=$2`,
		logID, userID,
	).Scan(&l.ID, &l.UserID, &l.PRURL, &l.PRTitle, &l.Summary, &l.RiskCount, &l.ResultJSON, &l.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("记录不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("查询审查记录失败: %w", err)
	}

	detail := &model.ReviewDetail{ReviewLog: l}
	if l.ResultJSON != "" {
		var result model.ReviewResponse
		if err := json.Unmarshal([]byte(l.ResultJSON), &result); err == nil {
			detail.Result = &result
		}
	}
	return detail, nil
}
