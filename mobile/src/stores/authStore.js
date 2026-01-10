import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { authAPI } from '../api/auth';

export const useAuthStore = create((set, get) => ({
  // State
  user: null,
  sessions: [],
  accessToken: null,
  sessionId: null,
  isAuthenticated: false,
  isLoading: true,
  isSessionsLoading: false,
  error: null,

  // Initialize auth state from secure storage
  initialize: async () => {
    try {
      const token = await SecureStore.getItemAsync('accessToken');
      const sessionId = await SecureStore.getItemAsync('sessionId');
      const userStr = await SecureStore.getItemAsync('user');

      if (token && userStr) {
        const user = JSON.parse(userStr);
        set({
          accessToken: token,
          sessionId,
          user,
          isAuthenticated: true,
          isLoading: false,
        });
      } else {
        set({ isLoading: false });
      }
    } catch (error) {
      console.error('Failed to initialize auth:', error);
      set({ isLoading: false });
    }
  },

  // Login
  login: async (email, password) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authAPI.login({ email, password });
      const data = response.data;

      console.log('Login response:', data); // Debug logging

      if (!data.access_token || !data.session_id) {
        throw new Error('Invalid response from server: Missing tokens');
      }

      // Store tokens securely (ensure strings)
      await SecureStore.setItemAsync('accessToken', String(data.access_token));
      await SecureStore.setItemAsync('sessionId', String(data.session_id));

      // Fetch user profile after successful login
      const profileResponse = await authAPI.getProfile();
      const user = profileResponse?.data || {};
      await SecureStore.setItemAsync('user', JSON.stringify(user));

      set({
        accessToken: data.access_token,
        sessionId: data.session_id,
        user: user,
        isAuthenticated: true,
        isLoading: false,
      });

      return { success: true };
    } catch (error) {
      console.error('Login error full:', error); // Debug logging
      console.error('Login response data:', error.response?.data);
      console.error('Login status:', error.response?.status);

      const message = error.response?.data?.message || error.message || 'Login failed';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Register
  register: async (name, email, password, timezone) => {
    set({ isLoading: true, error: null });
    try {
      await authAPI.register({ name, email, password, timezone });
      set({ isLoading: false });
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Registration failed';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Logout
  logout: async () => {
    try {
      const sessionId = get().sessionId;
      if (sessionId) {
        await authAPI.logout(sessionId);
      }
    } catch {
      // Ignore errors on logout
    }

    // Clear secure storage
    await SecureStore.deleteItemAsync('accessToken');
    await SecureStore.deleteItemAsync('sessionId');
    await SecureStore.deleteItemAsync('user');

    set({
      user: null,
      accessToken: null,
      sessionId: null,
      isAuthenticated: false,
    });
  },

  // Get profile
  getProfile: async () => {
    try {
      const response = await authAPI.getProfile();
      set({ user: response.data });
      await SecureStore.setItemAsync('user', JSON.stringify(response.data));
    } catch (error) {
      console.error('Failed to get profile:', error);
    }
  },

  // Update profile
  updateProfile: async (data) => {
    set({ isLoading: true, error: null });
    try {
      const response = await authAPI.updateProfile(data);
      const updatedUser = response.data;

      set({ user: updatedUser, isLoading: false });
      await SecureStore.setItemAsync('user', JSON.stringify(updatedUser));

      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to update profile';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Change password
  changePassword: async ({ current_password, new_password }) => {
    set({ isLoading: true, error: null });
    try {
      await authAPI.changePassword({ current_password, new_password });
      set({ isLoading: false });
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to change password';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Delete account
  deleteAccount: async (password) => {
    set({ isLoading: true, error: null });
    try {
      await authAPI.deleteAccount(password);
      await get().logout(); // Auto logout after deletion
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to delete account';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Session Management
  // Session Management
  fetchSessions: async () => {
    set({ isSessionsLoading: true, error: null });
    try {
      const response = await authAPI.getSessions();
      console.log('[authStore] fetchSessions raw response:', JSON.stringify(response, null, 2));

      // response is the response body object { data: [...], meta: ... }
      // So response.data is the sessions array
      const sessionsList = response.data;
      const safeSessions = Array.isArray(sessionsList) ? sessionsList : [];
      console.log('[authStore] Sessions set to store:', safeSessions.length);
      set({ sessions: safeSessions, isSessionsLoading: false });
    } catch (error) {
      console.error('Fetch sessions error:', error);
      const message = error.response?.data?.message || 'Failed to fetch sessions';
      set({ error: message, isSessionsLoading: false });
    }
  },

  revokeSession: async (sessionId) => {
    set({ isSessionsLoading: true, error: null });
    try {
      await authAPI.revokeSession(sessionId);
      // Refresh user sessions
      await get().fetchSessions();
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to revoke session';
      set({ error: message, isSessionsLoading: false });
      return { success: false, error: message };
    }
  },

  revokeOtherSessions: async () => {
    set({ isSessionsLoading: true, error: null });
    try {
      const response = await authAPI.revokeOtherSessions();
      await get().fetchSessions();
      return { success: true, count: response.data?.count };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to revoke other sessions';
      set({ error: message, isSessionsLoading: false });
      return { success: false, error: message };
    }
  },

  // Clear error
  clearError: () => set({ error: null }),
}));
