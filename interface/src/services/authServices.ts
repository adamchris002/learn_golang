import axios from "axios";
import { useAuthStore } from "@/stores/auth";
import router from "@/router";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
});

export async function signIn(data: { username: string; password: string }) {
  const authStore = useAuthStore();

  api({
    method: "POST",
    url: `/login`,
    data: data,
  })
    .then((result) => {
      if (result.status === 200) {
        authStore.setToken(result.data.token);
        authStore.setUser(result.data.user);
        const successData = {
          status: result.status,
          messageTitle: result.data.messageTitle,
          message: result.data.message,
        };
        authStore.setMessage(successData);
        router.push("/");
      }
    })
    .catch((err) => {
      const errorData = {
        status: err.response.status,
        messageTitle: err.response.data.messageTitle,
        message: err.response.data.message,
      };
      authStore.setMessage(errorData);
    });
}

export async function register(data: {
  firstName: string;
  lastName: string;
  username: string;
  password: string;
}) {
  const authStore = useAuthStore();

  api({
    method: "POST",
    url: `/register`,
    data: data,
  })
    .then((result) => {
      if (result.status === 200) {
        router.push("/Login");
        const successData = {
          status: result.status,
          messageTitle: result.data.messageTitle,
          message: result.data.message,
        };
        authStore.setMessage(successData);
      }
    })
    .catch((err) => {
      const errorData = {
        status: err.response.status,
        messageTitle: err.response.data.messageTitle,
        message: err.response.data.message,
      };
      authStore.setMessage(errorData);
    });
}

export async function verifyTokenAndLogin(username: string) {
  const authStore = useAuthStore();
  const token = localStorage.getItem("token");

  api({
    method: "POST",
    url: `/verifyTokenAndLogin?username=${username}`,
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
    .then((result) => {
      const messageData = {
        status: result.status,
        messageTitle: result.data.messageTitle,
        message: result.data.message
      }
      authStore.setMessage(messageData)
      router.push("/")
    })
    .catch((err) => {
      const errorData = {
        status: err.response.status,
        messageTitle: err.response.data.messageTitle,
        message: err.response.data.message,
      };
      authStore.setMessage(errorData);
      authStore.logout()
      router.push("/login");
    });
}
