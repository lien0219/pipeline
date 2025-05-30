<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="24" :lg="18" :xl="18">
        <el-card class="overview-card">
          <template #header>
            <div class="card-header">
              <span>流水线概览</span>
              <el-button type="primary" size="small" @click="refreshData">
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>

          <el-row :gutter="20">
            <el-col :span="6">
              <div class="stat-card success">
                <div class="stat-icon">
                  <el-icon :size="24"><Check /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ stats.success }}</div>
                  <div class="stat-label">成功</div>
                </div>
              </div>
            </el-col>

            <el-col :span="6">
              <div class="stat-card running">
                <div class="stat-icon">
                  <el-icon :size="24"><Loading /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ stats.running }}</div>
                  <div class="stat-label">运行中</div>
                </div>
              </div>
            </el-col>

            <el-col :span="6">
              <div class="stat-card failed">
                <div class="stat-icon">
                  <el-icon :size="24"><WarningFilled /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ stats.failed }}</div>
                  <div class="stat-label">失败</div>
                </div>
              </div>
            </el-col>

            <el-col :span="6">
              <div class="stat-card pending">
                <div class="stat-icon">
                  <el-icon :size="24"><Clock /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value">{{ stats.pending }}</div>
                  <div class="stat-label">等待中</div>
                </div>
              </div>
            </el-col>
          </el-row>

          <div class="chart-container">
            <div ref="pipelineChart" style="width: 100%; height: 300px"></div>
          </div>
        </el-card>

        <el-card class="recent-pipelines-card">
          <template #header>
            <div class="card-header">
              <span>最近流水线</span>
              <router-link to="/pipelines">
                <el-button link type="primary">查看全部</el-button>
              </router-link>
            </div>
          </template>

          <el-table
            :data="recentPipelines"
            style="width: 100%"
            v-loading="loading"
          >
            <el-table-column prop="name" label="名称" min-width="180">
              <template #default="{ row }">
                <router-link :to="`/pipelines/${row.id}`" class="pipeline-link">
                  {{ row.name }}
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

            <el-table-column prop="branch" label="分支" width="120" />

            <!-- <el-table-column prop="duration" label="耗时" width="120">
              <template #default="{ row }">
                {{ formatDuration(row.duration) }}
              </template>
            </el-table-column> -->

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
                  @click="triggerPipeline(row.id)"
                  :disabled="row.status === 'running'"
                >
                  运行
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="24" :md="24" :lg="6" :xl="6">
        <el-card class="activity-card">
          <template #header>
            <div class="card-header">
              <span>最近活动</span>
            </div>
          </template>

          <div class="activity-timeline">
            <el-timeline>
              <el-timeline-item
                v-for="(activity, index) in activities"
                :key="index"
                :type="getActivityType(activity.type)"
                :timestamp="formatDate(activity.timestamp)"
                :hollow="activity.hollow"
              >
                {{ activity.content }}
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-card>

        <el-card class="quick-actions-card">
          <template #header>
            <div class="card-header">
              <span>快捷操作</span>
            </div>
          </template>

          <div class="quick-actions">
            <el-button type="primary" @click="$router.push('/pipeline/create')">
              <el-icon><Plus /></el-icon>
              创建流水线
            </el-button>

            <el-button @click="$router.push('/builds/templates')">
              <el-icon><Document /></el-icon>
              构建模板
            </el-button>

            <el-button @click="$router.push('/deploy/environments')">
              <el-icon><SetUp /></el-icon>
              环境配置
            </el-button>

            <el-button @click="$router.push('/settings')">
              <el-icon><Setting /></el-icon>
              系统设置
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";
import { usePipelineStore } from "@/stores/pipeline";
import {
  Check,
  Loading,
  WarningFilled,
  Clock,
  Refresh,
  Plus,
  Document,
  SetUp,
  Setting,
} from "@element-plus/icons-vue";
import * as echarts from "echarts/core";
import { LineChart } from "echarts/charts";
import {
  GridComponent,
  TooltipComponent,
  TitleComponent,
  LegendComponent,
} from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import dayjs from "dayjs";
import { ElMessage } from "element-plus";

echarts.use([
  LineChart,
  GridComponent,
  TooltipComponent,
  TitleComponent,
  LegendComponent,
  CanvasRenderer,
]);

const router = useRouter();
const pipelineStore = usePipelineStore();
const pipelineChart = ref(null);
const chart = ref(null);
const loading = ref(false);

const stats = ref({
  success: 0,
  running: 0,
  failed: 0,
  pending: 0,
});

const recentPipelines = ref([]);
const activities = ref([]);

