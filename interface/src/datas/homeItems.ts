import MyDay from "@/components/MyDay.vue";
import Next7Days from "@/components/Next7Days.vue";
import AllMyTasks from "@/components/AllMyTasks.vue";
import MyCalendar from "@/components/MyCalendar.vue";

export const menuItems = [
  { title: "My day", component: MyDay },
  { title: "Next 7 days", component: Next7Days },
  { title: "All my tasks", component: AllMyTasks },
  { title: "My Calendar", component: MyCalendar },
];
