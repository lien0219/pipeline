<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>创建流水线</h2>
        <p>配置新的CI/CD流水线</p>
      </div>
    </div>

    <el-card>
      <el-form
          ref="pipelineForm"
          :model="formState"
          :rules="rules"
          label-position="top"
          @submit.prevent="submitForm"
      >
        <h3>基本信息</h3>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="流水线名称" prop="name">
              <el-input v-model="formState.name" placeholder="请输入流水线名称" />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="Git 仓库" prop="git_repo">
              <el-input v-model="formState.git_repo" placeholder="请输入Git仓库地址" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="Git 分支" prop="git_branch">
              <el-input v-model="formState.git_branch" placeholder="请输入Git分支，默认为main" />
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="描述">
              <el-input
                  v-model="formState.description"
                  type="textarea"
                  :rows="2"
                  placeholder="请输入流水线描述（可选）"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <h3>流水线阶段</h3>
        <p class="section-desc">流水线由多个阶段组成，每个阶段包含一个或多个作业</p>

        <div v-for="(stage, stageIndex) in formState.stages" :key="stageIndex" class="stage-container">
          <div class="stage-header">
            <h4>阶段 {{ stageIndex + 1 }}</h4>
            <el-button
                type="danger"
                size="small"
                circle
                @click="removeStage(stageIndex)"
                :disabled="formState.stages.length <= 1"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>

          <el-row :gutter="20">
            <el-col :span="24">
              <el-form-item
                  :label="`阶段名称`"
                  :prop="`stages.${stageIndex}.name`"
                  :rules="[{ required: true, message: '请输入阶段名称', trigger: 'blur' }]"
              >
                <el-input v-model="stage.name" placeholder="请输入阶段名称" />
              </el-form-item>
            </el-col>
          </el-row>

          <h5>作业列表</h5>

          <div
              v-for="(job, jobIndex) in stage.jobs"
              :key="jobIndex"
              class="job-container"
          >
            <div class="job-header">
              <h6>作业 {{ jobIndex + 1 }}</h6>
              <el-button
                  type="danger"
                  size="small"
                  circle
                  @click="removeJob(stageIndex, jobIndex)"
                  :disabled="stage.jobs.length <= 1"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item
                    :label="`作业名称`"
                    :prop="`stages.${stageIndex}.jobs.${jobIndex}.name`"
                    :rules="[{ required: true, message: '请输入作业名称', trigger: 'blur' }]"
                >
                  <el-input v-model="job.name" placeholder="请输入作业名称" />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item
                    :label="`Docker 镜像（可选）`"
                    :prop="`stages.${stageIndex}.jobs.${jobIndex}.image`"
                >
                  <el-input v-model="job.image" placeholder="请输入Docker镜像，例如：node:16" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="24">
                <el-form-item
                    :label="`命令`"
                    :prop="`stages.${stageIndex}.jobs.${jobIndex}.command`"
                    :rules="[{ required: true, message: '请输入命令', trigger: 'blur' }]"
                >
                  <el-input
                      v-model="job.command"
                      type="textarea"
                      :rows="3"
                      placeholder="请输入要执行的命令"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item
                    :label="`超时时间（秒）`"
                    :prop="`stages.${stageIndex}.jobs.${jobIndex}.timeout`"
                >
                  <el-input-number
                      v-model="job.timeout"
                      :min="0"
                      :step="60"
                      placeholder="默认3600秒"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </div>

          <div class="add-job-button">
            <el-button type="primary" plain @click="addJob(stageIndex)">
              <el-icon><Plus /></el-icon>
              添加作业
            </el-button>
          </div>
        </div>

        <div class="add-stage-button">
          <el-button type="primary" @click="addStage">
            <el-icon><Plus /></el-icon>
            添加阶段
          </el-button>
        </div>

        <div class="form-actions">
          <el-button @click="$router.push('/pipelines')">取消</el-button>
          <el-button type="primary" native-type="submit" :loading="loading">创建流水线</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { usePipelineStore } from '@/stores/pipeline';
import { Plus, Delete } from '@element-plus/icons-vue';

const router = useRouter();
const pipelineStore = usePipelineStore();
const pipelineForm = ref(null);
const loading = ref(false);

// 表单数据
const formState = reactive({
  name: '',
  description: '',
  git_repo: '',
  git_branch: 'main',
  stages: [
    {
      name: '构建',
      jobs: [
        {
          name: '编译',
          command: 'npm install && npm run build',
          image: 'node:16',
          timeout: 3600
        }
      ]
    }
  ]
});

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '请输入流水线名称', trigger: 'blur' },
    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  git_repo: [
    { required: true, message: '请输入Git仓库地址', trigger: 'blur' }
  ],
  git_branch: [
    { required: true, message: '请输入Git分支', trigger: 'blur' }
  ]
};

// 添加阶段
const addStage = () => {
  formState.stages.push({
    name: `阶段 ${formState.stages.length + 1}`,
    jobs: [
      {
        name: '默认作业',
        command: '',
        image: '',
        timeout: 3600
      }
    ]
  });
};

// 移除阶段
const removeStage = (index) => {
  if (formState.stages.length > 1) {
    formState.stages.splice(index, 1);
  }
};

// 添加作业
const addJob = (stageIndex) => {
  formState.stages[stageIndex].jobs.push({
    name: `作业 ${formState.stages[stageIndex].jobs.length + 1}`,
    command: '',
    image: '',
    timeout: 3600
  });
};

// 移除作业
const removeJob = (stageIndex, jobIndex) => {
  if (formState.stages[stageIndex].jobs.length > 1) {
    formState.stages[stageIndex].jobs.splice(jobIndex, 1);
  }
};

// 提交表单
const submitForm = async () => {
  if (!pipelineForm.value) return;

  await pipelineForm.value.validate(async (valid) => {
    if (valid) {
      loading.value = true;
      try {
        await pipelineStore.createPipeline(formState);
        ElMessage.success('流水线创建成功');
        router.push('/pipelines');
      } catch (error) {
        console.error('Failed to create pipeline:', error);
        ElMessage.error('创建流水线失败');
      } finally {
        loading.value = false;
      }
    } else {
      ElMessage.warning('请填写必填项');
      return false;
    }
  });
};
</script>

<style scoped>
.page-header {
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

h3 {
  margin-top: 0;
  margin-bottom: 20px;
  font-size: 18px;
  font-weight: 600;
}

.section-desc {
  margin-top: -15px;
  margin-bottom: 20px;
  color: #909399;
  font-size: 14px;
}

.stage-container {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
  background-color: #f8f9fa;
}

.stage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.stage-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

h5 {
  margin-top: 20px;
  margin-bottom: 15px;
  font-size: 15px;
  font-weight: 600;
}

.job-container {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
}

.job-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.job-header h6 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.add-job-button {
  margin-top: 15px;
  display: flex;
  justify-content: center;
}

.add-stage-button {
  margin: 20px 0;
  display: flex;
  justify-content: center;
}

.form-actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
