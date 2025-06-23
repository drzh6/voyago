<script setup lang="ts">
import { useTripsStore } from '@/stores/useTripsStore'
import { useRoute } from 'vue-router'

const route = useRoute()
const tripsStore = useTripsStore()

onMounted(() => {
  const tripId = route.params.tripId as string
  tripsStore.fetchTrip(tripId)
})
</script>

<template>
  <div class="flex flex-col">
    <div v-if="tripsStore.loading">Загрузка...</div>
    <div v-else-if="tripsStore.error">{{ tripsStore.error }}</div>
    <div v-else-if="tripsStore.activeTrip">
      <p class="text-2xl font-bold">{{ tripsStore.activeTrip.name }}</p>
      <p>{{ tripsStore.activeTrip.description }}</p>
    </div>
    <div v-else>Нет данных о маршруте</div>
  </div>
</template>

<style scoped></style>
