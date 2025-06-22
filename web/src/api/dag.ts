import http from "./http";

export const dagApi = {
  // 创建DAG
  create: (data: {
    name: string;
    description: string;
    pipeline_id: number;
    nodes: any[];
  }) => http.post("/v1/dag", data),

  // 获取DAG详情
  getById: (id: number) => http.get(`/v1/dag/${id}`),

  // 获取流水线的所有DAG
  getByPipelineId: (pipelineId: number) =>
    http.get(`/v1/dag/pipeline/${pipelineId}`),

  // 获取所有DAG的分页接口
  getAllDags: (params: { page: number; pageSize: number }) => {
    return http.get<{
      data: {
        items: any[];
        total: number;
        page: number;
        pageSize: number;
      };
    }>("/v1/dag/all", { params });
  },

  // 获取活动DAG
  getActive: (pipelineId: number) =>
    http.get(`/v1/dag/pipeline/${pipelineId}/active`),

  // 更新DAG
  update: (
    id: number,
    data: { name?: string; description?: string; nodes?: any[] }
  ) => http.put(`/v1/dag/${id}`, data),

  // 删除DAG
  delete: (id: number) => http.delete(`/v1/dag/${id}`),

  // 验证DAG
  validate: (nodes: any[]) => http.post("/v1/dag/validate", { nodes }),

  // 创建新版本
  createVersion: (id: number) => http.post(`/v1/dag/${id}/version`),

  // 获取历史版本
  getHistory: (pipelineId: number) =>
    http.get(`/v1/dag/pipeline/${pipelineId}/history`),

  // 激活DAG
  activate: (id: number) => http.post(`/v1/dag/${id}/activate`),
};
