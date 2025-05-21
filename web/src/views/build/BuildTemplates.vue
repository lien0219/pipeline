<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>构建模板</h2>
        <p>管理和使用预定义的构建模板</p>
      </div>

      <div class="header-actions">
        <el-button type="primary" @click="createTemplate">
          <el-icon><Plus /></el-icon>
          创建模板
        </el-button>
      </div>
    </div>

    <el-card>
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="名称">
            <el-input v-model="filterForm.name" placeholder="模板名称" clearable />
          </el-form-item>

          <el-form-item label="类型">
            <el-select v-model="filterForm.type" placeholder="全部类型" clearable>
              <el-option label="Docker" value="docker" />
              <el-option label="Docker" value="docker" />
              <el-option label="Node.js" value="nodejs" />
              <el-option label="Go" value="golang" />
              <el-option label="Python" value="python" />
              <el-option label="Java" value="java" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleFilter">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetFilter">
              <el-icon><RefreshRight /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
          :data="templates"
          style="width: 100%"
          v-loading="loading"
          @sort-change="handleSortChange"
      >
        <el-table-column prop="name" label="名称" min-width="150" sortable="custom" />

        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />

        <el-table-column prop="is_public" label="公开" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_public ? 'success' : 'info'" size="small">
              {{ row.is_public ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="usage_count" label="使用次数" width="100" sortable="custom" />

        <el-table-column prop="creator.name" label="创建者" width="120">
          <template #default="{ row }">
            {{ row.creator?.name || row.creator?.username || '-' }}
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="180" sortable="custom">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
                link
                type="primary"
                size="small"
                @click="viewTemplate(row.id)"
            >
              查看
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="useTemplate(row)"
            >
              使用
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="editTemplate(row.id)"
                v-if="row.creator?.id === currentUserId"
            >
              编辑
            </el-button>

            <el-button
                link
                type="danger"
                size="small"
                @click="deleteTemplate(row.id)"
                v-if="row.creator?.id === currentUserId"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
            v-model:current-page="pagination.currentPage"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 30, 50]"
            :total="pagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 模板详情对话框 -->
    <el-dialog
        v-model="templateDialogVisible"
        :title="templateDialogTitle"
        width="70%"
    >
      <div v-loading="templateLoading">
        <el-form :model="templateForm" label-position="top" :disabled="templateViewOnly">
          <el-form-item label="名称">
            <el-input v-model="templateForm.name" placeholder="请输入模板名称" />
          </el-form-item>

          <el-form-item label="类型">
            <el-select v-model="templateForm.type" placeholder="请选择模板类型">
              <el-option label="Docker" value="docker" />
              <el-option label="Node.js" value="nodejs" />
              <el-option label="Go" value="golang" />
              <el-option label="Python" value="python" />
              <el-option label="Java" value="java" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>

          <el-form-item label="描述">
            <el-input
                v-model="templateForm.description"
                type="textarea"
                :rows="2"
                placeholder="请输入模板描述"
            />
          </el-form-item>

          <el-form-item label="内容">
            <el-input
                v-model="templateForm.content"
                type="textarea"
                :rows="10"
                placeholder="请输入模板内容"
                font-family="monospace"
            />
          </el-form-item>

          <el-form-item label="公开">
            <el-switch v-model="templateForm.is_public" />
          </el-form-item>
        </el-form>

        <div class="dialog-footer">
          <el-button @click="templateDialogVisible = false">关闭</el-button>
          <el-button
              type="primary"
              @click="saveTemplate"
              v-if="!templateViewOnly"
              :loading="templateSaving"
          >
            保存
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 使用模板对话框 -->
    <el-dialog
        v-model="useTemplateDialogVisible"
        title="使用模板"
        width="70%"
    >
      <div>
        <p>您正在使用模板 <strong>{{ selectedTemplate?.name }}</strong></p>
        <p>此模板将用于创建新的流水线。请选择要应用此模板的流水线：</p>

        <el-form :model="useTemplateForm" label-position="top">
          <el-form-item label="选择操作">
            <el-radio-group v-model="useTemplateForm.action">
              <el-radio label="create">创建新流水线</el-radio>
              <el-radio label="update">更新现有流水线</el-radio>
            </el-radio-group>
          </el-form-item>

          <template v-if="useTemplateForm.action === 'create'">
            <el-form-item label="流水线名称">
              <el-input v-model="useTemplateForm.name" placeholder="请输入流水线名称" />
            </el-form-item>

            <el-form-item label="Git 仓库">
              <el-input v-model="useTemplateForm.git_repo" placeholder="请输入Git仓库地址" />
            </el-form-item>

            <el-form-item label="Git 分支">
              <el-input v-model="useTemplateForm.git_branch" placeholder="请输入Git分支，默认为main" />
            </el-form-item>
          </template>

          <template v-else>
            <el-form-item label="选择流水线">
              <el-select v-model="useTemplateForm.pipeline_id" placeholder="请选择流水线">
                <el-option
                    v-for="pipeline in pipelines"
                    :key="pipeline.id"
                    :label="pipeline.name"
                    :value="pipeline.id"
                />
              </el-select>
            </el-form-item>
          </template>
        </el-form>

        <div class="dialog-footer">
          <el-button @click="useTemplateDialogVisible = false">取消</el-button>
          <el-button
              type="primary"
              @click="applyTemplate"
              :loading="applyingTemplate"
          >
            应用模板
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { usePipelineStore } from '@/stores/pipeline';
import { useAuthStore } from '@/stores/auth';
import { Plus, Search, RefreshRight } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

const router = useRouter();
const pipelineStore = usePipelineStore();
const authStore = useAuthStore();
const loading = ref(false);
const templateLoading = ref(false);
const templateSaving = ref(false);
const applyingTemplate = ref(false);

const templates = ref([]);
const pipelines = ref([]);
const currentUserId = computed(() => authStore.currentUser?.id);

// 模板对话框
const templateDialogVisible = ref(false);
const templateDialogTitle = ref('');
const templateViewOnly = ref(false);
const templateForm = reactive({
  id: null,
  name: '',
  type: 'docker',
  description: '',
  content: '',
  is_public: true
});

// 使用模板对话框
const useTemplateDialogVisible = ref(false);
const selectedTemplate = ref(null);
const useTemplateForm = reactive({
  action: 'create',
  name: '',
  git_repo: '',
  git_branch: 'main',
  pipeline_id: null
});

// 筛选表单
const filterForm = reactive({
  name: '',
  type: ''
});

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
});

