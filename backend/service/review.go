package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"qiniudemo/config"
	"qiniudemo/llm"
	"qiniudemo/logger"
	"qiniudemo/model"
	"qiniudemo/prompt"
	"qiniudemo/repository"

	"go.uber.org/zap"
)

const maxDiffLines = 8000

// 需要跳过的文件后缀和文件名
var skipPatterns = []string{
	// lock 文件
	"package-lock.json", "yarn.lock", "pnpm-lock.yaml",
	"go.sum", "Cargo.lock", "Gemfile.lock", "composer.lock",
	// 二进制/资源文件
	".png", ".jpg", ".jpeg", ".gif", ".ico", ".svg",
	".woff", ".woff2", ".ttf", ".eot", ".otf",
	".pdf", ".zip", ".tar", ".gz", ".7z",
	// 生成文件
	".exe", ".dll", ".so", ".dylib", ".wasm",
	// 大型数据文件
	".sql", ".csv", ".tsv",
}

// filterDiff 去掉 lock 文件、二进制文件、资源文件的 diff
func filterDiff(diff string) string {
	parts := strings.Split(diff, "\ndiff --git ")
	if len(parts) <= 1 {
		return diff
	}

	kept := []string{parts[0]} // diff 开头的内容（如有）
	for _, block := range parts[1:] {
		// 提取文件路径: diff --git a/path/to/file b/path/to/file
		firstLine := "diff --git " + block
		fileName := extractFileName(firstLine)
		if shouldSkip(fileName) {
			continue
		}
		kept = append(kept, "diff --git "+block)
	}

	if len(kept) <= 1 {
		return diff
	}
	return strings.Join(kept, "\n")
}

// extractFileName 从 diff 块首行提取文件名
func extractFileName(firstLine string) string {
	lower := strings.ToLower(firstLine)
	return lower
}

// shouldSkip 判断文件是否应跳过分析
func shouldSkip(fileName string) bool {
	for _, p := range skipPatterns {
		if strings.Contains(fileName, p) {
			return true
		}
	}
	return false
}

// AnalyzePR 完整 PR 分析流程：解析 URL → 拉取数据 → 查缓存 → LLM 分析 → 写缓存 → 结构化结果
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

	diff = filterDiff(diff)
	if diff == "" {
		return nil, fmt.Errorf("该 PR 的文件变更已被过滤（可能仅包含 lock 文件或二进制文件）")
	}

	// 缓存检查：同一 PR 同一 diff 不重复调用 LLM
	cacheKey := repository.CacheKey(pr.Owner, pr.Repo, pr.Number, diff)
	if cached, ok := repository.GetCachedResult(cacheKey); ok {
		var result model.ReviewResponse
		if err := json.Unmarshal(cached, &result); err == nil {
			logger.Log.Info("review cache hit", zap.String("key", cacheKey))
			result.Stats = buildStats(pr, result.RiskItems)
			return &result, nil
		}
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
		repository.SetCachedResult(cacheKey, result, 1*time.Hour)
	return &result, nil
}

// AnalyzePRStream 流式 PR 分析，onChunk 每收到一段 LLM 输出就回调（便于 SSE 推送）
func AnalyzePRStream(cfg *config.Config, prURL string, onChunk func(string)) (*model.ReviewResponse, error) {
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

	diff = filterDiff(diff)
	if diff == "" {
		return nil, fmt.Errorf("该 PR 的文件变更已被过滤")
	}

	pr.Description = formatDescription(pr.Description)

	// 大 PR 降级为非流式分片分析
	if strings.Count(diff, "\n") > maxDiffLines {
		return AnalyzePR(cfg, prURL)
	}

	cacheKey := repository.CacheKey(pr.Owner, pr.Repo, pr.Number, diff)
	if cached, ok := repository.GetCachedResult(cacheKey); ok {
		var result model.ReviewResponse
		if err := json.Unmarshal(cached, &result); err == nil {
			logger.Log.Info("review cache hit (stream)", zap.String("key", cacheKey))
			result.Stats = buildStats(pr, result.RiskItems)
			return &result, nil
		}
	}

	messages := prompt.BuildMessages(pr.Title, pr.Description, diff)

	var fullText strings.Builder
	for chunk := range llm.CallStream(cfg.DeepSeekAPIKey, cfg.DeepSeekModel, messages) {
		if chunk.Error != nil {
			return nil, chunk.Error
		}
		if chunk.Done {
			break
		}
		fullText.WriteString(chunk.Content)
		if onChunk != nil {
			onChunk(chunk.Content)
		}
	}

	jsonStr := extractJSON(fullText.String())
	var result model.ReviewResponse
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		logger.Log.Error("解析 LLM JSON 失败", zap.String("raw", fullText.String()[:min(len(fullText.String()), 500)]))
		return nil, fmt.Errorf("解析分析结果失败，模型返回格式异常: %w", err)
	}

	result.Stats = buildStats(pr, result.RiskItems)
		repository.SetCachedResult(cacheKey, result, 1*time.Hour)
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

// AskFollowUpStream 流式追问，onChunk 推送文本片段，返回完整回答
func AskFollowUpStream(cfg *config.Config, prURL, question, contextJSON string, onChunk func(string)) (string, error) {
	pr, err := repository.ParsePRURL(prURL)
	if err != nil {
		return "", fmt.Errorf("解析 PR URL 失败: %w", err)
	}

	if err := repository.FetchPRMeta(cfg.GitHubToken, pr); err != nil {
		return "", fmt.Errorf("获取 PR 信息失败: %w", err)
	}

	diff, _ := repository.FetchPRDiff(cfg.GitHubToken, pr)
	diff = filterDiff(diff)

	messages := prompt.BuildAskMessages(pr.Title, pr.Description, diff, question, contextJSON)

	var fullText strings.Builder
	for chunk := range llm.CallStream(cfg.DeepSeekAPIKey, cfg.DeepSeekModel, messages) {
		if chunk.Error != nil {
			return "", chunk.Error
		}
		if chunk.Done {
			break
		}
		fullText.WriteString(chunk.Content)
		if onChunk != nil {
			onChunk(chunk.Content)
		}
	}
	return fullText.String(), nil
}
