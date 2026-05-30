<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { ReviewResponse } from '../api/review'

const props = defineProps<{
  data: ReviewResponse
}>()

const router = useRouter()
const activeTab = defineModel('tab', { default: 'risks' })
const expanded = ref<Record<number, boolean>>({})

function toggle(i: number) {
  expanded.value[i] = !expanded.value[i]
}

function goAsk(index: number) {
  router.push({ path: '/ask', query: { i: String(index) } })
}

const tabs = [
  { key: 'summary', label: '变更总结' },
  { key: 'risks', label: `风险代码 (${props.data.risk_items.length})` },
  { key: 'suggestions', label: `改进建议 (${props.data.suggestions.length})` },
]

const hasRisks = computed(() => props.data.risk_items.length > 0)
const hasSuggestions = computed(() => props.data.suggestions.length > 0)
</script>

<template>
  <div v-if="data" class="result-panel">
    <!-- 统计栏 -->
    <div class="stats">
      <div class="stat">
        <span class="num">{{ data.stats.files_changed }}</span>
        <span class="label">文件变更</span>
      </div>
      <div class="stat add">+{{ data.stats.additions }}</div>
      <div class="stat del">-{{ data.stats.deletions }}</div>
      <div class="stat risk">
        <span class="dot" style="background:#f85149"></span>
        <span class="num">{{ data.stats.risk_high }}</span>
        <span class="label">高危</span>
      </div>
      <div class="stat risk">
        <span class="dot" style="background:#d29922"></span>
        <span class="num">{{ data.stats.risk_medium }}</span>
        <span class="label">中危</span>
      </div>
      <div class="stat risk">
        <span class="dot" style="background:#6e7681"></span>
        <span class="num">{{ data.stats.risk_low }}</span>
        <span class="label">低危</span>
      </div>
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
      <div v-for="(r, i) in data.risk_items" :key="i" class="risk-card" :class="'severity-' + r.severity">
        <div class="risk-header">
          <span class="severity-badge" :class="r.severity">
            {{ { high: 'HIGH', medium: 'MED', low: 'LOW' }[r.severity] }}
          </span>
          <span class="risk-title">{{ r.title }}</span>
        </div>
        <div class="risk-file">{{ r.file }}<span class="line-num">:{{ r.line }}</span></div>
        <p class="risk-detail">{{ r.detail }}</p>
        <div class="risk-fix">
          <div class="fix-label">修复建议</div>
          {{ r.suggestion }}
        </div>
        <div class="risk-actions">
          <button class="ask-btn" @click="goAsk(i)">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            追问
          </button>
        </div>
      </div>
    </div>

    <!-- 改进建议 -->
    <div v-if="activeTab === 'suggestions'" class="tab-content">
      <div v-if="!hasSuggestions" class="empty">暂无改进建议</div>
      <div v-for="(s, i) in data.suggestions" :key="i" class="suggestion-card" @click="toggle(i)">
        <div class="sug-header">
          <span class="priority-badge" :class="s.priority">
            {{ { high: 'HIGH', medium: 'MED', low: 'LOW' }[s.priority] }}
          </span>
          <span class="sug-content">{{ s.content }}</span>
          <svg
            class="chevron"
            :class="{ open: expanded[i] }"
            width="16" height="16" viewBox="0 0 24 24"
            fill="none" stroke="currentColor" stroke-width="2"
          >
            <path d="M6 9l6 6 6-6"/>
          </svg>
        </div>
        <div v-if="expanded[i]" class="sug-detail">
          <div class="code-block">
            <code>// 示例代码片段将在后续版本中根据具体上下文生成</code>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.result-panel {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px 64px;
  animation: slideUp 0.4s ease;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 统计栏 */
