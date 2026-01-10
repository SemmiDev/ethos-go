import { create } from 'zustand';
import { notificationsAPI } from '../api/notifications';

export const useNotificationsStore = create((set, get) => ({
  // State
  notifications: [],
  unreadCount: 0,
  isLoading: false,
  error: null,
  paging: { page: 1, per_page: 20, has_next_page: false },

  // Actions
  fetchNotifications: async (params = {}) => {
    set({ isLoading: true, error: null });
    try {
      const response = await notificationsAPI.list(params);
      if (response.success) {
        const isFirstPage = params.page === 1 || !params.page;

        set((state) => ({
          notifications: isFirstPage ? response.data : [...state.notifications, ...response.data],
          paging: {
            page: response.meta?.pagination?.current_page || 1,
            per_page: response.meta?.pagination?.per_page || 20,
            has_next_page: response.meta?.pagination?.has_next_page || false,
          },
          isLoading: false,
        }));
      } else {
        set({ isLoading: false, error: response.message || 'Failed to fetch notifications' });
      }
    } catch (error) {
      const message = error.response?.data?.message || error.message || 'Failed to fetch notifications';
      set({ isLoading: false, error: message });
    }
  },

  fetchUnreadCount: async () => {
    try {
      const response = await notificationsAPI.getUnreadCount();
      if (response.success) {
        set({ unreadCount: response.data.count });
      }
    } catch (error) {
      console.error('Failed to fetch unread count', error);
    }
  },

  markAsRead: async (id) => {
    // Optimistic update
    const notif = get().notifications.find((n) => n.id === id);
    const wasUnread = notif && !notif.is_read;

    set((state) => ({
      notifications: state.notifications.map((n) => (n.id === id ? { ...n, is_read: true } : n)),
      unreadCount: wasUnread ? Math.max(0, state.unreadCount - 1) : state.unreadCount,
    }));

    try {
      await notificationsAPI.markAsRead(id);
      return { success: true };
    } catch (error) {
      // Revert on failure
      if (wasUnread) {
        set((state) => ({
          notifications: state.notifications.map((n) => (n.id === id ? { ...n, is_read: false } : n)),
          unreadCount: state.unreadCount + 1,
        }));
      }
      const message = error.response?.data?.message || 'Failed to mark as read';
      return { success: false, error: message };
    }
  },

  markAllAsRead: async () => {
    const previousState = get().notifications;
    const previousCount = get().unreadCount;

    // Optimistic update
    set((state) => ({
      notifications: state.notifications.map((n) => ({ ...n, is_read: true })),
      unreadCount: 0,
    }));

    try {
      await notificationsAPI.markAllAsRead();
      return { success: true };
    } catch (error) {
      // Revert on failure
      set({ notifications: previousState, unreadCount: previousCount });
      const message = error.response?.data?.message || 'Failed to mark all as read';
      return { success: false, error: message };
    }
  },

  deleteNotification: async (id) => {
    const notif = get().notifications.find((n) => n.id === id);
    const wasUnread = notif && !notif.is_read;
    const previousNotifications = get().notifications;
    const previousCount = get().unreadCount;

    // Optimistic update
    set((state) => ({
      notifications: state.notifications.filter((n) => n.id !== id),
      unreadCount: wasUnread ? Math.max(0, state.unreadCount - 1) : state.unreadCount,
    }));

    try {
      await notificationsAPI.delete(id);
      return { success: true };
    } catch (error) {
      // Revert on failure
      set({ notifications: previousNotifications, unreadCount: previousCount });
      const message = error.response?.data?.message || 'Failed to delete notification';
      return { success: false, error: message };
    }
  },

  // Clear error
  clearError: () => set({ error: null }),
}));
