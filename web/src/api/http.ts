import axios from "axios";
import { ElMessage } from "element-plus";
import { useAuthStore } from "@/stores/auth";

const http = axios.create({
  baseURL: "/api",
  timeout: 10000,
});

http.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`;
    }
    return config;
  },
  (error) => {
    console.error("Request error:", error);
    return Promise.reject(error);
  }
);

http.interceptors.response.use(
  (response) => {
    return response.data;
  },
  (error) => {
    console.error("Response error:", error.response);

    if (error.response) {
      switch (error.response.status) {
        case 401:
          const authStore = useAuthStore();
          authStore.logout();
          ElMessage.error("登录已过期，请重新登录");
          window.location.href = "/login";
          break;
        case 403:
          ElMessage.error("没有权限执行此操作");
          break;
        case 404:
          ElMessage.error("资源未找到");
          break;
        case 500:
          ElMessage.error("服务器内部错误");
          break;
        default:
          ElMessage.error(`请求失败: ${error.response.status}`);
      }
    } else {
      ElMessage.error("请求失败，请检查网络");
    }

    return Promise.reject(error);
  }
);

export default http;
