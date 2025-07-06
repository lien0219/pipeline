<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>发布管理</h2>
      </div>

      <div class="header-actions">
        <el-button type="primary" @click="createRelease">
          <el-icon><Plus /></el-icon>
          创建发布
        </el-button>
      </div>
    </div>

    <el-card>
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="环境" style="width: 200px">
            <el-select
              v-model="filterForm.environment"
              placeholder="全部环境"
              clearable
            >
              <el-option label="开发环境" value="development" />
              <el-option label="测试环境" value="testing" />
              <el-option label="预发布环境" value="staging" />
              <el-option label="生产环境" value="production" />
            </el-select>
          </el-form-item>

          <el-form-item label="状态" style="width: 200px">
            <el-select
              v-model="filterForm.status"
              placeholder="全部状态"
              clearable
            >
              <el-option label="成功" value="success" />
              <el-option label="失败" value="failed" />
              <el-option label="进行中" value="in_progress" />
              <el-option label="已回滚" value="rolled_back" />
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
        :data="releaseStore.releases"
        style="width: 100%"
        v-loading="releaseStore.loading"
      >
        <el-table-column prop="version" label="版本" width="120" />

        <el-table-column prop="environment" label="环境" width="120">
          <template #default="{ row }">
            <el-tag :type="getEnvironmentType(row.environment)">
              {{ getEnvironmentText(row.environment) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column
          prop="description"
          label="描述"
          min-width="200"
          show-overflow-tooltip
        />

        <el-table-column
          prop="artifact"
          label="制品"
          width="150"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <el-link type="primary" @click="viewArtifact(row.artifact_id)">
              {{ row.artifact }}
            </el-link>
          </template>
        </el-table-column>

        <el-table-column prop="deployed_by" label="部署者" width="120" />

        <el-table-column prop="deployed_at" label="部署时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.deployed_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              size="small"
              @click="viewRelease(row.id)"
            >
              详情
            </el-button>

            <el-button
              link
              type="warning"
              size="small"
              @click="rollbackRelease(row.id)"
              v-if="row.status === 'success' && !row.is_rollback"
            >
              回滚
            </el-button>

            <el-button
              link
              type="danger"
              size="small"
              @click="deleteRelease(row.id)"
              v-if="row.status !== 'in_progress'"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="releaseStore.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="创建发布" width="50%">
      <el-form
        ref="releaseFormRef"
        :model="releaseForm"
        :rules="rules"
        label-position="top"
      >
        <el-form-item label="版本" prop="version">
          <el-input
            v-model="releaseForm.version"
            placeholder="请输入版本号，例如：v1.0.0"
          />
        </el-form-item>

        <el-form-item label="环境" prop="environment">
          <el-select
            v-model="releaseForm.environment"
            placeholder="请选择部署环境"
          >
            <el-option label="开发环境" value="development" />
            <el-option label="测试环境" value="testing" />
            <el-option label="预发布环境" value="staging" />
            <el-option label="生产环境" value="production" />
          </el-select>
        </el-form-item>

        <el-form-item label="制品" prop="artifact_id">
          <el-select
            v-model="releaseForm.artifact_id"
            placeholder="请选择部署制品"
          >
            <el-option
              v-for="artifact in artifacts"
              :key="artifact.id"
              :label="artifact.name"
              :value="artifact.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="releaseForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入发布描述"
          />
        </el-form-item>

        <el-form-item label="发布说明">
          <el-input
            v-model="releaseForm.release_notes"
            type="textarea"
            :rows="5"
            placeholder="请输入发布说明"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button
            type="primary"
            @click="submitForm"
            :loading="releaseStore.submitting"
          >
            创建并部署
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus, Search, RefreshRight } from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { useReleaseStore } from "@/stores/release";
import { useRouter } from "vue-router";

const releaseStore = useReleaseStore();
const router = useRouter();

const releases = ref([]);
const artifacts = ref([]);
const loading = ref(false);
const dialogVisible = ref(false);
const submitting = ref(false);
const releaseFormRef = ref(null);

// 筛选表单
const filterForm = reactive({
  environment: "",
  status: "",
  dateRange: [],
});

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0,
});

// 发布表单
const releaseForm = reactive({
  version: "",
  environment: "",
  artifact_id: null,
  description: "",
  release_notes: "",
});

// 表单验证规则
const rules = {
  version: [
    { required: true, message: "请输入版本号", trigger: "blur" },
    {
      pattern: /^v\d+\.\d+\.\d+(-[\w.]+)?$/,
      message: "版本号格式不正确，例如：v1.0.0",
      trigger: "blur",
    },
  ],
  environment: [
    { required: true, message: "请选择部署环境", trigger: "change" },
  ],
  artifact_id: [
    { required: true, message: "请选择部署制品", trigger: "change" },
  ],
};

// 创建发布
const createRelease = () => {
  Object.assign(releaseForm, {
    version: "",
    environment: "",
    artifact_id: null,
    description: "",
    release_notes: "",
  });
  dialogVisible.value = true;
};

// 查看发布详情
const viewRelease = (id) => {
  router.push(`/deploy/releases/${id}`);
};

// 查看制品
const viewArtifact = (id) => {
  router.push(`/artifacts/${id}`);
};

// 回滚发布
const rollbackRelease = async (id) => {
  try {
    await ElMessageBox.confirm("确定要回滚此发布吗？", "回滚确认", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });
    await releaseStore.rollbackRelease(id);
  } catch (error) {
    if (error !== "cancel") {
      console.error("回滚发布失败:", error);
    }
  }
};

// 删除发布
const deleteRelease = async (id) => {
  try {
    await ElMessageBox.confirm(
      "确定要删除此发布记录吗？此操作不可恢复。",
      "删除确认",
      { confirmButtonText: "确定", cancelButtonText: "取消", type: "warning" }
    );
    await releaseStore.deleteRelease(id);
  } catch (error) {
    if (error !== "cancel") {
      console.error("删除发布失败:", error);
    }
  }
};

// 筛选
const handleFilter = () => {
  releaseStore.setFilters({
    environment: filterForm.environment,
    status: filterForm.status,
    dateRange: filterForm.dateRange,
  });
};
// 重置筛选
const resetFilter = () => {
  filterForm.environment = "";
  filterForm.status = "";
  filterForm.dateRange = [];
  releaseStore.setFilters({
    environment: "",
    status: "",
    dateRange: [],
  });
};

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  releaseStore.setPageParams(1, size);
};

const handleCurrentChange = (page) => {
  pagination.currentPage = page;
  releaseStore.setPageParams(page, pagination.pageSize);
};

// 提交表单
const submitForm = async () => {
  if (!releaseFormRef.value) return;

  await releaseFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await releaseStore.createRelease(releaseForm);
        dialogVisible.value = false;
      } catch (error) {}
    }
  });
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
const formatDate = (date) => {
  if (!date) return "-";
  return dayjs(date).format("YYYY-MM-DD HH:mm:ss");
};

onMounted(async () => {
  await Promise.all([
    releaseStore.fetchReleases(),
    releaseStore.fetchArtifacts(),
  ]);
  pagination.total = releaseStore.total;
  artifacts.value = releaseStore.artifacts;
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

.filter-container {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
