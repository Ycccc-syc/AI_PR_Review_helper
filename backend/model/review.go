package model

// 前端请求体
type ReviewRequest struct {
	PRURL string `json:"pr_url" binding:"required"`
}

// 返回给前端的完整结果（总结 + 风险 + 建议 + 统计）
type ReviewResponse struct {
	Summary     string       `json:"summary"`
	RiskItems   []RiskItem   `json:"risk_items"`
	Suggestions []Suggestion `json:"suggestions"`
	Stats       ReviewStats  `json:"stats"`
}

// 单条风险项，含严重度/文件/行号/标题/详情/修复建议
type RiskItem struct {
	Severity   string `json:"severity"` // high / medium / low
	File       string `json:"file"`
	Line       int    `json:"line"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Suggestion string `json:"suggestion"`
}

// 单条改进建议，含优先级和内容
type Suggestion struct {
	Priority string `json:"priority"` // high / medium / low
	Content  string `json:"content"`
}

// PR 基本统计（文件数/增删行数/各级风险数量）
type ReviewStats struct {
	FilesChanged int `json:"files_changed"`
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
	RiskHigh     int `json:"risk_high"`
	RiskMedium   int `json:"risk_medium"`
	RiskLow      int `json:"risk_low"`
}
