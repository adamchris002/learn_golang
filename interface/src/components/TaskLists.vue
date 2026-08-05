<script setup lang="ts">
import { ref, watch } from 'vue';
import { type TaskResponse } from '@/services/taskServices';
import dayjs from 'dayjs';

import { NCheckbox } from 'naive-ui';
import { VueDraggable, type DraggableEvent } from 'vue-draggable-plus';

const props = defineProps<{ taskArray: TaskResponse[], ifEmptyString: string }>()
const emits = defineEmits<{
    (e: "toggle", id: number, completed: boolean): void;
    (e: "openTaskDetail", item: TaskResponse): void;
}>()

const completedTasks = ref<TaskResponse[]>([])
const incompleteTasks = ref<TaskResponse[]>([])

function onCompleteChange(event: DraggableEvent<TaskResponse>) {
    emits('toggle', event.data.ID, true)
}

function onIncompleteChange(event: DraggableEvent<TaskResponse>) {
    emits('toggle', event.data.ID, false)
}

watch(() => props.taskArray, (newValue) => {
    if (newValue.length > 0) {
        completedTasks.value = newValue.filter(data => data.completed)
        incompleteTasks.value = newValue.filter(data => !data.completed)
    }
}, { immediate: true, deep: true })
</script>
<template>
    <div>
        <div v-if="taskArray.length > 0 && taskArray" class="max-h-[55vh] overflow-y-auto custom-scroll">
            <div class="w-200 mt-4 flex justify-start">
                <p class="font-jakarta text-2xl text-white font-semibold">Completed Tasks</p>
            </div>
            <div v-if="completedTasks.length <= 0" class="flex justify-center mt-4">
                <p class="font-jakarta text-white text-base">No completed Tasks</p>
            </div>
            <vue-draggable v-model="completedTasks" @add="onCompleteChange" ghostClass="ghost" :animation="150"
                :group="{ name: 'tasks' }" :sort="false">
                <div v-for="items in completedTasks" :key="items.ID">
                    <div
                        class="mt-4 rounded-lg w-200 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-2 flex items-center">
                        <n-checkbox :checked="items.completed"
                            @update-checked="(value: boolean) => emits('toggle', items.ID, value)" />
                        <div @click="emits('openTaskDetail', items)" class="w-full cursor-pointer">
                            <p class="font-jakarta text-[#8a8888] text-base ml-4">Created at: {{
                                dayjs(items.CreatedAt).format("DD/MMMM/YYYY HH:mm:ss") }}</p>
                            <p class="font-jakarta text-white text-xl ml-4">{{ items.title }}</p>
                        </div>
                    </div>
                </div>

            </vue-draggable>
            <div class="w-200 mt-4 flex justify-start">
                <p class="font-jakarta text-2xl text-white font-semibold">Incomplete Tasks</p>
            </div>
            <div v-if="incompleteTasks.length <= 0" class="flex justify-center mt-4">
                <p class="font-jakarta text-white text-base">No Tasks to complete</p>
            </div>
            <vue-draggable v-model="incompleteTasks" @add="onIncompleteChange" ghostClass="ghost" :animation="150"
                :group="{ name: 'tasks' }" :sort="false">

                <div v-for="items in incompleteTasks" :key="items.ID">
                    <div
                        class="mt-4 rounded-lg w-200 backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-2 flex items-center">
                        <n-checkbox :checked="items.completed"
                            @update-checked="(value: boolean) => emits('toggle', items.ID, value)" />
                        <div @click="emits('openTaskDetail', items)" class="w-full cursor-pointer">
                            <p class="font-jakarta text-[#8a8888] text-base ml-4">Created at: {{
                                dayjs(items.CreatedAt).format("DD/MMMM/YYYY HH:mm:ss") }}</p>
                            <p class="font-jakarta text-white text-xl ml-4">{{ items.title }}</p>
                        </div>
                    </div>
                </div>
            </vue-draggable>
        </div>
        <div class="flex justify-center mt-8" v-else>
            <p class="font-jakarta text-xl text-white font-light">No tasks to complete</p>
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
</style>