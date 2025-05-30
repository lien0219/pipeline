import http from "./http";

export const getHPAPolicies = (params?: { page?: number; limit?: number }) => {
  return http.get("/api/v1/hpa", { params });
};

export const createHPAPolicy = (data: {
  name: string;
  namespace: string;
  minReplicas: number;
  maxReplicas: number;
  metrics: Array<{
    type: string;
    value: number;
  }>;
}) => {
  return http.post("/api/v1/hpa", data);
};
