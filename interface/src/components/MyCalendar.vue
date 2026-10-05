<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import { NButton, NDropdown, NIcon, NCalendar, NConfigProvider, darkTheme, NAlert, NCheckbox } from 'naive-ui';
import { VueDraggable, type DraggableEvent, type MoveEvent } from 'vue-draggable-plus';
import MyCalendarTaskDetail from '@/components/MyCalendarTaskDetail.vue'
import TaskDetail from './TaskDetail.vue';

import { changeTaskDueDate, deleteExistingSubtask, deleteTask, getDayTasks, getMonthTasks, getWeekTasks, updateTaskCompletion, updateTaskValues, type TaskResponse } from '@/services/taskServices';
import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat'
import { useCalendarSizeAdjuster } from '@/composable/pageAdjuster';

import chevronLeftIcon from '@/assets/icons/chevron-left.svg';
import chevronRightIcon from '@/assets/icons/chevron-right.svg';

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
const calendarData = ref<TaskResponse[]>([])
const weekData = ref<TaskResponse[]>([])
const dayData = ref<TaskResponse[]>([])

const selectedTaskInformation = ref<TaskResponse | null>(null)

//datas for props
const dateDatas = ref<string>('')
const typeDatas = ref<string>('')

//error message
const taskErrorMessage = ref<{ status: number; messageTitle: string; message: string } | null>(null)

const timeVariations = [0, 15, 30, 45]
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
    const selectedDate = dayjs(`${year}-${month}-${date}`, 'YYYY-M-D')

    if (selectedDate.isBefore(dayjs(), 'day')) {
        return
    }

    const formattedMonth = String(month).padStart(2, '0')
    const formattedDay = String(date).padStart(2, '0')
    dateDatas.value = `${formattedDay}/${formattedMonth}/${year}`
    typeDatas.value = 'month'

    openCalendarTaskDetail.value = !openCalendarTaskDetail.value
}

function openTaskModalWeekAndDate(day: string, time: string, minutes: number) {
    const selectedDate = dayjs(`${day}`, 'DD/MM/YYYY')

    if (selectedDate.isBefore(dayjs(), 'day')) {
        return
    }
    dateDatas.value = dayjs(day + " " + time, 'DD/MM/YYYY HH:mm').add(minutes, 'minutes').format("DD/MM/YYYY HH:mm")
    typeDatas.value = selectedOption.value.key

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
    typeDatas.value = ''
}

function getTasksForCalendarDate(year: number, month: number, date: number) {
    const calendarDate = `${String(date).padStart(2, '0')}/${String(month).padStart(2, '0')}/${String(year)}`

    return calendarData.value.filter(data => dayjs(data.due_date, 'DD/MM/YYYY HH:mm').format("DD/MM/YYYY") === calendarDate)
}

function getTasksForCalendarWeek(day: string, time: string, times: number) {
    const dayTimeString = `${day} ${time}`
    const start = dayjs(dayTimeString, "DD/MM/YYYY HH:mm").add(times, 'minutes')
    const end = start.add(15, 'minutes')

    return weekData.value.filter(data => {
        const dueDate = dayjs(data.due_date, "DD/MM/YYYY HH:mm")

        return dueDate.isSameOrAfter(start) && dueDate.isBefore(end)
    })
}

function getTasksForCalendarDay(time: string, times: number) {
    const day = dayjs(selectedDay.value, "DD/MM/YYYY HH:mm").format('DD/MM/YYYY')
    const dayTimeString = `${day} ${time}`
    const start = dayjs(dayTimeString, "DD/MM/YYYY HH:mm").add(times, 'minutes')
    const end = start.add(15, 'minutes')
    return dayData.value.filter(data => {
        const dueDate = dayjs(data.due_date, "DD/MM/YYYY HH:mm")

        return dueDate.isSameOrAfter(start) && dueDate.isBefore(end)
    })
}

function checkIfUpdateStartDateWeekPossible(event: MoveEvent) {
    const taskId = Number(event.dragged.dataset.taskId)
    const targetDay = event.to.dataset.dayTime

    let selectedData: TaskResponse | undefined

    selectedData = weekData.value.find((data) => data.ID === taskId)

    if (!selectedData) return false
    const selectedDate = dayjs(targetDay, "DD/MM/YYYY HH:mm")
    const taskStartDate = dayjs(selectedData.task_start, "DD/MM/YYYY HH:mm")

    if (taskStartDate.isAfter(selectedDate) || selectedDate.isBefore(dayjs().startOf('day'))) {
        // taskErrorMessage.value = {
        //     message: 'Task due date cannot be set before its start date.',
        //     messageTitle: 'Edit Due Date Failed',
        //     status: 400
        // }
        return false
    }

    return true
}

