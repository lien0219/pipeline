<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>{{ pipeline.name }}</h2>
        <p>{{ pipeline.description || '暂无描述' }}</p>
      </div>

      <div class="header-actions">
        <el-button
            type="primary"
            @click="triggerPipeline"
            :disabled="loading || pipeline.status === 'running'"
        >
          <el-icon><VideoPlay /></el-icon>
          运行流水线
        </el-button>

        <el-button @click="$router.push(`/pipelines/${pipelineId}/edit`)">
          <el-icon><Edit /></el-icon>
          编辑
        </el-button>
      </div>
    </div>

    <el-tabs v-model="activeTab" class="pipeline-tabs">
      <el-tab-pane label="概览" name="overview">
        <el-card v-loading="loading">
          <template #header>
            <div class="card-header">
              <span>流水线信息</span>
            </div>
          </template>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="ID">{{ pipeline.id }}</el-descriptions-item>
            <el-descriptions-item label="创建者">{{ pipeline.creator?.name || pipeline.creator?.username }}</el-descriptions-item>
            <el-descriptions-item label="Git 仓库">{{ pipeline.git_repo }}</el-descriptions-item>
            <el-descriptions-item label="Git 分支">{{ pipeline.git_branch }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(pipeline.status)">
                {{ getStatusText(pipeline.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="最近执行">
              {{ pipeline.last_run_at ? formatDate(pipeline.last_run_at) : '从未执行' }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatDate(pipeline.created_at) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间">{{ formatDate(pipeline.updated_at) }}</el-descriptions-item>
          </el-descriptions>

          <div class="pipeline-stages" v-if="pipeline.stages && pipeline.stages.length > 0">
            <h3>流水线阶段</h3>

            <div class="stage-flow">
              <div
                  v-for="(stage, index) in pipeline.stages"
                  :key="stage.id"
                  class="stage-node"
              >
                <div class="stage-card">
                  <div class="stage-header">
                    <span class="stage-name">{{ stage.name }}</span>
                    <span class="stage-order">#{{ stage.order }}</span>
                  </div>

                  <div class="stage-jobs">
                    <div
                        v-for="job in stage.jobs"
                        :key="job.id"
                        class="job-item"
                    >
                      <el-tooltip :content="job.command" placement="top">
                        <div class="job-name">
                          <el-icon><Terminal /></el-icon>
                          {{ job.name }}
                        </div>
                      </el-tooltip>
                    </div>
                  </div>
                </div>

                <div class="stage-arrow" v-if="index < pipeline.stages.length - 1">
                  <el-icon><ArrowRight /></el-icon>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="运行历史" name="history">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>运行历史</span>
              <el-button link type="primary" @click="fetchPipelineRuns">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>

          <el-table :data="pipelineRuns" style="width: 100%" v-loading="runsLoading">
            <el-table-column prop="id" label="运行ID" width="80" />

            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>

            <el-table-column prop="git_branch" label="分支" width="120" />

            <el-table-column prop="git_commit" label="提交" width="120">
              <template #default="{ row }">
                {{ row.git_commit ? row.git_commit.substring(0, 8) : '-' }}
              </template>
            </el-table-column>

            <el-table-column prop="duration" label="耗时" width="120">
              <template #default="{ row }">
                {{ formatDuration(row.duration) }}
              </template>
            </el-table-column>

            <el-table-column prop="trigger_by" label="触发者" width="120">
              <template #default="{ row }">
                {{ row.user?.name || row.user?.username || '-' }}
              </template>
            </el-table-column>

            <el-table-column prop="start_time" label="开始时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.start_time) }}
              </template>
            </el-table-column>

            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button
                    link
                    type="primary"
                    size="small"
                    @click="viewPipelineRun(row.id)"
                >
                  详情
                </el-button>

                <el-button
                    link
                    type="primary"
                    size="small"
                    @click="viewLogs(row.id)"
                >
                  日志
                </el-button>

                <el-button
                    link
                    type="danger"
                    size="small"
                    v-if="row.status === 'running' || row.status === 'pending'"
                    @click="cancelPipelineRun(row.id)"
                >
                  取消
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-container">
            <el-pagination
                v-model:current-page="pagination.currentPage"
                v-model:page-size="pagination.pageSize"
                :page-sizes="[10, 20, 30, 50]"
                :total="pagination.total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
            />
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="制品" name="artifacts">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>制品列表</span>
              <el-button link type="primary" @click="fetchArtifacts">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>

          <el-table :data="artifacts" style="width: 100%" v-loading="artifactsLoading">
            <el-table-column prop="name" label="名称" min-width="180" />

            <el-table-column prop="size" label="大小" width="120">
              <template #default="{ row }">
                {{ formatFileSize(row.size) }}
              </template>
            </el-table-column>

            <el-table-column prop="pipeline_run_id" label="运行ID" width="100" />

            <el-table-column prop="download_count" label="下载次数" width="100" />

            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">
                {{ formatDate(row.created_at) }}
              </template>
            </el-table-column>

            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button
                    link
                    type="primary"
                    size="small"
                    @click="downloadArtifact(row.id)"
                >
                  下载
                </el-button>

                <el-button
                    link
                    type="danger"
                    size="small"
                    @click="deleteArtifact(row.id)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination-container">
            <el-pagination
                v-model:current-page="artifactsPagination.currentPage"
                v-model:page-size="artifactsPagination.pageSize"
                :page-sizes="[10, 20, 30, 50]"
                :total="artifactsPagination.total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleArtifactsSizeChange"
                @current-change="handleArtifactsCurrentChange"
            />
          </div>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 日志对话框 -->
    <el-dialog
        v-model="logsDialogVisible"
        title="流水线运行日志"
        width="80%"
        :before-close="handleCloseLogsDialog"
    >
      <div v-loading="logsLoading">
        <div v-if="currentRun" class="run-info">
          <el-descriptions :column="3" border size="small">
            <el-descriptions-item label="运行ID">{{ currentRun.id }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(currentRun.status)">
                {{ getStatusText(currentRun.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="耗时">{{ formatDuration(currentRun.duration) }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="logs-container">
          <pre class="logs-content">{{ logs }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { usePipelineStore } from '@/stores/pipeline';
import { VideoPlay, Edit, Terminal, ArrowRight, Refresh } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

const route = useRoute();
const router = useRouter();
const pipelineStore = usePipelineStore();

const pipelineId = route.params.id;
const activeTab = ref('overview');
const loading = ref(false);
const runsLoading = ref(false);
const artifactsLoading = ref(false);
const logsLoading = ref(false);
const logsDialogVisible = ref(false);
const logs = ref('');
const currentRun = ref(null);

const pipeline = ref({
  id: '',
  name: '',
  description: '',
  git_repo: '',
  git_branch: '',
  creator: {},
  stages: [],
  status: '',
  last_run_at: null,
  created_at: null,
  updated_at: null
});

const pipelineRuns = ref([]);
const artifacts = ref([]);

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
});

const artifactsPagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
});

// 获取流水线详情
const fetchPipelineDetail = async () => {
  loading.value = true;
  try {
    const response = await pipelineStore.fetchPipelineById(pipelineId);
    pipeline.value = response.data;
  } catch (error) {
    console.error('Failed to fetch pipeline details:', error);
    ElMessage.error('获取流水线详情失败');
  } finally {
    loading.value = false;
  }
};

// 获取流水线运行记录
const fetchPipelineRuns = async () => {
  runsLoading.value = true;
  try {
    const params = {
      page: pagination.currentPage,
      limit: pagination.pageSize
    };

    const response = await pipelineStore.getPipelineRuns(pipelineId, params);
    pipelineRuns.value = response.data || [];
    pagination.total = response.total || 0;
  } catch (error) {
    console.error('Failed to fetch pipeline runs:', error);
    ElMessage.error('获取流水线运行记录失败');
  } finally {
    runsLoading.value = false;
  }
};

// 获取制品列表
const fetchArtifacts = async () => {
  artifactsLoading.value = true;
  try {
    const params = {
      page: artifactsPagination.currentPage,
      limit: artifactsPagination.pageSize,
      pipeline_id: pipelineId
    };

    const response = await pipelineStore.getArtifacts(params);
    artifacts.value = response.data || [];
    artifactsPagination.total = response.total || 0;
  } catch (error) {
    console.error('Failed to fetch artifacts:', error);
    ElMessage.error('获取制品列表失败');
  } finally {
    artifactsLoading.value = false;
  }
};

// 触发流水线
const triggerPipeline = async () => {
  try {
    await pipelineStore.triggerPipeline(pipelineId);
    ElMessage.success('流水线已触发');
    fetchPipelineDetail();
    if (activeTab.value === 'history') {
      fetchPipelineRuns();
    }
  } catch (error) {
    console.error('Failed to trigger pipeline:', error);
    ElMessage.error('触发流水线失败');
  }
};

// 查看流水线运行详情
const viewPipelineRun = (runId) => {
  router.push(`/pipelines/${pipelineId}/runs/${runId}`);
};

// 查看日志
const viewLogs = async (runId) => {
  logsLoading.value = true;
  logsDialogVisible.value = true;
  logs.value = '加载中...';

  try {
    // 获取运行详情
    const runResponse = await pipelineStore.getPipelineRunById(pipelineId, runId);
    currentRun.value = runResponse.data;

    // 获取日志
    const logsResponse = await pipelineStore.getPipelineRunLogs(pipelineId, runId);

    if (logsResponse.data.logs) {
      // 格式化日志
      let formattedLogs = '';
      for (const stageId in logsResponse.data.logs) {
        const stageLogs = logsResponse.data.logs[stageId];
        for (const jobId in stageLogs) {
          formattedLogs += `===== Job ${jobId} =====\n\n${stageLogs[jobId]}\n\n`;
        }
      }
      logs.value = formattedLogs || '暂无日志';
    } else {
      logs.value = '暂无日志';
    }
  } catch (error) {
    console.error('Failed to fetch logs:', error);
    logs.value = '获取日志失败';
  } finally {
    logsLoading.value = false;
  }
};

// 关闭日志对话框
const handleCloseLogsDialog = () => {
  logsDialogVisible.value = false;
  logs.value = '';
  currentRun.value = null;
};

// 取消流水线运行
const cancelPipelineRun = async (runId) => {
  try {
    await ElMessageBox.confirm('确定要取消此流水线运行吗？', '取消确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    await pipelineStore.cancelPipelineRun(pipelineId, runId);
    ElMessage.success('流水线运行已取消');
    fetchPipelineRuns();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to cancel pipeline run:', error);
      ElMessage.error('取消流水线运行失败');
    }
  }
};

// 下载制品
const downloadArtifact = (artifactId) => {
  window.open(`/api/artifacts/${artifactId}/download`, '_blank');
};

// 删除制品
const deleteArtifact = async (artifactId) => {
  try {
    await ElMessageBox.confirm('确定要删除此制品吗？此操作不可恢复。', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    await pipelineStore.deleteArtifact(artifactId);
    ElMessage.success('制品已删除');
    fetchArtifacts();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete artifact:', error);
      ElMessage.error('删除制品失败');
    }
  }
};

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  fetchPipelineRuns();
};

const handleCurrentChange = (page) => {
  pagination.currentPage = page;
  fetchPipelineRuns();
};

const handleArtifactsSizeChange = (size) => {
  artifactsPagination.pageSize = size;
  fetchArtifacts();
};

const handleArtifactsCurrentChange = (page) => {
  artifactsPagination.currentPage = page;
  fetchArtifacts();
};

// 格式化状态
const getStatusType = (status) => {
  switch (status) {
    case 'success': return 'success';
    case 'running': return 'primary';
    case 'failed': return 'danger';
    case 'pending': return 'info';
    case 'canceled': return 'warning';
    default: return 'info';
  }
};

const getStatusText = (status) => {
  switch (status) {
    case 'success': return '成功';
    case 'running': return '运行中';
    case 'failed': return '失败';
    case 'pending': return '等待中';
    case 'canceled': return '已取消';
    default: return '未知';
  }
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

// 格式化持续时间
const formatDuration = (seconds) => {
  if (!seconds) return '0s';

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;

  if (minutes === 0) {
    return `${remainingSeconds}s`;
  }

  return `${minutes}m ${remainingSeconds}s`;
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B';

  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 监听标签页变化
watch(activeTab, (newValue) => {
  if (newValue === 'history') {
    fetchPipelineRuns();
  } else if (newValue === 'artifacts') {
    fetchArtifacts();
  }
});

onMounted(() => {
  fetchPipelineDetail();
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
  color: #606266;
}

.pipeline-tabs {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pipeline-stages {
  margin-top: 30px;
}

.pipeline-stages h3 {
  margin-bottom: 20px;
  font-size: 18px;
  font-weight: 600;
}

.stage-flow {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
}

.stage-node {
  display: flex;
  align-items: center;
}

.stage-card {
  background-color: #f5f7fa;
  border-radius: 4px;
  padding: 15px;
  min-width: 200px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.stage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.stage-name {
  font-weight: bold;
}

.stage-order {
  color: #909399;
  font-size: 12px;
}

.stage-jobs {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.job-item {
  background-color: #fff;
  border-radius: 4px;
  padding: 8px 12px;
  border-left: 3px solid var(--el-color-primary);
}

.job-name {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 14px;
}

.stage-arrow {
  margin: 0 10px;
  color: #909399;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.run-info {
  margin-bottom: 20px;
}

.logs-container {
  background-color: #1e1e1e;
  color: #f8f8f8;
  border-radius: 4px;
  padding: 15px;
  height: 500px;
  overflow-y: auto;
}

.logs-content {
  font-family: 'Courier New', Courier, monospace;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}
</style>
