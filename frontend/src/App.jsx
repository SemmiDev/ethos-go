import { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { MainLayout, AuthLayout } from './components/layout';
import { LoginPage, RegisterPage, VerifyEmailPage, ForgotPasswordPage, ResetPasswordPage, GoogleCallbackPage } from './pages/auth';
// ... (imports)

// ...

{
  /* Auth routes */
}
<Route element={<AuthLayout />}>
  <Route path="/login" element={<LoginPage />} />
  <Route path="/register" element={<RegisterPage />} />
  <Route path="/verify-email" element={<VerifyEmailPage />} />
  <Route path="/forgot-password" element={<ForgotPasswordPage />} />
  <Route path="/reset-password" element={<ResetPasswordPage />} />
</Route>;

{
  /* Google Callback - Standalone */
}
<Route path="/auth/google/callback" element={<GoogleCallbackPage />} />;

{
  /* Protected routes */
}
import { LandingPage } from './pages/landing';
import { DashboardPage } from './pages/dashboard';
import { HabitsListPage, HabitDetailPage } from './pages/habits';
import { AnalyticsPage } from './pages/analytics';
import { SettingsPage } from './pages/settings';
import { NotificationsPage } from './pages/notifications/NotificationsPage';
import { HelpPage } from './pages/help';
import { useThemeStore } from './stores/themeStore';
import { useAuthStore } from './stores/authStore';

// Landing page wrapper - shows landing or redirects to dashboard if authenticated
function LandingRoute() {
  const { isAuthenticated } = useAuthStore();

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return <LandingPage />;
}

function App() {
  const initializeTheme = useThemeStore((state) => state.initializeTheme);

  useEffect(() => {
    // Add no-transition class to prevent flash on initial load
    document.documentElement.classList.add('no-transition');

    // Initialize theme from persisted state or system preference
    initializeTheme();

    // Remove no-transition class after a brief delay
    const timer = setTimeout(() => {
      document.documentElement.classList.remove('no-transition');
    }, 100);

    return () => clearTimeout(timer);
  }, [initializeTheme]);

  return (
    <BrowserRouter>
      <Routes>
        {/* Landing page - public */}
        <Route path="/" element={<LandingRoute />} />

        {/* Auth routes */}
        <Route element={<AuthLayout />}>
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/verify-email" element={<VerifyEmailPage />} />
          <Route path="/forgot-password" element={<ForgotPasswordPage />} />
          <Route path="/reset-password" element={<ResetPasswordPage />} />
        </Route>

        {/* Protected routes */}
        <Route element={<MainLayout />}>
          <Route path="/dashboard" element={<DashboardPage />} />
          <Route path="/habits" element={<HabitsListPage />} />
          <Route path="/habits/:habitId" element={<HabitDetailPage />} />
          <Route path="/analytics" element={<AnalyticsPage />} />
          <Route path="/settings" element={<SettingsPage />} />
          <Route path="/notifications" element={<NotificationsPage />} />
          <Route path="/help" element={<HelpPage />} />
        </Route>

        {/* Fallback - redirect to landing */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
