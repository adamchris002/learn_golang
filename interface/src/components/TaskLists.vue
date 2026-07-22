<script setup lang="ts">
import { type TaskResponse } from '@/services/taskServices';
import dayjs from 'dayjs';

import { NCheckbox } from 'naive-ui';

const props = defineProps<{ title: string, taskArray: TaskResponse[], ifEmptyString: string }>()
const emits = defineEmits<{
    (e: "toggle", id: number, completed: boolean): void;
    (e: "openTaskDetail", item: TaskResponse): void;
}>()
</script>
<template>
    <div>
        <div class="w-200 mt-4 flex justify-start">
            <p class="font-jakarta text-2xl text-white font-semibold">{{props.title}}</p>
        </div>
        <div v-if="taskArray.length > 0 && taskArray" v-for="items in taskArray">
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
        <div class="flex justify-center mt-8" v-else>
            <p class="font-jakarta text-xl text-white font-light">No tasks to complete</p>
        </div>
    </div>
</template>
<style scoped></style>