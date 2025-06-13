<template>
  <div class="designer-container">
    <div class="panel node-library">
      <h3>节点库</h3>
      <div class="node-list">
        <div
          v-for="type in nodeTypes"
          :key="type.type"
          class="node-item"
          draggable="true"
          @dragstart="handleDragStart(type)"
        >
          {{ type.label }}
        </div>
      </div>
    </div>
    <div
      class="panel canvas"
      ref="canvasRef"
      @dragover.prevent
      @drop="handleCanvasDrop"
    >
      <svg
        class="dag-edges"
        :style="{ width: `${canvasWidth}px`, height: `${canvasHeight}px` }"
      >
        <path
          v-for="(edge, idx) in currentDAG.edges"
          :key="idx"
          :d="getEdgePath(edge)"
          stroke="#3b82f6"
          stroke-width="1"
          fill="none"
        />
        <path
          v-for="(tempEdge, idx) in tempEdges"
          :key="`temp-${idx}`"
          :d="getEdgePath(tempEdge)"
          stroke="#9ca3af"
          stroke-width="2"
          stroke-dasharray="4 2"
          fill="none"
          opacity="0.8"
        />
      </svg>
      <div
        class="dag-canvas"
        :style="{ width: `${canvasWidth}px`, height: `${canvasHeight}px` }"
      >
        <div
          v-for="node in currentDAG.nodes"
          :key="node.id"
          :style="{ left: `${node.position.x}px`, top: `${node.position.y}px` }"
          class="dag-node"
          :class="{ selected: node.id === selectedNodeId }"
          @mousedown="handleNodeMouseDown(node)"
        >
          <div class="node-header">
            {{ node.type.label }}
            <el-button
              type="text"
              size="mini"
              @click.stop="deleteNode(node.id)"
              class="delete-btn"
              >×</el-button
            >
          </div>
          <div class="node-content">{{ node.name }}</div>
          <div class="node-ports">
            <div
              class="port in"
              v-if="node.type.acceptsInput"
              @mousedown="handlePortMouseDown('in', node)"
            ></div>
            <div
              class="port out"
              v-if="node.type.providesOutput"
              @mousedown="handlePortMouseDown('out', node)"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <div class="panel properties" v-if="selectedNode">
      <h3>节点属性</h3>
      <el-form :model="selectedNode" label-width="80px">
        <el-form-item label="名称">
          <el-input v-model="selectedNode.name" />
        </el-form-item>
        <el-form-item label="超时时间">
          <el-input-number v-model="selectedNode.timeout" min="0" />
        </el-form-item>
        <el-form-item v-if="selectedNode.type === 'task'">
          <label>脚本内容</label>
          <el-input type="textarea" v-model="selectedNode.script" />
        </el-form-item>
      </el-form>
      <el-button type="primary" @click="saveDAG">保存配置</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRouter, useRoute } from "vue-router";
import { pipelineApi } from "@/api/pipeline";
import { DAG, DAGNode, NodeType, DAGEdge } from "@/types/dag";
import { ElMessage } from "element-plus";

const nodeTypes: NodeType[] = [
  { type: "task", label: "任务节点", acceptsInput: true, providesOutput: true },
  { type: "condition", label: "条件节点", acceptsInput: true, providesOutput: true },
  { type: "parallel", label: "并行节点", acceptsInput: true, providesOutput: true },
  { type: "approval", label: "审批节点", acceptsInput: true, providesOutput: false },
];

const router = useRouter();
const route = useRoute();
const pipelineId = route.params.id as string;
const currentDAG = ref<DAG>({ nodes: [], edges: [] });
const selectedNodeId = ref<string>("");
const canvasRef = ref<HTMLDivElement>(null);
const draggingNode = ref<DAGNode | null>(null);
const canvasWidth = ref(1000);
const canvasHeight = ref(1000);
const tempEdges = ref<DAGEdge[]>([]);
const CONNECT_THRESHOLD = 200; 

const selectedNode = computed(() => {
  return currentDAG.value.nodes.find(
    (node) => node.id === selectedNodeId.value
  );
});

onMounted(async () => {
  try {
    const res = await pipelineApi.getDAGByPipelineID(pipelineId);
    currentDAG.value = res.data;
  } catch (error) {
    ElMessage.error("加载DAG失败");
  }
});

let dragType: NodeType | null = null;
const handleDragStart = (type: NodeType) => {
  dragType = type;
};

const handleCanvasDrop = (e: DragEvent) => {
  if (!dragType || !canvasRef.value) return;
  const rect = canvasRef.value.getBoundingClientRect();
  const x = e.clientX - rect.left;
  const y = e.clientY - rect.top;
  currentDAG.value.nodes.push({
    id: `node-${Date.now()}`,
    type: dragType.type,
    name: `${dragType.label}-${currentDAG.value.nodes.length + 1}`,
    position: { x, y },
    timeout: 300,
    script: "",
  });
  dragType = null;
};

