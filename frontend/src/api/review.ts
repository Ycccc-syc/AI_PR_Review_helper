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

export async function analyzePR(prUrl: string): Promise<ReviewResponse> {
  const controller = new AbortController()
  const timer = setTimeout(() => controller.abort(), 120_000)

  let resp: Response
  try {
    resp = await fetch('/api/v1/review', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
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
