package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"qiniudemo/config"
	"qiniudemo/llm"
	"qiniudemo/logger"
	"qiniudemo/model"
	"qiniudemo/prompt"
	"qiniudemo/repository"

	"go.uber.org/zap"
)

const maxDiffLines = 8000

// AnalyzePR 完整 PR 分析流程：解析 URL → 拉取数据 → LLM 分析 → 结构化结果
func AnalyzePR(cfg *config.Config, prURL string) (*model.ReviewResponse, error) {
	pr, err := repository.ParsePRURL(prURL)
	if err != nil {
		return nil, fmt.Errorf("解析 PR URL 失败: %w", err)
	}

	if err := repository.FetchPRMeta(cfg.GitHubToken, pr); err != nil {
		return nil, fmt.Errorf("获取 PR 信息失败: %w", err)
	}

	diff, err := repository.FetchPRDiff(cfg.GitHubToken, pr)
	if err != nil {
		return nil, fmt.Errorf("获取 PR diff 失败: %w", err)
	}

	if diff == "" {
		return nil, fmt.Errorf("该 PR 没有代码变更")
	}

	pr.Description = formatDescription(pr.Description)

	var result model.ReviewResponse
	if lines := strings.Count(diff, "\n"); lines > maxDiffLines {
		result, err = analyzeChunked(cfg, pr, diff)
	} else {
		result, err = analyzeSingle(cfg, pr, diff)
	}
	if err != nil {
		return nil, err
	}

	result.Stats = buildStats(pr, result.RiskItems)
	return &result, nil
}

func analyzeSingle(cfg *config.Config, pr *repository.PRInfo, diff string) (model.ReviewResponse, error) {
	messages := prompt.BuildMessages(pr.Title, pr.Description, diff)
	return callAndParse(cfg, messages)
}

func analyzeChunked(cfg *config.Config, pr *repository.PRInfo, diff string) (model.ReviewResponse, error) {
	chunks := splitDiffByFile(diff, 5)
	if len(chunks) == 0 {
		return analyzeSingle(cfg, pr, diff)
	}

	var merged model.ReviewResponse
	for i, chunk := range chunks {
		chunkTitle := fmt.Sprintf("%s (第 %d/%d 部分)", pr.Title, i+1, len(chunks))
		messages := prompt.BuildMessages(chunkTitle, pr.Description, chunk)
		r, err := callAndParse(cfg, messages)
		if err != nil {
			logger.Log.Warn("chunk analysis failed", zap.Int("chunk", i), zap.Error(err))
			continue
		}
		if i == 0 {
			merged.Summary = r.Summary
		}
		merged.RiskItems = append(merged.RiskItems, r.RiskItems...)
		merged.Suggestions = append(merged.Suggestions, r.Suggestions...)
	}
	return merged, nil
}

func callAndParse(cfg *config.Config, messages []llm.Message) (model.ReviewResponse, error) {
	raw, err := llm.Call(cfg.DeepSeekAPIKey, cfg.DeepSeekModel, messages)
	if err != nil {
		return model.ReviewResponse{}, err
	}

	jsonStr := extractJSON(raw)

	var result model.ReviewResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		logger.Log.Error("解析 LLM JSON 失败", zap.String("raw", raw[:min(len(raw), 500)]))
		return model.ReviewResponse{}, fmt.Errorf("解析分析结果失败，模型返回格式异常: %w", err)
	}
	return result, nil
}

func buildStats(pr *repository.PRInfo, risks []model.RiskItem) model.ReviewStats {
	s := model.ReviewStats{
		FilesChanged: pr.FilesChanged,
		Additions:    pr.Additions,
		Deletions:    pr.Deletions,
	}
	for _, r := range risks {
		switch r.Severity {
		case "high":
			s.RiskHigh++
		case "medium":
			s.RiskMedium++
		case "low":
			s.RiskLow++
		}
	}
	return s
}

func extractJSON(raw string) string {
	raw = strings.TrimSpace(raw)

	if i := strings.Index(raw, "```json"); i != -1 {
		raw = raw[i+7:]
		if j := strings.Index(raw, "```"); j != -1 {
			raw = raw[:j]
		}
	} else if strings.HasPrefix(raw, "```") {
		raw = raw[3:]
		if j := strings.Index(raw, "```"); j != -1 {
			raw = raw[:j]
		}
	}

	if i := strings.Index(raw, "{"); i != -1 {
		raw = raw[i:]
	}
	if j := strings.LastIndex(raw, "}"); j != -1 {
		raw = raw[:j+1]
	}

	return strings.TrimSpace(raw)
}

func formatDescription(body string) string {
	if body == "" {
		return "（无描述）"
	}
	if len(body) > 500 {
		return body[:500] + "..."
	}
	return body
}

func splitDiffByFile(diff string, filesPerChunk int) []string {
	parts := strings.Split(diff, "\ndiff --git ")
	if len(parts) <= 1 || len(parts) <= filesPerChunk {
		return []string{diff}
	}

	// parts[0] is empty or leading text, skip it
	fileBlocks := make([]string, 0, len(parts))
	for _, p := range parts[1:] {
		fileBlocks = append(fileBlocks, "diff --git "+p)
	}

	var chunks []string
	for i := 0; i < len(fileBlocks); i += filesPerChunk {
		end := i + filesPerChunk
		if end > len(fileBlocks) {
			end = len(fileBlocks)
		}
		chunks = append(chunks, strings.Join(fileBlocks[i:end], "\n"))
	}
	return chunks
}
