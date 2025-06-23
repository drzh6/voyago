<script setup lang="ts">
import { useMenuStore } from '@/stores/useMenuStore'
const menu = useMenuStore()
</script>

<template>
  <div
    v-show="menu.isOpen"
    class="fixed left-0 right-0 top-16 bg-background-alt z-40 shadow-md"
  >
    <ul class="flex flex-col space-y-4 py-4 pl-4 text-lg font-bold">
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="mingcute:home-3-fill" class="size-5" />
          <NuxtLink to="/" @click="menu.close()">Главная</NuxtLink>
        </div>
      </li>
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="mage:dashboard-bar-fill" class="size-5" />
          <NuxtLink to="/dashboard" @click="menu.close()">Дашборд</NuxtLink>
        </div>
      </li>
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="lets-icons:map-fill" class="size-5" />
          <NuxtLink to="/routes" @click="menu.close()">Маршрут</NuxtLink>
        </div>
      </li>
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="ic:baseline-wallet" class="size-5" />
          <NuxtLink to="/expenses" @click="menu.close()">Расходы</NuxtLink>
        </div>
      </li>
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="flowbite:users-solid" class="size-5" />
          <NuxtLink to="/participants" @click="menu.close()"
            >Участники</NuxtLink
          >
        </div>
      </li>
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="mingcute:comment-fill" class="size-5" />
          <NuxtLink to="/comments" @click="menu.close()">Комментарии</NuxtLink>
        </div>
      </li>
      <li>
        <div class="flex items-center gap-2">
          <UIcon name="iconamoon:settings-fill" class="size-5" />
          <NuxtLink to="/settings" @click="menu.close()">Настройки</NuxtLink>
        </div>
      </li>
    </ul>
  </div>
</template>
