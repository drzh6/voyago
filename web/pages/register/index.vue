<script setup lang="ts">
import { reactive, ref } from 'vue'
import type { FormError, FormSubmitEvent } from '@nuxt/ui'
import { useAuthStore } from '@/stores/useAuthStore'
const authStore = useAuthStore()

const msg = ref('')
const msgType = ref<'success' | 'error' | ''>('')

const state = reactive({
  name: '',
  surname: '',
  password: '',
  confirmPassword: '',
  login: '',
  email: '',
  avatar: '',
})

const validate = (state: any): FormError[] => {
  const errors: FormError[] = []
  if (!state.name) errors.push({ name: 'name', message: 'Обязательное поле' })
  if (!state.surname)
    errors.push({ name: 'surname', message: 'Обязательное поле' })
  if (!state.login) errors.push({ name: 'login', message: 'Обязательное поле' })
  if (!state.email) errors.push({ name: 'email', message: 'Обязательное поле' })
  if (!state.password)
    errors.push({ name: 'password', message: 'Обязательное поле' })
  // if (!state.avatar)
  //   errors.push({ name: 'avatar', message: 'Обязательное поле' })
  if (!state.confirmPassword)
    errors.push({ name: 'confirmPassword', message: 'Обязательное поле' })
  return errors
}

async function onSubmit(event: FormSubmitEvent<typeof state>) {
  // Сброс сообщений
  msg.value = ''
  msgType.value = ''

  // Проверка совпадения паролей
  if (
    state.password &&
    state.confirmPassword &&
    state.password !== state.confirmPassword
  ) {
    msg.value = 'Пароли не совпадают'
    msgType.value = 'error'
    return
  }

  // Собираем только те поля, которые ожидает сервер
  const payload = {
    name: state.name,
    surname: state.surname,
    login: state.login,
    email: state.email,
    password: state.password,
    avatar: state.avatar || undefined,
  }

  await authStore.register(payload)

  if (authStore.error) {
    msg.value = authStore.error
    msgType.value = 'error'
  } else {
    msg.value = 'Регистрация прошла успешно'
    msgType.value = 'success'
  }

  // очистка формы
  state.name = ''
  state.surname = ''
  state.login = ''
  state.email = ''
  state.password = ''
  state.confirmPassword = ''
  state.avatar = ''
}
</script>

<template>
  <div>
    <h2 class="mb-4 text-xl font-semibold">Регистрация</h2>
    <UForm
      @submit="onSubmit"
      :validate="validate"
      :state="state"
      class="max-w-[720px] flex flex-col space-y-8"
    >
      <UInput
        v-model="state.name"
        placeholder="Имя"
        size="lg"
        icon="tabler:user-square-rounded"
        required
      />

      <UInput
        v-model="state.surname"
        placeholder="Фамилия"
        size="lg"
        icon="tabler:user-square-rounded"
        required
      />

      <UInput
        v-model="state.login"
        placeholder="Логин"
        size="lg"
        icon="tabler:user"
        required
      />

      <UInput
        v-model="state.email"
        type="email"
        placeholder="Email"
        size="lg"
        icon="tabler:mail"
        required
      />

      <UInput
        v-model="state.password"
        type="password"
        placeholder="Пароль"
        size="lg"
        icon="tabler:lock-password"
        required
      />
      <UInput
        v-model="state.confirmPassword"
        type="password"
        placeholder="Подтвердите пароль"
        size="lg"
        icon="tabler:lock-password"
        required
      />

      <!--      <UInput-->
      <!--        v-model="state.avatar"-->
      <!--        placeholder="URL аватара"-->
      <!--        size="lg"-->
      <!--        icon="tabler:photo"-->
      <!--      />-->

      <UButton type="submit" block> Зарегистрироваться </UButton>

      <div
        v-if="msg"
        :class="[
          msgType === 'success' ? 'text-success' : 'text-error',
          'text-center mt-2',
        ]"
      >
        {{ msg }}
      </div>
    </UForm>
  </div>
</template>
