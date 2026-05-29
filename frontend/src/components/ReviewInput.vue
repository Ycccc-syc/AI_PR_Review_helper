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
  <div class="review-input">
    <h1>AI PR Review</h1>
    <p class="subtitle">输入 GitHub PR 链接，AI 自动分析代码变更</p>
    <form @submit.prevent="onSubmit" class="input-row">
      <input
        v-model="prUrl"
        type="text"
        placeholder="https://github.com/owner/repo/pull/123"
        :disabled="loading"
      />
      <button type="submit" :disabled="loading || !prUrl.trim()">
        <span v-if="loading" class="spinner"></span>
        {{ loading ? '分析中...' : '开始分析' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.review-input {
  text-align: center;
  padding: 48px 24px;
}

h1 {
  font-size: 28px;
  margin-bottom: 8px;
}

.subtitle {
  color: #666;
  margin-bottom: 24px;
}

.input-row {
  display: flex;
  gap: 12px;
  max-width: 600px;
  margin: 0 auto;
}

input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 14px;
  outline: none;
}

input:focus {
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

button {
  padding: 10px 24px;
  background: #6366f1;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid #fff;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
