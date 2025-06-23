// stores/useTripsStore.ts
import { defineStore } from 'pinia'

export interface Trip {
  id: string
  name: string
  description: string
  owner_id: string
  startDate: Date
  endDate: Date
  status: string
  is_public: boolean
  invite_code: string
  cover_image: string
  created_at: Date
  updated_at: Date
}

export const useTripsStore = defineStore('trips', () => {
  // state
  const trips = ref<Trip[]>([])
  const activeTripId = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // getters
  const activeTrip = computed(() => {
    return trips.value.find((t) => t.id === activeTripId.value) || null
  })

  // actions
  async function fetchTrip(id: string) {
    const existingTrip = trips.value.find((t) => t.id === id)
    if (existingTrip) {
      activeTripId.value = id
      return
    }

    try {
      loading.value = true
      error.value = null

      // Запрос к API
      const { $api } = useNuxtApp()
      const { data } = await $api.get(`/api/trips/${id}`)

      trips.value.push(data)
      activeTripId.value = id
    } catch (err) {
      error.value = 'Не удалось загрузить маршрут'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  function setActiveTrip(id: string) {
    activeTripId.value = id
  }

  // вернём всё, что нужно использовать в компонентах
  return {
    trips,
    activeTripId,
    activeTrip,
    loading,
    error,
    fetchTrip,
    setActiveTrip,
  }
})
