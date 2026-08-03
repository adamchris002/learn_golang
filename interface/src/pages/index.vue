<route lang="json">{
  "meta": {
    "requiresAuth": true
  }
}</route>
<script setup lang="ts">
import router from "@/router"
import { computed, onMounted, ref, shallowRef, type Component } from "vue"
import { useAuthStore } from "@/stores/auth"
import { NButton, NIcon } from "naive-ui"

import homeBackground from "@/assets/images/home_background.png"
import settingIcon from "@/assets/icons/settings.svg"
import thumbtackIcon from "@/assets/icons/thumbtack.svg"
import { useTimeNow } from "@/composable/timeNow"
import { drawerItemScale } from "@/composable/pageAdjuster"
import { menuItems } from "@/datas/homeItems"
import { verifyTokenAndLogin } from "@/services/authServices"
import { getOneWeekTasks, getTodaysTasks, getAllIncompleteTasks, type TaskResponse } from "@/services/taskServices"
import MyDay from "@/components/MyDay.vue"
import Next7Days from "@/components/Next7Days.vue"
import AllTaskDrawer from "@/components/AllTaskDrawer.vue"

const authStore = useAuthStore()
const scale = drawerItemScale()

const currentComponent = shallowRef<Component | null>(null);

const showSettings = ref<boolean>(true)
const settingsPinned = ref<boolean>(true)
const showTaskDrawer = ref<boolean>(false)

const user = JSON.parse(localStorage.getItem("user") || "{}")
const { currentTime } = useTimeNow()

const todayTaskArray = ref<TaskResponse[]>([])
const next7DaysTaskArray = ref<TaskResponse[]>([])
const allIncompleteTasks = ref<TaskResponse[]>([])

const componentProps = computed(() => {
  switch (currentComponent.value) {
    case MyDay:
      return {
        todaysTaskArray: todayTaskArray.value,
        onRequestCallTodaysTask: refreshAllTasks,
      }

    case Next7Days:
      return {
        nextSevenDaysTaskArray: next7DaysTaskArray.value,
        onRequestCallNextSevenDays: refreshAllTasks,
      }

    default:
      return {}
  }
})

function toggleTasksDrawer(value: boolean) {
  showTaskDrawer.value = value
}

function handleLogout() {
  authStore.logout()
  router.push('/Login')
}

const isSettingsVisible = computed(() => {
  return showSettings.value || settingsPinned.value
})

async function callAllIncompleteTasks() {
  const result = await getAllIncompleteTasks(user.id);

  if (result.success) {
    allIncompleteTasks.value = result.data.sort(
      (a: TaskResponse, b: TaskResponse) =>
        new Date(a.CreatedAt).getTime() -
        new Date(b.CreatedAt).getTime()
    );
  } else {
    authStore.setMessage(result.messageData);
  }
}

async function callTodaysTasks() {
  const result = await getTodaysTasks(user.id);

  if (result.success) {
    todayTaskArray.value = result.data.sort(
      (a: TaskResponse, b: TaskResponse) =>
        new Date(a.CreatedAt).getTime() -
        new Date(b.CreatedAt).getTime()
    );
  } else {
    authStore.setMessage(result.messageData);
  }
}

async function callOneWeekTask() {
  const result = await getOneWeekTasks(user.id);

  if (result.success) {
    next7DaysTaskArray.value = result.data.sort(
      (a: TaskResponse, b: TaskResponse) =>
        new Date(a.CreatedAt).getTime() -
        new Date(b.CreatedAt).getTime()
    );
  } else {
    authStore.setMessage(result.messageData);
  }
}

async function refreshAllTasks() {
  await Promise.all([
    callTodaysTasks(),
    callOneWeekTask(),
    callAllIncompleteTasks()
  ])
}

onMounted(async () => {
  currentComponent.value = menuItems[0]?.component ?? null
  if (authStore.token !== null) {
    await verifyTokenAndLogin(user.username)
  }
  await callTodaysTasks()
  await callOneWeekTask()
  await callAllIncompleteTasks()
})

</script>
<template>
  <div class="flex justify-start w-screen h-screen">
    <div class="absolute inset-0 bg-no-repeat pointer-events-none" :style="{
      backgroundImage: `url(${homeBackground})`,
      backgroundSize: '70%',
      backgroundPosition: '200% center'
    }" />
    <div
      :class="[`${isSettingsVisible ? 'w-[15%]' : 'w-[5%]'} ${isSettingsVisible ? 'bg-[#0d0d0d]' : 'bg-[#1c1c1c]'} ${scale < 1 ? '' : 'py-5 px-5'}`]"
      @mouseenter="showSettings = true" @mouseleave="showSettings = false">
      <n-button v-if="!isSettingsVisible"
        :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc' }" class="setting-btn"
        circle>
        <n-icon>
          <settingIcon />
        </n-icon>
      </n-button>
      <div v-else-if="isSettingsVisible" :style="{ transform: `scale(${scale})`, transformOrigin: 'center' }"
        class="scale-container flex flex-col justify-between h-full">
        <div>
          <div class="flex justify-between ">
            <n-button :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', }"
              class="setting-btn" circle>
              <n-icon>
                <settingIcon />
              </n-icon>
            </n-button>
            <div class="w-full ml-4">
              <div class="flex items-center justify-between">
                <p class="text-white font-jakarta text-lg font-semibold mr-8">{{ user.first_name }}</p>
                <n-button @click="settingsPinned = !settingsPinned"
                  :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc' }"
                  class="pin-btn" circle>
                  <thumbtackIcon />
                </n-button>
              </div>
              <p class="text-white opacity-75">
                {{ currentTime.format('dddd, MMMM D') }}
              </p>
              <p class="text-white opacity-75">
                {{ currentTime.format('HH:mm:ss') }}
              </p>

            </div>
          </div>
          <div v-for="item in menuItems" :key="item.title">
            <div class="mt-4">
              <n-button @click="currentComponent = item.component" block :bordered="false"
                :theme-overrides="{ textColor: 'white', textColorHover: '#0373fc', textColorFocus: '#0373fc', borderHover: '1px solid #0373fc' }"
                type="info" ghost class="item-btn !justify-start !text-left">
                {{ item.title }}
              </n-button>
            </div>
          </div>
        </div>
        <n-button @click="handleLogout" class="mt-4" block color="#003465">Logout</n-button>
      </div>
    </div>

    <div :class="[`${isSettingsVisible ? 'w-[85%]' : 'w-[95%]'} bg-[#1c1c1c]`]">
      <component :is="currentComponent" v-bind="componentProps" />
    </div>

    <AllTaskDrawer :show-drawer="showTaskDrawer" :all-incomplete-tasks="allIncompleteTasks" @toggle-drawer="toggleTasksDrawer" @request-refresh-all-tasks="refreshAllTasks" />
  </div>

</template>
<style scoped>
.setting-btn {
  color: white;
  border-color: rgba(255, 255, 255, 0.3);
}

.setting-btn:hover {
  color: #0373fc;
  border-color: #0373fc;
}

.setting-btn:focus {
  color: #0373fc;
  border-color: #0373fc;
}

.pin-btn {
  color: white;
  border-color: rgba(255, 255, 255, 0.3);
}

.pin-btn:hover {
  color: #0373fc;
  border-color: #0373fc;
}

.pin-btn:focus {
  color: #0373fc;
  border-color: #0373fc;
}

.scale-container {
  transition: transform 300ms ease-in-out;
  zoom: 100%;
}

.item-btn {
  color: white;
}

.item-btn:hover {
  color: #0373fc;
}
</style>