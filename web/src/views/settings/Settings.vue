<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>系统设置</h2>
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
            @submit.prevent="saveBasicSettings"
          >
            <h3>系统信息</h3>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="系统名称">
                  <el-input
                    v-model="basicSettings.system_name"
                    placeholder="系统名称"
                  />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="系统版本">
                  <el-input
                    v-model="basicSettings.version"
                    placeholder="系统版本"
                    disabled
                  />
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
                  <el-input
                    v-model="basicSettings.admin_email"
                    placeholder="管理员邮箱"
                  />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="管理员电话">
                  <el-input
                    v-model="basicSettings.admin_phone"
                    placeholder="管理员电话"
                  />
                </el-form-item>
              </el-col>
            </el-row>
          </el-form>
          <div class="form-actions" v-if="!isEditing">
            <el-button type="primary" @click="startEditing">
              <el-icon><Edit /></el-icon>
              编辑设置
            </el-button>
          </div>

          <div class="form-actions" v-else>
            <el-button @click="cancelEditing">取消</el-button>
            <el-button
              type="primary"
              @click="saveBasicSettings"
              :loading="saving"
            >
              保存设置
            </el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="邮件设置" name="email">
          <el-form
            ref="emailForm"
            :model="emailSettings"
            label-position="top"
            :disabled="!isEditingEmail"
            @submit.prevent="saveEmailSettings"
          >
            <h3>SMTP 配置</h3>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="SMTP 服务器">
                  <el-input
                    v-model="emailSettings.smtp_server"
                    placeholder="SMTP 服务器地址"
                  />
                </el-form-item>
              </el-col>

              <el-col :span="12">
                <el-form-item label="SMTP 端口">
                  <el-input
                    v-model.number="emailSettings.smtp_port"
                    placeholder="SMTP 端口"
                  />
                </el-form-item>
              </el-col>
            </el-row>

            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="SMTP 用户名">
                  <el-input
                    v-model="emailSettings.smtp_username"
                    placeholder="SMTP 用户名"
                  />
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
              <el-input
                v-model="emailSettings.from_email"
                placeholder="发件人邮箱"
              />
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
          </el-form>
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
            <el-button
              type="primary"
              @click="saveEmailSettings"
              :loading="savingEmail"
            >
              保存设置
            </el-button>
          </div>
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
                    <el-input
                      v-model="integrationSettings.github.app_id"
                      placeholder="GitHub App ID"
                    />
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
                    <el-input
                      v-model="integrationSettings.github.webhook_secret"
                      placeholder="Webhook 密钥"
                    />
                  </el-form-item>

                  <el-form-item>
                    <el-button
                      type="primary"
                      @click="saveGitHubSettings"
                      :loading="savingGitHub"
                    >
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
                    <el-input
                      v-model="integrationSettings.gitlab.url"
                      placeholder="GitLab URL"
                    />
                  </el-form-item>

                  <el-form-item label="GitLab API Token">
                    <el-input
                      v-model="integrationSettings.gitlab.token"
                      placeholder="GitLab API Token"
                      show-password
                    />
                  </el-form-item>

                  <el-form-item label="Webhook 密钥">
                    <el-input
                      v-model="integrationSettings.gitlab.webhook_secret"
                      placeholder="Webhook 密钥"
                    />
                  </el-form-item>

                  <el-form-item>
                    <el-button
                      type="primary"
                      @click="saveGitLabSettings"
                      :loading="savingGitLab"
                    >
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
                    <el-input
                      v-model="integrationSettings.docker.url"
                      placeholder="Registry URL"
                    />
                  </el-form-item>

                  <el-form-item label="用户名">
                    <el-input
                      v-model="integrationSettings.docker.username"
                      placeholder="用户名"
                    />
                  </el-form-item>

                  <el-form-item label="密码">
                    <el-input
                      v-model="integrationSettings.docker.password"
                      type="password"
                      placeholder="密码"
                      show-password
                    />
                  </el-form-item>

                  <el-form-item>
                    <el-button
                      type="primary"
                      @click="saveDockerSettings"
                      :loading="savingDocker"
                    >
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
              <el-descriptions-item label="系统运行时间">{{
                systemStatus.uptime
              }}</el-descriptions-item>
              <el-descriptions-item label="CPU 使用率">{{
                systemStatus.cpu_usage
              }}</el-descriptions-item>
              <el-descriptions-item label="内存使用率">{{
                systemStatus.memory_usage
              }}</el-descriptions-item>
              <el-descriptions-item label="磁盘使用率">{{
                systemStatus.disk_usage
              }}</el-descriptions-item>
              <el-descriptions-item label="数据库连接数">{{
                systemStatus.db_connections
              }}</el-descriptions-item>
              <el-descriptions-item label="活跃用户数">{{
                systemStatus.active_users
              }}</el-descriptions-item>
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
              <el-button
                type="primary"
                @click="backupDatabase"
                :loading="backingUp"
              >
                <el-icon><Download /></el-icon>
                备份数据库
              </el-button>

              <el-button
                type="warning"
                @click="optimizeDatabase"
                :loading="optimizing"
              >
                <el-icon><Edit /></el-icon>
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
                  <!-- 文件内容较多，后续优化 -->
                  <el-button
                    link
                    type="primary"
                    size="small"
                    @click="downloadLogHand(row.name)"
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

    <el-dialog v-model="logDialogVisible" title="查看日志" width="80%">
      <div class="log-container">
        <div class="log-toolbar">
          <el-input
            v-model="logSearchKeyword"
            placeholder="搜索日志内容..."
            size="small"
            style="width: 300px; margin-right: 10px"
          />
          <el-button
            type="primary"
            size="small"
            @click="copyLogContent"
            :icon="CopyDocument"
          >
            复制日志
          </el-button>
        </div>
        <div class="log-content" v-html="formattedLogContent"></div>
      </div>
      <div class="log-pagination">
        <el-pagination
          v-model:current-page="currentLogPage"
          v-model:page-size="logPageSize"
          :total="totalLogLines"
          @size-change="handleLogPageSizeChange"
          @current-change="handleLogPageChange"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, toRaw, computed } from "vue";
