<template>
  <div class="search-container">
    <!-- 搜索区域 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="请输入模板名称或描述"
          ></el-input>
        </el-form-item>
        <el-form-item label="分类">
          <el-select
            v-model="searchForm.category_id"
            placeholder="全部分类"
            style="width: 150px"
          >
            <el-option label="全部分类" value=""></el-option>
            <el-option
              v-for="category in categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 结果展示区域 -->
    <div class="templates-container">
      <el-card
        v-for="template in templates"
        :key="template.id"
        class="template-card"
      >
        <div class="template-header">
          <h3>{{ template.name }}</h3>
          <el-tag :type="template.is_public ? 'success' : 'info'">
            {{ template.is_public ? "公开模板" : "私有模板" }}
          </el-tag>
        </div>
        <div class="template-body">
          <p class="template-description">{{ template.description }}</p>
          <div class="template-meta">
            <span>最新版本: {{ template.latest_version }}</span>
            <span>下载次数: {{ template.download_count }}</span>
            <span>更新时间: {{ formatDate(template.updated_at) }}</span>
          </div>
        </div>
        <div class="template-footer">
          <el-button type="primary" @click="handleDownload(template.id)">
            下载模板
          </el-button>
        </div>
      </el-card>

      <!-- 无结果状态 -->
      <div v-if="templates.length === 0 && !loading" class="no-result">
        <el-empty description="没有找到匹配的模板"></el-empty>
      </div>

      <!-- 加载状态 -->
      <div
        v-loading="loading"
        class="loading-state"
        element-loading-text="加载中..."
      >
        <div class="loading-content"></div>
      </div>
    </div>

    <!-- 分页控件 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      ></el-pagination>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { searchTemplates, getTemplateCategories } from "@/api/templateMarket";
import { downloadTemplate } from "@/api/templateMarket";
import { ElMessage } from "element-plus";

// 搜索表单数据
const searchForm = reactive({
  keyword: "",
  category_id: "",
  tags: "",
});

// 分类列表
const categories = ref([]);

// 模板列表
const templates = ref([]);

// 分页数据
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
});

// 加载状态
const loading = ref(false);

// 获取分类列表
const loadCategories = async () => {
  try {
    const response = await getTemplateCategories();
    categories.value = response.data.list;
  } catch (error) {
    ElMessage.error("获取分类失败");
    console.error(error);
  }
};

// 搜索模板
const handleSearch = async () => {
  loading.value = true;
  try {
    const response = await searchTemplates({
      keyword: searchForm.keyword,
      category_id: searchForm.category_id,
      tags: searchForm.tags,
      page: pagination.page,
      pageSize: pagination.pageSize,
    });
    templates.value = response.data;
    pagination.total = response.total;
  } catch (error) {
    ElMessage.error("搜索模板失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 重置表单
const resetForm = () => {
  searchForm.keyword = "";
  searchForm.category_id = "";
  searchForm.tags = "";
  pagination.page = 1;
  handleSearch();
};

// 下载模板
const handleDownload = async (templateId: number) => {
  try {
    loading.value = true;
    const response = await downloadTemplate(templateId);
    // 处理下载逻辑，通常是创建下载链接
    const blob = new Blob([JSON.stringify(response.data)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `template-${templateId}-${response.data.version}.json`;
    a.click();
    URL.revokeObjectURL(url);
    ElMessage.success("模板下载成功");
  } catch (error) {
    ElMessage.error("模板下载失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};

// 分页事件处理
const handleSizeChange = (val: number) => {
  pagination.pageSize = val;
  handleSearch();
};

const handleCurrentChange = (val: number) => {
  pagination.page = val;
  handleSearch();
};

// 日期格式化
const formatDate = (dateString: string) => {
  if (!dateString) return "";
  const date = new Date(dateString);
  return date.toLocaleString();
};

// 初始化
onMounted(() => {
  loadCategories();
  handleSearch();
});
</script>

<style scoped>
.search-container {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.templates-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.template-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.template-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.template-description {
  color: #606266;
  margin-bottom: 15px;
  flex-grow: 1;
}

.template-meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: #909399;
  font-size: 12px;
}

.template-footer {
  margin-top: 20px;
  text-align: right;
}

.pagination-container {
  text-align: right;
  margin-top: 20px;
}

.no-result,
.loading-state {
  margin: 50px 0;
  text-align: center;
}
</style>
