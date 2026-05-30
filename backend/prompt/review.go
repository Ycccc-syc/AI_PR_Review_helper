package prompt

import (
	"fmt"

	"qiniudemo/llm"
)

const systemPrompt = `你是一位资深代码审查专家。你的任务是对 GitHub PR 的代码变更进行专业 Review。

请遵循以下原则：
1. 只报告有实质影响的问题，忽略命名偏好、注释格式等主观风格问题
2. 每条风险必须附带具体的 file 路径和 line number
3. 不确定的问题标注为 low 严重度
4. 改进建议要具体、可执行，不要泛泛而谈
5. 必须给出 health_score 综合评分，从安全/性能/可维护三个维度评估

你必须严格按以下 JSON 结构输出，不要输出任何其他内容：
{
  "summary": "PR 变更总结，简洁描述做了什么改动（中文，200字以内）",
  "health_score": {
    "total": 85,
    "security": 80,
    "performance": 90,
    "maintainability": 85
  },
  "risk_items": [
    {
      "severity": "high/medium/low",
      "file": "文件路径",
      "line": 行号,
      "title": "风险标题",
      "detail": "风险详细说明",
      "suggestion": "修复建议"
    }
  ],
  "suggestions": [
    {
      "priority": "high/medium/low",
      "content": "改进建议内容"
    }
  ]
}

health_score 评分标准（每项 0-100）：
- total: 综合评分，按安全30% + 性能30% + 可维护40% 加权计算
- security: 是否存在 SQL 注入、XSS、密钥泄露、权限绕过等安全漏洞
- performance: 是否存在 N+1 查询、内存泄漏、不必要的大对象分配等性能问题
- maintainability: 函数复杂度是否过高、是否有硬编码魔法数字、错误处理是否完善
- 无风险时该项给满分 100，少量低风险扣 5-15 分，有高危扣 20-40 分`

const askSystemPrompt = `你是一位资深代码审查专家。用户正在针对之前审查中发现的风险项进行追问。
请基于提供的风险项上下文，给出专业、具体、可执行的回答。

回答格式要求：
- 使用 Markdown 格式组织你的回答
- 使用 ## 小标题分隔不同的回答段落（如"## 问题分析"、"## 修复建议"、"## 代码示例"）
- 重要术语使用 **加粗** 标注
	- 代码示例请用三个反引号代码块包裹，标注语言类型
- 使用 - 或 * 列出要点
- 回答简洁直接，聚焦用户问题，不要泛泛而谈
- 如果用户问题不够具体，主动引导追问方向`

// BuildAskMessages 构建追问对话消息，contextJSON 为风险项的 JSON 字符串
func BuildAskMessages(prTitle, prDescription, diff, question, contextJSON string) []llm.Message {
	desc := prDescription
	if desc == "" {
		desc = "（无描述）"
	}
	if len(desc) > 300 {
		desc = desc[:300] + "..."
	}
	if len(diff) > 6000 {
		diff = diff[:6000] + "\n... (diff truncated)"
	}

	userPrompt := fmt.Sprintf(`## 原始 PR
标题: %s
描述: %s

## 相关风险项
%s

## 追问
%s

请针对以上追问给出详细解答。`, prTitle, desc, contextJSON, question)

	return []llm.Message{
		{Role: "system", Content: askSystemPrompt},
		{Role: "user", Content: userPrompt},
	}
}

func BuildMessages(prTitle, prDescription, diff string) []llm.Message {
	userPrompt := fmt.Sprintf(`## PR 标题
%s

## PR 描述
%s

## 代码变更 (unified diff)
%s`, prTitle, prDescription, diff)

	return []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}
}
