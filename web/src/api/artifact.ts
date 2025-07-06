import http from "./http";
import { type BaseResponse, type PageResult } from "./types";

// 制品类型定义
export interface Artifact {
  id: number;
  name: string;
  version: string;
  type: string;
  size: number;
  url: string;
  created_at: string;
  created_by: number;
  pipeline_id?: number;
  pipeline?: { name: string };
  pipeline_run_id?: number;
  download_count?: number;
  description?: string;
  docker_image?: string;
}

// 创建制品请求参数
export interface CreateArtifactRequest {
  name: string;
  version: string;
  type: string;
  pipeline_id?: number;
  description?: string;
  docker_image?: string;
  file?: File;
}

// 制品管理API
export const artifactApi = {
  // 获取制品列表（支持分页和筛选）
  getArtifacts: (page = 1, pageSize = 10, filters?: Record<string, any>) => {
    return http.get<BaseResponse<PageResult<Artifact>>>("/v1/artifact", {
      params: { page, pageSize, ...filters },
    });
  },

  // 获取制品详情
  getArtifactById: (id: number) => {
    return http.get<BaseResponse<Artifact>>(`/v1/artifact/${id}`);
  },

  // 创建制品
  createArtifact: (data: FormData) => {
    return http.post<BaseResponse<Artifact>>("/v1/artifact", data);
  },

  // 删除制品
  deleteArtifact: (id: number) => {
    return http.delete<BaseResponse<null>>(`/v1/artifact/${id}`);
  },

  // 下载制品
  downloadArtifact: (id: number) => {
    return http.get(`/v1/artifact/${id}/download`, {
      responseType: "blob",
    });
  },
};
