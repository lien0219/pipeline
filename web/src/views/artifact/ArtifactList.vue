<template>
  <div class="app-container">
    <div class="page-header">
      <div class="header-title">
        <h2>制品管理</h2>
      </div>

      <div class="header-actions">
        <el-button type="primary" @click="uploadArtifact">
          <el-icon><Upload /></el-icon>
          上传制品
        </el-button>
      </div>
    </div>

    <el-card>
      <div class="filter-container">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="名称">
            <el-input
              v-model="filterForm.name"
              placeholder="制品名称"
              clearable
            />
          </el-form-item>

          <el-form-item label="类型" style="width: 200px">
            <el-select
              v-model="filterForm.type"
              placeholder="全部类型"
              clearable
            >
              <el-option label="ZIP" value="zip" />
              <el-option label="TAR" value="tar" />
              <el-option label="JAR" value="jar" />
              <el-option label="WAR" value="war" />
              <el-option label="Docker镜像" value="docker" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>

          <el-form-item label="流水线" style="width: 200px">
            <el-select
              v-model="filterForm.pipeline_id"
              placeholder="全部流水线"
              clearable
            >
              <el-option
                v-for="pipeline in pipelines"
                :key="pipeline.id"
                :label="pipeline.name"
                :value="pipeline.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="时间范围">
            <el-date-picker
              v-model="filterForm.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="handleFilter">
              <el-icon><Search /></el-icon>
              搜索
            </el-button>
            <el-button @click="resetFilter">
              <el-icon><RefreshRight /></el-icon>
              重置
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <el-table
        :data="artifacts"
        style="width: 100%"
        v-loading="loading"
        @sort-change="handleSortChange"
      >
        <el-table-column
          prop="name"
          label="名称"
          min-width="200"
          sortable="custom"
          show-overflow-tooltip
        />

        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type.toUpperCase() }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="size" label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>

        <el-table-column
          prop="pipeline.name"
          label="流水线"
          width="150"
          show-overflow-tooltip
        >
          <template #default="{ row }">
            <router-link
              v-if="row.pipeline_id"
              :to="`/pipelines/${row.pipeline_id}`"
              class="pipeline-link"
            >
              {{ row.pipeline?.name }}
            </router-link>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column prop="pipeline_run_id" label="构建ID" width="100">
          <template #default="{ row }">
            <router-link
              v-if="row.pipeline_id && row.pipeline_run_id"
              :to="`/pipelines/${row.pipeline_id}/runs/${row.pipeline_run_id}`"
              class="pipeline-link"
            >
              {{ row.pipeline_run_id }}
            </router-link>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column
          prop="download_count"
          label="下载次数"
          width="100"
          sortable="custom"
        />

        <el-table-column prop="created_by" label="创建者" width="120" />

        <el-table-column
          prop="created_at"
          label="创建时间"
          width="180"
          sortable="custom"
        >
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              link
              type="primary"
              size="small"
              @click="downloadArtifact(row.id)"
            >
              下载
            </el-button>

            <el-button
              link
              type="primary"
              size="small"
              @click="viewArtifactDetail(row.id)"
            >
              详情
            </el-button>

            <el-button
              link
              type="danger"
              size="small"
              @click="deleteArtifact(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 30, 50]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <el-dialog v-model="uploadDialogVisible" title="上传制品" width="50%">
      <el-form
        ref="uploadFormRef"
        :model="uploadFormData"
        :rules="uploadRules"
        label-position="top"
      >
        <el-form-item label="制品名称" prop="name">
          <el-input
            v-model="uploadFormData.name"
            placeholder="请输入制品名称"
          />
        </el-form-item>

        <el-form-item label="制品类型" prop="type">
          <el-select v-model="uploadFormData.type" placeholder="请选择制品类型">
            <el-option label="RAR" value="rar" />
            <el-option label="ZIP" value="zip" />
            <el-option label="TAR" value="tar" />
            <el-option label="JAR" value="jar" />
            <el-option label="WAR" value="war" />
            <el-option label="Docker镜像" value="docker" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>

        <el-form-item label="关联流水线" prop="pipeline_id">
          <el-select
            v-model="uploadFormData.pipeline_id"
            placeholder="请选择关联流水线"
            clearable
          >
            <el-option
              v-for="pipeline in pipelines"
              :key="pipeline.id"
              :label="pipeline.name"
              :value="pipeline.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="版本" prop="version">
          <el-input
            v-model="uploadFormData.version"
            placeholder="请输入版本号，例如：v1.0.0"
          />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="uploadFormData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入制品描述"
          />
        </el-form-item>

        <el-form-item
          label="文件"
          prop="file"
          v-if="uploadFormData.type !== 'docker'"
        >
          <el-upload
            class="artifact-upload"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处，或 <em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">请上传制品文件，大小不超过500MB</div>
            </template>
          </el-upload>
        </el-form-item>

        <el-form-item
          label="Docker镜像"
          prop="docker_image"
          v-if="uploadFormData.type === 'docker'"
        >
          <el-input
            v-model="uploadFormData.docker_image"
            placeholder="请输入Docker镜像名称，例如：registry.example.com/app:v1.0.0"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitUpload" :loading="uploading">
            上传
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog v-model="detailDialogVisible" title="制品详情" width="60%">
      <div v-loading="detailLoading">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">{{
            currentArtifact.id
          }}</el-descriptions-item>
          <el-descriptions-item label="名称">{{
            currentArtifact.name
          }}</el-descriptions-item>
          <el-descriptions-item label="类型">
            <el-tag size="small">{{
              currentArtifact.type?.toUpperCase()
            }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="大小">{{
            formatFileSize(currentArtifact.size)
          }}</el-descriptions-item>
          <el-descriptions-item label="版本">{{
            currentArtifact.version
          }}</el-descriptions-item>
          <el-descriptions-item label="下载次数">{{
            currentArtifact.download_count
          }}</el-descriptions-item>
          <el-descriptions-item label="流水线">
            <router-link
              v-if="currentArtifact.pipeline_id"
              :to="`/pipelines/${currentArtifact.pipeline_id}`"
              class="pipeline-link"
            >
              {{ currentArtifact.pipeline?.name }}
            </router-link>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="构建ID">
            <router-link
              v-if="
                currentArtifact.pipeline_id && currentArtifact.pipeline_run_id
              "
              :to="`/pipelines/${currentArtifact.pipeline_id}/runs/${currentArtifact.pipeline_run_id}`"
              class="pipeline-link"
            >
              {{ currentArtifact.pipeline_run_id }}
            </router-link>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="创建者">{{
            currentArtifact.created_by
          }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{
            formatDate(currentArtifact.created_at)
          }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">{{
            currentArtifact.description || "无描述"
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="artifact-actions">
          <el-button
            type="primary"
            @click="downloadArtifact(currentArtifact.id)"
          >
            <el-icon><Download /></el-icon>
            下载制品
          </el-button>

          <el-button type="success" @click="deployArtifact(currentArtifact.id)">
            <el-icon><Position /></el-icon>
            部署制品
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, getCurrentInstance } from "vue";
import { ElMessage, ElMessageBox, ElLoading } from "element-plus";
import {
  Upload,
  Search,
  RefreshRight,
  UploadFilled,
  Download,
  Position,
} from "@element-plus/icons-vue";
import dayjs from "dayjs";
import { artifactApi } from "@/api/artifact";
import { pipelineApi } from "@/api/pipeline";

const artifacts = ref([]);
const pipelines = ref([]);
const loading = ref(false);
const detailLoading = ref(false);
const uploading = ref(false);
const currentArtifact = ref({});
const uploadDialogVisible = ref(false);
const detailDialogVisible = ref(false);
const uploadFormRef = ref(null);
const selectedFile = ref(null);

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0,
});
const filterForm = reactive({
  name: "",
  type: "",
  pipeline_id: "",
  dateRange: [],
});

