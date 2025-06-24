<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { DropdownMenuItem, FormSubmitEvent, FormError } from '@nuxt/ui'
import { useAuthStore } from '@/stores/useAuthStore'

const auth = useAuthStore()
const loading = computed(() => auth.loading)
const msg = computed(() => auth.error)

const state = reactive({
  login: '',
  password: '',
})

const validate = (state: any): FormError[] => {
  const errors: FormError[] = []
  if (!state.login) errors.push({ name: 'login', message: 'Обязательное поле' })
  if (!state.password)
    errors.push({ name: 'password', message: 'Обязательное поле' })
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  auth.error = ''
  const errors = validate(state)
  if (errors.length > 0) return

  await auth.login(state.login, state.password)
}
</script>

<template>
  <UDropdownMenu
    :content="{ align: 'end', sideOffset: 20 }"
    :ui="{ content: 'w-80' }"
  >
    <UButton
      label="Войти"
      color="text-primary"
      variant="ghost"
      trailing-icon="solar:arrow-right-linear"
      size="xl"
    />

    <template #content-top>
      <UForm
        @submit="onSubmit"
        :validate="validate"
        :state="state"
        class="flex flex-col gap-2 p-4 bg-background/80 backdrop-blur-sm"
      >
        <div class="text-base font-semibold mb-1 flex justify-center">
          Авторизация
        </div>

        <UInput
          name="login"
          v-model="state.login"
          placeholder="Логин"
          variant="outline"
          size="lg"
          icon="tabler:user"
          class="mb-2"
        />

        <UInput
          name="password"
          v-model="state.password"
          type="password"
          placeholder="Пароль"
          variant="outline"
          size="lg"
          icon="tabler:lock-password"
          class="mb-2"
        />

        <div class="flex items-center justify-between mb-2">
          <UButton
            :loading="loading"
            type="submit"
            label="Войти"
            color="text-primary"
            size="md"
            class="w-1/2 flex justify-center cursor-pointer"
          />
          <button
            type="button"
            class="text-xs text-text hover:underline opacity-75 cursor-pointer"
            @click="$emit('forgot-password')"
            tabindex="-1"
          >
            Забыли пароль?
          </button>
        </div>
        <NuxtLink to="/registration" class="block">
          <UButton
            label="Регистрация"
            color="text-primary"
            variant="ghost"
            size="sm"
            block
            class="cursor-pointer"
          />
        </NuxtLink>
        <div v-if="msg" class="text-center text-red-500 mt-2">{{ msg }}</div>
      </UForm>
    </template>
  </UDropdownMenu>
</template>
