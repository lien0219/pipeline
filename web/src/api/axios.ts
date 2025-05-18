import axios from 'axios';
import { useAuthStore } from '@/stores/auth';
import router from '@/router';

// 创建axios实例
const apiClient = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
    }
});

// 设置拦截器
export function setupAxiosInterceptors() {
    // 请求拦截器
    apiClient.interceptors.request.use(
        (config) => {
            const authStore = useAuthStore();

            // 如果有token，添加到请求头
            if (authStore.token) {
                config.headers.Authorization = `Bearer ${authStore.token}`;
            }

            return config;
        },
        (error) => {
            return Promise.reject(error);
        }
    );

    // 响应拦截器
    apiClient.interceptors.response.use(
        (response) => {
            return response;
        },
        (error) => {
            const { response } = error;

            // 处理401错误（未授权）
            if (response && response.status === 401) {
                const authStore = useAuthStore();

                // 清除认证状态
                authStore.logout();

                // 重定向到登录页
                router.push({
                    path: '/auth/login',
                    query: { redirect: router.currentRoute.value.fullPath }
                });
            }

            return Promise.reject(error);
        }
    );
}

export default apiClient;