import http from "./http";

export const pipelineApi = {
  getPipelines(params) {
    return http.get("/v1/pipeline", { params });
  },

  getPipelineById(id) {
    return http.get(`/v1/pipeline/${id}`);
  },

  createPipeline(data) {
    return http.post("/v1/pipeline", data);
  },

  updatePipeline(id, data) {
    return http.put(`/v1/pipeline/${id}`, data);
  },

  deletePipeline(id) {
    return http.delete(`/v1/pipeline/${id}`);
  },

  triggerPipeline(id) {
    return http.post(`/v1/pipeline/${id}/trigger`);
  },

  getPipelineRuns(id, params) {
    return http.get(`/v1/pipeline/${id}/runs`, { params });
  },

  getPipelineRunById(pipelineId, runId) {
    return http.get(`/v1/pipeline/${pipelineId}/runs/${runId}`);
  },

  getPipelineRunLogs(pipelineId, runId) {
    return http.get(`/v1/pipeline/${pipelineId}/runs/${runId}/logs`);
  },

  cancelPipelineRun(pipelineId, runId) {
    return http.post(`/v1/pipeline/${pipelineId}/runs/${runId}/cancel`);
  },

  getArtifacts(params) {
    return http.get("/v1/artifacts", { params });
  },

  getArtifactById(id) {
    return http.get(`/v1/artifacts/${id}`);
  },

  deleteArtifact(id) {
    return http.delete(`/v1/artifacts/${id}`);
  },
};
export const triggerPipeline = (id: string) => {
  return http.post(`/v1/pipeline/${id}/trigger`);
};

export const getPipelineRuns = (
  id: string,
  params: {
    page?: number;
    limit?: number;
  }
) => {
  return http.get(`/v1/pipeline/${id}/runs`, { params });
};
