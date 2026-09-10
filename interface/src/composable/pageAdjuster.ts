import { computed, onMounted, onUnmounted, ref } from "vue";

export function useItemScale() {
  const width = ref(window.innerWidth);

  const update = () => {
    width.value = window.innerWidth;
  };

  onMounted(() => {
    window.addEventListener("resize", update);
  });

  onUnmounted(() => {
    window.removeEventListener("resize", update);
  });

  const scale = computed(() => {
    if (width.value < 1200) return   {scale: 0.6, zoom: '140%'};
    if (width.value < 1400) return  {scale: 0.9, zoom: '110%'};
    return {scale: 1, zoom: '100%'};
  });

  return scale;
}

export function useItemScaleV2() {
  const width = ref(window.innerWidth);

  const update = () => {
    width.value = window.innerWidth;
  };

  onMounted(() => {
    window.addEventListener("resize", update);
  });

  onUnmounted(() => {
    window.removeEventListener("resize", update);
  });

const scale = computed(() => {
    if (width.value < 1200) return   {scale: 0.6, zoom: '140%'};
    if (width.value < 1300) return  {scale: 0.8, zoom: '110%'};
    return {scale: 1, zoom: '100%'};
  });

  return scale;
}

export function drawerItemScale() {
  const width = ref(window.innerWidth);

  const update = () => {
    width.value = window.innerWidth;
  };

  onMounted(() => {
    window.addEventListener("resize", update);
  });

  onUnmounted(() => {
    window.removeEventListener("resize", update);
  });

  const scale = computed(() => {
    if (width.value < 1100) return 0.65;
    if (width.value < 1200) return 0.7;
    if (width.value < 1350) return 0.8;
    return 1;
  });

  return scale;
}

export function useCalendarSizeAdjuster() {
  const width = ref(window.innerWidth);

  const update = () => {
    width.value = window.innerWidth;
  };

  onMounted(() => {
    window.addEventListener("resize", update);
  });

  onUnmounted(() => {
    window.removeEventListener("resize", update);
  });

  const scale = computed(() => {
    if (width.value < 1000) return { scaleY: 0.5, scaleX: 0.65 };
    if (width.value < 1100) return { scaleY: 0.55, scaleX: 0.7 };
    if (width.value < 1200) return { scaleY: 0.65, scaleX: 0.75 };
    if (width.value < 1300) return { scaleY: 0.7, scaleX: 0.8 };
    if (width.value < 1400) return { scaleY: 0.8, scaleX: 0.85 };
    if (width.value < 1710) return { scaleY: 0.85, scaleX: 0.9 };
    return { scaleY: 1, scaleX: 1 };
  });

  return scale;
}
