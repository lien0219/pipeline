import http from "./http";

// 系统状态接口定义
export interface SystemStatus {
  uptime: string;
  cpu_usage: string;
  memory_usage: string;
  disk_usage: string;
  db_connections: number;
  active_users: number;
}

// 日志文件接口定义
export interface LogFile {
  name: string;
  size: number;
  modified: string;
}

// 获取系统状态
export function getSystemStatus() {
  return http.get("/v1/system/status");
}

// 获取日志文件列表
export function getLogFiles() {
  return http.get("/v1/system/logs");
}
// 获取日志文件内容
export const getLogContent = (
  name: string,
  startLine = 0,
  limit = 100
): Promise<any> => {
  return http({
    url: `/v1/system/logs/${name}`,
    method: "get",
    params: {
      startLine,
      limit,
    },
  });
};
// 下载日志文件
export const downloadLog = (name: string): Promise<Blob> => {
  return http({
    url: `/v1/system/logs/${name}/download`,
    method: "get",
    responseType: "blob",
    params: {
      t: new Date().getTime(),
    },
  });
};
