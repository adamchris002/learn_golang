<script setup lang="ts">
import { ref, watch } from "vue"

import loginBackground from "@/assets/images/background.png"

import { NInput, NButton } from 'naive-ui'
import { useItemScale } from "@/composable/pageAdjuster"
import { register } from "@/services/authServices"
import router from "@/router"

const firstName = ref<string>("")
const lastName = ref<string>("")
const username = ref<string>("")
const password = ref<string>("")

const credentialIncomplete = ref(false)

const scale = useItemScale()

function handleRegister() {
    if (!username.value || !password.value) {
        credentialIncomplete.value = true
    } else {
        const data = { firstName: firstName.value, lastName: lastName.value, username: username.value, password: password.value }

        register(data)
    }
}

function goToLogin() {
    router.push('/Login')
}

watch([username, password, firstName, lastName], () => {
    if (credentialIncomplete.value) {
        credentialIncomplete.value = false
    }
})

</script>
<template>
    <div class="h-screen w-screen flex justify-center items-center bg-cover bg-center bg-no-repeat"
        :style="{ backgroundImage: `url(${loginBackground})` }">
        <div :style="{ transform: `scale(${scale})`, transformOrigin: 'center' }"
            class="px-120 py-20 backdrop-blur-sm rounded-2xl inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] scale-container">
            <div class="w-[430px] min-w-[430px] py-10 bg-white/10 inset-shadow-sm rounded-3xl">
                <div class="px-20">
                    <h1 class="font-jakarta mb-10 text-white text-2xl font-bold">Welcome! This is practice</h1>
                    <h1 class="font-jakarta mb-6 text-white text-2xl font-bold">Register</h1>
                    <div class="mb-4">
                        <p class="font-jakrta mb-2 font-light text-[13px] text-white">First Name: </p>
                        <n-input :status="credentialIncomplete && firstName.length === 0 ? 'error' : 'success'"
                            type="text" v-model:value="firstName" placeholder="First Name" />
                    </div>
                    <div class="mb-4">
                        <p class="font-jakrta mb-2 font-light text-[13px] text-white">Last Name: </p>
                        <n-input :status="credentialIncomplete && lastName.length === 0 ? 'error' : 'success'"
                            type="text" v-model:value="lastName" placeholder="Last Name" />
                    </div>
                    <div class="mb-4">
                        <p class="font-jakrta mb-2 font-light text-[13px] text-white">Username: </p>
                        <n-input :status="credentialIncomplete && username.length === 0 ? 'error' : 'success'"
                            type="text" v-model:value="username" placeholder="Username" />
                    </div>
                    <div class="mb-[11px]">
                        <p class="font-jakarta mb-2 font-light text-[13px] text-white">Password: </p>
                        <n-input :status="credentialIncomplete && password.length === 0 ? 'error' : 'success'"
                            type="password" v-model:value="password" placeholder="Password" show-password-on="click" />
                    </div>
                    <div class="my-6">
                        <n-button @click="handleRegister" block color="#003465">Sign in</n-button>
                    </div>
                    
                    <div class="flex justify-center">
                        <p class="mt-6 font-jakarta font-light text-[13px] text-white">
                            <span>Already have an account?</span>
                            <span v-on:click="goToLogin" class="font-semibold cursor-pointer"> Login here</span>
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