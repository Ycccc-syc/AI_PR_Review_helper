<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchHistory, fetchReviewDetail, formatTime, type ReviewLog } from '../api/history'
import { useReview } from '../composables/useReview'

const router = useRouter()
const reviewStore = useReview()
const logs = ref<ReviewLog[]>([])
const loading = ref(false)
const error = ref('')
const loadingId = ref<number | null>(null)

async function load() {
  loading.value = true
  error.value = ''
  try {
    logs.value = await fetchHistory()
  } catch (e: any) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

async function goDetail(log: ReviewLog) {
  if (loadingId.value) return
  loadingId.value = log.id
  try {
    const detail = await fetchReviewDetail(log.id)
    if (detail.result) {
      reviewStore.setReview(detail.result, log.pr_url)
      router.push('/')
    } else {
      error.value = '该记录未存储完整结果（可能是旧版数据）'
    }
  } catch (e: any) {
    error.value = e.message || '加载详情失败'
  } finally {
    loadingId.value = null
  }
}

onMounted(load)

defineExpose({ refresh: load })
</script>

<template>
  <div class="history-panel">
    <div class="history-header">
      <h3>历史记录</h3>
      <button class="refresh-btn" @click="load" :disabled="loading">刷新</button>
    </div>

    <div v-if="error" class="history-error">{{ error }}</div>

    <div v-if="loading && logs.length === 0" class="history-empty">加载中...</div>

    <div v-else-if="logs.length === 0" class="history-empty">暂无审查记录</div>

    <div v-else class="history-list">
      <div
        v-for="log in logs" :key="log.id"
        class="history-item"
        :class="{ loading: loadingId === log.id }"
        @click="goDetail(log)"
      >
        <div class="item-top">
          <span class="pr-title">{{ log.pr_title || log.pr_url }}</span>
          <span class="time">{{ formatTime(log.created_at) }}</span>
        </div>
        <p class="item-summary">{{ log.summary }}</p>
        <div class="item-meta">
          <span class="risk-badge">{{ log.risk_count }} 个风险项</span>
          <span class="expand-hint">
            {{ loadingId === log.id ? '加载中...' : '点击查看完整报告 →' }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-panel {
  max-width: 680px;
  margin: 0 auto;
  padding: 0 24px 64px;
}

.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.history-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.refresh-btn {
  padding: 6px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.refresh-btn:hover:not(:disabled) {
  border-color: var(--accent);
}

.history-error {
  padding: 10px 14px;
  background: rgba(248, 81, 73, 0.1);
  border-radius: 8px;
  color: var(--red);
  font-size: 13px;
  margin-bottom: 16px;
}

.history-empty {
  text-align: center;
  color: var(--text-muted);
  padding: 48px 0;
  font-size: 14px;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.history-item {
  padding: 16px 20px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  transition: border-color 0.15s, opacity 0.15s;
}

.history-item:hover {
  border-color: var(--accent);
}

.history-item.loading {
  opacity: 0.6;
  pointer-events: none;
}

.pr-title {
  color: var(--accent);
  font-size: 14px;
  font-weight: 600;
  font-family: var(--font-mono);
}

.item-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.time {
  font-size: 12px;
  color: var(--text-muted);
  flex-shrink: 0;
  margin-left: 12px;
}

.item-summary {
  margin: 0 0 8px;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.item-meta {
  display: flex;
  gap: 12px;
  align-items: center;
}

.risk-badge {
  font-size: 11px;
  padding: 3px 10px;
  background: rgba(210, 153, 34, 0.1);
  color: var(--orange);
  border-radius: 4px;
  font-family: var(--font-mono);
}

.expand-hint {
  font-size: 11px;
  color: var(--text-muted);
}
</style>
