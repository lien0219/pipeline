import http from "./http"

export const authApi = {
    login(data) {
        return http.post("/auth/login", data)
    },

    logout() {
        return http.post("/auth/logout")
    },

    getUserInfo() {
        return http.get("/auth/user")
    },

    updateUserInfo(data) {
        return http.put("/auth/user", data)
    },

    changePassword(data) {
        return http.put("/auth/password", data)
    },
}
