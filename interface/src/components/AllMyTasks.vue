<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { changeTaskStartDate, type TaskResponse } from "@/services/taskServices";
import dayjs from "dayjs";
import customParseFormat from 'dayjs/plugin/customParseFormat'
import isSameOrBefore from 'dayjs/plugin/isSameOrBefore'
import isSameOrAfter from 'dayjs/plugin/isSameOrAfter'

import { NIcon, NCollapse, NAlert, NDatePicker, NInput, NButton, NConfigProvider, darkTheme, NCheckbox } from "naive-ui"
import pageIcon from "@/assets/icons/pages.svg"
import addIcon from "@/assets/icons/add.svg"
import closeIcon from "@/assets/icons/close.svg"
import AllTasksItems from "./AllTasksItems.vue";
import { useItemScale } from "@/composable/pageAdjuster.ts";

dayjs.extend(customParseFormat)
dayjs.extend(isSameOrBefore)
dayjs.extend(isSameOrAfter)

const user = JSON.parse(localStorage.getItem("user") || "{}")
const scale = useItemScale()

const props = defineProps<{ allTasksArray: TaskResponse[] }>()
const emits = defineEmits(['requestCallAllTasks'])

const todaysTasks = ref<TaskResponse[]>([])
const tomorrowsTasks = ref<TaskResponse[]>([])
const pendingTasks = ref<TaskResponse[]>([])
const pastTasks = ref<TaskResponse[]>([])

const selectedTask = ref<TaskResponse | null>(null)

const expandedNames = ref<string[]>([])
const taskErrorMessage = ref<{ status: number; messageTitle: string; message: string } | null>(null)
const hasChanges = computed(() => {
    const selected = selectedTask.value
    if (!selected) return false

    const original = props.allTasksArray.find((data) => data.ID === selected.ID)
    if (!original) return false

    if (selected.due_date !== original.due_date) return true
    if (selected.description !== original.description) return true

    const selectedSubtasks = selected.subtasks ?? []
    const originalSubtasks = original.subtasks ?? []

    if (selectedSubtasks.length !== originalSubtasks.length) return true

    return selectedSubtasks.some(current => {
        const originalSub = originalSubtasks.find(s => s.ID === current.ID)
        if (!originalSub) return true // new subtask with an ID not in original (edge case)
        return current.title !== originalSub.title || current.completed !== originalSub.completed
    })
})

const formattedDueDate = computed({
    get: () => selectedTask.value?.due_date || null,
    set: (value: string | null) => {
        if (selectedTask.value) {
            selectedTask.value.due_date = value ?? ''
        }
    }
})

function handleSetTaskErrorMessage(data: { status: number; messageTitle: string; message: string }) {
    taskErrorMessage.value = data
}

function handleSetSelectedTask(data: TaskResponse) {
    selectedTask.value = JSON.parse(JSON.stringify(data))
}

function handleShouldExpand(name: string, shouldExpand: boolean) {
    const isExpanded = expandedNames.value.includes(name)
    if (shouldExpand && !isExpanded) {
        expandedNames.value.push(name)
    }
}

function disablePreviousDate(ts: number) {
    const today = new Date()
    today.setHours(0, 0, 0, 0)

    return ts < today.getTime()
}

async function handleUpdateTask(data: TaskResponse, name: string) {
    switch (name) {
        case 'today':
            const today = dayjs().startOf('day').format("DD/MM/YYYY")
            taskErrorMessage.value = await changeTaskStartDate(data.ID, user.id, today)
            if (taskErrorMessage.value.status === 200) {
                emits('requestCallAllTasks')
            }
            break
        case 'tomorrow':
            const tomorrow = dayjs().startOf('day').add(1, 'day').format("DD/MM/YYYY")
            taskErrorMessage.value = await changeTaskStartDate(data.ID, user.id, tomorrow)
            if (taskErrorMessage.value.status === 200) {
                emits('requestCallAllTasks')
            }
            break
        default: {
            console.error(`Unknown task category: ${name}`)
        }
    }
}

