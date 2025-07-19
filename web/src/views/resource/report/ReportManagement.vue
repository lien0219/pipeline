<template>
  <div class="resource-report-container">
    <el-card class="box-card">
      <div slot="header" class="clearfix">
        <h2>资源报告管理</h2>
        <el-button
          style="float: right"
          type="primary"
          icon="el-icon-plus"
          @click="handleAdd"
        >
          创建资源报告
        </el-button>
      </div>

      <div class="table-container">
        <el-row :gutter="20" class="mb-4">
          <el-col :span="8">
            <el-input
              v-model="searchForm.tenantId"
              placeholder="请输入租户ID搜索"
              prefix-icon="el-icon-search"
              @keyup.enter.native="handleSearch"
            />
          </el-col>
          <el-col :span="4">
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
          </el-col>
        </el-row>

        <el-table
          v-loading="loading"
          :data="reportList"
          border
          stripe
          style="width: 100%"
        >
          <el-table-column
            type="index"
            label="序号"
            width="80"
          ></el-table-column>
          <el-table-column
            prop="tenant_id"
            label="租户ID"
            width="180"
          ></el-table-column>
          <el-table-column
            prop="cpu_usage"
            label="CPU使用量"
            width="120"
          ></el-table-column>
          <el-table-column
            prop="cpu_quota"
            label="CPU配额"
            width="120"
          ></el-table-column>
          <el-table-column
            prop="memory_usage"
            label="内存使用量(MB)"
            width="150"
          ></el-table-column>
          <el-table-column
            prop="memory_quota"
            label="内存配额(MB)"
            width="150"
          ></el-table-column>
          <el-table-column
            prop="storage_usage"
            label="存储使用量(GB)"
            width="150"
          ></el-table-column>
          <el-table-column
            prop="storage_quota"
            label="存储配额(GB)"
            width="150"
          ></el-table-column>
          <el-table-column
            prop="created_at"
            label="创建时间"
            width="180"
          ></el-table-column>
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="scope">
              <el-button type="text" size="small" @click="handleView(scope.row)"
                >查看</el-button
              >
              <el-button type="text" size="small" @click="handleEdit(scope.row)"
                >编辑</el-button
              >
              <el-button
                type="text"
                size="small"
                @click="handleDelete(scope.row)"
                >删除</el-button
              >
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
            :current-page="pagination.page"
            :page-sizes="[10, 20, 50, 100]"
            :page-size="pagination.pageSize"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
          ></el-pagination>
        </div>
      </div>
    </el-card>
    <el-dialog :title="dialogTitle" v-model="dialogVisible" width="60%">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="租户ID" prop="tenantId">
          <el-input
            v-model="form.tenantId"
            :disabled="dialogType === 'view' || dialogType === 'edit'"
          ></el-input>
        </el-form-item>
        <el-form-item label="CPU使用量" prop="cpuUsage">
          <el-input
            v-model.number="form.cpuUsage"
            :disabled="dialogType === 'view'"
            placeholder="请输入CPU使用量"
          ></el-input>
        </el-form-item>
        <el-form-item label="CPU配额" prop="cpuQuota">
          <el-input
            v-model.number="form.cpuQuota"
            :disabled="dialogType === 'view'"
            placeholder="请输入CPU配额"
          ></el-input>
        </el-form-item>
        <el-form-item label="内存使用量(MB)" prop="memoryUsage">
          <el-input
            v-model.number="form.memoryUsage"
            :disabled="dialogType === 'view'"
            placeholder="请输入内存使用量"
          ></el-input>
        </el-form-item>
        <el-form-item label="内存配额(MB)" prop="memoryQuota">
          <el-input
            v-model.number="form.memoryQuota"
            :disabled="dialogType === 'view'"
            placeholder="请输入内存配额"
          ></el-input>
        </el-form-item>
        <el-form-item label="存储使用量(GB)" prop="storageUsage">
          <el-input
            v-model.number="form.storageUsage"
            :disabled="dialogType === 'view'"
            placeholder="请输入存储使用量"
          ></el-input>
        </el-form-item>
        <el-form-item label="存储配额(GB)" prop="storageQuota">
          <el-input
            v-model.number="form.storageQuota"
            :disabled="dialogType === 'view'"
            placeholder="请输入存储配额"
          ></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button
          v-if="dialogType !== 'view'"
          type="primary"
          @click="handleSubmit"
          >确定</el-button
        >
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, Ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  getResourceReports,
  getResourceReportById,
  createResourceReport,
  updateResourceReport,
  deleteResourceReport,
} from "@/api/system";

