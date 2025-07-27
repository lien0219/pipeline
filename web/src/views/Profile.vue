<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>个人资料</h2>
      </div>
    </div>

    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :md="8" :lg="6" :xl="6">
        <el-card class="profile-card">
          <div class="profile-avatar">
            <el-avatar :size="100" :src="userInfo.avatar">
              {{ userInfo.name ? userInfo.name.charAt(0).toUpperCase() : "U" }}
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
                <span class="stat-value">{{
                  userInfo.pipelines_count || 0
                }}</span>
                <span class="stat-label">流水线</span>
              </div>

              <div class="stat-item">
                <span class="stat-value">{{
                  userInfo.deployments_count || 0
                }}</span>
                <span class="stat-label">部署</span>
              </div>

              <div class="stat-item">
                <span class="stat-value">{{
                  userInfo.artifacts_count || 0
                }}</span>
                <span class="stat-label">制品</span>
              </div>
            </div>

            <div class="user-contact">
              <p>
                <el-icon><Message /></el-icon>
                {{ userInfo.email || "未设置" }}
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
                      <el-input
                        v-model="basicForm.name"
                        placeholder="请输入姓名"
                      />
                    </el-form-item>
                  </el-col>
                </el-row>

                <el-row :gutter="20">
                  <el-col :span="12">
                    <el-form-item label="邮箱" prop="email">
                      <el-input
                        v-model="basicForm.email"
                        placeholder="请输入邮箱"
                      />
                    </el-form-item>
                  </el-col>

                  <el-col :span="12">
                    <el-form-item label="电话" prop="phone">
                      <el-input
                        v-model="basicForm.phone"
                        placeholder="请输入电话"
                      />
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
                  <el-button
                    type="primary"
                    @click="saveBasicInfo"
                    :loading="savingBasic"
                  >
                    保存信息
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>

            <el-tab-pane label="密码修改" name="password">
              <el-form
                ref="passwordFormRef"
                :model="passwordForm"
                :rules="passwordRules"
                label-position="top"
              >
                <el-form-item label="当前密码" prop="currentPassword">
                  <el-input
                    v-model="passwordForm.currentPassword"
                    type="password"
                    placeholder="请输入当前密码"
                  />
                </el-form-item>

                <el-form-item label="新密码" prop="newPassword">
                  <el-input
                    v-model="passwordForm.newPassword"
                    type="password"
                    placeholder="请输入新密码"
                  />
                </el-form-item>

                <el-form-item label="确认新密码" prop="confirmPassword">
                  <el-input
                    v-model="passwordForm.confirmPassword"
                    type="password"
                    placeholder="请确认新密码"
                  />
                </el-form-item>

                <el-form-item>
                  <el-button
                    type="primary"
                    @click="changePassword"
                    :loading="savingPassword"
                  >
                    修改密码
                  </el-button>
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </el-tabs>
        </el-card>
      </el-col>
    </el-row>

    <!-- 创建 API 密钥对话框 -->
    <el-dialog v-model="apiKeyDialogVisible" title="创建 API 密钥" width="50%">
      <div class="api-key-dialog-content">
        <p>API 密钥创建功能尚未实现</p>
      </div>
      <template #footer>
        <el-button @click="apiKeyDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="apiKeyDialogVisible = false"
          >创建</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from "vue";
import { ElMessage, ElMessageBox, ElLoading } from "element-plus";
import {
  Upload,
  Message,
  Phone,
  Plus,
  Delete,
  Key,
  Bell,
  Clock,
} from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { authApi } from "@/api/auth";

// 状态管理
const userInfo = ref({});
const loading = ref(false);
const activeTab = ref("basic");
const apiKeyDialogVisible = ref(false);
const savingBasic = ref(false);
const savingPassword = ref(false);
const apiKeys = ref([]);
const notifications = ref([]);
const activities = ref([]);
const avatarFile = ref(null);

// 表单数据
const basicForm = reactive({
  username: "",
  name: "",
  email: "",
  phone: "",
  bio: "",
});

const passwordForm = reactive({
  currentPassword: "",
  newPassword: "",
  confirmPassword: "",
});

// 表单验证规则
const basicRules = ref({
  name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
  email: [
    { required: true, message: "请输入邮箱", trigger: "blur" },
    { type: "email", message: "请输入有效的邮箱地址", trigger: "blur" },
  ],
});

const passwordRules = ref({
  currentPassword: [
    { required: true, message: "请输入当前密码", trigger: "blur" },
  ],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    { min: 6, message: "密码长度不能少于6位", trigger: "blur" },
  ],
  confirmPassword: [
    { required: true, message: "请确认新密码", trigger: "blur" },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error("两次输入的密码不一致"));
        } else {
          callback();
        }
      },
      trigger: "blur",
    },
  ],
});

// 引用
const basicFormRef = ref(null);
const passwordFormRef = ref(null);

// 格式化日期
const formatDate = (date) => {
  if (!date) return "-";
  return dayjs(date).format("YYYY-MM-DD HH:mm:ss");
};

// 获取用户信息
const fetchUserInfo = async () => {
  try {
    loading.value = true;
    const response = await authApi.getUserInfo();
    userInfo.value = response.data;
    // 填充表单数据
    basicForm.username = userInfo.value.username || "";
    basicForm.name = userInfo.value.name || "";
    basicForm.email = userInfo.value.email || "";
    basicForm.phone = userInfo.value.phone || "";
    basicForm.bio = userInfo.value.bio || "";
  } catch (error) {
    ElMessage.error("获取用户信息失败: " + (error.message || "未知错误"));
  } finally {
    loading.value = false;
  }
};

// 保存基本信息
const saveBasicInfo = async () => {
  try {
    const form = basicFormRef.value;
    if (!form) return;

    await form.validate();
    savingBasic.value = true;

    const response = await authApi.updateUserInfo({
      name: basicForm.name,
      email: basicForm.email,
      phone: basicForm.phone,
      bio: basicForm.bio,
    });

    ElMessage.success("信息保存成功");
    // 更新本地用户信息
    userInfo.value.name = basicForm.name;
    userInfo.value.email = basicForm.email;
    userInfo.value.phone = basicForm.phone;
    userInfo.value.bio = basicForm.bio;
  } catch (error) {
    if (error.name !== "ValidationError") {
      ElMessage.error("保存信息失败: " + (error.message || "未知错误"));
    }
  } finally {
    savingBasic.value = false;
  }
};

// 更改密码
const changePassword = async () => {
  try {
    const form = passwordFormRef.value;
    if (!form) return;

    await form.validate();
    savingPassword.value = true;

    await authApi.changePassword({
      OldPassword: passwordForm.currentPassword,
      NewPassword: passwordForm.newPassword,
    });

    ElMessage.success("密码修改成功");
    passwordForm.currentPassword = "";
    passwordForm.newPassword = "";
    passwordForm.confirmPassword = "";
    form.resetFields();
  } catch (error) {
    if (error.name !== "ValidationError") {
      ElMessage.error("修改密码失败: " + (error.message || "未知错误"));
    }
  } finally {
    savingPassword.value = false;
  }
};

// 处理头像变更
const handleAvatarChange = (file) => {
  avatarFile.value = file.raw;
};

onMounted(() => {
  fetchUserInfo();
});
</script>