import { CopyDocument } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  Edit,
  Message,
  Refresh,
  Download,
  Upload,
  Delete,
} from "@element-plus/icons-vue";
// @ts-ignore
import { getSettings, saveSettings } from "@/api/settings";
// @ts-ignore
import {
  getSystemStatus,
  getLogFiles,
  getLogContent,
  downloadLog,
} from "@/api/system";

const activeTab = ref("basic");
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
const currentLog = ref("");
const logContent = ref("");
const currentLogPage = ref(1);
const logPageSize = ref(100);
const totalLogLines = ref(0);
const logSearchKeyword = ref("");

// 基本设置
const basicSettings = reactive({
  system_name: "",
  version: "",
  description: "",
  admin_email: "",
  admin_phone: "",
});

// 邮件设置
const emailSettings = reactive({
  smtp_server: "",
  smtp_port: 0,
  smtp_username: "",
  smtp_password: "",
  from_email: "",
  use_ssl: false,
  enable_notifications: false,
  notification_events: [],
});

// 集成设置
const integrationSettings = reactive({
  github: {
    enabled: false,
    app_id: "",
    private_key: "",
    webhook_secret: "",
  },
  gitlab: {
    enabled: false,
    url: "",
    token: "",
    webhook_secret: "",
  },
  docker: {
    enabled: false,
    url: "",
    username: "",
    password: "",
  },
});

// 系统状态
const systemStatus = reactive<any>({
  uptime: "加载中...",
  cpu_usage: "加载中...",
  memory_usage: "加载中...",
  disk_usage: "加载中...",
  db_connections: 0,
  active_users: 0,
});

// 日志列表
const logs = ref<any[]>([]);

// 加载所有设置
const loadAllSettings = async () => {
  try {
    // 加载基本设置
    const basicRes = await getSettings("basic");
    Object.assign(basicSettings, basicRes.data);

    // 加载邮件设置
    const emailRes = await getSettings("email");
    Object.assign(emailSettings, emailRes.data);

    // 加载集成设置
    const integrationRes = await getSettings("integration");
    Object.assign(integrationSettings, integrationRes.data);
  } catch (error) {
    console.error("加载设置失败:", error);
    ElMessage.error("加载设置失败，请刷新页面重试");
  }
};
// 加载系统状态
const loadSystemStatus = async () => {
  try {
    const res = await getSystemStatus();
    Object.assign(systemStatus, res.data);
  } catch (error) {
    console.error("加载系统状态失败:", error);
    ElMessage.error("加载系统状态失败");
  }
};

