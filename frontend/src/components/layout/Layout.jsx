import { useEffect } from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '../../stores/authStore';
import { Sidebar } from './Sidebar';
import { ToastContainer } from '../ui/Toast';
import { NotificationBell } from '../notifications/NotificationBell';
import { ThemeIconButton } from '../ui/ThemeToggle';
import { CommandPalette, CommandPaletteTrigger } from '../ui/CommandPalette';
import { Celebration } from '../ui/Celebration';
import { OnboardingTour, useAutoStartTour } from '../ui/OnboardingTour';
import { PWAInstallBanner, usePWA } from '../ui/PWAInstall';
import { KeyboardShortcutsModal, useKeyboardShortcuts, ShortcutIndicator } from '../ui/KeyboardShortcuts';
import { CreateHabitModal } from '../habits';
import { useThemeStore } from '../../stores/themeStore';
import { useUIStore } from '../../stores/uiStore';
import { Menu, Shield } from 'lucide-react';

export function MainLayout() {
  const { isAuthenticated } = useAuthStore();
  const { toggleTheme } = useThemeStore();
  const { isCreateHabitModalOpen, openCreateHabitModal, closeCreateHabitModal } = useUIStore();
  const { isHelpOpen, closeHelp } = useKeyboardShortcuts();

  // Initialize PWA
  usePWA();

  // Auto-start tour for new users
  useAutoStartTour();

  // Listen for keyboard shortcut events
  useEffect(() => {
    const handleCreateHabit = () => openCreateHabitModal?.();
    const handleToggleTheme = () => toggleTheme();

    window.addEventListener('shortcut:create-habit', handleCreateHabit);
    window.addEventListener('shortcut:toggle-theme', handleToggleTheme);

    return () => {
      window.removeEventListener('shortcut:create-habit', handleCreateHabit);
      window.removeEventListener('shortcut:toggle-theme', handleToggleTheme);
    };
  }, [openCreateHabitModal, toggleTheme]);

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return (
    <div className="drawer lg:drawer-open">
      <input id="my-drawer" type="checkbox" className="drawer-toggle" />

      <div className="drawer-content flex flex-col bg-base-200 min-h-screen">
        {/* Top Header Bar */}
        <header className="sticky top-0 z-40 bg-base-100/80 backdrop-blur-md border-b border-base-300/50">
          <div className="flex items-center justify-between h-14 px-4 lg:px-6">
            {/* Left Side - Mobile menu + Logo (mobile) */}
            <div className="flex items-center gap-3">
              {/* Mobile Menu Button */}
              <label
                htmlFor="my-drawer"
                className="lg:hidden p-2 -ml-2 rounded-lg text-base-content/70 hover:bg-base-200 hover:text-base-content cursor-pointer transition-colors"
              >
                <Menu size={20} />
              </label>

              {/* Mobile Logo */}
              <div className="lg:hidden flex items-center gap-2">
                <div className="w-7 h-7 rounded-md bg-primary flex items-center justify-center">
                  <Shield size={14} className="text-primary-content" />
                </div>
                <span className="text-base font-semibold text-base-content">Ethos</span>
              </div>
            </div>

            {/* Center - Command Palette Trigger (desktop) */}
            <div className="hidden lg:flex flex-1 justify-center max-w-md mx-4">
              <CommandPaletteTrigger className="w-full max-w-sm" />
            </div>

            {/* Right Side - Actions */}
            <div className="flex items-center gap-1">
              {/* Mobile search trigger */}
              <div className="lg:hidden">
                <CommandPaletteTrigger />
              </div>
              <ThemeIconButton data-tour="theme-toggle" />
              <NotificationBell />
            </div>
          </div>
        </header>

        {/* Main Content */}
        <main className="flex-1 p-4 md:p-6 lg:p-8">
          <div className="max-w-6xl mx-auto">
            <Outlet />
          </div>
        </main>
      </div>

      <Sidebar />
      <ToastContainer />

      {/* Global Components */}
      <CommandPalette />
      <Celebration />
      <OnboardingTour />
      <PWAInstallBanner />
      <KeyboardShortcutsModal isOpen={isHelpOpen} onClose={closeHelp} />
      <ShortcutIndicator />

      {/* Global Create Habit Modal */}
      <CreateHabitModal isOpen={isCreateHabitModalOpen} onClose={closeCreateHabitModal} />
    </div>
  );
}

export function AuthLayout() {
  const { isAuthenticated } = useAuthStore();

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return (
    <>
      <Outlet />
      <ToastContainer />
    </>
  );
}
