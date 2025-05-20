import axios from "axios"
import { ElMessage } from "element-plus"
import { useAuthStore } from "@/stores/auth"

const http = axios.create({
    baseURL: "/api", // 基础URL，根据实际情况修改
    timeout: 10000, // 超时时间
})

// 请求拦截器
http.interceptors.request.use(
    (config) => {
        // 在发送请求之前做些什么
        const authStore = useAuthStore()
        if (authStore.token) {
            config.headers.Authorization = `Bearer ${authStore.token}` // 设置token
        }
        return config
    },
    (error) => {
        // 处理请求错误
        console.error("Request error:", error)
        return Promise.reject(error)
    },
)

// 响应拦截器
http.interceptors.response.use(
    (response) => {
        // 2xx 范围内的状态码都会触发该函数。
        // 对响应数据做点什么
        return response.data
    },
    (error) => {
        // 超出 2xx 范围的状态码都会触发该函数。
        // 处理响应错误
        console.error("Response error:", error.response)

        if (error.response) {
            switch (error.response.status) {
                case 401:
                    // 未授权，清除token并跳转到登录页
                    const authStore = useAuthStore()
                    authStore.logout()
                    ElMessage.error("登录已过期，请重新登录")
                    window.location.href = "/login" // 替换为你的登录页路由
                    break
                case 403:
                    ElMessage.error("没有权限执行此操作")
                    break
                case 404:
                    ElMessage.error("资源未找到")
                    break
                case 500:
                    ElMessage.error("服务器内部错误")
                    break
                default:
                    ElMessage.error(`请求失败: ${error.response.status}`)
            }
        } else {
            ElMessage.error("请求失败，请检查网络")
        }

        return Promise.reject(error)
    },
)

export default http