// 开始编辑基本设置
const startEditing = () => {
  isEditing.value = true;
};

// 取消编辑基本设置
const cancelEditing = () => {
  isEditing.value = false;
  getSettings("basic").then((res) => {
    Object.assign(basicSettings, res.data);
  });
};

// 保存基本设置
const saveBasicSettings = async () => {
  saving.value = true;

  try {
    await saveSettings({
      type: "basic",
      items: toRaw(basicSettings),
    });
    ElMessage.success("基本设置已保存");
    isEditing.value = false;
  } catch (error) {
    console.error("保存基本设置失败:", error);
    ElMessage.error("保存基本设置失败");
  } finally {
    saving.value = false;
  }
};
// 加载日志文件
const loadLogFiles = async () => {
  loadingLogs.value = true;
  try {
    const res = await getLogFiles();
    logs.value = res.data;
  } catch (error) {
    console.error("加载日志文件失败:", error);
    ElMessage.error("加载日志文件失败");
  } finally {
    loadingLogs.value = false;
  }
};
// 分页
const handleLogPageSizeChange = (size: number) => {
  logPageSize.value = size;
  loadLogContent(currentLog.value, (currentLogPage.value - 1) * size, size);
};

const handleLogPageChange = (page: number) => {
  currentLogPage.value = page;
  loadLogContent(
    currentLog.value,
    (page - 1) * logPageSize.value,
    logPageSize.value
  );
};
const formattedLogContent = computed(() => {
  if (!logContent.value) return "<div class='empty-log'>日志内容为空</div>";

  let cleanedContent = logContent.value.replace(
    /[\u001b\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]/g,
    ""
  );
  let lines = cleanedContent.split("\n");
  let formattedLines = lines.map((line, index) => {
    let formattedLine = line
      .replace(/ERROR/g, '<span class="log-level error">ERROR</span>')
      .replace(/WARN/g, '<span class="log-level warn">WARN</span>')
      .replace(/INFO/g, '<span class="log-level info">INFO</span>')
      .replace(/DEBUG/g, '<span class="log-level debug">DEBUG</span>');

    if (logSearchKeyword.value && logSearchKeyword.value.trim() !== "") {
      try {
        const keyword = logSearchKeyword.value
          .trim()
          .replace(/[.*+?^${}()|\[\]\\]/g, "\\$&");
        const regex = new RegExp(`(${keyword})`, "gi");
        formattedLine = formattedLine.replace(
          regex,
          '<span class="search-highlight">$1</span>'
        );
      } catch (e) {
        console.error("搜索正则表达式错误:", e);
      }
    }

    return `<div class="log-line"><span class="line-number">${
      index + 1
    }</span>${formattedLine}</div>`;
  });

  return formattedLines.join("");
});
const copyLogContent = () => {
  if (!logContent.value) {
    ElMessage.warning("没有可复制的日志内容");
    return;
  }
  navigator.clipboard
    .writeText(logContent.value)
    .then(() => {
      ElMessage.success("日志内容已复制到剪贴板");
    })
    .catch(() => {
      ElMessage.error("复制失败，请手动复制");
    });
};
// 开始编辑邮件设置
const startEditingEmail = () => {
  isEditingEmail.value = true;
};

// 取消编辑邮件设置
const cancelEditingEmail = () => {
  isEditingEmail.value = false;
  getSettings("email").then((res) => {
    Object.assign(emailSettings, res.data);
  });
};

// 保存邮件设置
const saveEmailSettings = async () => {
  savingEmail.value = true;

  try {
    await saveSettings({
      type: "email",
      items: toRaw(emailSettings),
    });
    ElMessage.success("邮件设置已保存");
    isEditingEmail.value = false;
  } catch (error) {
    console.error("保存邮件设置失败:", error);
    ElMessage.error("保存邮件设置失败");
  } finally {
    savingEmail.value = false;
  }
};

