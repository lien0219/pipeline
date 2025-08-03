<template>
  <el-dialog
    v-model="props.visible"
    :title="schemaId ? '编辑Schema' : '新建Schema'"
    width="70%"
  >
    <el-form
      ref="schemaFormRef"
      :model="schemaForm"
      :rules="schemaRules"
      label-width="120px"
      class="schema-form"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="schemaForm.name" placeholder="请输入名称" />
      </el-form-item>

      <el-form-item label="类型" prop="type">
        <el-select
          v-model="schemaForm.type"
          placeholder="请选择类型"
          :disabled="!!schemaId"
        >
          <el-option label="Kubernetes" value="kubernetes" />
          <el-option label="自定义" value="custom" />
        </el-select>
      </el-form-item>

      <el-form-item label="版本" prop="version">
        <el-input v-model="schemaForm.version" placeholder="请输入版本" />
      </el-form-item>

      <el-form-item label="描述" prop="description">
        <el-input
          v-model="schemaForm.description"
          placeholder="请输入描述"
          type="textarea"
        />
      </el-form-item>

      <el-form-item label="Schema内容" prop="schema">
        <el-input
          v-model="schemaForm.schema"
          type="textarea"
          placeholder="请输入JSON Schema内容"
          :rows="15"
          class="monospace"
        />
        <div class="hint-text">
          请输入符合JSON Schema规范的内容，用于验证YAML文件格式
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <!-- @click="visible = false" -->
      <el-button>取消</el-button>
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, onMounted } from "vue";
import { ElForm, ElMessage } from "element-plus";
import {
  createYAMLSchema,
  updateYAMLSchema,
  getYAMLSchemaById,
} from "../../api/yamlValidator";

const props = defineProps<{
  visible: boolean;
  schemaId?: number | null;
}>();

const emit = defineEmits<{
  (e: "update:visible", value: boolean): void;
  (e: "success"): void;
}>();

const schemaFormRef = ref<InstanceType<typeof ElForm>>();
const schemaForm = reactive<{
  name: string;
  type: string;
  version: string;
  description?: string;
  schema: string;
}>({
  name: "",
  type: "kubernetes",
  version: "",
  description: "",
  schema: "",
});

const schemaRules = {
  name: [
    { required: true, message: "请输入名称", trigger: "blur" },
    {
      min: 2,
      max: 100,
      message: "名称长度在 2 到 100 个字符",
      trigger: "blur",
    },
  ],
  type: [{ required: true, message: "请选择类型", trigger: "change" }],
  version: [
    { required: true, message: "请输入版本", trigger: "blur" },
    { min: 1, max: 50, message: "版本长度在 1 到 50 个字符", trigger: "blur" },
  ],
  schema: [{ required: true, message: "请输入Schema内容", trigger: "blur" }],
};

const loading = ref(false);

// 监听visible变化
watch(
  () => props.visible,
  (newVal) => {
    if (newVal && props.schemaId) {
      // 编辑模式，加载数据
      fetchSchemaDetail();
    } else if (newVal && !props.schemaId) {
      // 新建模式，重置表单
      resetForm();
    }
  }
);

// 获取Schema详情
const fetchSchemaDetail = async () => {
  if (!props.schemaId) return;

  try {
    loading.value = true;
    const res = await getYAMLSchemaById(props.schemaId);
    const data = res.data;

    schemaForm.name = data.name;
    schemaForm.type = data.type;
    schemaForm.version = data.version;
    schemaForm.description = data.description || "";
    schemaForm.schema = data.schema;
  } catch (error) {
    ElMessage.error("获取Schema详情失败");
    console.error("Failed to fetch schema detail:", error);
    emit("update:visible", false);
  } finally {
    loading.value = false;
  }
};

// 处理提交
const handleSubmit = async () => {
  try {
    await schemaFormRef?.value?.validate();
    loading.value = true;

    // 验证Schema是否为有效JSON
    try {
      JSON.parse(schemaForm.schema);
    } catch (error) {
      ElMessage.error("Schema内容必须是有效的JSON格式");
      return;
    }

    if (props.schemaId) {
      // 更新Schema
      await updateYAMLSchema(props.schemaId, {
        name: schemaForm.name,
        version: schemaForm.version,
        schema: schemaForm.schema,
        description: schemaForm.description,
      });
      ElMessage.success("更新成功");
    } else {
      // 创建Schema
      await createYAMLSchema({
        name: schemaForm.name,
        type: schemaForm.type,
        version: schemaForm.version,
        schema: schemaForm.schema,
        description: schemaForm.description,
      });
      ElMessage.success("创建成功");
    }

    emit("success");
    emit("update:visible", false);
  } catch (error: any) {
    if (error.name === "ValidationError") return;
    ElMessage.error(`操作失败: ${error.message || "未知错误"}`);
  } finally {
    loading.value = false;
  }
};

// 重置表单
const resetForm = () => {
  schemaFormRef?.value?.resetFields();
  schemaForm.type = "kubernetes";
};

// 初始化
onMounted(() => {
  if (props.visible && props.schemaId) {
    fetchSchemaDetail();
  }
});
</script>

<style scoped>
.schema-form {
  margin-top: 20px;
}

.monospace {
  font-family: monospace;
}

.hint-text {
  color: #606266;
  font-size: 12px;
  margin-top: 5px;
}
</style>
