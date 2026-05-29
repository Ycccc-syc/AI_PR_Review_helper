<script setup lang="ts">
import { ref, computed, nextTick } from 'vue'
import ReviewInput from './components/ReviewInput.vue'
import ReviewResult from './components/ReviewResult.vue'
import AuthForm from './components/AuthForm.vue'
import ReviewHistory from './components/ReviewHistory.vue'
import { analyzePRStream, type ReviewResponse } from './api/review'
import { getToken, logout as doLogout } from './api/auth'

const reviewRef = ref<InstanceType<typeof ReviewInput>>()
const historyRef = ref<InstanceType<typeof ReviewHistory>>()
const result = ref<ReviewResponse | null>(null)
const error = ref('')
const activeTab = ref('risks')

const username = ref(localStorage.getItem('ai_pr_review_user') || '')
const showAuth = ref(false)
const showHistory = ref(false)

const isLoggedIn = computed(() => !!getToken())

const streamingText = ref('')
const isStreaming = ref(false)

function onLogin(name: string) {
  username.value = name
  localStorage.setItem('ai_pr_review_user', name)
  showAuth.value = false
}

function onLogout() {
  doLogout()
  username.value = ''
  localStorage.removeItem('ai_pr_review_user')
  showHistory.value = false
}

function toggleHistory() {
  showHistory.value = !showHistory.value
  if (showHistory.value) {
    historyRef.value?.refresh()
  }
}

async function onSubmit(url: string) {
  error.value = ''
  result.value = null
  streamingText.value = ''
  isStreaming.value = true
  reviewRef.value?.setLoading(true)

  try {
    result.value = await analyzePRStream(url, (chunk) => {
      streamingText.value += chunk
      nextTick(() => {
        const el = document.getElementById('stream-box')
        if (el) el.scrollTop = el.scrollHeight
      })
    })
  } catch (e: any) {
    error.value = e.message || '分析失败，请稍后重试'
  } finally {
    isStreaming.value = false
    reviewRef.value?.setLoading(false)
  }
}
</script>

<template>
  <div class="app">
    <!-- 导航栏 -->
    <nav class="navbar">
      <span class="brand">AI PR Review</span>
      <div class="nav-actions">
        <template v-if="isLoggedIn">
          <span class="user-tag">{{ username }}</span>
          <button class="nav-btn" @click="toggleHistory">
            {{ showHistory ? '返回审查' : '历史记录' }}
          </button>
          <button class="nav-btn logout" @click="onLogout">退出</button>
        </template>
        <template v-else>
          <button class="nav-btn" @click="showAuth = !showAuth">
            {{ showAuth ? '关闭' : '登录 / 注册' }}
          </button>
        </template>
      </div>
    </nav>

    <!-- 登录/注册面板 -->
    <AuthForm v-if="showAuth" @login="onLogin" />

    <!-- 历史记录面板 -->
    <ReviewHistory v-if="showHistory" ref="historyRef" />

    <!-- 主审查界面 -->
    <template v-if="!showHistory">
      <ReviewInput ref="reviewRef" @submit="onSubmit" />

      <div v-if="error" class="error">{{ error }}</div>

      <!-- 流式分析面板 -->
      <div v-if="streamingText && !result" class="stream-panel">
        <div class="stream-header">
          <span class="stream-dot"></span>
          AI 正在分析...
        </div>
        <div id="stream-box" class="stream-box">
          <pre>{{ streamingText }}</pre>
        </div>
      </div>

      <ReviewResult
        v-if="result"
        :data="result"
        v-model:tab="activeTab"
      />
    </template>
  </div>
</template>

<style scoped>
.app {
  max-width: 900px;
  margin: 0 auto;
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 32px;
}

.brand {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 700;
  color: var(--accent);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-tag {
  font-size: 13px;
  color: var(--text-secondary);
  padding: 4px 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  font-family: var(--font-mono);
}

.nav-btn {
  padding: 6px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.nav-btn:hover {
  border-color: var(--accent);
  color: var(--text-primary);
}

.nav-btn.logout:hover {
  border-color: var(--red);
  color: var(--red);
}

.error {
  max-width: 600px;
  margin: 0 auto 24px;
  padding: 12px 20px;
  background: #fef2f2;
  color: #dc2626;
  border-radius: 8px;
  font-size: 14px;
  text-align: center;
}

/* 流式面板 */
.stream-panel {
  max-width: 800px;
  margin: 0 auto 24px;
  padding: 0 24px;
}

.stream-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--accent);
  margin-bottom: 8px;
}

.stream-dot {
  width: 8px;
  height: 8px;
  background: var(--accent);
  border-radius: 50%;
  animation: pulse 1.2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.stream-box {
  max-height: 320px;
  overflow-y: auto;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px;
}

.stream-box pre {
  margin: 0;
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-word;
  line-height: 1.7;
}
</style>
