import axios from 'axios';

const API_URL = process.env.REACT_APP_API_URL ?? 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
    'ngrok-skip-browser-warning': '1',
  },
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('user');
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const aiAPI = {
  generateBackground: () => api.post('/api/ai/background'),
  generateProjects: () => api.post('/api/ai/projects'),
  generateSkills: () => api.post('/api/ai/skills'),
  generateDaily: () => api.post('/api/ai/daily'),
}

export const noteAPI = {
  getNotes: (category) =>
    api.get('/api/notes', { params: category ? { category } : {} }),
  getNote: (id) => api.get(`/api/notes/${id}`),
  createNote: (data) => api.post('/api/notes', data),
  updateNote: (id, data) => api.put(`/api/notes/${id}`, data),
  deleteNote: (id) => api.delete(`/api/notes/${id}`),
}

export const userAPI = {
  register: (userData) => api.post('/api/users/register', userData),
  login: (credentials) => api.post('/api/users/login', credentials),
  logout: () => api.post('/api/users/logout'),
  getUsers: () => api.get('/api/users'),
  getUser: (id) => api.get(`/api/users/${id}`),
  updateUser: (id, userData) => api.put(`/api/users/${id}`, userData),
  deleteUser: (id) => api.delete(`/api/users/${id}`),
  updateUserRole: (id, role) => api.patch(`/api/users/${id}/role`, { role }),
};

export default api;
