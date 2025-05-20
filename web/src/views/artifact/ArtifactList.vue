<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>制品管理</h2>
      </div>

      <div class="header-actions">
        <el-button type="primary" @click="uploadArtifact">
          <el-icon><Upload /></el-icon>
          上传制品
        </el-button>
      </div>
    </div>

    <el-card>
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="名称">
            <el-input v-model="filterForm.name" placeholder="制品名称" clearable />
          </el-form-item>

          <el-form-item label="类型">
            <el-select v-model="filterForm.type" placeholder="全部类型" clearable>
              <el-option label="ZIP" value="zip" />
              <el-option label="TAR" value="tar" />
              <el-option label="JAR" value="jar" />
              <el-option label="WAR" value="war" />
              <el-option label="Docker镜像" value="docker" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>

          <el-form-item label="流水线">
            <el-select v-model="filterForm.pipeline_id" placeholder="全部流水线" clearable>
              <el-option
                  v-for="pipeline in pipelines"
                  :key="pipeline.id"
                  :label="pipeline.name"
                  :value="pipeline.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="时间范围">
            <el-date-picker
                v-model="filterForm.dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
            />
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
          :data="artifacts"
          style="width: 100%"
          v-loading="loading"
          @sort-change="handleSortChange"
      >
        <el-table-column prop="name" label="名称" min-width="200" sortable="custom" show-overflow-tooltip />

        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>

        <el-table-column prop="pipeline.name" label="流水线" width="150" show-overflow-tooltip>
          <template #default="{ row }">
            <router-link
                v-if="row.pipeline_id"
                :to="`/pipelines/${row.pipeline_id}`"
                class="pipeline-link"
            >
              {{ row.pipeline?.name }}
            </router-link>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="pipeline_run_id" label="构建ID" width="100">
          <template #default="{ row }">
            <router-link
                v-if="row.pipeline_id && row.pipeline_run_id"
                :to="`/pipelines/${row.pipeline_id}/runs/${row.pipeline_run_id}`"
                class="pipeline-link"
            >
              {{ row.pipeline_run_id }}
            </router-link>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="download_count" label="下载次数" width="100" sortable="custom" />

        <el-table-column prop="created_by" label="创建者" width="120" />

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
                @click="downloadArtifact(row.id)"
            >
              下载
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="viewArtifactDetail(row.id)"
            >
              详情
            </el-button>

            <el-button
                link
                type="danger"
                size="small"
                @click="deleteArtifact(row.id)"
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

    <el-dialog
        v-model="uploadDialogVisible"
        title="上传制品"
        width="50%"
    >
      <el-form
          ref="uploadFormRef"
          :model="uploadFormData"
          :rules="uploadRules"
          label-position="top"
      >
        <el-form-item label="制品名称" prop="name">
          <el-input v-model="uploadFormData.name" placeholder="请输入制品名称" />
        </el-form-item>

        <el-form-item label="制品类型" prop="type">
          <el-select v-model="uploadFormData.type" placeholder="请选择制品类型">
            <el-option label="ZIP" value="zip" />
            <el-option label="TAR" value="tar" />
            <el-option label="JAR" value="jar" />
            <el-option label="WAR" value="war" />
            <el-option label="Docker镜像" value="docker" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="关联流水线" prop="pipeline_id">
          <el-select v-model="uploadFormData.pipeline_id" placeholder="请选择关联流水线" clearable>
            <el-option
                v-for="pipeline in pipelines"
                :key="pipeline.id"
                :label="pipeline.name"
                :value="pipeline.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="版本" prop="version">
          <el-input v-model="uploadFormData.version" placeholder="请输入版本号，例如：v1.0.0" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
              v-model="uploadFormData.description"
              type="textarea"
              :rows="3"
              placeholder="请输入制品描述"
          />
        </el-form-item>

        <el-form-item label="文件" prop="file" v-if="uploadFormData.type !== 'docker'">
          <el-upload
              class="artifact-upload"
              drag
              action="#"
              :auto-upload="false"
              :on-change="handleFileChange"
              :limit="1"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">拖拽文件到此处，或 <em>点击上传</em></div>
            <template #tip>
              <div class="el-upload__tip">
                请上传制品文件，大小不超过500MB
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item label="Docker镜像" prop="docker_image" v-if="uploadFormData.type === 'docker'">
          <el-input v-model="uploadFormData.docker_image" placeholder="请输入Docker镜像名称，例如：registry.example.com/app:v1.0.0" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitUpload" :loading="uploading">
            上传
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
        v-model="detailDialogVisible"
        title="制品详情"
        width="60%"
    >
      <div v-loading="detailLoading">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{ currentArtifact.id }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{ currentArtifact.name }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            <el-tag size="small">{{ currentArtifact.type?.toUpperCase() }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="大小">{{ formatFileSize(currentArtifact.size) }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{ currentArtifact.version }}</el-descriptions-item>
          <el-descriptions-item label="下载次数">{{ currentArtifact.download_count }}</el-descriptions-item>
          <el-descriptions-item label="流水线">
            <router-link
                v-if="currentArtifact.pipeline_id"
                :to="`/pipelines/${currentArtifact.pipeline_id}`"
                class="pipeline-link"
            >
              {{ currentArtifact.pipeline?.name }}
            </router-link>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="构建ID">
            <router-link
                v-if="currentArtifact.pipeline_id && currentArtifact.pipeline_run_id"
                :to="`/pipelines/${currentArtifact.pipeline_id}/runs/${currentArtifact.pipeline_run_id}`"
                class="pipeline-link"
            >
              {{ currentArtifact.pipeline_run_id }}
            </router-link>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="创建者">{{ currentArtifact.created_by }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(currentArtifact.created_at) }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{ currentArtifact.description || '无描述' }}</el-descriptions-item>
        </el-descriptions>

        <div class="artifact-actions">
          <el-button type="primary" @click="downloadArtifact(currentArtifact.id)">
            <el-icon><Download /></el-icon>
            下载制品
          </el-button>

          <el-button type="success" @click="deployArtifact(currentArtifact.id)">
            <el-icon><Position /></el-icon>
            部署制品
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Upload, Search, RefreshRight, UploadFilled, Download, Position } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

