import { defineStore } from 'pinia'
import { ref } from 'vue'
import { accountAPI } from '@/api'
import type { LoginResponse, Profile } from '@/types'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const refreshToken = ref(localStorage.getItem('refresh_token') || '')
  const accountId = ref(parseInt(localStorage.getItem('account_id') || '0'))
  const usernameStore = ref(localStorage.getItem('username') || '')
  const profile = ref<Profile | null>(null)

  const isLoggedIn = () => !!token.value && accountId.value > 0

  const login = async (username: string, password: string) => {
    const res: LoginResponse = await accountAPI.login(username, password)
    token.value = res.token
    refreshToken.value = res.refresh_token
    accountId.value = res.account_id
    usernameStore.value = res.username
    localStorage.setItem('token', res.token)
    localStorage.setItem('refresh_token', res.refresh_token)
    localStorage.setItem('account_id', res.account_id.toString())
    localStorage.setItem('username', res.username)
  }

  const logout = async () => {
    try {
      await accountAPI.logout()
    } catch {
    } finally {
      token.value = ''
      refreshToken.value = ''
      accountId.value = 0
      usernameStore.value = ''
      profile.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('account_id')
      localStorage.removeItem('username')
    }
  }

  const refresh = async () => {
    const res = await accountAPI.refresh(refreshToken.value)
    token.value = res.data.token
    refreshToken.value = res.data.refresh_token
    localStorage.setItem('token', res.data.token)
    localStorage.setItem('refresh_token', res.data.refresh_token)
  }

  const setProfile = (p: Profile) => {
    profile.value = p
  }

  return {
    token,
    refreshToken,
    accountId,
    username: usernameStore,
    profile,
    isLoggedIn,
    login,
    logout,
    refresh,
    setProfile
  }
})