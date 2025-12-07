import api from './api';

// Get current cluster metrics
export const getMetrics = async () => {
  const response = await api.get('/metrics');
  return response.data;
};

// Health check
export const checkHealth = async () => {
  const response = await api.get('/health');
  return response.data;
};