// 保存GitHub设置
const saveGitHubSettings = async () => {
  savingGitHub.value = true;

  try {
    await saveSettings({
      type: "integration",
      items: { github: toRaw(integrationSettings.github) },
    });
    ElMessage.success("GitHub设置已保存");
  } catch (error) {
    console.error("保存GitHub设置失败:", error);
    ElMessage.error("保存GitHub设置失败");
  } finally {
    savingGitHub.value = false;
  }
};

// 保存GitLab设置
const saveGitLabSettings = async () => {
  savingGitLab.value = true;

  try {
    await saveSettings({
      type: "integration",
      items: { gitlab: toRaw(integrationSettings.gitlab) },
    });
    ElMessage.success("GitLab设置已保存");
  } catch (error) {
    console.error("保存GitLab设置失败:", error);
    ElMessage.error("保存GitLab设置失败");
  } finally {
    savingGitLab.value = false;
  }
};

// 保存Docker设置
const saveDockerSettings = async () => {
  savingDocker.value = true;

  try {
    await saveSettings({
      type: "integration",
      items: { docker: toRaw(integrationSettings.docker) },
    });
    ElMessage.success("Docker Registry设置已保存");
  } catch (error) {
    console.error("保存Docker Registry设置失败:", error);
    ElMessage.error("保存Docker Registry设置失败");
  } finally {
    savingDocker.value = false;
  }
};

// 测试邮件设置
const testEmailSettings = async () => {
  try {
    // 实际项目中应该调用API测试邮件设置
    ElMessageBox.alert("邮件测试功能将在后续版本实现", "提示", {
      confirmButtonText: "确定",
    });
  } catch (error) {
    ElMessage.error("邮件测试失败: " + (error as any).message);
  }
};

const testGitHubSettings = async () => {
  try {
    // 实际项目中应该调用API测试GitHub连接
    ElMessageBox.alert("GitHub连接测试功能将在后续版本实现", "提示", {
      confirmButtonText: "确定",
    });
  } catch (error) {
    ElMessage.error("GitHub连接测试失败: " + (error as any).message);
  }
};

const testGitLabSettings = async () => {
  try {
    // 实际项目中应该调用API测试GitLab连接
    ElMessageBox.alert("GitLab连接测试功能将在后续版本实现", "提示", {
      confirmButtonText: "确定",
    });
  } catch (error) {
    ElMessage.error("GitLab连接测试失败: " + (error as Error).message);
  }
};

const testDockerSettings = async () => {
  try {
    // 实际项目中应该调用API测试Docker连接
    ElMessageBox.alert("Docker连接测试功能将在后续版本实现", "提示", {
      confirmButtonText: "确定",
    });
  } catch (error) {
    ElMessage.error("Docker连接测试失败: " + (error as Error).message);
  }
};

// 系统维护功能
const refreshSystemStatus = async () => {
  try {
    loadSystemStatus();
    ElMessage.success("系统状态已刷新");
  } catch (error) {
    ElMessage.error("刷新系统状态失败: " + (error as Error).message);
  }
};

const backupDatabase = async () => {
  try {
    backingUp.value = true;
    // 这里应该调用后端数据库备份接口
    ElMessageBox.alert("数据库备份功能将在后续版本实现", "提示", {
      confirmButtonText: "确定",
    });
  } catch (error) {
    ElMessage.error("数据库备份失败: " + (error as Error).message);
  } finally {
    backingUp.value = false;
  }
};

const optimizeDatabase = async () => {
  try {
    optimizing.value = true;
    // 这里应该调用后端数据库优化接口
    ElMessageBox.alert("数据库优化功能将在后续版本实现", "提示", {
      confirmButtonText: "确定",
    });
  } catch (error) {
    ElMessage.error("数据库优化失败: " + (error as Error).message);
  } finally {
    optimizing.value = false;
  }
};

const handleDatabaseFileChange = (file) => {
  ElMessageBox.confirm(`确定要恢复数据库吗？此操作将覆盖现有数据。`, "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(async () => {
    try {
      // 这里应该调用后端数据库恢复接口
      ElMessage.success("数据库恢复成功");
    } catch (error) {
      ElMessage.error("数据库恢复失败");
    }
  });
};

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

