import apiClient from './axios';

// 定义接口类型
interface Credentials {
    username: string;
    password: string;
}

interface UserData {
    username: string;
    email: string;
    [key: string]: any; // 其他可能的字段
}

interface PasswordData {
    oldPassword: string;
    newPassword: string;
}

interface UserResponse {
    id: number;
    username: string;
    email: string;
    [key: string]: any; // 其他可能的字段
}

export const authApi = {
    // 用户登录
    login(credentials: Credentials): Promise<any> {
        return apiClient.post('/auth/login', credentials);
    },

    // 用户注册
    register(userData: UserData): Promise<any> {
        return apiClient.post('/auth/register', userData);
    },

    // 获取当前用户信息
    getCurrentUser(): Promise<any> {
        return apiClient.get('/users/me');
    },

    // 更新用户信息
    updateUserInfo(userData: UserData): Promise<any> {
        return apiClient.put('/users/me', userData);
    },

    // 修改密码
    changePassword(passwordData: PasswordData): Promise<any> {
        return apiClient.put('/users/password', passwordData);
    }
};