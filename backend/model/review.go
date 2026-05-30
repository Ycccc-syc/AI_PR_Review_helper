package model

// 前端请求体
type ReviewRequest struct {
	PRURL string `json:"pr_url" binding:"required"`
}

// 追问请求：携带 PR URL、问题、风险项上下文
type AskRequest struct {
	PRURL    string `json:"pr_url" binding:"required"`
	Question string `json:"question" binding:"required"`
	Context  string `json:"context"` // 风险项的 JSON 摘要
}

// 返回给前端的完整结果（总结 + 风险 + 建议 + 统计 + 健康分）
type ReviewResponse struct {
	Summary     string       `json:"summary"`
	RiskItems   []RiskItem   `json:"risk_items"`
	Suggestions []Suggestion `json:"suggestions"`
	HealthScore HealthScore  `json:"health_score"`
	Stats       ReviewStats  `json:"stats"`
}

// PR 健康分（总分/安全/性能/可维护，0-100）
type HealthScore struct {
	Total           int `json:"total"`
	Security        int `json:"security"`
	Performance     int `json:"performance"`
	Maintainability int `json:"maintainability"`
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
