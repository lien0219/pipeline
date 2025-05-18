import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import router from './router';
import './style.css'
import { setupAxiosInterceptors } from './api/axios';

// 设置axios拦截器
setupAxiosInterceptors();

// 创建应用实例
const app = createApp(App);

// 使用Pinia状态管理
app.use(createPinia());

// 使用路由
app.use(router);

// 挂载应用
app.mount('#app');
