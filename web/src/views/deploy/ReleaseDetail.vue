<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>发布详情</h2>
      </div>
      <div class="header-actions">
        <el-button @click="goBack">
          <el-icon><ArrowLeft /></el-icon>
          返回列表
        </el-button>
      </div>
    </div>

    <el-card v-loading="loading">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ release.id }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{
          release.version
        }}</el-descriptions-item>
        <el-descriptions-item label="环境">
          <el-tag :type="getEnvironmentType(release.environment)">
            {{ getEnvironmentText(release.environment) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(release.status)">
            {{ getStatusText(release.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="制品">
          <router-link
            v-if="release.artifact_id"
            :to="`/artifacts/${release.artifact_id}`"
            class="artifact-link"
          >
            {{ release.artifact_name }}
          </router-link>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="部署者">{{
          release.deployed_by
        }}</el-descriptions-item>
        <el-descriptions-item label="部署时间">{{
          formatDate(release.deployed_at)
        }}</el-descriptions-item>
        <el-descriptions-item label="回滚状态">
          <el-tag v-if="release.is_rollback" type="warning">是</el-tag>
          <el-tag v-else type="success">否</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="回滚版本" v-if="release.rollback_to">
          <router-link
            :to="`/deploy/releases/${release.rollback_to}`"
            class="release-link"
          >
            {{ release.rollback_to }}
          </router-link>
        </el-descriptions-item>
        <el-descriptions-item label="回滚时间" v-if="release.rolled_back_at">{{
          formatDate(release.rolled_back_at)
        }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{
          release.description || "无描述"
        }}</el-descriptions-item>
      </el-descriptions>

      <div class="release-notes" style="margin-top: 20px">
        <h3 style="margin-bottom: 10px">发布说明</h3>
        <el-card>
          <div v-html="formatReleaseNotes(release.release_notes)"></div>
        </el-card>
      </div>

      <div class="release-actions" style="margin-top: 20px; text-align: center">
        <el-button
          type="warning"
          @click="rollbackRelease"
          v-if="release.status === 'success' && !release.is_rollback"
        >
          <el-icon><RefreshLeft /></el-icon>
          回滚发布
        </el-button>
        <el-button
          type="danger"
          @click="deleteRelease"
          v-if="release.status !== 'in_progress'"
        >
          <el-icon><Delete /></el-icon>
          删除发布记录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { ArrowLeft, RefreshLeft, Delete } from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { releaseApi } from "@/api/release";
import { useReleaseStore } from "@/stores/release";

const route = useRoute();
const router = useRouter();
const releaseStore = useReleaseStore();
const releaseId = route.params.id;
const release = ref({});
const loading = ref(true);

// 返回列表
const goBack = () => {
  router.push("/deploy/releases");
};

// 获取环境类型样式
const getEnvironmentType = (environment) => {
  switch (environment) {
    case "development":
      return "info";
    case "testing":
      return "warning";
    case "staging":
      return "success";
    case "production":
      return "danger";
    default:
      return "info";
  }
};

// 获取环境文本
const getEnvironmentText = (environment) => {
  switch (environment) {
    case "development":
      return "开发环境";
    case "testing":
      return "测试环境";
    case "staging":
      return "预发布环境";
    case "production":
      return "生产环境";
    default:
      return environment;
  }
};

// 获取状态样式
const getStatusType = (status) => {
  switch (status) {
    case "success":
      return "success";
    case "failed":
      return "danger";
    case "in_progress":
      return "primary";
    case "rolled_back":
      return "warning";
    default:
      return "info";
  }
};

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case "success":
      return "成功";
    case "failed":
      return "失败";
    case "in_progress":
      return "进行中";
    case "rolled_back":
      return "已回滚";
    default:
      return "未知";
  }
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return "-";
  return dayjs(dateString).format("YYYY-MM-DD HH:mm:ss");
};

// 格式化发布说明（将换行符转换为<br>）
const formatReleaseNotes = (notes) => {
  if (!notes) return "<p>无发布说明</p>";
  return notes.replace(/\n/g, "<br>");
};

// 回滚发布
const rollbackRelease = async () => {
  try {
    await ElMessageBox.confirm("确定要回滚此发布吗？", "回滚确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });
    await releaseStore.rollbackRelease(releaseId);
    router.push("/deploy/releases");
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("回滚发布失败");
      console.error("回滚发布失败:", error);
    }
  }
};

// 删除发布
const deleteRelease = async () => {
  try {
    await ElMessageBox.confirm(
      "确定要删除此发布记录吗？此操作不可恢复。",
      "删除确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }
    );
    await releaseStore.deleteRelease(releaseId);
    router.push("/deploy/releases");
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("删除发布失败");
      console.error("删除发布失败:", error);
    }
  }
};

// 获取发布详情
const fetchReleaseDetail = async () => {
  try {
    loading.value = true;
    const response = await releaseApi.getReleaseById(releaseId);
    release.value = response.data;
  } catch (error) {
    ElMessage.error("获取发布详情失败");
    console.error("获取发布详情失败:", error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchReleaseDetail();
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
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

.artifact-link,
.release-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.artifact-link:hover,
.release-link:hover {
  text-decoration: underline;
}

.release-notes h3 {
  font-size: 16px;
  color: #303133;
}
</style>
