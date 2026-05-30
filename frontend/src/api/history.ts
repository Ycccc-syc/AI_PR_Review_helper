import { getToken } from './auth'
import type { ReviewResponse } from './review'

export interface ReviewLog {
  id: number
  user_id: number
  pr_url: string
  pr_title: string
  summary: string
  risk_count: number
  created_at: string
}

export interface ReviewDetail {
  id: number
  user_id: number
  pr_url: string
  pr_title: string
  summary: string
  risk_count: number
  created_at: string
  result: ReviewResponse | null
}

export async function fetchHistory(limit = 20, offset = 0): Promise<ReviewLog[]> {
  const token = getToken()
  const resp = await fetch(`/api/v1/reviews?limit=${limit}&offset=${offset}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  })

  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || '获取历史记录失败')
  }
  return json.data as ReviewLog[]
}

export async function fetchReviewDetail(id: number): Promise<ReviewDetail> {
  const token = getToken()
  const resp = await fetch(`/api/v1/reviews/${id}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
  })

  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || '获取记录详情失败')
  }
  return json.data as ReviewDetail
}

export function formatTime(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins} 分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours} 小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days} 天前`
  return d.toLocaleDateString('zh-CN')
}
