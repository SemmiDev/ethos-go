import axios from 'axios';
import * as SecureStore from 'expo-secure-store';

// API base URL - change this to your backend URL
const API_BASE_URL = 'https://1a1e6ff05e3e.ngrok-free.app/api'; // Update with your server IP

// Create axios instance with base configuration
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Get access token from secure storage
const getAccessToken = async () => {
  try {
    return await SecureStore.getItemAsync('accessToken');
  } catch {
    return null;
  }
};

// Request interceptor to add auth token
apiClient.interceptors.request.use(
  async (config) => {
    const token = await getAccessToken();
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

    // Check if this is an auth endpoint
    const isAuthEndpoint = AUTH_ENDPOINTS.some((endpoint) => requestUrl.includes(endpoint));

    // If 401 and not already retrying and NOT an auth endpoint
    if (error.response?.status === 401 && !originalRequest._retry && !isAuthEndpoint) {
      originalRequest._retry = true;
      // Token expired, user will be logged out by the store
      await SecureStore.deleteItemAsync('accessToken');
    }

    return Promise.reject(error);
  }
);

// Helper to update base URL (for dev/prod switching)
export const setApiBaseUrl = (url) => {
  apiClient.defaults.baseURL = url;
};

export default apiClient;
