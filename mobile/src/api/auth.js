import apiClient from './client';

export const authAPI = {
  // Register a new user
  register: async (data) => {
    const response = await apiClient.post('/auth/register', data);
    return response.data;
  },

  // Login user
  login: async (data) => {
    const response = await apiClient.post('/auth/login', data);
    return response.data;
  },

  // Logout current session
  logout: async (sessionId) => {
    const response = await apiClient.post('/auth/logout', { session_id: sessionId });
    return response.data;
  },

  // Get user profile
  getProfile: async () => {
    const response = await apiClient.get('/auth/profile');
    return response.data;
  },

  // Update user profile
  updateProfile: async (data) => {
    const response = await apiClient.put('/auth/profile', data);
    return response.data;
  },

  // Change password
  changePassword: async (data) => {
    const response = await apiClient.post('/auth/change-password', data);
    return response.data;
  },

  // Delete account
  deleteAccount: async (password) => {
    const response = await apiClient.delete('/auth/account', { data: { password } });
    return response.data;
  },

  // Forgot password
  forgotPassword: async (data) => {
    const response = await apiClient.post('/auth/forgot-password', data);
    return response.data;
  },

  // Reset password
  resetPassword: async (data) => {
    const response = await apiClient.post('/auth/reset-password', data);
    return response.data;
  },

  // Google Auth
  getGoogleLoginURL: async () => {
    const response = await apiClient.get('/auth/google/login');
    return response.data;
  },

  googleCallback: async (code) => {
    const response = await apiClient.post('/auth/google/callback', { code });
    return response.data;
  },

  // Session Management
  getSessions: async () => {
    const response = await apiClient.get('/auth/sessions');
    return response.data;
  },

  revokeSession: async (sessionId) => {
    // Using logout endpoint to revoke a specific session
    const response = await apiClient.post('/auth/logout', { session_id: sessionId });
    return response.data;
  },

  revokeOtherSessions: async () => {
    const response = await apiClient.delete('/auth/sessions/other');
    return response.data;
  },
};
