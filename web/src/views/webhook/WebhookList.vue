<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>Webhook管理</h2>
      </div>
      <el-button type="primary" @click="openCreateDialog">
        <el-icon><Plus /></el-icon> 创建Webhook
      </el-button>
    </div>

    <el-card>
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="Webhook名称">
          <el-input v-model="searchForm.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="状态" style="width: 180px">
          <el-select v-model="searchForm.status" placeholder="请选择状态">
            <el-option label="全部" value="" />
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchWebhooks">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table
        v-loading="loading"
        :data="webhookList"
        border
        fit
        highlight-current-row
      >
        <el-table-column type="index" label="序号" width="80" align="center" />
        <el-table-column prop="name" label="名称" align="center" />
        <el-table-column prop="pipeline_id" label="流水线ID" align="center" />
        <el-table-column
          width="200"
          prop="url"
          label="URL"
          align="center"
          :show-overflow-tooltip="true"
        />
        <el-table-column label="事件" align="center">
          <template #default="scope">
            <el-tag
              v-for="event in scope.row.events.split(',')"
              :key="event"
              size="small"
              style="margin-right: 5px"
            >
              {{ event }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_active" label="状态" align="center">
          <template #default="scope">
            <el-switch
              disabled="true"
              v-model="scope.row.is_active"
              @change="handleStatusChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="created_at"
          label="创建时间"
          align="center"
          width="180"
        />
        <el-table-column label="操作" align="center" width="180">
          <template #default="scope">
            <el-button
              type="text"
              @click="openEditDialog(scope.row)"
              disabled="true"
              >编辑</el-button
            >
            <el-button
              type="text"
              danger
              @click="handleDelete(scope.row)"
              disabled="true"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑Webhook对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form
        ref="webhookFormRef"
        :model="webhookForm"
        :rules="webhookRules"
        label-width="120px"
      >
        <el-form-item label="Webhook名称" prop="name">
          <el-input
            v-model="webhookForm.name"
            placeholder="请输入Webhook名称"
          />
        </el-form-item>
        <el-form-item label="所属流水线ID" prop="pipeline_id">
          <el-input
            v-model="webhookForm.pipeline_id"
            placeholder="请输入所属流水线ID"
            type="number"
          />
        </el-form-item>
        <el-form-item label="URL地址" prop="url">
          <el-input
            v-model="webhookForm.url"
            placeholder="请输入Webhook URL地址"
          />
        </el-form-item>
        <el-form-item label="密钥" prop="secret">
          <el-input
            v-model="webhookForm.secret"
            placeholder="请输入签名密钥（可选）"
          />
        </el-form-item>
        <el-form-item label="触发事件" prop="events">
          <el-select
            v-model="webhookForm.events"
            multiple
            placeholder="请选择触发事件"
          >
            <el-option label="流水线开始" value="pipeline.start" />
            <el-option label="流水线成功" value="pipeline.success" />
            <el-option label="流水线失败" value="pipeline.failure" />
            <el-option label="任务开始" value="task.start" />
            <el-option label="任务成功" value="task.success" />
            <el-option label="任务失败" value="task.failure" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用状态" prop="is_active">
          <el-switch v-model="webhookForm.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitWebhook">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox, ElLoading } from "element-plus";
import { Plus, Edit, Delete } from "@element-plus/icons-vue";
import { webhookApi } from "@/api/webhook";

// 状态管理
const loading = ref(false);
const webhookList = ref([]);
const total = ref(0);
const dialogVisible = ref(false);
const dialogTitle = ref("创建Webhook");
const currentWebhookId = ref(0);
const webhookFormRef = ref(null);

// 搜索表单
const searchForm = reactive({
  pipelineId: "",
  name: "",
  status: "",
});

// 分页配置
const pagination = reactive({
  page: 1,
  pageSize: 10,
});

// Webhook表单
const webhookForm = reactive({
  name: "",
  url: "",
  secret: "",
  events: [],
  is_active: true,
  pipeline_id: 0,
});

// 表单验证规则
const webhookRules = ref({
  name: [{ required: true, message: "请输入Webhook名称", trigger: "blur" }],
  url: [
    { required: true, message: "请输入URL地址", trigger: "blur" },
    { type: "url", message: "请输入有效的URL地址", trigger: "blur" },
  ],
  events: [
    { required: true, message: "请至少选择一个触发事件", trigger: "change" },
  ],
});

// 获取Webhook列表
const fetchWebhooks = async () => {
  try {
    loading.value = true;
    const response = await webhookApi.getWebhooks({
      page: pagination.page,
      pageSize: pagination.pageSize,
      pipelineId: searchForm.pipelineId || undefined,
      name: searchForm.name,
      status: searchForm.status !== "" ? searchForm.status : undefined,
    });
    webhookList.value = response.data.list;
    total.value = response.data.total;
  } catch (error) {
    ElMessage.error("获取Webhook列表失败: " + (error.message || "未知错误"));
  } finally {
    loading.value = false;
  }
};
// 创建/编辑Webhook
const submitWebhook = async () => {
  try {
    const form = webhookFormRef.value;
    if (!form) return;

    await form.validate();
    const loadingInstance = ElLoading.service({ text: "提交中..." });

    const submitData = { ...webhookForm };
    submitData.events = submitData.events.join(",");
    submitData.pipeline_id = Number(submitData.pipeline_id);

    if (currentWebhookId.value) {
      // 更新Webhook
      await webhookApi.updateWebhook(currentWebhookId.value, submitData);
      ElMessage.success("Webhook更新成功");
    } else {
      // 创建Webhook
      await webhookApi.createWebhook(submitData);
      ElMessage.success("Webhook创建成功");
    }

    dialogVisible.value = false;
    fetchWebhooks();
    loadingInstance.close();
  } catch (error) {
    if (error.name !== "ValidationError") {
      ElMessage.error("操作失败: " + (error.message || "未知错误"));
    }
  }
};

// 打开创建对话框
const openCreateDialog = () => {
  resetForm();
  dialogTitle.value = "创建Webhook";
  currentWebhookId.value = 0;
  dialogVisible.value = true;
};

// 打开编辑对话框
const openEditDialog = (row) => {
  resetForm();
  dialogTitle.value = "编辑Webhook";
  currentWebhookId.value = row.id;
  // 填充表单数据
  webhookForm.name = row.name;
  webhookForm.url = row.url;
  webhookForm.secret = row.secret;
  webhookForm.events = row.events.split(",");
  webhookForm.is_active = row.is_active;
  dialogVisible.value = true;
};

// 删除Webhook
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm("确定要删除该Webhook吗？", "确认删除", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await webhookApi.deleteWebhook(row.id);
    ElMessage.success("删除成功");
    fetchWebhooks();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("删除失败: " + (error.message || "未知错误"));
    }
  }
};

// 重置表单
const resetForm = () => {
  const form = webhookFormRef.value;
  if (form) {
    form.resetFields();
  }
  webhookForm.name = "";
  webhookForm.url = "";
  webhookForm.secret = "";
  webhookForm.events = [];
  webhookForm.is_active = true;
};

// 分页相关方法
const handleSizeChange = (val) => {
  pagination.pageSize = val;
  fetchWebhooks();
};

const handleCurrentChange = (val) => {
  pagination.page = val;
  fetchWebhooks();
};

const resetSearch = () => {
  searchForm.pipelineId = "";
  searchForm.name = "";
  searchForm.status = "";
  fetchWebhooks();
};

onMounted(() => {
  fetchWebhooks();
});
</script>

<style scoped>
.pagination-container {
  text-align: right;
  margin-top: 16px;
}

.search-form {
  margin-bottom: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
