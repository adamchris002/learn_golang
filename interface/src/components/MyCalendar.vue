<script setup lang="ts">
import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat'
import { useCalendarSizeAdjuster } from '@/composable/pageAdjuster';

//components
import MyCalendarTaskDetail from '@/components/MyCalendarTaskDetail.vue'

import chevronLeftIcon from '@/assets/icons/chevron-left.svg';
import chevronRightIcon from '@/assets/icons/chevron-right.svg';
import { NButton, NDropdown, NIcon, NCalendar, NConfigProvider, darkTheme, NAlert } from 'naive-ui';
import { onMounted, ref, watch } from 'vue';
import { deleteExistingSubtask, deleteTask, getMonthTasks, getWeekTasks, updateTaskValues, type TaskResponse } from '@/services/taskServices';
import TaskDetail from './TaskDetail.vue';

dayjs.extend(customParseFormat)

type Option = {
    label: string;
    key: string;
};

const user = JSON.parse(localStorage.getItem("user") || "{}")
let alertTimeout: ReturnType<typeof setTimeout> | null = null

const calendarSize = useCalendarSizeAdjuster()
const selectedDay = ref(dayjs())
const selectedOption = ref<Option>({ label: 'Month', key: 'month' })

// label datas
const weekDate = ref<string[]>([])
const timeDate = ref<string[]>([])

//datas
const calendarData = ref<any[]>([])
const weekData = ref<any[]>([])
const dayData = ref<any[]>([])

const selectedTaskInformation = ref<TaskResponse | null>(null)


//datas for props
const dateDatas = ref<string>('')

//error message
const taskErrorMessage = ref<{ status: number; messageTitle: string; message: string } | null>(null)

const options = [{ label: 'Month', key: 'month' }, { label: 'Week', key: 'week' }, { label: 'Day', key: 'day' }, { label: 'Agenda', key: 'agenda' }]

const openTaskInformation = ref<boolean>(false)
const openCalendarTaskDetail = ref<boolean>(false)

function chevronRightAction() {
    switch (selectedOption.value?.key) {
        case 'week':
            selectedDay.value = selectedDay.value.add(1, 'week').startOf('week').add(1, 'day')
            break
        case 'day':
            selectedDay.value = selectedDay.value.add(1, 'day')
            break
    }
}

function chevronLeftAction() {
    switch (selectedOption.value?.key) {
        case 'week':
            selectedDay.value = selectedDay.value.subtract(1, 'week').startOf('week').add(1, 'day')
            break
        case 'day':
            selectedDay.value = selectedDay.value.subtract(1, 'day')
            break
    }
}

function handleOptionSelect(option: string | number) {
    const found = options.find((opt) => opt.key === option)
    if (found) {
        selectedOption.value = found
    }
    selectedDay.value = dayjs()
}

function getWeek() {
    const startOfWeek = selectedDay.value.startOf('week').add(1, 'day')
    weekDate.value = []
    for (let i = 0; i < 7; i++) {
        weekDate.value.push(startOfWeek.add(i, 'day').format('DD/MM/YYYY'))
    }
}

function generateTime() {
    const startOfDay = selectedDay.value.startOf('day')
    const timeArray: string[] = []
    for (let i = 0; i < 24; i++) {
        const time = startOfDay.add(i, 'hour').format('HH:mm')
        timeArray.push(time)
    }
    timeDate.value = timeArray
}

function openTaskModal(_: number,
    { year, month, date }: { year: number, month: number, date: number }) {
    if (selectedOption.value?.key === 'month') {
        const formattedMonth = String(month).padStart(2, '0')
        dateDatas.value = `${date}/${formattedMonth}/${year}`
    }
    openCalendarTaskDetail.value = !openCalendarTaskDetail.value
}

function openTaskInformationModal(data: TaskResponse) {
    openTaskInformation.value = true
    selectedTaskInformation.value = data
}

function closeSelectedTaskInformation() {
    openTaskInformation.value = false
    selectedTaskInformation.value = null
}

function closeTaskModal() {
    openCalendarTaskDetail.value = false
    dateDatas.value = ''
}

function getTasksForCalendarDate(year: number, month: number, date: number) {
    const calendarDate = `${String(date).padStart(2, '0')}/${String(month).padStart(2, '0')}/${String(year)}`

    return calendarData.value.filter(data => data.due_date === calendarDate)
}

async function fetchData(key: string) {
    switch (key) {
        case 'month': {
            const result = await getMonthTasks(selectedDay.value.format("DD/MM/YYYY"), user.id)
            calendarData.value = result.data
            return
        }
        case 'week': {
            const result = await getWeekTasks(selectedDay.value.format("DD/MM/YYYY"), user.id)
            weekData.value = result.data
            return
        }
        case 'day': {

        }
        default:
            return

    }
}

