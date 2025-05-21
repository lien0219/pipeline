<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>构建历史</h2>
        <p>查看所有流水线的构建历史记录</p>
      </div>
    </div>

    <el-card>
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="流水线">
            <el-select v-model="filterForm.pipeline_id" placeholder="选择流水线" clearable>
              <el-option
                  v-for="pipeline in pipelines"
                  :key="pipeline.id"
                  :label="pipeline.name"
                  :value="pipeline.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="状态">
            <el-select v-model="filterForm.status" placeholder="全部状态" clearable>
              <el-option label="成功" value="success" />
              <el-option label="运行中" value="running" />
              <el-option label="失败" value="failed" />
              <el-option label="等待中" value="pending" />
              <el-option label="已取消" value="canceled" />
            </el-select>
          </el-form-item>

          <el-form-item label="时间范围">
            <el-date-picker
                v-model="filterForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleFilter">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetFilter">
              <el-icon><RefreshRight /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
          :data="buildRuns"
          style="width: 100%"
          v-loading="loading"
          @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="ID" width="80" sortable="custom" />

        <el-table-column prop="pipeline.name" label="流水线" min-width="150">
          <template #default="{ row }">
            <router-link :to="`/pipelines/${row.pipeline_id}`" class="pipeline-link">
              {{ row.pipeline?.name }}
            </router-link>
          </template>
        </el-table-column>

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

        <el-table-column prop="user.name" label="触发者" width="120">
          <template #default="{ row }">
            {{ row.user?.name || row.user?.username || '-' }}
          </template>
        </el-table-column>

        <el-table-column prop="start_time" label="开始时间" width="180" sortable="custom">
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
                @click="viewBuildRun(row.pipeline_id, row.id)"
            >
              详情
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="viewLogs(row.pipeline_id, row.id)"
            >
              日志
            </el-button>

            <el-button
                link
                type="danger"
                size="small"
                v-if="row.status === 'running' || row.status === 'pending'"
                @click="cancelBuildRun(row.pipeline_id, row.id)"
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

    <!-- 日志对话框 -->
    <el-dialog
        v-model="logsDialogVisible"
        title="构建日志"
        width="80%"
        :before-close="handleCloseLogsDialog"
    >
      <div v-loading="logsLoading">
        <div v-if="currentRun" class="run-info">
          <el-descriptions :column="3" border size="small">
            <el-descriptions-item label="流水线">{{ currentRun.pipeline?.name }}</el-descriptions-item>
            <el-descriptions-item label="运行ID">{{ currentRun.id }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(currentRun.status)">
                {{ getStatusText(currentRun.status) }}
              </el-tag>
            </el-descriptions-item>
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
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { usePipelineStore } from '@/stores/pipeline';
import { Search, RefreshRight } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

const router = useRouter();
const pipelineStore = usePipelineStore();
const loading = ref(false);
const logsLoading = ref(false);
const logsDialogVisible = ref(false);
const logs = ref('');
const currentRun = ref(null);

const buildRuns = ref([]);
const pipelines = ref([]);

// 筛选表单
const filterForm = reactive({
  pipeline_id: '',
  status: '',
  dateRange: []
});

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
});

// 排序
const sortParams = reactive({
  prop: 'start_time',
  order: 'descending'
});

// 获取构建历史
const fetchBuildHistory = async () => {
  loading.value = true;

  try {
    const params = {
      page: pagination.currentPage,
      limit: pagination.pageSize,
      sort_by: sortParams.prop,
      sort_order: sortParams.order === 'ascending' ? 'asc' : 'desc',
      pipeline_id: filterForm.pipeline_id || undefined,
      status: filterForm.status || undefined,
      start_date: filterForm.dateRange && filterForm.dateRange[0] ? filterForm.dateRange[0] : undefined,
      end_date: filterForm.dateRange && filterForm.dateRange[1] ? filterForm.dateRange[1] : undefined
    };

    const response = await pipelineStore.getBuildHistory(params);
    buildRuns.value = response.data || [];
    pagination.total = response.total || 0;
  } catch (error) {
    console.error('Failed to fetch build history:', error);
    ElMessage.error('获取构建历史失败');
  } finally {
    loading.value = false;
  }
};

// 获取流水线列表（用于筛选）
const fetchPipelines = async () => {
  try {
    const response = await pipelineStore.fetchPipelines({ page:1,pageSize:10 });
    pipelines.value = response.data || [];
  } catch (error) {
    console.error('Failed to fetch pipelines:', error);
  }
};

// 筛选
const handleFilter = () => {
  pagination.currentPage = 1;
  fetchBuildHistory();
};

// 重置筛选
const resetFilter = () => {
  filterForm.pipeline_id = '';
  filterForm.status = '';
  filterForm.dateRange = [];
  pagination.currentPage = 1;
  fetchBuildHistory();
};

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  fetchBuildHistory();
};

const handleCurrentChange = (page) => {
  pagination.currentPage = page;
  fetchBuildHistory();
};

// 排序处理
const handleSortChange = ({ prop, order }) => {
  if (prop) {
    sortParams.prop = prop;
    sortParams.order = order;
  } else {
    sortParams.prop = 'start_time';
    sortParams.order = 'descending';
  }
  fetchBuildHistory();
};

// 查看构建详情
const viewBuildRun = (pipelineId, runId) => {
  router.push(`/pipelines/${pipelineId}/runs/${runId}`);
};

// 查看日志
const viewLogs = async (pipelineId, runId) => {
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

// 取消构建
const cancelBuildRun = async (pipelineId, runId) => {
  try {
    await ElMessageBox.confirm('确定要取消此构建吗？', '取消确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    await pipelineStore.cancelPipelineRun(pipelineId, runId);
    ElMessage.success('构建已取消');
    fetchBuildHistory();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to cancel build run:', error);
      ElMessage.error('取消构建失败');
    }
  }
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

onMounted(() => {
  fetchBuildHistory();
  fetchPipelines();
});
</script>

<style scoped>
.page-header {
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

.filter-container {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}

.pipeline-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.pipeline-link:hover {
  text-decoration: underline;
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
