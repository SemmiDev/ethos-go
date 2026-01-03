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

  // Logout from all devices
  logoutAll: async (userId) => {
    const response = await apiClient.post('/auth/logout-all', { user_id: userId });
    return response.data;
  },

  // Revoke all other sessions
  revokeAllOtherSessions: async () => {
    const response = await apiClient.delete('/auth/sessions/other');
    return response.data;
  },

  // Get active sessions
  getSessions: async (params = {}) => {
    const response = await apiClient.get('/auth/sessions', { params });
    return response.data;
  },

  // Get user profile
  getProfile: async () => {
    const response = await apiClient.get('/auth/profile');
    return response.data;
  },

  // Export User Data
  exportUserData: async () => {
    const response = await apiClient.get('/auth/export');
    return response.data;
  },

  // Delete Account
  deleteAccount: async (password) => {
    const response = await apiClient.delete('/auth/account', { data: { password } });
    return response.data;
  },

  // Update user profile
  updateProfile: async (data) => {
    const response = await apiClient.put('/auth/profile', data);
    return response.data;
  },

  changePassword: async (data) => {
    const response = await apiClient.post('/auth/change-password', data);
    return response.data;
  },

  // Verify email
  verifyEmail: async (data) => {
    const response = await apiClient.post('/auth/verify-email', data);
    return response.data;
  },

  // Resend verification
  resendVerification: async (data) => {
    const response = await apiClient.post('/auth/resend-verification', data);
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
};
