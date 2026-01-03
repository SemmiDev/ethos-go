import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// Theme constants
export const THEMES = {
  LIGHT: 'banking',
  DARK: 'banking-dark',
};

// Get system preference
const getSystemTheme = () => {
  if (typeof window !== 'undefined' && window.matchMedia) {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? THEMES.DARK : THEMES.LIGHT;
  }
  return THEMES.LIGHT;
};

// Apply theme to document
const applyTheme = (theme) => {
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', theme);
    // Also update meta theme-color for mobile browsers
    const metaThemeColor = document.querySelector('meta[name="theme-color"]');
    if (metaThemeColor) {
      metaThemeColor.setAttribute('content', theme === THEMES.DARK ? '#0F172A' : '#FFFFFF');
    }
  }
};

export const useThemeStore = create(
  persist(
    (set, get) => ({
      // State
      theme: THEMES.LIGHT,
      isSystemPreference: true, // Whether to follow system preference

      // Actions
      setTheme: (theme) => {
        set({ theme, isSystemPreference: false });
        applyTheme(theme);
      },

      toggleTheme: () => {
        const { theme } = get();
        const newTheme = theme === THEMES.LIGHT ? THEMES.DARK : THEMES.LIGHT;
        set({ theme: newTheme, isSystemPreference: false });
        applyTheme(newTheme);
      },

      setSystemPreference: (enabled) => {
        if (enabled) {
          const systemTheme = getSystemTheme();
          set({ theme: systemTheme, isSystemPreference: true });
          applyTheme(systemTheme);
        } else {
          set({ isSystemPreference: false });
        }
      },

      // Initialize theme (call on app mount)
      initializeTheme: () => {
        const { theme, isSystemPreference } = get();
        if (isSystemPreference) {
          const systemTheme = getSystemTheme();
          set({ theme: systemTheme });
          applyTheme(systemTheme);
        } else {
          applyTheme(theme);
        }
      },

      // Helpers
      isDarkMode: () => get().theme === THEMES.DARK,
    }),
    {
      name: 'ethos-theme',
      partialize: (state) => ({
        theme: state.theme,
        isSystemPreference: state.isSystemPreference,
      }),
      onRehydrateStorage: () => (state) => {
        // Apply theme after rehydration
        if (state) {
          if (state.isSystemPreference) {
            const systemTheme = getSystemTheme();
            state.theme = systemTheme;
            applyTheme(systemTheme);
          } else {
            applyTheme(state.theme);
          }
        }
      },
    }
  )
);

// Listen for system theme changes
if (typeof window !== 'undefined' && window.matchMedia) {
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', (e) => {
    const store = useThemeStore.getState();
    if (store.isSystemPreference) {
      const newTheme = e.matches ? THEMES.DARK : THEMES.LIGHT;
      store.setTheme(newTheme);
      // Keep system preference flag
      useThemeStore.setState({ isSystemPreference: true });
    }
  });
}