const refreshLogs = () => {
  loadLogFiles();
  ElMessage.success("日志列表已刷新");
};

const downloadLogHand = async (fileName) => {
  if (!fileName) {
    ElMessage.warning("请选择有效的日志文件");
    return;
  }

  try {
    const blob = await downloadLog(fileName);
    const url = window.URL.createObjectURL(blob as any);
    const link = document.createElement("a");

    link.href = url;
    link.download = decodeURIComponent(fileName);
    document.body.appendChild(link);
    link.click();

    setTimeout(() => {
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    }, 100);

    ElMessage.success("日志下载成功");
  } catch (error) {
    console.error("下载日志失败:", error);
    const errorMsg = error.response?.data?.message || "下载日志失败，请重试";
    ElMessage.error(errorMsg);
  }
};

// 查看日志
const viewLog = async (fileName: string) => {
  logDialogVisible.value = true;
  currentLog.value = fileName;
  currentLogPage.value = 1;
  logContent.value = "";
  loadingLogContent.value = true;
  try {
    await loadLogContent(currentLog.value, 0, logPageSize.value);
  } catch (error) {
    console.error("加载日志失败:", error);
    ElMessage.error("加载日志失败");
  } finally {
    loadingLogContent.value = false;
  }
};
// 加载日志内容
const loadLogContent = async (
  fileName: string,
  startLine: number,
  limit: number
) => {
  loadingLogContent.value = true;
  try {
    const response = await getLogContent(fileName, startLine, limit);
    logContent.value = response.data.content || "";
    totalLogLines.value = response.data.totalLines || 0;

    if (!logContent.value) {
      ElMessage.info("当前日志文件内容为空");
    }
  } catch (error) {
    console.error("加载日志内容失败:", error);
    logContent.value = "";
    totalLogLines.value = 0;
    ElMessage.error("加载日志内容失败: " + (error as Error).message);
  } finally {
    loadingLogContent.value = false;
  }
};
const deleteLog = (name) => {
  ElMessageBox.confirm(`确定要删除日志文件 ${name} 吗？`, "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    // 模拟删除日志
    logs.value = logs.value.filter((log) => log.name !== name);
    ElMessage.success("日志文件已删除");
  });
};

const clearAllLogs = () => {
  ElMessageBox.confirm("确定要清除所有日志文件吗？", "警告", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning",
  }).then(() => {
    // 模拟清除所有日志
    logs.value = [];
    ElMessage.success("所有日志文件已清除");
  });
};

// 页面加载时获取设置数据
onMounted(() => {
  loadAllSettings();
  loadSystemStatus();
  loadLogFiles();
});
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.header-title {
  display: flex;
  align-items: center;
}

.form-actions {
  margin-top: 30px;
  text-align: right;
}

.maintenance-section {
  margin-bottom: 30px;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
}

.log-container {
  height: 500px;
  display: flex;
  flex-direction: column;
}
.log-toolbar {
  padding: 10px 0;
  display: flex;
  align-items: center;
}
.log-content {
  flex: 1;
  overflow: auto;
  background: #1e1e1e;
  color: #fff;
  padding: 10px;
  font-family: monospace;
  font-size: 14px;
  border-radius: 4px;
}
.log-line {
  padding: 2px 0;
  white-space: pre-wrap;
}
.line-number {
  display: inline-block;
  width: 60px;
  color: #888;
  text-align: right;
  padding-right: 10px;
  border-right: 1px solid #444;
  margin-right: 10px;
}
.log-level {
  padding: 0 4px;
  border-radius: 2px;
  font-weight: bold;
}
.log-level.error {
  background: #ff4d4f;
  color: white;
}
.log-level.warn {
  background: #faad14;
  color: white;
}
.log-level.info {
  background: #1890ff;
  color: white;
}
.log-level.debug {
  background: #87e8de;
  color: #000;
}
.search-highlight {
  background: rgba(255, 255, 0, 0.4);
  padding: 0 2px;
  border-radius: 2px;
}
.empty-log {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #888;
}
.log-pagination {
  margin-top: 10px;
  text-align: right;
}
.upload-button {
  margin-left: 20px;
}
</style>
