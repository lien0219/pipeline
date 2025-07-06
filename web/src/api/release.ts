import http from "./http";
import { type BaseResponse, type PageResult } from "./types";

// 发布记录类型定义
export interface Release {
  id: number;
  version: string;
  environment: string;
  status: "success" | "failed" | "in_progress" | "rolled_back";
  description: string;
  artifact: string;
  artifact_id: number;
  deployed_by: string;
  deployed_at: string;
  is_rollback: boolean;
  release_notes?: string;
}

// 创建发布请求参数
export interface CreateReleaseRequest {
  version: string;
  environment: string;
  artifact_id: number;
  description?: string;
  release_notes?: string;
}

// 发布管理API
export const releaseApi = {
  // 获取发布列表（支持分页和筛选）
  getReleases: (
    page = 1,
    pageSize = 10,
    environment = "",
    status = "",
    startDate = "",
    endDate = ""
  ) => {
    return http.get<BaseResponse<PageResult<Release>>>("/v1/release", {
      params: {
        page,
        pageSize,
        environment,
        status,
        start_date: startDate,
        end_date: endDate,
      },
    });
  },

  // 获取发布详情
  getReleaseById: (id: number) => {
    return http.get<BaseResponse<Release>>(`/v1/release/${id}`);
  },

  // 创建发布
  createRelease: (data: CreateReleaseRequest) => {
    return http.post<BaseResponse<Release>>("/v1/release", data);
  },

  // 删除发布
  deleteRelease: (id: number) => {
    return http.delete<BaseResponse<void>>(`/v1/release/${id}`);
  },

  // 回滚发布
  rollbackRelease: (id: number) => {
    return http.post<BaseResponse<Release>>(`/v1/release/${id}/rollback`);
  },
};
