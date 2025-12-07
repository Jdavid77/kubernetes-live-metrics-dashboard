import api from './api';

// Get all services
export const getServices = async (namespace = 'all') => {
  const params = namespace && namespace !== 'all' ? { namespace } : {};
  const response = await api.get('/services', { params });
  return response.data;
};
