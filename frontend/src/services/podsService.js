import api from './api';

// Get all pods
export const getPods = async (namespace = 'all') => {
  const params = namespace && namespace !== 'all' ? { namespace } : {};
  const response = await api.get('/pods', { params });
  return response.data;
};
