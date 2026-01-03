import { useThemeStore, THEMES } from '../../stores/themeStore';
import { Sun, Moon, Monitor } from 'lucide-react';

/**
 * Animated Theme Toggle Button
 * A beautiful toggle switch with sun/moon icons and smooth animations
 */
export function ThemeToggle({ showLabel = false, size = 'md' }) {
  const { theme, toggleTheme, isSystemPreference, setSystemPreference } = useThemeStore();
  const isDark = theme === THEMES.DARK;

  // Size variants
  const sizes = {
    sm: {
      button: 'w-14 h-7',
      circle: 'w-5 h-5',
      translate: 'translate-x-7',
      icon: 14,
    },
    md: {
      button: 'w-16 h-8',
      circle: 'w-6 h-6',
      translate: 'translate-x-8',
      icon: 16,
    },
    lg: {
      button: 'w-20 h-10',
      circle: 'w-8 h-8',
      translate: 'translate-x-10',
      icon: 20,
    },
  };

  const s = sizes[size] || sizes.md;

  return (
    <div className="flex items-center gap-3">
      {showLabel && <span className="text-sm text-base-content/60 font-medium">{isDark ? 'Dark' : 'Light'}</span>}
      <button
        onClick={toggleTheme}
        className={`
                    relative ${s.button} rounded-full p-1
                    bg-base-300 hover:bg-base-300/80
                    transition-all duration-300 ease-in-out
                    focus:outline-none focus:ring-2 focus:ring-primary/30 focus:ring-offset-2 focus:ring-offset-base-100
                    group
                `}
        aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
      >
        {/* Background glow effect */}
        <div
          className={`
                        absolute inset-0 rounded-full opacity-0 group-hover:opacity-100
                        transition-opacity duration-300
                        ${isDark ? 'bg-gradient-to-r from-indigo-500/20 to-purple-500/20' : 'bg-gradient-to-r from-amber-500/20 to-orange-500/20'}
                    `}
        />

        {/* Toggle circle with icon */}
        <div
          className={`
                        ${s.circle} rounded-full
                        flex items-center justify-center
                        transform transition-all duration-500 ease-[cubic-bezier(0.68,-0.55,0.265,1.55)]
                        ${isDark ? s.translate : 'translate-x-0'}
                        ${isDark ? 'bg-slate-700 shadow-lg shadow-indigo-500/20' : 'bg-white shadow-lg shadow-amber-500/20'}
                    `}
        >
          {/* Sun icon */}
          <Sun
            size={s.icon}
            className={`
                            absolute transition-all duration-300
                            ${isDark ? 'opacity-0 rotate-90 scale-0' : 'opacity-100 rotate-0 scale-100'}
                            text-amber-500
                        `}
          />
          {/* Moon icon */}
          <Moon
            size={s.icon}
            className={`
                            absolute transition-all duration-300
                            ${isDark ? 'opacity-100 rotate-0 scale-100' : 'opacity-0 -rotate-90 scale-0'}
                            text-indigo-300
                        `}
          />
        </div>

        {/* Stars animation for dark mode */}
        <div
          className={`
                        absolute inset-0 overflow-hidden rounded-full pointer-events-none
                        transition-opacity duration-500
                        ${isDark ? 'opacity-100' : 'opacity-0'}
                    `}
        >
          <div className="absolute w-1 h-1 bg-white/60 rounded-full top-1.5 left-2 animate-pulse" style={{ animationDelay: '0ms' }} />
          <div className="absolute w-0.5 h-0.5 bg-white/40 rounded-full top-3 left-4 animate-pulse" style={{ animationDelay: '200ms' }} />
          <div className="absolute w-1 h-1 bg-white/50 rounded-full bottom-2 left-3 animate-pulse" style={{ animationDelay: '400ms' }} />
        </div>
      </button>
    </div>
  );
}

/**
 * Theme Selector with System Preference Option
 * A dropdown/segmented control for Light/Dark/System
 */
export function ThemeSelector() {
  const { theme, setTheme, isSystemPreference, setSystemPreference } = useThemeStore();

  const options = [
    { value: 'light', label: 'Light', icon: Sun },
    { value: 'dark', label: 'Dark', icon: Moon },
    { value: 'system', label: 'System', icon: Monitor },
  ];

  const getCurrentValue = () => {
    if (isSystemPreference) return 'system';
    return theme === THEMES.DARK ? 'dark' : 'light';
  };

  const handleSelect = (value) => {
    if (value === 'system') {
      setSystemPreference(true);
    } else if (value === 'light') {
      setTheme(THEMES.LIGHT);
    } else {
      setTheme(THEMES.DARK);
    }
  };

  const currentValue = getCurrentValue();

  return (
    <div className="flex bg-base-200 rounded-lg p-1 gap-1">
      {options.map((option) => {
        const Icon = option.icon;
        const isActive = currentValue === option.value;
        return (
          <button
            key={option.value}
            onClick={() => handleSelect(option.value)}
            className={`
                            flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium
                            transition-all duration-200 ease-out
                            ${isActive ? 'bg-base-100 text-base-content shadow-sm' : 'text-base-content/60 hover:text-base-content hover:bg-base-100/50'}
                        `}
          >
            <Icon size={16} className={`transition-transform duration-200 ${isActive ? 'scale-110' : 'scale-100'}`} />
            <span className="hidden sm:inline">{option.label}</span>
          </button>
        );
      })}
    </div>
  );
}

/**
 * Compact Theme Icon Button
 * A simple icon button that toggles between light and dark
 */
export function ThemeIconButton({ className = '' }) {
  const { theme, toggleTheme } = useThemeStore();
  const isDark = theme === THEMES.DARK;

  return (
    <button
      onClick={toggleTheme}
      className={`
                relative p-2 rounded-lg
                bg-base-200 hover:bg-base-300
                text-base-content/70 hover:text-base-content
                transition-all duration-200
                focus:outline-none focus:ring-2 focus:ring-primary/30 focus:ring-offset-2 focus:ring-offset-base-100
                overflow-hidden
                ${className}
            `}
      aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {/* Animated icon container */}
      <div className="relative w-5 h-5">
        {/* Sun icon */}
        <Sun
          size={20}
          className={`
                        absolute inset-0 transition-all duration-500
                        ${isDark ? 'opacity-0 rotate-180 scale-0' : 'opacity-100 rotate-0 scale-100'}
                        text-amber-500
                    `}
        />
        {/* Moon icon */}
        <Moon
          size={20}
          className={`
                        absolute inset-0 transition-all duration-500
                        ${isDark ? 'opacity-100 rotate-0 scale-100' : 'opacity-0 rotate-180 scale-0'}
                        text-indigo-400
                    `}
        />
      </div>
    </button>
  );
}
