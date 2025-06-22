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
              <div
                class="stat-card success"
                @mouseenter="cardHover($event)"
                @mouseleave="cardLeave($event)"
              >
                <div class="stat-icon">
                  <el-icon :size="24"><Check /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value" :data-value="stats.success">
                    {{ stats.success }}
                  </div>
                  <div class="stat-label">成功</div>
                </div>
                <div class="stat-bg"></div>
              </div>
            </el-col>

            <el-col :span="6">
              <div
                class="stat-card running"
                @mouseenter="cardHover($event)"
                @mouseleave="cardLeave($event)"
              >
                <div class="stat-icon">
                  <el-icon :size="24"><Loading /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value" :data-value="stats.running">
                    {{ stats.running }}
                  </div>
                  <div class="stat-label">运行中</div>
                </div>
                <div class="stat-bg"></div>
              </div>
            </el-col>

            <el-col :span="6">
              <div
                class="stat-card failed"
                @mouseenter="cardHover($event)"
                @mouseleave="cardLeave($event)"
              >
                <div class="stat-icon">
                  <el-icon :size="24"><WarningFilled /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value" :data-value="stats.failed">
                    {{ stats.failed }}
                  </div>
                  <div class="stat-label">失败</div>
                </div>
                <div class="stat-bg"></div>
              </div>
            </el-col>

            <el-col :span="6">
              <div
                class="stat-card pending"
                @mouseenter="cardHover($event)"
                @mouseleave="cardLeave($event)"
              >
                <div class="stat-icon">
                  <el-icon :size="24"><Clock /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-value" :data-value="stats.pending">
                    {{ stats.pending }}
                  </div>
                  <div class="stat-label">等待中</div>
                </div>
                <div class="stat-bg"></div>
              </div>
            </el-col>
          </el-row>

          <div class="chart-container">
            <div ref="pipelineChart" style="width: 100%; height: 400px"></div>
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

            <el-table-column prop="branch" label="分支" width="120">
              <template #default="{ row }">
                {{ row.git_branch }}
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
                {{ activity.Content }}
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
import { pipelineApi } from "@/api/pipeline";
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
    const [
      pipelinesResponse,
      statsResponse,
      activitiesResponse,
      chartDataResponse,
    ] = await Promise.all([
      pipelineStore.fetchPipelines({ page: 1, pageSize: 10 }),
      pipelineStore.getDashboardStats(),
      pipelineStore.getDashboardActivities({ limit: 10 }),
      pipelineApi.getPipelineChartData(),
    ]);

    recentPipelines.value = pipelinesResponse.data.list || [];
    stats.value = {
      success: statsResponse?.Success,
      running: statsResponse?.Running,
      failed: statsResponse?.Failed,
      pending: statsResponse?.Pending,
    };
    activities.value = activitiesResponse;

    const chartData = chartDataResponse.data;
    initChart(chartData);
    setTimeout(startValueAnimation, 100);
  } catch (error) {
    console.error("Failed to fetch dashboard data:", error);
  } finally {
    loading.value = false;
  }
};

// 初始化图表
const initChart = (chartData) => {
  if (!chartData || !Array.isArray(chartData.dates)) {
    console.error("Invalid chartData structure:", chartData);
    return;
  }
  if (!pipelineChart.value) return;

  if (chart.value) {
    chart.value.dispose();
  }

  chart.value = echarts.init(pipelineChart.value);

  const option = {
    backgroundColor: "#1a1a1a",
    title: {
      text: "流水线执行趋势",
      left: "center",
      textStyle: {
        color: "#fff",
        fontSize: 20,
      },
    },
    tooltip: {
      trigger: "axis",
      axisPointer: {
        type: "cross",
        animation: false,
        label: {
          backgroundColor: "#505765",
        },
      },
    },
    legend: {
      data: ["成功", "失败", "总数"],
      bottom: 20,
      textStyle: {
        color: "#fff",
      },
    },
    xAxis: {
      type: "category",
      boundaryGap: false,
      data: chartData.dates || [],
      axisLine: {
        lineStyle: {
          color: "#8392A5",
        },
      },
      axisLabel: {
        color: "#fff",
      },
    },
    yAxis: {
      type: "value",
      axisLine: {
        lineStyle: {
          color: "#8392A5",
        },
      },
      splitLine: {
        lineStyle: {
          color: "rgba(255,255,255,0.1)",
        },
      },
      axisLabel: {
        color: "#fff",
      },
    },
    series: [
      {
        name: "成功",
        type: "line",
        smooth: true,
        symbolSize: 10,
        itemStyle: {
          color: "#00ff00",
          shadowColor: "rgba(0, 255, 0, 0.5)",
          shadowBlur: 10,
          shadowOffsetY: 5,
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: "rgba(0, 255, 0, 0.3)",
            },
            {
              offset: 1,
              color: "rgba(0, 255, 0, 0.1)",
            },
          ]),
        },
        data: chartData.success || [],
        animationDelay: (idx) => idx * 100,
      },
      {
        name: "失败",
        type: "line",
        smooth: true,
        symbolSize: 10,
        itemStyle: {
          color: "#ff0000",
          shadowColor: "rgba(255, 0, 0, 0.5)",
          shadowBlur: 10,
          shadowOffsetY: 5,
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: "rgba(255, 0, 0, 0.3)",
            },
            {
              offset: 1,
              color: "rgba(255, 0, 0, 0.1)",
            },
          ]),
        },
        data: chartData.failed || [],
        animationDelay: (idx) => idx * 100 + 100,
      },
      {
        name: "总数",
        type: "line",
        smooth: true,
        symbolSize: 10,
        itemStyle: {
          color: "#0000ff",
          shadowColor: "rgba(0, 0, 255, 0.5)",
          shadowBlur: 10,
          shadowOffsetY: 5,
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            {
              offset: 0,
              color: "rgba(0, 0, 255, 0.3)",
            },
            {
              offset: 1,
              color: "rgba(0, 0, 255, 0.1)",
            },
          ]),
        },
        data: chartData.total || [],
        animationDelay: (idx) => idx * 100 + 200,
      },
    ],
    animationEasing: "elasticOut",
    animationDuration: 1500,
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
    case "inactive":
      return "info";
    case "active":
      return "success";
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
    case "active":
      return "已激活";
    case "inactive":
      return "未使用";
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