function checkIfUpdateStartDatePossible(event: MoveEvent) {
    const taskId = Number(event.dragged.dataset.taskId)
    const targetDay = event.to.dataset.day

    let selectedData: TaskResponse | undefined

    selectedData = calendarData.value.find((data) => data.ID === taskId)

    if (!selectedData) return false

    const selectedDate = dayjs(targetDay, "DD/MM/YYYY")
    const taskStartDate = dayjs(selectedData.task_start, "DD/MM/YYYY")

    if (taskStartDate.isAfter(selectedDate) || selectedDate.isBefore(dayjs().startOf('day'))) {
        // taskErrorMessage.value = {
        //     message: 'Task due date cannot be set before its start date.',
        //     messageTitle: 'Edit Due Date Failed',
        //     status: 400
        // }
        return false
    }
    return true
}

function checkIfUpdateStartDateDayPossible(event: MoveEvent) {
    const taskId = Number(event.dragged.dataset.taskId)
    const targetDay = event.to.dataset.dateTime

    let selectedData: TaskResponse | undefined

    selectedData = calendarData.value.find((data) => data.ID === taskId)

    if (!selectedData) return false
    const selectedDate = dayjs(targetDay, "DD/MM/YYYY HH:mm")
    const taskStartDate = dayjs(selectedData.task_start, "DD/MM/YYYY HH:mm")

    if (taskStartDate.isAfter(selectedDate) || selectedDate.isBefore(dayjs().startOf('day'))) {
        return false
    }
    return true
}

function handleChangeToToday() {
    selectedDay.value = dayjs()
}

async function updateTaskCheckbox(id: number, data: boolean) {
    taskErrorMessage.value = await updateTaskCompletion(id, user.id, data)
    if (taskErrorMessage.value.status === 200) {
        fetchData(selectedOption.value.key)
    }
}

async function handleUpdateTaskDueDateV2(event: DraggableEvent<TaskResponse>, dateTimeData: string) {
    const taskId = event.data.ID

    const getSelectedTask = weekData.value.find((data) => data.ID === taskId)

    if (!getSelectedTask) {
        return
    }

    taskErrorMessage.value = await changeTaskDueDate(
        taskId,
        user.id,
        dateTimeData
    )

    if (taskErrorMessage.value.status === 200) {
        fetchData(selectedOption.value.key)
    }
}

async function handleUpdateTaskDueDate(
    event: DraggableEvent<TaskResponse>,
    year: number,
    month: number,
    day: number
) {
    const taskId = event.data.ID
    const targetDay = `${String(day).padStart(2, '0')}/${String(month).padStart(2, '0')}/${year}`

    const getSelectedTask = calendarData.value.find(
        (data) => data.ID === taskId
    )

    if (!getSelectedTask) {
        return
    }

    const originalDueDate = dayjs(
        getSelectedTask.due_date,
        "DD/MM/YYYY HH:mm"
    )

    const newDueDate = dayjs(
        targetDay,
        "DD/MM/YYYY"
    )
        .hour(originalDueDate.hour())
        .minute(originalDueDate.minute()).format("DD/MM/YYYY HH:mm")

    taskErrorMessage.value = await changeTaskDueDate(
        taskId,
        user.id,
        newDueDate
    )

    if (taskErrorMessage.value.status === 200) {
        fetchData(selectedOption.value.key)
    }
}


async function handlePanelChange({ year, month }: { year: number, month: number }) {
    selectedDay.value = dayjs(`${year}-${String(month).padStart(2, '0')}-01`, 'YYYY-MM-DD')

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
            const result = await getDayTasks(selectedDay.value.format('DD/MM/YYYY'), user.id)
            dayData.value = result.data
            return
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
        fetchData(selectedOption.value.key)
    },
    { immediate: true }
)

watch(() => selectedOption.value, (newValue) => {
    if (newValue) {
        fetchData(newValue.key)
    }
}, { immediate: true })