// 状态定义
const loading = ref(false);
const reportList = ref([]);
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});
const searchForm = reactive({
  tenantId: "",
});
const dialogVisible = ref(false);
const dialogType = ref<string>("add");
const dialogTitle = ref("");
const form = reactive({
  id: 0,
  tenantId: "",
  cpuUsage: 0,
  cpuQuota: 0,
  memoryUsage: 0,
  memoryQuota: 0,
  storageUsage: 0,
  storageQuota: 0,
});
const formRef = ref();
const rules = reactive({
  tenantId: [{ required: true, message: "请输入租户ID", trigger: "blur" }],
  cpuUsage: [{ required: true, message: "请输入CPU使用量", trigger: "blur" }],
  cpuQuota: [{ required: true, message: "请输入CPU配额", trigger: "blur" }],
  memoryUsage: [
    { required: true, message: "请输入内存使用量", trigger: "blur" },
  ],
  memoryQuota: [{ required: true, message: "请输入内存配额", trigger: "blur" }],
  storageUsage: [
    { required: true, message: "请输入存储使用量", trigger: "blur" },
  ],
  storageQuota: [
    { required: true, message: "请输入存储配额", trigger: "blur" },
  ],
});

// 生命周期
onMounted(() => {
  fetchResourceReports();
});