// 模拟数据，实际项目中应该从API获取
const artifacts = ref([
  {
    id: 1,
    name: 'app-v1.0.0.zip',
    type: 'zip',
    size: 15728640, // 15MB
    version: 'v1.0.0',
    pipeline_id: 1,
    pipeline: { name: '主应用构建' },
    pipeline_run_id: 101,
    download_count: 5,
    created_by: 'admin',
    created_at: '2023-05-10T09:00:00',
    description: '初始版本发布包'
  },
  {
    id: 2,
    name: 'app-v1.1.0.zip',
    type: 'zip',
    size: 16777216, // 16MB
    version: 'v1.1.0',
    pipeline_id: 1,
    pipeline: { name: '主应用构建' },
    pipeline_run_id: 102,
    download_count: 8,
    created_by: 'admin',
    created_at: '2023-05-15T10:30:00',
    description: '新增功能和bug修复'
  },
  {
    id: 3,
    name: 'app-v1.1.1.zip',
    type: 'zip',
    size: 16252928, // 15.5MB
    version: 'v1.1.1',
    pipeline_id: 1,
    pipeline: { name: '主应用构建' },
    pipeline_run_id: 103,
    download_count: 3,
    created_by: 'admin',
    created_at: '2023-05-16T14:20:00',
    description: '修复v1.1.0中的关键bug'
  },
  {
    id: 4,
    name: 'app-v1.2.0.zip',
    type: 'zip',
    size: 17825792, // 17MB
    version: 'v1.2.0',
    pipeline_id: 1,
    pipeline: { name: '主应用构建' },
    pipeline_run_id: 104,
    download_count: 2,
    created_by: 'admin',
    created_at: '2023-05-18T11:45:00',
    description: '新版本测试'
  },
  {
    id: 5,
    name: 'api-service-v1.0.0.jar',
    type: 'jar',
    size: 10485760, // 10MB
    version: 'v1.0.0',
    pipeline_id: 2,
    pipeline: { name: 'API服务构建' },
    pipeline_run_id: 201,
    download_count: 4,
    created_by: 'admin',
    created_at: '2023-05-12T13:20:00',
    description: 'API服务初始版本'
  }
]);

// 模拟流水线数据
const pipelines = ref([
  { id: 1, name: '主应用构建' },
  { id: 2, name: 'API服务构建' },
  { id: 3, name: '前端构建' }
]);

const loading = ref(false);
const uploadDialogVisible = ref(false);
const detailDialogVisible = ref(false);
const detailLoading = ref(false);
const uploading = ref(false);
const uploadFormRef = ref(null);
const currentArtifact = ref({});

// 筛选表单
const filterForm = reactive({
  name: '',
  type: '',
  pipeline_id: '',
  dateRange: []
});

// 分页
const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: artifacts.value.length
});

// 排序
const sortParams = reactive({
  prop: 'created_at',
  order: 'descending'
});

// 上传表单
const uploadFormData = reactive({
  name: '',
  type: 'zip',
  pipeline_id: null,
  version: '',
  description: '',
  file: null,
  docker_image: ''
});

