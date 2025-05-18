<template>
  <div class="login-page">
    <h2 class="auth-title">登录</h2>
    <p class="auth-subtitle">欢迎回来！请登录您的账号</p>

    <form @submit.prevent="handleLogin" class="auth-form">
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
      </div>

      <div v-if="error" class="error-message">
        {{ error }}
      </div>

      <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
        {{ loading ? '登录中...' : '登录' }}
      </button>
    </form>

    <div class="auth-links">
      <p>还没有账号？ <router-link to="/auth/register">立即注册</router-link></p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const route = useRoute();
const authStore = useAuthStore();

const form = ref({
  username: '',
  password: ''
});

const loading = computed(() => authStore.loading);
const error = computed(() => authStore.error);

const handleLogin = async () => {
  try {
    await authStore.login(form.value);

    // 登录成功后重定向
    const redirectPath = route.query.redirect || '/dashboard';
    router.push(redirectPath);
  } catch (err) {
    // 错误已在store中处理
    console.error('Login failed:', err);
  }
};
</script>

<style scoped>
.login-page {
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