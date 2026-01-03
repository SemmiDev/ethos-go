import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authAPI } from '../api/auth';

export const useAuthStore = create(
  persist(
    (set, get) => ({
      // State
      user: null,
      accessToken: null,
      refreshToken: null,
      sessionId: null,
      sessions: [],
      isAuthenticated: false,
      isLoading: false,
      error: null,

      // Actions
      register: async (data) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.register(data);
          if (response.success) {
            set({ isLoading: false });
            return { success: true, data: response.data };
          }
          throw new Error(response.message || 'Registration failed');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Registration failed';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      login: async (data) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.login(data);
          if (response.success && response.data) {
            set({
              user: { id: response.data.user_id },
              accessToken: response.data.access_token,
              refreshToken: response.data.refresh_token,
              sessionId: response.data.session_id,
              isAuthenticated: true,
              isLoading: false,
              error: null,
            });
            // Fetch profile to get full user data
            get().fetchProfile();
            return { success: true };
          }
          throw new Error(response.message || 'Login failed');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Login failed';
          const errorCode = error.response?.data?.error?.code || null;
          set({ isLoading: false, error: message });
          return { success: false, error: message, errorCode };
        }
      },

      logout: async () => {
        const { sessionId } = get();
        try {
          if (sessionId) {
            await authAPI.logout(sessionId);
          }
        } catch (error) {
          console.error('Logout error:', error);
        } finally {
          set({
            user: null,
            accessToken: null,
            refreshToken: null,
            sessionId: null,
            sessions: [],
            isAuthenticated: false,
            error: null,
          });
        }
      },

      logoutAll: async () => {
        const { user } = get();
        try {
          if (user?.id) {
            await authAPI.logoutAll(user.id);
          }
        } catch (error) {
          console.error('Logout all error:', error);
        } finally {
          set({
            user: null,
            accessToken: null,
            refreshToken: null,
            sessionId: null,
            sessions: [],
            isAuthenticated: false,
            error: null,
          });
        }
      },

      revokeSession: async (targetSessionId) => {
        set({ isLoading: true, error: null });
        try {
          await authAPI.logout(targetSessionId);

          set((state) => {
            return {
              sessions: state.sessions.filter((s) => s.session_id !== targetSessionId),
              isLoading: false,
              // If current session revoked, we should technically clear auth state,
              // but usually we let the logout() calls do that.
              // However, since we are inside revokeSession, we might want to just handle lists.
              // If it IS current session, the user basically logged out.
            };
          });

          if (targetSessionId === get().sessionId) {
            get().logout(); // Perform full local logout
          }

          return { success: true };
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to revoke session';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      revokeOtherSessions: async () => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.revokeAllOtherSessions();
          if (response.success) {
            get().fetchSessions(); // Refresh list
            set({ isLoading: false });
            return { success: true, count: response.data?.revoked_count };
          }
          throw new Error(response.message || 'Failed to revoke other sessions');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to revoke other sessions';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      deleteAccount: async (password) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.deleteAccount(password);
          if (response.success) {
            get().logout(); // Logout after deletion
            return { success: true };
          }
          throw new Error(response.message || 'Failed to delete account');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to delete account';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      exportData: async () => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.exportUserData();
          // Note: The response should contain the JSON data directly or a success wrapper with data
          if (response.success) {
            set({ isLoading: false });
            return { success: true, data: response.data };
          }
          throw new Error(response.message || 'Failed to export data');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to export data';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      fetchSessions: async (params) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.getSessions(params);
          if (response.success) {
            set({ sessions: response.data || [], isLoading: false });
            return { success: true, data: response.data };
          }
          throw new Error(response.message || 'Failed to fetch sessions');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to fetch sessions';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      fetchProfile: async () => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.getProfile();
          if (response.success) {
            set({
              user: {
                id: response.data.user_id,
                name: response.data.name,
                email: response.data.email,
                timezone: response.data.timezone,
                createdAt: response.data.created_at,
              },
              isLoading: false,
            });
            return { success: true, data: response.data };
          }
          throw new Error(response.message || 'Failed to fetch profile');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to fetch profile';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      updateProfile: async (data) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.updateProfile(data);
          if (response.success) {
            set({
              user: {
                id: response.data.user_id,
                name: response.data.name,
                email: response.data.email,
                timezone: response.data.timezone,
                createdAt: response.data.created_at,
              },
              isLoading: false,
            });
            return { success: true, data: response.data };
          }
          throw new Error(response.message || 'Failed to update profile');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to update profile';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      changePassword: async (data) => {
        set({ isLoading: true, error: null });
        try {
          const response = await authAPI.changePassword(data);
          if (response.success) {
            set({ isLoading: false });
            return { success: true };
          }
          throw new Error(response.message || 'Failed to change password');
        } catch (error) {
          const message = error.response?.data?.message || error.message || 'Failed to change password';
          set({ isLoading: false, error: message });
          return { success: false, error: message };
        }
      },

      setAuth: (accessToken, refreshToken, userId, sessionId) => {
        set({
          accessToken,
          refreshToken,
          sessionId,
          user: { id: userId },
          isAuthenticated: true,
          isLoading: false,
          error: null,
        });
        get().fetchProfile();
      },

      clearError: () => set({ error: null }),
    }),
    {
      name: 'ethos-auth',
      partialize: (state) => ({
        user: state.user,
        accessToken: state.accessToken,
        refreshToken: state.refreshToken,
        sessionId: state.sessionId,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);
