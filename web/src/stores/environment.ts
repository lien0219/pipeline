import { defineStore } from "pinia";
import {
  environmentApi,
  type Environment,
  type CreateEnvironmentRequest,
  type UpdateEnvironmentRequest,
} from "@/api/environment";
import { ElMessage } from "element-plus";

interface EnvironmentState {
  environments: Environment[];
  currentEnvironment: Environment | null | any;
  total: number;
  loading: boolean;
  page: number;
  pageSize: number;
  filters: { name?: string; type?: string; status?: string };
}

export const useEnvironmentStore = defineStore("environment", {
  state: (): EnvironmentState => ({
    environments: [],
    currentEnvironment: null,
    total: 0,
    loading: false,
    page: 1,
    pageSize: 10,
    filters: { name: "", type: "", status: "" },
  }),

  actions: {
    // 获取环境列表
    async fetchEnvironments() {
      this.loading = true;
      try {
        const response: any = await environmentApi.getEnvironments(
          this.page,
          this.pageSize,
          this.filters.name,
          this.filters.type,
          this.filters.status
        );
        this.environments = response.data.list;
        this.total = response.data.total;
      } catch (error) {
        ElMessage.error("获取环境列表失败");
        console.error("Failed to fetch environments:", error);
        this.environments = [];
      } finally {
        this.loading = false;
      }
    },

    // 获取环境详情
    async fetchEnvironmentById(id: number) {
      this.loading = true;
      try {
        const response = await environmentApi.getEnvironmentById(id);
        this.currentEnvironment = response.data;
        return this.currentEnvironment;
      } catch (error) {
        ElMessage.error("获取环境详情失败");
        console.error(`Failed to fetch environment ${id}:`, error);
        return null;
      } finally {
        this.loading = false;
      }
    },

    // 创建环境
    async createEnvironment(data: CreateEnvironmentRequest) {
      this.loading = true;
      try {
        const response = await environmentApi.createEnvironment(data);
        ElMessage.success("环境创建成功");
        this.fetchEnvironments(); // 重新获取列表
        return response.data.data;
      } catch (error) {
        ElMessage.error("环境创建失败");
        console.error("Failed to create environment:", error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // 更新环境
    async updateEnvironment(id: number, data: UpdateEnvironmentRequest) {
      this.loading = true;
      try {
        if (data.variables && !Array.isArray(data.variables)) {
          data.variables = [];
        }
        const response = await environmentApi.updateEnvironment(id, data);
        // ElMessage.success("环境更新成功");
        this.fetchEnvironments();
        return response.data.data;
      } catch (error) {
        ElMessage.error("环境更新失败");
        console.error(`Failed to update environment ${id}:`, error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // 删除环境
    async deleteEnvironment(id: number) {
      this.loading = true;
      try {
        await environmentApi.deleteEnvironment(id);
        ElMessage.success("环境删除成功");
        this.fetchEnvironments(); // 重新获取列表
      } catch (error) {
        ElMessage.error("环境删除失败");
        console.error(`Failed to delete environment ${id}:`, error);
        throw error;
      } finally {
        this.loading = false;
      }
    },

    // 设置筛选条件
    setFilters(filters: { name?: string; type?: string }) {
      this.filters = { ...this.filters, ...filters };
      this.page = 1; // 重置页码
      this.fetchEnvironments();
    },

    // 设置分页参数
    setPageParams(page: number, pageSize: number) {
      this.page = page;
      this.pageSize = pageSize;
      this.fetchEnvironments();
    },
  },
});
