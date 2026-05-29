<script setup lang="ts">
import { useRouter } from 'vue-router'
import AuthForm from '../components/AuthForm.vue'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const auth = useAuth()

function onLogin(username: string, token: string) {
  auth.login(username, token)
  router.push('/')
}
</script>

<template>
  <div class="auth-page">
    <button class="back-btn" @click="router.push('/')">← 返回</button>
    <AuthForm @login="onLogin" />
  </div>
</template>

<style scoped>
.auth-page {
  max-width: 500px;
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
