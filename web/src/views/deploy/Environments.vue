<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>部署环境</h2>
      </div>

      <div class="header-actions">
        <el-button type="primary" @click="createEnvironment">
          <el-icon><Plus /></el-icon>
          创建环境
        </el-button>
      </div>
    </div>

    <el-card>
      <el-table :data="environments" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="环境名称" min-width="150" />

        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getEnvironmentType(row.type)">{{ row.type }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="url" label="URL" min-width="200" show-overflow-tooltip />

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="last_deployed_at" label="最近部署" width="180">
          <template #default="{ row }">
            {{ formatDate(row.last_deployed_at) }}
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="创建时间" width="180">
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
                @click="viewEnvironment(row.id)"
            >
              详情
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="deployToEnvironment(row.id)"
            >
              部署
            </el-button>

            <el-button
                link
                type="primary"
                size="small"
                @click="editEnvironment(row.id)"
            >
              编辑
            </el-button>

            <el-button
                link
                type="danger"
                size="small"
                @click="deleteEnvironment(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
        v-model="dialogVisible"
        :title="dialogTitle"
        width="50%"
    >
      <el-form
          ref="environmentForm"
          :model="environmentForm"
          :rules="rules"
          label-position="top"
      >
        <el-form-item label="环境名称" prop="name">
          <el-input v-model="environmentForm.name" placeholder="请输入环境名称" />
        </el-form-item>

        <el-form-item label="环境类型" prop="type">
          <el-select v-model="environmentForm.type" placeholder="请选择环境类型">
            <el-option label="开发环境" value="development" />
            <el-option label="测试环境" value="testing" />
            <el-option label="预发布环境" value="staging" />
            <el-option label="生产环境" value="production" />
          </el-select>
        </el-form-item>

        <el-form-item label="环境URL" prop="url">
          <el-input v-model="environmentForm.url" placeholder="请输入环境URL" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
              v-model="environmentForm.description"
              type="textarea"
              :rows="3"
              placeholder="请输入环境描述"
          />
        </el-form-item>

        <h4>配置变量</h4>
        <p class="section-desc">为此环境设置配置变量</p>

        <div
            v-for="(variable, index) in environmentForm.variables"
            :key="index"
            class="variable-item"
        >
          <el-row :gutter="10">
            <el-col :span="10">
              <el-form-item
                  :prop="`variables.${index}.key`"
                  :rules="[{ required: true, message: '请输入变量名', trigger: 'blur' }]"
              >
                <el-input v-model="variable.key" placeholder="变量名" />
              </el-form-item>
            </el-col>

            <el-col :span="10">
              <el-form-item :prop="`variables.${index}.value`">
                <el-input v-model="variable.value" placeholder="变量值" />
              </el-form-item>
            </el-col>

            <el-col :span="4">
              <el-button
                  type="danger"
                  icon="Delete"
                  circle
                  @click="removeVariable(index)"
              />
            </el-col>
          </el-row>
        </div>

        <div class="add-variable">
          <el-button type="primary" plain @click="addVariable">
            <el-icon><Plus /></el-icon>
            添加变量
          </el-button>
        </div>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            {{ isEdit ? '更新' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Plus, Delete } from '@element-plus/icons-vue';
import dayjs from 'dayjs';

// 模拟数据，实际项目中应该从API获取
const environments = ref([
  {
    id: 1,
    name: '开发环境',
    type: 'development',
    url: 'https://dev.example.com',
    status: 'active',
    last_deployed_at: '2023-05-15T10:30:00',
    created_at: '2023-01-10T08:00:00',
    description: '用于开发和测试的环境',
    variables: [
      { key: 'API_URL', value: 'https://api.dev.example.com' },
      { key: 'DEBUG', value: 'true' }
    ]
  },
  {
    id: 2,
    name: '测试环境',
    type: 'testing',
    url: 'https://test.example.com',
    status: 'active',
    last_deployed_at: '2023-05-14T15:45:00',
    created_at: '2023-01-15T09:30:00',
    description: 'QA团队使用的测试环境',
    variables: [
      { key: 'API_URL', value: 'https://api.test.example.com' },
      { key: 'DEBUG', value: 'true' }
    ]
  },
  {
    id: 3,
    name: '预发布环境',
    type: 'staging',
    url: 'https://staging.example.com',
    status: 'active',
    last_deployed_at: '2023-05-13T11:20:00',
    created_at: '2023-02-01T10:00:00',
    description: '用于最终测试的预发布环境',
    variables: [
      { key: 'API_URL', value: 'https://api.staging.example.com' },
      { key: 'DEBUG', value: 'false' }
    ]
  },
  {
    id: 4,
    name: '生产环境',
    type: 'production',
    url: 'https://www.example.com',
    status: 'active',
    last_deployed_at: '2023-05-10T09:00:00',
    created_at: '2023-02-15T14:00:00',
    description: '面向用户的生产环境',
    variables: [
      { key: 'API_URL', value: 'https://api.example.com' },
      { key: 'DEBUG', value: 'false' }
    ]
  }
]);

const loading = ref(false);
const dialogVisible = ref(false);
const dialogTitle = ref('创建环境');
const isEdit = ref(false);
const submitting = ref(false);

// 表单数据
const environmentForm = reactive({
  id: null,
  name: '',
  type: 'development',
  url: '',
  description: '',
  variables: [{ key: '', value: '' }]
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入环境名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择环境类型', trigger: 'change' }
  ],
  url: [
    { required: true, message: '请输入环境URL', trigger: 'blur' }
  ]
};

// 创建环境
const createEnvironment = () => {
  dialogTitle.value = '创建环境';
  isEdit.value = false;

  // 重置表单
  Object.assign(environmentForm, {
    id: null,
    name: '',
    type: 'development',
    url: '',
    description: '',
    variables: [{ key: '', value: '' }]
  });

  dialogVisible.value = true;
};

// 查看环境详情
const viewEnvironment = (id) => {
  // 实际项目中应该跳转到环境详情页或打开详情对话框
  console.log('查看环境详情:', id);
};

// 部署到环境
const deployToEnvironment = (id) => {
  // 实际项目中应该打开部署对话框
  console.log('部署到环境:', id);
};

// 编辑环境
const editEnvironment = (id) => {
  dialogTitle.value = '编辑环境';
  isEdit.value = true;

  // 查找环境数据
  const env = environments.value.find(e => e.id === id);
  if (env) {
    // 复制数据到表单
    Object.assign(environmentForm, {
      id: env.id,
      name: env.name,
      type: env.type,
      url: env.url,
      description: env.description,
      variables: [...(env.variables || [{ key: '', value: '' }])]
    });

    dialogVisible.value = true;
  }
};

// 删除环境
const deleteEnvironment = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除此环境吗？此操作不可恢复。', '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    // 实际项目中应该调用API删除环境
    environments.value = environments.value.filter(e => e.id !== id);
    ElMessage.success('环境已删除');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除环境失败:', error);
    }
  }
};

