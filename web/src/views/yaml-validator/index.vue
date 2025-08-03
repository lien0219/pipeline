<template>
  <div class="yaml-validator-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>YAML管理中心</span>
        </div>
      </template>
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane label="YAML验证" name="validate" />
        <el-tab-pane label="验证历史" name="history" />
        <el-tab-pane label="Schema管理" name="schema" />
      </el-tabs>
    </el-card>

    <component :is="currentComponent" class="component-container" />
  </div>
</template>

<script lang="ts" setup>
import { ref, computed } from "vue";
import YamlValidate from "./yaml-validate.vue";
import ValidationHistory from "./validation-history.vue";
import SchemaList from "./schema-list.vue";

const activeTab = ref("validate");

const currentComponent = computed(() => {
  switch (activeTab.value) {
    case "validate":
      return YamlValidate;
    case "history":
      return ValidationHistory;
    case "schema":
      return SchemaList;
    default:
      return YamlValidate;
  }
});
</script>

<style scoped>
.yaml-validator-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.component-container {
  margin-top: 20px;
}
</style>
