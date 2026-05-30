<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useReview } from '../composables/useReview'
import { askFollowUp } from '../api/review'
import type { RiskItem } from '../api/review'

const router = useRouter()
const route = useRoute()
const review = useReview()

const riskIndex = Number(route.query.i || 0)
const risk = ref<RiskItem | null>(null)

onMounted(() => {
  if (!review.result.value) {
    router.replace('/')
    return
  }
  risk.value = review.result.value.risk_items[riskIndex] || null
})

function goBack() {
  router.push('/')
}

// ---- chat state ----
interface Turn {
  question: string
  answer: string
  loading: boolean
}
const turns = ref<Turn[]>([])
const input = ref('')
const loading = ref(false)

async function send() {
  const q = input.value.trim()
  if (!q || loading.value || !risk.value) return

  input.value = ''
  const turn: Turn = { question: q, answer: '', loading: true }
  turns.value.push(turn)
  loading.value = true

  const ctx = JSON.stringify({
    severity: risk.value.severity,
    file: risk.value.file,
    line: risk.value.line,
    title: risk.value.title,
    detail: risk.value.detail,
    suggestion: risk.value.suggestion,
  })

  try {
    const answer = await askFollowUp(
      review.prUrl.value,
      q,
      ctx,
      (chunk) => {
        turn.answer += chunk
        nextTick(() => {
          const el = document.getElementById('chat-scroll')
          if (el) el.scrollTop = el.scrollHeight
        })
      },
    )
    turn.answer = answer
  } catch (e: any) {
    turn.answer = e.message || '追问失败'
  } finally {
    turn.loading = false
    loading.value = false
  }
}

function formatAnswer(text: string) {
  // 简单 markdown 风格：将 ##、**、- 等转为 HTML
  return text
    .replace(/^## (.+)$/gm, '<h4>$1</h4>')
    .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
    .replace(/^[•-] (.+)$/gm, '<li>$1</li>')
    .replace(/```(\w*)\n([\s\S]*?)```/g, '<pre><code>$2</code></pre>')
    .replace(/\n\n/g, '</p><p>')
}
</script>

<template>
  <div class="ask-page" v-if="risk">
    <button class="back-btn" @click="goBack">← 返回审查</button>

    <!-- 风险项摘要 -->
    <div class="risk-context" :class="'severity-' + risk.severity">
      <div class="rc-header">
        <span class="severity-badge" :class="risk.severity">
          {{ { high: 'HIGH', medium: 'MED', low: 'LOW' }[risk.severity] }}
        </span>
        <span class="rc-title">{{ risk.title }}</span>
      </div>
      <div class="rc-file">{{ risk.file }}<span class="line-num">:{{ risk.line }}</span></div>
    </div>

    <!-- 对话区 -->
    <div id="chat-scroll" class="chat-list">
      <div v-for="(t, i) in turns" :key="i" class="chat-turn">
        <div class="chat-q">
          <span class="chat-label">Q</span>
          <span>{{ t.question }}</span>
        </div>
        <div class="chat-a">
          <span class="chat-label">A</span>
          <div class="chat-body" v-html="formatAnswer(t.answer)"></div>
          <div v-if="t.loading" class="chat-typing">
            <span class="typing-dot"></span>AI 正在回答...
          </div>
        </div>
      </div>
    </div>

    <!-- 输入区 -->
    <div class="chat-input-row">
      <input
        v-model="input"
        type="text"
        placeholder="针对这个风险项继续提问..."
        :disabled="loading"
        @keydown.enter="send"
      />
      <button class="send-btn" :disabled="loading || !input.trim()" @click="send">
        发送
      </button>
    </div>
  </div>
</template>

<style scoped>
.ask-page {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 24px 40px;
}

.back-btn {
  background: none;
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  padding: 8px 16px;
  margin-bottom: 20px;
  transition: border-color 0.15s, color 0.15s;
}

.back-btn:hover {
  border-color: var(--accent);
  color: var(--text-primary);
}

/* 风险项摘要 */
.risk-context {
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 14px 18px;
  margin-bottom: 24px;
  border-left: 3px solid var(--border);
}

.risk-context.severity-high { border-left-color: var(--red); }
.risk-context.severity-medium { border-left-color: var(--orange); }
.risk-context.severity-low { border-left-color: var(--text-muted); }

.rc-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.rc-title {
  font-weight: 600;
  font-size: 15px;
  color: var(--text-primary);
}

.rc-file {
  font-family: var(--font-mono);
  font-size: 12px;
  color: var(--accent-hover);
}

.rc-file .line-num {
  color: var(--orange);
}

.severity-badge {
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  padding: 2px 7px;
  border-radius: 4px;
  letter-spacing: 0.5px;
}

.severity-badge.high { background: rgba(248, 81, 73, 0.15); color: var(--red); }
.severity-badge.medium { background: rgba(210, 153, 34, 0.15); color: var(--orange); }
.severity-badge.low { background: rgba(110, 118, 129, 0.15); color: var(--text-muted); }

/* 对话列表 */
.chat-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 480px;
  overflow-y: auto;
  margin-bottom: 20px;
  padding-right: 4px;
}

.chat-q, .chat-a {
  display: flex;
  gap: 10px;
}

.chat-a {
  margin-top: 4px;
}

.chat-label {
  flex-shrink: 0;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
  font-family: var(--font-mono);
}

.chat-q .chat-label {
  background: rgba(99, 102, 241, 0.12);
  color: var(--accent);
}

.chat-a .chat-label {
  background: rgba(63, 185, 80, 0.12);
  color: var(--green);
}

.chat-q span:last-child {
  padding-top: 3px;
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.6;
}

.chat-body {
  padding-top: 2px;
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.75;
  flex: 1;
  min-width: 0;
}

/* markdown-style formatting inside answers */
.chat-body :deep(h4) {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 8px 0 4px;
}

.chat-body :deep(li) {
  margin: 2px 0 2px 8px;
}

.chat-body :deep(pre) {
  background: #f5f5f5;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 10px 14px;
  margin: 8px 0;
  font-family: var(--font-mono);
  font-size: 12px;
  color: #333;
  overflow-x: auto;
  line-height: 1.6;
}

.chat-body :deep(p) {
  margin: 4px 0;
}

.chat-body :deep(strong) {
  color: var(--text-primary);
  font-weight: 600;
}

.chat-typing {
  margin-top: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--accent);
}

.typing-dot {
  width: 6px;
  height: 6px;
  background: var(--accent);
  border-radius: 50%;
  animation: pulse 1.2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

/* 输入区 */
.chat-input-row {
  display: flex;
  gap: 10px;
  position: sticky;
  bottom: 0;
  background: var(--bg-primary);
  padding-top: 12px;
}

.chat-input-row input {
  flex: 1;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  transition: border-color 0.15s;
}

.chat-input-row input:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.send-btn {
  padding: 12px 24px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  border: none;
  border-radius: 24px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s, transform 0.15s;
  white-space: nowrap;
}

.send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.35);
}

.send-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
</style>
