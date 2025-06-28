import http from "./http";
import { type PageResult, type BaseResponse } from "./types";

export interface EnvironmentVariable {
  key: string;
  value: string;
}

export interface Environment {
  id: number;
  name: string;
  type: "development" | "testing" | "staging" | "production";
  url: string;
  description: string;
  status: "active" | "inactive" | "";
  variables: EnvironmentVariable[];
  last_deployed_at?: string;
  created_at: string;
  updated_at: string;
  created_by: number;
  user?: { id: number; username: string };
}

export interface CreateEnvironmentRequest {
  name: string;
  type: "development" | "testing" | "staging" | "production";
  url: string;
  description?: string;
  variables: EnvironmentVariable[];
}

export interface UpdateEnvironmentRequest extends CreateEnvironmentRequest {
  status?: "active" | "inactive";
}

// 环境管理API
export const environmentApi = {
  // 获取环境列表（支持分页和筛选）
  getEnvironments: (
    page = 1,
    pageSize = 10,
    name = "",
    type = "",
    status = ""
  ) => {
    return http.get<BaseResponse<PageResult<Environment>>>(`/v1/environment`, {
      params: { page, pageSize, name, type, status },
    });
  },

  // 获取环境详情
  getEnvironmentById: (id: number) => {
    return http.get<BaseResponse<Environment>>(`/v1/environment/${id}`);
  },

  // 创建环境
  createEnvironment: (data: CreateEnvironmentRequest) => {
    return http.post<BaseResponse<Environment>>(`/v1/environment`, data);
  },

  // 更新环境
  updateEnvironment: (id: number, data: UpdateEnvironmentRequest) => {
    return http.put<BaseResponse<Environment>>(`/v1/environment/${id}`, data);
  },

  // 删除环境
  deleteEnvironment: (id: number) => {
    return http.delete<BaseResponse<void>>(`/v1/environment/${id}`);
  },
};