// 上传表单验证规则
const uploadRules = {
  name: [
    { required: true, message: '请输入制品名称', trigger: 'blur' },
    { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择制品类型', trigger: 'change' }
  ],
  version: [
    { required: true, message: '请输入版本号', trigger: 'blur' }
  ],
  file: [
    { required: true, message: '请上传文件', trigger: 'change' }
  ],
  docker_image: [
    { required: true, message: '请输入Docker镜像名称', trigger: 'blur' }
  ]
};

// 上传制品
const uploadArtifact = () => {
  // 重置表单
  Object.assign(uploadFormData, {
    name: '',
    type: 'zip',
    pipeline_id: null,
    version: '',
    description: '',
    file: null,
    docker_image: ''
  });

  uploadDialogVisible.value = true;
};

// 处理文件变更
const handleFileChange = (file) => {
  uploadFormData.file = file.raw;

  // 如果没有填写名称，使用文件名
  if (!uploadFormData.name) {
    uploadFormData.name = file.name;
  }
};

// 提交上传
const submitUpload = async () => {
  if (!uploadFormRef.value) return;

  // 根据制品类型动态设置验证规则
  const rules = { ...uploadRules };
  if (uploadFormData.type === 'docker') {
    delete rules.file;
  } else {
    delete rules.docker_image;
  }

  await uploadFormRef.value.validate(async (valid) => {
    if (valid) {
      uploading.value = true;

      try {
        // 实际项目中应该调用API上传制品
        const newId = Math.max(...artifacts.value.map(a => a.id)) + 1;

        const newArtifact = {
          id: newId,
          name: uploadFormData.name,
          type: uploadFormData.type,
          size: uploadFormData.file ? uploadFormData.file.size : 1024 * 1024, // 模拟大小
          version: uploadFormData.version,
          pipeline_id: uploadFormData.pipeline_id,
          pipeline: uploadFormData.pipeline_id ? pipelines.value.find(p => p.id === uploadFormData.pipeline_id) : null,
          pipeline_run_id: null,
          download_count: 0,
          created_by: 'admin',
          created_at: new Date().toISOString(),
          description: uploadFormData.description
        };

        artifacts.value.push(newArtifact);
        ElMessage.success('制品上传成功');
        uploadDialogVisible.value = false;
      } catch (error) {
        console.error('上传制品失败:', error);
        ElMessage.error('上传制品失败');
      } finally {
        uploading.value = false;
      }
    } else {
      ElMessage.warning('请填写必填项');
      return false;
    }
  });
};

// 下载制品
const downloadArtifact = (id) => {
  // 实际项目中应该调用API下载制品
  const artifact = artifacts.value.find(a => a.id === id);
  if (artifact) {
    // 模拟下载
    ElMessage.success(`开始下载制品: ${artifact.name}`);

    // 更新下载次数
    artifact.download_count += 1;
  }
};

// 查看制品详情
const viewArtifactDetail = (id) => {
  detailLoading.value = true;
  detailDialogVisible.value = true;

  // 实际项目中应该调用API获取制品详情
  setTimeout(() => {
    const artifact = artifacts.value.find(a => a.id === id);
    if (artifact) {
      currentArtifact.value = { ...artifact };
    }
    detailLoading.value = false;
  }, 500);
};

// 部署制品
const deployArtifact = (id) => {
  // 实际项目中应该跳转到部署页面或打开部署对话框
  const artifact = artifacts.value.find(a => a.id === id);
  if (artifact) {
    ElMessage.info(`准备部署制品: ${artifact.name}`);
    detailDialogVisible.value = false;
  }
};

// 删除制品
const deleteArtifact = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除此制品吗？此操作不可恢复。', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    // 实际项目中应该调用API删除制品
    artifacts.value = artifacts.value.filter(a => a.id !== id);
    ElMessage.success('制品已删除');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除制品失败:', error);
    }
  }
};

// 筛选
const handleFilter = () => {
  // 实际项目中应该调用API获取筛选后的数据
  console.log('筛选条件:', filterForm);
};

// 重置筛选
const resetFilter = () => {
  filterForm.name = '';
  filterForm.type = '';
  filterForm.pipeline_id = '';
  filterForm.dateRange = [];
  // 实际项目中应该重新获取数据
};

// 分页处理
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  // 实际项目中应该重新获取数据
};

const handleCurrentChange = (page) => {
  pagination.currentPage = page;
  // 实际项目中应该重新获取数据
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
  // 实际项目中应该重新获取数据
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B';

  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

onMounted(() => {
  // 实际项目中应该从API获取制品列表和流水线列表
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

.artifact-upload {
  width: 100%;
}

.artifact-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: center;
}
</style>
