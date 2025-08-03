<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>Schema详情</span>
        <el-button type="primary" @click="navigateBack" :icon="ArrowLeft">
          返回列表
        </el-button>
      </div>
    </template>

    <div v-loading="loading" class="schema-detail">
      <el-descriptions title="基本信息" :column="2">
        <el-descriptions-item label="名称">{{
          schema?.name
        }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{
          schema?.type
        }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{
          schema?.version
        }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{
          schema?.created_at
        }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{
          schema?.description || "-"
        }}</el-descriptions-item>
      </el-descriptions>

      <el-divider content-position="left">Schema内容</el-divider>
      <div class="schema-content">
        <pre v-if="formattedSchema" class="monospace">{{
          formattedSchema
        }}</pre>
        <p v-else-if="!loading" class="text-center text-gray-500">
          无效的JSON格式
        </p>
      </div>
    </div>
  </el-card>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from "vue";
import { ElMessage } from "element-plus";
import { ArrowLeft } from "@element-plus/icons-vue";
import { useRoute, useRouter } from "vue-router";
import { getYAMLSchemaById, YAMLSchema } from "../../api/yamlValidator";

const route = useRoute();
const router = useRouter();
const schema = ref<YAMLSchema | null>(null);
const loading = ref(false);

// 获取Schema详情
const fetchSchemaDetail = async () => {
  const id = Number(route.params.id);
  if (!id) {
    ElMessage.error("无效的Schema ID");
    navigateBack();
    return;
  }

  try {
    loading.value = true;
    const res = await getYAMLSchemaById(id);
    schema.value = res.data;
  } catch (error) {
    ElMessage.error("获取Schema详情失败");
    console.error("Failed to fetch schema detail:", error);
    navigateBack();
  } finally {
    loading.value = false;
  }
};

// 返回列表
const navigateBack = () => {
  router.push("/yaml-validator");
};

onMounted(() => {
  fetchSchemaDetail();
});

const formattedSchema = computed(() => {
  if (!schema.value?.schema) return null;

  try {
    const cleanSchema = schema.value.schema.replace(/`/g, "");
    const parsedSchema = JSON.parse(cleanSchema);
    return JSON.stringify(parsedSchema, null, 2);
  } catch (error) {
    console.error("Failed to parse schema:", error);
    return null;
  }
});
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.schema-detail {
  padding: 20px 0;
}

.schema-content {
  background-color: #f5f7fa;
  border-radius: 4px;
  padding: 16px;
  max-height: 500px;
  overflow-y: auto;
}

.monospace {
  font-family: monospace;
  white-space: pre-wrap;
  word-break: break-all;
  color: #333;
  margin: 0;
}

.text-center {
  text-align: center;
  padding: 20px 0;
}

.text-gray-500 {
  color: #909399;
}
</style>
