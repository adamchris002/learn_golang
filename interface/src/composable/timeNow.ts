import { ref, onBeforeUnmount, onMounted } from "vue";
import dayjs from "dayjs";

export function useTimeNow() {
  const currentTime = ref(dayjs());

  let timer: number | undefined;

  const updateCurrentTime = () => {
    currentTime.value = dayjs();
  };

  onMounted(() => {
    updateCurrentTime();

    timer = window.setInterval(updateCurrentTime, 1000);
  });

  onBeforeUnmount(() => {
    if (timer) {
      clearInterval(timer);
    }
  });

  return {
    currentTime,
  };
}
