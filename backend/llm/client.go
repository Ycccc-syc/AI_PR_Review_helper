package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"qiniudemo/logger"

	"go.uber.org/zap"
)

const deepseekBaseURL = "https://api.deepseek.com"

// Message DeepSeek API 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest DeepSeek API 请求体（兼容 OpenAI 格式）
type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Choice DeepSeek API 返回的选项
type choice struct {
	Message Message `json:"message"`
}

// ChatResponse DeepSeek API 返回体
type chatResponse struct {
	Choices []choice `json:"choices"`
}

// Call DeepSeek Chat API，返回模型输出的文本
func Call(apiKey, model string, messages []Message) (string, error) {
	reqBody := chatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.3,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", deepseekBaseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Log.Error("LLM 请求失败", zap.Error(err))
		return "", fmt.Errorf("LLM API 请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(io.LimitReader(resp.Body, 10*1024*1024))
	if err != nil {
		return "", fmt.Errorf("读取 LLM 响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		logger.Log.Error("LLM API 返回错误",
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(respBytes)),
		)
		return "", fmt.Errorf("LLM API 返回 %d: %s", resp.StatusCode, string(respBytes))
	}

	var result chatResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return "", fmt.Errorf("解析 LLM 响应失败: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("LLM 返回了空结果")
	}

	return result.Choices[0].Message.Content, nil
}
