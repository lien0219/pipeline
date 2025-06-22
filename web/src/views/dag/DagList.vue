<template>
  <div class="dag-list-container">
    <div class="page-header">
      <h1>DAG管理</h1>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon> 创建DAG
      </el-button>
    </div>

    <el-card>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索DAG名称"
        prefix-icon="Search"
        style="width: 300px; margin-bottom: 20px"
      />

      <el-table
        v-loading="dagStore.loading"
        :data="filteredDags"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" />
        <el-table-column label="描述" width="200">
          <template #default="scope">
            <el-tooltip
              :content="scope.row.description"
              placement="top"
              effect="light"
            >
              <div class="dag-description">
                {{ scope.row.description }}
              </div>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column prop="pipeline_id" label="所属流水线ID" width="120" />
        <el-table-column prop="is_active" label="状态" width="100">
          <template #default="scope">
            <el-tag v-if="scope.row.is_active" type="success">活动</el-tag>
            <el-tag v-else>非活动</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" align="center">
          <template #default="scope">
            <el-button size="small" @click="handleView(scope.row.id)"
              >查看</el-button
            >
            <el-button
              size="small"
              type="primary"
              @click="handleEdit(scope.row.id)"
              >编辑</el-button
            >
            <el-button
              size="small"
              :type="scope.row.is_active ? 'default' : 'success'"
              @click="handleActivate(scope.row.id, !scope.row.is_active)"
              :disabled="scope.row.is_active"
            >
              {{ scope.row.is_active ? "已激活" : "激活" }}
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click="handleDelete(scope.row.id)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>

      <!-- 添加分页控件 -->
      <div style="margin-top: 15px; text-align: right">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[5, 10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handlePageChange"
          @size-change="handleSizeChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { useDagStore } from "@/stores/dag";
import { Plus, Search } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { formatDate } from "@/utils/date";

const router = useRouter();
const dagStore = useDagStore();
const searchKeyword = ref("");

const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const tableData = ref([]);

onMounted(async () => {
  await fetchAllDags();
});

// 加载DAG数据
const fetchAllDags = async () => {
  try {
    const result = await dagStore.fetchAllDags(
      currentPage.value,
      pageSize.value
    );
    total.value = result.total;
    tableData.value = result.list;
  } catch (err) {
    ElMessage.error(`获取DAG列表失败: ${dagStore.error}`);
  }
};

// 搜索过滤
const filteredDags = computed(() => {
  if (!searchKeyword.value) return tableData.value;
  return tableData.value?.filter((dag) =>
    dag.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  );
});

// 分页事件处理
const handlePageChange = (page: number) => {
  currentPage.value = page;
  fetchAllDags();
};

// 每页条数变更事件
const handleSizeChange = (size: number) => {
  pageSize.value = size;
  currentPage.value = 1; // 重置到第一页
  fetchAllDags();
};

// 处理创建DAG
const handleCreate = () => {
  router.push("/dags/create");
};

// 处理查看DAG
const handleView = (id: number) => {
  router.push(`/dags/${id}`);
};

// 处理编辑DAG
const handleEdit = (id: number) => {
  router.push(`/dags/${id}/edit`);
};

// 处理激活DAG
const handleActivate = async (id: number, isActive: boolean) => {
  try {
    await dagStore.activateDag(id);
    ElMessage.success(`DAG已${isActive ? "激活" : "取消激活"}`);
    fetchAllDags(); // 刷新数据
  } catch (err) {
    ElMessage.error(`操作失败: ${dagStore.error}`);
  }
};

// 处理删除DAG
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm("确定要删除此DAG吗？", "警告", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await dagStore.deleteDag(id);
    ElMessage.success("删除成功");
    fetchAllDags(); // 刷新数据
  } catch (err) {
    if (err === "cancel") return;
    ElMessage.error(`删除失败: ${dagStore.error}`);
  }
};
</script>

<style scoped>
.dag-list-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.dag-description {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.5;
  max-height: 4.5em;
  word-break: break-word;
}
</style>
