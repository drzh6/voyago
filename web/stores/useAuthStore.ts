// stores/useAuthStore.ts
import { defineStore } from 'pinia'
import { useRouter } from 'vue-router'

interface User {
  id: string
  name: string
  surname: string
  login: string
  email: string
  avatar: string
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const router = useRouter()

  // Инициализация из sessionStorage
  if (typeof window !== 'undefined') {
    const storedToken = sessionStorage.getItem('token')
    const storedUser = sessionStorage.getItem('user')
    if (storedToken) {
      token.value = JSON.parse(storedToken)
    }
    if (storedUser && storedUser !== 'undefined') {
      user.value = JSON.parse(storedUser)
    }
  }

  // Логин
  async function login(login: string, password: string) {
    try {
      loading.value = true
      error.value = null

      const { $api } = useNuxtApp()
      const { data } = await $api.post('/login/', { login, password })

      token.value = data.token
      user.value = data.user

      if (typeof window !== 'undefined') {
        sessionStorage.setItem('token', JSON.stringify(token.value))
        if (user.value !== null) {
          sessionStorage.setItem('user', JSON.stringify(user.value))
        } else {
          sessionStorage.removeItem('user')
        }
      }
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Ошибка авторизации'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Регистрация
  async function register(payload: {
    name: string
    surname: string
    login: string
    email: string
    password: string
    avatar?: string
  }) {
    try {
      loading.value = true
      error.value = null

      const { $api } = useNuxtApp()
      const { data } = await $api.post('/register/', payload)

      token.value = data.token
      user.value = data.user

      if (typeof window !== 'undefined') {
        sessionStorage.setItem('token', JSON.stringify(token.value))
        if (user.value !== null) {
          sessionStorage.setItem('user', JSON.stringify(user.value))
        } else {
          sessionStorage.removeItem('user')
        }
      }
    } catch (err: any) {
      error.value = err.response?.data?.detail || 'Ошибка регистрации'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Выход
  function logout() {
    token.value = null
    user.value = null

    if (typeof window !== 'undefined') {
      sessionStorage.removeItem('token')
      sessionStorage.removeItem('user')
    }

    router.push('/')
  }

  return {
    user,
    token,
    loading,
    error,
    login,
    register,
    logout,
  }
})
