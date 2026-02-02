import { defineStore } from 'pinia'

interface User {
  id: number
  email: string
  name: string
  avatar?: string
}


export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)

  const setToken = (t: string) => {
    token.value = t
  }

  const setUser = (u: User) => {
    user.value = u
  }

  const clearAuth = () => {
    user.value = null
    token.value = null
  }

  const logout = () => {
    clearAuth()
    // Optional: redirect to login page
    // navigateTo('/login')
  }

  return { user, token, setToken, setUser, clearAuth, logout }
} , {
  persist: true
})