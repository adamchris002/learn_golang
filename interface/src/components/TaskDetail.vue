<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { type TaskResponse } from '@/services/taskServices';

import addIcon from '@/assets/icons/add.svg'
import closeIcon from '@/assets/icons/close.svg'
import { NModal, NCard, NButton, NIcon, NDatePicker, darkTheme, NConfigProvider, NInput, NCheckbox } from 'naive-ui';
import { sanitizeInput } from '@/composable/sanitizeInput';

const props = defineProps<{ open: boolean, task: TaskResponse | null }>()
const emits = defineEmits(["closeModal", "updateTask", "deleteTask", "updateTaskDatas", "deleteSubtask"])

type subtask = {
    id: number;
    title: string;
    completed: boolean;
}

const someAreaMissing = ref<boolean>(false)

const taskDueDate = ref<string | null>(null)
const taskDescription = ref<string>("")
const substaskArray = ref<subtask[] | null>([])

const originalTask = ref<{
    dueDate: string | null
    description: string
    subtasks: subtask[]
} | null>(null)

const hasChanges = computed(() => {
    if (!originalTask.value) return false

    if (taskDueDate.value !== originalTask.value.dueDate)
        return true

    if (taskDescription.value !== originalTask.value.description)
        return true

    const current = substaskArray.value ?? []
    const original = originalTask.value.subtasks

    if (current.length > original.length)
        return true

    for (const originalSubtask of original) {
        const currentSubtask = current.find(s => s.id === originalSubtask.id)

        if (!currentSubtask)
            continue

        if (
            currentSubtask.title !== originalSubtask.title ||
            currentSubtask.completed !== originalSubtask.completed
        ) {
            return true
        }
    }

    return false
})

function disablePreviousDate(ts: number) {
    const today = new Date()
    today.setHours(0, 0, 0, 0)

    return ts < today.getTime()
}

function deleteSubtasks(subTaskId: number, index: number) {
    if (subTaskId === 0) {
        substaskArray.value?.splice(index, 1)
    } else {
        if (props.task !== null) {
            const taskId = props.task.ID
            emits("deleteSubtask", subTaskId, taskId)
        }
    }
}

function addSubtask() {
    substaskArray.value?.push({ id: 0, title: "", completed: false })
}

function sendDataForUpdate() {
    if (!taskDescription.value || !taskDueDate.value || (substaskArray.value && substaskArray.value?.length > 0 && substaskArray.value.some((item) => item.title === ""))) {
        someAreaMissing.value = true
    } else {
        emits('updateTaskDatas', taskDueDate.value, taskDescription.value, substaskArray.value, props.task?.ID)
    }
}

function sanitizeDescription() {
    taskDescription.value = sanitizeInput(taskDescription.value)
}

function clearModalValues() {
    taskDueDate.value = ""
    taskDescription.value = ""
    substaskArray.value = []
    someAreaMissing.value = false
}
watch(
    () => props.task,
    (newValue) => {
        taskDueDate.value = newValue?.due_date?.slice(0, 10) || null
        taskDescription.value = newValue?.description ?? ""

        substaskArray.value = (newValue?.subtasks ?? []).map(subtask => ({
            id: subtask.ID,
            title: subtask.title,
            completed: subtask.completed,
        }))

        originalTask.value = {
            dueDate: taskDueDate.value,
            description: taskDescription.value,
            subtasks: (substaskArray.value ?? []).map(subtask => ({
                id: subtask.id,
                title: subtask.title,
                completed: subtask.completed
            }))
        }
    },
    { immediate: true, deep: true }
)

</script>
<template>
    <n-modal v-model:show="props.open" @after-leave="clearModalValues">
        <n-card
            :theme-overrides="{ colorModal: '#0d0d0d', titleTextColor: 'white', titleFontSizeHuge: '32px', titleFontWeight: '500' }"
            style="width: 600px" :title="task?.title" :bordered="false" size="huge" role="dialog" aria-modal="true">
            <template #header-extra>
                <n-button circle ghost type="info" @click="emits('closeModal')">
                    <n-icon>
                        <closeIcon />
                    </n-icon>
                </n-button>
            </template>
            <div class="flex items-center">
                <p class="text-white font-jakarta text-base mr-4">Task due for: </p>
                <n-config-provider :theme="darkTheme" :theme-overrides="{
                    DatePicker: {
                        itemTextColor: '#fff'
                    }
                }">
                    <n-date-picker :is-date-disabled="disablePreviousDate" v-model:formatted-value="taskDueDate"
                        value-format="dd/MM/yyyy"
                        :status="(taskDueDate?.length === 0 || taskDueDate === null) && someAreaMissing ? 'error' : 'success'" />
                </n-config-provider>
            </div>
            <div class="mt-4">
                <p class="text-white font-jakarta text-base mb-2">Description: </p>
                <n-config-provider :theme="darkTheme" :theme-overrides="{
                    Input: {
                        itemTextColor: '#fff'
                    }
                }">
                    <n-input type="textarea" v-model:value="taskDescription" @blur="sanitizeDescription"
                        :status="taskDescription.length === 0 && someAreaMissing ? 'error' : 'success'" />
                </n-config-provider>
            </div>
            <div class="flex justify-between mt-8">
                <p class="text-white font-jakarta text-sm mb-2">Add a subtask: </p>
                <n-button circle ghost type="info" @click="addSubtask">
                    <n-icon>
                        <addIcon />
                    </n-icon>
                </n-button>
            </div>
            <div class="subtask-list">
                <div v-if="substaskArray && substaskArray?.length > 0" v-for="(items, index) in substaskArray"
                    class="mt-2 overflow-x-auto">
                    <n-config-provider :theme="darkTheme" :theme-overrides="{
                        Input: {
                            itemTextColor: '#fff'
                        }
                    }">
                        <div class="flex justify-start items-center">
                            <n-checkbox v-model:checked="items.completed" class="mr-4" />
                            <n-input v-model:value="items.title" class="mr-4"
                                :status="items.title.length === 0 && someAreaMissing ? 'error' : 'success'" />
                            <n-button circle ghost type="error" class="!mr-4" @click="deleteSubtasks(items.id, index)">
                                <n-icon>
                                    <closeIcon />
                                </n-icon>
                            </n-button>
                        </div>
                    </n-config-provider>
                </div>
            </div>
            <div class="flex justify-end class mt-8">
                <n-button @click="emits('deleteTask', task?.ID)" type="error" class="!mr-4">Delete Task</n-button>
                <n-button :disabled="!hasChanges" type="info" @click="sendDataForUpdate">Update Task</n-button>
            </div>
        </n-card>
    </n-modal>
</template>
<style scoped>
.subtask-list {
    max-height: 300px;
    overflow-x: auto;
    scrollbar-width: thin;
    scrollbar-color: #4b5563 transparent;
}

.subtask-list::-webkit-scrollbar {
    width: 6px;
}

.subtask-list::-webkit-scrollbar-track {
    background: transparent;
}

.subtask-list::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 9999px;
}

.subtask-list::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
}
</style>