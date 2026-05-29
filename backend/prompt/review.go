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

你必须严格按以下 JSON 结构输出，不要输出任何其他内容：
{
  "summary": "PR 变更总结，简洁描述做了什么改动（中文，200字以内）",
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
}`

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
