<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import ReviewHistory from '../components/ReviewHistory.vue'

const router = useRouter()
const auth = useAuth()

if (!auth.isLoggedIn.value) {
  router.replace('/auth')
}

const historyRef = ref<InstanceType<typeof ReviewHistory>>()

onMounted(() => {
  historyRef.value?.refresh()
})
</script>

<template>
  <div class="history-page" v-if="auth.isLoggedIn.value">
    <button class="back-btn" @click="router.push('/')">← 返回审查</button>
    <ReviewHistory ref="historyRef" />
  </div>
</template>

<style scoped>
.history-page {
  max-width: 800px;
  margin: 0 auto;
}

.back-btn {
  background: none;
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  padding: 8px 16px;
  margin-bottom: 24px;
  transition: border-color 0.15s, color 0.15s;
}

.back-btn:hover {
  border-color: var(--accent);
  color: var(--text-primary);
}
</style>