watch(() => props.allTasksArray, (newValue) => {
    if (newValue.length > 0) {
        const today = dayjs().startOf('day')
        const tomorrow = dayjs().add(1, "day").startOf('day')

        todaysTasks.value = newValue.filter((data) => {
            const taskStart = dayjs(data.task_start, "DD/MM/YYYY").startOf('day')
            const hasDueDate = !!data.due_date

            if (taskStart.isSame(today) && !hasDueDate) return true

            if (taskStart.isSameOrBefore(today) && hasDueDate) {
                const dueDate = dayjs(data.due_date, "DD/MM/YYYY").startOf('day')
                if (dueDate.isSameOrAfter(today)) return true
            }

            return false
        })

        tomorrowsTasks.value = newValue.filter((data) => {
            const taskStart = dayjs(data.task_start, "DD/MM/YYYY").startOf('day')
            const hasDueDate = !!data.due_date

            if (taskStart.isSame(tomorrow) && !hasDueDate) return true

            if (taskStart.isSameOrBefore(tomorrow) && hasDueDate) {
                const dueDate = dayjs(data.due_date, "DD/MM/YYYY").startOf('day')
                if (dueDate.isSameOrAfter(tomorrow)) return true
            }

            return false
        })
        pendingTasks.value = newValue.filter((data) => {
            const taskStart = dayjs(data.task_start, "DD/MM/YYYY").startOf('day')

            if (!taskStart.isAfter(today)) return false

            if (!data.due_date) return true
            const dueDate = dayjs(data.due_date, "DD/MM/YYYY").startOf('day')
            return dueDate.isAfter(today)
        })

        pastTasks.value = newValue.filter((data) => {
            const taskStart = dayjs(data.task_start, "DD/MM/YYYY").startOf('day')

            if (!taskStart.isBefore(today)) return false

            if (!data.due_date) return true
            const dueDate = dayjs(data.due_date, "DD/MM/YYYY").startOf('day')
            return dueDate.isBefore(today)
        })
    }
}, { immediate: true, deep: true })
</script>
<template>
    <div class="relative z-10 w-full px-4 py-8 h-screen scale-container"
        :style="{ transform: `scale(${scale.scale})`, transformOrigin: 'top left', zoom: scale.zoom }">
        <div
            class="flex items-center backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] size-fit rounded-xl py-2 px-4">
            <n-icon size="16" class="mr-4">
                <pageIcon class="text-white " />
            </n-icon>
            <p class="font-jakarta text-white text-lg">All My Tasks</p>
        </div>
        <div class="flex items-center justify-start mt-8">
            <div class="min-h-[80vh] min-w-[30vw] rounded-lg backdrop-blur-sm bg-[#171717]/65 p-8">
                <div class="max-h-[80vh] min-h-[80vh] overflow-x-auto custom-scroll">
                    <n-collapse v-model:expanded-names="expandedNames"
                        :theme-overrides="{ titleTextColor: 'white', titleFontSize: '20px' }">
                        <AllTasksItems name="today" :title="'Today'" :tasksArray="todaysTasks"
                            @shouldExpand="handleShouldExpand" @update-task="handleUpdateTask"
                            @set-selected-task="handleSetSelectedTask" @send-error-message="handleSetTaskErrorMessage"/>
                        <AllTasksItems name="tomorrow" :title="'Tomorrow'" :tasksArray="tomorrowsTasks"
                            @shouldExpand="handleShouldExpand" @update-task="handleUpdateTask"
                            @set-selected-task="handleSetSelectedTask" />
                        <AllTasksItems name="pending" :title="'Pending'" :tasksArray="pendingTasks"
                            @shouldExpand="handleShouldExpand" @set-selected-task="handleSetSelectedTask" />
                        <AllTasksItems name="past" :title="'Past'" :tasksArray="pastTasks"
                            @shouldExpand="handleShouldExpand" @set-selected-task="handleSetSelectedTask" />
                    </n-collapse>

                </div>
            </div>
            <div class="min-h-[80vh] min-w-[30vw] rounded-lg backdrop-blur-sm bg-[#171717]/65 p-8 ml-8">
                <div v-if="selectedTask === null" class="min-h-[80vh] flex items-center justify-center">
                    <div class="mb-[120px]">
                        <p class="text-white font-jakarta text-4xl">Hey There!</p>
                        <p class="text-white font-jakarta text-xl">Please select a Task to show detail</p>
                    </div>
                </div>
                <div v-else
                    class="max-h-[80vh] min-h-[80vh] overflow-x-auto custom-scroll flex flex-col justify-between">
                    <div>
                        <p class="text-white font-jakarta text-3xl truncate w-80">{{ selectedTask.title }}</p>

                        <div class="flex items-center justify-start mt-8">
                            <p class="font-jakarta text-white text-base">Task due for: </p>
                            <n-config-provider :theme="darkTheme" :theme-overrides="{
                                DatePicker: {
                                    itemTextColor: '#fff'
                                }
                            }">
                                <n-date-picker class="ml-2" :is-date-disabled="disablePreviousDate"
                                    v-model:formatted-value="formattedDueDate"
                                    :input-props="{ id: 'my-custom-datepicker-id', name: 'myCustomDatePickerId' }"
                                    value-format="dd/MM/yyyy" format="dd/MM/yyyy" />
                            </n-config-provider>
                        </div>

                        <div class="mt-4">
                            <p class="font-jakarta text-white text-base">Description:</p>
                            <n-config-provider :theme="darkTheme" :theme-overrides="{
                                Input: {
                                    itemTextColor: '#fff'
                                }
                            }">
                                <n-input :input-props="{
                                    id: 'edit-task-description',
                                    name: 'editTaskDescription'
                                }" v-model:value="selectedTask.description" type="textarea" class="mt-2" />
                            </n-config-provider>
                        </div>

                        <div class="flex items-center justify-between mt-4">
                            <p class="font-jakarta text-white text-base">Add a subtask:</p>
                            <n-button circle ghost type="info">
                                <n-icon>
                                    <addIcon />
                                </n-icon>
                            </n-button>
                        </div>

                        <div class="subtask-list">
                            <div v-if="selectedTask.subtasks && selectedTask.subtasks?.length > 0"
                                v-for="(items, index) in selectedTask.subtasks" class="mt-2">
                                <n-config-provider :theme="darkTheme" :theme-overrides="{
                                    Input: {
                                        itemTextColor: '#fff'
                                    }
                                }">
                                    <div class="flex justify-start items-center">
                                        <n-checkbox :checked="items.completed" class="mr-4" />
                                        <n-input :input-props="{
                                            id: `subtask-title-${index}`,
                                            name: 'subtaskTitle'
                                        }" v-model:value="items.title" class="mr-4" />

                                        <n-button circle ghost type="error" class="!mr-4">
                                            <n-icon>
                                                <closeIcon />
                                            </n-icon>
                                        </n-button>
                                    </div>
                                </n-config-provider>
                            </div>
                            <div class="flex items-center justify-center" v-else>
                                <p class="text-white font-jakarta text-sm">No subtasks for this task</p>
                            </div>
                        </div>
                    </div>
                    <div class="flex justify-end">
                        <n-button type="error">Delete Task</n-button>
                        <n-button :disabled="!hasChanges" type="info" class="!ml-2">Update Task</n-button>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div class="absolute top-2 right-2 z-9999">
        <n-alert v-if="taskErrorMessage" :title="taskErrorMessage.messageTitle"
            :type="taskErrorMessage.status === 200 ? 'success' : 'error'" closable @close="taskErrorMessage = null">
            {{ taskErrorMessage.message }}
        </n-alert>
    </div>
</template>
<style scoped>
.subtask-list {
    max-height: 300px;
    overflow-x: auto;
    scrollbar-width: thin;
    scrollbar-color: #4b5563 transparent;
}

.scale-container {
    transition: transform 300ms ease-in-out;
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

:deep(.n-collapse-item-arrow) {
    display: none !important;
}

:deep(.n-collapse-item) {
    border: none !important;
}
</style>