<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuth } from './composables/useAuth'

const router = useRouter()
const auth = useAuth()
</script>

<template>
  <div class="app">
    <nav class="navbar">
      <span class="brand" @click="router.push('/')">AI PR Review</span>
      <div class="nav-actions">
        <template v-if="auth.isLoggedIn.value">
          <span class="user-tag">{{ auth.username.value }}</span>
          <button class="nav-btn" @click="router.push('/history')">历史记录</button>
          <button class="nav-btn logout" @click="auth.logout()">退出</button>
        </template>
        <template v-else>
          <div class="nav-auth-group">
            <button class="nav-btn" @click="router.push('/auth')">登录 / 注册</button>
            <span class="auth-hint">登录后可查看历史记录</span>
          </div>
        </template>
      </div>
    </nav>

    <router-view />
  </div>
</template>

<style scoped>
.app {
  max-width: 900px;
  margin: 0 auto;
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 32px;
}

.brand {
  font-family: var(--font-mono);
  font-size: 16px;
  font-weight: 700;
  color: var(--accent);
  cursor: pointer;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-tag {
  font-size: 13px;
  color: var(--text-secondary);
  padding: 4px 10px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  font-family: var(--font-mono);
}

.nav-btn {
  padding: 6px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.nav-btn:hover {
  border-color: var(--accent);
  color: var(--text-primary);
}

.nav-btn.logout:hover {
  border-color: var(--red);
  color: var(--red);
}

.nav-auth-group {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.auth-hint {
  font-size: 11px;
  color: var(--text-muted);
}
</style>
