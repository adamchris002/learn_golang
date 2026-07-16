<route lang="json">
{
  "meta": {
    "requiresAuth": true
  }
}
</route>
<script setup lang="ts">
import router from "@/router"
import { useAuthStore } from "@/stores/auth"
import {NButton} from "naive-ui"

const authStore = useAuthStore()

function handleLogout() {
  authStore.logout()
  router.push('/Login')
}
</script>

<template>
  <div>
    <h1 class="text-3xl font-bold text-blue-600">Home</h1>
    <n-button @click="handleLogout">Log out</n-button>
  </div>

</template>