<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';

import { NConfigProvider, NModal, NCard, NIcon, NButton, NDatePicker, NInput, darkTheme, NCheckbox } from 'naive-ui'
import closeIcon from '@/assets/icons/close.svg'
import addIcon from '@/assets/icons/add.svg'
import { sanitizeInput } from '@/composable/sanitizeInput';
import { postTodaysTaskV2 } from '@/services/taskServices';

type subtask = {
    id: number;
    title: string;
    completed: boolean;
}

const user = JSON.parse(localStorage.getItem("user") || "{}")

const props = defineProps<{ open: boolean, dateData: string }>()
const emits = defineEmits(["closeModal", 'handlePostNewData', "resultMessage"])

const cardTitle = ref<string>('')

//datas
const taskTitle = ref<string>('')
const taskDueDate = ref<string | null>(null)
const taskDescription = ref<string>('')

const subtaskArray = ref<subtask[] | null>([])

const someAreaMissing = computed(() => {
    return Boolean(
        taskTitle.value &&
        taskDescription.value &&
        taskDueDate.value &&
        !subtaskArray.value?.some(item => !item.title)
    )
})

function deleteSubtasks(index: number) {
    subtaskArray.value?.splice(index, 1)
}

function handleSanitizeTitle() {
    taskTitle.value = sanitizeInput(taskTitle.value)
}

function handleSanitizeDescription() {
    taskDescription.value = sanitizeInput(taskDescription.value)
}

function sanitizeSubtaskTitle(index: number) {
    if (!subtaskArray.value && subtaskArray.value === null) {
        return
    }
    const subTasksData = subtaskArray.value[index]
    if (!subTasksData) {
        return
    }
    subTasksData.title = sanitizeInput(subTasksData.title)
}

function disablePreviousDate(ts: number) {
    const today = new Date()
    today.setHours(0, 0, 0, 0)

    return ts < today.getTime()
}

function handleCloseModal() {
    cardTitle.value = ''
    taskTitle.value = ''
    taskDescription.value = ''
    taskDueDate.value = null
    emits('closeModal')
}

function handleAddSubtask() {
    subtaskArray.value?.push({ id: 0, title: '', completed: false })
}

async function handleCreateNewData() {
    if (!someAreaMissing.value) {
        return
    }

    if (!taskDueDate.value || !subtaskArray.value) {
        return
    }

    const data = {
        title: taskTitle.value,
        description: taskDescription.value,
        task_start: props.dateData,
        due_date: taskDueDate.value,
        completed: false,
        subtasks: subtaskArray.value,
        userId: user.id
    }
    // emits('handlePostNewData', data)
    const result = await postTodaysTaskV2(data)
    if (result.status === 200) {
        handleCloseModal()
    }
    emits('resultMessage', result)
}

watch(() => props.dateData, (newData) => {
    if (newData.length > 0) {
        cardTitle.value = `New Task ${newData}`;
    }
}, { immediate: true })

onMounted(() => {
    subtaskArray.value?.push({ id: 0, title: '', completed: false })
})

onUnmounted(() => {
})

</script>
<template>
    <n-modal v-model:show="props.open">
        <n-card
            :theme-overrides="{ colorModal: '#0d0d0d', titleTextColor: 'white', titleFontSizeHuge: '32px', titleFontWeight: '500' }"
            style="width: 600px" :title="cardTitle" :bordered="false" size="huge" role="dialog" aria-modal="true">
            <template #header-extra>
                <n-button circle ghost type="info" @click="handleCloseModal">
                    <n-icon>
                        <closeIcon />
                    </n-icon>
                </n-button>
            </template>
            <div class="mb-3 flex items-center justify-start">
                <p class="text-white font-jakarta text-base">Task Name: </p>
                <n-config-provider :theme="darkTheme" :theme-overrides="{
                    Input: {
                        itemTextColor: '#fff'
                    }
                }">
                    <n-input class="ml-2" v-model:value="taskTitle" @blur="handleSanitizeTitle" />
                </n-config-provider>
            </div>
            <div class="mb-3 flex items-center justify-start">
                <p class="text-white font-jakarta text-base">Task due for: </p>
                <n-config-provider :theme="darkTheme" :theme-overrides="{
                    DatePicker: {
                        itemTextColor: '#fff'
                    }
                }">
                    <n-date-picker v-model:formatted-value="taskDueDate" value-format="dd/MM/yyyy" format="dd/MM/yyyy"
                        :is-date-disabled="disablePreviousDate" class="ml-2" />
                </n-config-provider>
            </div>
            <div class="mb-2">
                <p class="text-white font-jakarta text-base mb-2">Description: </p>
                <n-config-provider :theme="darkTheme" :theme-overrides="{
                    Input: {
                        itemTextColor: '#fff'
                    }
                }">
                    <n-input type="textarea" v-model:value="taskDescription" @blur="handleSanitizeDescription" />
                </n-config-provider>
            </div>
            <div class="mb-2">
                <div class="flex justify-between">
                    <p class="text-white font-jakarta text-base mb-2">Add a subtask: </p>
                    <n-button circle ghost type="info" @click="handleAddSubtask">
                        <n-icon>
                            <addIcon />
                        </n-icon>
                    </n-button>
                </div>
                <div class="subtask-list">
                    <div v-if="subtaskArray && subtaskArray?.length > 0" v-for="(items, index) in subtaskArray"
                        class="mt-2 overflow-x-auto">
                        <n-config-provider :theme="darkTheme" :theme-overrides="{
                            Input: {
                                itemTextColor: '#fff'
                            }
                        }">
                            <div class="flex justify-start items-center">
                                <n-checkbox v-model:checked="items.completed" class="mr-4" />
                                <n-input :input-props="{
                                    id: `subtask-title-${index}`,
                                    name: 'subtaskTitle'
                                }" v-model:value="items.title" class="mr-4" @blur="sanitizeSubtaskTitle(index)" />
                                <n-button circle ghost type="error" class="!mr-4" @click="deleteSubtasks(index)">
                                    <n-icon>
                                        <closeIcon />
                                    </n-icon>
                                </n-button>
                            </div>
                        </n-config-provider>
                    </div>
                    <div class="flex justify-end class mt-8">
                        <n-button @click="handleCreateNewData" :disabled="!someAreaMissing" type="info">Add
                            Task</n-button>
                    </div>
                </div>
            </div>
        </n-card>
    </n-modal>
</template>
<style scoped lang="css"></style>