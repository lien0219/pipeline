<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>验证历史</span>
        <el-input
          v-model="searchKeyword"
          placeholder="搜索名称"
          style="width: 200px"
          :prefix-icon="Search"
        />
        <el-button
          type="primary"
          size="small"
          @click="handleSearch"
          style="margin-left: 10px"
        >
          搜索
        </el-button>
      </div>
    </template>

    <el-table v-loading="loading" :data="historyList" style="width: 100%">
      <el-table-column prop="name" label="名称" width="200" />
      <el-table-column prop="schema_type" label="Schema类型" width="120" />
      <el-table-column prop="schema_name" label="Schema名称" width="150" />
      <el-table-column prop="is_valid" label="验证结果" width="100">
        <template #default="scope">
          <span :class="scope.row.is_valid ? 'valid' : 'invalid'">
            {{ scope.row.is_valid ? "通过" : "失败" }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="scope">
          <el-button
            type="primary"
            size="small"
            @click="viewDetails(scope.row)"
            :icon="View"
          >
            查看
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.currentPage"
      v-model:page-size="pagination.pageSize"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      :total="pagination.total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      style="margin-top: 20px"
    />

    <!-- 详情对话框 -->
    <el-dialog v-model="dialogVisible" title="验证详情" width="80%">
      <div v-if="selectedHistory" class="history-details">
        <el-descriptions title="基本信息" :column="2">
          <el-descriptions-item label="名称">{{
            selectedHistory.name
          }}</el-descriptions-item>
          <el-descriptions-item label="Schema类型">{{
            selectedHistory.schema_type
          }}</el-descriptions-item>
          <el-descriptions-item label="Schema名称">{{
            selectedHistory.schema_name || "-"
          }}</el-descriptions-item>
          <el-descriptions-item label="验证结果">
            <span :class="selectedHistory.is_valid ? 'valid' : 'invalid'">
              {{ selectedHistory.is_valid ? "通过" : "失败" }}
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{
            selectedHistory.created_at
          }}</el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">YAML内容</el-divider>
        <el-input
          v-model="selectedHistory.content"
          type="textarea"
          :rows="10"
          readonly
          class="monospace"
        />

        <el-divider v-if="!selectedHistory.is_valid" content-position="left"
          >错误信息</el-divider
        >
        <el-alert
          v-if="!selectedHistory.is_valid"
          title="错误信息"
          type="error"
          :closable="false"
          :description="selectedHistory.errors"
        />
      </div>
    </el-dialog>
  </el-card>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, watch } from "vue";
import { ElMessage } from "element-plus";
import { Search, View } from "@element-plus/icons-vue";
import {
  getValidationHistory,
  ValidateYAMLResult,
  PageResponse,
} from "../../api/yamlValidator";

const historyList = ref<ValidateYAMLResult[]>([]);
const loading = ref(false);
const searchKeyword = ref("");
const dialogVisible = ref(false);
const selectedHistory = ref<ValidateYAMLResult | null>(null);

const pagination = reactive<{
  currentPage: number;
  pageSize: number;
  total: number;
}>({
  currentPage: 1,
  pageSize: 10,
  total: 0,
});

// 获取验证历史
const fetchHistory = async () => {
  try {
    loading.value = true;
    const res: any = await getValidationHistory({
      page: pagination.currentPage,
      pageSize: pagination.pageSize,
      name: searchKeyword.value,
    });
    historyList.value = res.data.list;
    pagination.total = res.data.total;
  } catch (error) {
    ElMessage.error("获取验证历史失败");
    console.error("Failed to fetch history:", error);
  } finally {
    loading.value = false;
  }
};
const handleSearch = () => {
  pagination.currentPage = 1;
  fetchHistory();
};
// 查看详情
const viewDetails = (item: ValidateYAMLResult) => {
  selectedHistory.value = { ...item };
  dialogVisible.value = true;
};

// 分页处理
const handleSizeChange = (size: number) => {
  pagination.pageSize = size;
  fetchHistory();
};

const handleCurrentChange = (current: number) => {
  pagination.currentPage = current;
  fetchHistory();
};
watch(searchKeyword, (newValue) => {
  const timer = setTimeout(() => {
    if (newValue !== searchKeyword.value) return;
    handleSearch();
  }, 500);
  return () => clearTimeout(timer);
});
// 初始化
onMounted(() => {
  fetchHistory();
});
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.valid {
  color: #67c23a;
}

.invalid {
  color: #f56c6c;
}

.history-details {
  padding: 10px;
}

.monospace {
  font-family: monospace;
}
</style>
