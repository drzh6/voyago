// plugins/axios.ts
import axios from 'axios'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()

  const api = axios.create({
    baseURL: config.public.apiBaseUrl,
    timeout: 10000,
  })

  return {
    provide: {
      api,
    },
  }
})
