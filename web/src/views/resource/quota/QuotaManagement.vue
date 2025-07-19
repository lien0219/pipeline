<template>
  <div class="quota-management-container">
    <el-card shadow="hover">
      <div class="page-header">
        <h2>资源配额管理</h2>
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
              创建资源配额
            </el-button>
          </el-col>
        </el-row>
      </div>

      <el-table
        v-loading="loading"
        :data="quotaList"
        border
        style="width: 100%; margin-top: 20px"
        empty-text="暂无数据"
      >
        <el-table-column
          prop="tenant_id"
          label="租户ID"
          width="200"
        ></el-table-column>
        <el-table-column prop="cpu_quota" label="CPU配额(核)" width="150">
          <template #default="scope">
            <span>{{ scope.row.cpu_quota.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="memory_quota"
          label="内存配额(MB)"
          width="180"
        ></el-table-column>
        <el-table-column
          prop="storage_quota"
          label="存储配额(MB)"
          width="180"
        ></el-table-column>
        <el-table-column
          prop="created_at"
          label="创建时间"
          width="220"
        ></el-table-column>
        <el-table-column
          prop="updated_at"
          label="更新时间"
          width="220"
        ></el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="openEditDialog(scope.row)"
            >
              编辑
            </el-button>
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(scope.row)"
            >
              删除
            </el-button>
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

    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="quotaFormRef"
        :model="quotaForm"
        :rules="quotaRules"
        label-width="120px"
      >
        <el-form-item label="租户ID" prop="tenant_id" v-if="!isEdit">
          <el-input
            v-model="quotaForm.tenant_id"
            placeholder="请输入租户ID"
            disabled
            v-if="isEdit"
          />
          <el-input
            v-model="quotaForm.tenant_id"
            placeholder="请输入租户ID"
            v-else
          />
        </el-form-item>
        <el-form-item label="CPU配额(核)" prop="cpu_quota">
          <el-input-number
            v-model.number="quotaForm.cpu_quota"
            :min="0"
            :step="0.1"
            placeholder="请输入CPU配额"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="内存配额(MB)" prop="memory_quota">
          <el-input-number
            v-model.number="quotaForm.memory_quota"
            :min="0"
            placeholder="请输入内存配额"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="存储配额(MB)" prop="storage_quota">
          <el-input-number
            v-model.number="quotaForm.storage_quota"
            :min="0"
            placeholder="请输入存储配额"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitQuotaForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElLoading, ElMessageBox } from "element-plus";
import {
  createResourceQuota,
  getResourceQuota,
  updateResourceQuota,
  deleteResourceQuota,
} from "@/api/system";

// 类型定义
interface ResourceQuota {
  tenant_id: string;
  cpu_quota: number;
  memory_quota: number;
  storage_quota: number;
  created_at: string;
  updated_at: string;
}

// 状态管理
const loading = ref(false);
const searchTenantId = ref("");
const quotaList = ref<ResourceQuota[]>([]);
const dialogVisible = ref(false);
const dialogTitle = ref("创建资源配额");
const isEdit = ref(false);
const quotaFormRef = ref<any>(null);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);

// 表单数据
const quotaForm = reactive<{
  tenant_id: string;
  cpu_quota: number;
  memory_quota: number;
  storage_quota: number;
}>({
  tenant_id: "",
  cpu_quota: 0,
  memory_quota: 0,
  storage_quota: 0,
});

// 表单验证规则
const quotaRules = reactive({
  tenant_id: [{ required: true, message: "请输入租户ID", trigger: "blur" }],
  cpu_quota: [
    { required: true, message: "请输入CPU配额", trigger: "blur" },
    { type: "number", min: 0, message: "CPU配额不能小于0", trigger: "blur" },
  ],
  memory_quota: [
    { required: true, message: "请输入内存配额", trigger: "blur" },
    { type: "number", min: 0, message: "内存配额不能小于0", trigger: "blur" },
  ],
  storage_quota: [
    { required: true, message: "请输入存储配额", trigger: "blur" },
    { type: "number", min: 0, message: "存储配额不能小于0", trigger: "blur" },
  ],
});

