export interface ReviewRequest {
  pr_url: string
}

export interface RiskItem {
  severity: 'high' | 'medium' | 'low'
  file: string
  line: number
  title: string
  detail: string
  suggestion: string
}

export interface Suggestion {
  priority: 'high' | 'medium' | 'low'
  content: string
}

export interface ReviewStats {
  files_changed: number
  additions: number
  deletions: number
  risk_high: number
  risk_medium: number
  risk_low: number
}

export interface ReviewResponse {
  summary: string
  risk_items: RiskItem[]
  suggestions: Suggestion[]
  stats: ReviewStats
}

// SSE 流式审查：onChunk 实时回调文本片段，返回最终解析结果
export async function analyzePRStream(
  prUrl: string,
  onChunk: (text: string) => void,
  token?: string | null,
): Promise<ReviewResponse> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) headers['Authorization'] = `Bearer ${token}`

  const resp = await fetch('/api/v1/review/stream', {
    method: 'POST',
    headers,
    body: JSON.stringify({ pr_url: prUrl }),
  })

  if (!resp.ok) {
    const text = await resp.text().catch(() => '')
    throw new Error(text || `请求失败 (${resp.status})`)
  }

  // 后端 API 约定始终返回 HTTP 200，需要检查 Content-Type 区分 SSE 和 JSON 错误
  const ct = resp.headers.get('content-type') || ''
  if (!ct.includes('text/event-stream')) {
    const json = await resp.json().catch(() => null)
    if (json && json.code !== 0) throw new Error(json.message || '请求失败')
    throw new Error('非预期的响应格式')
  }

  const reader = resp.body?.getReader()
  if (!reader) throw new Error('浏览器不支持流式响应')

  const decoder = new TextDecoder()
  let buffer = ''
  let resolved = false

  return new Promise((resolve, reject) => {
    function pump(): void {
      reader!.read().then(({ done, value }) => {
        if (done) {
          if (!resolved) reject(new Error('连接意外关闭'))
          return
        }

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() || ''

        for (const line of lines) {
          if (!line.startsWith('data: ')) continue
          try {
            const msg = JSON.parse(line.slice(6))
            if (msg.type === 'chunk') {
              onChunk(msg.content)
            } else if (msg.type === 'result') {
              resolved = true
              resolve(msg.data as ReviewResponse)
            } else if (msg.type === 'error') {
              resolved = true
              reject(new Error(msg.message))
            }
          } catch { /* skip malformed */ }
        }

        pump()
      }).catch(reject)
    }

    pump()
  })
}

export async function analyzePR(prUrl: string, token?: string | null): Promise<ReviewResponse> {
  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(), 120_000)

  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  let resp: Response
  try {
    resp = await fetch('/api/v1/review', {
      method: 'POST',
      headers,
      body: JSON.stringify({ pr_url: prUrl }),
      signal: controller.signal,
    })
  } finally {
    clearTimeout(timer)
  }

  // 网络错误或超时时 resp 未赋值，直接 throw
  if (!resp.ok) {
    const text = await resp.text().catch(() => '')
    throw new Error(text || `请求失败 (${resp.status})`)
  }

  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || '分析失败')
  }
  return json.data as ReviewResponse
}
