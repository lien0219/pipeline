/**
 * 基础响应类型
 */
export interface BaseResponse<T> {
  code: number;
  data: T;
  msg: string;
}

/**
 * 分页结果类型
 */
export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}

/**
 * API错误类型
 */
export interface ApiError {
  message: string;
  code: number;
  data?: any;
}

export interface TemplateCategory {
  id: number;
  name: string;
  description: string;
  icon: string;
  order: number;
  creator_id: number;
  created_at: string;
  updated_at: string;
}

export interface Template {
  id: number;
  name: string;
  description: string;
  category_id: number;
  icon: string;
  tags: string;
  is_public: boolean;
  download_count: number;
  creator_id: number;
  created_at: string;
  updated_at: string;
  category?: TemplateCategory;
  creator?: any;
  versions?: TemplateVersion[];
}

export interface TemplateVersion {
  id: number;
  template_id: number;
  version: string;
  content: string;
  changelog: string;
  is_latest: boolean;
  download_count: number;
  creator_id: number;
  created_at: string;
  updated_at: string;
}
export interface TemplateVersion {
  id: number;
  template_id: number;
  version: string;
  content: string;
  changelog: string;
  download_count: number;
  is_latest: boolean;
  creator_id: number;
  created_at: string;
  updated_at: string;
}
