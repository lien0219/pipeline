<template>
  <div class="dag-create-container">
    <div class="page-header">
      <h1>{{ isEdit ? "编辑DAG" : "创建DAG" }}</h1>
      <el-button @click="handleCancel">取消</el-button>
    </div>

    <el-card v-loading="dagStore.loading">
      <el-form
        ref="dagFormRef"
        :model="dagForm"
        :rules="formRules"
        label-width="120px"
        style="max-width: 800px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="dagForm.name" placeholder="请输入DAG名称" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="dagForm.description"
            placeholder="请输入DAG描述"
            type="textarea"
            rows="3"
          />
        </el-form-item>

        <el-form-item label="所属流水线ID" prop="pipelineId">
          <el-input
            v-model.number="dagForm.pipeline_id"
            placeholder="请输入所属流水线ID"
            type="number"
          />
        </el-form-item>

        <el-form-item label="DAG节点配置" prop="nodes">
          <el-input
            v-model="nodesJson"
            placeholder="请输入JSON格式的节点配置"
            type="textarea"
            rows="10"
          />
          <el-button type="text" @click="validateDag" style="margin-top: 10px">
            验证DAG配置
          </el-button>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit">提交</el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useDagStore } from "@/stores/dag";
import { ElMessage, ElMessageBox } from "element-plus";

const router = useRouter();
const route = useRoute();
const dagStore = useDagStore();
const dagFormRef = ref();
const isEdit = ref(!!route.params.id);

// 表单数据
const dagForm = ref({
  name: "",
  description: "",
  pipeline_id: 0,
  nodes: [],
});

const nodesJson = ref("[]");

// 表单验证规则
const formRules = ref({
  name: [
    { required: true, message: "请输入DAG名称", trigger: "blur" },
    { min: 1, max: 50, message: "名称长度在1到50个字符", trigger: "blur" },
  ],
  pipeline_id: [
    { required: true, message: "请输入所属流水线ID", trigger: "blur" },
    { type: "number", message: "必须是数字", trigger: "blur" },
  ],
  nodes: [
    {
      required: true,
      message: "请输入节点配置",
      trigger: ["blur", "change"],
      validator: (rule, value, callback) => {
        if (Array.isArray(value) && value.length === 0) {
          callback();
        } else if (!value) {
          callback(new Error("请输入节点配置"));
        } else {
          callback();
        }
      },
    },
  ],
});

// 监听nodesJson变化，同步到dagForm.nodes
watch(
  () => nodesJson.value,
  (val) => {
    if (val) {
      try {
        dagForm.value.nodes = JSON.parse(val);
      } catch (err) {
        // JSON格式错误，不更新
      }
    }
  }
);

onMounted(async () => {
  if (isEdit.value) {
    const dagId = Number(route.params.id);
    try {
      const dag = await dagStore.fetchDagDetail(dagId);
      dagForm.value = {
        name: dag.data.name,
        description: dag.data.description,
        pipeline_id: dag.data.pipeline_id,
        nodes: dag.data.nodes,
      };
      // 格式化JSON显示
      nodesJson.value = JSON.stringify(dag.data.nodes, null, 2);
    } catch (err) {
      ElMessage.error(`加载DAG数据失败: ${dagStore.error}`);
      router.push("/dags");
    }
  }
});

// 验证DAG配置
const validateDag = async () => {
  if (!nodesJson.value) {
    ElMessage.warning("请输入节点配置");
    return;
  }

  try {
    const nodes = JSON.parse(nodesJson.value);
    await dagStore.validateDag(nodes);
    ElMessage.success("DAG配置验证通过");
  } catch (err) {
    ElMessage.error(
      `验证失败: ${err instanceof Error ? err.message : dagStore.error}`
    );
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!dagFormRef.value) return;

  try {
    // 先验证JSON格式
    if (nodesJson.value) {
      dagForm.value.nodes = JSON.parse(nodesJson.value);
    }

    // 表单验证
    await dagFormRef.value.validate();

    if (isEdit.value) {
      // 编辑DAG
      await dagStore.updateDag(Number(route.params.id), dagForm.value);
      ElMessage.success("DAG更新成功");
    } else {
      // 创建DAG
      await dagStore.createDag(dagForm.value);
      ElMessage.success("DAG创建成功");
    }

    router.push("/dags");
  } catch (err) {
    if (err instanceof Error && err.message.includes("Validation failed")) {
      return;
    }
    ElMessage.error(`${isEdit.value ? "更新" : "创建"}失败: ${dagStore.error}`);
  }
};

// 取消操作
const handleCancel = () => {
  router.push("/dags");
};
</script>

<style scoped>
.dag-create-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
</style>