async function handleUpdateTasks(dueDate: string, description: string, subtasks: { id: number, title: string, completed: boolean }[], taskId: number) {
    taskErrorMessage.value = await updateTaskValues(dueDate, description, subtasks, taskId, user.id)
    if (taskErrorMessage.value.status === 200) {
        fetchData(selectedOption.value.key)
    }
}

async function handleResultMessage(data: { status: number, message: string, messageTitle: string }) {
    taskErrorMessage.value = data
    if (taskErrorMessage.value.status === 200) {
        closeTaskModal()
        fetchData(selectedOption.value.key)
    }
}

async function handleDeleteTask(taskId: number) {
        taskErrorMessage.value = await deleteTask(taskId, user.id)
    if (taskErrorMessage.value.status === 200) {
        openTaskInformation.value = false
        fetchData(selectedOption.value.key)
    }
}

async function handleDeleteSubTask(subtaskId: number, taskId: number) {
    taskErrorMessage.value = await deleteExistingSubtask(subtaskId, taskId)
    if (taskErrorMessage.value.status === 200) {
        fetchData(selectedOption.value.key)
    }
}

onMounted(() => {
    generateTime()
})

watch(
    () => selectedDay.value,
    () => {
        getWeek()
    },
    { immediate: true }
)

watch(() => selectedOption.value, (newValue) => {
    if (newValue) {
        fetchData(newValue.key)
    }
}, { immediate: true })