const cardHover = (e) => {
  e.currentTarget.classList.add("hover");
};

const cardLeave = (e) => {
  e.currentTarget.classList.remove("hover");
};

// 数字增长动画
const animateValue = (el, start, end, duration) => {
  let startTimestamp = null;
  const step = (timestamp) => {
    if (!startTimestamp) startTimestamp = timestamp;
    const progress = Math.min((timestamp - startTimestamp) / duration, 1);
    el.innerText = Math.floor(progress * (end - start) + start);
    if (progress < 1) {
      window.requestAnimationFrame(step);
    }
  };
  window.requestAnimationFrame(step);
};

// 在数据加载后执行数字动画
const startValueAnimation = () => {
  document.querySelectorAll(".stat-value").forEach((el) => {
    const value = parseInt(el.dataset.value);
    if (!isNaN(value)) {
      animateValue(el, 0, value, 1500);
    }
  });
};
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
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 20px;
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.stat-card::before {
  content: "";
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: inherit;
  z-index: -1;
  filter: blur(8px);
  opacity: 0.7;
  transform: scale(0.95);
}

.stat-card.hover {
  transform: translateY(-5px);
  box-shadow: 0 12px 20px rgba(0, 0, 0, 0.15);
}

.stat-card.success {
  background: linear-gradient(
    135deg,
    rgba(103, 194, 58, 0.2),
    rgba(103, 194, 58, 0.05)
  );
}

.stat-card.running {
  background: linear-gradient(
    135deg,
    rgba(64, 158, 255, 0.2),
    rgba(64, 158, 255, 0.05)
  );
}

.stat-card.failed {
  background: linear-gradient(
    135deg,
    rgba(245, 108, 108, 0.2),
    rgba(245, 108, 108, 0.05)
  );
}

.stat-card.pending {
  background: linear-gradient(
    135deg,
    rgba(144, 147, 153, 0.2),
    rgba(144, 147, 153, 0.05)
  );
}

.stat-icon {
  margin-right: 15px;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

.stat-card.hover .stat-icon {
  transform: scale(1.1) rotate(5deg);
}

.stat-info {
  flex: 1;
  position: relative;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  line-height: 1.2;
  margin-bottom: 5px;
  font-family: "Segoe UI", Roboto, sans-serif;
  background: linear-gradient(90deg, #fff, #e0e0e0);
  -webkit-background-clip: text;
  background-clip: text;
  color: #000 !important;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.stat-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-bg {
  position: absolute;
  top: 0;
  right: 0;
  width: 60px;
  height: 60px;
  opacity: 0.1;
}

.stat-card.success .stat-bg {
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%2367c23a' d='M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z'/%3E%3C/svg%3E")
    no-repeat center;
  background-size: contain;
}

.stat-card.running .stat-bg {
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23409eff' d='M15.07 1.01h-6v2h6v-2zm-4 13h2v-6h-2v6zm8.03-6.62l1.42-1.42c-.43-.51-.9-.99-1.41-1.41l-1.42 1.42C16.14 4.74 14.19 4 12.07 4c-4.97 0-9 4.03-9 9s4.02 9 9 9 9-4.03 9-9c0-2.11-.74-4.06-1.97-5.61z'/%3E%3C/svg%3E")
    no-repeat center;
  background-size: contain;
}

.stat-card.failed .stat-bg {
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23f56c6c' d='M12 2c5.52 0 10 4.48 10 10s-4.48 10-10 10S2 17.52 2 12 6.48 2 12 2zm-1 15h2v-2h-2v2zm0-4h2V7h-2v6z'/%3E%3C/svg%3E")
    no-repeat center;
  background-size: contain;
}

.stat-card.pending .stat-bg {
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath fill='%23909399' d='M12 22C6.47 22 2 17.53 2 12S6.47 2 12 2s10 4.47 10 10-4.47 10-10 10zm0-18c-4.41 0-8 3.59-8 8s3.59 8 8 8 8-3.59 8-8-3.59-8-8-8zm.5 13H11v-6h1.5v6zm0-8H11V7h1.5v2z'/%3E%3C/svg%3E")
    no-repeat center;
  background-size: contain;
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
  background: #434242;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
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
