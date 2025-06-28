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