// 方法定义
const fetchResourceReports = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      tenant_id: searchForm.tenantId,
    };
    const res = await getResourceReports(params);
    reportList.value = res.data.list;
    pagination.total = res.data.total;
  } catch (error) {
    ElMessage.error("获取资源报告列表失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  pagination.page = 1;
  fetchResourceReports();
};

const resetSearch = () => {
  searchForm.tenantId = "";
  pagination.page = 1;
  fetchResourceReports();
};

const handleSizeChange = (val) => {
  pagination.pageSize = val;
  fetchResourceReports();
};

const handleCurrentChange = (val) => {
  pagination.page = val;
  fetchResourceReports();
};

const handleAdd = () => {
  dialogType.value = "add";
  dialogTitle.value = "创建资源报告";
  Object.assign(form, {
    id: 0,
    tenantId: "",
    cpuUsage: 0,
    cpuQuota: 0,
    memoryUsage: 0,
    memoryQuota: 0,
    storageUsage: 0,
    storageQuota: 0,
  });
  dialogVisible.value = true;
};

const handleView = async (row) => {
  dialogType.value = "view";
  dialogTitle.value = "查看资源报告";
  loading.value = true;
  try {
    const res = await getResourceReportById(row.id);
    const data = res.data;
    form.id = data.id;
    form.tenantId = data.tenant_id;
    form.cpuUsage = data.cpu_usage;
    form.cpuQuota = data.cpu_quota;
    form.memoryUsage = data.memory_usage;
    form.memoryQuota = data.memory_quota;
    form.storageUsage = data.storage_usage;
    form.storageQuota = data.storage_quota;
    dialogVisible.value = true;
  } catch (error) {
    ElMessage.error("获取资源报告详情失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const handleEdit = async (row) => {
  dialogType.value = "edit";
  dialogTitle.value = "编辑资源报告";
  loading.value = true;
  try {
    const res = await getResourceReportById(row.id);
    const data = res.data;
    form.id = data.id;
    form.tenantId = data.tenant_id;
    form.cpuUsage = data.cpu_usage;
    form.cpuQuota = data.cpu_quota;
    form.memoryUsage = data.memory_usage;
    form.memoryQuota = data.memory_quota;
    form.storageUsage = data.storage_usage;
    form.storageQuota = data.storage_quota;
    dialogVisible.value = true;
  } catch (error) {
    ElMessage.error("获取资源报告详情失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

const handleDelete = (row) => {
  ElMessageBox.confirm("确定要删除该资源报告吗？", "提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    loading.value = true;
    try {
      await deleteResourceReport(row.id);
      ElMessage.success("删除成功");
      fetchResourceReports();
    } catch (error) {
      ElMessage.error("删除失败");
      console.error(error);
    } finally {
      loading.value = false;
    }
  });
};

const handleSubmit = async () => {
  const formEl = formRef.value;
  if (!formEl) return;
  await formEl.validate();

  loading.value = true;
  try {
    const reportData = {
      tenant_id: form.tenantId,
      cpu_usage: form.cpuUsage,
      cpu_quota: form.cpuQuota,
      memory_usage: form.memoryUsage,
      memory_quota: form.memoryQuota,
      storage_usage: form.storageUsage,
      storage_quota: form.storageQuota,
    };

    if (dialogType.value === "add") {
      await createResourceReport(reportData);
      ElMessage.success("创建成功");
    } else {
      await updateResourceReport(form.id, reportData);
      ElMessage.success("更新成功");
    }
    dialogVisible.value = false;
    fetchResourceReports();
  } catch (error) {
    ElMessage.error(dialogType.value === "add" ? "创建失败" : "更新失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* 全局容器样式 */
.resource-report-container {
  padding: 20px;
  min-height: calc(100vh - 60px);
  background-color: #f5f7fa;
}

/* 卡片样式优化 */
.box-card {
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  transition: all 0.3s ease;
}

.box-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

/* 卡片头部样式 */
.clearfix {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.clearfix span {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

/* 搜索区域样式 */
.table-container {
  padding: 20px;
}

.mb-4 {
  margin-bottom: 16px !important;
}

/* 表格样式优化 */
.el-table {
  border-radius: 6px;
  overflow: hidden;
}

.el-table th {
  background-color: #f7f8fa;
  font-weight: 500;
  color: #606266;
}

.el-table tr:hover > td {
  background-color: #f5f9ff !important;
}

.el-table__row--striped td {
  background-color: #fafafa;
}

/* 分页区域样式 */
.pagination-container {
  margin-top: 20px;
  text-align: right;
  padding-right: 10px;
}

/* 按钮样式优化 */
.el-button {
  border-radius: 4px;
  transition: all 0.2s ease;
}

.el-button--primary {
  background-color: #409eff;
  border-color: #409eff;
}

.el-button--primary:hover {
  background-color: #66b1ff;
  border-color: #66b1ff;
}

/* 对话框样式优化 */
.el-dialog {
  border-radius: 8px;
  overflow: hidden;
}

.el-dialog__header {
  background-color: #f7f8fa;
  padding: 15px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.el-dialog__title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

.el-dialog__body {
  padding: 25px 20px;
}

/* 表单样式优化 */
.el-form-item {
  margin-bottom: 20px;
}

.el-form-item__label {
  color: #606266;
  font-weight: 400;
}

.el-input {
  border-radius: 4px;
}

/* 操作按钮样式 */
.el-button--text {
  color: #409eff;
}

.el-button--text:hover {
  color: #66b1ff;
  background-color: rgba(64, 158, 255, 0.05);
}

/* 响应式调整 */
@media (max-width: 768px) {
  .resource-report-container {
    padding: 10px;
  }

  .table-container {
    padding: 10px;
  }

  .el-col {
    width: 100% !important;
    margin-bottom: 10px;
  }
}
</style>
