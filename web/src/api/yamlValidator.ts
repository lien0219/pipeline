import http from "./http";

// YAML验证请求参数
export interface ValidateYAMLRequest {
  name: string;
  content: string;
  schema_type: string;
  schema_name?: string;
}

// YAML验证结果
export interface ValidateYAMLResult {
  id: number;
  name: string;
  content: string;
  schema_type: string;
  schema_name?: string;
  is_valid: boolean;
  errors: string;
  creator_id: number;
  created_at: string;
  updated_at: string;
}

// Schema请求参数
export interface CreateYAMLSchemaRequest {
  name: string;
  type: string;
  version: string;
  schema: string;
  description?: string;
}

export interface UpdateYAMLSchemaRequest {
  name: string;
  version: string;
  schema: string;
  description?: string;
}

// Schema模型
export interface YAMLSchema {
  id: number;
  name: string;
  type: string;
  version: string;
  schema: string;
  description?: string;
  creator_id: number;
  created_at: string;
  updated_at: string;
}

// 验证YAML
export function validateYAML(data: ValidateYAMLRequest) {
  return http.post<ValidateYAMLResult>("/v1/yaml/validate", data);
}

// 获取验证历史
export function getValidationHistory(params: {
  page?: number;
  pageSize?: number;
  name?: string;
}) {
  const defaultParams = {
    page: 1,
    pageSize: 10,
    ...params,
  };
  return http.get<PageResponse<ValidateYAMLResult>>("/v1/yaml/history", {
    params: defaultParams,
  });
}

// 创建Schema
export function createYAMLSchema(data: CreateYAMLSchemaRequest) {
  return http.post<YAMLSchema>("/v1/yaml/schema", data);
}

export interface PageResponse<T> {
  data: {
    list: T[];
    total: number;
    page: number;
    pageSize: number;
  };
  code: number;
  message: string;
}

// 获取Schema列表
export function getYAMLSchemas(params: {
  type?: string;
  page?: number;
  pageSize?: number;
}) {
  const defaultParams = {
    page: 1,
    pageSize: 10,
    ...params,
  };
  return http.get<PageResponse<YAMLSchema>>("/v1/yaml/schema", {
    params: defaultParams,
  });
}

// 获取Schema详情
export function getYAMLSchemaById(id: number) {
  return http.get<YAMLSchema>(`/v1/yaml/schema/${id}`);
}

// 更新Schema
export function updateYAMLSchema(id: number, data: UpdateYAMLSchemaRequest) {
  return http.put<YAMLSchema>(`/v1/yaml/schema/${id}`, data);
}

// 删除Schema
export function deleteYAMLSchema(id: number) {
  return http.delete(`/v1/yaml/schema/${id}`);
}