const openCreateDialog = () => {
  isEdit.value = false;
  dialogTitle.value = "创建资源配额";
  quotaForm.tenant_id = "";
  quotaForm.cpu_quota = 0;
  quotaForm.memory_quota = 0;
  quotaForm.storage_quota = 0;
  dialogVisible.value = true;
  if (quotaFormRef.value) {
    quotaFormRef.value.resetFields();
  }
};

const openEditDialog = (row: ResourceQuota) => {
  isEdit.value = true;
  dialogTitle.value = "编辑资源配额";
  quotaForm.tenant_id = row.tenant_id;
  quotaForm.cpu_quota = row.cpu_quota;
  quotaForm.memory_quota = row.memory_quota;
  quotaForm.storage_quota = row.storage_quota;
  dialogVisible.value = true;
};

// 提交表单
const submitQuotaForm = async () => {
  if (!quotaFormRef.value) return;

  try {
    await quotaFormRef.value.validate();

    const loadingInstance = ElLoading.service({
      lock: true,
      text: isEdit.value ? "更新中..." : "创建中...",
      background: "rgba(0, 0, 0, 0.7)",
    });

    if (isEdit.value) {
      await updateResourceQuota(quotaForm.tenant_id, {
        cpu_quota: quotaForm.cpu_quota,
        memory_quota: quotaForm.memory_quota,
        storage_quota: quotaForm.storage_quota,
      });
      ElMessage.success("资源配额更新成功");
    } else {
      await createResourceQuota(quotaForm);
      ElMessage.success("资源配额创建成功");
    }

    dialogVisible.value = true;
    if (!isEdit.value) {
      searchTenantId.value = quotaForm.tenant_id;
      handleSearch();
    } else {
      handleSearch();
    }

    dialogVisible.value = false;
  } catch (error: any) {
    if (error.name === "ValidationError") return;
    ElMessage.error(
      error.msg || (isEdit.value ? "更新资源配额失败" : "创建资源配额失败")
    );
  } finally {
    ElLoading.service().close();
  }
};

// 获取资源配额列表
const fetchQuotaList = async () => {
  loading.value = true;
  try {
    const response = await getResourceQuota({
      tenant_id: searchTenantId.value,
      page: page.value,
      page_size: pageSize.value,
    });
    quotaList.value = response.data.list;
    total.value = response.data.total;
  } catch (error) {
    console.error("获取资源配额列表失败", error);
  } finally {
    loading.value = false;
  }
};
const handleDelete = async (row: ResourceQuota) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除租户ID为 ${row.tenant_id} 的资源配额吗？`,
      "删除确认",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }
    );

    const loadingInstance = ElLoading.service({
      lock: true,
      text: "删除中...",
      background: "rgba(0, 0, 0, 0.7)",
    });

    await deleteResourceQuota(row.tenant_id);
    ElMessage.success("资源配额删除成功");
    fetchQuotaList();
  } catch (error: any) {
    if (error.name !== "Error") {
      ElMessage.error(error.msg || "删除资源配额失败");
    }
  } finally {
    ElLoading.service().close();
  }
};
// 搜索
const handleSearch = () => {
  page.value = 1;
  fetchQuotaList();
};

// 分页大小变化
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchQuotaList();
};

// 当前页码变化
const handleCurrentChange = (val: number) => {
  page.value = val;
  fetchQuotaList();
};

onMounted(() => {
  fetchQuotaList();
});
</script>

<style scoped>
.pagination-container {
  margin-top: 20px;
  text-align: right;
}
.quota-management-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin-bottom: 5px;
  font-size: 20px;
}

.page-header p {
  color: #606266;
  margin: 0;
}

.search-and-action {
  margin-bottom: 10px;
}
</style>
