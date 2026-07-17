<script setup lang="ts">
import { onMounted, ref, watch } from "vue"
import { getTodaysTasks, postTodaysTask, updateTaskCompletion, type TaskResponse } from "@/services/taskServices"

import { useTimeNow } from "@/composable/timeNow"
import { useItemScale } from "@/composable/pageAdjuster"

import { NInput, NButton, NIcon, NAlert, NCheckbox } from "naive-ui"
import send from "@/assets/icons/send.svg"
import dayjs from "dayjs"

const { currentTime } = useTimeNow()
const scale = useItemScale()

const taskArray = ref<TaskResponse[]>([])
const taskErrorMessage = ref<{ status: number; messageTitle: string; message: string } | null>(null)
const textString = ref<string>("")

const refreshTask = ref(false)

const user = JSON.parse(localStorage.getItem("user") || "{}")

function getGreetings() {
    const hour = currentTime.value.hour()
    if (hour >= 5 && hour < 12) {
        return "Good Morning";
    } else if (hour >= 12 && hour < 18) {
        return "Good Afternoon";
    } else {
        return "Good Evening";
    }
}

function setTextInput(value: string) {
    //sanitize input hehehe
    textString.value = value.trim().replace(/\s+/g, " ");
}

async function handleUpdateTaskCompletion(id: number, data: boolean) {
    taskErrorMessage.value = await updateTaskCompletion(id, user.id, data)
}

async function sendTask() {
    if (textString.value.length > 0) {
        const newData = {
            title: textString.value,
            description: '',
            dueDate: '',
            completed: false,
            userId: user.id
        }
        taskErrorMessage.value = await postTodaysTask(newData)
        if (taskErrorMessage.value.status === 200) {
            refreshTask.value = true
        }
    } else {
        taskErrorMessage.value = { status: 400, messageTitle: "Save Task Failed", message: "Task cannot be empty." };
    }
}

async function callTodaysTasks() {
    const result = await getTodaysTasks(user.id)
    if (result) {
        taskArray.value = result.data.sort(
            (a: TaskResponse, b: TaskResponse) => new Date(a.CreatedAt).getTime() - new Date(b.CreatedAt).getTime()
        )
        taskErrorMessage.value = result.messageData ?? null
    } else {
        taskErrorMessage.value = { status: 400, messageTitle: "Fetch Tasks Failed", message: "Unable to fetch today's tasks." };
    }
}

onMounted(async () => {
    await callTodaysTasks()
})

watch((refreshTask), async (newValue) => {
    if (newValue) {
        await getTodaysTasks(user.id)
    }
}, { immediate: true })
</script>
<template>
    <div class="w-full h-screen flex flex-col items-center py-10 scale-container"
        :style="{ transform: `scale(${scale})`, transformOrigin: 'center' }">
        <div>
            <h1 class="text-4xl jakarta-font">
                <span class="text-white">
                    {{ getGreetings() }}, {{ user.first_name }}
                </span>
                <span class="text-[#0373fc]">.</span>
            </h1>
            <h1 class="jakarta-font text-[#8a8888] text-4xl">Remove doubts with action</h1>
            <div
                class="mt-8 rounded-lg w-200 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-2 flex justify-start">
                <div class="flex flex-col items-center">
                    <p class="jakarta-font text-base font-semibold text-white">{{
                        currentTime.format("ddd").toUpperCase() }}</p>
                    <p class="jakarta-font text-5xl font-bold text-white">{{ currentTime.format("DD") }}</p>
                    <p class="jakarta-font text-base font-normal text-white">{{ currentTime.format("MMMM") }}</p>
                </div>
                <div class="w-full ml-4 flex items-center justify-center">
                    <h1 class="jakarta-font text-4xl font-bold text-white">What's your tasks today?</h1>
                </div>
            </div>
        </div>
        <div v-if="taskArray.length > 0 && taskArray" v-for="items in taskArray">
            <div
                class="mt-8 rounded-lg w-200 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-2 flex items-center">
                <n-checkbox :checked="items.completed" @update-checked="(value) => handleUpdateTaskCompletion(items.ID, value)" />
                <div class="cursor-pointer" @click="">
                    <p class="jakarta-font text-[#8a8888] text-base ml-4">Created at: {{
                        dayjs(items.CreatedAt).format("DD/MMMM/YYYY HH:mm:ss") }}</p>
                    <p class="jakarta-font text-white text-xl ml-4">{{ items.title }}</p>
                </div>
            </div>
        </div>
        <div class="mt-auto w-200 flex">
            <n-input block round
                :theme-overrides="{ color: 'backdrop-blur-sm', borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', colorFocus: 'backdrop-blur-sm', textColor: 'white' }"
                @update:value="setTextInput" placeholder="What needs to be done?">
                <template #suffix>
                    <n-button circle :bordered="false" class="send-btn"
                        :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc' }"
                        @click="sendTask">
                        <n-icon>
                            <send />
                        </n-icon>
                    </n-button>
                </template>
            </n-input>
        </div>
        <div class="absolute top-2 right-2 z-9999">
            <n-alert v-if="taskErrorMessage" :title="taskErrorMessage.messageTitle"
                :type="taskErrorMessage.status === 200 ? 'success' : 'error'" closable @close="taskErrorMessage = null">
                {{ taskErrorMessage.message }}
            </n-alert>
        </div>
    </div>
</template>
<style scoped>
.scale-container {
    transition: transform 300ms ease-in-out;
}

:deep(.n-input__suffix) {
    margin-right: -12px !important;
}

.send-btn {
    color: white;
}

.send-btn:hover {
    color: #0373fc;
}

.send-btn:focus {
    color: #0373fc;
}
</style>