// 排序
const sortParams = reactive({
  prop: 'created_at',
  order: 'descending'
});

// 获取模板列表
const fetchTemplates = async () => {
  loading.value = true;

  try {
    const params = {
      page: pagination.currentPage,
      limit: pagination.pageSize,
      sort_by: sortParams.prop,
      sort_order: sortParams.order === 'ascending' ? 'asc' : 'desc',
      name: filterForm.name || undefined,
      type: filterForm.type || undefined
    };

    const response = await pipelineStore.getBuildTemplates(params);
    templates.value = response.data || [];
    pagination.total = response.total || 0;
  } catch (error) {
    console.error('Failed to fetch templates:', error);
    ElMessage.error('获取模板列表失败');
  } finally {
    loading.value = false;
  }
};

// 获取流水线列表（用于使用模板）
const fetchPipelines = async () => {
  try {
    const response = await pipelineStore.fetchPipelines({ page:1,pageSize:10 });
    pipelines.value = response.data || [];
  } catch (error) {
    console.error('Failed to fetch pipelines:', error);
  }
};

// 筛选
const handleFilter = () => {
  pagination.currentPage = 1;
  fetchTemplates();
};

// 重置筛选
const resetFilter = () => {
  filterForm.name = '';
  filterForm.type = '';
  pagination.currentPage = 1;
  fetchTemplates();
};

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  fetchTemplates();
};

const handleCurrentChange = (page) => {
  pagination.currentPage = page;
  fetchTemplates();
};

// 排序处理
const handleSortChange = ({ prop, order }) => {
  if (prop) {
    sortParams.prop = prop;
    sortParams.order = order;
  } else {
    sortParams.prop = 'created_at';
    sortParams.order = 'descending';
  }
  fetchTemplates();
};

// 创建模板
const createTemplate = () => {
  templateDialogTitle.value = '创建模板';
  templateViewOnly.value = false;

  // 重置表单
  Object.assign(templateForm, {
    id: null,
    name: '',
    type: 'docker',
    description: '',
    content: '',
    is_public: true
  });

  templateDialogVisible.value = true;
};

