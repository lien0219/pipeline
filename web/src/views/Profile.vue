<template>
  <div class="profile-page">
    <div class="page-header">
      <h1>个人资料</h1>
      <p>管理您的账号信息</p>
    </div>

    <div class="profile-content">
      <div class="profile-card">
        <h2>基本信息</h2>
        <form @submit.prevent="updateProfile" class="profile-form">
          <div class="form-group">
            <label for="username" class="form-label">用户名</label>
            <input
                id="username"
                v-model="userInfo.username"
                type="text"
                class="form-input"
                disabled
            />
            <small class="form-text">用户名不可修改</small>
          </div>

          <div class="form-group">
            <label for="email" class="form-label">电子邮箱</label>
            <input
                id="email"
                v-model="userInfo.email"
                type="email"
                class="form-input"
                placeholder="请输入电子邮箱"
                required
            />
          </div>

          <div class="form-group">
            <label for="name" class="form-label">姓名</label>
            <input
                id="name"
                v-model="userInfo.name"
                type="text"
                class="form-input"
                placeholder="请输入姓名"
            />
          </div>

          <div v-if="updateProfileError" class="error-message">
            {{ updateProfileError }}
          </div>

          <div v-if="updateProfileSuccess" class="success-message">
            {{ updateProfileSuccess }}
          </div>

          <button type="submit" class="btn btn-primary" :disabled="updateProfileLoading">
            {{ updateProfileLoading ? '保存中...' : '保存修改' }}
          </button>
        </form>
      </div>

      <div class="profile-card">
        <h2>修改密码</h2>
        <form @submit.prevent="changePassword" class="profile-form">
          <div class="form-group">
            <label for="currentPassword" class="form-label">当前密码</label>
            <input
                id="currentPassword"
                v-model="passwordForm.oldPassword"
                type="password"
                class="form-input"
                placeholder="请输入当前密码"
                required
            />
          </div>

          <div class="form-group">
            <label for="newPassword" class="form-label">新密码</label>
            <input
                id="newPassword"
                v-model="passwordForm.newPassword"
                type="password"
                class="form-input"
                placeholder="请输入新密码"
                required
            />
            <small class="form-text">密码长度至少为6个字符</small>
          </div>

          <div class="form-group">
            <label for="confirmPassword" class="form-label">确认新密码</label>
            <input
                id="confirmPassword"
                v-model="passwordForm.confirmPassword"
                type="password"
                class="form-input"
                placeholder="请再次输入新密码"
                required
            />
          </div>

          <div v-if="changePasswordError" class="error-message">
            {{ changePasswordError }}
          </div>

          <div v-if="changePasswordSuccess" class="success-message">
            {{ changePasswordSuccess }}
          </div>

          <button type="submit" class="btn btn-primary" :disabled="changePasswordLoading || !isPasswordFormValid">
            {{ changePasswordLoading ? '修改中...' : '修改密码' }}
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { authApi } from '@/api/auth';

const authStore = useAuthStore();

// 用户信息表单
const userInfo = ref({
  username: '',
  email: '',
  name: ''
});

// 密码表单
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// 表单状态
const updateProfileLoading = ref(false);
const updateProfileError = ref('');
const updateProfileSuccess = ref('');

const changePasswordLoading = ref(false);
const changePasswordError = ref('');
const changePasswordSuccess = ref('');

// 密码表单验证
const isPasswordFormValid = computed(() => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    changePasswordError.value = '两次输入的密码不一致';
    return false;
  }

  if (passwordForm.value.newPassword.length < 6) {
    changePasswordError.value = '密码长度必须至少为6个字符';
    return false;
  }

  changePasswordError.value = '';
  return true;
});

const user = ref(null);

// 加载用户信息
onMounted(async () => {
  user.value = authStore.currentUser;

  if (user.value) {
    userInfo.value = {
      username: user.value.username || '',
      email: user.value.email || '',
      name: user.value.name || ''
    };
  } else {
    // 如果store中没有用户信息，尝试重新获取
    try {
      await authStore.fetchUserInfo();
      const updatedUser = authStore.currentUser;

      if (updatedUser) {
        userInfo.value = {
          username: updatedUser.username || '',
          email: updatedUser.email || '',
          name: updatedUser.name || ''
        };
      }
    } catch (error) {
      console.error('Failed to fetch user info:', error);
    }
  }
});

// 更新个人资料
const updateProfile = async () => {
  updateProfileLoading.value = true;
  updateProfileError.value = '';
  updateProfileSuccess.value = '';

  try {
    // 只更新email和name
    const { email, name } = userInfo.value;
    await authStore.updateUserInfo({ email, name });

    updateProfileSuccess.value = '个人资料更新成功';
  } catch (error) {
    updateProfileError.value = error.response?.data?.message || '更新个人资料失败';
    console.error('Failed to update profile:', error);
  } finally {
    updateProfileLoading.value = false;
  }
};

// 修改密码
const changePassword = async () => {
  if (!isPasswordFormValid.value) return;

  changePasswordLoading.value = true;
  changePasswordError.value = '';
  changePasswordSuccess.value = '';

  try {
    const { oldPassword, newPassword } = passwordForm.value;
    await authStore.changePassword({ oldPassword, newPassword });

    // 清空表单
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    };

    changePasswordSuccess.value = '密码修改成功';
  } catch (error) {
    changePasswordError.value = error.response?.data?.message || '修改密码失败';
    console.error('Failed to change password:', error);
  } finally {
    changePasswordLoading.value = false;
  }
};
</script>

<style scoped>
.profile-page {
  width: 100%;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.page-header p {
  color: var(--text-light);
}

.profile-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
}

.profile-card {
  background-color: var(--surface-color);
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.profile-card h2 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: var(--text-color);
}

.profile-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.form-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: var(--text-light);
}

.error-message {
  background-color: rgba(211, 47, 47, 0.1);
  color: var(--error-color);
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.success-message {
  background-color: rgba(56, 142, 60, 0.1);
  color: var(--success-color);
  padding: 0.75rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .profile-content {
    grid-template-columns: 1fr;
  }
}
</style>