// 排序
const sortParams = reactive({
  prop: "created_at",
  order: "descending",
});

// 上传表单
const uploadFormData = reactive({
  name: "",
  type: "",
  pipeline_id: "",
  version: "",
  description: "",
  docker_image: "",
  file: null,
});

const uploadRules = computed(() => ({
  name: [{ required: true, message: "请输入制品名称", trigger: "blur" }],
  type: [{ required: true, message: "请选择制品类型", trigger: "change" }],
  version: [{ required: true, message: "请输入版本号", trigger: "blur" }],
  ...(uploadFormData.type === "docker"
    ? {
        docker_image: [
          { required: true, message: "请输入Docker镜像名称", trigger: "blur" },
        ],
      }
    : {
        file: [
          {
            required: true,
            type: "file",
            message: "请选择文件",
            trigger: "change",
          },
          {
            type: "file",
            message: "文件格式不正确",
            trigger: "change",
            validator: (rule, value, callback) => {
              if (!value) return callback();
              const allowedTypes = ["rar", "zip", "tar", "jar", "war"];
              const fileExt = value.name.split(".").pop().toLowerCase();
              if (allowedTypes.includes(fileExt)) {
                callback();
              } else if (value.size > 500 * 1024 * 1024) {
                callback(new Error("文件大小不能超过500MB"));
              } else {
                callback(new Error(`仅支持${allowedTypes.join(", ")}格式`));
              }
            },
          },
        ],
      }),
}));

