<template>
  <el-breadcrumb separator="/">
    <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path">
      {{ item.meta.title }}
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const breadcrumbs = ref([]);

const getBreadcrumbs = () => {
  const matched = route.matched.filter(item => item.meta && item.meta.title);
  breadcrumbs.value = matched;
};

onMounted(() => {
  getBreadcrumbs();
  watch(
      () => route.path,
      () => {
        getBreadcrumbs();
      }
  );
});
</script>
