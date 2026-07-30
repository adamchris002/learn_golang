<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { NIcon, NInput, NButton, NCheckbox, NAlert } from 'naive-ui';
import TaskDetail from './TaskDetail.vue';
import dayjs from 'dayjs';

import { useItemScale } from '@/composable/pageAdjuster';
import { deleteExistingSubtask, deleteTask, postTodaysTask, updateTaskCompletion, updateTaskValues, type TaskResponse } from '@/services/taskServices';
import customParseFormat from 'dayjs/plugin/customParseFormat'

import send from "@/assets/icons/send.svg"
import calendarIcon from '@/assets/icons/calendar.svg'
import { sanitizeInput } from '@/composable/sanitizeInput';

dayjs.extend(customParseFormat)

type dayTemplate = {
    day: string,
    taskDatas: TaskResponse[]
}

const props = defineProps<{ nextSevenDaysTaskArray: TaskResponse[] }>()
const emits = defineEmits(['requestCallNextSevenDays'])

const openDetails = ref<boolean>(false)
const scale = useItemScale()
const user = JSON.parse(localStorage.getItem("user") || "{}")
let alertTimeout: ReturnType<typeof setTimeout> | null = null

const taskInputArray = ref<string[]>([])
const selectedTaskInformation = ref<TaskResponse | null>(null)

const taskErrorMessage = ref<{ status: number; messageTitle: string; message: string } | null>(null)
const days = computed<dayTemplate[]>(() => {
    const today = dayjs()

    return Array.from({ length: 7 }, (_, i) => {
        const day = today.add(i, "day")

        return {
            day: day.format("DD/MM/YYYY"),
            taskDatas: props.nextSevenDaysTaskArray.filter(
                (data) => data.task_start === day.format("DD/MM/YYYY")
            ),
        }
    })
})

function setSelectedTaskInformation(data: any) {
    openDetails.value = true
    selectedTaskInformation.value = data
}

function closeSelectedTaskInformation() {
    openDetails.value = false
    selectedTaskInformation.value = null
}

function todayOrTomorrowLabel(value: string) {
    const dateValue = dayjs(value, "DD/MM/YYYY").format('dddd')
    if (dateValue === dayjs().format('dddd')) {
        return "Today"
    } else if (dateValue === dayjs().add(1, 'day').format('dddd')) {
        return "Tomorrow"
    } else {
        return
    }
}

function sanitizeTextInput(index: number) {
    const textInput = taskInputArray.value[index]
    if (textInput) {
        const sanitizedText = sanitizeInput(textInput)
        taskInputArray.value[index] = sanitizedText
    }
}

function initializeTaskInputArray() {
    for (let i = 0; i < 7; i++) {
        taskInputArray.value.push("")
    }
}

async function addNewTask(index: number) {
    sanitizeTextInput(index);
    const textString = taskInputArray.value[index]


    const selectedDay = days.value[index]?.day

    if (!selectedDay) {
        return
    }

    if (textString && textString.length > 0) {
        const newData = {
            title: textString,
            description: '',
            task_start: selectedDay,
            due_date: '',
            completed: false,
            userId: user.id
        }
        taskErrorMessage.value = await postTodaysTask(newData)
        if (taskErrorMessage.value.status === 200) {
            emits('requestCallNextSevenDays')
        }
    }
    taskInputArray.value[index] = ""
}

async function handleDeleteTasks(taskId: number) {
    taskErrorMessage.value = await deleteTask(taskId, user.id)
    if (taskErrorMessage.value.status === 200) {
        openDetails.value = false
        emits('requestCallNextSevenDays')
    }
}

async function handleUpdateTasks(dueDate: string, description: string, subtasks: { id: number, title: string, completed: boolean }[], taskId: number) {
    taskErrorMessage.value = await updateTaskValues(dueDate, description, subtasks, taskId, user.id)
    if (taskErrorMessage.value.status === 200) {
        emits('requestCallNextSevenDays')
    }
}

async function handleDeleteSubstasks(subtaskId: number, taskId: number) {
    taskErrorMessage.value = await deleteExistingSubtask(subtaskId, taskId)
    if (taskErrorMessage.value.status === 200) {
        emits('requestCallNextSevenDays')
    }
}

async function updateTaskCheckbox(id: number, data: boolean) {
    taskErrorMessage.value = await updateTaskCompletion(id, user.id, data)
    if (taskErrorMessage.value.status === 200) {
        emits('requestCallNextSevenDays')
    }
}