// 上传制品
const uploadArtifact = () => {
  // 重置表单
  Object.assign(uploadFormData, {
    name: "",
    type: "zip",
    pipeline_id: null,
    version: "",
    description: "",
    file: null,
    docker_image: "",
  });

  uploadDialogVisible.value = true;
};

// 处理文件变更
const handleFileChange = (file, fileList) => {
  if (fileList.length > 0) {
    selectedFile.value = file.raw;
    uploadFormData.file = file.raw;
    uploadFormRef.value?.validateField("file");
  } else {
    selectedFile.value = null;
    uploadFormData.file = null;
    uploadFormRef.value?.validateField("file");
  }
};

// 提交上传
const submitUpload = async () => {
  // const form = uploadFormRef.value;
  // if (!form) return;

  // const validateField =
  //   uploadFormData.type === "docker" ? "docker_image" : "file";
  // const validateResult = await form.validateField(validateField);
  if (uploadFormData.type !== "docker") {
    delete uploadFormData.docker_image;
  }

  // console.log("表单数据:", uploadFormData);
  const validateResult = await uploadFormRef.value.validate();
  // console.log("验证结果:", validateResult);

  // if (validateResult) return;

  uploading.value = true;
  try {
    const formData = new FormData();
    formData.append("name", uploadFormData.name);
    formData.append("version", uploadFormData.version);
    formData.append("type", uploadFormData.type);
    formData.append("description", uploadFormData.description || "");
    if (uploadFormData.pipeline_id) {
      formData.append("pipeline_id", uploadFormData.pipeline_id);
    }

    if (uploadFormData.type === "docker") {
      formData.append("docker_image", uploadFormData.docker_image);
    } else if (selectedFile.value) {
      formData.append("file", selectedFile.value);
    }

    await artifactApi.createArtifact(formData);
    ElMessage.success("制品上传成功");
    uploadDialogVisible.value = false;
    fetchArtifacts();

    Object.assign(uploadFormData, {
      name: "",
      version: "",
      type: "",
      pipeline_id: "",
      description: "",
      docker_image: "",
      file: null,
    });
    selectedFile.value = null;
  } catch (error) {
    ElMessage.error("制品上传失败");
    console.error(error);
  } finally {
    uploading.value = false;
  }
};

