<template>
  <div class="sidebar">
    <div class="logo-container">
      <router-link to="/">
        <img
          v-if="!collapse"
          src="@/assets/vue.svg"
          alt="CICD Pipeline"
          class="logo-img"
        />
        <img v-else src="@/assets/vue.svg" alt="CICD" class="logo-icon" />
      </router-link>
    </div>

    <el-scrollbar>
      <el-menu
        :default-active="activeMenu"
        :collapse="collapse"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        unique-opened
        router
      >
        <el-menu-item index="/">
          <el-icon><Monitor /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>

        <el-menu-item index="/pipelines">
          <el-icon><Connection /></el-icon>
          <template #title>流水线</template>
        </el-menu-item>

        <el-sub-menu index="builds">
          <template #title>
            <el-icon><Box /></el-icon>
            <span>构建管理</span>
          </template>
          <el-menu-item index="/builds/history">构建历史</el-menu-item>
          <el-menu-item index="/builds/templates">构建模板</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="deploy">
          <template #title>
            <el-icon><Upload /></el-icon>
            <span>部署管理</span>
          </template>
          <el-menu-item index="/deploy/environments">环境管理</el-menu-item>
          <el-menu-item index="/deploy/releases">发布记录</el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/artifacts">
          <el-icon><Files /></el-icon>
          <template #title>制品管理</template>
        </el-menu-item>

        <el-menu-item index="/dags">
          <el-icon><DataAnalysis /></el-icon>
          <template #title>DAG管理</template>
        </el-menu-item>
        <el-sub-menu index="/resource">
          <template #title>
            <el-icon><TrendCharts /></el-icon>
            <span>资源管理</span>
          </template>
          <el-menu-item index="/resource/quota">资源配额管理</el-menu-item>
          <el-menu-item index="/resource/request">资源请求管理</el-menu-item>
          <el-menu-item index="/resource/report">资源报告管理</el-menu-item>
        </el-sub-menu>

        <el-sub-menu index="/template-market">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>模板市场</span>
          </template>
          <el-menu-item index="/template-market/category"
            >分类管理</el-menu-item
          >
          <el-menu-item index="/template-market/template"
            >模板管理</el-menu-item
          >
          <!-- <el-menu-item index="/template-market/version">版本管理</el-menu-item> -->
          <el-menu-item index="/template-market/search">搜索下载</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/webhook">
          <el-icon><Files /></el-icon>
          <template #title>webhook</template>
        </el-menu-item>
        <el-sub-menu index="/yaml-validator">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>YAML管理</span>
          </template>
          <el-menu-item index="/yaml-validator">YAML验证中心</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>系统设置</template>
        </el-menu-item>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useRoute } from "vue-router";
import { ref } from "vue";
import {
  Monitor,
  Connection,
  Box,
  Upload,
  Files,
  Setting,
  DataAnalysis,
  TrendCharts,
  Document,
} from "@element-plus/icons-vue";

const props = defineProps({
  collapse: {
    type: Boolean,
    default: false,
  },
});

const route = useRoute();
const activeMenu = computed(() => {
  return route.path;
});
</script>

<style scoped>
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.logo-container {
  height: 60px;
  padding: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b2f3a;
}

.logo-img {
  height: 32px;
  max-width: 180px;
}

.logo-icon {
  height: 32px;
  width: 32px;
}

.el-menu {
  border-right: none;
}
</style>
