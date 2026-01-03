import axios from 'axios';
import { useAuthStore } from '../stores/authStore';

// Create axios instance with base configuration
const apiClient = axios.create({
  baseURL: '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().accessToken;
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Auth endpoints that should NOT trigger auto-logout on 401
const AUTH_ENDPOINTS = ['/auth/login', '/auth/register', '/auth/verify-email', '/auth/resend-verification', '/auth/forgot-password', '/auth/reset-password'];

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const requestUrl = originalRequest?.url || '';

    // Check if this is an auth endpoint that should handle its own 401 errors
    const isAuthEndpoint = AUTH_ENDPOINTS.some((endpoint) => requestUrl.includes(endpoint));

    // If 401 and not already retrying and NOT an auth endpoint, logout
    if (error.response?.status === 401 && !originalRequest._retry && !isAuthEndpoint) {
      originalRequest._retry = true;

      // Logout on 401 for protected routes
      useAuthStore.getState().logout();
      window.location.href = '/login';
    }

    return Promise.reject(error);
  }
);

export default apiClient;
