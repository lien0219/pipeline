<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>制品详情</h2>
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
        <el-descriptions-item label="ID">{{
          artifact.id
        }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{
          artifact.name
        }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag size="small">{{ artifact.type?.toUpperCase() }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="大小">{{
          formatFileSize(artifact.size)
        }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{
          artifact.version
        }}</el-descriptions-item>
        <el-descriptions-item label="下载次数">{{
          artifact.download_count
        }}</el-descriptions-item>
        <el-descriptions-item label="流水线">
          <router-link
            v-if="artifact.pipeline_id"
            :to="`/pipelines/${artifact.pipeline_id}`"
            class="pipeline-link"
          >
            {{ artifact.pipeline?.name }}
          </router-link>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="构建ID">
          <router-link
            v-if="artifact.pipeline_id && artifact.pipeline_run_id"
            :to="`/pipelines/${artifact.pipeline_id}/runs/${artifact.pipeline_run_id}`"
            class="pipeline-link"
          >
            {{ artifact.pipeline_run_id }}
          </router-link>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="创建者">{{
          artifact.created_by
        }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{
          formatDate(artifact.created_at)
        }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{
          artifact.description || "无描述"
        }}</el-descriptions-item>
        <el-descriptions-item label="存储路径" :span="2">{{
          artifact.storage_path || "-"
        }}</el-descriptions-item>
      </el-descriptions>

      <div
        class="artifact-actions"
        style="margin-top: 20px; text-align: center"
      >
        <el-button type="primary" @click="downloadArtifact">
          <el-icon><Download /></el-icon>
          下载制品
        </el-button>
        <el-button type="success" @click="deployArtifact">
          <el-icon><Position /></el-icon>
          部署制品
        </el-button>
        <el-button type="danger" @click="deleteArtifact">
          <el-icon><Delete /></el-icon>
          删除制品
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { ArrowLeft, Download, Position, Delete } from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { artifactApi } from "@/api/artifact";

const route = useRoute();
const router = useRouter();
const artifactId = route.params.id;
const artifact = ref({});
const loading = ref(true);

// 返回列表
const goBack = () => {
  router.push("/artifacts");
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

// 格式化日期
const formatDate = (dateString) => {
  return dayjs(dateString).format("YYYY-MM-DD HH:mm:ss");
};

// 下载制品
const downloadArtifact = async () => {
  try {
    const response = await artifactApi.downloadArtifact(artifactId);
    const blob = new Blob([response.data]);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    const fileName = `${artifact.value.name}-${artifact.value.version}.${artifact.value.type}`;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
  } catch (error) {
    ElMessage.error("制品下载失败");
    console.error(error);
  }
};

// 部署制品
const deployArtifact = () => {
  router.push({ name: "Releases", query: { artifact_id: artifactId } });
};

// 删除制品
const deleteArtifact = async () => {
  try {
    await ElMessageBox.confirm("确定要删除此制品吗？此操作不可撤销。", "警告", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await artifactApi.deleteArtifact(artifactId);
    ElMessage.success("制品删除成功");
    router.push("/artifacts");
  } catch (error) {
    if (error === "cancel") return;
    ElMessage.error("制品删除失败");
    console.error(error);
  }
};

// 获取制品详情
const fetchArtifactDetail = async () => {
  try {
    loading.value = true;
    const response = await artifactApi.getArtifactById(artifactId);
    artifact.value = response.data;
  } catch (error) {
    ElMessage.error("获取制品详情失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchArtifactDetail();
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

.pipeline-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.pipeline-link:hover {
  text-decoration: underline;
}
</style>
