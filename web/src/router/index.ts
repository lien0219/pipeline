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
        path: "pipelines/:id/runs/:runId",
        name: "PipelineRunDetail",
        component: () => import("@/views/pipeline/PipelineRunDetail.vue"),
        meta: { requiresAuth: true, title: "流水线运行详情" },
      },
      {
        path: "pipeline/create",
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
        path: "pipelines/:id/designer",
        name: "PipelineDesigner",
        component: () => import("@/views/pipeline/PipelineDesigner.vue"),
        meta: { requiresAuth: true, title: "流水线设计器" },
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
        path: "/deploy/releases/:id",
        name: "ReleaseDetail",
        component: () => import("@/views/deploy/ReleaseDetail.vue"),
        meta: { requiresAuth: true, title: "发布详情" },
      },
      {
        path: "artifacts",
        name: "Artifacts",
        component: () => import("@/views/artifact/ArtifactList.vue"),
        meta: { requiresAuth: true, title: "制品管理" },
      },
      {
        path: "/artifacts/:id",
        name: "ArtifactDetail",
        component: () => import("@/views/artifact/ArtifactDetail.vue"),
        meta: { requiresAuth: true, title: "制品详情" },
      },
      {
        path: "dags",
        name: "DagList",
        component: () => import("@/views/dag/DagList.vue"),
        meta: { requiresAuth: true, title: "DAG管理" },
      },
      {
        path: "dags/:id",
        name: "DagDetail",
        component: () => import("@/views/dag/DagDetail.vue"),
        meta: { requiresAuth: true, title: "DAG详情" },
      },
      {
        path: "dags/create",
        name: "DagCreate",
        component: () => import("@/views/dag/DagCreate.vue"),
        meta: { requiresAuth: true, title: "创建DAG" },
      },
      {
        path: "dags/:id/edit",
        name: "DagEdit",
        component: () => import("@/views/dag/DagCreate.vue"),
        meta: { requiresAuth: true, title: "编辑DAG" },
      },
      {
        path: "settings",
        name: "Settings",
        component: () => import("@/views/settings/Settings.vue"),
        meta: { requiresAuth: true, title: "系统设置" },
      },
      {
        path: "resource",
        name: "Resource",
        redirect: "/resource/quota",
        meta: { requiresAuth: true, title: "资源管理" },
        children: [
          {
            path: "quota",
            name: "ResourceQuota",
            component: () =>
              import("@/views/resource/quota/QuotaManagement.vue"),
            meta: { requiresAuth: true, title: "资源配额管理" },
          },
          {
            path: "request",
            name: "ResourceRequest",
            component: () =>
              import("@/views/resource/request/RequestManagement.vue"),
            meta: { requiresAuth: true, title: "资源请求管理" },
          },
          {
            path: "report",
            name: "ResourceReport",
            component: () =>
              import("@/views/resource/report/ReportManagement.vue"),
            meta: { requiresAuth: true, title: "资源报告管理" },
          },
        ],
      },
      {
        path: "template-market",
        name: "TemplateMarket",
        redirect: "/template-market/category",
        meta: { requiresAuth: true, title: "模板市场" },
        children: [
          {
            path: "category",
            name: "TemplateCategory",
            component: () => import("@/views/template-market/category.vue"),
            meta: { requiresAuth: true, title: "分类管理" },
          },
          {
            path: "template",
            name: "TemplateManagement",
            component: () => import("@/views/template-market/template.vue"),
            meta: { requiresAuth: true, title: "模板管理" },
          },
          {
            path: "version",
            name: "VersionManagement",
            component: () => import("@/views/template-market/version.vue"),
            meta: { requiresAuth: true, title: "版本管理" },
          },
          {
            path: "search",
            name: "TemplateSearch",
            component: () => import("@/views/template-market/search.vue"),
            meta: { requiresAuth: true, title: "搜索下载" },
          },
        ],
      },
      {
        path: "webhook",
        name: "webhook",
        component: () => import("@/views/webhook/WebhookList.vue"),
        meta: { requiresAuth: true, title: "webhook" },
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
