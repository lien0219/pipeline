<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h1>仪表盘</h1>
      <p>欢迎回来，{{ user?.username }}</p>
    </div>

    <div class="dashboard-stats">
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <h3>总访问量</h3>
          <p class="stat-value">1,234</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-content">
          <h3>用户数</h3>
          <p class="stat-value">567</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📝</div>
        <div class="stat-content">
          <h3>内容数</h3>
          <p class="stat-value">89</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🔄</div>
        <div class="stat-content">
          <h3>转化率</h3>
          <p class="stat-value">12.3%</p>
        </div>
      </div>
    </div>

    <div class="dashboard-content">
      <div class="dashboard-card">
        <h2>最近活动</h2>
        <div class="activity-list">
          <div class="activity-item">
            <div class="activity-time">今天 10:30</div>
            <div class="activity-details">
              <h4>系统更新</h4>
              <p>系统已更新到最新版本 v1.2.0</p>
            </div>
          </div>

          <div class="activity-item">
            <div class="activity-time">昨天 15:45</div>
            <div class="activity-details">
              <h4>新用户注册</h4>
              <p>有5名新用户完成了注册</p>
            </div>
          </div>

          <div class="activity-item">
            <div class="activity-time">2天前</div>
            <div class="activity-details">
              <h4>数据备份</h4>
              <p>系统数据已成功备份</p>
            </div>
          </div>
        </div>
      </div>

      <div class="dashboard-card">
        <h2>系统状态</h2>
        <div class="system-status">
          <div class="status-item">
            <div class="status-label">CPU使用率</div>
            <div class="status-bar">
              <div class="status-progress" style="width: 45%"></div>
            </div>
            <div class="status-value">45%</div>
          </div>

          <div class="status-item">
            <div class="status-label">内存使用率</div>
            <div class="status-bar">
              <div class="status-progress" style="width: 60%"></div>
            </div>
            <div class="status-value">60%</div>
          </div>

          <div class="status-item">
            <div class="status-label">磁盘使用率</div>
            <div class="status-bar">
              <div class="status-progress" style="width: 25%"></div>
            </div>
            <div class="status-value">25%</div>
          </div>

          <div class="status-item">
            <div class="status-label">网络流量</div>
            <div class="status-bar">
              <div class="status-progress" style="width: 70%"></div>
            </div>
            <div class="status-value">70%</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const user = computed(() => authStore.currentUser);

onMounted(() => {
  // 可以在这里加载仪表盘数据
});
</script>

<style scoped>
.dashboard-page {
  width: 100%;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.page-header p {
  color: var(--text-light);
}

.dashboard-stats {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background-color: var(--surface-color);
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  font-size: 2rem;
  margin-right: 1rem;
}

.stat-content h3 {
  font-size: 1rem;
  font-weight: 500;
  color: var(--text-light);
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--primary-color);
}

.dashboard-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
}

.dashboard-card {
  background-color: var(--surface-color);
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.dashboard-card h2 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: var(--text-color);
}

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.activity-item {
  display: flex;
  padding-bottom: 1rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.activity-item:last-child {
  border-bottom: none;
  padding-bottom: 0;
}

.activity-time {
  width: 100px;
  color: var(--text-light);
  font-size: 0.875rem;
}

.activity-details h4 {
  font-size: 1rem;
  font-weight: 500;
  margin-bottom: 0.25rem;
}

.activity-details p {
  color: var(--text-light);
  font-size: 0.875rem;
}

.system-status {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.status-item {
  display: flex;
  align-items: center;
}

.status-label {
  width: 120px;
  font-size: 0.875rem;
  color: var(--text-color);
}

.status-bar {
  flex: 1;
  height: 8px;
  background-color: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
  overflow: hidden;
  margin: 0 1rem;
}

.status-progress {
  height: 100%;
  background-color: var(--primary-color);
  border-radius: 4px;
}

.status-value {
  width: 50px;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-color);
  text-align: right;
}

@media (max-width: 768px) {
  .dashboard-content {
    grid-template-columns: 1fr;
  }

  .activity-item {
    flex-direction: column;
  }

  .activity-time {
    margin-bottom: 0.5rem;
  }

  .status-item {
    flex-wrap: wrap;
  }

  .status-label {
    width: 100%;
    margin-bottom: 0.25rem;
  }

  .status-bar {
    margin: 0 0.5rem;
  }
}
</style>