.stats {
  display: flex;
  gap: 4px;
  padding: 16px 24px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.stat {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  font-size: 13px;
  color: var(--text-secondary);
}

.stat .num {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 16px;
}

.stat .label {
  font-size: 12px;
  color: var(--text-muted);
}

.stat.add { color: var(--green); }
.stat.del { color: var(--red); }

.stat.risk .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

/* Tabs */
.tabs {
  display: flex;
  gap: 0;
  border-bottom: 1px solid var(--border);
  margin-bottom: 24px;
}

.tabs button {
  padding: 10px 20px;
  border: none;
  background: none;
  font-size: 14px;
  color: var(--text-muted);
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -1px;
  transition: color 0.15s;
}

.tabs button:hover {
  color: var(--text-primary);
}

.tabs button.active {
  color: var(--accent);
  border-bottom-color: var(--accent);
}

.tab-content {
  min-height: 120px;
}

/* Summary */
.summary p {
  line-height: 1.8;
  color: var(--text-secondary);
}

/* Empty */
.empty {
  text-align: center;
  color: var(--text-muted);
  padding: 40px;
  font-size: 14px;
}

/* Risk card */
.risk-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 16px 20px;
  margin-bottom: 12px;
  border-left: 3px solid var(--border);
  transition: border-color 0.2s;
}

.risk-card.severity-high { border-left-color: var(--red); }
.risk-card.severity-medium { border-left-color: var(--orange); }
.risk-card.severity-low { border-left-color: var(--text-muted); }

.risk-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.severity-badge {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 4px;
  flex-shrink: 0;
  letter-spacing: 0.5px;
}

.severity-badge.high { background: rgba(248, 81, 73, 0.15); color: var(--red); }
.severity-badge.medium { background: rgba(210, 153, 34, 0.15); color: var(--orange); }
.severity-badge.low { background: rgba(110, 118, 129, 0.15); color: var(--text-muted); }

.risk-title {
  font-weight: 600;
  font-size: 15px;
  color: var(--text-primary);
}

.risk-file {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--accent-hover);
  margin-bottom: 8px;
}

.risk-file .line-num {
  color: var(--orange);
}

.risk-detail {
  color: var(--text-secondary);
  font-size: 14px;
  line-height: 1.7;
  margin-bottom: 10px;
}

.risk-fix {
  background: rgba(63, 185, 80, 0.06);
  border: 1px solid rgba(63, 185, 80, 0.12);
  padding: 10px 14px;
  border-radius: 6px;
  font-size: 13px;
  color: #7ee787;
  line-height: 1.6;
}

.fix-label {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: var(--green);
  margin-bottom: 4px;
}

/* Suggestion card */
.suggestion-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 14px 18px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.suggestion-card:hover {
  border-color: var(--accent);
}

.sug-header {
  display: flex;
  align-items: flex-start;
  gap: 10px;
}

.priority-badge {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  padding: 3px 8px;
  border-radius: 4px;
  flex-shrink: 0;
  letter-spacing: 0.5px;
}

.priority-badge.high { background: rgba(248, 81, 73, 0.15); color: var(--red); }
.priority-badge.medium { background: rgba(210, 153, 34, 0.15); color: var(--orange); }
.priority-badge.low { background: rgba(110, 118, 129, 0.15); color: var(--text-muted); }

.sug-content {
  flex: 1;
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.6;
}

.chevron {
  color: var(--text-muted);
  flex-shrink: 0;
  transition: transform 0.2s;
  margin-top: 2px;
}

.chevron.open {
  transform: rotate(180deg);
}

.sug-detail {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border);
}

.code-block {
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 14px 16px;
  font-family: var(--font-mono);
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.6;
  overflow-x: auto;
}

/* 追问按钮 */
.risk-actions {
  margin-top: 10px;
}

.ask-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  background: rgba(99, 102, 241, 0.06);
  border: 1px solid rgba(99, 102, 241, 0.15);
  border-radius: 8px;
  font-size: 13px;
  color: var(--accent);
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s, transform 0.15s;
}

.ask-btn:hover {
  background: rgba(99, 102, 241, 0.12);
  border-color: rgba(99, 102, 241, 0.3);
  transform: translateY(-1px);
}
</style>
