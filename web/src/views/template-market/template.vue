<template>
  <div class="template-management-container">
    <div class="page-header">
      <h2>模板管理</h2>
      <el-button type="primary" @click="handleCreate">
        <el-icon-plus /> 新建模板
      </el-button>
    </div>

    <!-- 搜索和筛选区 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-input
            v-model="searchForm.keyword"
            placeholder="输入模板名称搜索"
            prefix-icon="Search"
          />
        </el-col>
        <el-col :span="8">
          <el-select
            v-model="searchForm.category_id"
            placeholder="选择分类"
            filterable
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-col>
        <el-col :span="8">
          <el-select v-model="searchForm.public" placeholder="模板状态">
            <el-option label="全部" value="" />
            <el-option label="公开" :value="true" />
            <el-option label="私有" :value="false" />
          </el-select>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 15px">
        <el-col :span="24" class="text-right">
          <el-button @click="resetFilter">重置</el-button>
          <el-button type="primary" @click="fetchTemplates">查询</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 模板列表 -->
    <el-card class="template-table-card">
      <el-table v-loading="loading" :data="templates" border>
        <el-table-column type="index" label="序号" width="80" />
        <el-table-column prop="name" label="模板名称" width="120" />
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column prop="description" label="描述" width="300" />
        <el-table-column
          prop="is_public"
          label="状态"
          width="100"
          :formatter="formatStatus"
        />
        <el-table-column prop="download_count" label="下载次数" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="300" fixed="right">
          <template #default="scope">
            <el-button
              type="text"
              size="small"
              @click="handleVersion(scope.row.id)"
            >
              版本管理
            </el-button>
            <el-button
              type="text"
              size="small"
              @click="handleDetail(scope.row)"
            >
              详情
            </el-button>
            <el-button type="text" size="small" @click="handleEdit(scope.row)">
              编辑
            </el-button>
            <el-button
              type="text"
              size="small"
              danger
              @click="handleDelete(scope.row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          background
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑模板对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建模板' : '编辑模板'"
      :width="'600px'"
    >
      <el-form
        ref="templateFormRef"
        :model="templateForm"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="模板名称" prop="name">
          <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
        </el-form-item>
        <el-form-item label="所属分类" prop="category_id">
          <el-select
            v-model="templateForm.category_id"
            placeholder="请选择分类"
          >
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="模板描述" prop="description">
          <el-input
            v-model="templateForm.description"
            placeholder="请输入模板描述"
            type="textarea"
            :rows="4"
          />
        </el-form-item>
        <el-form-item label="标签" prop="tags">
          <el-input
            v-model="templateForm.tags"
            placeholder="请输入标签，多个标签用逗号分隔"
          />
        </el-form-item>
        <el-form-item label="版本号" prop="version">
          <el-input
            v-model="templateForm.version"
            placeholder="请输入版本号，如1.0.0"
          />
        </el-form-item>
        <el-form-item label="模板内容" prop="content">
          <el-input
            v-model="templateForm.content"
            placeholder="请输入模板内容"
            type="textarea"
            :rows="6"
          />
        </el-form-item>
        <el-form-item label="是否公开" prop="public">
          <el-switch v-model="templateForm.public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
    <!-- 模板详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="模板详情" :width="'600px'">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="模板名称">{{
          templateDetail.name
        }}</el-descriptions-item>
        <el-descriptions-item label="所属分类">{{
          templateDetail.category?.name || "-"
        }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{
          templateDetail.description
        }}</el-descriptions-item>
        <el-descriptions-item label="标签">{{
          templateDetail.tags
        }}</el-descriptions-item>
        <el-descriptions-item label="版本号">{{
          templateDetail.versions[0].version
        }}</el-descriptions-item>
        <el-descriptions-item label="是否公开">{{
          templateDetail.is_public ? "公开" : "私有"
        }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{
          templateDetail.created_at
        }}</el-descriptions-item>
        <el-descriptions-item label="模板内容">
          <el-input
            v-model="templateDetail.versions[0].content"
            type="textarea"
            :rows="6"
            readonly
          />
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, getCurrentInstance } from "vue";
import { useRouter } from "vue-router";
import {
  getTemplates,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  getTemplateCategories,
  getTemplateById,
} from "@/api/templateMarket";
import type { Template, TemplateCategory } from "../../api/types";
import { ElMessage, ElMessageBox } from "element-plus";

