<template>
  <div class="resource-request-container">
    <el-card shadow="hover">
      <div class="page-header">
        <h2>资源请求管理</h2>
      </div>

      <div class="search-and-action">
        <el-row :gutter="20" align="middle">
          <el-col :span="10">
            <el-input
              v-model="searchTenantId"
              placeholder="输入租户ID搜索"
              clearable
              prefix-icon="Search"
              @keyup.enter="handleSearch"
            >
              <template #append>
                <el-button type="primary" @click="handleSearch">搜索</el-button>
              </template>
            </el-input>
          </el-col>
          <el-col :span="4" :offset="10">
            <el-button type="primary" @click="openCreateDialog" icon="Plus">
              创建资源请求
            </el-button>
          </el-col>
        </el-row>
      </div>

      <el-table
        v-loading="loading"
        :data="requestList"
        border
        style="width: 100%; margin-top: 20px"
        empty-text="暂无数据"
      >
        <el-table-column prop="id" label="请求ID" width="180"></el-table-column>
        <el-table-column
          prop="tenant_id"
          label="租户ID"
          width="180"
        ></el-table-column>
        <el-table-column prop="cpu_request" label="CPU请求(核)" width="150">
          <template #default="scope">{{
            scope.row.cpu_request.toFixed(2)
          }}</template>
        </el-table-column>
        <el-table-column
          prop="memory_request"
          label="内存请求(MB)"
          width="180"
        ></el-table-column>
        <el-table-column
          prop="storage_request"
          label="存储请求(MB)"
          width="180"
        ></el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)">
              {{ formatStatus(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          label="提交时间"
          width="220"
        ></el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="scope">
            <el-button
              v-if="scope.row.status === 'pending'"
              type="success"
              size="small"
              @click="handleApprove(scope.row)"
              style="margin-right: 8px"
            >
              批准
            </el-button>
            <el-button
              v-if="scope.row.status === 'pending'"
              type="danger"
              size="small"
              @click="handleReject(scope.row)"
            >
              拒绝
            </el-button>
            <span v-else>无操作</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建资源请求对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="'创建资源请求'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="requestFormRef"
        :model="requestForm"
        :rules="requestRules"
        label-width="120px"
      >
        <el-form-item label="租户ID" prop="tenant_id">
          <el-input
            v-model="requestForm.tenant_id"
            placeholder="请输入租户ID"
          />
        </el-form-item>
        <el-form-item label="CPU请求(核)" prop="cpu_request">
          <el-input-number
            v-model.number="requestForm.cpu_request"
            :min="0.1"
            :step="0.1"
            placeholder="请输入CPU请求"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="内存请求(MB)" prop="memory_request">
          <el-input-number
            v-model.number="requestForm.memory_request"
            :min="1"
            placeholder="请输入内存请求"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="存储请求(MB)" prop="storage_request">
          <el-input-number
            v-model.number="requestForm.storage_request"
            :min="1"
            placeholder="请输入存储请求"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="请求原因" prop="reason">
          <el-input
            v-model="requestForm.reason"
            placeholder="请简要说明资源请求原因"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitRequestForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElLoading, ElMessageBox } from "element-plus";
import {
  createResourceRequest,
  getResourceRequests,
  approveResourceRequest,
  rejectResourceRequest,
} from "@/api/system";

// 类型定义
interface ResourceRequest {
  id: string;
  tenant_id: string;
  cpu_request: number;
  memory_request: number;
  storage_request: number;
  status: "pending" | "approved" | "rejected";
  reason: string;
  created_at: string;
  approved_at?: string;
  approved_by?: string;
}

// 状态管理
const loading = ref(false);
const searchTenantId = ref("");
const requestList = ref<ResourceRequest[]>([]);
const dialogVisible = ref(false);
const requestFormRef = ref<any>(null);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);

// 表单数据
const requestForm = reactive<{
  tenant_id: string;
  cpu_request: number;
  memory_request: number;
  storage_request: number;
  reason: string;
}>({
  tenant_id: "",
  cpu_request: 0.1,
  memory_request: 128,
  storage_request: 1024,
  reason: "",
});

// 表单验证规则
const requestRules = reactive({
  tenant_id: [{ required: true, message: "请输入租户ID", trigger: "blur" }],
  cpu_request: [
    { required: true, message: "请输入CPU请求", trigger: "blur" },
    {
      type: "number",
      min: 0.1,
      message: "CPU请求不能小于0.1核",
      trigger: "blur",
    },
  ],
  memory_request: [
    { required: true, message: "请输入内存请求", trigger: "blur" },
    { type: "number", min: 1, message: "内存请求不能小于1MB", trigger: "blur" },
  ],
  storage_request: [
    { required: true, message: "请输入存储请求", trigger: "blur" },
    { type: "number", min: 1, message: "存储请求不能小于1MB", trigger: "blur" },
  ],
  reason: [
    { required: true, message: "请输入请求原因", trigger: "blur" },
    { min: 10, message: "请至少输入10个字符", trigger: "blur" },
  ],
});

// 打开创建对话框
const openCreateDialog = () => {
  dialogVisible.value = true;
  if (requestFormRef.value) {
    requestFormRef.value.resetFields();
    // 设置默认值
    requestForm.cpu_request = 0.1;
    requestForm.memory_request = 128;
    requestForm.storage_request = 1024;
  }
};

// 提交资源请求表单
const submitRequestForm = async () => {
  if (!requestFormRef.value) return;

  try {
    await requestFormRef.value.validate();

    const loadingInstance = ElLoading.service({
      lock: true,
      text: "提交中...",
      background: "rgba(0, 0, 0, 0.7)",
    });

    await createResourceRequest(requestForm);
    ElMessage.success("资源请求创建成功");
    dialogVisible.value = false;
    fetchRequestList();
  } catch (error: any) {
    if (error.name === "ValidationError") return;
    ElMessage.error(error.msg || "创建资源请求失败");
  } finally {
    ElLoading.service().close();
  }
};

// 获取资源请求列表
const fetchRequestList = async () => {
  loading.value = true;
  try {
    const response = await getResourceRequests({
      tenant_id: searchTenantId.value,
      page: page.value,
      page_size: pageSize.value,
    });
    requestList.value = response.data.list;
    total.value = response.data.total;
  } catch (error) {
    console.error("获取资源请求列表失败", error);
    ElMessage.error("获取资源请求列表失败");
  } finally {
    loading.value = false;
  }
};

// 处理搜索
const handleSearch = () => {
  page.value = 1;
  fetchRequestList();
};

// 处理分页大小变化
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchRequestList();
};

