<template>
  <div class="dashboard-layout">
    <el-container class="container">
      <!-- 侧边栏 -->
      <el-aside width="220px" class="aside">
        <div class="logo">
          <img src="/logo.svg" alt="Logo" height="40" />
          <span>CICD 平台</span>
        </div>

        <el-scrollbar>
          <el-menu
              :default-active="activeMenu"
              class="sidebar-menu"
              :collapse="isCollapse"
              background-color="#304156"
              text-color="#bfcbd9"
              active-text-color="#409EFF"
              router
          >
            <el-menu-item index="/dashboard">
              <el-icon><el-icon-monitor /></el-icon>
              <span>概览</span>
            </el-menu-item>

            <el-sub-menu index="/pipelines">
              <template #title>
                <el-icon><el-icon-connection /></el-icon>
                <span>流水线管理</span>
              </template>
              <el-menu-item index="/pipelines/list">流水线列表</el-menu-item>
              <el-menu-item index="/pipelines/create">创建流水线</el-menu-item>
              <el-menu-item index="/pipelines/templates">流水线模板</el-menu-item>
            </el-sub-menu>

            <el-sub-menu index="/jobs">
              <template #title>
                <el-icon><el-icon-tickets /></el-icon>
                <span>任务管理</span>
              </template>
              <el-menu-item index="/jobs/list">任务列表</el-menu-item>
              <el-menu-item index="/jobs/history">执行历史</el-menu-item>
            </el-sub-menu>

            <el-sub-menu index="/environments">
              <template #title>
                <el-icon><el-icon-box /></el-icon>
                <span>环境管理</span>
              </template>
              <el-menu-item index="/environments/list">环境列表</el-menu-item>
              <el-menu-item index="/environments/variables">环境变量</el-menu-item>
            </el-sub-menu>

            <el-menu-item index="/artifacts">
              <el-icon><el-icon-files /></el-icon>
              <span>制品管理</span>
            </el-menu-item>

            <el-menu-item index="/settings">
              <el-icon><el-icon-setting /></el-icon>
              <span>系统设置</span>
            </el-menu-item>
          </el-menu>
        </el-scrollbar>
      </el-aside>

      <el-container class="main-container">
        <!-- 头部 -->
        <el-header class="header">
          <div class="header-left">
            <el-button @click="toggleSidebar" text>
              <el-icon :size="20">
                <el-icon-fold v-if="!isCollapse" />
                <el-icon-expand v-else />
              </el-icon>
            </el-button>

            <el-breadcrumb separator="/">
              <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path">
                {{ item.title }}
              </el-breadcrumb-item>
            </el-breadcrumb>
          </div>

          <div class="header-right">
            <el-tooltip content="查看文档" placement="bottom">
              <el-button :icon="ElIconQuestion" circle />
            </el-tooltip>

            <el-dropdown trigger="click">
              <span class="user-dropdown">
                <el-avatar :size="32" src="https://placeholder.pics/svg/32x32" />
                <span class="username">{{ username }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="goToProfile">个人资料</el-dropdown-item>
                  <el-dropdown-item divided @click="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <!-- 主内容区 -->
        <el-main class="main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import {
  Monitor as ElIconMonitor,
  Connection as ElIconConnection,
  Tickets as ElIconTickets,
  Box as ElIconBox,
  Files as ElIconFiles,
  Setting as ElIconSetting,
  Fold as ElIconFold,
  Expand as ElIconExpand,
  Question as ElIconQuestion
} from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const isCollapse = ref(false);
const username = ref('');

// 激活的菜单
const activeMenu = computed(() => {
  return route.path;
});

// 面包屑导航
const breadcrumbs = computed(() => {
  const { meta, path, matched } = route;

  if (!matched) return [];

  return matched.filter(item => item.meta && item.meta.title).map(item => ({
    title: item.meta.title,
    path: item.path
  }));
});

// 折叠/展开侧边栏
const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value;
};

// 跳转到个人资料页面
const goToProfile = () => {
  router.push('/profile');
};

// 退出登录
const logout = async () => {
  try {
    await authStore.logout();
    router.push('/login');
    ElMessage.success('退出登录成功');
  } catch (error) {
    console.error('Failed to logout:', error);
    ElMessage.error('退出登录失败');
  }
};

onMounted(async () => {
  // 获取用户信息
  try {
    await authStore.fetchUserInfo();
  } catch (error) {
    console.error('Failed to fetch user info:', error);
  }
  const user = authStore.currentUser;

  if (user) {
    username.value = user.name || user.username;
  }
});
</script>

<style scoped>
.dashboard-layout {
  height: 100vh;
  width: 100%;
}

.container {
  height: 100%;
}

.aside {
  height: 100%;
  background-color: #304156;
  color: #bfcbd9;
  transition: width 0.3s;
  overflow: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 20px;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  border-bottom: 1px solid #1f2d3d;
}

.logo img {
  margin-right: 10px;
}

.sidebar-menu {
  border-right: none;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0;
}

.header-left, .header-right {
  display: flex;
  align-items: center;
}

.header-left {
  gap: 15px;
}

.header-right {
  gap: 10px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 8px;
  font-size: 14px;
}

.main-container {
  flex-direction: column;
  height: 100%;
}

.main {
  padding: 20px;
  background-color: #f5f7fa;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .aside {
    width: 64px !important;
  }

  .sidebar-menu {
    width: 64px;
  }

  .username {
    display: none;
  }
}
</style>