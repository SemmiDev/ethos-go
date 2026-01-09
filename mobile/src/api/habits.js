import apiClient from './client';

export const habitsAPI = {
  // List all habits
  list: async (params = {}) => {
    const response = await apiClient.get('/habits', { params });
    return response.data;
  },

  // Get a single habit
  get: async (habitId) => {
    const response = await apiClient.get(`/habits/${habitId}`);
    return response.data;
  },

  // Create a new habit
  create: async (data) => {
    const response = await apiClient.post('/habits', data);
    return response.data;
  },

  // Update a habit
  update: async (habitId, data) => {
    const response = await apiClient.put(`/habits/${habitId}`, data);
    return response.data;
  },

  // Delete a habit
  delete: async (habitId) => {
    const response = await apiClient.delete(`/habits/${habitId}`);
    return response.data;
  },

  // Activate a habit
  activate: async (habitId) => {
    const response = await apiClient.post(`/habits/${habitId}/activate`);
    return response.data;
  },

  // Deactivate a habit
  deactivate: async (habitId) => {
    const response = await apiClient.post(`/habits/${habitId}/deactivate`);
    return response.data;
  },

  // Get habit statistics
  getStats: async (habitId) => {
    const response = await apiClient.get(`/habits/${habitId}/stats`);
    return response.data;
  },

  // Log a habit
  log: async (habitId, data) => {
    const response = await apiClient.post(`/habits/${habitId}/logs`, data);
    return response.data;
  },

  // Get habit logs
  getLogs: async (habitId, params = {}) => {
    const response = await apiClient.get(`/habits/${habitId}/logs`, { params });
    return response.data;
  },

  // Update a habit log
  updateLog: async (logId, data) => {
    const response = await apiClient.put(`/habit-logs/${logId}`, data);
    return response.data;
  },

  // Delete a habit log
  deleteLog: async (logId) => {
    const response = await apiClient.delete(`/habit-logs/${logId}`);
    return response.data;
  },

  // Get dashboard data
  getDashboard: async () => {
    const response = await apiClient.get('/dashboard');
    return response.data;
  },

  // Get weekly analytics data
  getWeeklyAnalytics: async () => {
    const response = await apiClient.get('/analytics/weekly');
    return response.data;
  },
};
