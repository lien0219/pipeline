<template>
  <div class="default-layout">
    <header class="header">
      <div class="container header-container">
        <div class="logo">
          <router-link to="/">企业应用</router-link>
        </div>
        <nav class="nav">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link v-if="isAuthenticated" to="/dashboard" class="nav-link">仪表盘</router-link>
        </nav>
        <div class="user-menu">
          <template v-if="isAuthenticated">
            <div class="user-dropdown" @click="toggleDropdown" ref="dropdown">
              <div class="user-info">
                <span class="username">{{ user?.username }}</span>
                <i class="icon">▼</i>
              </div>
              <div class="dropdown-menu" v-show="dropdownOpen">
                <router-link to="/profile" class="dropdown-item">个人资料</router-link>
                <div class="dropdown-divider"></div>
                <button @click="handleLogout" class="dropdown-item">退出登录</button>
              </div>
            </div>
          </template>
          <template v-else>
            <router-link to="/auth/login" class="btn btn-secondary">登录</router-link>
            <router-link to="/auth/register" class="btn btn-primary">注册</router-link>
          </template>
        </div>
      </div>
    </header>

    <main class="main">
      <div class="container">
        <router-view />
      </div>
    </main>

    <footer class="footer">
      <div class="container">
        <p>&copy; {{ currentYear }} 企业应用. 保留所有权利.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const authStore = useAuthStore();
const isAuthenticated = computed(() => authStore.isAuthenticated);
const user = computed(() => authStore.currentUser);
const currentYear = computed(() => new Date().getFullYear());
const dropdown = ref(null);
const dropdownOpen = ref(false);



// 切换下拉菜单
const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value;
};

// 点击外部关闭下拉菜单
const handleClickOutside = (event) => {
  if (dropdown.value && !dropdown.value.contains(event.target)) {
    dropdownOpen.value = false;
  }
};

// 退出登录
const handleLogout = () => {
  authStore.logout();
  router.push('/auth/login');
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
.default-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  background-color: var(--surface-color);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 64px;
}

.logo a {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--primary-color);
  text-decoration: none;
}

.nav {
  display: flex;
  gap: 1.5rem;
}

.nav-link {
  color: var(--text-color);
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--primary-color);
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-dropdown {
  position: relative;
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
}

.username {
  font-weight: 500;
}

.icon {
  font-size: 0.75rem;
  color: var(--text-light);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  width: 200px;
  background-color: var(--surface-color);
  border-radius: 4px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  z-index: 10;
}

.dropdown-item {
  display: block;
  padding: 0.75rem 1rem;
  color: var(--text-color);
  text-decoration: none;
  transition: background-color 0.3s;
  text-align: left;
  width: 100%;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 1rem;
}

.dropdown-item:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.dropdown-divider {
  height: 1px;
  background-color: rgba(0, 0, 0, 0.1);
  margin: 0.25rem 0;
}

.main {
  flex: 1;
  padding: 2rem 0;
}

.footer {
  background-color: var(--surface-color);
  padding: 1.5rem 0;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
  text-align: center;
  color: var(--text-light);
}
</style>