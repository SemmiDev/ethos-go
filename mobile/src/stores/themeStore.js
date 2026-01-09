import { create } from 'zustand';
import * as SecureStore from 'expo-secure-store';
import { lightTheme, darkTheme } from '../theme/theme';

export const useThemeStore = create((set, get) => ({
  // State
  isDark: false,
  theme: lightTheme,

  // Initialize theme from storage
  initialize: async () => {
    try {
      const isDarkStr = await SecureStore.getItemAsync('isDark');
      const isDark = isDarkStr === 'true';
      set({
        isDark,
        theme: isDark ? darkTheme : lightTheme,
      });
    } catch {
      // Default to light theme
    }
  },

  // Toggle theme
  toggleTheme: async () => {
    const newIsDark = !get().isDark;
    await SecureStore.setItemAsync('isDark', String(newIsDark));
    set({
      isDark: newIsDark,
      theme: newIsDark ? darkTheme : lightTheme,
    });
  },

  // Set specific theme
  setTheme: async (isDark) => {
    await SecureStore.setItemAsync('isDark', String(isDark));
    set({
      isDark,
      theme: isDark ? darkTheme : lightTheme,
    });
  },
}));
