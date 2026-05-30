<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import ReviewInput from '../components/ReviewInput.vue'
import ReviewResult from '../components/ReviewResult.vue'
import { analyzePRStream, type ReviewResponse } from '../api/review'
import { useReview } from '../composables/useReview'

const reviewStore = useReview()
const reviewRef = ref<InstanceType<typeof ReviewInput>>()
const result = ref<ReviewResponse | null>(null)
const error = ref('')
const activeTab = ref('risks')

const streamingText = ref('')
const isStreaming = ref(false)
const prUrl = ref('')
let errorTimer: ReturnType<typeof setTimeout> | null = null

// 从历史记录跳转过来时，直接展示已保存的结果
onMounted(() => {
  if (reviewStore.result.value) {
    result.value = reviewStore.result.value
    prUrl.value = reviewStore.prUrl.value
  }
})

function showError(msg: string) {
  error.value = msg
  if (errorTimer) clearTimeout(errorTimer)
  errorTimer = setTimeout(() => { error.value = '' }, 3000)
}

async function onSubmit(url: string) {
  prUrl.value = url
  error.value = ''
  result.value = null
  streamingText.value = ''
  isStreaming.value = true
  reviewRef.value?.setLoading(true)

  const token = localStorage.getItem('ai_pr_review_token')
  try {
    result.value = await analyzePRStream(url, (chunk) => {
      streamingText.value += chunk
      nextTick(() => {
        const el = document.getElementById('stream-box')
        if (el) el.scrollTop = el.scrollHeight
      })
    }, token)
    reviewStore.setReview(result.value, url)
  } catch (e: any) {
    showError(e.message || '分析失败，请稍后重试')
  } finally {
    isStreaming.value = false
    reviewRef.value?.setLoading(false)
  }
}
</script>

<template>
  <div class="home">
    <ReviewInput ref="reviewRef" @submit="onSubmit" />

    <Transition name="slide-fade">
      <div v-if="error" class="error">{{ error }}</div>
    </Transition>

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
  </div>
</template>

<style scoped>
.error {
  max-width: 600px;
  margin: 0 auto 24px;
  padding: 10px 24px;
  background: #fef2f2;
  color: #dc2626;
  border-radius: 24px;
  font-size: 13px;
  text-align: center;
}

.slide-fade-enter-active {
  animation: fadeIn 0.4s ease;
}

.slide-fade-leave-active {
  animation: fadeIn 0.3s ease reverse;
}

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