// 获取数据
const fetchData = async () => {
  loading.value = true;
  try {
    // 获取流水线数据
    const response = await pipelineStore.fetchPipelines({
      page: 1,
      pageSize: 10,
    });
    recentPipelines.value = response.data.list || [];

    // 统计数据
    stats.value = {
      success: 12,
      running: 3,
      failed: 2,
      pending: 5,
    };

    // 模拟活动数据
    activities.value = [
      {
        type: "success",
        content: '流水线 "Frontend Deploy" 构建成功',
        timestamp: new Date(),
        hollow: false,
      },
      {
        type: "warning",
        content: '流水线 "Backend API" 构建失败',
        timestamp: new Date(Date.now() - 3600000),
        hollow: false,
      },
      {
        type: "primary",
        content: '用户 admin 创建了新的流水线 "Database Migration"',
        timestamp: new Date(Date.now() - 7200000),
        hollow: false,
      },
      {
        type: "info",
        content: "系统更新完成",
        timestamp: new Date(Date.now() - 86400000),
        hollow: true,
      },
      {
        type: "success",
        content: '流水线 "Mobile App" 部署成功',
        timestamp: new Date(Date.now() - 172800000),
        hollow: true,
      },
    ];

    initChart();
  } catch (error) {
    console.error("Failed to fetch dashboard data:", error);
  } finally {
    loading.value = false;
  }
};

// 初始化图表
const initChart = () => {
  if (!pipelineChart.value) return;

  if (chart.value) {
    chart.value.dispose();
  }

  chart.value = echarts.init(pipelineChart.value);

  const option = {
    title: {
      text: "流水线执行趋势",
      left: "center",
    },
    tooltip: {
      trigger: "axis",
    },
    legend: {
      data: ["成功", "失败", "总数"],
      bottom: 0,
    },
    grid: {
      left: "3%",
      right: "4%",
      bottom: "10%",
      top: "15%",
      containLabel: true,
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: ["周一", "周二", "周三", "周四", "周五", "周六", "周日"],
    },
    yAxis: {
      type: "value",
    },
    series: [
      {
        name: "成功",
        type: "line",
        data: [5, 7, 6, 9, 8, 7, 10],
        itemStyle: {
          color: "#67C23A",
        },
        areaStyle: {
          color: {
            type: "linear",
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: "rgba(103, 194, 58, 0.3)" },
              { offset: 1, color: "rgba(103, 194, 58, 0.1)" },
            ],
          },
        },
      },
      {
        name: "失败",
        type: "line",
        data: [2, 1, 3, 1, 2, 0, 1],
        itemStyle: {
          color: "#F56C6C",
        },
        areaStyle: {
          color: {
            type: "linear",
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: "rgba(245, 108, 108, 0.3)" },
              { offset: 1, color: "rgba(245, 108, 108, 0.1)" },
            ],
          },
        },
      },
      {
        name: "总数",
        type: "line",
        data: [7, 8, 9, 10, 10, 7, 11],
        itemStyle: {
          color: "#409EFF",
        },
        areaStyle: {
          color: {
            type: "linear",
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: "rgba(64, 158, 255, 0.3)" },
              { offset: 1, color: "rgba(64, 158, 255, 0.1)" },
            ],
          },
        },
      },
    ],
  };

  chart.value.setOption(option);
};

// 刷新数据
const refreshData = () => {
  fetchData();
};

// 触发流水线
const triggerPipeline = async (id) => {
  try {
    await pipelineStore.triggerPipeline(id);
    ElMessage.success("流水线已触发");
    refreshData();
  } catch (error) {
    console.error("Failed to trigger pipeline:", error);
  }
};

// 格式化状态
const getStatusType = (status) => {
  switch (status) {
    case "success":
      return "success";
    case "running":
      return "primary";
    case "failed":
      return "danger";
    case "pending":
      return "info";
    default:
      return "info";
  }
};

const getStatusText = (status) => {
  switch (status) {
    case "success":
      return "成功";
    case "running":
      return "运行中";
    case "failed":
      return "失败";
    case "pending":
      return "等待中";
    default:
      return "未知";
  }
};

// 格式化活动类型
const getActivityType = (type) => {
  switch (type) {
    case "success":
      return "success";
    case "warning":
      return "warning";
    case "primary":
      return "primary";
    default:
      return "info";
  }
};

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format("YYYY-MM-DD HH:mm:ss");
};

// 格式化持续时间
const formatDuration = (seconds) => {
  if (!seconds) return "0s";

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;

  if (minutes === 0) {
    return `${remainingSeconds}s`;
  }

  return `${minutes}m ${remainingSeconds}s`;
};

// 监听窗口大小变化
const handleResize = () => {
  if (chart.value) {
    chart.value.resize();
  }
};

onMounted(() => {
  fetchData();
  window.addEventListener("resize", handleResize);
});

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
  if (chart.value) {
    chart.value.dispose();
    chart.value = null;
  }
});
</script>

<style scoped>
.overview-card,
.recent-pipelines-card,
.activity-card,
.quick-actions-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-card {
  display: flex;
  align-items: center;
  padding: 15px;
  border-radius: 4px;
  margin-bottom: 20px;
}

.stat-card.success {
  background-color: rgba(103, 194, 58, 0.1);
}

.stat-card.running {
  background-color: rgba(64, 158, 255, 0.1);
}

.stat-card.failed {
  background-color: rgba(245, 108, 108, 0.1);
}

.stat-card.pending {
  background-color: rgba(144, 147, 153, 0.1);
}

.stat-icon {
  margin-right: 15px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  line-height: 1;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
}

.chart-container {
  margin-top: 20px;
}

.pipeline-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.pipeline-link:hover {
  text-decoration: underline;
}

.activity-timeline {
  padding: 10px 0;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quick-actions .el-button {
  width: 100%;
  justify-content: flex-start;
}
</style>
