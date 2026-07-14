<script setup lang="ts">
import loginBackground from "@/assets/images/background.png"
import facebookLogo from '@/assets/logos/facebook_logo.png'
import githubLogo from '@/assets/logos/github_logo.png'
import googleLogo from '@/assets/logos/google_logo.png'
import { NInput, NButton } from 'naive-ui'

</script>
<template>
  <div class="h-screen w-screen flex justify-center items-center bg-cover bg-center bg-no-repeat"
    :style="{ backgroundImage: `url(${loginBackground})` }">
    <div class="px-120 py-40 backdrop-blur-sm rounded-2xl inset-shadow-[0_0_80px_rgba(0,0,0,0.25)]">
      <div class="px-20 py-10 bg-white/10 inset-shadow-sm rounded-3xl">
        <h1 class="font-jakarta mb-10 text-white text-2xl font-bold">Welcome! This is practice</h1>
        <h1 class="font-jakarta mb-6 text-white text-2xl font-bold">Login</h1>
        <div class="mb-4">
          <p class="font-jakrta mb-2 font-light text-[13px] text-white">Username: </p>
          <n-input placeholder="username" />
        </div>
        <div class="mb-[11px]">
          <p class="font-jakrta mb-2 font-light text-[13px] text-white">Password: </p>
          <n-input placeholder="Password" />
        </div>
        <p class="font-jakarta font-light text-[13px] text-white">Forgot Password?</p>
        <div class="my-6">
          <n-button block color="#003465">Sign in</n-button>
        </div>
        <div class="mb-6 flex items-center">
          <p class="font-jakarta font-light text-[13px] text-white">or continue with</p>
        </div>
        <div class="flex justify-between">
          <div class="py-2 px-7 rounded-xl bg-white"><img class="w-6 h-6" :src="googleLogo" /></div>
          <div class="py-2 px-7 rounded-xl bg-white"><img class="w-6 h-6" :src="githubLogo" /></div>
          <div class="py-2 px-7 rounded-xl bg-white"><img class="w-6 h-6" :src="facebookLogo" /></div>
        </div>
        <p class="mt-6 font-jakarta font-light text-[13px] text-white">
          <span>Don’t have an account yet?</span>
          <span class="font-semibold"> Register for free</span>
        </p>
      </div>
    </div>
  </div>
</template>
<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');
</style>