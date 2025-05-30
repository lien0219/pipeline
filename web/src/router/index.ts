import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import AppLayout from "@/components/layout/AppLayout.vue";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: { requiresAuth: false, title: "登录" },
  },
  {
    path: "/",
    component: AppLayout,
    redirect: "/dashboard",
    children: [
      {
        path: "dashboard",
        name: "Dashboard",
        component: () => import("@/views/Dashboard.vue"),
        meta: { requiresAuth: true, title: "仪表盘" },
      },
      {
        path: "pipelines",
        name: "Pipelines",
        component: () => import("@/views/pipeline/PipelineList.vue"),
        meta: { requiresAuth: true, title: "流水线列表" },
      },
      {
        path: "pipelines/:id",
        name: "PipelineDetail",
        component: () => import("@/views/pipeline/PipelineDetail.vue"),
        meta: { requiresAuth: true, title: "流水线详情" },
      },
      {
        path: "pipelines/create",
        name: "PipelineCreate",
        component: () => import("@/views/pipeline/PipelineCreate.vue"),
        meta: { requiresAuth: true, title: "创建流水线" },
      },
      {
        path: "pipelines/:id/edit",
        name: "PipelineEdit",
        component: () => import("@/views/pipeline/PipelineCreate.vue"),
        meta: { requiresAuth: true, title: "编辑流水线" },
      },
      {
        path: "builds/history",
        name: "BuildHistory",
        component: () => import("@/views/build/BuildHistory.vue"),
        meta: { requiresAuth: true, title: "构建历史" },
      },
      {
        path: "builds/templates",
        name: "BuildTemplates",
        component: () => import("@/views/build/BuildTemplates.vue"),
        meta: { requiresAuth: true, title: "构建模板" },
      },
      {
        path: "deploy/environments",
        name: "Environments",
        component: () => import("@/views/deploy/Environments.vue"),
        meta: { requiresAuth: true, title: "环境管理" },
      },
      {
        path: "deploy/releases",
        name: "Releases",
        component: () => import("@/views/deploy/Releases.vue"),
        meta: { requiresAuth: true, title: "发布记录" },
      },
      {
        path: "artifacts",
        name: "Artifacts",
        component: () => import("@/views/artifact/ArtifactList.vue"),
        meta: { requiresAuth: true, title: "制品管理" },
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/settings/Settings.vue"),
        meta: { requiresAuth: true, title: "系统设置" },
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/Profile.vue"),
        meta: { requiresAuth: true, title: "个人资料" },
      },
    ],
  },
  {
    path: "/:pathMatch(.*)*",
    name: "NotFound",
    component: () => import("@/views/NotFound.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  const requiresAuth = to.matched.some(
    (record) => record.meta.requiresAuth !== false
  );

  if (requiresAuth && !authStore.isAuthenticated) {
    next("/login");
  } else {
    next();
  }
});

export default router;
