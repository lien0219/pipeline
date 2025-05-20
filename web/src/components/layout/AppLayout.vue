<template>
  <el-container class="app-layout">
    <el-aside :width="isCollapse ? '64px' : '220px'" :class="{ 'is-collapsed': isCollapse }">
      <sidebar :collapse="isCollapse" />
    </el-aside>

    <el-container>
      <el-header>
        <navbar @toggle-sidebar="toggleSidebar" />
      </el-header>

      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref } from 'vue';
import Sidebar from './Sidebar.vue';
import Navbar from './Navbar.vue';

const isCollapse = ref(false);

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value;
};
</script>