watch(() => [calendarData.value, weekData.value, dayData.value], () => {
    const dataMap: Record<string, TaskResponse[]> = {
        month: calendarData.value,
        week: weekData.value,
        day: dayData.value
    }

    const activeData = dataMap[selectedOption.value.key]
    if (!activeData || selectedTaskInformation.value === null) return

    const updated = activeData.find(data => data.ID === selectedTaskInformation.value?.ID)
    if (updated) {
        selectedTaskInformation.value = updated
    }
}, { immediate: true, deep: true })

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
                <div @click="handleChangeToToday"
                    class="cursor-pointer rounded-4xl border-1 border-[#a3a3a3] bg-[#1c1c1c] px-4 py-2">
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
        <div v-if="selectedOption?.key === 'month'" class="w-full backdrop-blur-sm"
            :style="{ transform: `scaleX(${calendarSize.scaleX}) scaleY(${calendarSize.scaleY})`, transformOrigin: 'top' }">
            <n-config-provider :theme="darkTheme">
                <n-calendar :value="selectedDay.valueOf()" @update:value="openTaskModal"
                    @panel-change="handlePanelChange">
                    <template #="{ year, month, date }">
                        <div class="max-h-[100px] overflow-y-auto custom-calendar-scroll">
                            <vue-draggable class="min-h-[90px]"
                                :model-value="getTasksForCalendarDate(year, month, date)" :sort="false" :animation="150"
                                :group="{
                                    name: 'tasks',
                                    pull: true,
                                    put: ['tasks', 'active', 'pending', 'past']
                                }" @add="(event) => handleUpdateTaskDueDate(event, year, month, date)"
                                :onMove="(event) => checkIfUpdateStartDatePossible(event)"
                                :data-day="`${String(date).padStart(2, '0')}/${String(month).padStart(2, '0')}/${year}`">
                                <div v-for="task in getTasksForCalendarDate(year, month, date)" :key="task.ID"
                                    :data-task-id="task.ID"
                                    class=" relative z-9999 border-1 border-[#3a3a3a] rounded-lg mb-2 max-w-[130px] flex items-center justify-start  p-2 "
                                    @click.stop="openTaskInformationModal(task)">
                                    <n-checkbox @click.stop :checked="task.completed"
                                        @update-checked="(value: boolean) => updateTaskCheckbox(task.ID, value)" />
                                    <p class="font-jakarta ml-1 truncate w-32">{{ task.title }}</p>
                                </div>
                            </vue-draggable>
                        </div>
                    </template>
                </n-calendar>
            </n-config-provider>
        </div>
        <div v-if="selectedOption?.key === 'week'"
            class="w-full max-h-[75vh] overflow-y-auto bg-[#1c1c1c] backdrop-blur-sm custom-scroll">
            <table class="w-full table-fixed">
                <thead>
                    <tr class="sticky top-0 z-10 bg-[#1c1c1c]">
                        <th class="w-16 p-4"></th>
                        <th class="p-4" v-for="day in weekDate" :key="day">
                            <div :class="[
                                'flex flex-col items-center justify-center w-16 h-16 mx-auto text-white',
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
                            </p>c
                        </td>
                        <td v-for="day in weekDate" :key="day" class="border border-gray-600 h-25 max-h-25">
                            <div v-for="times in timeVariations"
                                class="min-h-[25px] max-h-[25px] overflow-x-auto custom-scroll">
                                <vue-draggable @click="openTaskModalWeekAndDate(day, time, times)"
                                    @add="(event) => handleUpdateTaskDueDateV2(event, dayjs(day + ' ' + time, 'DD/MM/YYYY HH:mm').add(times, 'minutes').format('DD/MM/YYYY HH:mm'))"
                                    :onMove="(event) => checkIfUpdateStartDateWeekPossible(event)"
                                    class="min-h-[25px] max-h-[25px] flex items-center overflow-x-auto whitespace-nowrap custom-scroll"
                                    :data-day-time="dayjs(day + ' ' + time, 'DD/MM/YYYY HH:mm').add(times, 'minutes').format('DD/MM/YYYY HH:mm')"
                                    :model-value="getTasksForCalendarWeek(day, time, times)" :sort="false"
                                    :swap-threshold="0.65" :invert-swap="true" :animation="150" :group="{
                                        name: 'tasks',
                                        pull: true,
                                        put: ['tasks', 'active', 'pending', 'past']
                                    }">
                                    <div @click.stop="openTaskInformationModal(task)"
                                        class="max-h-[22px] min-h-[22px] max-w-[100px] flex items-center justify-start px-2 py-[2px] box-border rounded-sm border-1 border-[#a3a3a3] ml-2 cursor-pointer"
                                        :data-task-id="task.ID"
                                        v-for="task in getTasksForCalendarWeek(day, time, times)">
                                        <n-config-provider :theme="darkTheme">
                                            <n-checkbox :checked="task.completed" @click.stop
                                                @update-checked="(value: boolean) => updateTaskCheckbox(task.ID, value)" />
                                        </n-config-provider>
                                        <p class="font-jakarta text-white text-xs w-20 truncate ml-2">{{ task.title }}
                                        </p>
                                    </div>
                                </vue-draggable>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div v-if="selectedOption?.key === 'day'"
            class="w-full max-h-[75vh] overflow-y-auto bg-[#1c1c1c] backdrop-blur-sm custom-scroll">
            <div class="w-full flex justify-center items-center sticky top-0 z-10 bg-[#1c1c1c]">
                <div :class="['w-16 h-16 mx-auto flex flex-col items-center justify-center',
                    dayjs(selectedDay).isSame(dayjs(), 'day') ? 'bg-blue-500 rounded-full' : '']">
                    <p class="text-white font-bold">{{ dayjs(selectedDay).format('ddd') }}</p>
                    <p class="text-white font-bold text-xl">{{ dayjs(selectedDay).format('DD') }}</p>
                </div>
            </div>
            <table class="w-full table-fixed">
                <thead>
                    <tr class="sticky top-0 z-10 bg-[#1c1c1c]">
                        <th class="w-16 p-4">

                        </th>
                    </tr>
                </thead>
                <tbody>
                    <!-- row to give space for time labels -->
                    <tr>
                        <td class="h-4"></td>
                    </tr>
                    <tr v-for="time in timeDate" :key="time">
                        <td class="relative">
                            <p class="absolute text-white text-sm -top-3">
                                {{ dayjs(time, 'HH:mm').format("h:mm A") }}
                            </p>
                        </td>
                        <td class="border-t border-gray-600 h-25 max-h-25">
                            <div v-for="times in timeVariations"
                                class="min-h-[25px] max-h-[25px] overflow-x-auto custom-scroll">
                                <vue-draggable @click="openTaskModalWeekAndDate(selectedDay.format('DD/MM/YYYY'), time, times)"
                                    @add="(event) => handleUpdateTaskDueDateV2(event, dayjs(selectedDay.format('DD/MM/YYYY') + ' ' + time, 'DD/MM/YYYY HH:mm').add(times, 'minutes').format('DD/MM/YYYY HH:mm'))"
                                    :onMove="(event) => checkIfUpdateStartDateDayPossible(event)"
                                    :model-value="getTasksForCalendarDay(time, times)"
                                    class="min-h-[25px] max-h-[25px] flex items-center overflow-x-auto whitespace-nowrap custom-scroll"
                                    :data-date-time="dayjs(selectedDay.format('DD/MM/YYYY') + ' ' + time, 'DD/MM/YYYY HH:mm').add(times, 'minutes').format('DD/MM/YYYY HH:mm')"
                                    :animation="150" :sort="false" :swap-threshold="0.65" :invert-swap="true" :group="{
                                        name: 'tasks',
                                        pull: true,
                                        put: ['tasks', 'active', 'pending', 'past']
                                    }">
                                    <div @click.stop="openTaskInformationModal(task)"
                                        class="max-h-[22px] min-h-[22px] max-w-[100px] flex items-center justify-start px-2 py-[2px] box-border rounded-sm border-1 border-[#a3a3a3] ml-2 cursor-pointer"
                                        :data-task-id="task.ID" v-for="task in getTasksForCalendarDay(time, times)">
                                        <n-config-provider :theme="darkTheme">
                                            <n-checkbox :checked="task.completed" @click.stop
                                                @update-checked="(value: boolean) => updateTaskCheckbox(task.ID, value)" />
                                        </n-config-provider>
                                        <p class="font-jakarta text-white text-xs w-20 truncate ml-2">{{ task.title }}
                                        </p>
                                    </div>
                                </vue-draggable>
                            </div>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div v-if="selectedOption?.key === 'agenda'">

        </div>
    </div>
    <MyCalendarTaskDetail :open="openCalendarTaskDetail" :type="typeDatas" @close-modal="closeTaskModal"
        :dateData="dateDatas" @result-message="handleResultMessage" />
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

:deep(.n-calendar-dates) {
    max-height: 1150px !important;
    overflow-y: auto !important;
}

:deep(.n-calendar-dates)::-webkit-scrollbar {
    width: 6px;
}

:deep(.n-calendar-dates)::-webkit-scrollbar-track {
    background: transparent;
}

:deep(.n-calendar-dates)::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 9999px;
}

:deep(.n-calendar-dates)::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
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