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
// 创建资源配额
export const createResourceQuota = (data: {
  tenant_id: string;
  cpu_quota: number;
  memory_quota: number;
  storage_quota: number;
}) => {
  return http({
    url: "/v1/resource-quota",
    method: "post",
    data,
  });
};

// 获取资源配额
export const getResourceQuota = (params: {
  tenant_id?: string;
  page?: number;
  page_size?: number;
}) => {
  return http({
    url: "/v1/resource-quota",
    method: "get",
    params,
  });
};

// 更新资源配额
export const updateResourceQuota = (
  tenantId: string,
  data: {
    cpu_quota: number;
    memory_quota: number;
    storage_quota: number;
  }
) => {
  return http({
    url: `/v1/resource-quota/${tenantId}`,
    method: "put",
    data,
  });
};

export const deleteResourceQuota = (tenantId: string) => {
  return http.delete(`/v1/resource-quota/${tenantId}`);
};
// 创建资源请求
export const createResourceRequest = (data: {
  tenant_id: string;
  cpu_request: number;
  memory_request: number;
  storage_request: number;
  reason: string;
}) => {
  return http.post("/v1/resource-requests", data);
};

// 获取资源请求列表
export const getResourceRequests = (params: {
  tenant_id?: string;
  page?: number;
  page_size?: number;
}) => {
  return http.get("/v1/resource-requests", { params });
};

// 批准资源请求
export const approveResourceRequest = (requestId: string) => {
  return http.post(`/v1/resource-requests/${requestId}/approve`);
};

// 拒绝资源请求
export const rejectResourceRequest = (
  requestId: string,
  data: {
    reason: string;
  }
) => {
  return http.post(`/v1/resource-requests/${requestId}/reject`, data);
};
// 资源报告相关接口
export const getResourceReports = (params?: any) => {
  return http({
    url: "/v1/resource-report",
    method: "get",
    params,
  });
};

export const getResourceReportById = (id: number) => {
  return http({
    url: `/v1/resource-report/${id}`,
    method: "get",
  });
};

export const createResourceReport = (data: any) => {
  return http({
    url: "/v1/resource-report",
    method: "post",
    data,
  });
};

export const updateResourceReport = (id: number, data: any) => {
  return http({
    url: `/v1/resource-report/${id}`,
    method: "put",
    data,
  });
};

export const deleteResourceReport = (id: number) => {
  return http({
    url: `/v1/resource-report/${id}`,
    method: "delete",
  });
};
