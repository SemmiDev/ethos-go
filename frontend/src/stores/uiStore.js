import { create } from 'zustand';

export const useUIStore = create((set) => ({
    // Toast notifications
    toasts: [],

    // Modals
    isCreateHabitModalOpen: false,
    isEditHabitModalOpen: false,
    isDeleteConfirmModalOpen: false,
    isLogHabitModalOpen: false,

    // Sidebar
    isSidebarOpen: true,
    isMobileSidebarOpen: false,

    // Actions
    addToast: (toast) => {
        const id = Date.now();
        set((state) => ({
            toasts: [...state.toasts, { ...toast, id }],
        }));

        // Auto remove after duration
        setTimeout(() => {
            set((state) => ({
                toasts: state.toasts.filter((t) => t.id !== id),
            }));
        }, toast.duration || 4000);
    },

    removeToast: (id) => {
        set((state) => ({
            toasts: state.toasts.filter((t) => t.id !== id),
        }));
    },

    // Modal actions
    openCreateHabitModal: () => set({ isCreateHabitModalOpen: true }),
    closeCreateHabitModal: () => set({ isCreateHabitModalOpen: false }),

    openEditHabitModal: () => set({ isEditHabitModalOpen: true }),
    closeEditHabitModal: () => set({ isEditHabitModalOpen: false }),

    openDeleteConfirmModal: () => set({ isDeleteConfirmModalOpen: true }),
    closeDeleteConfirmModal: () => set({ isDeleteConfirmModalOpen: false }),

    openLogHabitModal: () => set({ isLogHabitModalOpen: true }),
    closeLogHabitModal: () => set({ isLogHabitModalOpen: false }),

    // Sidebar actions
    toggleSidebar: () => set((state) => ({ isSidebarOpen: !state.isSidebarOpen })),
    toggleMobileSidebar: () => set((state) => ({ isMobileSidebarOpen: !state.isMobileSidebarOpen })),
    closeMobileSidebar: () => set({ isMobileSidebarOpen: false }),
}));
