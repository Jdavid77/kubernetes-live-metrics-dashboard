import api from './api';

// Get all applications
export const getApplications = async (namespace = 'all') => {
  const params = namespace && namespace !== 'all' ? { namespace } : {};
  const response = await api.get('/applications', { params });
  return response.data;
};

// Get application detail
export const getApplicationDetail = async (namespace, name) => {
  const response = await api.get(`/applications/${name}`, {
    params: { namespace },
  });
  return response.data;
};
