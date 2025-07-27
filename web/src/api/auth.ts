import http from "./http";

export const authApi = {
  login(data: any) {
    return http.post("/v1/user/login", data);
  },
  getUserInfo() {
    return http.get("/v1/user/info");
  },

  updateUserInfo(data: any) {
    return http.put("/v1/user/info ", data);
  },

  changePassword(data: any) {
    return http.put("/v1/user/password", data);
  },
};
