<template>
  <div class="app-container">
    <!-- 页面标题和操作按钮 -->
    <div class="header-actions">
      <h1>环境管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon-plus /> 创建环境
      </el-button>
    </div>

    <!-- 筛选条件 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="environmentStore.filters.name"
            placeholder="环境名称"
            clearable
            @clear="handleFilterChange"
            @input="handleFilterChange"
          />
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="environmentStore.filters.type"
            placeholder="环境类型"
            clearable
            @change="handleFilterChange"
          >
            <el-option label="开发环境" value="development" />
            <el-option label="测试环境" value="testing" />
            <el-option label="生产环境" value="production" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="environmentStore.filters.status"
            placeholder="环境状态"
            clearable
            @change="handleFilterChange"
            :model-value="environmentStore.filters.status"
          >
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="inactive" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="default" @click="resetFilter">重置筛选</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 环境列表表格 -->
    <el-card class="mt-4">
      <el-table
        v-loading="environmentStore.loading"
        :data="environmentStore.environments || []"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="name" label="环境名称" />
        <el-table-column prop="type" label="环境类型">
          <template #default="scope">
            <el-tag
              :type="
                scope.row.type === 'production'
                  ? 'danger'
                  : scope.row.type === 'testing'
                  ? 'warning'
                  : 'success'
              "
            >
              {{
                scope.row.type === "production"
                  ? "生产环境"
                  : scope.row.type === "testing"
                  ? "测试环境"
                  : "开发环境"
              }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="url" label="环境URL" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-switch
              v-model="scope.row.status"
              active-value="active"
              inactive-value="inactive"
              @change="handleStatusChange(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row)"
              >查看</el-button
            >
            <el-button
              size="small"
              type="primary"
              @click="handleEdit(scope.row)"
              >编辑</el-button
            >
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(scope.row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页控件 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="environmentStore.filters.page"
          v-model:page-size="environmentStore.filters.pageSize"
          :total="environmentStore.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePageSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑环境弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '创建环境' : '编辑环境'"
      :width="'600px'"
    >
      <el-form
        ref="environmentFormRef"
        :model="environmentForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="环境名称" prop="name">
          <el-input
            v-model="environmentForm.name"
            placeholder="请输入环境名称"
          />
        </el-form-item>
        <el-form-item label="环境类型" prop="type">
          <el-select
            v-model="environmentForm.type"
            placeholder="请选择环境类型"
          >
            <el-option label="开发环境" value="development" />
            <el-option label="测试环境" value="testing" />
            <el-option label="生产环境" value="production" />
          </el-select>
        </el-form-item>
        <el-form-item label="环境URL" prop="url">
          <el-input v-model="environmentForm.url" placeholder="请输入环境URL" />
        </el-form-item>
        <el-form-item label="环境描述" prop="description">
          <el-input
            v-model="environmentForm.description"
            placeholder="请输入环境描述"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="环境状态" prop="status">
          <el-switch
            v-model="environmentForm.status"
            active-value="active"
            inactive-value="inactive"
          />
        </el-form-item>
        <el-form-item label="环境变量">
          <div
            v-for="(variable, index) in environmentForm.variables"
            :key="index"
            class="variable-item"
          >
            <el-input
              v-model="variable.key"
              placeholder="变量名"
              style="width: 45%; margin-right: 10px"
            />
            <el-input
              v-model="variable.value"
              placeholder="变量值"
              style="width: 45%"
            />
            <el-button
              type="text"
              icon="el-icon-delete"
              @click="removeVariable(index)"
              v-if="index > 0"
            />
          </div>
          <el-button type="text" @click="addVariable">+ 添加环境变量</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">
          <el-icon-loading v-if="submitting" class="loading-icon" />
          {{ submitting ? "提交中..." : "提交" }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 环境详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="环境详情"
      :width="'600px'"
      style="overflow: hidden"
    >
      <el-descriptions :column="1" border v-if="currentEnvironment">
        <el-descriptions-item label="环境名称">{{
          currentEnvironment.name || "-"
        }}</el-descriptions-item>
        <el-descriptions-item label="环境类型">
          <el-tag
            :type="
              currentEnvironment.type === 'production'
                ? 'danger'
                : currentEnvironment.type === 'testing'
                ? 'warning'
                : 'success'
            "
          >
            {{
              currentEnvironment.type === "production"
                ? "生产环境"
                : currentEnvironment.type === "testing"
                ? "测试环境"
                : "开发环境"
            }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="环境URL">{{
          currentEnvironment.url || "-"
        }}</el-descriptions-item>
        <el-descriptions-item label="环境状态">{{
          currentEnvironment.status === "active" ? "启用" : "禁用"
        }}</el-descriptions-item>
        <el-descriptions-item label="环境描述">
          <div style="word-break: break-all; word-wrap: break-word">
            {{ currentEnvironment.description || "-" }}
          </div></el-descriptions-item
        >
        <el-descriptions-item label="创建时间">{{
          formatDate(currentEnvironment.created_at) || "-"
        }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{
          formatDate(currentEnvironment.updated_at) || "-"
        }}</el-descriptions-item>
        <el-descriptions-item label="环境变量">
          <el-table
            :data="currentEnvironment.variables || []"
            border
            style="width: 100%; max-width: 500px"
            :cell-style="{
              'white-space': 'normal',
              'word-wrap': 'break-word',
              'word-break': 'break-all',
              padding: '8px',
            }"
          >
            <el-table-column prop="key" label="变量名" />
            <el-table-column prop="value" label="变量值" />
          </el-table>
        </el-descriptions-item>
      </el-descriptions>
      <div v-else>加载环境详情失败</div>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
// @ts-ignore
import { useEnvironmentStore } from "@/stores/environment";
// @ts-ignore
import { formatDate } from "@/utils/date";

// 状态管理
const environmentStore = useEnvironmentStore();
const statusUpdating = ref<string>("");

// 弹窗控制
const dialogVisible = ref(false);
const detailVisible = ref(false);
const dialogType = ref<"create" | "edit">("create");
const submitting = ref(false);
const environmentFormRef = ref();

// 当前选中的环境
const currentEnvironment = ref<any | null>(null);

// 表单数据
const environmentForm = reactive<any>({
  name: "",
  type: "",
  url: "",
  description: "",
  status: "active",
  variables: [] as Array<{ key: string; value: string }>,
});
const environmentStatus = ref<Record<number, string>>({});

// 表单验证规则
const formRules = reactive({
  name: [
    { required: true, message: "请输入环境名称", trigger: ["blur", "input"] },
  ],
  type: [{ required: true, message: "请选择环境类型", trigger: "change" }],
  url: [
    { required: true, message: "请输入环境URL", trigger: ["blur", "input"] },
    { type: "url", message: "请输入有效的URL", trigger: ["blur", "input"] },
  ],
});

onMounted(async () => {
  try {
    await environmentStore.fetchEnvironments();
    environmentStore.environments.forEach((env) => {
      environmentStatus.value[env.id] = env.status;
    });
  } catch (error) {
    console.error("Failed to fetch environments:", error);
    environmentStore.environments = [];
  }
});

// 筛选条件变化
const handleFilterChange = () => {
  // 重置到第一页
  environmentStore.filters.page = 1;
  console.log("筛选状态值:", environmentStore.filters.status);
  environmentStore.fetchEnvironments();
};

// 重置筛选条件
const resetFilter = () => {
  environmentStore.filters = {
    name: "",
    type: "",
    status: "",
    page: 1,
    pageSize: 10,
  };
  environmentStore.fetchEnvironments();
};

// 分页大小变化
const handlePageSizeChange = (pageSize: number) => {
  environmentStore.pagination.pageSize = pageSize;
  environmentStore.fetchEnvironments();
};

// 页码变化
const handlePageChange = (page: number) => {
  environmentStore.pagination.page = page;
  environmentStore.fetchEnvironments();
};

// 打开创建环境弹窗
const handleCreate = () => {
  dialogType.value = "create";
  // 重置表单
  Object.assign(environmentForm, {
    name: "",
    type: "",
    url: "",
    description: "",
    status: "active",
    variables: [{ key: "", value: "" }],
  });
  dialogVisible.value = true;
};

// 打开编辑环境弹窗
const handleEdit = async (environment: any) => {
  dialogType.value = "edit";
  currentEnvironment.value = environment;
  try {
    const detail = await environmentStore.fetchEnvironmentById(environment.id);
    if (!detail) {
      ElMessage.error("获取环境详情失败: 数据为空");
      return;
    }
    let variables = [{ key: "", value: "" }];
    if (detail.variables) {
      if (typeof detail.variables === "string") {
        try {
          variables = JSON.parse(detail.variables);
          if (!Array.isArray(variables)) {
            throw new Error("解析结果不是数组");
          }
          variables = variables.filter(
            (v) => v && typeof v === "object" && "key" in v
          );
        } catch (e) {
          ElMessage.error(`变量解析失败: ${(e as Error).message}`);
          variables = [{ key: "", value: "" }];
        }
      } else if (Array.isArray(detail.variables)) {
        variables = detail.variables;
      }
    }
    Object.assign(environmentForm, {
      name: detail.name,
      type: detail.type,
      url: detail.url,
      description: detail.description,
      status: detail.status,
      variables: variables,
    });
    dialogVisible.value = true;
  } catch (error) {
    ElMessage.error("加载环境详情失败");
    currentEnvironment.value = null;
  }
};
// 查看环境详情
const handleView = async (environment: any) => {
  try {
    const detail = await environmentStore.fetchEnvironmentById(environment.id);
    if (!detail) {
      ElMessage.error("获取环境详情失败: 数据为空");
      return;
    }
    currentEnvironment.value = detail;
    detailVisible.value = true;
  } catch (error) {
    ElMessage.error("加载环境详情失败");
    currentEnvironment.value = null;
  }
};

// 删除环境
const handleDelete = async (environment: any) => {
  try {
    await ElMessageBox.confirm(
      "确定要删除这个环境吗？此操作不可撤销。",
      "警告",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      }
    );

    await environmentStore.deleteEnvironment(environment.id);
    ElMessage.success("环境删除成功");
    environmentStore.fetchEnvironments();
  } catch (error) {
    // 如果用户取消，则不显示错误信息
    if (error !== "cancel") {
      ElMessage.error("环境删除失败");
    }
  }
};

// 添加环境变量
const addVariable = () => {
  environmentForm.variables.push({ key: "", value: "" });
};

// 移除环境变量
const removeVariable = (index: number) => {
  environmentForm.variables.splice(index, 1);
};

// 切换环境状态
const handleStatusChange = async (environment: any) => {
  const currentStatus = environmentStatus.value[environment.id];
  const newStatus = currentStatus === "active" ? "inactive" : "active";
  environmentStatus.value[environment.id] = newStatus;

  const originalStatus = newStatus;

  try {
    await environmentStore.updateEnvironment(environment.id, {
      name: environment.name,
      type: environment.type,
      status: environment.status,
    });
  } catch (error) {
    environment.status = originalStatus;
    ElMessage.error(
      "环境状态更新失败: " +
        (error instanceof Error ? error.message : String(error))
    );
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!environmentFormRef.value) return;

  try {
    await environmentFormRef.value.validate();
    submitting.value = true;

    if (!Array.isArray(environmentForm.variables)) {
      environmentForm.variables = [];
    }
    const validVariables = environmentForm.variables.filter(
      (v) => v.key && v.key.trim() !== ""
    );

    if (dialogType.value === "create") {
      // 创建环境
      await environmentStore.createEnvironment({
        ...environmentForm,
        variables: validVariables,
      });
    } else {
      // 编辑环境
      await environmentStore.updateEnvironment(currentEnvironment.value.id, {
        ...environmentForm,
        variables: validVariables,
      });
    }

    dialogVisible.value = false;
    ElMessage.success(
      `${dialogType.value === "create" ? "创建" : "更新"}环境成功`
    );
  } catch (error) {
    if (typeof error === "string") {
      ElMessage.error(error);
    } else {
      console.error("提交失败:", error);
      ElMessage.error(
        `${dialogType.value === "create" ? "创建" : "更新"}环境失败`
      );
    }
  } finally {
    submitting.value = false;
  }
};
</script>

<style scoped>
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 16px;
  text-align: right;
}

.variable-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.loading-icon {
  margin-right: 5px;
}
</style>