onMounted(() => {
    emits('requestCallNextSevenDays')
    initializeTaskInputArray()
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

watch(() => props.nextSevenDaysTaskArray, (newValue) => {
    if (selectedTaskInformation.value !== null && newValue.some(data => data.ID === selectedTaskInformation.value?.ID)) {
        const getSelectedTaskInformation = newValue.find(data => data.ID === selectedTaskInformation.value?.ID)
        if (getSelectedTaskInformation) {
            selectedTaskInformation.value = getSelectedTaskInformation
        }
    }
}, { immediate: true, deep: true })
</script>
<template>
    <div class="relative z-10 py-8 px-4 w-full h-screen overflow-x-auto custom-scroll">
        <div
            class="flex items-center backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] size-fit rounded-xl py-2 px-4">
            <n-icon size="16" class="mr-4">
                <calendarIcon class="text-white " />
            </n-icon>
            <p class="font-jakarta text-white text-lg">Next 7 Days</p>
        </div>
        <div class="flex items-start w-full pt-8 scale-container"
            :style="{ transform: `scale(${scale})`, transformOrigin: 'top left' }">
            <div v-for="(value, index) in days" :key="index"
                class="w-[300px] shrink-0 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] p-4 rounded-xl mr-8">
                <div class="flex justify-start">
                    <p
                        :class="[typeof (todayOrTomorrowLabel(value.day)) === 'string' ? 'text-white mr-2 text-base' : '']">
                        {{
                            todayOrTomorrowLabel(value.day) }}</p>
                    <p
                        :class="[typeof (todayOrTomorrowLabel(value.day)) === 'string' ? 'text-[#9c9c9c] text-base' : 'text-white text-base']">
                        {{
                            dayjs(value.day, "DD/MM/YYYY").format('dddd') }}</p>
                </div>
                <div class="my-4 overflow-y-auto max-h-[600px] custom-scroll">
                    <div v-if="value.taskDatas.length > 0" v-for="(items, childIndex) in value.taskDatas"
                        :key="childIndex" :class="[childIndex > 0 ? 'mt-2' : '']">
                        <div @click="setSelectedTaskInformation(value.taskDatas[childIndex])"
                            class="flex justify-start items-center cursor-pointer">
                            <n-checkbox :checked="items.completed"
                                @update-checked="(value: boolean) => updateTaskCheckbox(items.ID, value)" />
                            <div class="ml-2 p-2 bg-[#262626] rounded-md">
                                <p class="font-jakarta text-[#8a8888] text-sm">Created At: {{
                                    dayjs(items.CreatedAt).format("DD/MM/YYYY HH:mm:ss") }}</p>
                                <p class="text-white font-jakarta text-base">{{ items.title }}</p>
                            </div>
                        </div>
                    </div>
                </div>
                <n-input block v-model:value="taskInputArray[index]" @blur="sanitizeTextInput(index)" type="textarea"
                    :resizable="false" placeholder="+ Add Task"
                    :theme-overrides="{ color: 'backdrop-blur-sm', borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', colorFocus: 'backdrop-blur-sm', textColor: 'white' }">
                    <template #suffix>
                        <n-button :bordered="false" @click="addNewTask(index)"
                            :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', rippleColor: 'none' }"
                            class="send-btn">
                            <n-icon>
                                <send />
                            </n-icon>
                        </n-button>
                    </template>
                </n-input>
            </div>
        </div>
        <div class="absolute top-2 right-2 z-9999">
            <n-alert v-if="taskErrorMessage" :title="taskErrorMessage.messageTitle"
                :type="taskErrorMessage.status === 200 ? 'success' : 'error'" closable @close="taskErrorMessage = null">
                {{ taskErrorMessage.message }}
            </n-alert>
        </div>
    </div>
    <TaskDetail :open="openDetails" :task="selectedTaskInformation" @close-modal="closeSelectedTaskInformation"
        @delete-task="handleDeleteTasks" @update-task-datas="handleUpdateTasks"
        @delete-subtask="handleDeleteSubstasks" />
</template>
<style scoped>
.scale-container {
    transition: transform 300ms ease-in-out;
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

:deep(.n-input__suffix) {
    margin-right: -12px !important;
}

:deep(.n-input__textarea.n-scrollbar) {
    height: 34px !important
}

.custom-scroll::-webkit-scrollbar {
    width: 6px;
}

.custom-scroll::-webkit-scrollbar-track {
    background: transparent;
}

.custom-scroll::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 9999px;
}

.custom-scroll::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
}
</style>