// 处理页码变化
const handleCurrentChange = (val: number) => {
  page.value = val;
  fetchRequestList();
};

// 处理批准操作
const handleApprove = async (row: ResourceRequest) => {
  try {
    await ElMessageBox.confirm(
      `确定要批准租户ID为 ${row.tenant_id} 的资源请求吗？`,
      "批准确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "success",
      }
    );

    const loadingInstance = ElLoading.service({
      lock: true,
      text: "处理中...",
      background: "rgba(0, 0, 0, 0.7)",
    });

    await approveResourceRequest(row.id);
    ElMessage.success("资源请求已批准");
    fetchRequestList();
  } catch (error: any) {
    if (error.name !== "Error") {
      ElMessage.error(error.msg || "批准资源请求失败");
    }
  } finally {
    ElLoading.service().close();
  }
};

// 处理拒绝操作
const handleReject = async (row: ResourceRequest) => {
  try {
    const { value: reason } = await ElMessageBox.prompt(
      "请输入拒绝原因:",
      "拒绝确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        inputPlaceholder: "请输入拒绝原因",
        inputRules: [
          {
            required: true,
            message: "请输入拒绝原因",
          },
        ],
      }
    );

    const loadingInstance = ElLoading.service({
      lock: true,
      text: "处理中...",
      background: "rgba(0, 0, 0, 0.7)",
    });

    await rejectResourceRequest(row.id, { reason });
    ElMessage.success("资源请求已拒绝");
    fetchRequestList();
  } catch (error: any) {
    if (error.name !== "Error" && error !== "cancel") {
      ElMessage.error(error.msg || "拒绝资源请求失败");
    }
  } finally {
    ElLoading.service().close();
  }
};

// 格式化状态显示
const formatStatus = (status: string): string => {
  const statusMap: Record<string, string> = {
    pending: "待审批",
    approved: "已批准",
    rejected: "已拒绝",
  };
  return statusMap[status] || status;
};

// 获取状态标签类型
const getStatusTagType = (status: string): string => {
  const typeMap: Record<string, string> = {
    pending: "warning",
    approved: "success",
    rejected: "danger",
  };
  return typeMap[status] || "default";
};

// 页面加载时获取列表数据
onMounted(() => {
  fetchRequestList();
});
</script>

<style scoped>
.pagination-container {
  margin-top: 20px;
  text-align: right;
}
.resource-request-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin-bottom: 5px;
  font-size: 20px;
}

.search-and-action {
  margin-bottom: 10px;
}
</style>
