import { useState } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { Mail, Lock, Key, ArrowRight, CheckCircle, RefreshCw } from 'lucide-react';
import { Button } from '../../components/ui/Button';
import { Input } from '../../components/ui/Input';
import { useUIStore } from '../../stores/uiStore';
import { authAPI } from '../../api/auth';
import { AuthLayout } from './AuthPages';

export function VerifyEmailPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { addToast } = useUIStore();
  const [email, setEmail] = useState(searchParams.get('email') || '');
  const [code, setCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isResendLoading, setIsResendLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!email || !code) return;

    setIsLoading(true);
    try {
      await authAPI.verifyEmail({ email, code: code });
      addToast({ type: 'success', title: 'Verified!', message: 'Email verified successfully. Please login.' });
      navigate('/login');
    } catch (error) {
      addToast({ type: 'error', title: 'Verification Failed', message: error.response?.data?.message || 'Invalid code' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleResend = async () => {
    if (!email) {
      addToast({ type: 'error', title: 'Error', message: 'Please enter your email first' });
      return;
    }
    setIsResendLoading(true);
    try {
      await authAPI.resendVerification({ email });
      addToast({ type: 'success', title: 'Sent!', message: 'Verification code resent to your email.' });
    } catch (error) {
      addToast({ type: 'error', title: 'Failed', message: error.response?.data?.message || 'Could not resend code' });
    } finally {
      setIsResendLoading(false);
    }
  };

  return (
    <AuthLayout
      title="Verify Your Email"
      subtitle="Enter the verification code sent to your email address to activate your account."
      features={['Secure account activation', 'Instant access', 'Account protection']}
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Email Address"
          type="email"
          placeholder="name@example.com"
          icon={Mail}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          disabled={!!searchParams.get('email')}
        />
        <Input label="Verification Code" type="text" placeholder="123456" icon={Key} value={code} onChange={(e) => setCode(e.target.value)} maxLength={6} />
        <Button type="submit" className="w-full h-11" isLoading={isLoading}>
          Verify Account <CheckCircle size={16} />
        </Button>
      </form>
      <div className="mt-6 text-center">
        <button
          type="button"
          onClick={handleResend}
          className="text-sm text-primary hover:text-primary/80 font-medium flex items-center justify-center gap-2 mx-auto transition-colors"
          disabled={isResendLoading}
        >
          {isResendLoading ? <RefreshCw size={14} className="animate-spin" /> : <RefreshCw size={14} />}
          Resend Verification Code
        </button>
      </div>
      <p className="text-center text-sm text-base-content/60 mt-8">
        Back to{' '}
        <Link to="/login" className="text-primary hover:text-primary/80 font-medium transition-colors">
          Login
        </Link>
      </p>
    </AuthLayout>
  );
}

export function ForgotPasswordPage() {
  const navigate = useNavigate();
  const { addToast } = useUIStore();
  const [email, setEmail] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSent, setIsSent] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!email) {
      addToast({ type: 'error', title: 'Required', message: 'Email address is required' });
      return;
    }
    setIsLoading(true);
    try {
      await authAPI.forgotPassword({ email });
      setIsSent(true);
      addToast({ type: 'success', title: 'Email Sent', message: 'If an account exists, we sent a reset code.' });
    } catch (error) {
      addToast({ type: 'error', title: 'Error', message: error.message });
    } finally {
      setIsLoading(false);
    }
  };

  if (isSent) {
    return (
      <AuthLayout title="Check your email" subtitle={`We have sent password reset instructions to ${email}`}>
        <div className="text-center">
          <div className="w-16 h-16 bg-success/10 text-success rounded-full flex items-center justify-center mx-auto mb-6">
            <Mail size={32} />
          </div>
          <p className="text-base-content/70 mb-8 leading-relaxed">Did not receive the email? Check your spam folder or try another email address.</p>
          <Link to={`/reset-password?email=${email}`}>
            <Button className="w-full h-11">Enter Reset Code</Button>
          </Link>
          <button onClick={() => setIsSent(false)} className="mt-6 text-sm text-base-content/60 hover:text-primary transition-colors">
            Try a different email
          </button>
        </div>
      </AuthLayout>
    );
  }

  return (
    <AuthLayout title="Reset Password" subtitle="Enter your email address and we'll send you instructions to reset your password.">
      <form onSubmit={handleSubmit} className="space-y-5">
        <Input label="Email Address" type="email" placeholder="name@example.com" icon={Mail} value={email} onChange={(e) => setEmail(e.target.value)} />
        <Button type="submit" className="w-full h-11" isLoading={isLoading}>
          Send Reset Instructions <ArrowRight size={16} />
        </Button>
      </form>
      <p className="text-center text-sm text-base-content/60 mt-8">
        Remember your password?{' '}
        <Link to="/login" className="text-primary hover:text-primary/80 font-medium transition-colors">
          Sign in
        </Link>
      </p>
    </AuthLayout>
  );
}

export function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { addToast } = useUIStore();
  const [email, setEmail] = useState(searchParams.get('email') || '');
  const [code, setCode] = useState(searchParams.get('code') || '');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!email || !code || !password) {
      addToast({ type: 'error', title: 'Required', message: 'All fields are required' });
      return;
    }

    setIsLoading(true);
    try {
      await authAPI.resetPassword({ email, code, new_password: password });
      addToast({ type: 'success', title: 'Success', message: 'Password has been reset. Please login.' });
      navigate('/login');
    } catch (error) {
      addToast({ type: 'error', title: 'Failed', message: error.response?.data?.message || 'Reset failed' });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <AuthLayout title="Set New Password" subtitle="Your new password must be different to previously used passwords.">
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input label="Email Address" type="email" value={email} onChange={(e) => setEmail(e.target.value)} disabled={!!searchParams.get('email')} icon={Mail} />
        <Input label="Reset Code" type="text" placeholder="123456" icon={Key} value={code} onChange={(e) => setCode(e.target.value)} />
        <Input
          label="New Password"
          type="password"
          placeholder="••••••••"
          icon={Lock}
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          helperText="Must be at least 8 characters"
        />
        <Button type="submit" className="w-full h-11" isLoading={isLoading}>
          Reset Password <CheckCircle size={16} />
        </Button>
      </form>
      <p className="text-center text-sm text-base-content/60 mt-8">
        Back to{' '}
        <Link to="/login" className="text-primary hover:text-primary/80 font-medium transition-colors">
          Login
        </Link>
      </p>
    </AuthLayout>
  );
}
