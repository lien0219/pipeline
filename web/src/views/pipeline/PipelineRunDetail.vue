<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>{{ pipelineName }} - 运行详情 #{{ runId }}</h2>
        <p>
          状态:
          <el-tag :type="getStatusType(status)">{{
            getStatusText(status)
          }}</el-tag>
        </p>
      </div>
      <div class="header-actions">
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回列表
        </el-button>
        <el-button
          v-if="['running', 'pending'].includes(status)"
          type="danger"
          @click="cancelRun"
        >
          <el-icon><Close /></el-icon>
          取消运行
        </el-button>
      </div>
    </div>

    <el-card v-loading="loading">
      <el-descriptions :column="3" border class="run-info">
        <el-descriptions-item label="运行ID">{{ runId }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{
          formatDate(createdAt)
        }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{
          formatDate(updatedAt)
        }}</el-descriptions-item>
        <el-descriptions-item label="Git仓库">{{
          gitRepo
        }}</el-descriptions-item>
        <el-descriptions-item label="Git分支">{{
          gitBranch
        }}</el-descriptions-item>
        <el-descriptions-item label="最后运行时间">{{
          formatDate(lastRunAt)
        }}</el-descriptions-item>
        <el-descriptions-item label="创建者">{{
          creatorName
        }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{
          getStatusText(status)
        }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{
          description || "-"
        }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="mt-4">
      <template #header>
        <div class="card-header">
          <span>执行流程</span>
        </div>
      </template>

      <div class="execution-flow">
        <div
          v-for="(stage, index) in stages"
          :key="stage.id"
          class="stage-item"
        >
          <div class="stage-header">
            <div class="stage-icon">
              <el-icon v-if="stage.status === 'success'"><Check /></el-icon>
              <el-icon v-else-if="stage.status === 'failed'"><Close /></el-icon>
              <el-icon v-else-if="stage.status === 'running'"
                ><Loading
              /></el-icon>
              <el-icon v-else><Clock /></el-icon>
            </div>
            <div class="stage-info">
              <h4>
                {{ stage.name }} <small>#{{ stage.order }}</small>
              </h4>
              <p class="stage-time">
                {{ formatDate(stage.start_time) }} -
                {{ stage.end_time ? formatDate(stage.end_time) : "进行中" }}
                <span v-if="stage.duration"
                  >({{ formatDuration(stage.duration) }})</span
                >
              </p>
            </div>
          </div>

          <div class="stage-jobs">
            <div
              v-for="job in stage.jobs"
              :key="job.id"
              class="job-item"
              :class="{
                'job-success': job.status === 'success',
                'job-failed': job.status === 'failed',
              }"
            >
              <div class="job-header">
                <div class="job-icon">
                  <el-icon v-if="job.status === 'success'"><Check /></el-icon>
                  <el-icon v-else-if="job.status === 'failed'"
                    ><Close
                  /></el-icon>
                  <el-icon v-else-if="job.status === 'running'"
                    ><Loading
                  /></el-icon>
                  <el-icon v-else><Clock /></el-icon>
                </div>
                <div class="job-info">
                  <h5>{{ job.name }}</h5>
                  <p class="job-time">
                    {{ formatDate(job.start_time) }} -
                    {{ job.end_time ? formatDate(job.end_time) : "进行中" }}
                    <span v-if="job.duration"
                      >({{ formatDuration(job.duration) }})</span
                    >
                  </p>
                  <div v-if="job.command" class="job-command">
                    <strong>命令:</strong> {{ job.command }}
                  </div>
                  <div v-if="job.image" class="job-image">
                    <strong>镜像:</strong> {{ job.image }}
                  </div>
                </div>
                <el-tag :type="getStatusType(job.status)" size="small">{{
                  getStatusText(job.status)
                }}</el-tag>
              </div>

              <div class="job-logs" v-if="expandedJobId === job.id">
                <pre class="logs-content">{{ job.logs || "暂无日志" }}</pre>
              </div>

              <el-button
                v-if="job.logs"
                type="text"
                size="small"
                class="view-logs-btn"
                @click="toggleJobLogs(job.id)"
              >
                {{ expandedJobId === job.id ? "收起日志" : "查看日志" }}
                <el-icon v-if="expandedJobId === job.id"><ArrowUp /></el-icon>
                <el-icon v-else><ArrowDown /></el-icon>
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <el-card class="mt-4" v-if="artifacts.length > 0">
      <template #header>
        <div class="card-header">
          <span>制品列表</span>
        </div>
      </template>

      <el-table :data="artifacts" style="width: 100%">
        <el-table-column prop="name" label="名称" min-width="200" />
        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              size="small"
              @click="downloadArtifact(row.id)"
            >
              <el-icon><Download /></el-icon> 下载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="mt-4" v-if="config">
      <template #header>
        <div class="card-header">
          <span>配置信息</span>
        </div>
      </template>
      <pre class="config-content">{{ config }}</pre>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, reactive } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  ArrowLeft,
  Check,
  Close,
  Clock,
  Download,
  Loading,
  ArrowUp,
  ArrowDown,
} from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { pipelineApi } from "@/api/pipeline";
import { usePipelineStore } from "@/stores/pipeline";

const route = useRoute();
const router = useRouter();
const pipelineStore = usePipelineStore();

// 路由参数
const pipelineId = route.params.id as string;
const runId = route.params.runId as string;

// 状态管理
const loading = ref(true);
const pipelineName = ref("");
const gitBranch = ref("");
const stages = ref([]);
const artifacts = ref([]);
const expandedJobId = ref("");
const createdAt = ref("");
const updatedAt = ref("");
const description = ref("");
const gitRepo = ref("");
const lastRunAt = ref("");
const creatorName = ref("");
const config = ref("");
const status = ref("");

