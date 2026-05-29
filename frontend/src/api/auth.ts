const TOKEN_KEY = 'ai_pr_review_token'

export interface User {
  id: number
  username: string
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}

function saveToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

async function authRequest(path: string, body: Record<string, string>): Promise<AuthResponse> {
  const resp = await fetch(`/api/v1/auth/${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  const json = await resp.json()
  if (json.code !== 0) {
    throw new Error(json.message || '操作失败')
  }
  return json.data as AuthResponse
}

export async function register(username: string, password: string): Promise<AuthResponse> {
  const data = await authRequest('register', { username, password })
  return data
}

export async function login(username: string, password: string): Promise<AuthResponse> {
  const data = await authRequest('login', { username, password })
  if (data.token) {
    saveToken(data.token)
  }
  return data
}

export function logout() {
  clearToken()
}
