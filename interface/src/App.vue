<script setup lang="ts">
import { watch } from 'vue'
import { useAuthStore } from './stores/auth';
import { NAlert } from 'naive-ui';

const authStore = useAuthStore()

let alertTimeout: ReturnType<typeof setTimeout> | null = null
watch(
  () => authStore.message,
  (message) => {
    if (alertTimeout) {
      clearTimeout(alertTimeout)
    }

    if (message) {
      alertTimeout = setTimeout(() => {
        authStore.setMessage(null)
      }, 5000)
    }
  },
  { immediate: true }
)
</script>

<template>
  <div v-if="authStore.message !== null" class="absolute z-9999 top-3 right-3">
    <n-alert :title="authStore.message.messageTitle" :type="authStore.message.status === 200 ? 'success' : 'error'"
      closable @close="authStore.setMessage(null)">
      {{ authStore.message?.message }}
    </n-alert>
  </div>
  <router-view />
</template>