// 返回列表
const goBack = () => {
  router.push(`/pipelines/${pipelineId}`);
};

// 取消运行
const cancelRun = async () => {
  try {
    await ElMessageBox.confirm("确定要取消此流水线运行吗？", "取消确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await pipelineStore.cancelPipelineRun(pipelineId, runId);
    ElMessage.success("流水线运行已取消");
    fetchRunDetail();
  } catch (error) {
    if (error !== "cancel") {
      console.error("取消流水线运行失败:", error);
      ElMessage.error("取消流水线运行失败");
    }
  }
};

// 切换作业日志显示
const toggleJobLogs = (jobId: string) => {
  if (expandedJobId.value === jobId) {
    expandedJobId.value = "";
  } else {
    expandedJobId.value = jobId;
  }
};

// 下载制品
const downloadArtifact = (artifactId: string) => {
  window.open(`/api/artifacts/${artifactId}/download`, "_blank");
};

// 获取运行详情
const fetchRunDetail = async () => {
  loading.value = true;
  try {
    // 获取流水线基本信息
    const pipelineResponse = await pipelineStore.fetchPipelineById(pipelineId);
    const pipelineData = pipelineResponse.data;
    pipelineName.value = pipelineData.name;
    createdAt.value = pipelineData.created_at;
    updatedAt.value = pipelineData.updated_at;
    description.value = pipelineData.description;
    gitRepo.value = pipelineData.git_repo;
    gitBranch.value = pipelineData.git_branch;
    status.value = pipelineData.status;
    lastRunAt.value = pipelineData.last_run_at;
    creatorName.value =
      pipelineData.creator?.name || pipelineData.creator?.username || "-";
    config.value = pipelineData.Config;
    stages.value = pipelineData.stages || [];

    // 获取制品列表
    const artifactsResponse = await pipelineApi.getArtifacts({
      pipeline_id: pipelineId,
      pipeline_run_id: runId,
    });
    artifacts.value = artifactsResponse.data?.list || [];

    // 如果运行中，设置定时刷新
    if (["running", "pending"].includes(status.value)) {
      const refreshInterval = setInterval(() => {
        fetchRunDetail();
      }, 5000);

      // 组件卸载时清除定时器
      onUnmounted(() => clearInterval(refreshInterval));
    }
  } catch (error) {
    console.error("获取流水线运行详情失败:", error);
    ElMessage.error("获取流水线运行详情失败");
  } finally {
    loading.value = false;
  }
};

// 格式化状态类型
const getStatusType = (status: string) => {
  switch (status) {
    case "success":
      return "success";
    case "running":
      return "primary";
    case "failed":
      return "danger";
    case "pending":
      return "info";
    case "canceled":
      return "warning";
    default:
      return "default";
  }
};

// 格式化状态文本
const getStatusText = (status: string) => {
  switch (status) {
    case "success":
      return "成功";
    case "running":
      return "运行中";
    case "failed":
      return "失败";
    case "pending":
      return "等待中";
    case "canceled":
      return "已取消";
    default:
      return "未知";
  }
};

// 格式化日期
const formatDate = (date: string) => {
  if (!date) return "-";
  return dayjs(date).format("YYYY-MM-DD HH:mm:ss");
};

// 格式化持续时间
const formatDuration = (seconds: number) => {
  if (!seconds || seconds < 0) return "0s";

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;

  if (minutes === 0) {
    return `${remainingSeconds}s`;
  }

  return `${minutes}m ${remainingSeconds}s`;
};

// 格式化文件大小
const formatFileSize = (bytes: number) => {
  if (bytes === 0) return "0 B";

  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

// 页面加载时获取数据
onMounted(() => {
  fetchRunDetail();
});
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-title h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.header-title p {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.run-info {
  margin-bottom: 20px;
}

.execution-flow {
  padding: 10px 0;
}

.stage-item {
  margin-bottom: 30px;
  position: relative;
}

.stage-item:not(:last-child):after {
  content: "";
  position: absolute;
  top: 30px;
  left: 14px;
  bottom: -30px;
  width: 2px;
  background-color: #e5e7eb;
}

.stage-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.stage-icon {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background-color: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stage-info {
  flex-grow: 1;
}

.stage-info h4 {
  margin: 0 0 4px 0;
  font-size: 16px;
}

.stage-time {
  margin: 0;
  font-size: 12px;
  color: #606266;
}

.stage-jobs {
  margin-left: 45px;
  border-left: 2px solid #e5e7eb;
  padding-left: 25px;
}

.job-item {
  margin-bottom: 20px;
  padding: 15px;
  border-radius: 4px;
  background-color: #f9fafb;
  transition: all 0.3s;
}

.job-item.job-success {
  border-left: 3px solid #67c23a;
}

.job-item.job-failed {
  border-left: 3px solid #f56c6c;
}

.job-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.job-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background-color: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.job-info {
  flex-grow: 1;
}

.job-info h5 {
  margin: 0 0 4px 0;
  font-size: 14px;
}

.job-time {
  margin: 0;
  font-size: 12px;
  color: #606266;
}

.job-logs {
  margin-top: 10px;
  background-color: #1e1e1e;
  color: #f8f8f8;
  border-radius: 4px;
  padding: 15px;
  max-height: 300px;
  overflow-y: auto;
}

.logs-content {
  margin: 0;
  font-family: monospace;
  font-size: 13px;
  white-space: pre-wrap;
}

.view-logs-btn {
  padding: 0;
  height: auto;
  font-size: 12px;
  color: #409eff;
}

.mt-4 {
  margin-top: 20px;
}
.config-content {
  margin: 0;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
  font-family: monospace;
  font-size: 13px;
  white-space: pre-wrap;
  max-height: 400px;
  overflow-y: auto;
}
</style>
