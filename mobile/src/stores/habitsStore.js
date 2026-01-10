import { create } from 'zustand';
import { habitsAPI } from '../api/habits';

export const useHabitsStore = create((set, get) => ({
  // State
  habits: [],
  dashboard: null,
  analytics: null,
  habitLogs: [],
  isLoading: false,
  error: null,

  // Fetch all habits
  fetchHabits: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await habitsAPI.list();
      // API returns wrapped response: { data: { data: [], meta: {} } }
      // We need the inner data array
      const habitsData = response.data?.data || response.data || [];
      set({ habits: habitsData, isLoading: false });
    } catch (error) {
      console.error('Fetch habits error:', error);
      const message = error.response?.data?.message || 'Failed to fetch habits';
      set({ error: message, isLoading: false });
    }
  },

  // Fetch dashboard
  fetchDashboard: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await habitsAPI.getDashboard();
      set({ dashboard: response.data, isLoading: false });
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to fetch dashboard';
      set({ error: message, isLoading: false });
    }
  },

  // Fetch weekly analytics
  fetchAnalytics: async () => {
    set({ isLoading: true, error: null });
    try {
      const response = await habitsAPI.getWeeklyAnalytics();
      set({ analytics: response.data, isLoading: false });
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to fetch analytics';
      set({ error: message, isLoading: false });
    }
  },

  // Fetch habit logs
  fetchHabitLogs: async (habitId) => {
    set({ isLoading: true, error: null });
    try {
      const response = await habitsAPI.getLogs(habitId);
      const logsData = response.data?.data || response.data || [];
      set({ habitLogs: logsData, isLoading: false });
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to fetch habit logs';
      set({ error: message, isLoading: false });
    }
  },

  // Create habit
  createHabit: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await habitsAPI.create(data);
      await get().fetchHabits();
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to create habit';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Update habit
  updateHabit: async (habitId, data) => {
    set({ isLoading: true, error: null });
    try {
      await habitsAPI.update(habitId, data);
      await get().fetchHabits();
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to update habit';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Delete habit
  deleteHabit: async (habitId) => {
    set({ isLoading: true, error: null });
    try {
      await habitsAPI.delete(habitId);
      await get().fetchHabits();
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to delete habit';
      set({ error: message, isLoading: false });
      return { success: false, error: message };
    }
  },

  // Activate habit
  activateHabit: async (habitId) => {
    try {
      await habitsAPI.activate(habitId);
      await get().fetchHabits();
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to activate habit';
      return { success: false, error: message };
    }
  },

  // Deactivate habit
  deactivateHabit: async (habitId) => {
    try {
      await habitsAPI.deactivate(habitId);
      await get().fetchHabits();
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to deactivate habit';
      return { success: false, error: message };
    }
  },

  // Log habit
  logHabit: async (habitId, { count = 1, note = '', log_date } = {}) => {
    try {
      await habitsAPI.log(habitId, {
        count,
        note: note || undefined, // Send undefined if empty string
        log_date: log_date ? new Date(log_date).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
      });
      await get().fetchDashboard();
      await get().fetchHabitLogs(habitId); // Refresh logs after logging
      return { success: true };
    } catch (error) {
      const message = error.response?.data?.message || 'Failed to log habit';
      return { success: false, error: message };
    }
  },

  // Clear error
  clearError: () => set({ error: null }),
}));
