<script setup lang="ts">
import dayjs from 'dayjs';
import customParseFormat from 'dayjs/plugin/customParseFormat'
import { useCalendarSizeAdjuster } from '@/composable/pageAdjuster';

import chevronLeftIcon from '@/assets/icons/chevron-left.svg';
import chevronRightIcon from '@/assets/icons/chevron-right.svg';
import { NButton, NDropdown, NIcon, NCalendar, NConfigProvider, darkTheme, NDataTable, type DataTableColumns, } from 'naive-ui';
import { computed, onMounted, ref, watch } from 'vue';

dayjs.extend(customParseFormat)

type Option = {
    label: string;
    key: string;
};

type WeekRow = {
    label: string
    [key: string]: string
}

const calendarSize = useCalendarSizeAdjuster()
const selectedDay = dayjs()
const selectedOption = ref<Option | undefined>({ label: 'Month', key: 'month' })

//week datas
const weekDate = ref<string[]>([])
const timeDate = ref<string[]>([])

const options = [{ label: 'Month', key: 'month' }, { label: 'Week', key: 'week' }, { label: 'Day', key: 'day' }, { label: 'Agenda', key: 'agenda' }]

function handleOptionSelect(option: string | number) {
    selectedOption.value = options.find((opt) => opt.key === option)
}

function getWeek() {
    const startOfWeek = selectedDay.startOf('week').add(1, 'day')
    weekDate.value = []
    for (let i = 0; i < 7; i++) {
        weekDate.value.push(startOfWeek.add(i, 'day').format('DD/MM/YYYY'))
    }
}

function generateTime() {
    const startOfDay = selectedDay.startOf('day')
    const timeArray: string[] = []
    for (let i = 0; i < 24; i++) {
        const time = startOfDay.add(i, 'hour').format('HH:mm')
        timeArray.push(time)
    }
    timeDate.value = timeArray
}

onMounted(() => {
    generateTime()
})

watch(() => selectedDay, (newValue) => {
    getWeek()
}, { immediate: true })
</script>
<template>
    <div class="z-10 w-full px-4 py-8 h-screen">
        <div class="flex justify-between items-center mb-8">
            <div class="flex items-center">
                <div class="cursor-pointer rounded-4xl border-1 border-[#a3a3a3] bg-[#1c1c1c] px-4 py-2">
                    <p class="text-[#a3a3a3] text-xl">Today</p>
                </div>
                <n-button circle ghost :bordered="false" class="chevron-btn size-fit"
                    :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', rippleColor: 'none' }">
                    <n-icon size="36">
                        <chevronLeftIcon class="text-4xl" />
                    </n-icon>
                </n-button>
                <n-button circle ghost :bordered="false" class="chevron-btn size-fit"
                    :theme-overrides="{ borderHover: '1px solid #0373fc', borderFocus: '1px solid #0373fc', rippleColor: 'none' }">
                    <n-icon size="36">
                        <chevronRightIcon class="text-4xl" />
                    </n-icon>
                </n-button>
                <p class="text-white font-jakarta text-4xl">{{ selectedDay.format('MMMM YYYY') }}</p>
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
        <div v-show="selectedOption?.key === 'month'" class="w-full backdrop-blur-sm"
            :style="{ transform: `scaleX(${calendarSize.scaleX}) scaleY(${calendarSize.scaleY})`, transformOrigin: 'top' }">
            <n-config-provider :theme="darkTheme">
                <n-calendar />
            </n-config-provider>
        </div>
        <div v-show="selectedOption?.key === 'week'"
            class="w-full max-h-[75vh] overflow-y-auto bg-[#1c1c1c] backdrop-blur-sm custom-scroll">
            <table class="w-full table-fixed">
                <thead>
                    <tr class="sticky top-0 z-10 bg-[#1c1c1c]">
                        <th class="w-16 p-4"></th>
                        <th class="p-4" v-for="day in weekDate" :key="day">
                            <div :class="[
                                'flex flex-col items-center justify-center w-12 h-12 mx-auto text-white',
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
                            </p>
                        </td>
                        <td v-for="day in weekDate" :key="day" class="border border-gray-600 h-16">
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
        <div v-show="selectedOption?.key === 'day'"
            class="w-full max-h-[75vh] overflow-y-auto bg-[#1c1c1c] backdrop-blur-sm custom-scroll">
            <div class="w-full flex justify-center items-center sticky top-0 z-10 bg-[#1c1c1c]"">
                <div class=" bg-blue-500 rounded-full w-12 h-12 mx-auto flex flex-col items-center justify-center">
                <p class="text-white">{{ dayjs().format('ddd') }}</p>
                <p class="text-white text-xl">{{ dayjs().format('DD') }}</p>
            </div>
        </div>
        <table class="w-full table-fixed">
            <thead>
                <tr class="sticky top-0 z-10 bg-[#1c1c1c]">
                    <th class="w-16 p-4"></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(time, index) in timeDate" :key="time">
                    <td class="relative">
                        <p :class="['absolute text-white text-sm', index !== 0 ? '-top-3' : '-top-1']">
                            {{ dayjs(time, 'HH:mm').format("h:mm A") }}
                        </p>
                    </td>
                    <td v-for="day in weekDate" :key="day" class="border-t border-gray-600 h-16">
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
    <div v-show="selectedOption?.key === 'agenda'">

    </div>
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

:deep(.n-calendar-header) {
    display: none !important;
}

:deep(.week-table th) {
    white-space: pre-line;
}

:deep(.n-data-table-th__title) {
    text-align: center !important;
}

.custom-scroll::-webkit-scrollbar {
    display: none;
}
</style>