<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useMessage } from 'naive-ui'

const router = useRouter()
const auth = useAuthStore()
const message = useMessage()

const username = ref('')
const password = ref('')
const loading = ref(false)

async function handleLogin() {
  if (!username.value || !password.value) {
    message.error('Please enter username and password')
    return
  }
  loading.value = true
  try {
    await auth.login(username.value, password.value)
    router.push('/dashboard')
  } catch (err: any) {
    message.error(err.response?.data?.error || 'Login failed')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1 class="title">goPanel</h1>
        <p class="subtitle">Server Management Panel</p>
      </div>
      <n-form @submit.prevent="handleLogin">
        <n-form-item label="Username">
          <n-input v-model:value="username" placeholder="admin" />
        </n-form-item>
        <n-form-item label="Password">
          <n-input v-model:value="password" type="password" placeholder="••••••" @keyup.enter="handleLogin" />
        </n-form-item>
        <n-button type="primary" block :loading="loading" attr-type="submit" size="large">
          Sign In
        </n-button>
      </n-form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: var(--claude-bg, #f5f4ed);
  padding: 16px;
}

.login-card {
  background: var(--claude-surface, #faf9f5);
  border: 1px solid var(--claude-border-light, #f0eee6);
  border-radius: var(--radius-lg, 12px);
  padding: 40px;
  width: 380px;
  max-width: 100%;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.04);
}

@media (max-width: 420px) {
  .login-card {
    padding: 24px;
  }
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.title {
  font-family: var(--font-heading, 'Playfair Display', serif);
  font-size: 1.75rem;
  font-weight: 500;
  color: var(--claude-text-primary, #141413);
  letter-spacing: -0.02em;
  margin: 0 0 4px 0;
}

.subtitle {
  color: var(--claude-text-secondary, #5e5d59);
  margin: 0;
  font-size: 14px;
}
</style>
