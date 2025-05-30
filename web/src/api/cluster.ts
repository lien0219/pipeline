import http from "./http";

export const getAllClusters = () => {
  return http.get("/api/v1/clusters");
};

export const addCluster = (data: {
  name: string;
  kubeconfig: string;
  description?: string;
}) => {
  return http.post("/api/v1/clusters", data);
};
