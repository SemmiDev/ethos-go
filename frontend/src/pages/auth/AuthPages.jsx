import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Mail, Lock, User, ArrowRight, Shield, CheckCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { LanguageToggle } from '../../components/ui/LanguageToggle';
import { useAuthStore } from '../../stores/authStore';
import { useUIStore } from '../../stores/uiStore';
import { authAPI } from '../../api/auth';

export const AuthLayout = ({ children, title, subtitle, features = [] }) => (
  <div className="flex min-h-screen">
    {/* Left Side - Branding */}
    <div className="hidden lg:flex flex-col w-1/2 bg-primary p-12 xl:p-16 relative overflow-hidden">
      {/* Subtle Background Pattern */}
      <div className="absolute inset-0 opacity-10">
        <div className="absolute top-0 right-0 w-96 h-96 bg-white/20 rounded-full blur-3xl" />
        <div className="absolute bottom-0 left-0 w-96 h-96 bg-white/10 rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 flex-1 flex flex-col justify-center">
        <div className="flex items-center gap-3 mb-10">
          <img src="/logo.jpg" alt="Ethos Logo" className="w-12 h-12 rounded-xl object-cover" />
          <div>
            <h2 className="text-xl font-semibold text-white">Ethos</h2>
            <p className="text-xs text-white/60 uppercase tracking-wider">Habit Tracker</p>
          </div>
        </div>

        <h1 className="text-4xl xl:text-5xl font-semibold text-white leading-tight mb-4 tracking-tight">{title}</h1>
        <p className="text-lg text-white/70 max-w-md leading-relaxed mb-10">{subtitle}</p>

        {features.length > 0 && (
          <ul className="space-y-4">
            {features.map((feature, index) => (
              <li key={index} className="flex items-center gap-3 text-white/80">
                <CheckCircle size={18} className="text-accent shrink-0" />
                <span className="text-sm">{feature}</span>
              </li>
            ))}
          </ul>
        )}
      </div>

      <p className="relative z-10 text-sm text-white/40">© {new Date().getFullYear()} Ethos Inc. All rights reserved.</p>
    </div>

    {/* Right Side - Form */}
    <div className="flex-1 flex flex-col px-6 py-12 sm:px-12 lg:px-16 xl:px-24 bg-base-100">
      {/* Language Toggle */}
      <div className="flex justify-end mb-4">
        <LanguageToggle size="sm" />
      </div>

      <div className="flex-1 flex flex-col justify-center">
        <div className="w-full max-w-sm mx-auto">
          {/* Mobile Logo */}
          <div className="lg:hidden flex flex-col items-center gap-3 mb-10">
            <img src="/logo.jpg" alt="Ethos Logo" className="w-12 h-12 rounded-xl object-cover" />
            <div className="text-center">
              <h1 className="text-xl font-semibold text-base-content">Ethos</h1>
              <p className="text-xs text-base-content/50 uppercase tracking-wider">Habit Tracker</p>
            </div>
          </div>

          {children}
        </div>
      </div>
    </div>
  </div>
);

