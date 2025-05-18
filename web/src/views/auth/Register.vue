<template>
  <div class="register-page">
    <h2 class="auth-title">注册</h2>
    <p class="auth-subtitle">创建您的账号</p>

    <form @submit.prevent="handleRegister" class="auth-form">
      <div class="form-group">
        <label for="username" class="form-label">用户名</label>
        <input
            id="username"
            v-model="form.username"
            type="text"
            class="form-input"
            placeholder="请输入用户名"
            required
        />
        <small class="form-text">用户名长度为3-32个字符</small>
      </div>

      <div class="form-group">
        <label for="email" class="form-label">电子邮箱</label>
        <input
            id="email"
            v-model="form.email"
            type="email"
            class="form-input"
            placeholder="请输入电子邮箱"
            required
        />
      </div>

      <div class="form-group">
        <label for="password" class="form-label">密码</label>
        <input
            id="password"
            v-model="form.password"
            type="password"
            class="form-input"
            placeholder="请输入密码"
            required
        />
        <small class="form-text">密码长度至少为6个字符</small>
      </div>

      <div class="form-group">
        <label for="confirmPassword" class="form-label">确认密码</label>
        <input
            id="confirmPassword"
            v-model="form.confirmPassword"
            type="password"
            class="form-input"
            placeholder="请再次输入密码"
            required
        />
      </div>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <button type="submit" class="btn btn-primary btn-block" :disabled="loading || !isFormValid">
        {{ loading ? '注册中...' : '注册' }}
      </button>
    </form>

    <div class="auth-links">
      <p>已有账号？ <router-link to="/auth/login">立即登录</router-link></p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
});

const localError = ref('');

const isFormValid = computed(() => {
  let isValid = true;
  localError.value = '';

  if (form.value.password !== form.value.confirmPassword) {
    localError.value = '两次输入的密码不一致';
    isValid = false;
  }

  if (form.value.username.length < 3 || form.value.username.length > 32) {
    localError.value = '用户名长度必须在3-32个字符之间';
    isValid = false;
  }

  if (form.value.password.length < 6) {
    localError.value = '密码长度必须至少为6个字符';
    isValid = false;
  }

  return isValid;
});

const loading = computed(() => authStore.loading);
const error = computed(() => localError.value || authStore.error);

const handleRegister = async () => {
  if (!isFormValid.value) return;

  try {
    // 注册请求不需要确认密码字段
    const { username, email, password } = form.value;
    await authStore.register({ username, email, password });

    // 注册成功后跳转到登录页
    router.push({
      path: '/auth/login',
      query: { registered: 'true' }
    });
  } catch (err) {
    // 错误已在store中处理
    console.error('Registration failed:', err);
  }
};
</script>

<style scoped>
.register-page {
  width: 100%;
}

.auth-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.auth-subtitle {
  color: var(--text-light);
  margin-bottom: 2rem;
}

.auth-form {
  margin-bottom: 1.5rem;
}

.form-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: var(--text-light);
}

.btn-block {
  width: 100%;
  margin-top: 1.5rem;
  padding: 0.75rem;
}

.error-message {
  background-color: rgba(211, 47, 47, 0.1);
  color: var(--error-color);
  padding: 0.75rem;
  border-radius: 4px;
  margin-top: 1rem;
  font-size: 0.875rem;
}

.auth-links {
  text-align: center;
  margin-top: 1.5rem;
  color: var(--text-light);
}

.auth-links a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.auth-links a:hover {
  text-decoration: underline;
}
</style>