<template>
  <div class="home-page">
    <section class="hero">
      <div class="hero-content">
        <h1>企业级应用解决方案</h1>
        <p>基于Gin+GORM+Redis+MySQL+Vue3的全栈应用框架</p>
        <div class="hero-buttons">
          <router-link to="/auth/register" class="btn btn-primary">立即开始</router-link>
          <a href="#features" class="btn btn-secondary">了解更多</a>
        </div>
      </div>
    </section>

    <section id="features" class="features">
      <h2 class="section-title">核心特性</h2>
      <div class="feature-grid">
        <div class="feature-card">
          <div class="feature-icon">🚀</div>
          <h3>高性能后端</h3>
          <p>基于Gin框架构建的高性能API服务，支持中间件、路由分组等特性</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">🔄</div>
          <h3>数据持久化</h3>
          <p>使用GORM和MySQL实现强大的数据持久化层，支持复杂查询和事务</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">⚡</div>
          <h3>缓存支持</h3>
          <p>集成Redis缓存，提高应用响应速度和并发处理能力</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">📱</div>
          <h3>现代化前端</h3>
          <p>基于Vue 3和组合式API构建的响应式用户界面</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">🔒</div>
          <h3>安全认证</h3>
          <p>完整的JWT认证系统，保护API安全</p>
        </div>
        <div class="feature-card">
          <div class="feature-icon">📚</div>
          <h3>API文档</h3>
          <p>集成Swagger文档，方便API调试和协作开发</p>
        </div>
      </div>
    </section>

    <section class="cta">
      <h2>准备好开始构建了吗？</h2>
      <p>立即注册并体验企业级应用框架的强大功能</p>
      <router-link to="/auth/register" class="btn btn-primary">免费注册</router-link>
    </section>
  </div>
</template>

<script setup>
// 首页不需要特殊逻辑
</script>

<style scoped>
.home-page {
  padding-bottom: 2rem;
}

.hero {
  padding: 4rem 0;
  text-align: center;
  background-color: var(--surface-color);
  border-radius: 8px;
  margin-bottom: 3rem;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 1rem;
}

.hero h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: var(--primary-color);
}

.hero p {
  font-size: 1.25rem;
  color: var(--text-light);
  margin-bottom: 2rem;
}

.hero-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.section-title {
  text-align: center;
  font-size: 2rem;
  margin-bottom: 2rem;
  color: var(--text-color);
}

.features {
  margin-bottom: 3rem;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.feature-card {
  background-color: var(--surface-color);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s, box-shadow 0.3s;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.feature-card h3 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  color: var(--primary-color);
}

.feature-card p {
  color: var(--text-light);
}

.cta {
  background-color: var(--surface-color);
  border-radius: 8px;
  padding: 3rem 2rem;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.cta h2 {
  font-size: 2rem;
  margin-bottom: 1rem;
  color: var(--primary-color);
}

.cta p {
  font-size: 1.25rem;
  color: var(--text-light);
  margin-bottom: 2rem;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
}

@media (max-width: 768px) {
  .hero h1 {
    font-size: 2rem;
  }

  .hero p {
    font-size: 1rem;
  }

  .feature-grid {
    grid-template-columns: 1fr;
  }

  .cta h2 {
    font-size: 1.5rem;
  }

  .cta p {
    font-size: 1rem;
  }
}
</style>