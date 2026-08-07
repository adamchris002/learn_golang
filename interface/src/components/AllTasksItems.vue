<script setup lang="ts">
import { VueDraggable, type MoveEvent } from "vue-draggable-plus";
import { NCollapseItem, NCheckbox } from "naive-ui";
import { type TaskResponse } from '@/services/taskServices';
import dayjs from 'dayjs';
import { computed, ref, watch } from "vue";

const props = defineProps<{ title: string, tasksArray: TaskResponse[], name: string }>();
const emits = defineEmits<{
  (e: 'shouldExpand', name: string, expand: boolean): void,
  (e: 'updateTask', data: TaskResponse, name: string): void
  (e: 'setSelectedTask', data: TaskResponse): void
  (e: 'sendErrorMessage', data: { status: number; messageTitle: string; message: string }): void
}>()

const tasksArrayList = ref<TaskResponse[]>([])

const putProperties = computed<boolean | string[]>(() => {
  switch (props.name) {
    case 'today':
      return ['tomorrow', 'pending', 'past']
    case 'tomorrow':
      return ['today', 'pending', 'past']
    case 'pending':
      return false
    case 'past':
      return false
    default:
      return []
  }
})


function checkIfUpdateStartDatePossible(event: MoveEvent) {
  const taskId = Number(event.dragged.dataset.taskId)
  const target = event.to.dataset.name

  const findData = props.tasksArray.find((data) => data.ID === taskId)

  if (target === 'tomorrow') {
    const tomorrow = dayjs().startOf('day').add(1, 'day')

    if (findData && dayjs(findData.due_date, "DD/MM/YYYY").isBefore(tomorrow)) {
      const errorMessage = {
        message: 'Cannot change task if pass due date',
        messageTitle: 'Edit Start Date Failed',
        status: 400
      }
      emits('sendErrorMessage', errorMessage)
      return false
    }
    return true
  }
}

watch(() => props.tasksArray, (newValue) => {
  tasksArrayList.value = newValue
  emits('shouldExpand', props.name, newValue.length === 0)
}, { immediate: true, deep: true })

</script>
<template>
  <n-collapse-item :title="props.title" :name="props.name">
    <template #header-extra>
      <div v-if="tasksArray.length > 0"
        class="flex items-center justify-center h-8 w-8 rounded-[50%] text-sm bg-[#3a3a3a] p-2">
        <p class="text-xs text-white">{{ tasksArray.length }}</p>
      </div>
    </template>
    <vue-draggable v-model="tasksArrayList" ghostClass="ghost" :animation="150"
      :group="{ name: props.name, pull: true, put: putProperties }" :sort="false"
      @add="(event) => emits('updateTask', event.data, props.name)"
      :onMove="(event) => checkIfUpdateStartDatePossible(event)" :data-name="props.name">
      <div v-for="(tasks, index) in tasksArrayList" :key="tasks.ID" :data-task-id="tasks.ID">
        <div @click="emits('setSelectedTask', tasks)"
          :class="['flex justify-start items-center', index === 0 ? '' : 'mt-4']">
          <n-checkbox :checked="tasks.completed" />
          <div class="ml-4 rounded-lg backdrop-blur-sm inset-shadow-[0_0_80px_rgba(0,0,0,0.25)] px-4 py-1 cursor-grab">
            <p class="text-[#8a8888] font-jakarta">
              Created At:
              {{ dayjs(tasks.CreatedAt).format("DD/MM/YYYY HH:mm:ss") }}
            </p>
            <p class="text-white text-base font-jakarta truncate w-45">{{ tasks.title }}</p>
          </div>
        </div>
      </div>
    </vue-draggable>
    <div v-if="tasksArray.length === 0" class="flex justify-center">
      <p class="text-white text-base font-jakarta">No tasks available</p>
    </div>
  </n-collapse-item>
</template>
<style scoped></style>
