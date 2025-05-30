import { defineStore } from "pinia";
import { ElMessage } from "element-plus";
import * as hpaApi from "@/api/hpa";

export const useHPAStore = defineStore("hpa", {
  state: () => ({
    policies: [],
    loading: false,
  }),

  actions: {
    async fetchHPAPolicies(params?: any) {
      this.loading = true;
      try {
        const response = await hpaApi.getHPAPolicies(params);
        this.policies = response.data;
      } catch (error) {
        ElMessage.error("获取HPA策略失败");
        throw error;
      } finally {
        this.loading = false;
      }
    },
  },
});
