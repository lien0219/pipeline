<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>流水线管理</h2>
        <p>管理和监控所有CI/CD流水线</p>
      </div>

      <div class="header-actions">
        <el-button type="primary" @click="$router.push('/pipelines/create')">
          <el-icon><Plus /></el-icon>
          创建流水线
        </el-button>
      </div>
    </div>

    <el-card>
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="名称">
            <el-input v-model="filterForm.name" placeholder="流水线名称" clearable />
          </el-form-item>

          <el-form-item label="状态">
            <el-select v-model="filterForm.status" placeholder="全部状态" clearable>
              <el-option label="成功" value="success" />
              <el-option label="运行中" value="running" />
              <el-option label="失败" value="failed" />
              <el-option label="等待中" value="pending" />
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
          :data="pipelines"
          style="width: 100%"
          v-loading="loading"
          @sort-change="handleSortChange"
      >
        <el-table-column prop="name" label="名称" sortable="custom" min-width="180">
          <template #default="{ row }">
            <router-link :to="`/pipelines/${row.id}`" class="pipeline-link">
              {{ row.name }}
            </router-link>
          </template>
        </el-table-column>

        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="branch" label="分支" width="120" />

        <el-table-column prop="last_run" label="最近执行" width="180">
          <template #default="{ row }">
            {{ row.last_run ? formatDate(row.last_run) : '从未执行' }}
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
                @click="triggerPipeline(row.id)"
                :disabled="row.status === 'running'"
            >
              运行
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="$router.push(`/pipelines/${row.id}`)"
            >
              详情
            </el-button>

            <el-dropdown>
              <el-button link type="primary" size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="$router.push(`/pipelines/${row.id}/edit`)">
                    编辑
                  </el-dropdown-item>
                  <el-dropdown-item @click="clonePipeline(row)">
                    克隆
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="confirmDelete(row)">
                    <span style="color: var(--el-color-danger)">删除</span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { usePipelineStore } from '@/stores/pipeline';
import { Plus, Search, RefreshRight, ArrowDown } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

const pipelineStore = usePipelineStore();
const loading = ref(false);
const pipelines = ref([]);

// 筛选表单
const filterForm = reactive({
  name: '',
  status: ''
});

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
});

// 排序
const sortParams = reactive({
  prop: 'created_at',
  order: 'descending'
});

// 获取流水线列表
const fetchPipelines = async () => {
  loading.value = true;

  try {
    const params = {
      page: pagination.currentPage,
      limit: pagination.pageSize,
      sort_by: sortParams.prop,
      sort_order: sortParams.order === 'ascending' ? 'asc' : 'desc',
      name: filterForm.name || undefined,
      status: filterForm.status || undefined
    };

    const response = await pipelineStore.fetchPipelines(params);
    pipelines.value = response.data || [];
    pagination.total = response.total || 0;
  } catch (error) {
    console.error('Failed to fetch pipelines:', error);
  } finally {
    loading.value = false;
  }
};

// 筛选
const handleFilter = () => {
  pagination.currentPage = 1;
  fetchPipelines();
};

// 重置筛选
const resetFilter = () => {
  filterForm.name = '';
  filterForm.status = '';
  pagination.currentPage = 1;
  fetchPipelines();
};

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  fetchPipelines();
};

const handleCurrentChange = (page) => {
  pagination.currentPage = page;
  fetchPipelines();
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
  fetchPipelines();
};

// 触发流水线
const triggerPipeline = async (id) => {
  try {
    await pipelineStore.triggerPipeline(id);
    ElMessage.success('流水线已触发');
    fetchPipelines();
  } catch (error) {
    console.error('Failed to trigger pipeline:', error);
  }
};

// 克隆流水线
const clonePipeline = (pipeline) => {
  ElMessage.info('克隆功能开发中');
};

// 确认删除
const confirmDelete = (pipeline) => {
  ElMessageBox.confirm(
      `确定要删除流水线 "${pipeline.name}" 吗？此操作不可恢复。`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
  )
      .then(async () => {
        try {
          await pipelineStore.deletePipeline(pipeline.id);
          ElMessage.success('流水线已删除');
          fetchPipelines();
        } catch (error) {
          console.error('Failed to delete pipeline:', error);
        }
      })
      .catch(() => {
        // 取消删除
      });
};

// 格式化状态
const getStatusType = (status) => {
  switch (status) {
    case 'success': return 'success';
    case 'running': return 'primary';
    case 'failed': return 'danger';
    case 'pending': return 'info';
    default: return 'info';
  }
};

const getStatusText = (status) => {
  switch (status) {
    case 'success': return '成功';
    case 'running': return '运行中';
    case 'failed': return '失败';
    case 'pending': return '等待中';
    default: return '未知';
  }
};

// 格式化日期
const formatDate = (date) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

onMounted(() => {
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

.pipeline-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.pipeline-link:hover {
  text-decoration: underline;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
