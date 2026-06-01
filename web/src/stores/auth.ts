import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as apiLogin, getMe } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<{ user_id: number; username: string; role: string; totp_enabled?: boolean } | null>(
    JSON.parse(localStorage.getItem('user') || 'null'),
  )
  const isAuthenticated = ref(!!token.value)

  async function login(username: string, password: string, code?: string) {
    const res = await apiLogin({ username, password, code })
    if (res.status === 'require_2fa') {
      return { require_2fa: true }
    }
    token.value = res.token || ''
    user.value = {
      user_id: res.user_id || 0,
      username: res.username || '',
      role: res.role || '',
      totp_enabled: res.status !== 'require_2fa' // will update via checkAuth
    }
    localStorage.setItem('token', token.value)
    localStorage.setItem('user', JSON.stringify(user.value))
    isAuthenticated.value = true
    return { require_2fa: false }
  }

  function logout() {
    token.value = ''
    user.value = null
    isAuthenticated.value = false
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function checkAuth() {
    if (!token.value) return
    try {
      const res = await getMe()
      // getMe returns totp_enabled
      user.value = {
        user_id: res.user_id || 0,
        username: res.username || '',
        role: res.role || '',
        totp_enabled: (res as any).totp_enabled
      }
      localStorage.setItem('user', JSON.stringify(user.value))
      isAuthenticated.value = true
    } catch {
      logout()
    }
  }

  return { token, user, isAuthenticated, login, logout, checkAuth }
})