// 下载制品
const downloadArtifact = async (id) => {
  try {
    const response = await artifactApi.downloadArtifact(id);
    const blob = new Blob([response.data]);
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    const fileName = `artifact-${id}.${currentArtifact.value?.type || "bin"}`;
    a.download = fileName;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);

    // 更新下载次数
  } catch (error) {
    ElMessage.error("制品下载失败");
    console.error(error);
  }
};

// 查看制品详情
const viewArtifactDetail = async (id) => {
  detailLoading.value = true;
  try {
    const response = await artifactApi.getArtifactById(id);
    currentArtifact.value = response.data.list;
    detailDialogVisible.value = true;
  } catch (error) {
    ElMessage.error("获取制品详情失败");
    console.error(error);
  } finally {
    detailLoading.value = false;
  }
};

// 部署制品
const deployArtifact = (id) => {
  // 可以跳转到部署页面或打开部署对话框
  ElMessage.info("部署功能将在后续实现");
};

// 删除制品
const deleteArtifact = async (id) => {
  try {
    await ElMessageBox.confirm("确定要删除此制品吗？此操作不可撤销。", "警告", {
      confirmButtonText: "确定",
      cancelButtonText: "取消",
      type: "warning",
    });

    await artifactApi.deleteArtifact(id);
    ElMessage.success("制品删除成功");
    fetchArtifacts();
  } catch (error) {
    if (error === "cancel") return;
    ElMessage.error("制品删除失败");
    console.error(error);
  }
};

// 筛选
const handleFilter = () => {
  pagination.currentPage = 1;
  fetchArtifacts();
};

// 重置筛选
const resetFilter = () => {
  Object.assign(filterForm, {
    name: "",
    type: "",
    pipeline_id: "",
    dateRange: [],
  });
  pagination.currentPage = 1;
  fetchArtifacts();
};
// 分页处理
const handleSizeChange = (val) => {
  pagination.pageSize = val;
  fetchArtifacts();
};

const handleCurrentChange = (val) => {
  pagination.currentPage = val;
  fetchArtifacts();
};

// 排序处理
const handleSortChange = ({ prop, order }) => {
  if (prop) {
    sortParams.prop = prop;
    sortParams.order = order;
  } else {
    sortParams.prop = "created_at";
    sortParams.order = "descending";
  }
};

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return "0 Bytes";
  const k = 1024;
  const sizes = ["Bytes", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

// 格式化日期
const formatDate = (dateString) => {
  return dayjs(dateString).format("YYYY-MM-DD HH:mm:ss");
};
const fetchPipelines = async () => {
  try {
    const response = await pipelineApi.getPipelines({
      page: 1,
      limit: 50,
    });
    pipelines.value = response.data.list;
  } catch (error) {
    ElMessage.error("获取流水线列表失败");
    console.error(error);
  }
};
const fetchArtifacts = async () => {
  loading.value = true;
  try {
    const filters = {
      name: filterForm.name,
      type: filterForm.type,
      pipeline_id: filterForm.pipeline_id,
      start_date: filterForm.dateRange[0] || "",
      end_date: filterForm.dateRange[1] || "",
    };

    const response = await artifactApi.getArtifacts(
      pagination.currentPage,
      pagination.pageSize,
      filters
    );

    artifacts.value = response.data.list;
    pagination.total = response.data.total;
  } catch (error) {
    ElMessage.error("获取制品列表失败");
    console.error(error);
  } finally {
    loading.value = false;
  }
};
onMounted(() => {
  fetchPipelines();
  fetchArtifacts();
});
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header-title h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  font-weight: 600;
}

.header-title p {
  margin: 0;
  color: #606266;
}

.filter-container {
  margin-bottom: 20px;
}

.filter-form {
  display: flex;
  flex-wrap: wrap;
}

.pipeline-link {
  color: var(--el-color-primary);
  text-decoration: none;
}

.pipeline-link:hover {
  text-decoration: underline;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.artifact-upload {
  width: 100%;
}

.artifact-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
  justify-content: center;
}
</style>
