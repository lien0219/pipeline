<template>
  <div class="dag-detail-container">
    <div class="page-header">
      <h1>DAG详情</h1>
      <div>
        <el-button @click="handleBack">返回列表</el-button>
        <el-button type="primary" @click="handleEdit">编辑</el-button>
      </div>
    </div>

    <el-card v-loading="dagStore.loading">
      <div class="dag-info">
        <div class="info-item">
          <span class="label">ID:</span>
          <span class="value">{{ currentDag?.data?.id }}</span>
        </div>
        <div class="info-item">
          <span class="label">名称:</span>
          <span class="value">{{ currentDag?.data?.name }}</span>
        </div>
        <div class="info-item">
          <span class="label">描述:</span>
          <span class="value">{{ currentDag?.data?.description || "-" }}</span>
        </div>
        <div class="info-item">
          <span class="label">所属流水线ID:</span>
          <span class="value">{{ currentDag?.data?.pipeline_id }}</span>
        </div>
        <div class="info-item">
          <span class="label">状态:</span>
          <span class="value">
            <el-tag v-if="currentDag?.data?.is_active" type="success"
              >活动</el-tag
            >
            <el-tag v-else>非活动</el-tag>
          </span>
        </div>
        <div class="info-item">
          <span class="label">创建时间:</span>
          <span class="value">{{
            formatDate(currentDag?.data?.created_at)
          }}</span>
        </div>
        <div class="info-item">
          <span class="label">更新时间:</span>
          <span class="value">{{
            formatDate(currentDag?.data?.updated_at)
          }}</span>
        </div>
      </div>

      <div class="dag-nodes">
        <h2>DAG节点配置</h2>
        <el-divider />
        <pre>{{ JSON.stringify(currentDag?.data?.nodes || [], null, 2) }}</pre>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useDagStore } from "@/stores/dag";
import { ElMessage } from "element-plus";

const router = useRouter();
const route = useRoute();
const dagStore = useDagStore();
const currentDag = computed(() => dagStore.currentDag);

onMounted(async () => {
  const dagId = Number(route.params.id);
  if (dagId) {
    try {
      await dagStore.fetchDagDetail(dagId);
    } catch (err) {
      ElMessage.error(`获取DAG详情失败: ${dagStore.error}`);
      router.push("/dags");
    }
  } else {
    router.push("/dags");
  }
});

// 格式化日期
const formatDate = (dateString: any) => {
  if (!dateString) return "-";
  const date = new Date(dateString);
  return date.toLocaleString();
};

// 返回列表
const handleBack = () => {
  router.push("/dags");
};

// 编辑DAG
const handleEdit = () => {
  router.push(`/dags/${currentDag.value?.data?.id}/edit`);
};
</script>

<style scoped>
.dag-detail-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.dag-info {
  margin-bottom: 30px;
}

.info-item {
  display: flex;
  margin-bottom: 16px;
}

.label {
  width: 120px;
  font-weight: bold;
  color: #606266;
}

.value {
  flex: 1;
  color: #303133;
}

.dag-nodes {
  margin-top: 30px;
}

pre {
  background-color: #f5f5f5;
  padding: 16px;
  border-radius: 4px;
  overflow-x: auto;
  font-size: 14px;
}
</style>
