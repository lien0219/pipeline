import http from "./http";

// 配置类型定义
export interface SettingItem {
  [key: string]: any;
}

export interface SaveSettingsParams {
  type: "basic" | "email" | "integration" | "system";
  items: SettingItem;
}

// 保存配置
export function saveSettings(data: SaveSettingsParams) {
  return http.post("/v1/settings", data);
}

// 获取配置
export function getSettings(type?: string) {
  return http.get("/v1/settings", { params: { type } });
}
