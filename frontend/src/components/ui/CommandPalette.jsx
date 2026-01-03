import { useState, useEffect, useRef, useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Search,
  LayoutDashboard,
  Target,
  BarChart3,
  Settings,
  HelpCircle,
  Bell,
  LogOut,
  Sun,
  Moon,
  Plus,
  User,
  Command,
  ArrowRight,
  Globe,
} from 'lucide-react';
import { useAuthStore } from '../../stores/authStore';
import { useUIStore } from '../../stores/uiStore';
import { useThemeStore, THEMES } from '../../stores/themeStore';
import { useHabitsStore } from '../../stores/habitsStore';
import { useLanguageStore } from '../../stores/languageStore';

// Command types for categorization
const COMMAND_TYPES = {
  NAVIGATION: 'navigation',
  ACTION: 'action',
  HABIT: 'habit',
  QUICK: 'quick',
};

// Animation variants
const backdropVariants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
  exit: { opacity: 0 },
};

const paletteVariants = {
  hidden: {
    opacity: 0,
    scale: 0.95,
    y: -20,
  },
  visible: {
    opacity: 1,
    scale: 1,
    y: 0,
    transition: {
      type: 'spring',
      damping: 25,
      stiffness: 300,
    },
  },
  exit: {
    opacity: 0,
    scale: 0.95,
    y: -10,
    transition: {
      duration: 0.15,
    },
  },
};

const itemVariants = {
  hidden: { opacity: 0, x: -10 },
  visible: (i) => ({
    opacity: 1,
    x: 0,
    transition: {
      delay: i * 0.03,
      duration: 0.2,
    },
  }),
};

