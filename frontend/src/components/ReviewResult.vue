<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { ReviewResponse } from '../api/review'

const props = defineProps<{
  data: ReviewResponse
}>()

const router = useRouter()
const activeTab = defineModel('tab', { default: 'risks' })

const animated = ref(false)

const C = 2 * Math.PI * 52 // 圆环周长

function ringStyle(score: number) {
  const offset = animated.value ? C * (1 - score / 100) : C
  return {
    strokeDasharray: C,
    strokeDashoffset: offset,
    transition: animated.value ? 'stroke-dashoffset 1s cubic-bezier(0.34, 1.56, 0.64, 1)' : 'none',
  }
}

onMounted(() => {
  requestAnimationFrame(() => { animated.value = true })
})


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
    <!-- 健康分：霓虹细环 -->
    <div class="score-panel">
      <!-- 综合评分 -->
      <div class="score-ring-wrap">
        <div class="ring-icon check">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5"><path d="M5 13l4 4L19 7"/></svg>
        </div>
        <div class="score-ring">
          <svg viewBox="0 0 120 120">
            <circle class="ring-track" cx="60" cy="60" r="52" />
            <circle
              class="ring-fill purple"
              cx="60" cy="60" r="52"
              :style="ringStyle(data.health_score.total)"
            />
          </svg>
          <div class="ring-center">
            <span class="ring-num">{{ data.health_score.total }}</span>
            <span class="ring-label">综合评分</span>
          </div>
        </div>
      </div>
      <!-- 安全 -->
      <div class="score-ring-wrap">
        <div class="ring-icon shield">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
        </div>
        <div class="score-ring">
          <svg viewBox="0 0 120 120">
            <circle class="ring-track" cx="60" cy="60" r="52" />
            <circle
              class="ring-fill red"
              cx="60" cy="60" r="52"
              :style="ringStyle(data.health_score.security)"
            />
          </svg>
          <div class="ring-center">
            <span class="ring-num">{{ data.health_score.security }}</span>
            <span class="ring-label">安全</span>
          </div>
        </div>
      </div>
      <!-- 性能 -->
      <div class="score-ring-wrap">
        <div class="ring-icon bolt">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/></svg>
        </div>
        <div class="score-ring">
          <svg viewBox="0 0 120 120">
            <circle class="ring-track" cx="60" cy="60" r="52" />
            <circle
              class="ring-fill gold"
              cx="60" cy="60" r="52"
              :style="ringStyle(data.health_score.performance)"
            />
          </svg>
          <div class="ring-center">
            <span class="ring-num">{{ data.health_score.performance }}</span>
            <span class="ring-label">性能</span>
          </div>
        </div>
      </div>
      <!-- 可维护 -->
      <div class="score-ring-wrap">
        <div class="ring-icon code">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M16 18l6-6-6-6M8 6l-6 6 6 6"/></svg>
        </div>
        <div class="score-ring">
          <svg viewBox="0 0 120 120">
            <circle class="ring-track" cx="60" cy="60" r="52" />
            <circle
              class="ring-fill green"
              cx="60" cy="60" r="52"
              :style="ringStyle(data.health_score.maintainability)"
            />
          </svg>
          <div class="ring-center">
            <span class="ring-num">{{ data.health_score.maintainability }}</span>
            <span class="ring-label">可维护</span>
          </div>
        </div>
      </div>
    </div>

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
      <div v-for="(s, i) in data.suggestions" :key="i" class="suggestion-card">
        <div class="sug-header">
          <span class="priority-badge" :class="s.priority">
            {{ { high: 'HIGH', medium: 'MED', low: 'LOW' }[s.priority] }}
          </span>
          <span class="sug-content">{{ s.content }}</span>
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

/* 健康分面板：毛玻璃卡片 */
.score-panel {
  display: flex;
  gap: 20px;
  padding: 28px 20px;
  background: rgba(22, 27, 34, 0.65);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid rgba(48, 54, 61, 0.6);
  border-radius: 16px;
  margin-bottom: 24px;
  flex-wrap: wrap;
  justify-content: center;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.25);
}

/* 单个圆环容器 */
.score-ring-wrap {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 140px;
}

/* 图标 */
.ring-icon {
  position: absolute;
  top: 18px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2;
  width: 18px;
  height: 18px;
}
.ring-icon svg { width: 100%; height: 100%; }
.ring-icon.check  { color: #a78bfa; }
.ring-icon.shield { color: #f87171; }
.ring-icon.bolt   { color: #fbbf24; }
.ring-icon.code   { color: #4ade80; }

/* SVG 圆环 */
.score-ring {
  position: relative;
  width: 130px;
  height: 130px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.score-ring svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

/* 底层轨道 */
.ring-track {
  fill: none;
  stroke: rgba(48, 54, 61, 0.5);
  stroke-width: 3;
}

/* 进度条 */
.ring-fill {
  fill: none;
  stroke-width: 3;
  stroke-linecap: round;
}
.ring-fill.purple { stroke: #a78bfa; filter: drop-shadow(0 0 8px rgba(167, 139, 250, 0.6)); }
.ring-fill.red    { stroke: #f87171; filter: drop-shadow(0 0 8px rgba(248, 113, 113, 0.6)); }
.ring-fill.gold   { stroke: #fbbf24; filter: drop-shadow(0 0 8px rgba(251, 191, 36, 0.6)); }
.ring-fill.green  { stroke: #4ade80; filter: drop-shadow(0 0 8px rgba(74, 222, 128, 0.6)); }

/* 环中心数字 */
.ring-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}
.ring-num {
  font-family: var(--font-mono);
  font-size: 30px;
  font-weight: 800;
  color: #e6edf3;
  line-height: 1;
  text-shadow: 0 0 12px rgba(255, 255, 255, 0.08);
}
.ring-label {
  font-size: 12px;
  color: #8b949e;
  letter-spacing: 0.5px;
}

/* 呼吸动画 */
@keyframes breathe {
  0%, 100% { opacity: 0.7; }
  50%      { opacity: 1; }
}
.score-ring-wrap {
  animation: breathe 3s ease-in-out infinite;
}
.score-ring-wrap:nth-child(1) { animation-delay: 0s; }
.score-ring-wrap:nth-child(2) { animation-delay: 0.75s; }
.score-ring-wrap:nth-child(3) { animation-delay: 1.5s; }
.score-ring-wrap:nth-child(4) { animation-delay: 2.25s; }

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
