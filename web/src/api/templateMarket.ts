import http from "./http";
import {
  type TemplateCategory,
  type Template,
  type TemplateVersion,
} from "./types";

// 分类管理
export const createTemplateCategory = (data: any) =>
  http.post("/v1/template-market/category", data);
export const getTemplateCategories = (params?: {
  page?: number;
  pageSize?: number;
  keyword?: string;
}) =>
  http.get<{
    data: {
      list: TemplateCategory[];
      total: number;
      page: number;
      pageSize: number;
    };
  }>("/v1/template-market/category", { params });
export const updateTemplateCategory = (id: number, data: any) =>
  http.put(`/v1/template-market/category/${id}`, data);
export const deleteTemplateCategory = (id: number) =>
  http.delete(`/v1/template-market/category/${id}`);

// 模板管理
export const createTemplate = (data: any) =>
  http.post("/v1/template-market/template", data);
export const getTemplates = (params?: {
  category_id?: number;
  public?: boolean;
  page?: number;
  pageSize?: number;
}) =>
  http.get<{
    data: Template[];
    total: number;
    page: number;
    pageSize: number;
  }>("/v1/template-market/template", { params });
export const getTemplateById = (id: number) =>
  http.get<{ data: Template }>(`/v1/template-market/template/${id}`);
export const updateTemplate = (id: number, data: any) =>
  http.put(`/v1/template-market/template/${id}`, data);
export const deleteTemplate = (id: number) =>
  http.delete(`/v1/template-market/template/${id}`);

// 版本管理
export const createTemplateVersion = (templateId: number, data: any) =>
  http.post(`/v1/template-market/template/${templateId}/version`, data);
export const getTemplateVersions = (templateId: number) =>
  http.get<{ data: TemplateVersion[] }>(
    `/v1/template-market/template/${templateId}/version`
  );
export const getTemplateVersionById = (templateId: number, versionId: number) =>
  http.get<{ data: TemplateVersion }>(
    `/v1/template-market/template/${templateId}/version/${versionId}`
  );
export const deleteTemplateVersion = (templateId: number, versionId: number) =>
  http.delete(
    `/v1/template-market/template/${templateId}/version/${versionId}`
  );
export const setVersionAsLatest = (templateId: number, versionId: number) =>
  http.post(
    `/v1/template-market/template/${templateId}/version/${versionId}/latest`
  );

// 搜索和下载
export const searchTemplates = (params?: {
  keyword?: string;
  category_id?: number;
  tags?: string;
}) => http.get<{ data: Template[] }>("/v1/template-market/search", { params });
export const downloadTemplate = (id: number) =>
  http.get<{ data: TemplateVersion }>(
    `/v1/template-market/template/${id}/download`
  );
