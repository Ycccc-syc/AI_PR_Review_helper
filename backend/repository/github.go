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

var prURLPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/pull/(\d+)/?$`)

type PRInfo struct {
	Owner         string
	Repo          string
	Number        int
	Title         string
	Description   string
	FilesChanged  int
	Additions     int
	Deletions     int
}

type githubPRResponse struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	ChangedFiles int `json:"changed_files"`
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
}

func ParsePRURL(raw string) (*PRInfo, error) {
	raw = strings.TrimSpace(raw)
	matches := prURLPattern.FindStringSubmatch(raw)
	if matches == nil {
		return nil, fmt.Errorf("invalid GitHub PR URL: %s", raw)
	}
	num, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid PR number: %s", matches[3])
	}
	return &PRInfo{
		Owner:  matches[1],
		Repo:   matches[2],
		Number: num,
	}, nil
}

func FetchPRMeta(token string, pr *PRInfo) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", pr.Owner, pr.Repo, pr.Number)
	body, err := githubGet(token, url)
	if err != nil {
		return err
	}
	var resp githubPRResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("failed to parse PR metadata: %w", err)
	}
	pr.Title = resp.Title
	pr.Description = resp.Body
	pr.FilesChanged = resp.ChangedFiles
	pr.Additions = resp.Additions
	pr.Deletions = resp.Deletions
	return nil
}

func FetchPRDiff(token string, pr *PRInfo) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d", pr.Owner, pr.Repo, pr.Number)
	return githubGet(token, url, "application/vnd.github.v3.diff")
}

func githubGet(token string, url string, extraHeaders ...string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
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
		return "", fmt.Errorf("GitHub API request failed: %w", err)
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
		return "", fmt.Errorf("GitHub API rate limit exceeded, try again later")
	}
	if resp.StatusCode == 404 {
		return "", fmt.Errorf("PR not found (404), check the URL and make sure the repo is not private")
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return "", fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	return string(body), nil
}
