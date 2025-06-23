<script setup lang="ts">
import { ref } from 'vue'
import Header from '@/components/layout/Header.vue'
import MobileMenu from '@/components/layout/MobileMenu.vue'
import Footer from '@/components/layout/Footer.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <Header />
    <MobileMenu />

    <main class="flex-1 pt-20">
      <slot />
    </main>

    <Footer />
  </div>
</template>
