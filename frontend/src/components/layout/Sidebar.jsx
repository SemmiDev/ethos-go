import { useState, useRef, useEffect } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { LayoutDashboard, Target, BarChart3, Settings, LogOut, Shield, ChevronUp, User, HelpCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useAuthStore } from '../../stores/authStore';
import { useUIStore } from '../../stores/uiStore';
import { LanguageToggleAnimated } from '../ui/LanguageToggle';

const getNavItems = (t) => [
  { icon: LayoutDashboard, label: t('nav.dashboard'), path: '/dashboard' },
  { icon: Target, label: t('nav.habits'), path: '/habits' },
  { icon: BarChart3, label: t('nav.analytics'), path: '/analytics' },
  { icon: Settings, label: t('nav.settings'), path: '/settings' },
];

export function Sidebar() {
  const { t } = useTranslation();
  const location = useLocation();
  const navigate = useNavigate();
  const { logout, user } = useAuthStore();
  const { addToast } = useUIStore();
  const [isProfileMenuOpen, setIsProfileMenuOpen] = useState(false);
  const profileMenuRef = useRef(null);
  const navItems = getNavItems(t);

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (profileMenuRef.current && !profileMenuRef.current.contains(event.target)) {
        setIsProfileMenuOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = async () => {
    setIsProfileMenuOpen(false);
    await logout();
    addToast({
      type: 'success',
      title: 'Logged out',
      message: 'You have been logged out successfully.',
    });
    navigate('/login');
  };

  const closeDrawer = () => {
    const elem = document.getElementById('my-drawer');
    if (elem) elem.checked = false;
    setIsProfileMenuOpen(false);
  };

  const handleProfileMenuItemClick = (path) => {
    setIsProfileMenuOpen(false);
    closeDrawer();
    navigate(path);
  };

  return (
    <div className="drawer-side z-50">
      <label htmlFor="my-drawer" aria-label="close sidebar" className="drawer-overlay"></label>
      <aside className="w-64 min-h-full bg-base-100 border-r border-base-300 flex flex-col">
        {/* Logo */}
        <div className="px-5 py-6 border-b border-base-200">
          <Link to="/" className="flex items-center gap-3" onClick={closeDrawer}>
            <img src="/logo.jpg" alt="Ethos Logo" className="w-10 h-10 rounded-lg object-cover" />
            <div>
              <h1 className="text-lg font-semibold text-base-content tracking-tight">Ethos</h1>
              <p className="text-[11px] text-base-content/50 font-medium uppercase tracking-wider">Habit Tracker</p>
            </div>
          </Link>
        </div>

        {/* Navigation */}
        <nav className="flex-1 px-3 py-4">
          <p className="px-3 mb-2 text-[10px] font-semibold text-base-content/40 uppercase tracking-wider">Menu</p>
          <ul className="space-y-1">
            {navItems.map((item) => {
              const isActive = location.pathname.startsWith(item.path);
              return (
                <li key={item.path}>
                  <Link
                    to={item.path}
                    onClick={closeDrawer}
                    className={`
                      flex items-center gap-3 px-3 py-2.5 rounded-md text-sm font-medium
                      transition-colors duration-150
                      ${isActive ? 'bg-primary text-primary-content' : 'text-base-content/70 hover:bg-base-200 hover:text-base-content'}
                    `}
                  >
                    <item.icon size={18} />
                    {item.label}
                  </Link>
                </li>
              );
            })}
          </ul>
        </nav>

        {/* Language Toggle */}
        <div className="px-3 pb-3">
          <div className="flex items-center justify-between px-3 py-2 bg-base-200/50 rounded-lg">
            <span className="text-xs font-medium text-base-content/60">{t('settings.language.title')}</span>
            <LanguageToggleAnimated />
          </div>
        </div>

        {/* User Profile Section with Dropdown */}
        <div className="p-3 border-t border-base-200" ref={profileMenuRef}>
          <div className="relative">
            {/* Profile Dropdown Menu - Positioned above */}
            <div
              className={`
                absolute bottom-full left-0 right-0 mb-2
                bg-base-100 border border-base-300 rounded-lg shadow-lg
                overflow-hidden
                transition-all duration-200 ease-out origin-bottom
                ${isProfileMenuOpen ? 'opacity-100 scale-100 translate-y-0' : 'opacity-0 scale-95 translate-y-2 pointer-events-none'}
              `}
            >
              {/* Menu Header */}
              <div className="px-4 py-3 bg-base-200/50 border-b border-base-200">
                <p className="text-xs font-medium text-base-content/50 uppercase tracking-wider">Account</p>
              </div>

              {/* Menu Items */}
              <div className="py-1">
                <button
                  onClick={() => handleProfileMenuItemClick('/settings')}
                  className="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-base-content/80 hover:bg-base-200 hover:text-base-content transition-colors"
                >
                  <User size={16} />
                  <span>{t('settings.profile.title')}</span>
                </button>
                <button
                  onClick={() => handleProfileMenuItemClick('/help')}
                  className="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-base-content/80 hover:bg-base-200 hover:text-base-content transition-colors"
                >
                  <HelpCircle size={16} />
                  <span>Help & Support</span>
                </button>
              </div>

              {/* Divider */}
              <div className="border-t border-base-200" />

              {/* Logout */}
              <div className="py-1">
                <button onClick={handleLogout} className="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-error hover:bg-error/5 transition-colors">
                  <LogOut size={16} />
                  <span>{t('nav.logout')}</span>
                </button>
              </div>
            </div>

            {/* Profile Button */}
            <button
              onClick={() => setIsProfileMenuOpen(!isProfileMenuOpen)}
              className={`
                w-full flex items-center gap-3 p-2.5 rounded-lg
                transition-all duration-150
                ${isProfileMenuOpen ? 'bg-base-200' : 'hover:bg-base-200/70'}
              `}
            >
              {/* Avatar */}
              <div className="w-9 h-9 rounded-full bg-gradient-to-br from-primary to-secondary flex items-center justify-center shrink-0 shadow-sm">
                <span className="text-sm font-semibold text-white">{user?.name?.charAt(0)?.toUpperCase() || 'U'}</span>
              </div>

              {/* User Info */}
              <div className="flex-1 min-w-0 text-left">
                <p className="text-sm font-medium text-base-content truncate">{user?.name || 'User'}</p>
                <p className="text-xs text-base-content/50 truncate">{user?.email || 'user@example.com'}</p>
              </div>

              {/* Chevron */}
              <ChevronUp
                size={16}
                className={`
                  text-base-content/40 shrink-0
                  transition-transform duration-200
                  ${isProfileMenuOpen ? 'rotate-180' : 'rotate-0'}
                `}
              />
            </button>
          </div>
        </div>
      </aside>
    </div>
  );
}

export function Header({ title, subtitle, actions }) {
  return (
    <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 pb-6 border-b border-base-200 mb-8">
      <div>
        <h1 className="text-xl font-semibold text-base-content tracking-tight">{title}</h1>
        {subtitle && <p className="text-sm text-base-content/50 mt-1">{subtitle}</p>}
      </div>
      {actions && <div className="flex gap-3 w-full sm:w-auto">{actions}</div>}
    </div>
  );
}
