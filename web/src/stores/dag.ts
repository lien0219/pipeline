import { defineStore } from "pinia";
import { dagApi } from "@/api/dag";

export const useDagStore = defineStore("dag", {
  state: () => ({
    currentDag: null as any | null,
    dagList: [] as any[],
    activeDag: null as any | null,
    historyVersions: [] as any[],
    loading: false,
    error: null as string | null,
  }),

  actions: {
    async fetchDagsByPipeline(pipelineId: number) {
      this.loading = true;
      try {
        this.dagList = await dagApi.getByPipelineId(pipelineId);
        this.error = null;
      } catch (err) {
        this.error = "获取DAG列表失败";
        console.error(err);
      } finally {
        this.loading = false;
      }
    },

    async fetchActiveDag(pipelineId: number) {
      this.loading = true;
      try {
        this.activeDag = await dagApi.getActive(pipelineId);
        this.error = null;
      } catch (err) {
        this.error = "获取活动DAG失败";
        console.error(err);
      } finally {
        this.loading = false;
      }
    },

    async createDag(data: any) {
      this.loading = true;
      try {
        const newDag = await dagApi.create(data);
        this.dagList.push(newDag);
        this.error = null;
        return newDag;
      } catch (err) {
        this.error = "创建DAG失败";
        console.error(err);
        throw err;
      } finally {
        this.loading = false;
      }
    },

    // 添加新的action方法
    async fetchDagDetail(id: number) {
      this.loading = true;
      try {
        this.currentDag = await dagApi.getById(id);
        this.error = null;
        return this.currentDag;
      } catch (err) {
        this.error = "获取DAG详情失败";
        console.error(err);
        throw err;
      } finally {
        this.loading = false;
      }
    },
    // 获取所有DAG的分页方法
    async fetchAllDags(page = 1, pageSize = 10) {
      this.loading = true;
      try {
        const response = await dagApi.getAllDags({ page, pageSize });
        return response.data;
      } catch (error) {
        this.error = "获取DAG列表失败";
        console.error("Failed to fetch DAGs:", error);
      } finally {
        this.loading = false;
      }
    },

    async updateDag(id: number, data: any) {
      this.loading = true;
      try {
        const updatedDag = await dagApi.update(id, data);
        const index = this.dagList.findIndex((d) => d.id === id);
        if (index !== -1) {
          this.dagList[index] = updatedDag;
        }
        if (this.currentDag && this.currentDag.id === id) {
          this.currentDag = updatedDag;
        }
        this.error = null;
        return updatedDag;
      } catch (err) {
        this.error = "更新DAG失败";
        console.error(err);
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async deleteDag(id: number) {
      this.loading = true;
      try {
        await dagApi.delete(id);
        this.dagList = this.dagList.filter((d) => d.id !== id);
        if (this.currentDag && this.currentDag.id === id) {
          this.currentDag = null;
        }
        this.error = null;
      } catch (err) {
        this.error = "删除DAG失败";
        console.error(err);
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async activateDag(id: number) {
      this.loading = true;
      try {
        await dagApi.activate(id);
        // 更新活动DAG状态
        this.activeDag = await dagApi.getActive(this.currentDag?.pipelineId);
        this.error = null;
      } catch (err) {
        this.error = "激活DAG失败";
        console.error(err);
        throw err;
      } finally {
        this.loading = false;
      }
    },

    async fetchDagHistory(pipelineId: number) {
      this.loading = true;
      try {
        this.historyVersions = await dagApi.getHistory(pipelineId);
        this.error = null;
      } catch (err) {
        this.error = "获取DAG历史版本失败";
        console.error(err);
      } finally {
        this.loading = false;
      }
    },

    async validateDag(nodes: any[]) {
      try {
        return await dagApi.validate(nodes);
      } catch (err) {
        this.error = "DAG验证失败";
        console.error(err);
        throw err;
      }
    },
  },
});
