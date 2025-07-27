import http from "./http";

export const webhookApi = {
  // 创建webhook
  createWebhook(data: any) {
    return http.post("/v1/webhook", data);
  },

  // 获取流水线的webhooks（分页）
  getWebhooks(params: any) {
    return http.get("/v1/webhook", {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        pipelineId: params.pipelineId,
        name: params.name,
        status: params.status,
      },
    });
  },
};