const router = useRouter();
const loading = ref<boolean>(false);
const dialogVisible = ref<boolean>(false);
const dialogType = ref<"create" | "edit">("create");
const templateForm = ref<any>({
  id: 0,
  name: "",
  category_id: 0,
  description: "",
  tags: "",
  public: true,
  version: "",
  content: "",
});
const templateFormRef = ref<InstanceType<any>>();
const searchForm = ref<any>({
  keyword: "",
  category_id: "",
  public: "",
});
const categories = ref<TemplateCategory[]>([]);
const templates = ref<Template[]>([]);
const pagination = ref({
  current: 1,
  size: 10,
  total: 0,
});

const detailDialogVisible = ref<boolean>(false);
const templateDetail = ref<any>({});

// 表单验证规则
const formRules = ref({
  name: [
    { required: true, message: "请输入模板名称", trigger: "blur" },
    {
      min: 2,
      max: 50,
      message: "模板名称长度在 2 到 50 个字符",
      trigger: "blur",
    },
  ],
  category_id: [{ required: true, message: "请选择分类", trigger: "change" }],
  description: [{ required: true, message: "请输入模板描述", trigger: "blur" }],
  version: [
    { required: true, message: "请输入版本号", trigger: "blur" },
    {
      pattern: /^\d+\.\d+\.\d+$/,
      message: "版本号格式应为x.y.z",
      trigger: "blur",
    },
  ],
  content: [{ required: true, message: "请输入模板内容", trigger: "blur" }],
});

// 获取分类列表
const fetchCategories = async () => {
  try {
    const res = await getTemplateCategories();
    categories.value = res.data.list;
  } catch (error) {
    ElMessage.error("获取分类列表失败");
    console.error(error);
  }
};

// 获取模板列表
const fetchTemplates = async () => {
  try {
    loading.value = true;
    const params = {
      ...searchForm.value,
      page: pagination.value.current,
      pageSize: pagination.value.size,
    };
    const res: any = await getTemplates(params);
    const responseData = res.data || { list: [], total: 0 };
    templates.value = responseData.list || [];
    pagination.value.total = responseData.total || 0;
  } catch (error) {
    ElMessage.error("获取模板列表失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 格式化状态显示
const formatStatus = (row: any) => {
  return row.is_public ? "公开" : "私有";
};

// 分页相关方法
const handleSizeChange = (val: number) => {
  pagination.value.size = val;
  fetchTemplates();
};

const handleCurrentChange = (val: number) => {
  pagination.value.current = val;
  fetchTemplates();
};

// 重置筛选条件
const resetFilter = () => {
  searchForm.value = {
    keyword: "",
    category_id: "",
    public: "",
  };
  fetchTemplates();
};
const handleDetail = async (row: Template) => {
  try {
    loading.value = true;
    const res = await getTemplateById(row.id);
    templateDetail.value = res.data;
    detailDialogVisible.value = true;
  } catch (error) {
    ElMessage.error("获取模板详情失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};
// 创建模板
const handleCreate = () => {
  dialogType.value = "create";
  templateForm.value = {
    id: 0,
    name: "",
    category_id: undefined,
    description: "",
    tags: "",
    public: true,
    version: "",
    content: "",
  };
  dialogVisible.value = true;
};

// 编辑模板
const handleEdit = (row: any) => {
  dialogType.value = "edit";
  templateForm.value = {
    id: row.id,
    name: row.name,
    category_id: row.category_id,
    description: row.description,
    tags: row.tags || "",
    public: row.is_public,
    version: row.versions[0].version || "",
    content: row.versions[0].content || "",
  };
  dialogVisible.value = true;
};

// 删除模板
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm("确定要删除该模板吗？", "确认删除", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });
    await deleteTemplate(id);
    ElMessage.success("删除成功");
    fetchTemplates();
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("删除失败");
      console.error(error);
    }
  }
};

const handleVersion = (templateId: number) => {
  router.push(`/template-market/version?templateId=${templateId}`);
};

// 提交表单
const submitForm = async () => {
  if (!templateFormRef.value) return;

  try {
    await templateFormRef.value.validate();
    const formData = {
      ...templateForm.value,
      is_public: templateForm.value.public,
      tags: templateForm.value.tags
        ? templateForm.value.tags
            .split(",")
            .map((tag) => tag.trim())
            .join(",")
        : "",
    };

    if (dialogType.value === "create") {
      await createTemplate(formData);
      ElMessage.success("模板创建成功");
    } else {
      await updateTemplate(templateForm.value.id, formData);
      ElMessage.success("模板更新成功");
    }

    dialogVisible.value = false;
    fetchTemplates();
  } catch (error) {
    console.error("表单验证失败:", error);
  }
};

// 初始化
onMounted(() => {
  fetchCategories();
  fetchTemplates();
});
</script>

<style scoped>
.template-management-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.filter-card {
  margin-bottom: 20px;
  padding: 15px;
}

.pagination-container {
  margin-top: 15px;
  text-align: right;
}
</style>
