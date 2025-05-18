import { defineStore } from 'pinia';
import { authApi } from '@/api/auth';

// 定义接口类型
interface User {
    id: number;
    username: string;
    email: string;
    [key: string]: any; // 其他可能的字段
}

interface AuthData {
    token: string;
    user: User;
}

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null as User | null,
        token: null as string | null,
        loading: false as boolean,
        error: null as string | null
    }),

    getters: {
        isAuthenticated: (state) => !!state.token,
        currentUser: (state) => state.user
    },

    actions: {
        // 登录
        async login(credentials: { username: string; password: string }): Promise<any> {
            this.loading = true;
            this.error = null;

            try {
                const response = await authApi.login(credentials);
                this.setAuthData(response.data.data);
                return response;
            } catch (error: any) {
                this.error = error.response?.data?.message || '登录失败';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        // 注册
        async register(userData: User): Promise<any> {
            this.loading = true;
            this.error = null;

            try {
                const response = await authApi.register(userData);
                return response;
            } catch (error: any) {
                this.error = error.response?.data?.message || '注册失败';
                throw error;
            } finally {
                this.loading = false;
            }
        },

        // 获取用户信息
        async fetchUserInfo(): Promise<any> {
            if (!this.token) return;

            this.loading = true;

            try {
                const response = await authApi.getCurrentUser();
                this.user = response.data.data;
                return response;
            } catch (error: any) {
                console.error('获取用户信息失败:', error);
                // 如果获取用户信息失败，可能是token无效
                if (error.response?.status === 401) {
                    this.logout();
                }
            } finally {
                this.loading = false;
            }
        },

        // 更新用户信息
        async updateUserInfo(userData: User): Promise<any> {
            this.loading = true;

            try {
                const response = await authApi.updateUserInfo(userData);
                this.user = { ...this.user, ...response.data.data };
                return response;
            } catch (error: any) {
                console.error('更新用户信息失败:', error);
                throw error;
            } finally {
                this.loading = false;
            }
        },

        // 修改密码
        async changePassword(passwordData: { oldPassword: string; newPassword: string }): Promise<any> {
            this.loading = true;

            try {
                const response = await authApi.changePassword(passwordData);
                return response;
            } catch (error: any) {
                console.error('修改密码失败:', error);
                throw error;
            } finally {
                this.loading = false;
            }
        },

        // 设置认证数据
        setAuthData(data: AuthData): void {
            this.token = data.token;
            this.user = data.user;

            // 保存到本地存储
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
        },

        // 恢复会话
        restoreSession(): void {
            const token = localStorage.getItem('token');
            const user = localStorage.getItem('user');

            if (token) {
                this.token = token;
                this.user = user ? JSON.parse(user) : null;

                // 验证token有效性
                this.fetchUserInfo();
            }
        },

        // 登出
        logout(): void {
            this.user = null;
            this.token = null;

            // 清除本地存储
            localStorage.removeItem('token');
            localStorage.removeItem('user');
        }
    }
});