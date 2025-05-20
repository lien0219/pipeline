import http from "./http"

export const pipelineApi = {
    getPipelines(params) {
        return http.get("/pipelines", { params })
    },

    getPipelineById(id) {
        return http.get(`/pipelines/${id}`)
    },

    createPipeline(data) {
        return http.post("/pipelines", data)
    },

    updatePipeline(id, data) {
        return http.put(`/pipelines/${id}`, data)
    },

    deletePipeline(id) {
        return http.delete(`/pipelines/${id}`)
    },

    triggerPipeline(id) {
        return http.post(`/pipelines/${id}/trigger`)
    },

    getPipelineRuns(id, params) {
        return http.get(`/pipelines/${id}/runs`, { params })
    },

    getPipelineRunById(pipelineId, runId) {
        return http.get(`/pipelines/${pipelineId}/runs/${runId}`)
    },

    getPipelineRunLogs(pipelineId, runId) {
        return http.get(`/pipelines/${pipelineId}/runs/${runId}/logs`)
    },

    cancelPipelineRun(pipelineId, runId) {
        return http.post(`/pipelines/${pipelineId}/runs/${runId}/cancel`)
    },

    getArtifacts(params) {
        return http.get("/artifacts", { params })
    },

    getArtifactById(id) {
        return http.get(`/artifacts/${id}`)
    },

    deleteArtifact(id) {
        return http.delete(`/artifacts/${id}`)
    },
}