watch(() => calendarData.value, (newValue) => {
    if (selectedTaskInformation.value !== null && newValue.some(data => data.ID === selectedTaskInformation.value?.ID)) {
        const getSelectedTaskInformation = newValue.find(data => data.ID === selectedTaskInformation.value?.ID)
        if (getSelectedTaskInformation) {
            selectedTaskInformation.value = getSelectedTaskInformation
        }
    }
}, { immediate: true })

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
</script>
<template>
    <div class="z-10 w-full px-4 py-8 h-screen">
        <div :class="['flex items-center', selectedOption?.key !== 'month' ? 'justify-between' : 'justify-end']">
            <div v-if="selectedOption?.key !== 'month'" class="flex items-center">
                <div class="cursor-pointer rounded-4xl border-1 border-[#a3a3a3] bg-[#1c1c1c] px-4 py-2">
                    <p class="text-[#a3a3a3] text-xl">Today</p>
                </div>
                <n-button @click="chevronLeftAction" circle ghost :bordered="false" class="chevron-btn size-fit"
                    :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', rippleColor: 'none' }">
                    <n-icon size="36">
                        <chevronLeftIcon class="text-4xl" />
                    </n-icon>
                </n-button>
                <n-button @click="chevronRightAction" circle ghost :bordered="false" class="chevron-btn size-fit"
                    :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', rippleColor: 'none' }">
                    <n-icon size="36">
                        <chevronRightIcon class="text-4xl" />
                    </n-icon>
                </n-button>
                <p :class="['text-white font-jakarta text-4xl', selectedOption?.key !== 'month' ? '' : 'pl-4']">{{
                    selectedDay.format('MMMM YYYY') }}</p>
            </div>
            <div>
                <n-dropdown trigger="click" :options="options" @select="(option) => handleOptionSelect(option)">
                    <n-button>
                        <p class="text-[#a3a3a3]">
                            {{ selectedOption?.label ? selectedOption.label : 'Select View' }}
                        </p>
                    </n-button>
                </n-dropdown>
            </div>
        </div>
        <div v-show="selectedOption?.key === 'month'" class="w-full backdrop-blur-sm"
            :style="{ transform: `scaleX(${calendarSize.scaleX}) scaleY(${calendarSize.scaleY})`, transformOrigin: 'top' }">
            <n-config-provider :theme="darkTheme">
                <n-calendar :key="selectedDay.format('YYYY-MM')" :value="selectedDay.valueOf()"
                    @update:value="openTaskModal">
                    <template #="{ year, month, date }">
                        <div class="max-h-[100px] overflow-y-auto custom-calendar-scroll">
                            <div v-for="task in getTasksForCalendarDate(year, month, date)" :key="task.title"
                                class=" relative z-9999 border-1 border-[#3a3a3a] rounded-lg mb-2 max-w-[130px]"
                                @click.stop="openTaskInformationModal(task)">
                                <p class="font-jakarta p-2 truncate w-32">{{ task.title }}</p>
                            </div>
                        </div>
                    </template>
                </n-calendar>
            </n-config-provider>
        </div>
        <div v-show="selectedOption?.key === 'week'"
            class="w-full max-h-[75vh] overflow-y-auto bg-[#1c1c1c] backdrop-blur-sm custom-scroll">
            <table class="w-full table-fixed">
                <thead>
                    <tr class="sticky top-0 z-10 bg-[#1c1c1c]">
                        <th class="w-16 p-4"></th>
                        <th class="p-4" v-for="day in weekDate" :key="day">
                            <div :class="[
                                'flex flex-col items-center justify-center w-12 h-12 mx-auto text-white',
                                dayjs(day, 'DD/MM/YYYY').isSame(dayjs(), 'day') ? 'bg-blue-500 rounded-full' : ''
                            ]">
                                <p>
                                    {{ dayjs(day, 'DD/MM/YYYY').format('ddd') }}
                                </p>

                                <p :class="[dayjs(day, 'DD/MM/YYYY').isSame(dayjs(), 'day') ? 'text-xl' : '']">
                                    {{ dayjs(day, 'DD/MM/YYYY').format('DD') }}
                                </p>
                            </div>
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="(time, index) in timeDate" :key="time">
                        <td class="relative">
                            <p v-if="index !== 0" class="absolute -top-3 text-white text-sm">
                                {{ dayjs(time, 'HH:mm').format("h:mm A") }}
                            </p>
                        </td>
                        <td v-for="day in weekDate" :key="day" class="border border-gray-600 h-16">
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div v-show="selectedOption?.key === 'day'"
            class="w-full max-h-[75vh] overflow-y-auto bg-[#1c1c1c] backdrop-blur-sm custom-scroll">
            <div class="w-full flex justify-center items-center sticky top-0 z-10 bg-[#1c1c1c]"">
                <div :class="['w-12 h-12 mx-auto flex flex-col items-center justify-center',
                    dayjs(selectedDay).isSame(dayjs(), 'day') ? 'bg-blue-500 rounded-full' : '']">
                <p class="text-white">{{ dayjs(selectedDay).format('ddd') }}</p>
                <p class="text-white text-xl">{{ dayjs(selectedDay).format('DD') }}</p>
            </div>
        </div>
        <table class="w-full table-fixed">
            <thead>
                <tr class="sticky top-0 z-10 bg-[#1c1c1c]">
                    <th class="w-16 p-4"></th>
                </tr>
            </thead>
            <tbody>
                <!-- row to give space for time labels -->
                <tr>
                    <td class="h-4"></td>
                </tr>
                <tr v-for="(time, index) in timeDate" :key="time">
                    <td class="relative">
                        <p class="absolute text-white text-sm -top-3">
                            {{ dayjs(time, 'HH:mm').format("h:mm A") }}
                        </p>
                    </td>
                    <td v-for="day in weekDate" :key="day" class="border-t border-gray-600 h-16">
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    <div v-show="selectedOption?.key === 'agenda'">

    </div>
    </div>
    <MyCalendarTaskDetail :open="openCalendarTaskDetail" @close-modal="closeTaskModal" :dateData="dateDatas"
        @result-message="handleResultMessage" />
    <TaskDetail :open="openTaskInformation" :task="selectedTaskInformation" @close-modal="closeSelectedTaskInformation"
        @update-task-datas="handleUpdateTasks" @delete-task="handleDeleteTask" @delete-subtask="handleDeleteSubTask" />
    <div class="absolute top-2 right-2 z-9999">
        <n-alert v-if="taskErrorMessage" :title="taskErrorMessage.messageTitle"
            :type="taskErrorMessage.status === 200 ? 'success' : 'error'" closable @close="taskErrorMessage = null">
            {{ taskErrorMessage.message }}
        </n-alert>
    </div>
</template>
<style scoped>
.scale-container {
    transition: transform 300ms ease-in-out;
}

.chevron-btn {
    color: white;
}

.chevron-btn:hover {
    color: #0373fc;
}

.chevron-btn:focus {
    color: #0373fc;
}

:deep(.n-calendar-header__title) {
    font-size: 3rem;
}

:deep(.week-table th) {
    white-space: pre-line;
}

:deep(.n-data-table-th__title) {
    text-align: center !important;
}

.custom-scroll::-webkit-scrollbar {
    display: none;
}

.custom-calendar-scroll::-webkit-scrollbar {
    width: 6px;
}

.custom-calendar-scroll::-webkit-scrollbar-track {
    background: transparent;
}

.custom-calendar-scroll::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 9999px;
}

.custom-calendar-scroll::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
}
</style>