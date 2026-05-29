package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// streamRequest 流式请求体
type streamRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// Choice DeepSeek API 返回的选项
type choice struct {
	Message Message `json:"message"`
	Delta   struct {
		Content string `json:"content"`
	} `json:"delta"`
}

// ChatResponse DeepSeek API 返回体
type chatResponse struct {
	Choices []choice `json:"choices"`
}

// StreamChunk 流式响应的单个 chunk
type StreamChunk struct {
	Content string
	Done    bool
	Error   error
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

// CallStream DeepSeek Chat API 流式调用，返回 chunk channel
func CallStream(apiKey, model string, messages []Message) <-chan StreamChunk {
	ch := make(chan StreamChunk, 100)

	go func() {
		defer close(ch)

		reqBody := streamRequest{
			Model:       model,
			Messages:    messages,
			Temperature: 0.3,
			Stream:      true,
		}

		bodyBytes, err := json.Marshal(reqBody)
		if err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("序列化请求失败: %w", err)}
			return
		}

		req, err := http.NewRequest("POST", deepseekBaseURL+"/v1/chat/completions", bytes.NewReader(bodyBytes))
		if err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("创建请求失败: %w", err)}
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		client := &http.Client{Timeout: 180 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("LLM API 请求失败: %w", err)}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
			ch <- StreamChunk{Error: fmt.Errorf("LLM API 返回 %d: %s", resp.StatusCode, string(body))}
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				ch <- StreamChunk{Done: true}
				return
			}

			var streamResp chatResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				continue
			}
			if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
				ch <- StreamChunk{Content: streamResp.Choices[0].Delta.Content}
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("读取流式响应失败: %w", err)}
			return
		}
		ch <- StreamChunk{Done: true}
	}()

	return ch
}
