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
    if (width.value < 1200) return 0.65;
    if (width.value < 1480) return 0.7;
    if (width.value < 1550) return 0.8;
    return 1;
  });

  return scale;
}
