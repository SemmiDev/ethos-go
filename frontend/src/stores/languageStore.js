import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import i18n from '../i18n';

export const useLanguageStore = create(
  persist(
    (set) => ({
      language: i18n.language || 'en',

      setLanguage: (lang) => {
        i18n.changeLanguage(lang);
        set({ language: lang });
        // Update HTML lang attribute
        document.documentElement.lang = lang;
      },
    }),
    {
      name: 'ethos-language-storage',
      onRehydrateStorage: () => (state) => {
        // When store is rehydrated, sync i18n with stored language
        if (state?.language) {
          i18n.changeLanguage(state.language);
          document.documentElement.lang = state.language;
        }
      },
    }
  )
);
