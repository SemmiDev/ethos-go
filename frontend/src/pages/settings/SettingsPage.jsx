import { useState, useEffect } from 'react';
import {
  User,
  Bell,
  Shield,
  Trash2,
  Save,
  Monitor,
  LogOut,
  Key,
  ChevronDown,
  ChevronUp,
  Sun,
  Moon,
  Palette,
  Smartphone,
  Download,
  CheckCircle,
  Globe,
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Header } from '../../components/layout/Sidebar';
import { Card } from '../../components/ui/Card';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { ConfirmModal, Modal } from '../../components/ui/Modal';
import { Badge } from '../../components/ui/Badge';
import { ThemeSelector } from '../../components/ui/ThemeToggle';
import { LanguageToggleAnimated } from '../../components/ui/LanguageToggle';
import { useAuthStore } from '../../stores/authStore';
import { useUIStore } from '../../stores/uiStore';
import { useThemeStore, THEMES } from '../../stores/themeStore';
import { useLanguageStore } from '../../stores/languageStore';

import { languages } from '../../i18n';
import { format } from 'date-fns';

// Move components OUTSIDE the main component to prevent re-creation on every render
const SettingSection = ({ icon: Icon, title, description, children, collapsible = false, defaultOpen = true }) => {
  const [isOpen, setIsOpen] = useState(defaultOpen);

  return (
    <div className="bg-base-100 border border-base-300 rounded-lg p-6 mb-6">
      <div className="flex items-start gap-4">
        <div className="p-2.5 bg-primary/5 rounded-lg shrink-0">
          <Icon size={20} className="text-primary" />
        </div>
        <div className="flex-1 min-w-0">
          <div
            className={`flex items-start justify-between ${collapsible ? 'cursor-pointer select-none group' : ''}`}
            onClick={collapsible ? () => setIsOpen(!isOpen) : undefined}
          >
            <div>
              <h3 className={`text-base font-semibold text-base-content mb-1 ${collapsible ? 'group-hover:text-primary transition-colors' : ''}`}>{title}</h3>
              <p className="text-sm text-base-content/50 mb-2">{description}</p>
            </div>
            {collapsible && (
              <div className="mt-1 text-base-content/40 group-hover:text-primary transition-colors">
                {isOpen ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
              </div>
            )}
          </div>

          {(!collapsible || isOpen) && <div className={collapsible ? 'mt-4 animate-in fade-in slide-in-from-top-2 duration-200' : 'mt-4'}>{children}</div>}
        </div>
      </div>
    </div>
  );
};

const Toggle = ({ checked, onChange, label }) => (
  <div className="flex items-center justify-between py-3 border-b border-base-200 last:border-0">
    <span className="text-sm font-medium text-base-content">{label}</span>
    <button
      type="button"
      onClick={() => onChange(!checked)}
      className={`
                relative inline-flex h-6 w-11 items-center rounded-full transition-colors
                ${checked ? 'bg-primary' : 'bg-base-300'}
            `}
    >
      <span
        className={`
                    inline-block h-4 w-4 rounded-full bg-white transition-transform shadow-sm
                    ${checked ? 'translate-x-6' : 'translate-x-1'}
                `}
      />
    </button>
  </div>
);

export function SettingsPage() {
  const { t } = useTranslation();
  const {
    user,
    logout,
    logoutAll,
    sessions,
    revokeSession,
    revokeOtherSessions,
    deleteAccount,
    exportData,
    fetchSessions,
    fetchProfile,
    updateProfile,
    changePassword,
  } = useAuthStore();
  const { addToast } = useUIStore();
  const { theme, isSystemPreference, setSystemPreference } = useThemeStore();
  const { language, setLanguage } = useLanguageStore();

  const [profileData, setProfileData] = useState({
    name: '',
    email: '',
    timezone: 'Asia/Jakarta',
  });

  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });

  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isPasswordModalOpen, setIsPasswordModalOpen] = useState(false);
  const [isSavingProfile, setIsSavingProfile] = useState(false);
  const [isSavingPassword, setIsSavingPassword] = useState(false);

  useEffect(() => {
    fetchProfile();
    fetchSessions();
  }, [fetchProfile, fetchSessions]);

  useEffect(() => {
    if (user) {
      setProfileData({
        name: user.name || '',
        email: user.email || '',
        timezone: user.timezone || 'Asia/Jakarta',
      });
    }
  }, [user]);

  const handleSaveProfile = async () => {
    setIsSavingProfile(true);
    const result = await updateProfile({
      name: profileData.name,
      email: profileData.email,
      timezone: profileData.timezone,
    });
    setIsSavingProfile(false);

    if (result.success) {
      addToast({ type: 'success', title: t('common.success'), message: t('toast.profileUpdated') });
    } else {
      addToast({ type: 'error', title: t('common.error'), message: result.error });
    }
  };

  const handleChangePassword = async () => {
    if (passwordData.newPassword !== passwordData.confirmPassword) {
      addToast({ type: 'error', title: t('common.error'), message: t('auth.validation.passwordMismatch') });
      return;
    }

    if (passwordData.newPassword.length < 8) {
      addToast({ type: 'error', title: t('common.error'), message: t('auth.validation.passwordMin') });
      return;
    }

    setIsSavingPassword(true);
    const result = await changePassword({
      current_password: passwordData.currentPassword,
      new_password: passwordData.newPassword,
    });
    setIsSavingPassword(false);

    if (result.success) {
      addToast({ type: 'success', title: t('common.success'), message: t('toast.passwordChanged') });
      setIsPasswordModalOpen(false);
      setPasswordData({ currentPassword: '', newPassword: '', confirmPassword: '' });
    } else {
      addToast({ type: 'error', title: t('common.error'), message: result.error });
    }
  };

  const [deletePassword, setDeletePassword] = useState('');

  const handleDeleteAccount = async () => {
    if (!deletePassword) return;

    const result = await deleteAccount(deletePassword);
    if (result.success) {
      addToast({ type: 'success', title: t('common.success'), message: 'Account permanently deleted' });
      setIsDeleteModalOpen(false);
    } else {
      addToast({ type: 'error', title: t('common.error'), message: result.error });
    }
  };

  const handleRevokeSession = async (sessionId) => {
    const result = await revokeSession(sessionId);
    if (result.success) {
      addToast({ type: 'success', message: t('toast.sessionRevoked') });
    } else {
      addToast({ type: 'error', message: result.error });
    }
  };

  const handleRevokeOtherSessions = async () => {
    const result = await revokeOtherSessions();
    if (result.success) {
      addToast({ type: 'success', message: `${result.count || 0} other sessions revoked` });
    } else {
      addToast({ type: 'error', message: result.error });
    }
  };

  const handleExportData = async () => {
    try {
      const result = await exportData();
      if (result.success) {
        const blob = new Blob([JSON.stringify(result.data, null, 2)], { type: 'application/json' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `ethos-data-${new Date().toISOString().split('T')[0]}.json`;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
        addToast({ type: 'success', title: 'Export Ready', message: 'Data download started' });
      } else {
        addToast({ type: 'error', title: 'Export Failed', message: result.error });
      }
    } catch (error) {
      addToast({ type: 'error', title: 'Export Failed', message: error.message });
    }
  };

  const handleLogoutAll = async () => {
    await logoutAll();
    addToast({ type: 'success', message: t('toast.loggedOutAll') });
  };

  return (
    <div className="space-y-6">
      <Header title={t('settings.title')} subtitle={t('settings.subtitle')} />

      <div>
        {/* Profile Section */}
        <SettingSection icon={User} title={t('settings.profile.title')} description={t('settings.profile.subtitle')}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              label={t('settings.profile.name')}
              value={profileData.name}
              onChange={(e) => setProfileData((prev) => ({ ...prev, name: e.target.value }))}
            />
            <Input
              label={t('settings.profile.email')}
              type="email"
              value={profileData.email}
              onChange={(e) => setProfileData((prev) => ({ ...prev, email: e.target.value }))}
            />
          </div>
          <div className="mt-4">
            <label className="block text-sm font-medium text-base-content mb-1.5">{t('settings.profile.timezone')}</label>
            <select
              className="w-full px-3 py-2.5 text-sm bg-base-100 border border-base-300 rounded-md text-base-content focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20"
              value={profileData.timezone}
              onChange={(e) => setProfileData((prev) => ({ ...prev, timezone: e.target.value }))}
            >
              <option value="Asia/Jakarta">Asia/Jakarta (WIB, UTC+7)</option>
              <option value="Asia/Makassar">Asia/Makassar (WITA, UTC+8)</option>
              <option value="Asia/Jayapura">Asia/Jayapura (WIT, UTC+9)</option>
              <option value="Asia/Singapore">Asia/Singapore (SGT, UTC+8)</option>
              <option value="Asia/Tokyo">Asia/Tokyo (JST, UTC+9)</option>
              <option value="Asia/Seoul">Asia/Seoul (KST, UTC+9)</option>
              <option value="America/New_York">America/New_York (EST, UTC-5)</option>
              <option value="America/Los_Angeles">America/Los_Angeles (PST, UTC-8)</option>
              <option value="Europe/London">Europe/London (GMT/BST)</option>
              <option value="UTC">UTC</option>
            </select>
            <p className="mt-1.5 text-xs text-base-content/50">{t('settings.profile.timezoneHelp')}</p>
          </div>
          <div className="mt-5 flex justify-end">
            <Button variant="primary" onClick={handleSaveProfile} loading={isSavingProfile}>
              <Save size={16} />
              {t('settings.profile.saveChanges')}
            </Button>
          </div>
        </SettingSection>

        {/* Appearance Section */}
        <SettingSection icon={Palette} title={t('settings.appearance.title')} description={t('settings.appearance.subtitle')}>
          <div className="space-y-5">
            {/* Theme Selection */}
            <div>
              <label className="block text-sm font-medium text-base-content mb-3">{t('settings.appearance.theme')}</label>
              <ThemeSelector />
            </div>

            {/* System Preference Toggle */}
            <div className="pt-2">
              <Toggle label={t('settings.appearance.systemPreference')} checked={isSystemPreference} onChange={setSystemPreference} />
              <p className="text-xs text-base-content/50 mt-2">{t('settings.appearance.systemPreferenceHelp')}</p>
            </div>

            {/* Current Theme Preview */}
            <div className="p-4 bg-base-200 rounded-lg border border-base-300">
              <div className="flex items-center gap-3">
                <div
                  className={`w-10 h-10 rounded-lg flex items-center justify-center ${
                    theme === THEMES.DARK ? 'bg-slate-700' : 'bg-white border border-base-300'
                  }`}
                >
                  {theme === THEMES.DARK ? <Moon size={20} className="text-indigo-400" /> : <Sun size={20} className="text-amber-500" />}
                </div>
                <div>
                  <p className="text-sm font-medium text-base-content">
                    {theme === THEMES.DARK ? t('settings.appearance.dark') : t('settings.appearance.light')}
                  </p>
                  <p className="text-xs text-base-content/50">
                    {isSystemPreference ? t('settings.appearance.followingSystem') : t('settings.appearance.manuallySelected')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </SettingSection>

        {/* Language Section */}
        <SettingSection icon={Globe} title={t('settings.language.title')} description={t('settings.language.subtitle')}>
          <div className="space-y-4">
            {/* Language Toggle */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-3">
                <span className="text-sm font-medium text-base-content">{t('settings.language.current')}</span>
                <Badge variant="primary">{languages.find((l) => l.code === language)?.name || 'English'}</Badge>
              </div>
              <LanguageToggleAnimated />
            </div>

            {/* Language Options */}
            <div className="p-4 bg-base-200 rounded-lg border border-base-300">
              <div className="space-y-3">
                {languages.map((lang) => (
                  <button
                    key={lang.code}
                    onClick={() => setLanguage(lang.code)}
                    className={`
                      w-full flex items-center gap-3 p-3 rounded-lg transition-all
                      ${language === lang.code ? 'bg-primary/10 border border-primary/30' : 'bg-base-100 border border-base-300 hover:bg-base-100/50'}
                    `}
                  >
                    <span className="text-2xl">{lang.flag}</span>
                    <div className="flex-1 text-left">
                      <p className={`text-sm font-medium ${language === lang.code ? 'text-primary' : 'text-base-content'}`}>{lang.name}</p>
                      <p className="text-xs text-base-content/50">{lang.code.toUpperCase()}</p>
                    </div>
                    {language === lang.code && <CheckCircle size={18} className="text-primary" />}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </SettingSection>

        {/* Install App Section */}
        <SettingSection icon={Smartphone} title={t('settings.installApp.title')} description={t('settings.installApp.subtitle')}>
          <div className="space-y-4">
            {/* PWA Status */}
            <div className="p-4 bg-base-200 rounded-lg border border-base-300">
              <div className="flex items-start gap-4">
                <div className="p-3 bg-primary/10 rounded-xl">
                  <Download size={24} className="text-primary" />
                </div>
                <div className="flex-1">
                  <h4 className="font-semibold text-base-content mb-1">{t('settings.installApp.pwaTitle')}</h4>
                  <p className="text-sm text-base-content/60 mb-3">{t('settings.installApp.pwaDesc')}</p>

                  {/* Browser-specific instructions */}
                  <div className="space-y-3">
                    <div className="text-sm">
                      <p className="font-medium text-base-content mb-2">{t('settings.installApp.howToInstall')}</p>
                      <div className="space-y-2 text-base-content/70">
                        <div className="flex items-start gap-2">
                          <span className="font-medium text-primary shrink-0">Chrome/Edge:</span>
                          <span>{t('settings.installApp.chrome')}</span>
                        </div>
                        <div className="flex items-start gap-2">
                          <span className="font-medium text-primary shrink-0">Safari (Mac):</span>
                          <span>{t('settings.installApp.safariMac')}</span>
                        </div>
                        <div className="flex items-start gap-2">
                          <span className="font-medium text-primary shrink-0">Safari (iOS):</span>
                          <span>{t('settings.installApp.safariIOS')}</span>
                        </div>
                        <div className="flex items-start gap-2">
                          <span className="font-medium text-primary shrink-0">Firefox:</span>
                          <span>{t('settings.installApp.firefox')}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* App Info */}
            <div className="flex items-center justify-between p-3 bg-base-50 border border-base-200 rounded-lg">
              <div className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-primary flex items-center justify-center">
                  <Shield size={18} className="text-primary-content" />
                </div>
                <div>
                  <p className="text-sm font-medium text-base-content">Ethos - Habit Tracker</p>
                  <p className="text-xs text-base-content/50">Version 1.0.0</p>
                </div>
              </div>
              <Badge variant="success" size="sm">
                <CheckCircle size={12} className="mr-1" />
                {t('settings.installApp.readyToInstall')}
              </Badge>
            </div>

            {/* Note about development */}
            <p className="text-xs text-base-content/40 italic">{t('settings.installApp.devNote')}</p>
          </div>
        </SettingSection>

        {/* Active Sessions Section */}
        <SettingSection icon={Monitor} title={t('settings.sessions.title')} description={t('settings.sessions.subtitle')} collapsible defaultOpen={true}>
          <div className="space-y-3">
            <div className="flex justify-end gap-2 mb-2">
              <Button size="sm" variant="outline" onClick={() => fetchSessions()}>
                Refresh
              </Button>
              <Button size="sm" variant="secondary" onClick={handleRevokeOtherSessions}>
                Revoke Others
              </Button>
            </div>
            {sessions.length > 0 ? (
              sessions.map((session) => (
                <div key={session.session_id} className="flex items-center justify-between p-3 bg-base-50 border border-base-200 rounded-lg">
                  <div>
                    <p className="text-sm font-medium text-base-content">{session.user_agent || 'Unknown Device'}</p>
                    <p className="text-xs text-base-content/50 mt-0.5">
                      {session.client_ip} • {session.created_at ? format(new Date(session.created_at), 'PP p') : 'Unknown Date'}
                    </p>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge variant={session.is_current ? 'primary' : 'success'} size="xs">
                      {session.is_current ? t('settings.sessions.current') : t('settings.sessions.active')}
                    </Badge>
                    {!session.is_current && (
                      <button
                        onClick={() => handleRevokeSession(session.session_id)}
                        className="btn btn-ghost btn-xs btn-square text-error hover:bg-error/10"
                        title={t('settings.sessions.revoke')}
                      >
                        <Trash2 size={14} />
                      </button>
                    )}
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-base-content/50">{t('settings.sessions.noSessions')}</p>
            )}
            <div className="pt-3 border-t border-base-100 mt-2">
              <Button variant="danger" size="sm" onClick={handleLogoutAll}>
                <LogOut size={14} />
                {t('settings.sessions.logoutAll')}
              </Button>
            </div>
          </div>
        </SettingSection>

        {/* Security Section (Updated with Export) */}
        <SettingSection icon={Shield} title={t('settings.security.title')} description={t('settings.security.subtitle')}>
          <div className="flex flex-wrap gap-3">
            <Button variant="secondary" onClick={() => setIsPasswordModalOpen(true)}>
              <Key size={16} />
              {t('settings.security.changePassword')}
            </Button>
            <Button variant="outline" onClick={handleExportData}>
              <Download size={16} />
              Export My Data
            </Button>
          </div>
        </SettingSection>

        {/* Danger Zone */}
        <div className="bg-error/5 border border-error/20 rounded-lg p-6">
          <div className="flex items-start gap-4">
            <div className="p-2.5 bg-error/10 rounded-lg shrink-0">
              <Trash2 size={20} className="text-error" />
            </div>
            <div className="flex-1">
              <h3 className="text-base font-semibold text-error mb-1">{t('settings.danger.title')}</h3>
              <p className="text-sm text-error/70 mb-5">{t('settings.danger.subtitle')}</p>
              <Button variant="danger" onClick={() => setIsDeleteModalOpen(true)}>
                <Trash2 size={16} />
                {t('settings.danger.deleteAccount')}
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Delete Account Modal with Password Confirmation */}
      <Modal isOpen={isDeleteModalOpen} onClose={() => setIsDeleteModalOpen(false)} title="Delete Account">
        <div className="space-y-4">
          <div className="p-4 bg-error/10 text-error rounded-lg text-sm">
            Warning: This action is irreversible. All your habits, logs, and data will be permanently deleted.
          </div>
          <p className="text-sm text-base-content/80">Please enter your password to confirm deletion.</p>
          <Input type="password" placeholder="Enter your password" value={deletePassword} onChange={(e) => setDeletePassword(e.target.value)} />
          <div className="flex justify-end gap-3 pt-4">
            <Button onClick={() => setIsDeleteModalOpen(false)} variant="secondary">
              Cancel
            </Button>
            <Button onClick={handleDeleteAccount} variant="destructive" disabled={!deletePassword}>
              Confirm Delete
            </Button>
          </div>
        </div>
      </Modal>

      {/* Change Password Modal */}
      <Modal isOpen={isPasswordModalOpen} onClose={() => setIsPasswordModalOpen(false)} title={t('settings.security.changePassword')}>
        <div className="space-y-4">
          <Input
            label={t('settings.security.currentPassword')}
            type="password"
            value={passwordData.currentPassword}
            onChange={(e) => setPasswordData((prev) => ({ ...prev, currentPassword: e.target.value }))}
            placeholder={t('settings.security.currentPassword')}
          />
          <Input
            label={t('settings.security.newPassword')}
            type="password"
            value={passwordData.newPassword}
            onChange={(e) => setPasswordData((prev) => ({ ...prev, newPassword: e.target.value }))}
            placeholder={t('settings.security.passwordHint')}
          />
          <Input
            label={t('settings.security.confirmPassword')}
            type="password"
            value={passwordData.confirmPassword}
            onChange={(e) => setPasswordData((prev) => ({ ...prev, confirmPassword: e.target.value }))}
            placeholder={t('settings.security.confirmPassword')}
          />
        </div>
        <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-base-200">
          <Button variant="ghost" onClick={() => setIsPasswordModalOpen(false)}>
            {t('common.cancel')}
          </Button>
          <Button variant="primary" onClick={handleChangePassword} loading={isSavingPassword}>
            {t('settings.security.changePassword')}
          </Button>
        </div>
      </Modal>

      <ConfirmModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={handleDeleteAccount}
        title={t('settings.danger.deleteAccount')}
        message={t('settings.danger.deleteConfirm')}
        confirmText={t('settings.danger.deleteAccount')}
        variant="danger"
      />
    </div>
  );
}
