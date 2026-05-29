package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// GitHub PR URL 解析正则: https://github.com/{owner}/{repo}/pull/{number}
var prURLPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/pull/(\d+)/?$`)

// PRInfo 从 PR URL 解析出的信息 + API 返回的元数据
type PRInfo struct {
	Owner        string
	Repo         string
	Number       int
	Title        string
	Description  string
	FilesChanged int
	Additions    int
	Deletions    int
}

// githubPRResponse GitHub API 返回的 PR JSON 结构
type githubPRResponse struct {
	Title        string `json:"title"`
	Body         string `json:"body"`
	ChangedFiles int    `json:"changed_files"`
	Additions    int    `json:"additions"`
	Deletions    int    `json:"deletions"`
}

// ParsePRURL 解析 GitHub PR URL，提取 owner/repo/number
func ParsePRURL(raw string) (*PRInfo, error) {
	raw = strings.TrimSpace(raw)
	matches := prURLPattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf("无效的 GitHub PR 链接: %s", raw)
	}
	num, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("无效的 PR 编号: %s", matches[3])
	}
	return &PRInfo{
		Owner:  matches[1],
		Repo:   matches[2],
		Number: num,
	}, nil
}

// FetchPRMeta 获取 PR 元数据（标题、描述、变更统计）
func FetchPRMeta(token string, pr *PRInfo) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", pr.Owner, pr.Repo, pr.Number)
	body, err := githubGet(token, url)
	if err != nil {
		return err
	}
	var resp githubPRResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return fmt.Errorf("解析 PR 元数据失败: %w", err)
	}
	pr.Title = resp.Title
	pr.Description = resp.Body
	pr.FilesChanged = resp.ChangedFiles
	pr.Additions = resp.Additions
	pr.Deletions = resp.Deletions
	return nil
}

// FetchPRDiff 获取 PR 的 unified diff 文本
func FetchPRDiff(token string, pr *PRInfo) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", pr.Owner, pr.Repo, pr.Number)
	return githubGet(token, url, "application/vnd.github.v3.diff")
}

// githubGet 封装 GitHub API 的 GET 请求，统一处理认证、限流、错误
func githubGet(token string, url string, extraHeaders ...string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "ai-pr-review")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	for _, h := range extraHeaders {
		req.Header.Set("Accept", h)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GitHub API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			sec, _ := strconv.Atoi(retryAfter)
			if sec > 0 && sec <= 60 {
				time.Sleep(time.Duration(sec) * time.Second)
				return githubGet(token, url, extraHeaders...)
			}
		}
		return "", fmt.Errorf("GitHub API 限流，请稍后重试")
	}
	if resp.StatusCode == 404 {
		return "", fmt.Errorf("PR 不存在 (404)，请检查链接或确认仓库不是私有")
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("GitHub API 返回 %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %w", err)
	}
	return string(body), nil
}
