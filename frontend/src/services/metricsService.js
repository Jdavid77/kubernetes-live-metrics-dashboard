import api from './api';

// Get current cluster metrics
export const getMetrics = async () => {
  const response = await api.get('/api/metrics');
  return response.data;
};

// Health check
export const checkHealth = async () => {
  const response = await api.get('/api/health');
  return response.data;
};
