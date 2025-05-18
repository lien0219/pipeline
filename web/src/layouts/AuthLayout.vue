<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-logo">
        <router-link to="/">
          <h1>企业应用</h1>
        </router-link>
      </div>
      <div class="auth-card">
        <router-view />
      </div>
      <div class="auth-footer">
        <p>&copy; {{ currentYear }} 企业应用. 保留所有权利.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const currentYear = computed(() => new Date().getFullYear());
</script>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--background-color);
}

.auth-container {
  width: 100%;
  max-width: 400px;
  padding: 2rem 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.auth-logo {
  margin-bottom: 2rem;
  text-align: center;
}

.auth-logo a {
  text-decoration: none;
}

.auth-logo h1 {
  color: var(--primary-color);
  font-size: 2rem;
  font-weight: 700;
}

.auth-card {
  width: 100%;
  background-color: var(--surface-color);
  border-radius: 8px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  margin-bottom: 2rem;
}

.auth-footer {
  text-align: center;
  color: var(--text-light);
  font-size: 0.875rem;
}
</style>