export function LoginPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { login, isLoading } = useAuthStore();
  const { addToast } = useUIStore();
  const [formData, setFormData] = useState({ email: '', password: '' });
  const [errors, setErrors] = useState({});

  const validate = () => {
    const newErrors = {};
    if (!formData.email) {
      newErrors.email = t('auth.validation.emailRequired');
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = t('auth.validation.emailInvalid');
    }
    if (!formData.password) {
      newErrors.password = t('auth.validation.passwordRequired');
    }
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validate()) return;

    const result = await login(formData);
    if (result.success) {
      addToast({ type: 'success', title: t('auth.welcomeBack'), message: t('auth.loginSuccess') });
      navigate('/dashboard');
    } else {
      // Check if the error is due to email not being verified
      const isEmailNotVerified = result.errorCode === 'AUTH_UNAUTHORIZED' && result.error?.toLowerCase().includes('verify');

      if (isEmailNotVerified) {
        addToast({
          type: 'warning',
          title: t('auth.emailNotVerified'),
          message: t('auth.verifyEmailPrompt'),
        });
        // Small delay to allow user to see the toast before redirect
        setTimeout(() => {
          navigate(`/verify-email?email=${encodeURIComponent(formData.email)}`);
        }, 1500);
      } else {
        addToast({ type: 'error', title: t('auth.loginFailed'), message: result.error });
      }
    }
  };

  const handleChange = (field) => (e) => {
    setFormData({ ...formData, [field]: e.target.value });
    if (errors[field]) setErrors({ ...errors, [field]: null });
  };

  const handleGoogleLogin = async () => {
    try {
      const { data } = await authAPI.getGoogleLoginURL();
      window.location.href = data.url;
    } catch (error) {
      addToast({ type: 'error', title: t('common.error'), message: t('auth.loginFailed') });
    }
  };

  return (
    <AuthLayout
      title={t('landing.hero.title') + ' ' + t('landing.hero.titleHighlight')}
      subtitle={t('landing.hero.subtitle')}
      features={[t('auth.features.unlimited'), t('auth.features.analytics'), t('auth.features.secure')]}
    >
      <div className="mb-8">
        <h2 className="text-2xl font-semibold text-base-content tracking-tight">{t('auth.login.title')}</h2>
        <p className="text-base-content/60 mt-1.5">{t('auth.login.subtitle')}</p>
      </div>

      <button type="button" className="w-full h-11 mb-6 flex items-center justify-center gap-2 btn-banking btn-banking-secondary" onClick={handleGoogleLogin}>
        <svg className="w-5 h-5" viewBox="0 0 24 24">
          <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
          <path
            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
            fill="#34A853"
          />
          <path
            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.26.81-.58z"
            fill="#FBBC05"
          />
          <path
            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
            fill="#EA4335"
          />
        </svg>
        {t('auth.google.signIn')}
      </button>

      <div className="relative my-6">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-base-300"></div>
        </div>
        <div className="relative flex justify-center text-sm">
          <span className="px-2 bg-base-100 text-base-content/60">{t('auth.google.continueWith')}</span>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        <Input
          label={t('auth.login.email')}
          type="email"
          placeholder="name@example.com"
          icon={Mail}
          value={formData.email}
          onChange={handleChange('email')}
          error={errors.email}
        />

        <div className="space-y-1">
          <Input
            label={t('auth.login.password')}
            type="password"
            placeholder="••••••••"
            icon={Lock}
            value={formData.password}
            onChange={handleChange('password')}
            error={errors.password}
          />
          <div className="flex justify-between items-center pt-1">
            <label className="flex items-center gap-2 cursor-pointer">
              <input type="checkbox" className="w-4 h-4 rounded border-base-300 text-primary focus:ring-primary/20" />
              <span className="text-sm text-base-content/60">{t('auth.rememberMe')}</span>
            </label>
            <Link to="/forgot-password" className="text-sm text-primary hover:text-primary/80 font-medium transition-colors">
              {t('auth.login.forgotPassword')}
            </Link>
          </div>
        </div>

        <Button type="submit" className="w-full h-11" isLoading={isLoading}>
          {t('auth.login.submitButton')}
          <ArrowRight size={16} />
        </Button>
      </form>

      <p className="text-center text-sm text-base-content/60 mt-8">
        {t('auth.login.noAccount')}{' '}
        <Link to="/register" className="text-primary hover:text-primary/80 font-medium transition-colors">
          {t('auth.login.signUp')}
        </Link>
      </p>
    </AuthLayout>
  );
}

