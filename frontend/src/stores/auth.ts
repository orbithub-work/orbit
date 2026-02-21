// 认证 Store - 用户认证状态管理
// 
// 提供用户登录、登出、认证状态初始化等功能

import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as { id: string; name: string; email: string } | null,
    token: localStorage.getItem('token') || null as string | null,
    isAuthenticated: false,
  }),

  getters: {
    isLoggedIn: (state) => !!state.token,
    currentUser: (state) => state.user,
  },

  actions: {
    // 用户登录
    login(credentials: { email: string; password: string }) {
      // 实际实现中应调用 API 进行认证
      // 这里模拟成功的登录
      this.token = 'fake-jwt-token-' + Date.now()
      this.user = {
        id: 'user-' + Date.now(),
        name: credentials.email.split('@')[0],
        email: credentials.email,
      }
      this.isAuthenticated = true
      
      // 将 token 存储到 localStorage
      localStorage.setItem('token', this.token)
    },

    // 用户登出
    logout() {
      this.token = null
      this.user = null
      this.isAuthenticated = false
      
      // 从 localStorage 移除 token
      localStorage.removeItem('token')
    },

    // 从存储的 token 初始化认证状态
    initializeAuth() {
      const storedToken = localStorage.getItem('token')
      if (storedToken) {
        this.token = storedToken
        this.isAuthenticated = true
        // 实际应用中应与后端验证 token
      }
    }
  },

  // 将 store 持久化到 localStorage
  persist: {
    key: 'auth-store',
    storage: localStorage,
  }
})