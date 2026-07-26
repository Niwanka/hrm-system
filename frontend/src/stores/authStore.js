import { defineStore } from 'pinia'
import axios from 'axios'

let isInterceptorSetup = false

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('hrm_user') || 'null'),
    token: localStorage.getItem('hrm_token') || '',
    refreshToken: localStorage.getItem('hrm_refresh_token') || '',
    isLoading: false,
    error: null,
  }),

  getters: {
    isAuthenticated: (state) => !!state.user,
    accessLevel: (state) => state.user?.access_level || 0,
    userRole: (state) => state.user?.role_name || 'Guest',
    fullName: (state) => state.user ? `${state.user.first_name} ${state.user.last_name}` : '',
    
    // RBAC Permission Getters
    isAdmin: (state) => (state.user?.access_level || 0) >= 100,
    isHR: (state) => (state.user?.access_level || 0) >= 80,
    isManager: (state) => (state.user?.access_level || 0) >= 50,
    hasAccessLevel: (state) => (minLevel) => (state.user?.access_level || 0) >= minLevel,
  },

  actions: {
    setupAxiosInterceptors() {
      if (isInterceptorSetup) return
      isInterceptorSetup = true

      // Attach Authorization Header to all outgoing requests
      axios.interceptors.request.use((config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`
        }
        return config
      })

      // Intercept 401 responses and attempt transparent token refresh
      axios.interceptors.response.use(
        (response) => response,
        async (error) => {
          const originalRequest = error.config

          if (
            error.response?.status === 401 &&
            !originalRequest._retry &&
            !originalRequest.url.includes('/api/login') &&
            !originalRequest.url.includes('/api/refresh')
          ) {
            originalRequest._retry = true
            const refreshed = await this.refreshAccessToken()
            if (refreshed) {
              originalRequest.headers.Authorization = `Bearer ${this.token}`
              return axios(originalRequest)
            }
          }
          return Promise.reject(error)
        }
      )
    },

    async login(credentials) {
      this.setupAxiosInterceptors()
      this.isLoading = true
      this.error = null
      try {
        const response = await axios.post('/api/login', credentials, {
          withCredentials: true,
        })

        const { user, token, refresh_token } = response.data
        this.user = user
        this.token = token
        this.refreshToken = refresh_token

        localStorage.setItem('hrm_user', JSON.stringify(user))
        localStorage.setItem('hrm_token', token)
        localStorage.setItem('hrm_refresh_token', refresh_token)
        
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`

        return { success: true }
      } catch (err) {
        this.error = err.response?.data?.error || 'Invalid credentials or server error'
        return { success: false, error: this.error }
      } finally {
        this.isLoading = false
      }
    },

    async refreshAccessToken() {
      try {
        const response = await axios.post(
          '/api/refresh',
          { refresh_token: this.refreshToken },
          { withCredentials: true }
        )

        const { token } = response.data
        this.token = token
        localStorage.setItem('hrm_token', token)
        axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
        return true
      } catch (err) {
        console.warn('Session expired or refresh token revoked. Logging out.', err)
        this.clearAuth()
        return false
      }
    },

    async fetchCurrentUser() {
      this.setupAxiosInterceptors()
      if (!this.token && !this.user) return
      try {
        const response = await axios.get('/api/me', {
          headers: { Authorization: `Bearer ${this.token}` },
          withCredentials: true,
        })
        this.user = response.data.user
        localStorage.setItem('hrm_user', JSON.stringify(this.user))
      } catch (err) {
        console.warn('Session check failed or expired token', err)
        this.clearAuth()
      }
    },

    async logout() {
      try {
        await axios.post('/api/logout', {}, { withCredentials: true })
      } catch (err) {
        console.error('Logout error:', err)
      } finally {
        this.clearAuth()
      }
    },

    clearAuth() {
      this.user = null
      this.token = ''
      this.refreshToken = ''
      localStorage.removeItem('hrm_user')
      localStorage.removeItem('hrm_token')
      localStorage.removeItem('hrm_refresh_token')
      delete axios.defaults.headers.common['Authorization']
    },
  },
})