export function RegisterPage() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const { register } = useAuthStore();
  const { addToast } = useUIStore();
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState({ name: '', email: '', password: '', confirmPassword: '' });
  const [errors, setErrors] = useState({});

  const validate = () => {
    const newErrors = {};
    if (!formData.name) newErrors.name = t('auth.validation.nameRequired');
    else if (formData.name.length < 2) newErrors.name = t('auth.validation.nameMin');
    if (!formData.email) newErrors.email = t('auth.validation.emailRequired');
    else if (!/\S+@\S+\.\S+/.test(formData.email)) newErrors.email = t('auth.validation.emailInvalid');
    if (!formData.password) newErrors.password = t('auth.validation.passwordRequired');
    else if (formData.password.length < 8) newErrors.password = t('auth.validation.passwordMin');
    if (formData.password !== formData.confirmPassword) newErrors.confirmPassword = t('auth.validation.passwordMismatch');
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!validate()) return;

    setIsLoading(true);
    const result = await register({ name: formData.name, email: formData.email, password: formData.password });
    setIsLoading(false);

    if (result.success) {
      addToast({ type: 'success', title: t('auth.accountCreated'), message: t('auth.verifyEmailMessage') });
      navigate(`/verify-email?email=${encodeURIComponent(formData.email)}`);
    } else {
      addToast({ type: 'error', title: t('auth.registrationFailed'), message: result.error });
    }
  };

  const handleChange = (field) => (e) => {
    setFormData({ ...formData, [field]: e.target.value });
    if (errors[field]) setErrors({ ...errors, [field]: null });
  };

  return (
    <AuthLayout
      title={t('auth.register.title')}
      subtitle={t('auth.register.subtitle')}
      features={[t('auth.features.free'), t('auth.features.noCreditCard'), t('auth.features.cancelAnytime')]}
    >
      <div className="mb-8">
        <h2 className="text-2xl font-semibold text-base-content tracking-tight">{t('auth.register.title')}</h2>
        <p className="text-base-content/60 mt-1.5">{t('auth.register.subtitle')}</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label={t('auth.register.name')}
          type="text"
          placeholder="John Doe"
          icon={User}
          value={formData.name}
          onChange={handleChange('name')}
          error={errors.name}
        />

        <Input
          label={t('auth.register.email')}
          type="email"
          placeholder="name@example.com"
          icon={Mail}
          value={formData.email}
          onChange={handleChange('email')}
          error={errors.email}
        />

        <Input
          label={t('auth.register.password')}
          type="password"
          placeholder="••••••••"
          icon={Lock}
          value={formData.password}
          onChange={handleChange('password')}
          error={errors.password}
          helperText={t('auth.register.passwordHint')}
        />

        <Input
          label={t('auth.register.confirmPassword')}
          type="password"
          placeholder="••••••••"
          icon={Lock}
          value={formData.confirmPassword}
          onChange={handleChange('confirmPassword')}
          error={errors.confirmPassword}
        />

        <div className="flex items-start gap-3 py-2">
          <input type="checkbox" id="terms" className="w-4 h-4 mt-0.5 rounded border-base-300 text-primary focus:ring-primary/20" />
          <label htmlFor="terms" className="text-sm text-base-content/60 cursor-pointer leading-relaxed">
            {t('auth.terms')}{' '}
            <a href="#" className="text-primary hover:text-primary/80 font-medium">
              {t('auth.termsOfService')}
            </a>{' '}
            {t('auth.and')}{' '}
            <a href="#" className="text-primary hover:text-primary/80 font-medium">
              {t('auth.privacyPolicy')}
            </a>
          </label>
        </div>

        <Button type="submit" className="w-full h-11" isLoading={isLoading}>
          {t('auth.register.submitButton')}
          <ArrowRight size={16} />
        </Button>
      </form>

      <p className="text-center text-sm text-base-content/60 mt-8">
        {t('auth.register.hasAccount')}{' '}
        <Link to="/login" className="text-primary hover:text-primary/80 font-medium transition-colors">
          {t('auth.register.signIn')}
        </Link>
      </p>
    </AuthLayout>
  );
}
