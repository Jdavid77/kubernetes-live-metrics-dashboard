import api from './api';

// Get all nodes
export const getNodes = async () => {
  const response = await api.get('/nodes');
  return response.data;
};
