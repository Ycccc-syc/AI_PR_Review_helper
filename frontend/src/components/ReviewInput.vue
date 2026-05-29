<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  submit: [url: string]
}>()

const prUrl = ref('')
const loading = ref(false)

function onSubmit() {
  const url = prUrl.value.trim()
  if (!url) return
  emit('submit', url)
}

function setLoading(v: boolean) {
  loading.value = v
}

defineExpose({ setLoading })
</script>

<template>
  <div class="hero">
    <div class="logo">
      <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
        <rect width="40" height="40" rx="10" fill="url(#g)" />
        <path d="M12 20l5 5 11-11" stroke="#fff" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
        <defs><linearGradient id="g" x1="0" y1="0" x2="40" y2="40"><stop stop-color="#6366f1"/><stop offset="1" stop-color="#a855f7"/></linearGradient></defs>
      </svg>
    </div>
    <h1>AI PR Review</h1>
    <p class="desc">输入 GitHub PR 链接，AI 自动审查代码变更，发现潜在风险</p>
    <form @submit.prevent="onSubmit" class="input-row">
      <div class="glass-wrap">
        <svg class="icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/></svg>
        <input
          v-model="prUrl"
          type="text"
          placeholder="https://github.com/owner/repo/pull/123"
          :disabled="loading"
        />
      </div>
      <button type="submit" :disabled="loading || !prUrl.trim()" class="submit-btn">
        <span v-if="loading" class="spinner"></span>
        {{ loading ? '分析中...' : '开始分析' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.hero {
  text-align: center;
  padding: 64px 24px 40px;
  animation: fadeIn 0.6s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

.logo {
  margin-bottom: 16px;
}

h1 {
  font-family: var(--font-mono);
  font-size: 32px;
  font-weight: 700;
  margin: 0 0 8px;
  letter-spacing: -0.5px;
}

.desc {
  color: var(--text-secondary);
  font-size: 15px;
  margin: 0 0 32px;
}

.input-row {
  display: flex;
  gap: 10px;
  max-width: 680px;
  margin: 0 auto;
}

.glass-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border);
  border-radius: 12px;
  backdrop-filter: blur(8px);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.glass-wrap:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.glass-wrap .icon {
  color: var(--text-muted);
  flex-shrink: 0;
}

.glass-wrap input {
  flex: 1;
  padding: 12px 0;
  background: none;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-size: 14px;
  font-family: var(--font-mono);
}

.glass-wrap input::placeholder {
  color: var(--text-muted);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.submit-btn {
  padding: 12px 28px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  border: none;
  border-radius: 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
  transition: opacity 0.2s, transform 0.15s, box-shadow 0.2s;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.35);
}

.submit-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