// 添加变量
const addVariable = () => {
  environmentForm.variables.push({ key: '', value: '' });
};

// 移除变量
const removeVariable = (index) => {
  if (environmentForm.variables.length > 1) {
    environmentForm.variables.splice(index, 1);
  } else {
    ElMessage.warning('至少需要一个变量');
  }
};

// 提交表单
const submitForm = async () => {
  if (!environmentForm.value) return;

  await environmentForm.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true;

      try {
        // 实际项目中应该调用API创建或更新环境
        if (isEdit.value) {
          // 更新环境
          const index = environments.value.findIndex(e => e.id === environmentForm.id);
          if (index !== -1) {
            environments.value[index] = {
              ...environments.value[index],
              ...environmentForm,
              updated_at: new Date().toISOString()
            };
          }
          ElMessage.success('环境已更新');
        } else {
          // 创建环境
          const newId = Math.max(...environments.value.map(e => e.id), 0) + 1;
          environments.value.push({
            ...environmentForm,
            id: newId,
            status: 'active',
            created_at: new Date().toISOString(),
            last_deployed_at: null
          });
          ElMessage.success('环境已创建');
        }

        dialogVisible.value = false;
      } catch (error) {
        console.error('提交表单失败:', error);
        ElMessage.error('操作失败，请重试');
      } finally {
        submitting.value = false;
      }
    } else {
      ElMessage.warning('请填写必填项');
      return false;
    }
  });
};

// 获取环境类型样式
const getEnvironmentType = (type) => {
  switch (type) {
    case 'development': return 'info';
    case 'testing': return 'warning';
    case 'staging': return 'success';
    case 'production': return 'danger';
    default: return 'info';
  }
};

// 获取状态样式
const getStatusType = (status) => {
  switch (status) {
    case 'active': return 'success';
    case 'inactive': return 'info';
    case 'error': return 'danger';
    default: return 'info';
  }
};

// 获取状态文本
const getStatusText = (status) => {
  switch (status) {
    case 'active': return '活跃';
    case 'inactive': return '不活跃';
    case 'error': return '错误';
    default: return '未知';
  }
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
};

onMounted(() => {
  // 实际项目中应该从API获取环境列表
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

h4 {
  margin-top: 20px;
  margin-bottom: 10px;
  font-size: 16px;
  font-weight: 600;
}

.section-desc {
  margin-top: -5px;
  margin-bottom: 15px;
  color: #909399;
  font-size: 14px;
}

.variable-item {
  margin-bottom: 10px;
}

.add-variable {
  margin-top: 15px;
  display: flex;
  justify-content: center;
}
</style>
