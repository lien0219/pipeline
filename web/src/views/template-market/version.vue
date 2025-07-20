<template>
  <div class="version-management-container">
    <div class="page-header">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item>版本管理</el-breadcrumb-item>
      </el-breadcrumb>
      <el-button type="primary" @click="handleCreateVersion">
        <el-icon-plus /> 新建版本
      </el-button>
    </div>

    <!-- 版本列表 -->
    <el-card class="version-table-card">
      <el-table v-loading="loading" :data="versions" border>
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="version" label="版本号" width="120" />
        <el-table-column prop="changelog" label="更新日志" width="250" />
        <el-table-column prop="download_count" label="下载次数" width="120" />
        <el-table-column prop="is_latest" label="是否最新" width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.is_latest" type="success">是</el-tag>
            <el-tag v-else type="info">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="scope">
            <el-button
              type="text"
              size="small"
              @click="handleViewContent(scope.row)"
              v-if="scope.row.content"
            >
              查看内容
            </el-button>
            <el-button
              type="text"
              size="small"
              @click="handleSetLatest(scope.row)"
              v-if="!scope.row.is_latest"
            >
              设为最新
            </el-button>
            <el-button
              type="text"
              size="small"
              danger
              @click="handleDeleteVersion(scope.row.id)"
              v-if="!scope.row.is_latest"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 创建/编辑版本对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建版本' : '编辑版本'"
      :width="'600px'"
    >
      <el-form
        ref="versionFormRef"
        :model="versionForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="版本号" prop="version">
          <el-input
            v-model="versionForm.version"
            placeholder="请输入语义化版本号，如1.0.0"
          />
          <div class="form-hint">
            版本号格式应为x.y.z，其中x、y、z为非负整数，且禁止在数字前补零
          </div>
        </el-form-item>
        <el-form-item label="更新日志" prop="changelog">
          <el-input
            v-model="versionForm.changelog"
            placeholder="请输入更新日志"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <el-input
            v-model="versionForm.content"
            placeholder="请输入模板内容"
            type="textarea"
            :rows="6"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>

    <!-- 版本内容查看对话框 -->
    <el-dialog v-model="contentDialogVisible" title="模板内容" :width="'800px'">
      <el-input v-model="currentContent" type="textarea" :rows="20" readonly />
      <template #footer>
        <el-button @click="contentDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  getTemplateById,
  getTemplateVersions,
  createTemplateVersion,
  deleteTemplateVersion,
  setVersionAsLatest,
} from "@/api/templateMarket";
import type { Template, TemplateVersion } from "../..//api/types";
import { ElMessage, ElMessageBox } from "element-plus";
import semver from "semver";

const route = useRoute();
const router = useRouter();
const templateId = ref<number>(Number(route.query.templateId));
const templateName = ref<string>("");
const loading = ref<boolean>(false);
const versions = ref<TemplateVersion[]>([]);

// 对话框状态
const dialogVisible = ref<boolean>(false);
const contentDialogVisible = ref<boolean>(false);
const dialogType = ref<"create" | "edit">("create");
const currentContent = ref<string>("");

// 表单数据
const versionForm = ref<any>({
  id: 0,
  version: "",
  changelog: "",
  content: "",
});
const versionFormRef = ref<InstanceType<any>>();

// 表单验证规则
const formRules = ref({
  version: [
    { required: true, message: "请输入版本号", trigger: "blur" },
    {
      validator: (rule: any, value: string, callback: any) => {
        if (semver.valid(value)) {
          callback();
        } else {
          callback(new Error("请输入有效的语义化版本号，如1.0.0"));
        }
      },
      trigger: "blur",
    },
  ],
  content: [{ required: true, message: "请输入模板内容", trigger: "blur" }],
});

// 获取模板信息
const fetchTemplateInfo = async () => {
  try {
    const res = await getTemplateById(templateId.value);
    templateName.value = res.data.name;
  } catch (error) {
    ElMessage.error("获取模板信息失败");
    console.error(error);
  }
};

// 获取版本列表
const fetchVersions = async () => {
  try {
    loading.value = true;
    const res = await getTemplateVersions(templateId.value);
    versions.value = res.data || [];
  } catch (error) {
    ElMessage.error("获取版本列表失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 创建版本
const handleCreateVersion = () => {
  dialogType.value = "create";
  versionForm.value = {
    id: 0,
    version: "",
    changelog: "",
    content: "",
  };
  dialogVisible.value = true;
};

// 查看版本内容
const handleViewContent = (row: TemplateVersion) => {
  currentContent.value = row.content || "";
  contentDialogVisible.value = true;
};

// 设置为最新版本
const handleSetLatest = async (row: TemplateVersion) => {
  try {
    await ElMessageBox.confirm("确定要将此版本设为最新版本吗？", "确认操作", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await setVersionAsLatest(templateId.value, row.id);
    ElMessage.success("设置成功");
    fetchVersions();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("设置失败");
      console.error(error);
    }
  }
};

// 删除版本
const handleDeleteVersion = async (versionId: number) => {
  try {
    await ElMessageBox.confirm("确定要删除此版本吗？此操作不可恢复！", "警告", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "error",
    });

    await deleteTemplateVersion(templateId.value, versionId);
    ElMessage.success("删除成功");
    fetchVersions();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("删除失败");
      console.error(error);
    }
  }
};

// 提交表单
const submitForm = async () => {
  try {
    await versionFormRef.value.validate();
    loading.value = true;

    if (dialogType.value === "create") {
      await createTemplateVersion(templateId.value, versionForm.value);
      ElMessage.success("创建版本成功");
    }

    dialogVisible.value = false;
    fetchVersions();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error(
        dialogType.value === "create" ? "创建版本失败" : "更新版本失败"
      );
      console.error(error);
    }
  } finally {
    loading.value = false;
  }
};

// 初始化
onMounted(async () => {
  if (!templateId.value) {
    ElMessage.error("模板ID不存在");
    return;
  }

  await Promise.all([fetchTemplateInfo(), fetchVersions()]);
});

// 监听模板ID变化
watch(
  () => route.params.templateId,
  (newVal) => {
    templateId.value = Number(newVal);
    fetchTemplateInfo();
    fetchVersions();
  }
);
</script>

<style scoped>
.version-management-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.version-table-card {
  margin-bottom: 20px;
}

.form-hint {
  color: #606266;
  font-size: 12px;
  margin-top: 5px;
}

.pagination-container {
  margin-top: 16px;
  text-align: right;
}
</style>
