import http from "./http";

export const pipelineApi = {
  getPipelines(params) {
    return http.get("/v1/pipeline", { params });
  },

  getPipelineById(id) {
    return http.get(`/v1/pipelines/${id}`);
  },

  createPipeline(data) {
    return http.post("/v1/pipelines", data);
  },

  updatePipeline(id, data) {
    return http.put(`/v1/pipelines/${id}`, data);
  },

  deletePipeline(id) {
    return http.delete(`/v1/pipelines/${id}`);
  },

  triggerPipeline(id) {
    return http.post(`/v1/pipelines/${id}/trigger`);
  },

  getPipelineRuns(id, params) {
    return http.get(`/v1/pipelines/${id}/runs`, { params });
  },

  getPipelineRunById(pipelineId, runId) {
    return http.get(`/v1/pipelines/${pipelineId}/runs/${runId}`);
  },

  getPipelineRunLogs(pipelineId, runId) {
    return http.get(`/v1/pipelines/${pipelineId}/runs/${runId}/logs`);
  },

  cancelPipelineRun(pipelineId, runId) {
    return http.post(`/v1/pipelines/${pipelineId}/runs/${runId}/cancel`);
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
