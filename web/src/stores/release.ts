import { defineStore } from "pinia";
import {
  releaseApi,
  type Release,
  type CreateReleaseRequest,
} from "@/api/release";
import { ElMessage } from "element-plus";
import { artifactApi } from "../../src/api/artifact";

interface ReleaseState {
  releases: Release[];
  artifacts: any[];
  currentRelease: Release | null;
  total: number;
  loading: boolean;
  submitting: boolean;
  page: number;
  pageSize: number;
  filters: {
    environment: string;
    status: string;
    dateRange: [string, string] | [];
  };
}

export const useReleaseStore = defineStore("release", {
  state: (): ReleaseState => ({
    releases: [],
    artifacts: [],
    currentRelease: null,
    total: 0,
    loading: false,
    submitting: false,
    page: 1,
    pageSize: 10,
    filters: {
      environment: "",
      status: "",
      dateRange: [],
    },
  }),

  actions: {
    // 获取发布列表
    async fetchReleases() {
      this.loading = true;
      try {
        const [startDate, endDate] = this.filters.dateRange;
        const response: any = await releaseApi.getReleases(
          this.page,
          this.pageSize,
          this.filters.environment,
          this.filters.status,
          startDate || "",
          endDate || ""
        );
        this.releases = response.data.list;
        this.total = response.data.total;
      } catch (error) {
        ElMessage.error("获取发布列表失败");
        console.error("Failed to fetch releases:", error);
      } finally {
        this.loading = false;
      }
    },

    // 获取制品列表
    async fetchArtifacts() {
      try {
        const response: any = await artifactApi.getArtifacts();
        this.artifacts = response.data.list;
      } catch (error) {
        ElMessage.error("获取制品列表失败");
        console.error("Failed to fetch artifacts:", error);
      }
    },

    // 创建发布
    async createRelease(data: CreateReleaseRequest) {
      this.submitting = true;
      try {
        const response = await releaseApi.createRelease(data);
        ElMessage.success("发布创建成功");
        this.fetchReleases();
        return response.data.data;
      } catch (error) {
        ElMessage.error("创建发布失败");
        console.error("Failed to create release:", error);
        throw error;
      } finally {
        this.submitting = false;
      }
    },

    // 删除发布
    async deleteRelease(id: number) {
      try {
        await releaseApi.deleteRelease(id);
        ElMessage.success("发布记录已删除");
        this.fetchReleases();
      } catch (error) {
        ElMessage.error("删除发布失败");
        console.error(`Failed to delete release ${id}`, error);
        throw error;
      }
    },

    // 回滚发布
    async rollbackRelease(id: number) {
      this.submitting = true;
      try {
        const response = await releaseApi.rollbackRelease(id);
        ElMessage.success("发布回滚成功");
        this.fetchReleases();
        return response.data.data;
      } catch (error) {
        ElMessage.error("回滚发布失败");
        console.error(`Failed to rollback release ${id}`, error);
        throw error;
      } finally {
        this.submitting = false;
      }
    },

    // 设置筛选条件
    setFilters(filters: Partial<ReleaseState["filters"]>) {
      this.filters = { ...this.filters, ...filters };
      this.page = 1;
      this.fetchReleases();
    },

    // 设置分页参数
    setPageParams(page: number, pageSize: number) {
      this.page = page;
      this.pageSize = pageSize;
      this.fetchReleases();
    },
  },
});
