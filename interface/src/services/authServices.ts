import axios from "axios";
import { useAuthStore } from "@/stores/auth";
import router from "@/router";

const authStore = useAuthStore();
const token = localStorage.getItem("token");

export function signIn(data: { username: string; password: string }) {
  axios({
    method: "POST",
    url: "http://localhost:8080/login",
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

export function register(data: {
  firstName: string;
  lastName: string;
  username: string;
  password: string;
}) {
  axios({
    method: "POST",
    url: "http://localhost:8080/register",
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

export function verifyTokenAndLogin(username: string) {
  axios({
    method: "POST",
    url: `http://localhost:8080/verifyTokenAndLogin?username=${username}`,
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
