import { create } from 'zustand';
import { notificationsAPI } from '../api/notifications';

export const useNotificationStore = create((set, get) => ({
  // State
  notifications: [],
  unreadCount: 0,
  isLoading: false,
  error: null,
  paging: { page: 1, per_page: 20, has_next_page: false },

  // Actions
  fetchNotifications: async (params = {}) => {
    set({ isLoading: true });
    try {
      const response = await notificationsAPI.list(params);
      if (response.success) {
        // If loading page 1, replace list, otherwise append (if infinite scroll)
        // For now, let's just replace or handle paging externally
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
      }
    } catch (error) {
      set({ isLoading: false, error: error.message });
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
    set((state) => ({
      notifications: state.notifications.map((n) => (n.id === id ? { ...n, is_read: true } : n)),
      unreadCount: Math.max(0, state.unreadCount - 1),
    }));

    try {
      await notificationsAPI.markAsRead(id);
    } catch (error) {
      // Revert on failure? For read status it's usually fine to ignore or refetch
      console.error('Failed to mark as read', error);
    }
  },

  markAllAsRead: async () => {
    // Optimistic update
    set((state) => ({
      notifications: state.notifications.map((n) => ({ ...n, is_read: true })),
      unreadCount: 0,
    }));

    try {
      await notificationsAPI.markAllAsRead();
    } catch (error) {
      console.error('Failed to mark all as read', error);
    }
  },

  deleteNotification: async (id) => {
    const notif = get().notifications.find((n) => n.id === id);
    const wasUnread = notif && !notif.is_read;

    // Optimistic update
    set((state) => ({
      notifications: state.notifications.filter((n) => n.id !== id),
      unreadCount: wasUnread ? Math.max(0, state.unreadCount - 1) : state.unreadCount,
    }));

    try {
      await notificationsAPI.delete(id);
    } catch (error) {
      console.error('Failed to delete notification', error);
    }
  },
}));