// 查看模板
const viewTemplate = async (id) => {
  templateLoading.value = true;
  templateDialogTitle.value = '查看模板';
  templateViewOnly.value = true;
  templateDialogVisible.value = true;

  try {
    const response = await pipelineStore.getBuildTemplateById(id);
    Object.assign(templateForm, response.data);
  } catch (error) {
    console.error('Failed to fetch template:', error);
    ElMessage.error('获取模板详情失败');
    templateDialogVisible.value = false;
  } finally {
    templateLoading.value = false;
  }
};

// 编辑模板
const editTemplate = async (id) => {
  templateLoading.value = true;
  templateDialogTitle.value = '编辑模板';
  templateViewOnly.value = false;
  templateDialogVisible.value = true;

  try {
    const response = await pipelineStore.getBuildTemplateById(id);
    Object.assign(templateForm, response.data);
  } catch (error) {
    console.error('Failed to fetch template:', error);
    ElMessage.error('获取模板详情失败');
    templateDialogVisible.value = false;
  } finally {
    templateLoading.value = false;
  }
};

// 保存模板
const saveTemplate = async () => {
  // 表单验证
  if (!templateForm.name) {
    ElMessage.warning('请输入模板名称');
    return;
  }

  if (!templateForm.content) {
    ElMessage.warning('请输入模板内容');
    return;
  }

  templateSaving.value = true;

  try {
    if (templateForm.id) {
      // 更新模板
      await pipelineStore.updateBuildTemplate(templateForm.id, templateForm);
      ElMessage.success('模板更新成功');
    } else {
      // 创建模板
      await pipelineStore.createBuildTemplate(templateForm);
      ElMessage.success('模板创建成功');
    }

    templateDialogVisible.value = false;
    fetchTemplates();
  } catch (error) {
    console.error('Failed to save template:', error);
    ElMessage.error('保存模板失败');
  } finally {
    templateSaving.value = false;
  }
};

// 删除模板
const deleteTemplate = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除此模板吗？此操作不可恢复。', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    await pipelineStore.deleteBuildTemplate(id);
    ElMessage.success('模板已删除');
    fetchTemplates();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete template:', error);
      ElMessage.error('删除模板失败');
    }
  }
};

// 使用模板
const useTemplate = (template) => {
  selectedTemplate.value = template;

  // 重置表单
  Object.assign(useTemplateForm, {
    action: 'create',
    name: '',
    git_repo: '',
    git_branch: 'main',
    pipeline_id: null
  });

  useTemplateDialogVisible.value = true;
};

// 应用模板
const applyTemplate = async () => {
  // 表单验证
  if (useTemplateForm.action === 'create') {
    if (!useTemplateForm.name) {
      ElMessage.warning('请输入流水线名称');
      return;
    }

    if (!useTemplateForm.git_repo) {
      ElMessage.warning('请输入Git仓库地址');
      return;
    }
  } else {
    if (!useTemplateForm.pipeline_id) {
      ElMessage.warning('请选择流水线');
      return;
    }
  }

  applyingTemplate.value = true;

  try {
    if (useTemplateForm.action === 'create') {
      // 创建新流水线
      const pipelineData = {
        name: useTemplateForm.name,
        git_repo: useTemplateForm.git_repo,
        git_branch: useTemplateForm.git_branch,
        template_id: selectedTemplate.value.id
      };

      const response = await pipelineStore.createPipelineFromTemplate(pipelineData);
      ElMessage.success('流水线创建成功');
      useTemplateDialogVisible.value = false;

      // 跳转到新创建的流水线
      router.push(`/pipelines/${response.data.id}`);
    } else {
      // 更新现有流水线
      await pipelineStore.applyTemplateToExistingPipeline(useTemplateForm.pipeline_id, selectedTemplate.value.id);
      ElMessage.success('模板应用成功');
      useTemplateDialogVisible.value = false;
    }
  } catch (error) {
    console.error('Failed to apply template:', error);
    ElMessage.error('应用模板失败');
  } finally {
    applyingTemplate.value = false;
  }
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

onMounted(() => {
  fetchTemplates();
  fetchPipelines();
});
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-title h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.header-title p {
  margin: 0;
  color: #606266;
}

.filter-container {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  gap: 10px;
}
</style>
