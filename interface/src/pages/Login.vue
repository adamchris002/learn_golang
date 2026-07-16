<script setup lang="ts">
import { ref, watch } from "vue"

import loginBackground from "@/assets/images/background.png"
import facebookLogo from '@/assets/logos/facebook_logo.png'
import githubLogo from '@/assets/logos/github_logo.png'
import googleLogo from '@/assets/logos/google_logo.png'

import { NInput, NButton } from 'naive-ui'
import { useItemScale } from "@/composable/pageAdjuster"
import { signIn } from "@/services/authServices"
import router from "@/router"

const username = ref<string>("")
const password = ref<string>("")

const credentialIncomplete = ref(false)

const scale = useItemScale()

function handleSignIn() {
  if (!username.value || !password.value) {
    credentialIncomplete.value = true
  } else {
    const data = { username: username.value, password: password.value }

    signIn(data)
  }
}

function goToRegister() {
  router.push('/Register')
}

watch([username, password], () => {
  if (credentialIncomplete.value) {
    credentialIncomplete.value = false
  }
})
// watch((username) , (newVal) => {console.log(newVal)}, {deep: true, immediate: true})

</script>
<template>
  <div class="h-screen w-screen flex justify-center items-center bg-cover bg-center bg-no-repeat"
    :style="{ backgroundImage: `url(${loginBackground})` }">
    <div :style="{ transform: `scale(${scale})`, transformOrigin: 'center' }"
      class="px-120 py-20 backdrop-blur-sm rounded-2xl inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] scale-container">
      <div class="w-[430px] min-w-[430px] py-10 bg-white/10 inset-shadow-sm rounded-3xl">
        <div class="px-20">
          <h1 class="font-jakarta mb-10 text-white text-2xl font-bold">Welcome! This is practice</h1>
          <h1 class="font-jakarta mb-6 text-white text-2xl font-bold">Login</h1>
          <div class="mb-4">
            <p class="font-jakrta mb-2 font-light text-[13px] text-white">Username: </p>
            <n-input :status="credentialIncomplete && username.length === 0 ? 'error' : 'success'" type="text"
              v-model:value="username" placeholder="Username" />
          </div>
          <div class="mb-[11px]">
            <p class="font-jakarta mb-2 font-light text-[13px] text-white">Password: </p>
            <n-input :status="credentialIncomplete && password.length === 0 ? 'error' : 'success'" type="password"
              v-model:value="password" placeholder="Password" show-password-on="click" />
          </div>
          <p class="font-jakarta font-light text-[13px] text-white">Forgot Password?</p>
          <div class="my-6">
            <n-button @click="handleSignIn" block color="#003465">Sign in</n-button>
          </div>
          <div class="mb-6 flex justify-center">
            <p class="font-jakarta font-light text-[13px] text-white">or continue with</p>
          </div>
          <div class="flex justify-between">
            <div class="py-2 px-7 rounded-xl bg-white"><img class="w-6 h-6" :src="googleLogo" /></div>
            <div class="py-2 px-7 rounded-xl bg-white"><img class="w-6 h-6" :src="githubLogo" /></div>
            <div class="py-2 px-7 rounded-xl bg-white"><img class="w-6 h-6" :src="facebookLogo" /></div>
          </div>
          <div class="flex justify-center">
            <p class="mt-6 font-jakarta font-light text-[13px] text-white">
              <span>Don’t have an account yet?</span>
              <span v-on:click="goToRegister" class="font-semibold cursor-pointer"> Register for free</span>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700&display=swap');

.scale-container {
  transition: transform 300ms ease-in-out;
}
</style>