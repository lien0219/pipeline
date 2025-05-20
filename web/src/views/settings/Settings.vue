<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>系统设置</h2>
        <p>配置系统参数和全局设置</p>
      </div>
    </div>

    <el-card>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本设置" name="basic">
          <el-form
              ref="basicForm"
              :model="basicSettings"
              label-position="top"
              :disabled="!isEditing"
          >
            <h3>系统信息</h3>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="系统名称">
                  <el-input v-model="basicSettings.system_name" placeholder="系统名称" />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="系统版本">
                  <el-input v-model="basicSettings.version" placeholder="系统版本" disabled />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="系统描述">
              <el-input
                  v-model="basicSettings.description"
                  type="textarea"
                  :rows="3"
                  placeholder="系统描述"
              />
            </el-form-item>

            <h3>管理员联系信息</h3>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="管理员邮箱">
                  <el-input v-model="basicSettings.admin_email" placeholder="管理员邮箱" />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="管理员电话">
                  <el-input v-model="basicSettings.admin_phone" placeholder="管理员电话" />
                </el-form-item>
              </el-col>
            </el-row>

            <div class="form-actions" v-if="!isEditing">
              <el-button type="primary" @click="startEditing">
                <el-icon><Edit /></el-icon>
                编辑设置
              </el-button>
            </div>

            <div class="form-actions" v-else>
              <el-button @click="cancelEditing">取消</el-button>
              <el-button type="primary" @click="saveBasicSettings" :loading="saving">
                保存设置
              </el-button>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="邮件设置" name="email">
          <el-form
              ref="emailForm"
              :model="emailSettings"
              label-position="top"
              :disabled="!isEditingEmail"
          >
            <h3>SMTP 配置</h3>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="SMTP 服务器">
                  <el-input v-model="emailSettings.smtp_server" placeholder="SMTP 服务器地址" />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="SMTP 端口">
                  <el-input v-model="emailSettings.smtp_port" placeholder="SMTP 端口" />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="SMTP 用户名">
                  <el-input v-model="emailSettings.smtp_username" placeholder="SMTP 用户名" />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="SMTP 密码">
                  <el-input
                      v-model="emailSettings.smtp_password"
                      type="password"
                      placeholder="SMTP 密码"
                      show-password
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-form-item label="发件人邮箱">
              <el-input v-model="emailSettings.from_email" placeholder="发件人邮箱" />
            </el-form-item>

            <el-form-item label="启用 SSL/TLS">
              <el-switch v-model="emailSettings.use_ssl" />
            </el-form-item>

            <h3>通知设置</h3>

            <el-form-item label="启用邮件通知">
              <el-switch v-model="emailSettings.enable_notifications" />
            </el-form-item>

            <el-form-item label="通知事件">
              <el-checkbox-group v-model="emailSettings.notification_events">
                <el-checkbox label="pipeline_success">流水线成功</el-checkbox>
                <el-checkbox label="pipeline_failure">流水线失败</el-checkbox>
                <el-checkbox label="deployment_success">部署成功</el-checkbox>
                <el-checkbox label="deployment_failure">部署失败</el-checkbox>
                <el-checkbox label="system_error">系统错误</el-checkbox>
              </el-checkbox-group>
            </el-form-item>

            <div class="form-actions" v-if="!isEditingEmail">
              <el-button type="primary" @click="startEditingEmail">
                <el-icon><Edit /></el-icon>
                编辑设置
              </el-button>
              <el-button type="success" @click="testEmailSettings">
                <el-icon><Message /></el-icon>
                测试邮件设置
              </el-button>
            </div>

            <div class="form-actions" v-else>
              <el-button @click="cancelEditingEmail">取消</el-button>
              <el-button type="primary" @click="saveEmailSettings" :loading="savingEmail">
                保存设置
              </el-button>
            </div>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="集成设置" name="integrations">
          <el-form
              ref="integrationsForm"
              :model="integrationSettings"
              label-position="top"
          >
            <h3>Git 集成</h3>

            <el-collapse>
              <el-collapse-item title="GitHub 集成" name="github">
                <el-form-item label="启用 GitHub 集成">
                  <el-switch v-model="integrationSettings.github.enabled" />
                </el-form-item>

                <template v-if="integrationSettings.github.enabled">
                  <el-form-item label="GitHub App ID">
                    <el-input v-model="integrationSettings.github.app_id" placeholder="GitHub App ID" />
                  </el-form-item>

                  <el-form-item label="GitHub App 私钥">
                    <el-input
                        v-model="integrationSettings.github.private_key"
                        type="textarea"
                        :rows="3"
                        placeholder="GitHub App 私钥"
                    />
                  </el-form-item>

                  <el-form-item label="Webhook 密钥">
                    <el-input v-model="integrationSettings.github.webhook_secret" placeholder="Webhook 密钥" />
                  </el-form-item>

                  <el-form-item>
                    <el-button type="primary" @click="saveGitHubSettings" :loading="savingGitHub">
                      保存 GitHub 设置
                    </el-button>
                    <el-button type="success" @click="testGitHubSettings">
                      测试连接
                    </el-button>
                  </el-form-item>
                </template>
              </el-collapse-item>

              <el-collapse-item title="GitLab 集成" name="gitlab">
                <el-form-item label="启用 GitLab 集成">
                  <el-switch v-model="integrationSettings.gitlab.enabled" />
                </el-form-item>

                <template v-if="integrationSettings.gitlab.enabled">
                  <el-form-item label="GitLab URL">
                    <el-input v-model="integrationSettings.gitlab.url" placeholder="GitLab URL" />
                  </el-form-item>

                  <el-form-item label="GitLab API Token">
                    <el-input v-model="integrationSettings.gitlab.token" placeholder="GitLab API Token" show-password />
                  </el-form-item>

                  <el-form-item label="Webhook 密钥">
                    <el-input v-model="integrationSettings.gitlab.webhook_secret" placeholder="Webhook 密钥" />
                  </el-form-item>

                  <el-form-item>
                    <el-button type="primary" @click="saveGitLabSettings" :loading="savingGitLab">
                      保存 GitLab 设置
                    </el-button>
                    <el-button type="success" @click="testGitLabSettings">
                      测试连接
                    </el-button>
                  </el-form-item>
                </template>
              </el-collapse-item>
            </el-collapse>

            <h3>容器仓库集成</h3>

            <el-collapse>
              <el-collapse-item title="Docker Registry 集成" name="docker">
                <el-form-item label="启用 Docker Registry 集成">
                  <el-switch v-model="integrationSettings.docker.enabled" />
                </el-form-item>

                <template v-if="integrationSettings.docker.enabled">
                  <el-form-item label="Registry URL">
                    <el-input v-model="integrationSettings.docker.url" placeholder="Registry URL" />
                  </el-form-item>

                  <el-form-item label="用户名">
                    <el-input v-model="integrationSettings.docker.username" placeholder="用户名" />
                  </el-form-item>

                  <el-form-item label="密码">
                    <el-input v-model="integrationSettings.docker.password" type="password" placeholder="密码" show-password />
                  </el-form-item>

                  <el-form-item>
                    <el-button type="primary" @click="saveDockerSettings" :loading="savingDocker">
                      保存 Docker Registry 设置
                    </el-button>
                    <el-button type="success" @click="testDockerSettings">
                      测试连接
                    </el-button>
                  </el-form-item>
                </template>
              </el-collapse-item>
            </el-collapse>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="系统维护" name="maintenance">
          <div class="maintenance-section">
            <h3>系统状态</h3>

            <el-descriptions :column="2" border>
              <el-descriptions-item label="系统运行时间">{{ systemStatus.uptime }}</el-descriptions-item>
              <el-descriptions-item label="CPU 使用率">{{ systemStatus.cpu_usage }}</el-descriptions-item>
              <el-descriptions-item label="内存使用率">{{ systemStatus.memory_usage }}</el-descriptions-item>
              <el-descriptions-item label="磁盘使用率">{{ systemStatus.disk_usage }}</el-descriptions-item>
              <el-descriptions-item label="数据库连接数">{{ systemStatus.db_connections }}</el-descriptions-item>
              <el-descriptions-item label="活跃用户数">{{ systemStatus.active_users }}</el-descriptions-item>
            </el-descriptions>

            <div class="action-buttons">
              <el-button type="primary" @click="refreshSystemStatus">
                <el-icon><Refresh /></el-icon>
                刷新状态
              </el-button>
            </div>
          </div>

          <div class="maintenance-section">
            <h3>数据库维护</h3>

            <div class="action-buttons">
              <el-button type="primary" @click="backupDatabase" :loading="backingUp">
                <el-icon><Download /></el-icon>
                备份数据库
              </el-button>

              <el-button type="warning" @click="optimizeDatabase" :loading="optimizing">
                <el-icon><Magic /></el-icon>
                优化数据库
              </el-button>

              <el-upload
                  class="upload-button"
                  action="#"
                  :auto-upload="false"
                  :on-change="handleDatabaseFileChange"
                  :limit="1"
              >
                <el-button type="danger">
                  <el-icon><Upload /></el-icon>
                  恢复数据库
                </el-button>
              </el-upload>
            </div>
          </div>

          <div class="maintenance-section">
            <h3>日志管理</h3>

            <el-table :data="logs" style="width: 100%" v-loading="loadingLogs">
              <el-table-column prop="name" label="日志文件" min-width="200" />
              <el-table-column prop="size" label="大小" width="120">
                <template #default="{ row }">
                  {{ formatFileSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="modified" label="修改时间" width="180" />
              <el-table-column label="操作" width="200" fixed="right">
                <template #default="{ row }">
                  <el-button
                      link
                      type="primary"
                      size="small"
                      @click="downloadLog(row.name)"
                  >
                    下载
                  </el-button>

                  <el-button
                      link
                      type="primary"
                      size="small"
                      @click="viewLog(row.name)"
                  >
                    查看
                  </el-button>

                  <el-button
                      link
                      type="danger"
                      size="small"
                      @click="deleteLog(row.name)"
                  >
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <div class="action-buttons">
              <el-button type="primary" @click="refreshLogs">
                <el-icon><Refresh /></el-icon>
                刷新日志
              </el-button>

              <el-button type="danger" @click="clearAllLogs">
                <el-icon><Delete /></el-icon>
                清除所有日志
              </el-button>
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    &lt;!-- 查看日志对话框 -->
    <el-dialog
        v-model="logDialogVisible"
        :title="`查看日志: ${currentLog}`"
        width="80%"
    >
      <div v-loading="loadingLogContent">
        <div class="log-container">
          <pre class="log-content">{{ logContent }}</pre>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Edit, Message, Refresh, Download, Upload, Magic, Delete } from '@element-plus/icons-vue';

const activeTab = ref('basic');
const isEditing = ref(false);
const isEditingEmail = ref(false);
const saving = ref(false);
const savingEmail = ref(false);
const savingGitHub = ref(false);
const savingGitLab = ref(false);
const savingDocker = ref(false);
const backingUp = ref(false);
const optimizing = ref(false);
const loadingLogs = ref(false);
const logDialogVisible = ref(false);
const loadingLogContent = ref(false);
const currentLog = ref('');
const logContent = ref('');

// 基本设置
const basicSettings = reactive({
  system_name: 'CI/CD 流水线平台',
  version: 'v1.0.0',
  description: '一个强大的持续集成和持续部署平台，支持多种构建和部署场景。',
  admin_email: 'admin@example.com',
  admin_phone: '13800138000'
});

// 邮件设置
const emailSettings = reactive({
  smtp_server: 'smtp.example.com',
  smtp_port: '587',
  smtp_username: 'notification@example.com',
  smtp_password: 'password123',
  from_email: 'notification@example.com',
  use_ssl: true,
  enable_notifications: true,
  notification_events: ['pipeline_failure', 'deployment_failure', 'system_error']
});

// 集成设置
const integrationSettings = reactive({
  github: {
    enabled: false,
    app_id: '',
    private_key: '',
    webhook_secret: ''
  },
  gitlab: {
    enabled: true,
    url: 'https://gitlab.example.com',
    token: 'glpat-XXXXXXXXXXXXXXXXXXXX',
    webhook_secret: 'webhook-secret-123'
  },
  docker: {
    enabled: true,
    url: 'https://registry.example.com',
    username: 'registry-user',
    password: 'registry-password'
  }
});

// 系统状态
const systemStatus = reactive({
  uptime: '10天 5小时 30分钟',
  cpu_usage: '25%',
  memory_usage: '40%',
  disk_usage: '60%',
  db_connections: '15',
  active_users: '8'
});

// 日志列表
const logs = ref([
  { name: 'system.log', size: 1048576, modified: '2023-05-18 10:30:00' },
  { name: 'access.log', size: 5242880, modified: '2023-05-18 10:30:00' },
  { name: 'error.log', size: 524288, modified: '2023-05-18 10:30:00' },
  { name: 'pipeline.log', size: 2097152, modified: '2023-05-18 10:30:00' },
  { name: 'deployment.log', size: 1572864, modified: '2023-05-18 10:30:00' }
]);

// 开始编辑基本设置
const startEditing = () => {
  isEditing.value = true;
};

// 取消编辑基本设置
const cancelEditing = () => {
  isEditing.value = false;
  // 实际项目中应该重新获取设置数据
};

// 保存基本设置
const saveBasicSettings = async () => {
  saving.value = true;

  try {
    // 实际项目中应该调用API保存设置
    await new Promise(resolve => setTimeout(resolve, 1000));
    ElMessage.success('基本设置已保存');
    isEditing.value = false;
  } catch (error) {
    console.error('保存基本设置失败:', error);
    ElMessage.error('保存基本设置失败');
  } finally {
    saving.value = false;
  }
};

// 开始编辑邮件设置
const startEditingEmail = () => {
  isEditingEmail.value = true;
};

// 取消编辑邮件设置
const cancelEditingEmail = () => {
  isEditingEmail.value = false;
  // 实际项目中应该重新获取设置数据
};

// 保存邮件设置
const saveEmailSettings = async () => {
  savingEmail.value = true;

  try {
    // 实际项目中应该调用API保存设置
    await new Promise(resolve => setTimeout(resolve, 1000));
    ElMessage.success('邮件设置已保存');
    isEditingEmail.value = false;
  } catch (error) {
    console.error('保存邮件设置失败:', error);
    ElMessage.error('保存邮件设置失败');
  } finally {
    savingEmail.value = false;
  }
};

// 测试邮件设置
const testEmailSettings = async () => {
  try {
    // 实际项目中应该调用API测试邮件设置
    await new Promise(resolve => setTimeout(resolve, 2000));
    ElMessage.success('测试邮件已发送');
  } catch (error) {
    console.error('测试邮件发送失败:', error);
    ElMessage.error('测试邮件发送失败');
  }
};

// 保存 GitHub 设置
const saveGitHubSettings = async () => {
  savingGitHub.value = true;

  try {
    // 实际项目中应该调用API保存设置
    await new Promise(resolve => setTimeout(resolve, 1000));
    ElMessage.success('GitHub 设置已保存');
  } catch (error) {
    console.error('保存 GitHub 设置失败:', error);
    ElMessage.error('保存 GitHub 设置失败');
  } finally {
    savingGitHub.value = false;
  }
};

// 测试 GitHub 设置
const testGitHubSettings = async () => {
  try {
    // 实际项目中应该调用API测试 GitHub 设置
    await new Promise(resolve => setTimeout(resolve, 2000));
    ElMessage.success('GitHub 连接测试成功');
  } catch (error) {
    console.error('GitHub 连接测试失败:', error);
    ElMessage.error('GitHub 连接测试失败');
  }
};

// 保存 GitLab 设置
const saveGitLabSettings = async () => {
  savingGitLab.value = true;

  try {
    // 实际项目中应该调用API保存设置
    await new Promise(resolve => setTimeout(resolve, 1000));
    ElMessage.success('GitLab 设置已保存');
  } catch (error) {
    console.error('保存 GitLab 设置失败:', error);
    ElMessage.error('保存 GitLab 设置失败');
  } finally {
    savingGitLab.value = false;
  }
};

// 测试 GitLab 设置
const testGitLabSettings = async () => {
  try {
    // 实际项目中应该调用API测试 GitLab 设置
    await new Promise(resolve => setTimeout(resolve, 2000));
    ElMessage.success('GitLab 连接测试成功');
  } catch (error) {
    console.error('GitLab 连接测试失败:', error);
    ElMessage.error('GitLab 连接测试失败');
  }
};

// 保存 Docker Registry 设置
const saveDockerSettings = async () => {
  savingDocker.value = true;

  try {
    // 实际项目中应该调用API保存设置
    await new Promise(resolve => setTimeout(resolve, 1000));
    ElMessage.success('Docker Registry 设置已保存');
  } catch (error) {
    console.error('保存 Docker Registry 设置失败:', error);
    ElMessage.error('保存 Docker Registry 设置失败');
  } finally {
    savingDocker.value = false;
  }
};

// 测试 Docker Registry 设置
const testDockerSettings = async () => {
  try {
    // 实际项目中应该调用API测试 Docker Registry 设置
    await new Promise(resolve => setTimeout(resolve, 2000));
    ElMessage.success('Docker Registry 连接测试成功');
  } catch (error) {
    console.error('Docker Registry 连接测试失败:', error);
    ElMessage.error('Docker Registry 连接测试失败');
  }
};

// 刷新系统状态
const refreshSystemStatus = async () => {
  try {
    // 实际项目中应该调用API获取系统状态
    await new Promise(resolve => setTimeout(resolve, 1000));

    // 模拟更新系统状态
    systemStatus.cpu_usage = `${Math.floor(Math.random() * 50) + 10}%`;
    systemStatus.memory_usage = `${Math.floor(Math.random() * 50) + 20}%`;
    systemStatus.active_users = `${Math.floor(Math.random() * 20) + 5}`;

    ElMessage.success('系统状态已刷新');
  } catch (error) {
    console.error('刷新系统状态失败:', error);
    ElMessage.error('刷新系统状态失败');
  }
};

// 备份数据库
const backupDatabase = async () => {
  backingUp.value = true;

  try {
    // 实际项目中应该调用API备份数据库
    await new Promise(resolve => setTimeout(resolve, 3000));
    ElMessage.success('数据库备份成功');
  } catch (error) {
    console.error('数据库备份失败:', error);
    ElMessage.error('数据库备份失败');
  } finally {
    backingUp.value = false;
  }
};

// 优化数据库
const optimizeDatabase = async () => {
  optimizing.value = true;

  try {
    // 实际项目中应该调用API优化数据库
    await new Promise(resolve => setTimeout(resolve, 3000));
    ElMessage.success('数据库优化成功');
  } catch (error) {
    console.error('数据库优化失败:', error);
    ElMessage.error('数据库优化失败');
  } finally {
    optimizing.value = false;
  }
};

// 处理数据库文件变更
const handleDatabaseFileChange = async (file) => {
  try {
    await ElMessageBox.confirm('确定要恢复数据库吗？此操作将覆盖当前数据库，且不可恢复。', '恢复确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    // 实际项目中应该调用API上传并恢复数据库
    await new Promise(resolve => setTimeout(resolve, 3000));
    ElMessage.success('数据库恢复成功');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('数据库恢复失败:', error);
      ElMessage.error('数据库恢复失败');
    }
  }
};

// 刷新日志列表
const refreshLogs = async () => {
  loadingLogs.value = true;

  try {
    // 实际项目中应该调用API获取日志列表
    await new Promise(resolve => setTimeout(resolve, 1000));

    // 模拟更新日志大小和修改时间
    logs.value = logs.value.map(log => ({
      ...log,
      size: Math.floor(Math.random() * 5242880) + 524288,
      modified: new Date().toLocaleString()
    }));

    ElMessage.success('日志列表已刷新');
  } catch (error) {
    console.error('刷新日志列表失败:', error);
    ElMessage.error('刷新日志列表失败');
  } finally {
    loadingLogs.value = false;
  }
};

// 下载日志
const downloadLog = (name) => {
  // 实际项目中应该调用API下载日志
  ElMessage.success(`开始下载日志: ${name}`);
};

// 查看日志
const viewLog = async (name) => {
  currentLog.value = name;
  logDialogVisible.value = true;
  loadingLogContent.value = true;

  try {
    // 实际项目中应该调用API获取日志内容
    await new Promise(resolve => setTimeout(resolve, 1000));

    // 模拟日志内容
    logContent.value = `[2023-05-18 10:00:00] [INFO] 系统启动\n[2023-05-18 10:01:00] [INFO] 用户 admin 登录\n[2023-05-18 10:05:00] [INFO] 创建流水线 #123\n[2023-05-18 10:10:00] [INFO] 流水线 #123 开始运行\n[2023-05-18 10:15:00] [INFO] 流水线 #123 构建阶段完成\n[2023-05-18 10:20:00] [INFO] 流水线 #123 测试阶段完成\n[2023-05-18 10:25:00] [INFO] 流水线 #123 部署阶段完成\n[2023-05-18 10:30:00] [INFO] 流水线 #123 运行成功`;
  } catch (error) {
    console.error('获取日志内容失败:', error);
    logContent.value = '获取日志内容失败';
  } finally {
    loadingLogContent.value = false;
  }
};

// 删除日志
const deleteLog = async (name) => {
  try {
    await ElMessageBox.confirm(`确定要删除日志 ${name} 吗？此操作不可恢复。`, '删除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    // 实际项目中应该调用API删除日志
    logs.value = logs.value.filter(log => log.name !== name);
    ElMessage.success(`日志 ${name} 已删除`);
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除日志失败:', error);
      ElMessage.error('删除日志失败');
    }
  }
};

// 清除所有日志
const clearAllLogs = async () => {
  try {
    await ElMessageBox.confirm('确定要清除所有日志吗？此操作不可恢复。', '清除确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });

    // 实际项目中应该调用API清除所有日志
    logs.value = [];
    ElMessage.success('所有日志已清除');
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清除日志失败:', error);
      ElMessage.error('清除日志失败');
    }
  }
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '0 B';

  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

onMounted(() => {
  // 实际项目中应该从API获取设置数据
});
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
  margin-top: 20px;
  margin-bottom: 20px;
  font-size: 18px;
  font-weight: 600;
}

.form-actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.maintenance-section {
  margin-bottom: 30px;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.upload-button {
  display: inline-block;
}

.log-container {
  background-color: #1e1e1e;
  color: #f8f8f8;
  border-radius: 4px;
  padding: 15px;
  height: 500px;
  overflow-y: auto;
}

.log-content {
  font-family: 'Courier New', Courier, monospace;
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
}
</style>
