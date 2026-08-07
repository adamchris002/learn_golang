<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { ref, watch } from "vue";
import { VueDraggable, type DraggableEvent } from "vue-draggable-plus";
import { NIcon, NButton } from "naive-ui"
import { changeTaskToActive, type TaskResponse } from "@/services/taskServices";

import chevronLeft from "@/assets/icons/chevron-left.svg"
import closeIcon from "@/assets/icons/close.svg"
import dayjs from "dayjs";
import customParseFormat from 'dayjs/plugin/customParseFormat'
import isSameOrBefore from 'dayjs/plugin/isSameOrBefore'
import isSameOrAfter from 'dayjs/plugin/isSameOrAfter'

dayjs.extend(customParseFormat)
dayjs.extend(isSameOrBefore)
dayjs.extend(isSameOrAfter)

const auth = useAuthStore()
const user = JSON.parse(localStorage.getItem("user") || "{}")

const props = defineProps<{ showDrawer: boolean, allIncompleteTasks: TaskResponse[], taskDrawerPinned: boolean }>()
const emits = defineEmits(['toggleDrawer', 'requestRefreshAllTasks'])

const activeTasks = ref<TaskResponse[]>([])
const pendingTasks = ref<TaskResponse[]>([])
const pastTasks = ref<TaskResponse[]>([])

const toggleDrawer = () => {
    emits('toggleDrawer', !props.showDrawer)
}

async function handleChangeToActiveTasks(event: DraggableEvent<TaskResponse>) {
    auth.message = await changeTaskToActive(event.data.ID, user.id, event.data.due_date)
    emits('requestRefreshAllTasks')
}

// async function handleChangeToPendingTasks(event: DraggableEvent<TaskResponse>) {
//     auth.message = await changeTaskToPending(event.data.ID, user.id)
//     emits('requestRefreshAllTasks')
// }

watch(() => props.allIncompleteTasks, (newValue) => {
    if (newValue.length === 0) {
        activeTasks.value = []
        pendingTasks.value = []
        pastTasks.value = []
        return
    }

    const today = dayjs().startOf('day')

    activeTasks.value = newValue.filter((data) => {
        const taskStart = dayjs(data.task_start, "DD/MM/YYYY").startOf('day')
        const hasDueDate = !!data.due_date

        if (taskStart.isSame(today) && !hasDueDate) return true

        if (taskStart.isSameOrBefore(today) && hasDueDate) {
            const dueDate = dayjs(data.due_date, "DD/MM/YYYY").startOf('day')
            if (dueDate.isSameOrAfter(today)) return true
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
}, { immediate: true, deep: true })
</script>
<template>
    <div v-show="taskDrawerPinned" class="fixed z-20 right-0 top-1/2 -translate-y-1/2 flex items-center">
        <Transition name="drawer">
            <div v-if="props.showDrawer">
                <div class="h-screen bg-black/40 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] p-4">
                    <div class="flex items-center justify-between mb-8">
                        <p class="text-white text-xl font-jakarta mr-4">All Incomplete Tasks</p>
                        <n-button ghost circle type="info" @click="toggleDrawer">
                            <n-icon>
                                <closeIcon class="text-white" />
                            </n-icon>
                        </n-button>
                    </div>

                    <p class="text-white text-base font-jakarta">Active Tasks ({{ activeTasks.length }})</p>
                    <div v-if="activeTasks.length <= 0" class="flex justify-center mt-4">
                        <p class="text-white font-jakarta">No Active Tasks</p>
                    </div>
                    <div class="max-h-[22vh] overflow-y-auto custom-scroll">
                        <vue-draggable v-model="activeTasks" :sort="false" :animation="150" :group="{
                            name: 'active',
                            pull: true,
                            put: ['pending', 'past']
                        }" @add="handleChangeToActiveTasks">
                            <div v-for="items in activeTasks" :key="items.ID" class="max-h-[100px] overflow-y-auto">
                                <div class="flex justify-start items-center mt-4">
                                    <div
                                        class="rounded-lg backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-1 cursor-grab">
                                        <p class="text-[#8a8888]">Created At: {{
                                            dayjs(items.CreatedAt).format("DD/MM/YYYY HH:mm") }}</p>
                                        <p class="text-white truncate w-45">{{ items.title }}</p>
                                    </div>
                                </div>
                            </div>
                        </vue-draggable>

                    </div>

                    <p class="text-white text-base font-jakarta mt-8">Pending Tasks ({{ pendingTasks.length }})</p>
                    <div v-if="pendingTasks.length <= 0" class="flex justify-center mt-4">
                        <p class="text-white font-jakarta">No Pending Tasks</p>
                    </div>
                    <div class="max-h-[22vh] overflow-y-auto custom-scroll">
                        <vue-draggable v-model="pendingTasks" :sort="false" :animation="150" :group="{
                            name: 'pending',
                            pull: true,
                            put: false
                        }">
                            <div v-for="items in pendingTasks" :key="items.ID" class="">
                                <div class="flex justify-start items-center mt-4">
                                    <div
                                        class="rounded-lg backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 cursor-grab">
                                        <p class="text-[#8a8888]">Created At: {{
                                            dayjs(items.CreatedAt).format("DD/MM/YYYY HH:mm") }}</p>
                                        <p class="text-white truncate w-45">{{ items.title }}</p>
                                    </div>
                                </div>
                            </div>
                        </vue-draggable>
                    </div>

                    <p class="text-white text-base font-jakarta mt-8">Past Tasks ({{ pastTasks.length }})</p>
                    <div v-if="pastTasks.length <= 0" class="flex justify-center mt-4">
                        <p class="text-white font-jakarta">No Past Tasks</p>
                    </div>
                    <div class="max-h-[22vh] overflow-y-auto custom-scroll">
                        <vue-draggable v-model="pastTasks" @add="" :sort="false" :animation="150" :group="{
                            name: 'past',
                            pull: true,
                            put: false
                        }">
                            <div v-for="items in pastTasks" :key="items.ID" class="">
                                <div class="flex justify-start items-center mt-4">
                                    <div
                                        class="rounded-lg backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 cursor-grab">
                                        <p class="text-[#8a8888]">Created At: {{
                                            dayjs(items.CreatedAt).format("DD/MM/YYYY HH:mm") }}</p>
                                        <p class="text-white truncate w-45">{{ items.title }}</p>
                                    </div>
                                </div>
                            </div>
                        </vue-draggable>
                    </div>
                </div>
            </div>
        </Transition>
        <div v-if="!props.showDrawer" @click="toggleDrawer"
            class="cursor-pointer rounded-l-lg p-2 bg-black/40 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)]">
            <NIcon :size="36">
                <chevronLeft class="text-white text-4xl" />
                <!-- <chevronRight v-else class="text-white text-4xl" /> -->
            </NIcon>
        </div>
    </div>
</template>
<style scoped>
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

.drawer-enter-active,
.drawer-leave-active {
    transition: transform 0.3s ease, opacity 0.3s ease;
    transform-origin: right;
    overflow: hidden;
}

.drawer-enter-from,
.drawer-leave-to {
    transform: translateX(100%);
    opacity: 0;
}

.drawer-enter-to,
.drawer-leave-from {
    transform: translateX(0);
    opacity: 1;
}
</style>