let animationFrameId: number | null = null;
const handleNodeMouseMove = (e: MouseEvent) => {
  if (!draggingNode.value) return;
  if (animationFrameId) cancelAnimationFrame(animationFrameId);
  animationFrameId = requestAnimationFrame(() => {
    const deltaX = e.clientX - startX;
    const deltaY = e.clientY - startY;
    draggingNode.value!.position = {
      x: initialPosition.x + deltaX,
      y: initialPosition.y + deltaY,
    };
    calculateTempEdges(); 
  });
};

const calculateTempEdges = () => {
  tempEdges.value = [];
  const nodes = currentDAG.value.nodes;
  for (let i = 0; i < nodes.length; i++) {
    for (let j = i + 1; j < nodes.length; j++) {
      const nodeA = nodes[i];
      const nodeB = nodes[j];
      const distance = Math.hypot(
        nodeA.position.x - nodeB.position.x,
        nodeA.position.y - nodeB.position.y
      );
      if (distance < CONNECT_THRESHOLD) {
        tempEdges.value.push({
          source: nodeA.id,
          target: nodeB.id,
        });
      }
    }
  }
};

let startX = 0;
let startY = 0;
let initialPosition = { x: 0, y: 0 };
const handleNodeMouseDown = (node: DAGNode) => {
  selectedNodeId.value = node.id;
  draggingNode.value = node;
  initialPosition = { ...node.position };
  startX = window.event!.clientX;
  startY = window.event!.clientY;
  document.addEventListener("mousemove", handleNodeMouseMove);
  document.addEventListener("mouseup", handleNodeMouseUp);
};

const handleNodeMouseUp = () => {
  draggingNode.value = null;
  if (animationFrameId) cancelAnimationFrame(animationFrameId);
  document.removeEventListener("mousemove", handleNodeMouseMove);
  document.removeEventListener("mouseup", handleNodeMouseUp);
  calculateTempEdges(); 
};

let edgeStart: { nodeId: string; port: "in" | "out" } | null = null;
const handlePortMouseDown = (port: "in" | "out", node: DAGNode) => {
  if (port === "out") {
    edgeStart = { nodeId: node.id, port };
  }
};

const getEdgePath = (edge: DAGEdge) => {
  const startNode = currentDAG.value.nodes.find((n) => n.id === edge.source)!;
  const endNode = currentDAG.value.nodes.find((n) => n.id === edge.target)!;
  const startX = startNode.position.x + 150; 
  const startY = startNode.position.y + 12; 
  const endX = endNode.position.x;
  const endY = endNode.position.y + 12;
  return `M ${startX} ${startY} C ${startX + 50} ${startY}, ${endX - 50} ${endY}, ${endX} ${endY}`;
};

const deleteNode = (nodeId: string) => {
  currentDAG.value.nodes = currentDAG.value.nodes.filter((n) => n.id !== nodeId);
  currentDAG.value.edges = currentDAG.value.edges.filter(
    (e) => e.source !== nodeId && e.target !== nodeId
  );
  selectedNodeId.value = "";
};

const saveDAG = async () => {
  try {
    await pipelineApi.updateDAG(pipelineId, currentDAG.value);
    ElMessage.success("DAG保存成功");
  } catch (error) {
    ElMessage.error("DAG保存失败");
  }
};

onUnmounted(() => {
  document.removeEventListener("mousemove", handleNodeMouseMove);
  document.removeEventListener("mouseup", handleNodeMouseUp);
});
</script>

<style scoped>
.designer-container {
  display: flex;
  height: 100vh;
  padding: 20px;
  gap: 20px;
}

.panel {
  flex: 1;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  overflow: auto;
  position: relative;
}

.node-library {
  flex: 0 0 200px;
}

.node-item {
  padding: 8px 12px;
  margin: 8px 0;
  background: #f3f4f6;
  border-radius: 4px;
  cursor: move;
  user-select: none;
}

.canvas {
  flex: 3;
}

.dag-canvas {
  position: absolute;
  top: 16px;
  left: 16px;
}

.dag-node {
  position: absolute;
  width: 150px;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  cursor: move;
  transition: all 0.1s;
}

.dag-node.selected {
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.node-header {
  padding: 4px 8px;
  background: #3b82f6;
  color: white;
  border-radius: 4px 4px 0 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.delete-btn {
  color: white;
  padding: 0 4px;
}

.node-content {
  padding: 8px;
}

.node-ports {
  display: flex;
  justify-content: space-between;
  padding: 0 8px 4px;
}

.port {
  width: 12px;
  height: 12px;
  background: #e5e7eb;
  border-radius: 50%;
  cursor: crosshair;
}

.properties {
  flex: 0 0 300px;
}
.dag-edges path[stroke-dasharray] {
  opacity: 0.8;
}
</style>