package model

import "time"

type ReviewLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PRURL     string    `json:"pr_url"`
	PRTitle   string    `json:"pr_title"`
	Summary   string    `json:"summary"`
	RiskCount int       `json:"risk_count"`
	ResultJSON string   `json:"result_json,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ReviewDetail 历史记录详情（含完整分析结果）
type ReviewDetail struct {
	ReviewLog
	Result *ReviewResponse `json:"result,omitempty"`
}