// Main Command Palette Component
export function CommandPalette() {
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef(null);
  const listRef = useRef(null);

  const { t } = useTranslation();
  const navigate = useNavigate();
  const { logout } = useAuthStore();
  const { addToast, openCreateHabitModal } = useUIStore();
  const { theme, toggleTheme } = useThemeStore();
  const { habits } = useHabitsStore();
  const { language, setLanguage } = useLanguageStore();
  const isDark = theme === THEMES.DARK;

  // Define all available commands
  const allCommands = useMemo(
    () => [
      // Quick Actions
      {
        id: 'new-habit',
        type: COMMAND_TYPES.QUICK,
        icon: Plus,
        title: t('commands.actions.createHabit'),
        subtitle: t('commands.actions.createHabitDesc'),
        keywords: ['new', 'add', 'create', 'habit', 'baru', 'tambah', 'kebiasaan'],
        action: () => {
          openCreateHabitModal();
        },
      },
      {
        id: 'toggle-theme',
        type: COMMAND_TYPES.QUICK,
        icon: isDark ? Sun : Moon,
        title: t('commands.actions.toggleTheme', { mode: isDark ? t('settings.appearance.light') : t('settings.appearance.dark') }),
        subtitle: t('commands.actions.toggleThemeDesc'),
        keywords: ['theme', 'dark', 'light', 'mode', 'toggle', 'tema', 'gelap', 'terang'],
        action: () => {
          toggleTheme();
          addToast({
            type: 'success',
            title: t('toast.habitUpdated'), // Reusing a success toast or general success
            message: `Switched to ${isDark ? 'light' : 'dark'} mode`,
          });
        },
      },
      {
        id: 'switch-language',
        type: COMMAND_TYPES.QUICK,
        icon: Globe,
        title: t('commands.actions.switchLanguage'),
        subtitle: t('commands.actions.switchLanguageDesc', { lang: language === 'en' ? 'Bahasa Indonesia' : 'English' }),
        keywords: ['language', 'bahasa', 'english', 'indonesia', 'switch', 'ganti'],
        action: () => {
          const newLang = language === 'en' ? 'id' : 'en';
          setLanguage(newLang);
          addToast({
            type: 'success',
            title: t('common.success'),
            message: newLang === 'en' ? 'Language switched to English' : 'Bahasa diganti ke Indonesia',
          });
        },
      },

      // Navigation
      {
        id: 'nav-dashboard',
        type: COMMAND_TYPES.NAVIGATION,
        icon: LayoutDashboard,
        title: t('commands.actions.goDashboard'),
        subtitle: t('nav.dashboard'),
        keywords: ['dashboard', 'home', 'overview', 'beranda'],
        action: () => navigate('/dashboard'),
      },
      {
        id: 'nav-habits',
        type: COMMAND_TYPES.NAVIGATION,
        icon: Target,
        title: t('commands.actions.goHabits'),
        subtitle: t('nav.habits'),
        keywords: ['habits', 'list', 'track', 'kebiasaan'],
        action: () => navigate('/habits'),
      },
      {
        id: 'nav-analytics',
        type: COMMAND_TYPES.NAVIGATION,
        icon: BarChart3,
        title: t('commands.actions.goAnalytics'),
        subtitle: t('nav.analytics'),
        keywords: ['analytics', 'stats', 'charts', 'progress', 'analitik', 'statistik'],
        action: () => navigate('/analytics'),
      },
      {
        id: 'nav-settings',
        type: COMMAND_TYPES.NAVIGATION,
        icon: Settings,
        title: t('commands.actions.goSettings'),
        subtitle: t('nav.settings'),
        keywords: ['settings', 'preferences', 'config', 'pengaturan'],
        action: () => navigate('/settings'),
      },
      {
        id: 'nav-notifications',
        type: COMMAND_TYPES.NAVIGATION,
        icon: Bell,
        title: t('commands.actions.goNotifications'),
        subtitle: t('nav.notifications'),
        keywords: ['notifications', 'alerts', 'messages', 'notifikasi'],
        action: () => navigate('/notifications'),
      },
      {
        id: 'nav-help',
        type: COMMAND_TYPES.NAVIGATION,
        icon: HelpCircle,
        title: t('commands.actions.goHelp'),
        subtitle: t('help.title'),
        keywords: ['help', 'support', 'faq', 'docs', 'bantuan'],
        action: () => navigate('/help'),
      },

      // Actions
      {
        id: 'action-profile',
        type: COMMAND_TYPES.ACTION,
        icon: User,
        title: t('commands.actions.editProfile'),
        subtitle: t('commands.actions.editProfileDesc'),
        keywords: ['profile', 'account', 'user', 'edit', 'profil', 'akun'],
        action: () => navigate('/settings'),
      },
      {
        id: 'action-logout',
        type: COMMAND_TYPES.ACTION,
        icon: LogOut,
        title: t('commands.actions.logout'),
        subtitle: t('commands.actions.logoutDesc'),
        keywords: ['logout', 'signout', 'exit', 'leave', 'keluar'],
        action: async () => {
          await logout();
          addToast({
            type: 'success',
            title: t('auth.logout'),
            message: t('auth.loginSuccess'), // Usually "Logged out successfully" but reusing available key or generic
          });
          navigate('/login');
        },
        danger: true,
      },

      // Habits as searchable items
      ...habits.map((habit) => ({
        id: `habit-${habit.id}`,
        type: COMMAND_TYPES.HABIT,
        icon: Target,
        title: habit.name,
        subtitle: habit.description || t('habits.detail.title'),
        keywords: [habit.name.toLowerCase(), 'habit', 'kebiasaan'],
        action: () => navigate(`/habits/${habit.id}`),
      })),
    ],
    [isDark, habits, navigate, logout, addToast, toggleTheme, openCreateHabitModal, language, setLanguage, t]
  );

  // Filter commands based on search
  const filteredCommands = useMemo(() => {
    if (!search.trim()) {
      return allCommands.slice(0, 9); // Increased to include language switch
    }

    const searchLower = search.toLowerCase();
    return allCommands.filter((cmd) => {
      const titleMatch = cmd.title.toLowerCase().includes(searchLower);
      const subtitleMatch = cmd.subtitle?.toLowerCase().includes(searchLower);
      const keywordMatch = cmd.keywords?.some((k) => k.includes(searchLower));
      return titleMatch || subtitleMatch || keywordMatch;
    });
  }, [search, allCommands]);

  // Group commands by type
  const groupedCommands = useMemo(() => {
    const groups = {
      [COMMAND_TYPES.QUICK]: { title: t('commands.groups.quick'), items: [] },
      [COMMAND_TYPES.NAVIGATION]: { title: t('commands.groups.navigation'), items: [] },
      [COMMAND_TYPES.HABIT]: { title: t('commands.groups.habits'), items: [] },
      [COMMAND_TYPES.ACTION]: { title: t('commands.groups.actions'), items: [] },
    };

    filteredCommands.forEach((cmd) => {
      if (groups[cmd.type]) {
        groups[cmd.type].items.push(cmd);
      }
    });

    return Object.entries(groups).filter(([_, group]) => group.items.length > 0);
  }, [filteredCommands, t]);

  // Flatten for keyboard navigation
  const flatCommands = useMemo(() => {
    return groupedCommands.flatMap(([_, group]) => group.items);
  }, [groupedCommands]);

  // Open/Close handlers
  const openPalette = useCallback(() => {
    setIsOpen(true);
    setSearch('');
    setSelectedIndex(0);
  }, []);

  const closePalette = useCallback(() => {
    setIsOpen(false);
    setSearch('');
    setSelectedIndex(0);
  }, []);

  // Execute selected command
  const executeCommand = useCallback(
    (command) => {
      closePalette();
      command.action();
    },
    [closePalette]
  );

  // Keyboard shortcut listener (Cmd+K / Ctrl+K)
  useEffect(() => {
    const handleKeyDown = (e) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        if (isOpen) {
          closePalette();
        } else {
          openPalette();
        }
      }

      if (e.key === 'Escape' && isOpen) {
        e.preventDefault();
        closePalette();
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, openPalette, closePalette]);

  // Focus input when opened
  useEffect(() => {
    if (isOpen && inputRef.current) {
      setTimeout(() => inputRef.current?.focus(), 50);
    }
  }, [isOpen]);

  // Keyboard navigation within palette
  const handleKeyDown = (e) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex((prev) => Math.min(prev + 1, flatCommands.length - 1));
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex((prev) => Math.max(prev - 1, 0));
        break;
      case 'Enter':
        e.preventDefault();
        if (flatCommands[selectedIndex]) {
          executeCommand(flatCommands[selectedIndex]);
        }
        break;
    }
  };

  // Scroll selected item into view
  useEffect(() => {
    if (listRef.current && flatCommands[selectedIndex]) {
      const selectedElement = listRef.current.querySelector(`[data-index="${selectedIndex}"]`);
      if (selectedElement) {
        selectedElement.scrollIntoView({ block: 'nearest' });
      }
    }
  }, [selectedIndex, flatCommands]);

  // Reset selected index when search changes
  useEffect(() => {
    setSelectedIndex(0);
  }, [search]);

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Backdrop */}
          <motion.div
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
            variants={backdropVariants}
            initial="hidden"
            animate="visible"
            exit="exit"
            onClick={closePalette}
          />

          {/* Palette Container */}
          <div className="fixed inset-0 z-50 flex items-start justify-center pt-[12vh] px-4 pointer-events-none">
            <motion.div
              className="w-full max-w-2xl bg-base-100 rounded-2xl shadow-2xl border border-base-300/50 pointer-events-auto overflow-hidden"
              variants={paletteVariants}
              initial="hidden"
              animate="visible"
              exit="exit"
              onClick={(e) => e.stopPropagation()}
            >
              {/* Search Input */}
              <div className="relative">
                <div className="absolute left-5 top-1/2 -translate-y-1/2">
                  <Search size={20} className="text-primary" />
                </div>
                <input
                  ref={inputRef}
                  type="text"
                  placeholder={t('commands.placeholder')}
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  onKeyDown={handleKeyDown}
                  className="w-full h-14 pl-14 pr-24 bg-transparent text-base-content text-base placeholder:text-base-content/40 outline-none border-b border-base-200 font-medium"
                />
                <div className="absolute right-4 top-1/2 -translate-y-1/2 flex items-center gap-2">
                  <kbd className="px-2 py-1 text-xs font-semibold text-base-content/50 bg-base-200 rounded-md border border-base-300">ESC</kbd>
                </div>
              </div>

              {/* Results */}
              <div ref={listRef} className="max-h-[50vh] overflow-y-auto py-2">
                {flatCommands.length === 0 ? (
                  <motion.div className="px-6 py-12 text-center" initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
                    <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-base-200 flex items-center justify-center">
                      <Search size={24} className="text-base-content/30" />
                    </div>
                    <p className="text-base-content/60 font-medium">{t('commands.noResults')}</p>
                    <p className="text-base-content/40 text-sm mt-1">{t('commands.tryAgain')}</p>
                  </motion.div>
                ) : (
                  groupedCommands.map(([type, group], groupIndex) => (
                    <div key={type} className="mb-2">
                      {/* Group Header */}
                      <div className="px-5 py-2">
                        <p className="text-xs font-semibold text-base-content/40 uppercase tracking-wider">{group.title}</p>
                      </div>

                      {/* Group Items */}
                      {group.items.map((command, itemIndex) => {
                        const flatIndex = groupedCommands.slice(0, groupIndex).reduce((acc, [_, g]) => acc + g.items.length, 0) + itemIndex;

                        const isSelected = selectedIndex === flatIndex;
                        const Icon = command.icon;

                        return (
                          <motion.button
                            key={command.id}
                            data-index={flatIndex}
                            onClick={() => executeCommand(command)}
                            onMouseEnter={() => setSelectedIndex(flatIndex)}
                            variants={itemVariants}
                            initial="hidden"
                            animate="visible"
                            custom={flatIndex}
                            className={`
                                                            w-full flex items-center gap-4 px-5 py-3 text-left
                                                            transition-all duration-100 mx-2 rounded-xl
                                                            ${
                                                              isSelected
                                                                ? command.danger
                                                                  ? 'bg-error/10 text-error'
                                                                  : 'bg-primary/10 text-base-content'
                                                                : 'text-base-content/70 hover:bg-base-200'
                                                            }
                                                        `}
                            style={{ width: 'calc(100% - 16px)' }}
                          >
                            <motion.div
                              className={`
                                                                p-2.5 rounded-xl shrink-0 transition-colors
                                                                ${isSelected ? (command.danger ? 'bg-error/20' : 'bg-primary/20') : 'bg-base-200'}
                                                            `}
                              whileHover={{ scale: 1.05 }}
                              whileTap={{ scale: 0.95 }}
                            >
                              <Icon size={18} className={isSelected ? (command.danger ? 'text-error' : 'text-primary') : 'text-base-content/60'} />
                            </motion.div>
                            <div className="flex-1 min-w-0">
                              <p className={`text-sm font-semibold truncate ${command.danger && isSelected ? 'text-error' : ''}`}>{command.title}</p>
                              {command.subtitle && <p className="text-xs text-base-content/50 truncate mt-0.5">{command.subtitle}</p>}
                            </div>
                            {isSelected && (
                              <motion.div initial={{ opacity: 0, x: -10 }} animate={{ opacity: 1, x: 0 }} className="flex items-center gap-2">
                                <span className="text-xs text-base-content/40">{t('commands.tips.select')}</span>
                                <ArrowRight size={14} className="text-base-content/40" />
                              </motion.div>
                            )}
                          </motion.button>
                        );
                      })}
                    </div>
                  ))
                )}
              </div>

              {/* Footer with keyboard hints */}
              <div className="flex items-center justify-between px-5 py-3 border-t border-base-200 bg-base-200/30">
                <div className="flex items-center gap-4 text-xs text-base-content/50">
                  <span className="flex items-center gap-1.5">
                    <kbd className="px-1.5 py-0.5 bg-base-300 rounded font-semibold">↑↓</kbd>
                    <span>{t('commands.tips.navigate')}</span>
                  </span>
                  <span className="flex items-center gap-1.5">
                    <kbd className="px-1.5 py-0.5 bg-base-300 rounded font-semibold">↵</kbd>
                    <span>{t('commands.tips.select')}</span>
                  </span>
                  <span className="flex items-center gap-1.5">
                    <kbd className="px-1.5 py-0.5 bg-base-300 rounded font-semibold">ESC</kbd>
                    <span>{t('commands.tips.close')}</span>
                  </span>
                </div>
                <div className="flex items-center gap-1.5 text-xs text-base-content/50">
                  <Command size={12} />
                  <span className="font-semibold">K</span>
                </div>
              </div>
            </motion.div>
          </div>
        </>
      )}
    </AnimatePresence>
  );
}

// Trigger button component for header/navbar
export function CommandPaletteTrigger({ className = '' }) {
  const { t } = useTranslation();

  const openPalette = () => {
    const event = new KeyboardEvent('keydown', {
      key: 'k',
      metaKey: true,
      bubbles: true,
    });
    document.dispatchEvent(event);
  };

  return (
    <motion.button
      onClick={openPalette}
      className={`
                flex items-center gap-3 px-4 py-2.5
                bg-base-100 hover:bg-base-200
                rounded-xl border border-base-300
                text-sm text-base-content/60 hover:text-base-content
                transition-colors duration-150 shadow-sm
                ${className}
            `}
      whileHover={{ scale: 1.01 }}
      whileTap={{ scale: 0.99 }}
    >
      <Search size={16} className="text-base-content/40" />
      <span className="hidden sm:inline text-base-content/50">{t('commands.placeholder')}</span>
      <kbd className="hidden sm:flex items-center gap-1 ml-auto px-2 py-1 text-xs font-semibold bg-base-200 text-base-content/50 rounded-md border border-base-300">
        <Command size={11} />
        <span>K</span>
      </kbd>
    </motion.button>
  );
}
