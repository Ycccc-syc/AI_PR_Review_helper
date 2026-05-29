<script setup lang="ts">
import { computed } from 'vue'
import type { ReviewResponse } from '../api/review'

const props = defineProps<{
  data: ReviewResponse
}>()

const tabs = [
  { key: 'summary', label: '变更总结' },
  { key: 'risks', label: `风险代码 (${props.data.risk_items.length})` },
  { key: 'suggestions', label: `改进建议 (${props.data.suggestions.length})` },
]

const activeTab = defineModel('tab', { default: 'risks' })

const severityLabel = { high: '高', medium: '中', low: '低' }
const severityColor = { high: '#dc2626', medium: '#d97706', low: '#6b7280' }

const hasRisks = computed(() => props.data.risk_items.length > 0)
const hasSuggestions = computed(() => props.data.suggestions.length > 0)
</script>

<template>
  <div v-if="data" class="result-panel">
    <!-- 统计卡片 -->
    <div class="stats">
      <div class="stat"><span class="num">{{ data.stats.files_changed }}</span> 文件变更</div>
      <div class="stat add">+{{ data.stats.additions }}</div>
      <div class="stat del">-{{ data.stats.deletions }}</div>
      <div class="stat" :style="{ color: severityColor.high }">🔴 {{ data.stats.risk_high }}</div>
      <div class="stat" :style="{ color: severityColor.medium }">🟡 {{ data.stats.risk_medium }}</div>
      <div class="stat" :style="{ color: severityColor.low }">⚪ {{ data.stats.risk_low }}</div>
    </div>

    <!-- Tab 切换 -->
    <div class="tabs">
      <button
        v-for="t in tabs"
        :key="t.key"
        :class="{ active: activeTab === t.key }"
        @click="activeTab = t.key"
      >
        {{ t.label }}
      </button>
    </div>

    <!-- 总结 -->
    <div v-if="activeTab === 'summary'" class="tab-content summary">
      <p>{{ data.summary }}</p>
    </div>

    <!-- 风险列表 -->
    <div v-if="activeTab === 'risks'" class="tab-content">
      <div v-if="!hasRisks" class="empty">未发现风险代码</div>
      <div v-for="(r, i) in data.risk_items" :key="i" class="risk-card">
        <div class="risk-header">
          <span class="severity" :style="{ background: severityColor[r.severity] }">
            {{ severityLabel[r.severity] }}
          </span>
          <span class="risk-title">{{ r.title }}</span>
        </div>
        <div class="risk-location">{{ r.file }}:{{ r.line }}</div>
        <p class="risk-detail">{{ r.detail }}</p>
        <p class="risk-fix">💡 {{ r.suggestion }}</p>
      </div>
    </div>

    <!-- 建议列表 -->
    <div v-if="activeTab === 'suggestions'" class="tab-content">
      <div v-if="!hasSuggestions" class="empty">暂无改进建议</div>
      <div v-for="(s, i) in data.suggestions" :key="i" class="suggestion-card">
        <span class="priority" :style="{ background: severityColor[s.priority] }">
          {{ severityLabel[s.priority] }}
        </span>
        <p>{{ s.content }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.result-panel {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px 48px;
}

.stats {
  display: flex;
  gap: 16px;
  padding: 16px 24px;
  background: #f9fafb;
  border-radius: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.stat {
  font-size: 13px;
  color: #374151;
}

.stat .num {
  font-weight: 700;
  font-size: 18px;
}

.stat.add { color: #16a34a; font-weight: 700; }
.stat.del { color: #dc2626; font-weight: 700; }

.tabs {
  display: flex;
  gap: 0;
  border-bottom: 2px solid #e5e7eb;
  margin-bottom: 24px;
}

.tabs button {
  padding: 10px 20px;
  border: none;
  background: none;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
}

.tabs button.active {
  color: #6366f1;
  border-bottom-color: #6366f1;
}

.tab-content {
  min-height: 120px;
}

.summary p {
  line-height: 1.8;
  color: #374151;
}

.empty {
  text-align: center;
  color: #9ca3af;
  padding: 32px;
}

.risk-card {
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 10px;
  padding: 16px 20px;
  margin-bottom: 12px;
}

.risk-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.severity,
.priority {
  color: #fff;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  flex-shrink: 0;
}

.risk-title {
  font-weight: 600;
  font-size: 15px;
}

.risk-location {
  font-size: 12px;
  color: #6366f1;
  font-family: monospace;
  margin-bottom: 8px;
}

.risk-detail {
  color: #4b5563;
  font-size: 14px;
  line-height: 1.7;
  margin-bottom: 8px;
}

.risk-fix {
  background: #f0fdf4;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 13px;
  color: #166534;
}

.suggestion-card {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  padding: 14px 0;
  border-bottom: 1px solid #f0f0f0;
}

.suggestion-card p {
  color: #374151;
  line-height: 1.7;
}
</style>
