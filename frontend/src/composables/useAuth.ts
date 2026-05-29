import { ref, computed } from 'vue'

const token = ref(localStorage.getItem('ai_pr_review_token') || '')
const username = ref(localStorage.getItem('ai_pr_review_user') || '')

export function useAuth() {
  const isLoggedIn = computed(() => !!token.value)

  function login(name: string, t: string) {
    token.value = t
    username.value = name
    localStorage.setItem('ai_pr_review_token', t)
    localStorage.setItem('ai_pr_review_user', name)
  }

  function logout() {
    token.value = ''
    username.value = ''
    localStorage.removeItem('ai_pr_review_token')
    localStorage.removeItem('ai_pr_review_user')
  }

  return { username, isLoggedIn, login, logout }
}
