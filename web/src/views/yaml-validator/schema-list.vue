<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>Schema管理</span>
        <el-button type="primary" @click="openSchemaForm()" :icon="Plus">
          新建Schema
        </el-button>
      </div>
    </template>

    <el-form :model="searchForm" class="search-form" label-width="100px">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-form-item label="Schema类型">
            <el-select v-model="searchForm.type" placeholder="请选择类型">
              <el-option label="全部" value="" />
              <el-option label="Kubernetes" value="kubernetes" />
              <el-option label="自定义" value="custom" />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="名称">
            <el-input v-model="searchForm.name" placeholder="请输入名称" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item>
            <el-button type="primary" @click="fetchSchemas()" :icon="Search">
              搜索
            </el-button>
            <el-button @click="resetSearch()">重置</el-button>
          </el-form-item>
        </el-col>
      </el-row>
    </el-form>

    <el-table v-loading="loading" :data="schemaList" style="width: 100%">
      <el-table-column prop="name" label="名称" width="200" />
      <el-table-column prop="type" label="类型" width="100" />
      <el-table-column prop="version" label="版本" width="100" />
      <el-table-column prop="description" label="描述" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="scope">
          <el-button
            type="primary"
            size="small"
            @click="viewSchema(scope.row.id)"
            :icon="View"
          >
            查看
          </el-button>
          <el-button
            type="success"
            size="small"
            @click="openSchemaForm(scope.row.id)"
            :icon="Edit"
          >
            编辑
          </el-button>
          <el-button
            type="danger"
            size="small"
            @click="deleteSchema(scope.row.id)"
            :icon="Delete"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="pagination.currentPage"
      v-model:page-size="pagination.pageSize"
      :page-sizes="[10, 20, 50, 100]"
      :default-page-size="10"
      layout="total, sizes, prev, pager, next, jumper"
      :total="pagination.total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      style="margin-top: 20px"
    />

    <!-- Schema表单对话框 -->
    <SchemaForm
      v-model:visible="formVisible"
      :schema-id="schemaId"
      @success="fetchSchemas"
    />
  </el-card>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Search, Plus, Edit, Delete, View } from "@element-plus/icons-vue";
import {
  getYAMLSchemas,
  deleteYAMLSchema,
  YAMLSchema,
  PageResponse,
} from "../../api/yamlValidator";
import SchemaForm from "./schema-form.vue";

const schemaList = ref<YAMLSchema[]>([]);
const loading = ref(false);
const formVisible = ref(false);
const schemaId = ref<number | null>(null);

const searchForm = reactive<{
  type: string;
  name: string;
}>({
  type: "",
  name: "",
});

const pagination = reactive<{
  currentPage: number;
  pageSize: number;
  total: number;
}>({
  currentPage: 1,
  pageSize: 10,
  total: 0,
});

// 获取Schema列表
const fetchSchemas = async () => {
  try {
    loading.value = true;
    const res: any = await getYAMLSchemas({
      type: searchForm.type,
      page: pagination.currentPage,
      pageSize: pagination.pageSize || 10,
    });
    schemaList.value = res?.data?.list;
    pagination.total = res?.data?.total;
    pagination.currentPage = res?.data?.page;
    pagination.pageSize = res?.data?.pageSize;
  } catch (error) {
    ElMessage.error("获取Schema列表失败");
    console.error("Failed to fetch schemas:", error);
  } finally {
    loading.value = false;
  }
};

// 打开Schema表单
const openSchemaForm = (id?: number) => {
  schemaId.value = id || null;
  formVisible.value = true;
};

// 查看Schema详情
const viewSchema = (id: number) => {
  window.open(`/yaml-validator/schema/${id}`, "_blank");
};

// 删除Schema
const deleteSchema = async (id: number) => {
  try {
    await ElMessageBox.confirm("确定要删除这个Schema吗？", "警告", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await deleteYAMLSchema(id);
    ElMessage.success("删除成功");
    fetchSchemas();
  } catch (error) {
    // 忽略取消操作的错误
    if (error !== "cancel") {
      ElMessage.error("删除失败");
      console.error("Failed to delete schema:", error);
    }
  }
};

// 重置搜索
const resetSearch = () => {
  searchForm.type = "";
  searchForm.name = "";
  fetchSchemas();
};

// 分页处理
const handleSizeChange = (size: number) => {
  pagination.pageSize = size;
  fetchSchemas();
};

const handleCurrentChange = (current: number) => {
  pagination.currentPage = current;
  fetchSchemas();
};

// 初始化
onMounted(() => {
  fetchSchemas();
});
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}
</style>
