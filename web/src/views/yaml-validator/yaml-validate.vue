<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>YAML验证</span>
      </div>
    </template>

    <el-form
      ref="validateFormRef"
      :model="validateForm"
      :rules="validateRules"
      label-width="120px"
      class="validate-form"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="validateForm.name" placeholder="请输入名称" />
      </el-form-item>

      <el-form-item label="YAML内容" prop="content">
        <el-input
          v-model="validateForm.content"
          type="textarea"
          placeholder="示例: apiVersion: v1\nkind: Pod\nmetadata:\n  name: example\nspec:\n  containers:\n  - name: example\n    image: nginx:latest"
          :rows="10"
        />
      </el-form-item>

      <el-form-item label="Schema类型" prop="schema_type">
        <el-select
          v-model="validateForm.schema_type"
          placeholder="请选择Schema类型"
        >
          <el-option label="Kubernetes" value="kubernetes" />
          <el-option label="自定义" value="custom" />
        </el-select>
      </el-form-item>

      <el-form-item label="Schema名称" prop="schema_name" width="100%">
        <div class="schema-select-container">
          <el-select
            style="width: 900px"
            v-model="validateForm.schema_name"
            placeholder="请选择Schema名称"
            class="schema-select"
            ref="schemaSelectRef"
            @visible-change="handleSelectVisibleChange"
          >
            <el-option
              v-for="schema in schemas"
              :key="schema.id"
              :label="schema.name"
              :value="schema.name"
            />
            <div v-if="loadingMore" class="loading-more">
              <span>加载中...</span>
            </div>
          </el-select>
        </div>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" @click="handleValidate">验证</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>

    <div v-if="validationResult" class="validation-result">
      <el-divider content-position="left">验证结果</el-divider>
      <div class="result-header">
        <span :class="validationResult.is_valid ? 'valid' : 'invalid'">
          {{ validationResult.is_valid ? "验证通过" : "验证失败" }}
        </span>
      </div>
      <div v-if="!validationResult.is_valid" class="error-message">
        <el-alert
          title="错误信息"
          type="error"
          :closable="false"
          :description="validationResult.errors"
        />
      </div>
    </div>
  </el-card>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, watch, nextTick } from "vue";
import { ElForm, ElMessage, ElLoading, ElSelect } from "element-plus";
import {
  validateYAML,
  getYAMLSchemas,
  YAMLSchema,
  PageResponse,
} from "../../api/yamlValidator";

const validateFormRef = ref<InstanceType<typeof ElForm>>();
const validateForm = reactive<{
  name: string;
  content: string;
  schema_type: string;
  schema_name?: string;
}>({
  name: "",
  content: "",
  schema_type: "kubernetes",
});

const schemas = ref<YAMLSchema[]>([]);
const validationResult = ref<any>(null);
const loading = ref(false);
const loadingSchemas = ref(false);
const loadingMore = ref(false);
const schemaSelectRef = ref<any>();
const hasMore = ref(true);
const currentPage = ref(1);
const pageSize = ref(10);

const validateRules = {
  name: [
    { required: true, message: "请输入名称", trigger: "blur" },
    {
      min: 2,
      max: 100,
      message: "名称长度在 2 到 100 个字符",
      trigger: "blur",
    },
  ],
  content: [{ required: true, message: "请输入YAML内容", trigger: "blur" }],
  schema_type: [
    { required: true, message: "请选择Schema类型", trigger: "change" },
  ],
};

// 获取Schema列表
const fetchSchemas = async (append = false) => {
  if (!hasMore.value && append) return;

  try {
    if (!append) {
      loadingSchemas.value = true;
      schemas.value = [];
      currentPage.value = 1;
      hasMore.value = true;
    } else {
      loadingMore.value = true;
    }

    const res: any = await getYAMLSchemas({
      type: validateForm.schema_type,
      page: currentPage.value,
      pageSize: pageSize.value,
    });

    if (append) {
      schemas.value = [...schemas.value, ...res.data.list];
    } else {
      schemas.value = res.data.list;
    }

    hasMore.value = schemas.value.length < res.data.total;
    currentPage.value++;
  } catch (error) {
    ElMessage.error("获取Schema列表失败");
    console.error("Failed to fetch schemas:", error);
  } finally {
    loadingSchemas.value = false;
    loadingMore.value = false;
  }
};
const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement;
  if (
    target.scrollHeight - target.scrollTop - target.clientHeight < 10 &&
    !loadingMore.value &&
    hasMore.value
  ) {
    fetchSchemas(true);
  }
};

watch(
  () => validateForm.schema_type,
  () => {
    fetchSchemas();
  }
);
const handleSelectVisibleChange = (visible: boolean) => {
  if (visible) {
    nextTick(() => {
      const selectDropdown = document.querySelector(
        ".el-select-dropdown__wrap"
      );
      if (selectDropdown) {
        selectDropdown.addEventListener("scroll", handleScroll);
      }
    });
  } else {
    const selectDropdown = document.querySelector(".el-select-dropdown__wrap");
    if (selectDropdown) {
      selectDropdown.removeEventListener("scroll", handleScroll);
    }
  }
};

// 处理验证
const handleValidate = async () => {
  try {
    await validateFormRef.value?.validate();
    loading.value = true;

    const res = await validateYAML({
      name: validateForm.name,
      content: validateForm.content,
      schema_type: validateForm.schema_type,
      schema_name: validateForm.schema_name,
    });

    validationResult.value = res.data;
    ElMessage.success("验证完成");
  } catch (error: any) {
    if (error.name === "ValidationError") return;
    ElMessage.error(`验证失败: ${error.message || "未知错误"}`);
  } finally {
    loading.value = false;
  }
};

// 重置表单
const resetForm = () => {
  validateFormRef.value?.resetFields();
  validationResult.value = null;
};

// 初始化
onMounted(() => {
  fetchSchemas();
});
</script>

<style scoped>
.validate-form {
  margin-top: 20px;
}

.validation-result {
  margin-top: 30px;
}

.result-header {
  margin-bottom: 10px;
}

.valid {
  color: #67c23a;
  font-weight: bold;
}

.invalid {
  color: #f56c6c;
  font-weight: bold;
}

.error-message {
  margin-top: 10px;
}

.schema-select-container {
  position: relative;
}

.schema-select {
  width: 100%;
}

.loading-more {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  color: #606266;
}

.pagination-container {
  display: none;
}
</style>
