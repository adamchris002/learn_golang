<script setup lang="ts">
import { useAuthStore } from "@/stores/auth";
import { ref, watch } from "vue";
import { VueDraggable, type DraggableEvent } from "vue-draggable-plus";
import { NIcon, NCheckbox } from "naive-ui"
import { changeActiveTaskValue, type TaskResponse } from "@/services/taskServices";

import chevronLeft from "@/assets/icons/chevron-left.svg"
import chevronRight from "@/assets/icons/chevron-right.svg"
import dayjs from "dayjs";

const auth = useAuthStore()
const user = JSON.parse(localStorage.getItem("user") || "{}")

const props = defineProps<{ showDrawer: boolean, allIncompleteTasks: TaskResponse[] }>()
const emits = defineEmits(['toggleDrawer', 'requestCallAllIncompleteTasks'])

const activeTasks = ref<TaskResponse[]>([])
const pendingTasks = ref<TaskResponse[]>([])

const toggleDrawer = () => {
    emits('toggleDrawer', !props.showDrawer)
}

async function handleChangeActiveTasks(event: DraggableEvent<TaskResponse>) {
    auth.message = await changeActiveTaskValue(event.data.ID, user.id)
    emits('requestCallAllIncompleteTasks')
}

watch(() => props.allIncompleteTasks, (newValue) => {
    const today = dayjs().format("DD/MM/YYYY");
    console.log(today)
    if (newValue.length > 0) {
        activeTasks.value = newValue.filter(data => data.task_start === today || (data.due_date))
        pendingTasks.value = newValue.filter(data => data.task_start !== today && (!data.due_date))
    }
}, { immediate: true, deep: true })
</script>
<template>
    <div class="fixed right-0 top-1/2 -translate-y-1/2 flex items-center">
        <div @click="toggleDrawer"
            class="cursor-pointer rounded-l-lg p-2 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)]">
            <NIcon :size="36">
                <chevronLeft v-if="!props.showDrawer" class="text-white text-4xl" />
                <chevronRight v-else class="text-white text-4xl" />
            </NIcon>
        </div>
        <Transition name="drawer">
            <div v-if="props.showDrawer">
                <div class="h-screen backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] p-4">
                    <p class="text-white text-xl font-jakarta mb-8">All My Tasks</p>

                    <p class="text-white text-base font-jakarta">Active Tasks</p>
                    <div class="max-h-[35vh] overflow-y-auto custom-scroll">
                        <vue-draggable v-model="activeTasks" :sort="false" :animation="150" :group="{
                            name: 'tasks',
                            pull: true,
                            put: true
                        }">
                            <div v-for="items in activeTasks" :key="items.ID" class="max-h-[100px] overflow-y-auto">
                                <div class="flex justify-start items-center mt-4">
                                    <n-checkbox />
                                    <div
                                        class="rounded-lg backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 ml-2 cursor-grab">
                                        <p class="text-[#8a8888]">Created At: {{
                                            dayjs(items.CreatedAt).format("DD/MM/YYYY HH:mm") }}</p>
                                        <p class="text-white truncate w-45">{{ items.title }}</p>
                                    </div>
                                </div>
                            </div>
                        </vue-draggable>
                    </div>

                    <p class="text-white text-base font-jakarta mt-8">Pending Tasks</p>
                    <div class="max-h-[35vh] overflow-y-auto custom-scroll">
                        <vue-draggable v-model="pendingTasks" :sort="false" :animation="150" :group="{
                            name: 'tasks',
                            pull: false,
                            put: true
                        }" @add="handleChangeActiveTasks">
                            <div v-for="items in pendingTasks" :key="items.ID" class="">
                                <div class="flex justify-start items-center mt-4">
                                    <n-checkbox />
                                    <div
                                        class="rounded-lg backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 ml-2 cursor-grab">
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