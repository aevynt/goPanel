import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login as apiLogin, getMe } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<{ user_id: number; username: string; role: string } | null>(
    JSON.parse(localStorage.getItem('user') || 'null'),
  )
  const isAuthenticated = ref(!!token.value)

  async function login(username: string, password: string) {
    const res = await apiLogin({ username, password })
    token.value = res.token
    user.value = { user_id: res.user_id, username: res.username, role: res.role }
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(user.value))
    isAuthenticated.value = true
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
      user.value = { user_id: res.user_id, username: res.username, role: res.role }
      localStorage.setItem('user', JSON.stringify(user.value))
      isAuthenticated.value = true
    } catch {
      logout()
    }
  }

  return { token, user, isAuthenticated, login, logout, checkAuth }
})
