<script setup lang="ts">
import { ref } from 'vue'
import { login as doLogin, register as doRegister } from '../api/auth'

const emit = defineEmits<{
  login: [username: string, token: string]
}>()

const isLogin = ref(true)
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!username.value.trim() || !password.value.trim()) {
    error.value = '请填写用户名和密码'
    return
  }

  loading.value = true
  try {
    if (isLogin.value) {
      const data = await doLogin(username.value, password.value)
      emit('login', data.user.username, data.token)
    } else {
      await doRegister(username.value, password.value)
      isLogin.value = true
      error.value = ''
      username.value = ''
      password.value = ''
    }
  } catch (e: any) {
    error.value = e.message || '操作失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-card">
    <h2>{{ isLogin ? '登录' : '注册' }}</h2>

    <form @submit.prevent="submit" class="auth-form">
      <div class="field">
        <label>用户名</label>
        <input
          v-model="username"
          type="text"
          placeholder="输入用户名"
          :disabled="loading"
          autocomplete="username"
        />
      </div>

      <div class="field">
        <label>密码</label>
        <input
          v-model="password"
          type="password"
          placeholder="输入密码"
          :disabled="loading"
          autocomplete="current-password"
        />
      </div>

      <div v-if="error" class="auth-error">{{ error }}</div>
      <div v-if="!isLogin && !error" class="auth-hint">注册成功后需登录获取令牌</div>

      <button type="submit" :disabled="loading" class="auth-btn">
        <span v-if="loading" class="spinner"></span>
        {{ isLogin ? '登录' : '注册' }}
      </button>
    </form>

    <p class="switch-mode">
      {{ isLogin ? '还没有账号？' : '已有账号？' }}
      <button class="link-btn" @click="isLogin = !isLogin; error = ''">
        {{ isLogin ? '去注册' : '去登录' }}
      </button>
    </p>
  </div>
</template>

<style scoped>
.auth-card {
  max-width: 400px;
  margin: 0 auto 48px;
  padding: 32px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 12px;
}

h2 {
  margin: 0 0 24px;
  font-size: 20px;
  text-align: center;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field label {
  font-size: 13px;
  color: var(--text-secondary);
}

.field input {
  padding: 10px 14px;
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.field input:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.field input::placeholder {
  color: var(--text-muted);
}

.auth-error {
  padding: 10px 14px;
  background: rgba(248, 81, 73, 0.1);
  border: 1px solid rgba(248, 81, 73, 0.2);
  border-radius: 8px;
  color: var(--red);
  font-size: 13px;
}

.auth-hint {
  padding: 10px 14px;
  background: rgba(99, 102, 241, 0.08);
  border: 1px solid rgba(99, 102, 241, 0.15);
  border-radius: 8px;
  color: var(--accent-hover);
  font-size: 13px;
}

.auth-btn {
  padding: 12px;
  margin-top: 4px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: opacity 0.2s;
}

.auth-btn:hover:not(:disabled) {
  opacity: 0.9;
}

.auth-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.switch-mode {
  margin: 20px 0 0;
  text-align: center;
  font-size: 13px;
  color: var(--text-muted);
}

.link-btn {
  background: none;
  border: none;
  color: var(--accent);
  font-size: 13px;
  cursor: pointer;
  padding: 0;
}

.link-btn:hover {
  color: var(--accent-hover);
  text-decoration: underline;
}

.spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
