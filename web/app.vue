<script setup lang="ts">
const colorMode = useColorMode()
onMounted(() => {
  colorMode.preference = 'light'
})
</script>

<template>
  <UApp>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>
