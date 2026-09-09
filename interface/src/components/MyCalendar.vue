<script setup lang="ts">
import dayjs from 'dayjs';
import { useItemScale } from '@/composable/pageAdjuster';

import chevronLeftIcon from '@/assets/icons/chevron-left.svg';
import chevronRightIcon from '@/assets/icons/chevron-right.svg';
import { NButton, NDropdown, NIcon, NCalendar } from 'naive-ui';
import { ref } from 'vue';

type Option = {
    label: string;
    key: string;
};

const scale = useItemScale()
const selectedDay = dayjs()
const selectedOption = ref<Option | undefined>({ label: 'Month', key: 'month' })

const options = [{ label: 'Month', key: 'month' }, { label: 'Week', key: 'week' }, { label: 'Day', key: 'day' }, { label: 'Agenda', key: 'agenda' }]

function handleOptionSelect(option: string | number) {
    selectedOption.value = options.find((opt) => opt.key === option)
}
</script>
<template>
    <div class="relative z-10 min-w-full px-4 py-8 h-screen scale-container"
        :style="{ transform: `scale(${scale.scale})`, transformOrigin: 'top left', zoom: scale.zoom }">
        <div class="flex flex-col">
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
                    <p class="text-white font-jakarta text-4xl">{{ selectedDay.format('DD MMMM YYYY') }}</p>
                </div>
                <div>
                    <n-dropdown trigger="click" :options="options" @select="(option) => handleOptionSelect(option)">
                        <n-button>
                            {{ selectedOption?.label ? selectedOption.label : 'Select View' }}
                        </n-button>
                    </n-dropdown>
                </div>
            </div>
            <div>
                <n-calendar />
            </div>
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
</style>