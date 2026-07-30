<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue"
import { deleteExistingSubtask, deleteTask, getTodaysTasks, postTodaysTask, updateTaskCompletion, updateTaskValues, type TaskResponse } from "@/services/taskServices"

import { useTimeNow } from "@/composable/timeNow"
import { useItemScale } from "@/composable/pageAdjuster"

import { NInput, NButton, NIcon, NAlert } from "naive-ui"
import send from "@/assets/icons/send.svg"
import TaskLists from "./TaskLists.vue"
import TaskDetail from "./TaskDetail.vue"
import { sanitizeInput } from "@/composable/sanitizeInput.ts"
import dayjs from "dayjs"

const props = defineProps<{ todaysTaskArray: TaskResponse[] }>()
const emits = defineEmits(['requestCallTodaysTask'])

let alertTimeout: ReturnType<typeof setTimeout> | null = null

const { currentTime } = useTimeNow()
const scale = useItemScale()

const taskErrorMessage = ref<{ status: number; messageTitle: string; message: string } | null>(null)
const textString = ref<string>("")
const openDetails = ref<boolean>(false)
const taskInformation = ref<TaskResponse | null>(null)

const user = JSON.parse(localStorage.getItem("user") || "{}")

const completedTasks = computed(() => { return props.todaysTaskArray.filter((item) => item.completed) })
const incompleteTasks = computed(() => { return props.todaysTaskArray.filter((item) => !item.completed) })

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

function openTaskInformation(taskItem: TaskResponse) {
    openDetails.value = true
    taskInformation.value = taskItem
}

function closeTaskInformation() {
    openDetails.value = false
    taskInformation.value = null
}

function sanitizeTaskInput() {
    //sanitize input hehehe
    textString.value = sanitizeInput(textString.value)
}

async function handleDeleteSubtask(subtaskId: number, taskId: number) {
    taskErrorMessage.value = await deleteExistingSubtask(subtaskId, taskId)
    emits('requestCallTodaysTask')
}

async function handleDeleteTask(id: number) {
    taskErrorMessage.value = await deleteTask(id, user.id)
    closeTaskInformation()
    emits('requestCallTodaysTask')
}

async function handleUpdateTaskValues(taskDueDate: string, taskDescription: string, taskSubArray: { id: number, title: string, completed: boolean }[], taskId: number) {
    taskErrorMessage.value = await updateTaskValues(taskDueDate, taskDescription, taskSubArray, taskId, user.id)
    emits('requestCallTodaysTask')
}

async function handleUpdateTaskCompletion(id: number, data: boolean) {
    taskErrorMessage.value = await updateTaskCompletion(id, user.id, data)
    if (taskErrorMessage.value.status === 200) {
        emits('requestCallTodaysTask')
    }
}

async function sendTask() {
    sanitizeTaskInput()
    if (textString.value.length > 0) {
        const newData = {
            title: textString.value,
            description: '',
            task_start: dayjs().format("DD/MM/YYYY"),
            due_date: '',
            completed: false,
            userId: user.id
        }
        taskErrorMessage.value = await postTodaysTask(newData)
        if (taskErrorMessage.value.status === 200) {
            emits('requestCallTodaysTask')
            textString.value = ""
        }
    } else {
        taskErrorMessage.value = { status: 400, messageTitle: "Save Task Failed", message: "Task cannot be empty." };
    }
}

onMounted(() => {
    emits('requestCallTodaysTask')
})

watch(() => taskErrorMessage.value, (message) => {
    if (alertTimeout) {
        clearTimeout(alertTimeout)
    }

    if (message) {
        alertTimeout = setTimeout(() => {
            taskErrorMessage.value = null
        }, 5000)
    }
},
    { immediate: true })

watch(() => props.todaysTaskArray, (newValue) => {
    if (taskInformation.value !== null && newValue.some(data => data.ID === taskInformation.value?.ID)) {
        const getSelectedTaskInformation = newValue.find(data => data.ID === taskInformation.value?.ID)
        if (getSelectedTaskInformation) {
            taskInformation.value = getSelectedTaskInformation
        }
    }
}, { immediate: true, deep: true })

</script>
<template>
    <div class="w-full h-screen flex flex-col items-center py-10 scale-container"
        :style="{ transform: `scale(${scale})`, transformOrigin: scale > 0.7 ? 'center' : 'top' }">
        <div>
            <h1 class="text-4xl font-jakarta">
                <span class="text-white">
                    {{ getGreetings() }}, {{ user.first_name }}
                </span>
                <span class="text-[#0373fc]">.</span>
            </h1>
            <h1 class="font-jakarta text-[#8a8888] text-4xl">Remove doubts with action</h1>
            <div
                class="mt-8 rounded-lg w-200 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-2 flex justify-start">
                <div class="flex flex-col items-center">
                    <p class="font-jakarta text-base font-semibold text-white">{{
                        currentTime.format("ddd").toUpperCase() }}</p>
                    <p class="font-jakarta text-5xl font-bold text-white">{{ currentTime.format("DD") }}</p>
                    <p class="font-jakarta text-base font-normal text-white">{{ currentTime.format("MMMM") }}</p>
                </div>
                <div class="w-full ml-4 flex items-center justify-center">
                    <h1 class="font-jakarta text-4xl font-bold text-white">What's your tasks today?</h1>
                </div>
            </div>
        </div>
        <TaskLists title="Completed Tasks" if-empty-string="No completed Tasks yet" :task-array="completedTasks"
            @toggle="handleUpdateTaskCompletion" @open-task-detail="openTaskInformation" :scale="scale" />
        <TaskLists title="Not Completed Tasks" if-empty-string="No Tasks to complete" :task-array="incompleteTasks"
            @toggle="handleUpdateTaskCompletion" @open-task-detail="openTaskInformation" :scale="scale" />
        <div class="mt-auto w-200 flex">
            <n-input block round v-model:value="textString" @blur="sanitizeTaskInput"
                :theme-overrides="{ color: 'backdrop-blur-sm', borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', colorFocus: 'backdrop-blur-sm', textColor: 'white' }"
                placeholder="What needs to be done?">
                <template #suffix>
                    <n-button circle :bordered="false" class="send-btn"
                        :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', rippleColor: 'none' }"
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
    <TaskDetail :open="openDetails" :task="taskInformation" @close-modal="closeTaskInformation"
        @delete-task="handleDeleteTask" @update-task-datas="handleUpdateTaskValues"
        @delete-subtask="handleDeleteSubtask" />
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