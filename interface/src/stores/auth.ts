import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useAuthStore = defineStore("auth", () => {
  const token = ref<string | null>(localStorage.getItem("token"));
  const user = ref<{
    id: string;
    firstName: string;
    lastName: string;
    username: string;
    email: string;
  } | null>(null);

  const message = ref<{ status: Number; messageTitle: string; message: string } | null>(null);

  const isAuthenticated = computed(() => !!token.value);

  function setToken(newToken: string) {
    token.value = newToken;
    localStorage.setItem("token", newToken);
  }

  function setUser(newUser: {
    id: string;
    firstName: string;
    lastName: string;
    username: string;
    email: string;
  }) {
    user.value = newUser;
  }

  function setMessage(err: { status: Number; messageTitle: string; message: string } | null) {
    message.value = err;
  }

  function logout() {
    token.value = null;
    user.value = null;
    localStorage.removeItem("token");
  }

  return {
    token,
    user,
    message,
    isAuthenticated,
    setToken,
    setUser,
    setMessage,
    logout,
  };
});
