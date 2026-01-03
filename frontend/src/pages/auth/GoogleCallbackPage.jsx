import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { authAPI } from '../../api/auth';
import { useAuthStore } from '../../stores/authStore';
import { useUIStore } from '../../stores/uiStore';

const GoogleCallbackPage = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);
  const addToast = useUIStore((state) => state.addToast);

  useEffect(() => {
    const code = searchParams.get('code');
    if (!code) {
      addToast({ type: 'error', title: 'Error', message: 'No authorization code found' });
      navigate('/login');
      return;
    }

    const handleCallback = async () => {
      try {
        const response = await authAPI.googleCallback(code);
        if (response.success) {
          const { access_token, refresh_token, user_id, session_id, expires_at } = response.data;
          setAuth(access_token, refresh_token, user_id, session_id, expires_at);
          addToast({ type: 'success', title: 'Success', message: 'Successfully logged in with Google' });
          navigate('/dashboard');
        }
      } catch (error) {
        console.error(error);
        const msg = error.response?.data?.message || 'Google login failed';
        addToast({ type: 'error', title: 'Login Failed', message: msg });
        navigate('/login');
      }
    };

    handleCallback();
  }, [searchParams, navigate, setAuth, addToast]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="text-center">
        <svg className="animate-spin h-10 w-10 text-primary mx-auto mb-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          ></path>
        </svg>
        <h2 className="text-xl font-semibold mb-2">Authenticating...</h2>
        <p className="text-gray-500">Please wait while we log you in with Google.</p>
      </div>
    </div>
  );
};

export default GoogleCallbackPage;
