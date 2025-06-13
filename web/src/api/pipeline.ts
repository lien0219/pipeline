import http from "./http";

export const pipelineApi = {
  getPipelines(params:any) {
    return http.get("/v1/pipeline", { params });
  },

  getPipelineById(id:string|number) {
    return http.get(`/v1/pipeline/${id}`);
  },

  createPipeline(data:any) {
    return http.post("/v1/pipeline", data);
  },

  updatePipeline(id:number|string, data:any) {
    return http.put(`/v1/pipeline/${id}`, data);
  },

  deletePipeline(id:string|number) {
    return http.delete(`/v1/pipeline/${id}`);
  },

  triggerPipeline(id:string|number) {
    return http.post(`/v1/pipeline/${id}/trigger`);
  },

  getPipelineRuns(id:string|number, params:any) {
    return http.get(`/v1/pipeline/${id}/runs`, { params });
  },

  getPipelineRunById(pipelineId:string|number, runId:string|number) {
    return http.get(`/v1/pipeline/${pipelineId}/runs/${runId}`);
  },

  getPipelineRunLogs(pipelineId:string|number, runId:string|number) {
    return http.get(`/v1/pipeline/${pipelineId}/runs/${runId}/logs`);
  },

  cancelPipelineRun(pipelineId:string|number, runId:number|string) {
    return http.post(`/v1/pipeline/${pipelineId}/runs/${runId}/cancel`);
  },

  getArtifacts(params:any) {
    return http.get("/v1/artifacts", { params });
  },

  getArtifactById(id:number|string) {
    return http.get(`/v1/artifacts/${id}`);
  },

  deleteArtifact(id:number|string) {
    return http.delete(`/v1/artifacts/${id}`);
  },
  getDashboardStats() {
    return http.get("/v1/dashboard/stats");
  },
  getDashboardActivities(data:any) {
    return http.post("/v1/dashboard/activities", data);
  },
  getDAGByPipelineID(id: number) {
    return http.get(`/v1/dag/pipeline/${id}`); 
  },
  
  updateDAG(dagId: number, data: any) {
    return http.put(`/v1/dag/${dagId}`, data); 
  },
};
