// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2025-05-15',
  devtools: { enabled: true },
  modules: [
    '@nuxt/eslint',
    '@nuxt/fonts',
    '@nuxt/image',
    '@nuxt/test-utils',
    '@nuxt/ui',
    '@pinia/nuxt',
    'pinia-plugin-persistedstate',
    '@nuxtjs/color-mode',
  ],
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.API_BASE_URL,
    },
  },
  css: ['@/assets/css/main.css'],
  vite: {
    plugins: [tailwindcss()],
  },
  colorMode: {
    preference: 'light', // светлая тема по умолчанию
    fallback: 'light', // если нет сохранённого значения в localStorage
    classSuffix: '',
  },
  ui: {
    theme: {
      colors: [
        'primary',
        'secondary',
        'background',
        'background-alt',
        'text-primary',
        'text-secondary',
        'text-on-primary',
        'success',
        'error',
        'warning',
        'border',
      ],
    },
  },
})
