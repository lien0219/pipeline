import { defineStore } from "pinia";
import { ref } from "vue";
import { pipelineApi } from "@/api/pipeline";

export const usePipelineStore = defineStore("pipeline", () => {
  const pipelines = ref([]);
  const currentPipeline = ref(null);
  const loading = ref(false);
  const error = ref(null);

  async function fetchPipelines(params = {}) {
    loading.value = true;
    error.value = null;

    try {
      const response = await pipelineApi.getPipelines(params);
      pipelines.value = response.data;
      return response;
    } catch (err) {
      error.value = err.message || "Failed to fetch pipeline";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchPipelineById(id) {
    loading.value = true;
    error.value = null;

    try {
      const response = await pipelineApi.getPipelineById(id);
      currentPipeline.value = response.data;
      return response;
    } catch (err) {
      error.value = err.message || "Failed to fetch pipeline";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function createPipeline(pipelineData) {
    loading.value = true;
    error.value = null;

    try {
      const response = await pipelineApi.createPipeline(pipelineData);
      pipelines.value.push(response.data);
      return response;
    } catch (err) {
      error.value = err.message || "Failed to create pipeline";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function updatePipeline(id, pipelineData) {
    loading.value = true;
    error.value = null;

    try {
      const response = await pipelineApi.updatePipeline(id, pipelineData);
      const index = pipelines.value.findIndex((p) => p.id === id);
      if (index !== -1) {
        pipelines.value[index] = response.data;
      }
      if (currentPipeline.value && currentPipeline.value.id === id) {
        currentPipeline.value = response.data;
      }
      return response;
    } catch (err) {
      error.value = err.message || "Failed to update pipeline";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function deletePipeline(id) {
    loading.value = true;
    error.value = null;

    try {
      await pipelineApi.deletePipeline(id);
      pipelines.value = pipelines.value.filter((p) => p.id !== id);
      if (currentPipeline.value && currentPipeline.value.id === id) {
        currentPipeline.value = null;
      }
    } catch (err) {
      error.value = err.message || "Failed to delete pipeline";
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function getDashboardStats() {
    try {
      const response = await pipelineApi.getDashboardStats();
      return response.data;
    } catch (err) {
      error.value = err.message || "Failed to get dashboard stats";
      throw err;
    }
  }

  async function getDashboardActivities(limit) {
    try {
      const response = await pipelineApi.getDashboardActivities(limit);
      return response.data;
    } catch (err) {
      error.value = err.message || "Failed to get dashboard activities";
      throw err;
    }
  }
  async function triggerPipeline(id) {
    loading.value = true;
    error.value = null;

    try {
      const response = await pipelineApi.triggerPipeline(id);
      return response;
    } catch (err) {
      error.value = err.message || "Failed to trigger pipeline";
      throw err;
    } finally {
      loading.value = false;
    }
  }
  return {
    pipelines,
    currentPipeline,
    loading,
    error,
    fetchPipelines,
    fetchPipelineById,
    createPipeline,
    updatePipeline,
    deletePipeline,
    getDashboardStats,
    getDashboardActivities,
    triggerPipeline,
  };
});
