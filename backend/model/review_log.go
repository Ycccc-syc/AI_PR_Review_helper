package model

import "time"

type ReviewLog struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PRURL     string    `json:"pr_url"`
	PRTitle   string    `json:"pr_title"`
	Summary   string    `json:"summary"`
	RiskCount int       `json:"risk_count"`
	CreatedAt time.Time `json:"created_at"`
}
