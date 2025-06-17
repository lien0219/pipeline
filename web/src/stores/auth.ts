import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { authApi } from "@/api/auth";

export const useAuthStore = defineStore("auth", () => {
  const currentUser = ref<any>(null);
  const token = ref(localStorage.getItem("token") || "");

  const isAuthenticated = computed(() => !!token.value);

  async function login(credentials: any) {
    try {
      const response = await authApi.login(credentials);
      token.value = response.data.token;
      localStorage.setItem("token", token.value);
      await fetchUserInfo();
      return response;
    } catch (error) {
      throw error;
    }
  }

  async function fetchUserInfo() {
    try {
      const response = await authApi.getUserInfo();
      currentUser.value = response.data;
      return response;
    } catch (error) {
      throw error;
    }
  }

  async function logout() {
    token.value = "";
    currentUser.value = null;
    localStorage.removeItem("token");
  }

  async function updateUserInfo(userInfo: any) {
    try {
      const response = await authApi.updateUserInfo(userInfo);
      currentUser.value = { ...currentUser.value, ...userInfo };
      return response;
    } catch (error) {
      throw error;
    }
  }

  async function changePassword(passwordData: any) {
    try {
      const response = await authApi.changePassword(passwordData);
      return response;
    } catch (error) {
      throw error;
    }
  }

  return {
    currentUser,
    token,
    isAuthenticated,
    login,
    logout,
    fetchUserInfo,
    updateUserInfo,
    changePassword,
  };
});
