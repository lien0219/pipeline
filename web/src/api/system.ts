import apiClient, { AxiosResponse } from './axios';

// 定义接口类型
interface SystemConfigResponse {
    config: {
        version: string;
        environment: string;
        [key: string]: any; // 其他可能的字段
    };
}

interface HealthCheckResponse {
    status: string;
    uptime: number;
    [key: string]: any; // 其他可能的字段
}

export const systemApi = {
    // 获取系统配置
    getSystemConfig(): Promise<AxiosResponse<SystemConfigResponse>> {
        return apiClient.get('/system/config');
    },

    // 健康检查
    healthCheck(): Promise<AxiosResponse<HealthCheckResponse>> {
        return apiClient.get('/health');
    }
};