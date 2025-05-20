<template>
  <div class="navbar">
    <div class="left-menu">
      <div class="hamburger-container" @click="toggleSidebar">
        <el-icon :size="20">
          <Fold v-if="!isMobile" />
          <Expand v-else />
        </el-icon>
      </div>

      <breadcrumb class="breadcrumb-container" />
    </div>

    <div class="right-menu">
      <el-tooltip content="全屏" placement="bottom">
        <div class="right-menu-item" @click="toggleFullScreen">
          <el-icon :size="18"><FullScreen /></el-icon>
        </div>
      </el-tooltip>

      <el-dropdown trigger="click" class="avatar-container">
        <div class="avatar-wrapper">
          <el-avatar :size="30" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
          <span class="user-name">管理员</span>
          <el-icon><CaretBottom /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item>
              <router-link to="/profile">个人资料</router-link>
            </el-dropdown-item>
            <el-dropdown-item divided @click="logout">
              <span style="display: block">退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useRouter } from 'vue-router';
import { Fold, Expand, FullScreen, CaretBottom } from '@element-plus/icons-vue';
import Breadcrumb from './Breadcrumb.vue';
import { useAuthStore } from '@/stores/auth';

const emit = defineEmits(['toggle-sidebar']);
const router = useRouter();
const authStore = useAuthStore();
const isMobile = ref(false);

onMounted(() => {
  // Check screen size on mount and update isMobile
  const checkScreenSize = () => {
    isMobile.value = window.innerWidth <= 768; // Adjust the breakpoint as needed
  };

  // Initial check
  checkScreenSize();

  // Listen for window resize events
  window.addEventListener('resize', checkScreenSize);

  // Clean up the event listener on unmount
  onBeforeUnmount(() => {
    window.removeEventListener('resize', checkScreenSize);
  });
});

const toggleSidebar = () => {
  emit('toggle-sidebar');
};

const toggleFullScreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen();
  } else {
    if (document.exitFullscreen) {
      document.exitFullscreen();
    }
  }
};

const logout = async () => {
  await authStore.logout();
  router.push('/login');
};
</script>

<style scoped>
.navbar {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
}

.left-menu {
  display: flex;
  align-items: center;
}

.hamburger-container {
  padding: 0 15px;
  cursor: pointer;
  transition: background 0.3s;
  height: 100%;
  display: flex;
  align-items: center;
}

.hamburger-container:hover {
  background: rgba(0, 0, 0, 0.025);
}

.breadcrumb-container {
  margin-left: 8px;
}

.right-menu {
  display: flex;
  align-items: center;
}

.right-menu-item {
  padding: 0 12px;
  cursor: pointer;
  transition: background 0.3s;
  height: 100%;
  display: flex;
  align-items: center;
}

.right-menu-item:hover {
  background: rgba(0, 0, 0, 0.025);
}

.avatar-container {
  margin-left: 10px;
  cursor: pointer;
}

.avatar-wrapper {
  display: flex;
  align-items: center;
}

.user-name {
  margin: 0 5px;
}
</style>
