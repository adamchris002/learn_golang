<script setup lang="ts">
import { useAuthStore } from './stores/auth';
import { NAlert } from 'naive-ui';

const authStore = useAuthStore()
</script>

<template>
  <div v-if="authStore.message !== null" class="absolute z-9999 top-3 right-3">
    <n-alert :title="authStore.message.messageTitle" :type="authStore.message.status === 400 ? 'error' : 'success'"
      closable @close="authStore.setMessage(null)">
      {{ authStore.message?.message }}
    </n-alert>
  </div>
  <router-view />
</template>