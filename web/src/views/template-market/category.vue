<template>
  <div class="category-management">
    <div class="page-header">
      <h1>分类管理</h1>
      <div class="search-container">
        <el-input
          v-model="searchKeyword"
          placeholder="请输入分类名称或描述"
          style="width: 300px"
          @keyup.enter.prevent="handleSearch"
        />
        <el-button
          type="primary"
          @click="handleSearch"
          style="margin-left: 10px"
          >搜索</el-button
        >
        <el-button @click="handleReset" style="margin-left: 10px"
          >重置</el-button
        >
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon-plus /> 添加分类
      </el-button>
    </div>

    <el-card class="content-card">
      <el-table :data="categories" border stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="分类名称" align="center" />
        <el-table-column prop="description" label="分类描述" align="center" />
        <el-table-column prop="created_at" label="创建时间" align="center" />
        <el-table-column label="操作" align="center" width="180">
          <template #default="scope">
            <el-button
              type="primary"
              size="small"
              @click="handleEdit(scope.row)"
              >编辑</el-button
            >
            <el-button
              type="danger"
              size="small"
              @click="handleDelete(scope.row.id)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加/编辑分类弹窗 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="分类名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="分类描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            rows="4"
            placeholder="请输入分类描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  createTemplateCategory,
  getTemplateCategories,
  updateTemplateCategory,
  deleteTemplateCategory,
} from "@/api/templateMarket";
import type { TemplateCategory } from "../../api/types";

// 状态定义
const categories = ref<TemplateCategory[]>([]);
const dialogVisible = ref(false);
const dialogTitle = ref("添加分类");
const formRef = ref();
const loading = ref(false);

const searchKeyword = ref("");
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 表单数据
const formData = reactive({
  id: 0,
  name: "",
  description: "",
});

// 表单验证规则
const formRules = reactive({
  name: [
    { required: true, message: "请输入分类名称", trigger: "blur" },
    {
      min: 2,
      max: 20,
      message: "分类名称长度应在2-20个字符之间",
      trigger: "blur",
    },
  ],
  description: [
    { max: 200, message: "分类描述不能超过200个字符", trigger: "blur" },
  ],
});

// 获取分类列表
const fetchCategories = async () => {
  try {
    loading.value = true;
    const response: any = await getTemplateCategories({
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchKeyword.value,
    });
    categories.value = response?.data?.list || [];
    total.value = response?.data?.total || 0;
  } catch (error) {
    ElMessage.error("获取分类列表失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};
const handleSearch = () => {
  currentPage.value = 1;
  fetchCategories();
};

const handleReset = () => {
  searchKeyword.value = "";
  currentPage.value = 1;
  fetchCategories();
};

const handleSizeChange = (val: number) => {
  pageSize.value = val;
  currentPage.value = 1;
  fetchCategories();
};

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  fetchCategories();
};
// 添加分类
const handleAdd = () => {
  dialogTitle.value = "添加分类";
  formData.id = 0;
  formData.name = "";
  formData.description = "";
  dialogVisible.value = true;
};

// 编辑分类
const handleEdit = (row: TemplateCategory) => {
  dialogTitle.value = "编辑分类";
  formData.id = row.id;
  formData.name = row.name;
  formData.description = row.description || "";
  dialogVisible.value = true;
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    loading.value = true;

    if (formData.id) {
      // 更新分类
      await updateTemplateCategory(formData.id, {
        name: formData.name,
        description: formData.description,
      });
      ElMessage.success("分类更新成功");
    } else {
      // 创建分类
      await createTemplateCategory({
        name: formData.name,
        description: formData.description,
      });
      ElMessage.success("分类创建成功");
    }

    dialogVisible.value = false;
    fetchCategories();
  } catch (error) {
    ElMessage.error(formData.id ? "分类更新失败" : "分类创建失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 删除分类
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm("确定要删除该分类吗？", "确认删除", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await deleteTemplateCategory(id);
    ElMessage.success("分类删除成功");
    fetchCategories();
  } catch (error) {
    // 用户取消或操作失败
    if (error !== "cancel") {
      ElMessage.error("分类删除失败");
      console.error(error);
    }
  }
};

// 页面加载时获取分类列表
onMounted(() => {
  fetchCategories();
});
</script>

<style scoped>
.category-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.content-card {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  border-radius: 8px;
  overflow: hidden;
}
.pagination-container {
  margin-top: 16px;
  text-align: right;
  padding: 10px;
}

.search-container {
  display: flex;
  align-items: center;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 10px;
}
</style>
