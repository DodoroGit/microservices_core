import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const aiAPI = {
  sendMessage: (message, userId) =>
    api.post('/api/ai/chat', { message, user_id: userId }),
  getHistory: (userId) =>
    api.get('/api/ai/chat/history', { params: { user_id: userId } }),
  clearHistory: (userId) =>
    api.delete('/api/ai/chat/history', { params: { user_id: userId } }),
}

export const userAPI = {
  register: (userData) => api.post('/api/users/register', userData),
  login: (credentials) => api.post('/api/users/login', credentials),
  getUsers: () => api.get('/api/users'),
  getUser: (id) => api.get(`/api/users/${id}`),
  updateUser: (id, userData) => api.put(`/api/users/${id}`, userData),
  deleteUser: (id) => api.delete(`/api/users/${id}`),
  updateUserRole: (id, role) => api.patch(`/api/users/${id}/role`, { role }),
};

export default api;
