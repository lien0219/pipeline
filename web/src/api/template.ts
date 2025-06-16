import http from "./http";

// 定义构建模板管理的 API 对象
export const templateApi = {
  // 获取所有构建模板
  getBuildTemplates(params: any) {
    return http.get("/v1/build-template", { params });
  },

  // 根据 ID 获取单个构建模板
  getBuildTemplateById(id: string | number) {
    return http.get(`/v1/build-template/${id}`);
  },

  // 创建新的构建模板
  createBuildTemplate(data: any) {
    return http.post("/v1/build-template", data);
  },

  // 更新现有构建模板
  updateBuildTemplate(id: number | string, data: any) {
    return http.put(`/v1/build-template/${id}`, data);
  },

  // 删除构建模板
  deleteBuildTemplate(id: string | number) {
    return http.delete(`/v1/build-template/${id}`);
  },
  // 从模板创建流水线
  createPipelineFromTemplate(templateId: string | number, data: any) {
    return http.post(`/v1/build-template/${templateId}/apply`, data);
  },
};
