<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>个人资料</h2>
        <p>查看和编辑您的个人信息</p>
      </div>
    </div>

    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="8" :lg="6" :xl="6">
        <el-card class="profile-card">
          <div class="profile-avatar">
            <el-avatar :size="100" :src="userInfo.avatar">
              {{ userInfo.name ? userInfo.name.charAt(0).toUpperCase() : 'U' }}
            </el-avatar>

            <div class="avatar-actions">
              <el-upload
                  class="avatar-uploader"
                  action="#"
                  :auto-upload="false"
                  :on-change="handleAvatarChange"
                  :show-file-list="false"
              >
                <el-button size="small" type="primary">
                  <el-icon><Upload /></el-icon>
                  更换头像
                </el-button>
              </el-upload>
            </div>
          </div>

          <div class="profile-info">
            <h3>{{ userInfo.name }}</h3>
            <p class="user-role">{{ userInfo.role

                ```vue file="src/views/Profile.vue"
<template>
              <div class="app-container">
                <div class="page-header">
                  <div class="header-title">
                    <h2>个人资料</h2>
                    <p>查看和编辑您的个人信息</p>
                  </div>
                </div>

                <el-row :gutter="20">
                  <el-col :xs="24" :sm="24" :md="8" :lg="6" :xl="6">
                    <el-card class="profile-card">
                      <div class="profile-avatar">
                        <el-avatar :size="100" :src="userInfo.avatar">
                          {{ userInfo.name ? userInfo.name.charAt(0).toUpperCase() : 'U' }}
                        </el-avatar>

                        <div class="avatar-actions">
                          <el-upload
                              class="avatar-uploader"
                              action="#"
                              :auto-upload="false"
                              :on-change="handleAvatarChange"
                              :show-file-list="false"
                          >
                            <el-button size="small" type="primary">
                              <el-icon><Upload /></el-icon>
                              更换头像
                            </el-button>
                          </el-upload>
                        </div>
                      </div>

                      <div class="profile-info">
                        <h3>{{ userInfo.name }}</h3>
                        <p class="user-role">{{ userInfo.role }}</p>

                        <div class="user-stats">
                          <div class="stat-item">
                            <span class="stat-value">{{ userInfo.pipelines_count }}</span>
                            <span class="stat-label">流水线</span>
                          </div>

                          <div class="stat-item">
                            <span class="stat-value">{{ userInfo.deployments_count }}</span>
                            <span class="stat-label">部署</span>
                          </div>

                          <div class="stat-item">
                            <span class="stat-value">{{ userInfo.artifacts_count }}</span>
                            <span class="stat-label">制品</span>
                          </div>
                        </div>

                        <div class="user-contact">
                          <p>
                            <el-icon><Message /></el-icon>
                            {{ userInfo.email }}
                          </p>

                          <p v-if="userInfo.phone">
                            <el-icon><Phone /></el-icon>
                            {{ userInfo.phone }}
                          </p>
                        </div>

                        <div class="last-login">
                          <p>上次登录: {{ formatDate(userInfo.last_login) }}</p>
                          <p>注册时间: {{ formatDate(userInfo.created_at) }}</p>
                        </div>
                      </div>
                    </el-card>
                  </el-col>

                  <el-col :xs="24" :sm="24" :md="16" :lg="18" :xl="18">
                    <el-card>
                      <el-tabs v-model="activeTab">
                        <el-tab-pane label="基本信息" name="basic">
                          <el-form
                              ref="basicFormRef"
                              :model="basicForm"
                              :rules="basicRules"
                              label-position="top"
                          >
                            <el-row :gutter="20">
                              <el-col :span="12">
                                <el-form-item label="用户名" prop="username">
                                  <el-input v-model="basicForm.username" disabled />
                                </el-form-item>
                              </el-col>

                              <el-col :span="12">
                                <el-form-item label="姓名" prop="name">
                                  <el-input v-model="basicForm.name" placeholder="请输入姓名" />
                                </el-form-item>
                              </el-col>
                            </el-row>

                            <el-row :gutter="20">
                              <el-col :span="12">
                                <el-form-item label="邮箱" prop="email">
                                  <el-input v-model="basicForm.email" placeholder="请输入邮箱" />
                                </el-form-item>
                              </el-col>

                              <el-col :span="12">
                                <el-form-item label="电话" prop="phone">
                                  <el-input v-model="basicForm.phone" placeholder="请输入电话" />
                                </el-form-item>
                              </el-col>
                            </el-row>

                            <el-form-item label="个人简介" prop="bio">
                              <el-input
                                  v-model="basicForm.bio"
                                  type="textarea"
                                  :rows="3"
                                  placeholder="请输入个人简介"
                              />
                            </el-form-item>

                            <el-form-item>
                              <el-button type="primary" @click="saveBasicInfo" :loading="savingBasic">
                                保存信息
                              </el-button>
                            </el-form-item>
                          </el-form>
                        </el-tab-pane>

                        <el-tab-pane label="修改密码" name="password">
                          <el-form
                              ref="passwordFormRef"
                              :model="passwordForm"
                              :rules="passwordRules"
                              label-position="top"
                          >
                            <el-form-item label="当前密码" prop="current_password">
                              <el-input
                                  v-model="passwordForm.current_password"
                                  type="password"
                                  placeholder="请输入当前密码"
                                  show-password
                              />
                            </el-form-item>

                            <el-form-item label="新密码" prop="new_password">
                              <el-input
                                  v-model="passwordForm.new_password"
                                  type="password"
                                  placeholder="请输入新密码"
                                  show-password
                              />
                            </el-form-item>

                            <el-form-item label="确认新密码" prop="confirm_password">
                              <el-input
                                  v-model="passwordForm.confirm_password"
                                  type="password"
                                  placeholder="请再次输入新密码"
                                  show-password
                              />
                            </el-form-item>

                            <el-form-item>
                              <el-button type="primary" @click="changePassword" :loading="changingPassword">
                                修改密码
                              </el-button>
                            </el-form-item>
                          </el-form>
                        </el-tab-pane>

                        <el-tab-pane label="通知设置" name="notifications">
                          <el-form
                              ref="notificationFormRef"
                              :model="notificationForm"
                              label-position="top"
                          >
                            <h3>邮件通知</h3>

                            <el-form-item label="启用邮件通知">
                              <el-switch v-model="notificationForm.email_enabled" />
                            </el-form-item>

                            <el-form-item label="通知事件">
                              <el-checkbox-group v-model="notificationForm.email_events">
                                <el-checkbox label="pipeline_success">流水线成功</el-checkbox>
                                <el-checkbox label="pipeline_failure">流水线失败</el-checkbox>
                                <el-checkbox label="deployment_success">部署成功</el-checkbox>
                                <el-checkbox label="deployment_failure">部署失败</el-checkbox>
                                <el-checkbox label="system_notification">系统通知</el-checkbox>
                              </el-checkbox-group>
                            </el-form-item>

                            <h3>站内通知</h3>

                            <el-form-item label="启用站内通知">
                              <el-switch v-model="notificationForm.web_enabled" />
                            </el-form-item>

                            <el-form-item label="通知事件">
                              <el-checkbox-group v-model="notificationForm.web_events">
                                <el-checkbox label="pipeline_success">流水线成功</el-checkbox>
                                <el-checkbox label="pipeline_failure">流水线失败</el-checkbox>
                                <el-checkbox label="deployment_success">部署成功</el-checkbox>
                                <el-checkbox label="deployment_failure">部署失败</el-checkbox>
                                <el-checkbox label="system_notification">系统通知</el-checkbox>
                                <el-checkbox label="mention">@提及</el-checkbox>
                              </el-checkbox-group>
                            </el-form-item>

                            <el-form-item>
                              <el-button type="primary" @click="saveNotificationSettings" :loading="savingNotifications">
                                保存设置
                              </el-button>
                            </el-form-item>
                          </el-form>
                        </el-tab-pane>

                        <el-tab-pane label="API 密钥" name="api_keys">
                          <div class="api-keys-header">
                            <h3>API 密钥</h3>
                            <el-button type="primary" @click="createApiKey">
                              <el-icon><Plus /></el-icon>
                              创建 API 密钥
                            </el-button>
                          </div>

                          <el-table :data="apiKeys" style="width: 100%">
                            <el-table-column prop="name" label="名称" min-width="150" />

                            <el-table-column prop="key" label="密钥" min-width="200">
                              <template #default="{ row }">
                                <span v-if="row.is_new">{{ row.key }}</span>
                                <span v-else>••••••••••••••••••••••</span>
                              </template>
                            </el-table-column>

                            <el-table-column prop="created_at" label="创建时间" width="180">
                              <template #default="{ row }">
                                {{ formatDate(row.created_at) }}
                              </template>
                            </el-table-column>

                            <el-table-column prop="last_used_at" label="最后使用" width="180">
                              <template #default="{ row }">
                                {{ row.last_used_at ? formatDate(row.last_used_at) : '从未使用' }}
                              </template>
                            </el-table-column>

                            <el-table-column label="操作" width="120" fixed="right">
                              <template #default="{ row }">
                                <el-button
                                    link
                                    type="danger"
                                    size="small"
                                    @click="deleteApiKey(row.id)"
                                >
                                  删除
                                </el-button>
                              </template>
                            </el-table-column>
                          </el-table>

                          <div class="api-keys-note">
                            <p>注意：API 密钥只会在创建时显示一次，请妥善保存。</p>
                          </div>
                        </el-tab-pane>

                        <el-tab-pane label="活动日志" name="activity">
                          <el-timeline>
                            <el-timeline-item
                                v-for="(activity, index) in activities"
                                :key="index"
                                :timestamp="formatDate(activity.timestamp)"
                                :type="getActivityType(activity.type)"
                            >
                              {{ activity.description }}
                            </el-timeline-item>
                          </el-timeline>

                          <div class="load-more">
                            <el-button @click="loadMoreActivities" :loading="loadingActivities">
                              加载更多
                            </el-button>
                          </div>
                        </el-tab-pane>
                      </el-tabs>
                    </el-card>
                  </el-col>
                </el-row>

                &lt;!-- 创建 API 密钥对话框 -->
                <el-dialog
                    v-model="apiKeyDialogVisible"
                    title="创建 API 密钥"
                    width="50%"
                >
                  <el-form
                      ref="apiKeyFormRef"
                      :model="apiKeyForm"
                      :rules="apiKeyRules"
                      label-position="top"
                  >
                    <el-form-item label="密钥名称" prop="name">
                      <el-input v-model="apiKeyForm.name" placeholder="请输入密钥名称" />
                    </el-form-item>

                    <el-form-item label="过期时间" prop="expires_at">
                      <el-radio-group v-model="apiKeyForm.expiration_type">
                        <el-radio label="never">永不过期</el-radio>
                        <el-radio label="custom">自定义</el-radio>
                      </el-radio-group>

                      <el-date-picker
                          v-if="apiKeyForm.expiration_type === 'custom'"
                          v-model="apiKeyForm.expires_at"
                          type="datetime"
                          placeholder="选择过期时间"
                          style="margin-top: 10px; width: 100%;"
                      />
                    </el-form-item>
                  </el-form>

                  <template #footer>
        <span class="dialog-footer">
          <el-button @click="apiKeyDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitApiKeyForm" :loading="creatingApiKey">
            创建
          </el-button>
        </span>
                  </template>
                </el-dialog>
              </div>
            </template>

              <script setup>
                import { ref, reactive, onMounted } from 'vue';
                import { ElMessage, ElMessageBox } from 'element-plus';
                import { Upload, Message, Phone, Plus } from '@element-plus/icons-vue';
                import dayjs from 'dayjs';

                const activeTab = ref('basic');
                const savingBasic = ref(false);
                const changingPassword = ref(false);
                const savingNotifications = ref(false);
                const apiKeyDialogVisible = ref(false);
                const creatingApiKey = ref(false);
                const loadingActivities = ref(false);
                const basicFormRef = ref(null);
                const passwordFormRef = ref(null);
                const notificationFormRef = ref(null);
                const apiKeyFormRef = ref(null);

                // 用户信息
                const userInfo = reactive({
                  id: 1,
                  username: 'admin',
                  name: '管理员',
                  email: 'admin@example.com',
                  phone: '13800138000',
                  role: '系统管理员',
                  avatar: '',
                  bio: '系统管理员，负责平台的维护和管理。',
                  pipelines_count: 12,
                  deployments_count: 45,
                  artifacts_count: 28,
                  last_login: '2023-05-18T08:30:00',
                  created_at: '2023-01-01T00:00:00'
                });

                // 基本信息表单
                const basicForm = reactive({
                  username: userInfo.username,
                  name: userInfo.name,
                  email: userInfo.email,
                  phone: userInfo.phone,
                  bio: userInfo.bio
                });

                // 密码表单
                const passwordForm = reactive({
                  current_password: '',
                  new_password: '',
                  confirm_password: ''
                });

                // 通知设置表单
                const notificationForm = reactive({
                  email_enabled: true,
                  email_events: ['pipeline_failure', 'deployment_failure', 'system_notification'],
                  web_enabled: true,
                  web_events: ['pipeline_failure', 'deployment_failure', 'system_notification', 'mention']
                });

                // API 密钥表单
                const apiKeyForm = reactive({
                  name: '',
                  expiration_type: 'never',
                  expires_at: null
                });

                // API 密钥列表
                const apiKeys = ref([
                  {
                    id: 1,
                    name: '开发环境',
                    key: '••••••••••••••••••••••',
                    created_at: '2023-04-10T10:00:00',
                    last_used_at: '2023-05-15T14:30:00',
                    is_new: false
                  },
                  {
                    id: 2,
                    name: 'CI 服务器',
                    key: '••••••••••••••••••••••',
                    created_at: '2023-03-20T09:15:00',
                    last_used_at: '2023-05-17T16:45:00',
                    is_new: false
                  }
                ]);

                // 活动日志
                const activities = ref([
                  {
                    type: 'info',
                    description: '登录系统',
                    timestamp: '2023-05-18T08:30:00'
                  },
                  {
                    type: 'success',
                    description: '创建流水线 "前端构建"',
                    timestamp: '2023-05-17T14:20:00'
                  },
                  {
                    type: 'success',
                    description: '触发流水线 "API服务构建" 运行',
                    timestamp: '2023-05-17T11:05:00'
                  },
                  {
                    type: 'warning',
                    description: '修改个人资料',
                    timestamp: '2023-05-16T16:30:00'
                  },
                  {
                    type: 'danger',
                    description: '删除流水线 "测试项目"',
                    timestamp: '2023-05-15T09:45:00'
                  }
                ]);

                // 表单验证规则
                const basicRules = {
                  name: [
                    { required: true, message: '请输入姓名', trigger: 'blur' },
                    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
                  ],
                  email: [
                    { required: true, message: '请输入邮箱', trigger: 'blur' },
                    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
                  ]
                };

                const passwordRules = {
                  current_password: [
                    { required: true, message: '请输入当前密码', trigger: 'blur' }
                  ],
                  new_password: [
                    { required: true, message: '请输入新密码', trigger: 'blur' },
                    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
                  ],
                  confirm_password: [
                    { required: true, message: '请再次输入新密码', trigger: 'blur' },
                    {
                      validator: (rule, value, callback) => {
                        if (value !== passwordForm.new_password) {
                          callback(new Error('两次输入的密码不一致'));
                        } else {
                          callback();
                        }
                      },
                      trigger: 'blur'
                    }
                  ]
                };

                const apiKeyRules = {
                  name: [
                    { required: true, message: '请输入密钥名称', trigger: 'blur' },
                    { min: 2, max: 50, message: '长度在 2 到 50 个字符', trigger: 'blur' }
                  ]
                };

                // 处理头像变更
                const handleAvatarChange = async (file) => {
                  // 实际项目中应该调用API上传头像
                  try {
                    // 模拟上传
                    await new Promise(resolve => setTimeout(resolve, 1000));

                    // 使用 FileReader 读取文件并预览
                    const reader = new FileReader();
                    reader.onload = (e) => {
                      userInfo.avatar = e.target.result;
                    };
                    reader.readAsDataURL(file.raw);

                    ElMessage.success('头像上传成功');
                  } catch (error) {
                    console.error('头像上传失败:', error);
                    ElMessage.error('头像上传失败');
                  }
                };

                // 保存基本信息
                const saveBasicInfo = async () => {
                  if (!basicFormRef.value) return;

                  await basicFormRef.value.validate(async (valid) => {
                    if (valid) {
                      savingBasic.value = true;

                      try {
                        // 实际项目中应该调用API保存基本信息
                        await new Promise(resolve => setTimeout(resolve, 1000));

                        // 更新用户信息
                        Object.assign(userInfo, {
                          name: basicForm.name,
                          email: basicForm.email,
                          phone: basicForm.phone,
                          bio: basicForm.bio
                        });

                        ElMessage.success('基本信息已保存');
                      } catch (error) {
                        console.error('保存基本信息失败:', error);
                        ElMessage.error('保存基本信息失败');
                      } finally {
                        savingBasic.value = false;
                      }
                    } else {
                      ElMessage.warning('请填写必填项');
                      return false;
                    }
                  });
                };

                // 修改密码
                const changePassword = async () => {
                  if (!passwordFormRef.value) return;

                  await passwordFormRef.value.validate(async (valid) => {
                    if (valid) {
                      changingPassword.value = true;

                      try {
                        // 实际项目中应该调用API修改密码
                        await new Promise(resolve => setTimeout(resolve, 1000));

                        // 重置表单
                        passwordForm.current_password = '';
                        passwordForm.new_password = '';
                        passwordForm.confirm_password = '';

                        ElMessage.success('密码修改成功');
                      } catch (error) {
                        console.error('修改密码失败:', error);
                        ElMessage.error('修改密码失败');
                      } finally {
                        changingPassword.value = false;
                      }
                    } else {
                      ElMessage.warning('请填写必填项');
                      return false;
                    }
                  });
                };

                // 保存通知设置
                const saveNotificationSettings = async () => {
                  savingNotifications.value = true;

                  try {
                    // 实际项目中应该调用API保存通知设置
                    await new Promise(resolve => setTimeout(resolve, 1000));
                    ElMessage.success('通知设置已保存');
                  } catch (error) {
                    console.error('保存通知设置失败:', error);
                    ElMessage.error('保存通知设置失败');
                  } finally {
                    savingNotifications.value = false;
                  }
                };

                // 创建 API 密钥
                const createApiKey = () => {
                  // 重置表单
                  Object.assign(apiKeyForm, {
                    name: '',
                    expiration_type: 'never',
                    expires_at: null
                  });

                  apiKeyDialogVisible.value = true;
                };

                // 提交 API 密钥表单
                const submitApiKeyForm = async () => {
                  if (!apiKeyFormRef.value) return;

                  await apiKeyFormRef.value.validate(async (valid) => {
                    if (valid) {
                      creatingApiKey.value = true;

                      try {
                        // 实际项目中应该调用API创建 API 密钥
                        await new Promise(resolve => setTimeout(resolve, 1000));

                        // 生成随机密钥
                        const randomKey = Array.from(Array(32), () => Math.floor(Math.random() * 36).toString(36)).join('');

                        // 添加到列表
                        const newApiKey = {
                          id: Math.max(...apiKeys.value.map(key => key.id)) + 1,
                          name: apiKeyForm.name,
                          key: randomKey,
                          created_at: new Date().toISOString(),
                          last_used_at: null,
                          is_new: true
                        };

                        apiKeys.value.unshift(newApiKey);
                        apiKeyDialogVisible.value = false;

                        ElMessage.success('API 密钥创建成功，请妥善保存');

                        // 5秒后隐藏密钥
                        setTimeout(() => {
                          const index = apiKeys.value.findIndex(key => key.id === newApiKey.id);
                          if (index !== -1) {
                            apiKeys.value[index].is_new = false;
                          }
                        }, 5000);
                      } catch (error) {
                        console.error('创建 API 密钥失败:', error);
                        ElMessage.error('创建 API 密钥失败');
                      } finally {
                        creatingApiKey.value = false;
                      }
                    } else {
                      ElMessage.warning('请填写必填项');
                      return false;
                    }
                  });
                };

                // 删除 API 密钥
                const deleteApiKey = async (id) => {
                  try {
                    await ElMessageBox.confirm('确定要删除此 API 密钥吗？此操作不可恢复。', '删除确认', {
                      confirmButtonText: '确定',
                      cancelButtonText: '取消',
                      type: 'warning'
                    });

                    // 实际项目中应该调用API删除 API 密钥
                    apiKeys.value = apiKeys.value.filter(key => key.id !== id);
                    ElMessage.success('API 密钥已删除');
                  } catch (error) {
                    if (error !== 'cancel') {
                      console.error('删除 API 密钥失败:', error);
                      ElMessage.error('删除 API 密钥失败');
                    }
                  }
                };

                // 加载更多活动
                const loadMoreActivities = async () => {
                  loadingActivities.value = true;

                  try {
                    // 实际项目中应该调用API加载更多活动
                    await new Promise(resolve => setTimeout(resolve, 1000));

                    // 模拟加载更多
                    const moreActivities = [
                      {
                        type: 'success',
                        description: '部署应用到生产环境',
                        timestamp: '2023-05-14T15:20:00'
                      },
                      {
                        type: 'info',
                        description: '创建制品 "app-v1.0.0.zip"',
                        timestamp: '2023-05-14T14:10:00'
                      },
                      {
                        type: 'warning',
                        description: '修改流水线 "主应用构建"',
                        timestamp: '2023-05-13T11:30:00'
                      }
                    ];

                    activities.value.push(...moreActivities);
                  } catch (error) {
                    console.error('加载更多活动失败:', error);
                    ElMessage.error('加载更多活动失败');
                  } finally {
                    loadingActivities.value = false;
                  }
                };

                // 获取活动类型
                const getActivityType = (type) => {
                  switch (type) {
                    case 'success': return 'success';
                    case 'warning': return 'warning';
                    case 'danger': return 'danger';
                    default: return 'primary';
                  }
                };

                // 格式化日期
                const formatDate = (date) => {
                  if (!date) return '-';
                  return dayjs(date).format('YYYY-MM-DD HH:mm:ss');
                };

                onMounted(() => {
                  // 实际项目中应该从API获取用户信息
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

                .profile-card {
                  margin-bottom: 20px;
                }

                .profile-avatar {
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                  margin-bottom: 20px;
                }

                .avatar-actions {
                  margin-top: 10px;
                }

                .profile-info {
                  text-align: center;
                }

                .profile-info h3 {
                  margin: 0 0 5px 0;
                  font-size: 18px;
                  font-weight: 600;
                }

                .user-role {
                  color: #909399;
                  margin-bottom: 15px;
                }

                .user-stats {
                  display: flex;
                  justify-content: space-around;
                  margin-bottom: 20px;
                }

                .stat-item {
                  display: flex;
                  flex-direction: column;
                  align-items: center;
                }

                .stat-value {
                  font-size: 18px;
                  font-weight: 600;
                  color: #409EFF;
                }

                .stat-label {
                  font-size: 12px;
                  color: #909399;
                }

                .user-contact {
                  text-align: left;
                  margin-bottom: 15px;
                }

                .user-contact p {
                  display: flex;
                  align-items: center;
                  gap: 5px;
                  margin: 5px 0;
                  color: #606266;
                }

                .last-login {
                  font-size: 12px;
                  color: #909399;
                  text-align: left;
                }

                .last-login p {
                  margin: 5px 0;
                }

                h3 {
                  margin-top: 20px;
                  margin-bottom: 15px;
                  font-size: 16px;
                  font-weight: 600;
                }

                .api-keys-header {
                  display: flex;
                  justify-content: space-between;
                  align-items: center;
                  margin-bottom: 20px;
                }

                .api-keys-header h3 {
                  margin: 0;
                }

                .api-keys-note {
                  margin-top: 15px;
                  font-size: 12px;
                  color: #909399;
                }

                .load-more {
                  margin-top: 20px;
                  display: flex;
                  justify-content: center;
                }
              </style>
